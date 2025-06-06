// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=ec37c232-0b4b-4a91-af2c-ef680eaa123b
// Last Updated: Sat May 17 08:11:08 AM CEST 2025

// Package handlers provides HTTP request handlers
package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"

	"github.com/lyffseba/ana/internal/ai"
)

// CerebrasAIRequest represents an incoming request to the Cerebras AI assistant
type CerebrasAIRequest struct {
	Query     string                `form:"query" binding:"required"`
	ModelType string                `form:"model_type" binding:"required"`
	Image     *multipart.FileHeader `form:"image"`
}

// CerebrasAIResponse represents the response from the Cerebras AI assistant
type CerebrasAIResponse struct {
	Response     string  `json:"response"`
	HasImage     bool    `json:"has_image,omitempty"`
	FromCache    bool    `json:"from_cache,omitempty"`
	ResponseTime float64 `json:"response_time_ms,omitempty"`
}

// CerebrasHealthResponse represents the health check response
type CerebrasHealthResponse struct {
	Status      string `json:"status"`
	ApiStatus   string `json:"api_status"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

// CerebrasStatsResponse represents cache and performance statistics
type CerebrasStatsResponse struct {
	CacheSize       int     `json:"cache_size"`
	CacheHitRate    float64 `json:"cache_hit_rate"`
	AvgResponseTime float64 `json:"avg_response_time_ms"`
	RequestCount    int64   `json:"request_count"`
	ErrorCount      int64   `json:"error_count"`
	CircuitState    string  `json:"circuit_state"`
}

// UserRateLimiter manages rate limiting per IP
type UserRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
}

// NewUserRateLimiter creates a new rate limiter manager
func NewUserRateLimiter() *UserRateLimiter {
	return &UserRateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

// GetLimiter gets or creates a rate limiter for a given IP
func (u *UserRateLimiter) GetLimiter(ip string) *rate.Limiter {
	u.mu.RLock()
	limiter, exists := u.limiters[ip]
	u.mu.RUnlock()

	if !exists {
		// Create a new limiter that allows 5 requests per minute with a burst of 10
		limiter = rate.NewLimiter(rate.Limit(5/60.0), 10)
		u.mu.Lock()
		u.limiters[ip] = limiter
		u.mu.Unlock()
	}

	return limiter
}

// Global instances
var (
	cerebrasClient *ai.CerebrasClient
	clientOnce     sync.Once
	rateLimiter    = NewUserRateLimiter()

	// Statistics for monitoring
	requestCount    int64
	errorCount      int64
	cacheHits       int64
	cacheMisses     int64
	totalResponseMs float64
	statsMutex      sync.RWMutex
)

// getCerebrasClient returns the Cerebras client, initializing it if needed
func getCerebrasClient() *ai.CerebrasClient {
	clientOnce.Do(func() {
		cerebrasClient = ai.NewCerebrasClient()
	})
	return cerebrasClient
}

// RegisterCerebrasRoutes registers all Cerebras AI-related routes
func RegisterCerebrasRoutes(router *gin.Engine) {
	// AI assistant endpoint
	router.POST("/api/cerebras/assistant", GetCerebrasAIAssistance)

	// Monitoring endpoints
	router.GET("/api/cerebras/health", GetCerebrasHealth)
	router.GET("/api/cerebras/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/api/cerebras/stats", GetCerebrasStats)
}

// GetCerebrasHealth provides health check information
func GetCerebrasHealth(c *gin.Context) {
	// Simple health check that confirms API client is initialized
	apiStatus := "healthy"

	// Check if API key is configured
	if cerebrasClient.GetAPIStatus() != "ok" {
		apiStatus = "degraded"
	}

	c.JSON(http.StatusOK, CerebrasHealthResponse{
		Status:      "ok",
		ApiStatus:   apiStatus,
		Version:     "1.2.0", // Should be configured elsewhere
		Environment: getEnvironment(),
	})
}

// GetCerebrasStats provides cache and performance statistics
func GetCerebrasStats(c *gin.Context) {
	statsMutex.RLock()
	defer statsMutex.RUnlock()

	var hitRate float64 = 0
	totalRequests := cacheHits + cacheMisses
	if totalRequests > 0 {
		hitRate = float64(cacheHits) / float64(totalRequests) * 100
	}

	var avgResponseTime float64 = 0
	if requestCount > 0 {
		avgResponseTime = totalResponseMs / float64(requestCount)
	}

	client := getCerebrasClient()
	c.JSON(http.StatusOK, CerebrasStatsResponse{
		CacheSize:       client.GetCacheSize(),
		CacheHitRate:    hitRate,
		AvgResponseTime: avgResponseTime,
		RequestCount:    requestCount,
		ErrorCount:      errorCount,
		CircuitState:    client.GetCircuitState(),
	})
}

// updateStats updates the global statistics
func updateStats(fromCache bool, responseTimeMs float64, isError bool) {
	statsMutex.Lock()
	defer statsMutex.Unlock()

	requestCount++
	totalResponseMs += responseTimeMs

	if fromCache {
		cacheHits++
	} else {
		cacheMisses++
	}

	if isError {
		errorCount++
	}
}

// getEnvironment returns the current environment (dev, test, prod)
func getEnvironment() string {
	// This should be configured via environment variable
	env := "dev" // Default
	if envVar := os.Getenv("ANA_ENVIRONMENT"); envVar != "" {
		env = envVar
	}
	return env
}

// GetCerebrasAIAssistance handles requests to the Cerebras AI assistant
// Supports both text models (QWen-3B-32B) and vision models (QWen-2.5-Vision)
func GetCerebrasAIAssistance(c *gin.Context) {
	startTime := time.Now()
	fromCache := false
	var responseTimeMs float64
	var response string

	// Apply rate limiting
	clientIP := c.ClientIP()
	limiter := rateLimiter.GetLimiter(clientIP)
	if !limiter.Allow() {
		log.Printf("Rate limit exceeded for IP: %s", clientIP)
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Has excedido el límite de solicitudes. Por favor, intenta de nuevo en un momento."})

		// Update stats
		responseTimeMs = float64(time.Since(startTime).Milliseconds())
		updateStats(false, responseTimeMs, true)
		return
	}

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
		modelType = "qwen-3-32b" // Default to text model
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

	// Construct the system prompt
	systemPrompt := fmt.Sprintf(
		"Eres un asistente especializado en arquitectura para la plataforma ana.world de gestión de proyectos arquitectónicos. %s%sTu conocimiento incluye: 1) Normativas colombianas: NSR-10 (Norma Sismo Resistente), POT de Bogotá, Decreto 1077 de 2015, normas urbanísticas locales; 2) Diseño arquitectónico: metodología BIM, diseño paramétrico, estilos arquitectónicos latinoamericanos, soluciones para clima tropical; 3) Gestión de proyectos: metodologías PMI/PRINCE2 adaptadas a construcción, control de cronogramas, gestión de contratistas, licencias de construcción; 4) Materiales sostenibles: guadua, tierra compactada, sistemas pasivos de climatización, certificación LEED/EDGE para Colombia; 5) Presupuestos: estimación de costos por m², control de presupuestos, análisis de precios unitarios (APU). %s Si te preguntan en inglés, comprende la consulta pero responde en español.",
		noThinkInstructions,
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

	// Check circuit breaker state before making request
	client := getCerebrasClient()
	if err := client.CheckCircuitBreaker(); err != nil {
		log.Printf("Circuit breaker is open: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "El servicio de IA está experimentando problemas temporales. Por favor, intenta de nuevo en unos minutos.",
		})

		// Update stats
		responseTimeMs = float64(time.Since(startTime).Milliseconds())
		updateStats(false, responseTimeMs, true)
		return
	}

	// Use text model only
	modelName := "qwen-3-32b"

	// Check cache first
	cachedResponse, isCached := client.GetCachedResponse(modelName, systemContext)
	if isCached {
		response = cachedResponse
		fromCache = true
		log.Printf("Returning cached response for query")

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
	} else {
		// Generate new response
		var err error
		response, err = client.GenerateTextResponse(query, modelName, systemContext)
		if err != nil {
			log.Printf("Error getting text response: %v", err)
			errorMsg := "Error en el procesamiento de la consulta. Intenta reformularla."

			// Record failure in circuit breaker
			client.RecordFailure()
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
			return
		}
	}

	// Return the response
	responseTimeMs = float64(time.Since(startTime).Milliseconds())
	updateStats(fromCache, responseTimeMs, false)

	c.JSON(http.StatusOK, CerebrasAIResponse{
		Response:     response,
		FromCache:    fromCache,
		ResponseTime: responseTimeMs,
	})
}
