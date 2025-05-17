// Package ai provides AI integration services
package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultCerebrasAPIURL = "https://inference.cerebras.ai/v1/chat/completions"
	defaultTimeout        = 30 * time.Second
)

// CerebrasClient handles communication with the Cerebras AI API
type CerebrasClient struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest represents a request to the Cerebras chat API
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// ChatCompletionResponse represents a response from the Cerebras chat API
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// NewCerebrasClient creates a new client for the Cerebras API
func NewCerebrasClient() *CerebrasClient {
	apiKey := getEnv("CEREBRAS_API_KEY", "")
	if apiKey == "" {
		log.Println("Warning: CEREBRAS_API_KEY not set. AI assistant functionality will not work")
	}

	return &CerebrasClient{
		apiKey: apiKey,
		apiURL: getEnv("CEREBRAS_API_URL", defaultCerebrasAPIURL),
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// GenerateTextResponse generates a response to a text-only query
func (c *CerebrasClient) GenerateTextResponse(userQuery string, model string, conversationContext []Message) (string, error) {
	if c.apiKey == "" {
		return "Lo sentimos, el asistente de arquitectura no está disponible en este momento. Por favor contacta al administrador para activar esta funcionalidad.", nil
	}

	// Create a context with user's query
	messages := append(conversationContext, Message{
		Role:    "user",
		Content: userQuery,
	})

	// Create the request body
	requestBody := ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   500,
	}

	// Convert to JSON
	requestBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(requestBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("API error: status code %d, body: %s", resp.StatusCode, string(body))
		
		// Return user-friendly error messages based on status code
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return "No se pudo autenticar con el servicio de IA. Por favor verifica la configuración del asistente arquitectónico.", nil
		case http.StatusForbidden:
			return "No tienes permisos para utilizar el asistente arquitectónico. Por favor contacta al administrador.", nil
		case http.StatusTooManyRequests:
			return "El asistente arquitectónico está experimentando mucho tráfico. Por favor intenta de nuevo en unos momentos.", nil
		case http.StatusServiceUnavailable:
			return "El servicio de asistencia arquitectónica no está disponible temporalmente. Por favor intenta más tarde.", nil
		default:
			return "Hubo un problema al procesar tu consulta arquitectónica. Por favor intenta reformularla o contacta soporte técnico.", nil
		}
	}

	var completionResponse ChatCompletionResponse
	if err := json.Unmarshal(body, &completionResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(completionResponse.Choices) == 0 {
		return "", fmt.Errorf("no completions returned")
	}

	return completionResponse.Choices[0].Message.Content, nil
}

// GenerateVisionResponse generates a response to a query with an image
func (c *CerebrasClient) GenerateVisionResponse(userQuery string, imageBase64 string, conversationContext []Message) (string, error) {
	if c.apiKey == "" {
		return "Lo sentimos, el asistente de visión arquitectónica no está disponible en este momento. Por favor contacta al administrador para activar esta funcionalidad.", nil
	}

	if imageBase64 == "" {
		return "", fmt.Errorf("image data is required for vision model")
	}

	// Determine image format (simple heuristic)
	var imageFormat string
	if len(imageBase64) > 0 {
		imageFormat = "jpeg" // default assumption
	}

	// Create content with image
	imageContent := fmt.Sprintf(
		`%s\n\n![Imagen arquitectónica](data:image/%s;base64,%s)`,
		userQuery,
		imageFormat,
		imageBase64,
	)

	// Create a context with user's query and image
	messages := append(conversationContext, Message{
		Role:    "user",
		Content: imageContent,
	})

	// Create the request body - specifically use vision model
	requestBody := ChatCompletionRequest{
		Model:       "cerebras/QWen-2.5-Vision",
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   800, // Higher for vision descriptions
	}

	// Convert to JSON
	requestBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal vision request: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(requestBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create vision request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send vision request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read vision response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Vision API error: status code %d, body: %s", resp.StatusCode, string(body))
		
		// Return user-friendly error messages based on status code
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return "No se pudo autenticar con el servicio de visión AI. Por favor verifica la configuración del asistente.", nil
		case http.StatusForbidden:
			return "No tienes permisos para utilizar el análisis de imágenes. Por favor contacta al administrador.", nil
		case http.StatusTooManyRequests:
			return "El servicio de análisis de imágenes está experimentando mucho tráfico. Por favor intenta de nuevo en unos momentos.", nil
		case http.StatusServiceUnavailable:
			return "El servicio de análisis de imágenes no está disponible temporalmente. Por favor intenta más tarde.", nil
		case http.StatusRequestEntityTooLarge:
			return "La imagen es demasiado grande. Por favor utiliza una imagen más pequeña (máximo 5MB).", nil
		default:
			return "Hubo un problema al procesar tu imagen arquitectónica. Por favor intenta con otra imagen o contacta soporte técnico.", nil
		}
	}

	var completionResponse ChatCompletionResponse
	if err := json.Unmarshal(body, &completionResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal vision response: %w", err)
	}

	if len(completionResponse.Choices) == 0 {
		return "No se pudo generar un análisis de la imagen. Por favor intenta con otra imagen más clara o con mejor iluminación.", nil
	}

	return completionResponse.Choices[0].Message.Content, nil
}

// GenerateAssistantResponse is a legacy function that calls GenerateTextResponse with default model
// Kept for backward compatibility
func (c *CerebrasClient) GenerateAssistantResponse(userQuery string, conversationContext []Message) (string, error) {
	return c.GenerateTextResponse(userQuery, "cerebras/QWen-3B-32B", conversationContext)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

