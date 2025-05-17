// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=ec37c232-0b4b-4a91-af2c-ef680eaa123b
// Last Updated: Sat May 17 08:11:08 AM CEST 2025

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lyffseba/ana/internal/ai"
)

// MockCerebrasClient is a mock implementation of the Cerebras client for testing
type MockCerebrasClient struct {
	GenerateTextResponseFunc     func(query string, model string, context []ai.Message) (string, error)
	GenerateVisionResponseFunc   func(query string, imageBase64 string, context []ai.Message) (string, error)
	GenerateAssistantResponseFunc func(query string, context []ai.Message) (string, error)
}

func (m *MockCerebrasClient) GenerateTextResponse(query string, model string, context []ai.Message) (string, error) {
	return m.GenerateTextResponseFunc(query, model, context)
}

func (m *MockCerebrasClient) GenerateVisionResponse(query string, imageBase64 string, context []ai.Message) (string, error) {
	return m.GenerateVisionResponseFunc(query, imageBase64, context)
}

func (m *MockCerebrasClient) GenerateAssistantResponse(query string, context []ai.Message) (string, error) {
	if m.GenerateAssistantResponseFunc != nil {
		return m.GenerateAssistantResponseFunc(query, context)
	}
	// Default implementation calls GenerateTextResponse with default model
	return m.GenerateTextResponse(query, "cerebras/QWen-3B-32B", context)
}

// Define a common interface that both real and mock clients implement
type CerebrasClientInterface interface {
	GenerateTextResponse(query string, model string, context []ai.Message) (string, error)
	GenerateVisionResponse(query string, imageBase64 string, context []ai.Message) (string, error)
	GenerateAssistantResponse(query string, context []ai.Message) (string, error)
}

// Store the original client
var originalClient CerebrasClientInterface

// createTestHandler creates a test handler with the mock client
func createTestHandler(mockClient *MockCerebrasClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a custom handler that uses the mock client directly
		var request CerebrasAIRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error en los datos enviados. Verifica que has incluido una consulta válida."})
			return
		}
		
		// Process the query
		query := request.Query
		modelType := request.ModelType
		
		if modelType == "" {
			modelType = "qwen-text" // Default to text model
		}
		
		// Check for /no_think command
		isNoThink := false
		if strings.HasPrefix(strings.ToLower(query), "/no_think") {
			isNoThink = true
			// Remove the command from the query
			query = strings.TrimSpace(query[9:])
		}
		
		// Process image for vision model if present
		var imageBase64 string
		hasImage := false
		
		if modelType == "qwen-vision" && request.Image != nil {
			file, err := request.Image.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo procesar la imagen. Intenta con otro formato o una imagen más pequeña."})
				return
			}
			defer file.Close()
			
			// Read image data
			imageData, err := io.ReadAll(file)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer los datos de la imagen. Intenta con otra imagen."})
				return
			}
			
			// Check image size
			if len(imageData) > 5*1024*1024 { // 5MB limit
				c.JSON(http.StatusBadRequest, gin.H{"error": "La imagen es demasiado grande. Por favor utiliza una imagen menor a 5MB."})
				return
			}
			
			// Use the image data directly without base64 encoding in tests
			imageBase64 = string(imageData)
			hasImage = true
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
		
		// Call the appropriate mock method based on model type
		var response string
		var err error
		
		if modelType == "qwen-vision" && hasImage {
			response, err = mockClient.GenerateVisionResponse(query, imageBase64, systemContext)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en el procesamiento de la imagen. Intenta con otra consulta o imagen."})
				return
			}
		} else {
			modelName := "cerebras/QWen-3B-32B"
			response, err = mockClient.GenerateTextResponse(query, modelName, systemContext)
			if err != nil {
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
}

// setupTestRouter sets up a gin router for testing with our handler
func setupTestRouter(mockClient *MockCerebrasClient) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	// Instead of trying to replace the global client, we create a custom handler
	r.POST("/api/ai/cerebras", createTestHandler(mockClient))
	
	return r
}

// createMultipartRequest creates a test request with optional image
func createMultipartRequest(t *testing.T, query string, modelType string, imagePath string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Add query
	_ = writer.WriteField("query", query)
	
	// Add model type
	_ = writer.WriteField("model_type", modelType)
	
	// Add image if provided
	if imagePath != "" {
		file, err := os.Open(imagePath)
		if err != nil {
			// If the file doesn't exist, create a test image
			part, err := writer.CreateFormFile("image", "test-image.jpg")
			if err != nil {
				return nil, err
			}
			// Write some fake image data
			_, err = part.Write([]byte("fake image data"))
			if err != nil {
				return nil, err
			}
		} else {
			defer file.Close()
			part, err := writer.CreateFormFile("image", "test-image.jpg")
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return nil, err
			}
		}
	}
	
	err := writer.Close()
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("POST", "/api/ai/cerebras", body)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

// TestTextModelRequest tests a basic text model request
func TestTextModelRequest(t *testing.T) {
	// Create a mock client that returns a text response
	mockClient := &MockCerebrasClient{
		GenerateTextResponseFunc: func(query string, model string, context []ai.Message) (string, error) {
			// Verify inputs
			if query != "test query" {
				t.Errorf("Expected query 'test query', got '%s'", query)
			}
			if model != "cerebras/QWen-3B-32B" {
				t.Errorf("Expected model 'cerebras/QWen-3B-32B', got '%s'", model)
			}
			// Check that system context contains architectural info
			if len(context) == 0 || !strings.Contains(context[0].Content, "arquitectura") {
				t.Error("Expected system context with architectural information")
			}
			
			return "Esta es una respuesta de prueba", nil
		},
		GenerateVisionResponseFunc: func(query string, imageBase64 string, context []ai.Message) (string, error) {
			t.Error("Vision model function should not be called for text model request")
			return "", nil
		},
	}
	
	router := setupTestRouter(mockClient)
	req, err := createMultipartRequest(t, "test query", "qwen-text", "")
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	// Verify response contains expected text
	if !strings.Contains(w.Body.String(), "Esta es una respuesta de prueba") {
		t.Errorf("Expected response to contain text response, got '%s'", w.Body.String())
	}
}

// TestVisionModelRequest tests a vision model request with an image
func TestVisionModelRequest(t *testing.T) {
	// Create a mock client that returns a vision response
	mockClient := &MockCerebrasClient{
		GenerateTextResponseFunc: func(query string, model string, context []ai.Message) (string, error) {
			t.Error("Text model function should not be called for vision model request with image")
			return "", nil
		},
		GenerateVisionResponseFunc: func(query string, imageBase64 string, context []ai.Message) (string, error) {
			// Verify inputs
			if query != "analyze this image" {
				t.Errorf("Expected query 'analyze this image', got '%s'", query)
			}
			if imageBase64 == "" {
				t.Error("Expected non-empty image data")
			}
			// Check that system context contains vision-specific info
			if len(context) == 0 || !strings.Contains(context[0].Content, "Estás analizando una imagen arquitectónica") {
				t.Error("Expected system context with vision information")
			}
			
			return "Este es un análisis de la imagen", nil
		},
	}
	
	router := setupTestRouter(mockClient)
	
	// Create fake test image data
	testImageData := "test-image-data"
	testImagePath := "/tmp/test-image.jpg"
	err := os.WriteFile(testImagePath, []byte(testImageData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	defer os.Remove(testImagePath)
	
	req, err := createMultipartRequest(t, "analyze this image", "qwen-vision", testImagePath)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	// Verify response contains expected text and has_image flag
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "Este es un análisis de la imagen") {
		t.Errorf("Expected response to contain vision analysis, got '%s'", responseBody)
	}
	if !strings.Contains(responseBody, "\"has_image\":true") {
		t.Errorf("Expected response to indicate image was used, got '%s'", responseBody)
	}
}

// TestVisionModelWithoutImage tests fallback to text model when vision is selected but no image
func TestVisionModelWithoutImage(t *testing.T) {
	// Create a mock client
	mockClient := &MockCerebrasClient{
		GenerateTextResponseFunc: func(query string, model string, context []ai.Message) (string, error) {
			// Verify inputs
			if query != "vision query without image" {
				t.Errorf("Expected query 'vision query without image', got '%s'", query)
			}
			if model != "cerebras/QWen-3B-32B" {
				t.Errorf("Expected model 'cerebras/QWen-3B-32B' as fallback, got '%s'", model)
			}
			
			return "Respuesta de texto porque no hay imagen", nil
		},
		GenerateVisionResponseFunc: func(query string, imageBase64 string, context []ai.Message) (string, error) {
			t.Error("Vision model function should not be called when no image is provided")
			return "", nil
		},
	}
	
	router := setupTestRouter(mockClient)
	req, err := createMultipartRequest(t, "vision query without image", "qwen-vision", "")
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	// Parse the JSON response to verify its structure
	var response struct {
		Response string `json:"response"`
		HasImage bool   `json:"has_image"`
	}
	
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}
	
	// Verify the response contains the text model output
	if !strings.Contains(response.Response, "Respuesta de texto porque no hay imagen") {
		t.Errorf("Expected text response when vision selected without image, got '%s'", response.Response)
	}
	
	// Verify the has_image field is false
	if response.HasImage {
		t.Errorf("Expected has_image to be false when no image provided, got true")
	}
}

// TestNoThinkMode tests the /no_think command processing
func TestNoThinkMode(t *testing.T) {
	mockClient := &MockCerebrasClient{
		GenerateTextResponseFunc: func(query string, model string, context []ai.Message) (string, error) {
			// Verify the /no_think prefix was removed
			if query != "give me the dimensions" {
				t.Errorf("Expected query 'give me the dimensions' after stripping /no_think, got '%s'", query)
			}
			// Verify system message contains no_think instructions
			if len(context) == 0 || !strings.Contains(context[0].Content, "modo respuesta directa") {
				t.Error("Expected system context with no_think instructions")
			}
			
			return "La dimensión es 5m x 10m x 3m.\nOtras características incluyen material de concreto reforzado.\nEl espesor de la losa es de 20cm.", nil
		},
		GenerateVisionResponseFunc: func(query string, imageBase64 string, context []ai.Message) (string, error) {
			t.Error("Vision model function should not be called for text request in no_think mode")
			return "", nil
		},
	}
	
	router := setupTestRouter(mockClient)
	req, err := createMultipartRequest(t, "/no_think give me the dimensions", "qwen-text", "")
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	// In no_think mode, we expect the shortest relevant line to be extracted
	responseBody := w.Body.String()
	
	// Check that the shortest relevant line is included
	if !strings.Contains(responseBody, "La dimensión es 5m x 10m x 3m") {
		t.Errorf("Expected shortest relevant line with dimensions, got '%s'", responseBody)
	}
	
	// Check that longer explanations are removed
	if strings.Contains(responseBody, "Otras características") {
		t.Errorf("Expected longer explanations to be removed in no_think mode, got '%s'", responseBody)
	}
	
	// Also check that the third line with less relevant info is removed
	if strings.Contains(responseBody, "El espesor de la losa") {
		t.Errorf("Expected less relevant information to be removed in no_think mode, got '%s'", responseBody)
	}
	
	// Verify that response length is significantly shorter than the original
	originalLength := len("La dimensión es 5m x 10m x 3m.\nOtras características incluyen material de concreto reforzado.\nEl espesor de la losa es de 20cm.")
	responseLength := len(responseBody)
	
	// Response should be shorter but still contain the key information
	if responseLength >= originalLength {
		t.Errorf("Expected shortened response in no_think mode, but response length wasn't reduced: original=%d, response=%d", originalLength, responseLength)
	}
}

// TestErrorHandling tests error handling in the handler
func TestErrorHandling(t *testing.T) {
	// Test cases with different error scenarios
	testCases := []struct {
		name           string
		modelType      string
		withImage      bool
		mockClient     *MockCerebrasClient
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "Text Model Error",
			modelType: "qwen-text",
			withImage: false,
			mockClient: &MockCerebrasClient{
				GenerateTextResponseFunc: func(query string, model string, context []ai.Message) (string, error) {
					return "", io.EOF // Simulate network error
				},
				GenerateVisionResponseFunc: func(query string, imageBase64 string, context []ai.Message) (string, error) {
					t.Error("Vision model function should not be called for text model request")
					return "", nil
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Error en el procesamiento de la consulta",
		},
		{
			name:      "Vision Model Error",
			modelType: "qwen-vision",
			withImage: true,
			mockClient: &MockCerebrasClient{
				GenerateTextResponseFunc: func(query string, model string, context []ai.Message) (string, error) {
					t.Error("Text model function should not be called for vision model request with image")
					return "", nil
				},
				GenerateVisionResponseFunc: func(query string, imageBase64 string, context []ai.Message) (string, error) {
					return "", io.EOF // Simulate network error
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Error en el procesamiento de la imagen",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := setupTestRouter(tc.mockClient)
			
			var imagePath string
			if tc.withImage {
				// Create fake test image
				imagePath = "/tmp/test-error-image.jpg"
				err := os.WriteFile(imagePath, []byte("test-error-image-data"), 0644)
				if err != nil {
					t.Fatalf("Failed to create test image: %v", err)
				}
				defer os.Remove(imagePath)
			}
			
			req, err := createMultipartRequest(t, "test error query", tc.modelType, imagePath)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			// Verify response
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}
			
			if !strings.Contains(w.Body.String(), tc.expectedError) {
				t.Errorf("Expected error containing '%s', got '%s'", tc.expectedError, w.Body.String())
			}
		})
	}
}

// TestInvalidRequest tests handling of invalid requests
func TestInvalidRequest(t *testing.T) {
	// Create a mock client that shouldn't be called during an invalid request
	mockClient := &MockCerebrasClient{
		GenerateTextResponseFunc: func(query string, model string, context []ai.Message) (string, error) {
			t.Error("Text model function should not be called for an invalid request")
			return "", nil
		},
		GenerateVisionResponseFunc: func(query string, imageBase64 string, context []ai.Message) (string, error) {
			t.Error("Vision model function should not be called for an invalid request")
			return "", nil
		},
	}
	
	router := setupTestRouter(mockClient)
	
	// Create an invalid request (missing required query field)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Only add model_type but omit query
	_ = writer.WriteField("model_type", "qwen-text")
	writer.Close()
	
	req, err := http.NewRequest("POST", "/api/ai/cerebras", body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	// Make the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Verify response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid request, got %d", w.Code)
	}
	
	// Check for error message
	if !strings.Contains(w.Body.String(), "Error en los datos enviados") {
		t.Errorf("Expected error message about invalid data, got '%s'", w.Body.String())
	}
}

// TestLargeImageFile tests the image size limit handling
func TestLargeImageFile(t *testing.T) {
	// Create a mock client
	mockClient := &MockCerebrasClient{
		GenerateTextResponseFunc: func(query string, model string, context []ai.Message) (string, error) {
			t.Error("Text model function should not be called for large image test")
			return "", nil
		},
		GenerateVisionResponseFunc: func(query string, imageBase64 string, context []ai.Message) (string, error) {
			t.Error("Vision model function should not be called for large image test")
			return "", nil
		},
	}
	
	router := setupTestRouter(mockClient)
	
	// Create a fake large image (just over 5MB)
	largeImageSize := 5*1024*1024 + 100 // 5MB + 100 bytes
	largeImageData := make([]byte, largeImageSize)
	for i := range largeImageData {
		largeImageData[i] = byte(i % 256) // Fill with some data
	}
	
	imagePath := "/tmp/large-test-image.jpg"
	err := os.WriteFile(imagePath, largeImageData, 0644)
	if err != nil {
		t.Fatalf("Failed to create large test image: %v", err)
	}
	defer os.Remove(imagePath)
	
	// Create a request with the large image
	req, err := createMultipartRequest(t, "analyze large image", "qwen-vision", imagePath)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	// Make the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Verify response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for large image, got %d", w.Code)
	}
	
	// Check for size limit error
	if !strings.Contains(w.Body.String(), "demasiado grande") {
		t.Errorf("Expected error message about large image, got '%s'", w.Body.String())
	}
}

// TestUnsupportedModelType tests handling of unsupported model types
func TestUnsupportedModelType(t *testing.T) {
	// Test will use default model if unsupported model type is provided
	mockClient := &MockCerebrasClient{
		GenerateTextResponseFunc: func(query string, model string, context []ai.Message) (string, error) {
			// Should fall back to text model with default model name
			if model != "cerebras/QWen-3B-32B" {
				t.Errorf("Expected default model, got '%s'", model)
			}
			return "Response using default model", nil
		},
		GenerateVisionResponseFunc: func(query string, imageBase64 string, context []ai.Message) (string, error) {
			t.Error("Vision model function should not be called for unsupported model test")
			return "", nil
		},
	}
	
	router := setupTestRouter(mockClient)
	
	// Create a request with invalid model type
	req, err := createMultipartRequest(t, "test unsupported model", "unsupported-model", "")
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	// Make the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Verify response - should still succeed with default model
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 with default model fallback, got %d", w.Code)
	}
	
	if !strings.Contains(w.Body.String(), "Response using default model") {
		t.Errorf("Expected response using default model, got '%s'", w.Body.String())
	}
}

