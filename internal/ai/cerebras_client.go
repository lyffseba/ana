// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
// Last Updated: Wed May 28 08:12:00 AM CEST 2025

// Package ai provides AI integration services
package ai

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/sync/semaphore"
)

const (
	defaultCerebrasAPIURL = "https://api.cerebras.ai/v1/chat/completions"
	defaultTimeout        = 30 * time.Second
	defaultCacheTTL       = 15 * time.Minute
	defaultMaxRetries     = 3
	defaultMaxConcurrent  = 10
)

// CachedResponse represents a cached API response
type CachedResponse struct {
	Response string
	Expiry   time.Time
}

// CircuitBreakerState tracks the circuit breaker state
type CircuitBreakerState struct {
	Failures       int
	LastFailure    time.Time
	Open           bool
	OpenUntil      time.Time
	ThresholdCount int
	ResetTimeout   time.Duration
	mutex          sync.RWMutex
}

// CerebrasClient handles communication with the Cerebras AI API
type CerebrasClient struct {
	apiKey             string
	apiURL             string
	httpClient         *retryablehttp.Client
	cache              map[string]CachedResponse
	cacheTTL           time.Duration
	cacheMutex         sync.RWMutex
	circuitBreaker     CircuitBreakerState
	concurrencyLimiter *semaphore.Weighted
	metricsEnabled     bool
}

// metrics for monitoring client performance
var (
	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cerebras_request_duration_seconds",
			Help:    "Duration of requests to Cerebras API",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint", "status"},
	)

	tokenUsage = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cerebras_token_usage_total",
			Help: "Token usage stats for Cerebras API requests",
		},
		[]string{"type"},
	)

	cacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cerebras_cache_hits_total",
			Help: "Number of cache hits",
		},
	)

	cacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cerebras_cache_misses_total",
			Help: "Number of cache misses",
		},
	)

	errorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cerebras_errors_total",
			Help: "Number of errors when calling Cerebras API",
		},
		[]string{"type"},
	)

	batchSize = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "cerebras_batch_size",
			Help:    "Size of batched requests to Cerebras API",
			Buckets: []float64{1, 2, 5, 10, 20},
		},
	)
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// FunctionParameters represents the parameters for a function
type FunctionParameters struct {
	Type       string                 `json:"type,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

// Function represents a function that can be called by the model
type Function struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Parameters  FunctionParameters `json:"parameters,omitempty"`
}

// Tool represents a tool that can be used by the model
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// ToolChoice represents the tool choice configuration
type ToolChoice struct {
	Type     string   `json:"type,omitempty"`
	Function Function `json:"function,omitempty"`
}

// JSONSchema represents a JSON schema for structured outputs
type JSONSchema struct {
	Name string `json:"name,omitempty"`

	Strict bool                   `json:"strict"`
	Schema map[string]interface{} `json:"schema"`
}

// ResponseFormat represents the format of the model response
type ResponseFormat struct {
	Type       string      `json:"type"`
	JSONSchema *JSONSchema `json:"json_schema,omitempty"`
}

// ChatCompletionRequest represents a request to the Cerebras chat API
type ChatCompletionRequest struct {
	Model          string          `json:"model"`
	Messages       []Message       `json:"messages"`
	Temperature    float64         `json:"temperature,omitempty"`
	MaxTokens      int             `json:"max_completion_tokens,omitempty"`
	Stream         bool            `json:"stream,omitempty"`
	Stop           []string        `json:"stop,omitempty"`
	TopP           float64         `json:"top_p,omitempty"`
	Seed           int             `json:"seed,omitempty"`
	Tools          []Tool          `json:"tools,omitempty"`
	ToolChoice     interface{}     `json:"tool_choice,omitempty"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
	User           string          `json:"user,omitempty"`
	LogProbs       bool            `json:"logprobs,omitempty"`
	TopLogProbs    int             `json:"top_logprobs,omitempty"`
}

// TimeInfo represents timing information in the response
type TimeInfo struct {
	QueueTime      float64 `json:"queue_time"`
	PromptTime     float64 `json:"prompt_time"`
	CompletionTime float64 `json:"completion_time"`
	TotalTime      float64 `json:"total_time"`
	Created        int64   `json:"created"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// FunctionCall represents a function call in the response
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ToolCall represents a tool call in the response
type ToolCall struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	FunctionCall FunctionCall `json:"function"`
}

// ResponseMessage represents a message in the response
type ResponseMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// Choice represents a completion choice in the response
type Choice struct {
	Index        int             `json:"index"`
	Message      ResponseMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
}

// ChatCompletionResponse represents a response from the Cerebras chat API
type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	TimeInfo          TimeInfo `json:"time_info"`
}

// NewCerebrasClient creates a new client for the Cerebras API
func NewCerebrasClient() *CerebrasClient {
	apiKey := getEnv("CEREBRAS_API_KEY", "")
	if apiKey == "" {
		log.Println("Warning: CEREBRAS_API_KEY not set. AI assistant functionality will not work")
	}

	// Create retry client for more robust error handling
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = defaultMaxRetries
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 5 * time.Second
	retryClient.Logger = nil // Disable default logger
	standardClient := retryClient.StandardClient()
	standardClient.Timeout = defaultTimeout

	// Create circuit breaker
	circuitBreaker := CircuitBreakerState{
		ThresholdCount: 5,
		ResetTimeout:   60 * time.Second,
	}

	// Enable metrics based on env var
	metricsEnabled := getEnv("ENABLE_CEREBRAS_METRICS", "false") == "true"

	return &CerebrasClient{
		apiKey:             apiKey,
		apiURL:             getEnv("CEREBRAS_API_URL", defaultCerebrasAPIURL),
		httpClient:         retryClient,
		cache:              make(map[string]CachedResponse),
		cacheTTL:           parseDuration(getEnv("CEREBRAS_CACHE_TTL", "15m"), defaultCacheTTL),
		circuitBreaker:     circuitBreaker,
		concurrencyLimiter: semaphore.NewWeighted(int64(defaultMaxConcurrent)),
		metricsEnabled:     metricsEnabled,
	}
}

// parseDuration parses a duration string and returns a fallback if parsing fails
func parseDuration(durationStr string, fallback time.Duration) time.Duration {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return fallback
	}
	return duration
}

// computeCacheKey generates a unique key for caching based on the request
func computeCacheKey(model string, messages []Message) string {
	// Marshal only what matters for caching
	data := struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}{
		Model:    model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		// If marshaling fails, use a simpler approach
		var key string
		key = model
		for _, msg := range messages {
			key += ":" + msg.Role + ":" + msg.Content
		}
		return key
	}

	hash := md5.Sum(jsonData)
	return hex.EncodeToString(hash[:])
}

// GetCachedResponse retrieves a response from cache if available
func (c *CerebrasClient) GetCachedResponse(model string, messages []Message) (string, bool) {
	if !isCacheable(messages) {
		return "", false
	}

	cacheKey := computeCacheKey(model, messages)

	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()

	if cached, exists := c.cache[cacheKey]; exists {
		if time.Now().Before(cached.Expiry) {
			if c.metricsEnabled {
				cacheHits.Inc()
			}
			return cached.Response, true
		}
	}

	if c.metricsEnabled {
		cacheMisses.Inc()
	}
	return "", false
}

// SetCachedResponse stores a response in the cache
func (c *CerebrasClient) SetCachedResponse(model string, messages []Message, response string) {
	if !isCacheable(messages) {
		return
	}

	cacheKey := computeCacheKey(model, messages)

	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	c.cache[cacheKey] = CachedResponse{
		Response: response,
		Expiry:   time.Now().Add(c.cacheTTL),
	}

	// Prune expired entries occasionally
	if rand.Intn(100) < 5 { // 5% chance to clean up on each set
		c.pruneExpiredCache()
	}
}

// pruneExpiredCache removes expired entries from the cache
func (c *CerebrasClient) pruneExpiredCache() {
	now := time.Now()
	for key, cached := range c.cache {
		if now.After(cached.Expiry) {
			delete(c.cache, key)
		}
	}
}

// isCacheable determines if a request should be cached
// Some requests shouldn't be cached, like those containing
// timestamps or randomness instructions
func isCacheable(messages []Message) bool {
	// Don't cache empty requests
	if len(messages) == 0 {
		return false
	}

	// Check last user message for non-cacheable patterns
	lastUserMsgIdx := -1
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			lastUserMsgIdx = i
			break
		}
	}

	if lastUserMsgIdx == -1 {
		return false
	}

	userMsg := messages[lastUserMsgIdx].Content
	nonCacheablePatterns := []string{
		"random",
		"current time",
		"current date",
		"today",
		"time is",
		"date is",
	}

	for _, pattern := range nonCacheablePatterns {
		if strings.Contains(strings.ToLower(userMsg), pattern) {
			return false
		}
	}

	return true
}

// CheckCircuitBreaker determines if requests should be allowed based on failure history
func (c *CerebrasClient) CheckCircuitBreaker() error {
	c.circuitBreaker.mutex.RLock()
	defer c.circuitBreaker.mutex.RUnlock()

	if c.circuitBreaker.Open {
		if time.Now().Before(c.circuitBreaker.OpenUntil) {
			return fmt.Errorf("circuit breaker is open until %v", c.circuitBreaker.OpenUntil)
		}
	}

	return nil
}

// RecordSuccess resets the failure counter on successful API calls
func (c *CerebrasClient) RecordSuccess() {
	c.circuitBreaker.mutex.Lock()
	defer c.circuitBreaker.mutex.Unlock()

	c.circuitBreaker.Failures = 0
	c.circuitBreaker.Open = false
}

// RecordFailure tracks API failures and opens circuit breaker if threshold is exceeded
func (c *CerebrasClient) RecordFailure() {
	c.circuitBreaker.mutex.Lock()
	defer c.circuitBreaker.mutex.Unlock()

	c.circuitBreaker.Failures++
	c.circuitBreaker.LastFailure = time.Now()

	if c.circuitBreaker.Failures >= c.circuitBreaker.ThresholdCount {
		c.circuitBreaker.Open = true
		c.circuitBreaker.OpenUntil = time.Now().Add(c.circuitBreaker.ResetTimeout)
		log.Printf("Circuit breaker opened until %v after %d failures",
			c.circuitBreaker.OpenUntil, c.circuitBreaker.Failures)
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

	// Convert standard request to retryable request
	retryReq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return "", fmt.Errorf("failed to create retry request: %w", err)
	}

	// Make the request
	resp, err := c.httpClient.Do(retryReq)
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

	// Extract the response content
	responseContent := completionResponse.Choices[0].Message.Content

	// Remove thinking tags if present
	responseContent = removeThinkingTags(responseContent)

	return responseContent, nil
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
	// Note: Check Cerebras documentation for current vision model availability
	requestBody := ChatCompletionRequest{
		Model:       "llama-4-scout-17b-16e-instruct", // Use appropriate vision-capable model
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

	// Convert standard request to retryable request
	retryReq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return "", fmt.Errorf("failed to create retry request for vision: %w", err)
	}

	// Make the request
	resp, err := c.httpClient.Do(retryReq)
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
	return c.GenerateTextResponse(userQuery, "qwen-3-32b", conversationContext)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetAPIStatus returns the status of the API authentication
func (c *CerebrasClient) GetAPIStatus() string {
	if c.apiKey == "" {
		return "missing_key"
	}
	return "ok"
}

// GetCacheSize returns the current number of items in the cache
func (c *CerebrasClient) GetCacheSize() int {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()
	return len(c.cache)
}

// GetCircuitState returns the current state of the circuit breaker
func (c *CerebrasClient) GetCircuitState() string {
	c.circuitBreaker.mutex.RLock()
	defer c.circuitBreaker.mutex.RUnlock()

	if c.circuitBreaker.Open {
		return "open"
	}
	return "closed"
}

// removeThinkingTags removes <think>...</think> tags and their content from the response
func removeThinkingTags(content string) string {
	// Log the original content for debugging
	log.Printf("Original content length: %d, starts with: %.50s", len(content), content)

	// Check if the content contains thinking tags
	thinkStart := strings.Index(content, "<think>")
	if thinkStart == -1 {
		return content // No thinking tags found
	}

	thinkEnd := strings.Index(content, "</think>")
	if thinkEnd == -1 {
		return content // No closing tag found
	}

	// Log the thinking section for debugging
	log.Printf("Found thinking tags from position %d to %d", thinkStart, thinkEnd+8)

	// Remove the thinking section (including tags)
	beforeThink := content[:thinkStart]
	afterThink := ""
	if thinkEnd+8 < len(content) {
		afterThink = content[thinkEnd+8:] // 8 is the length of "</think>"
	}

	result := beforeThink + afterThink

	// Log the result for debugging
	log.Printf("After removing thinking tags, content length: %d", len(result))

	// Recursively check for more thinking tags
	return removeThinkingTags(result)
}
