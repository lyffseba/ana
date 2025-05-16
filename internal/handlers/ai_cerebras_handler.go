// Package handlers provides HTTP request handlers
package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sebae/ana/internal/ai"
)

// CerebrasAIRequest represents an incoming request to the Cerebras AI assistant
type CerebrasAIRequest struct {
	Query     string                `form:"query" binding:"required"`
	ModelType string                `form:"model_type" binding:"required"`
	Image     *multipart.FileHeader `form:"image"`
}

// CerebrasAIResponse represents the response from the Cerebras AI assistant
type CerebrasAIResponse struct {
	Response string `json:"response"`
	HasImage bool   `json:"has_image,omitempty"`
}

// Global Cerebras AI client instance
var cerebrasClient = ai.NewCerebrasClient()

// GetCerebrasAIAssistance handles requests to the Cerebras AI assistant
// Supports both text models (QWen-3B-32B) and vision models (QWen-2.5-Vision)
func GetCerebrasAIAssistance(c *gin.Context) {
	// Use form binding for multipart form data
	var request CerebrasAIRequest
	if err := c.ShouldBind(&request); err != nil {
		log.Printf("Error binding request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error en los datos enviados. Verifica que has incluido una consulta válida."})
		return
	}
	
	// Process the query
	query := request.Query
	modelType := request.ModelType
	
	if modelType == "" {
		modelType = "qwen-text" // Default to text model
	}
	
	// Log request info
	log.Printf("Processing AI request: model=%s, query_length=%d", modelType, len(query))
	
	// Check for /no_think command
	isNoThink := false
	if strings.HasPrefix(strings.ToLower(query), "/no_think") {
		isNoThink = true
		// Remove the command from the query
		query = strings.TrimSpace(query[9:])
		log.Printf("No-think mode detected, processing query: %s", query)
	}
	
	// Process image for vision model if present
	var imageBase64 string
	hasImage := false
	
	if modelType == "qwen-vision" && request.Image != nil {
		file, err := request.Image.Open()
		if err != nil {
			log.Printf("Error opening uploaded image: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo procesar la imagen. Intenta con otro formato o una imagen más pequeña."})
			return
		}
		defer file.Close()
		
		// Read image data
		imageData, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Error reading image data: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer los datos de la imagen. Intenta con otra imagen."})
			return
		}
		
		// Check image size
		if len(imageData) > 5*1024*1024 { // 5MB limit
			log.Printf("Image too large: %d bytes", len(imageData))
			c.JSON(http.StatusBadRequest, gin.H{"error": "La imagen es demasiado grande. Por favor utiliza una imagen menor a 5MB."})
			return
		}
		
		// Convert to base64
		imageBase64 = base64.StdEncoding.EncodeToString(imageData)
		hasImage = true
		log.Printf("Image processed successfully, size: %d bytes", len(imageData))
	} else if modelType == "qwen-vision" && request.Image == nil {
		// Vision model selected but no image provided
		log.Printf("Vision model selected but no image provided")
		// We'll continue with text-only query, but log this situation
	}

	// System message for context with enhanced architectural domain knowledge
	noThinkInstructions := ""
	responseStyleInstructions := "Responde siempre en español con terminología técnica precisa."
	technicalInfoInstructions := "Cuando proporciones información técnica, incluye referencias a códigos específicos, ejemplos prácticos y consideraciones para el contexto colombiano."
	
	// Adjust instructions for no_think mode
	if isNoThink {
		noThinkInstructions = "Estás en modo respuesta directa. Responde con la información precisa sin explicaciones adicionales, usando el mínimo de palabras posible."
		responseStyleInstructions = "Responde solo con datos concretos sin introducción ni explicación."
		technicalInfoInstructions = ""
	}
	
	// Adjust instructions for vision model
	visionSpecificInstructions := ""
	if modelType == "qwen-vision" {
		if hasImage {
			visionSpecificInstructions = "Estás analizando una imagen arquitectónica. "
			if !isNoThink {
				visionSpecificInstructions += "Proporciona un análisis detallado que incluya: estilo arquitectónico, elementos estructurales prominentes, aspectos de diseño notables, posibles problemas o consideraciones técnicas, y recomendaciones basadas en el código colombiano de construcción cuando sea relevante."
			}
		} else {
			visionSpecificInstructions = "Aunque tienes capacidad de análisis visual, no se proporcionó una imagen con esta consulta. "
		}
	}
	
	// Construct the system prompt
	systemPrompt := fmt.Sprintf(
		"Eres un asistente especializado en arquitectura para la plataforma ana.world de gestión de proyectos arquitectónicos. %s%sTu conocimiento incluye: 1) Normativas colombianas: NSR-10 (Norma Sismo Resistente), POT de Bogotá, Decreto 1077 de 2015, normas urbanísticas locales; 2) Diseño arquitectónico: metodología BIM, diseño paramétrico, estilos arquitectónicos latinoamericanos, soluciones para clima tropical; 3) Gestión de proyectos: metodologías PMI/PRINCE2 adaptadas a construcción, control de cronogramas, gestión de contratistas, licencias de construcción; 4) Materiales sostenibles: guadua, tierra compactada, sistemas pasivos de climatización, certificación LEED/EDGE para Colombia; 5) Presupuestos: estimación de costos por m², control de presupuestos, análisis de precios unitarios (APU). %s Si te preguntan en inglés, comprende la consulta pero responde en español. %s",
		noThinkInstructions,
		visionSpecificInstructions,
		responseStyleInstructions,
		technicalInfoInstructions,
	)
	
	// Create system context
	systemContext := []ai.Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
	}

	// Call the appropriate Cerebras client method based on model type
	var response string
	var err error
	
	if modelType == "qwen-vision" && hasImage {
		// Process vision request with image
		response, err = cerebrasClient.GenerateVisionResponse(query, imageBase64, systemContext)
		if err != nil {
			log.Printf("Error getting vision response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en el procesamiento de la imagen. Intenta con otra consulta o imagen."})
			return
		}
	} else {
		// Use text model (for both qwen-text and qwen-vision without image)
		modelName := "cerebras/QWen-3B-32B"
		if modelType == "qwen-vision" {
			log.Printf("Vision model selected but using text model as no image was provided")
		}
		
		response, err = cerebrasClient.GenerateTextResponse(query, modelName, systemContext)
		if err != nil {
			log.Printf("Error getting text response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en el procesamiento de la consulta. Intenta reformularla."})
			return
		}
	}
	
	// Process response for no_think mode
	if isNoThink && len(response) > 0 {
		// For no_think mode, extract the most concise answer
		lines := strings.Split(response, "\n")
		shortestRelevantLine := response // Default to full response
		
		// Find the shortest line that contains relevant information
		minLength := 1000
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if len(trimmed) > 10 && len(trimmed) < minLength {
				// Check if line contains technical information (numbers, units, etc.)
				if strings.ContainsAny(trimmed, "0123456789") || 
				   strings.Contains(trimmed, "kg") ||
				   strings.Contains(trimmed, "m²") ||
				   strings.Contains(trimmed, "mm") {
					shortestRelevantLine = trimmed
					minLength = len(trimmed)
				}
			}
		}
		response = shortestRelevantLine
	}

	// Return the response with additional metadata
	c.JSON(http.StatusOK, CerebrasAIResponse{
		Response: response,
		HasImage: hasImage,
	})
}

