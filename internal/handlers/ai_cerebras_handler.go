// Package handlers provides HTTP request handlers
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sebae/ana/internal/ai"
)

// CerebrasAIRequest represents an incoming request to the Cerebras AI assistant
type CerebrasAIRequest struct {
	Query string `json:"query" binding:"required"`
}

// CerebrasAIResponse represents the response from the Cerebras AI assistant
type CerebrasAIResponse struct {
	Response string `json:"response"`
}

// Global Cerebras AI client instance
var cerebrasClient = ai.NewCerebrasClient()

// GetCerebrasAIAssistance handles requests to the Cerebras AI assistant
// Supports special commands like /no_think for direct answers
func GetCerebrasAIAssistance(c *gin.Context) {
	var request CerebrasAIRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check for /no_think command
	isNoThink := false
	query := request.Query
	if strings.HasPrefix(strings.ToLower(query), "/no_think") {
		isNoThink = true
		// Remove the command from the query
		query = strings.TrimSpace(query[9:])
		log.Printf("No-think mode detected, processing query: %s", query)
	}

	// System message for context with enhanced architectural domain knowledge
	systemContext := []ai.Message{
		{
			Role:    "system",
			Content: fmt.Sprintf(
				"Eres un asistente especializado en arquitectura para la plataforma ana.world de gestión de proyectos arquitectónicos. %s Tu conocimiento incluye: 1) Normativas colombianas: NSR-10 (Norma Sismo Resistente), POT de Bogotá, Decreto 1077 de 2015, normas urbanísticas locales; 2) Diseño arquitectónico: metodología BIM, diseño paramétrico, estilos arquitectónicos latinoamericanos, soluciones para clima tropical; 3) Gestión de proyectos: metodologías PMI/PRINCE2 adaptadas a construcción, control de cronogramas, gestión de contratistas, licencias de construcción; 4) Materiales sostenibles: guadua, tierra compactada, sistemas pasivos de climatización, certificación LEED/EDGE para Colombia; 5) Presupuestos: estimación de costos por m², control de presupuestos, análisis de precios unitarios (APU). %s Si te preguntan en inglés, comprende la consulta pero responde en español. %s",
				// Add no_think instructions if applicable
				isNoThink ? "Estás en modo respuesta directa. Responde con la información precisa sin explicaciones adicionales, usando el mínimo de palabras posible." : "",
				// Regular response style
				isNoThink ? "Responde solo con datos concretos sin introducción ni explicación." : "Responde siempre en español con terminología técnica precisa.",
				// Technical information handling
				isNoThink ? "" : "Cuando proporciones información técnica, incluye referencias a códigos específicos, ejemplos prácticos y consideraciones para el contexto colombiano.",
			),
		},
	}

	// Generate response using Cerebras API
	response, err := cerebrasClient.GenerateAssistantResponse(query, systemContext)
	if err != nil {
		log.Printf("Cerebras AI assistant error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Cerebras AI response"})
		return
	}
	
	// For no_think mode, we may want to format the response differently
	if isNoThink {
		// Remove any explanatory text and keep only the direct answer
		// This is a simple heuristic that might need refinement based on actual responses
		lines := strings.Split(response, "\n")
		if len(lines) > 0 {
			// If multi-line, try to find the most concise, relevant line
			shortestRelevantLine := response
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				// Skip empty lines or very short lines that might be headings
				if len(trimmed) > 5 && len(trimmed) < len(shortestRelevantLine) {
					// Check if it contains numbers or specific units which often indicate an answer
					if strings.ContainsAny(trimmed, "0123456789") || 
					   strings.Contains(trimmed, "kg") ||
					   strings.Contains(trimmed, "m²") ||
					   strings.Contains(trimmed, "mm") {
						shortestRelevantLine = trimmed
					}
				}
			}
			response = shortestRelevantLine
		}
	}

	c.JSON(http.StatusOK, CerebrasAIResponse{
		Response: response,
	})
}

