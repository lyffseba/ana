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

// GenerateAssistantResponse generates a response to a user query
func (c *CerebrasClient) GenerateAssistantResponse(userQuery string, conversationContext []Message) (string, error) {
	if c.apiKey == "" {
		return "AI assistance is not available. Please contact the administrator.", nil
	}

	// Prepare the messages
	messages := append(conversationContext, Message{
		Role:    "user",
		Content: userQuery,
	})

	// Create the request body
	requestBody := ChatCompletionRequest{
		Model:       "cerebras/Cerebras-GPT-4o",
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   500,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: status code %d, body: %s", resp.StatusCode, string(body))
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

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

