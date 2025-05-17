// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
// Last Updated: Sat May 17 07:34:44 AM CEST 2025

// Package monitoring provides centralized monitoring functionality for ana.world
package monitoring

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// RequestDuration tracks the duration of API requests
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "ana",
			Name:      "request_duration_seconds",
			Help:      "Duration of API requests",
			Buckets:   []float64{0.1, 0.5, 1, 2, 5, 10, 30},
		},
		[]string{"service", "endpoint", "status"},
	)

	// TokenUsage tracks token usage for AI models
	TokenUsage = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "ana",
			Name:      "token_usage_total",
			Help:      "Token usage for AI models",
		},
		[]string{"service", "model", "type"},
	)

	// CacheMetrics tracks cache performance
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "ana",
			Name:      "cache_hits_total",
			Help:      "Number of cache hits",
		},
		[]string{"service", "cache_type"},
	)

	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "ana",
			Name:      "cache_misses_total",
			Help:      "Number of cache misses",
		},
		[]string{"service", "cache_type"},
	)

	// ErrorCounter tracks errors
	ErrorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "ana",
			Name:      "errors_total",
			Help:      "Number of errors",
		},
		[]string{"service", "error_type"},
	)

	// CircuitBreakerState tracks circuit breaker state changes
	CircuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ana",
			Name:      "circuit_breaker_state",
			Help:      "Circuit breaker state (0=closed, 1=open)",
		},
		[]string{"service"},
	)

	// RateLimiterRejections tracks rate limiter rejections
	RateLimiterRejections = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "ana",
			Name:      "rate_limiter_rejections_total",
			Help:      "Number of requests rejected by rate limiter",
		},
		[]string{"service", "endpoint"},
	)

	// SystemInfo provides system-level metrics
	SystemInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ana",
			Name:      "system_info",
			Help:      "System information (1=always set, value is not significant)",
		},
		[]string{"version", "environment", "go_version"},
	)

	// AggregatedStats stores stats that don't need to be in Prometheus
	aggregatedStats     = make(map[string]interface{})
	aggregatedStatsMutex sync.RWMutex
)

// ServiceStats provides service-specific statistics
type ServiceStats struct {
	CacheStats      CacheStats      `json:"cache_stats"`
	CircuitStats    CircuitStats    `json:"circuit_breaker_stats"`
	PerformanceStats PerformanceStats `json:"performance_stats"`
	ErrorStats      ErrorStats      `json:"error_stats"`
}

// CacheStats represents cache statistics
type CacheStats struct {
	Size         int     `json:"size"`
	HitRate      float64 `json:"hit_rate"`
	HitCount     int64   `json:"hit_count"`
	MissCount    int64   `json:"miss_count"`
	AvgValueSize int     `json:"avg_value_size_bytes"`
}

// CircuitStats represents circuit breaker statistics
type CircuitStats struct {
	State         string    `json:"state"`
	FailureCount  int       `json:"failure_count"`
	LastFailure   time.Time `json:"last_failure,omitempty"`
	ResetTimeout  string    `json:"reset_timeout,omitempty"`
}

// PerformanceStats represents performance statistics
type PerformanceStats struct {
	AvgResponseTime float64 `json:"avg_response_time_ms"`
	RequestCount    int64   `json:"request_count"`
	P95ResponseTime float64 `json:"p95_response_time_ms,omitempty"`
	P99ResponseTime float64 `json:"p99_response_time_ms,omitempty"`
}

// ErrorStats represents error statistics
type ErrorStats struct {
	ErrorCount           int64 `json:"error_count"`
	CircuitBreakerOpens  int64 `json:"circuit_breaker_opens"`
	RateLimitRejections  int64 `json:"rate_limit_rejections"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status      string            `json:"status"`
	Version     string            `json:"version"`
	Environment string            `json:"environment"`
	Services    map[string]string `json:"services"`
	Timestamp   time.Time         `json:"timestamp"`
}

// Init initializes the monitoring system
func Init() {
	// Set system info metric - this metric is always 1 and used for version/env labels
	SystemInfo.WithLabelValues(
		getVersion(),
		getEnvironment(),
		getGoVersion(),
	).Set(1)
}

// RegisterMetricsEndpoint registers the metrics endpoint with the given router
func RegisterMetricsEndpoint(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// RegisterHealthEndpoint registers the health endpoint with the given router
func RegisterHealthEndpoint(router *gin.Engine) {
	router.GET("/health", HealthCheckHandler)
}

// HealthCheckHandler provides a health check endpoint
func HealthCheckHandler(c *gin.Context) {
	// Check services health
	services := map[string]string{
		"api":        "ok",
		"database":   getDatabaseStatus(),
		"ai_service": getAIServiceStatus(),
	}
	
	// Determine overall status
	status := "ok"
	for _, serviceStatus := range services {
		if serviceStatus != "ok" {
			status = "degraded"
			break
		}
	}
	
	c.JSON(http.StatusOK, HealthResponse{
		Status:      status,
		Version:     getVersion(),
		Environment: getEnvironment(),
		Services:    services,
		Timestamp:   time.Now(),
	})
}

// UpdateServiceStats updates the stats for a specific service
func UpdateServiceStats(service string, stats ServiceStats) {
	aggregatedStatsMutex.Lock()
	defer aggregatedStatsMutex.Unlock()
	
	aggregatedStats[service] = stats
}

// GetServiceStats retrieves the stats for a specific service
func GetServiceStats(service string) (ServiceStats, bool) {
	aggregatedStatsMutex.RLock()
	defer aggregatedStatsMutex.RUnlock()
	
	stats, exists := aggregatedStats[service]
	if !exists {
		return ServiceStats{}, false
	}
	
	// Type assertion
	serviceStats, ok := stats.(ServiceStats)
	if !ok {
		return ServiceStats{}, false
	}
	
	return serviceStats, true
}

// StatsHandler handles stats requests
func StatsHandler(c *gin.Context) {
	service := c.Param("service")
	
	if service == "" {
		// Return all stats
		aggregatedStatsMutex.RLock()
		defer aggregatedStatsMutex.RUnlock()
		
		c.JSON(http.StatusOK, aggregatedStats)
		return
	}
	
	// Return stats for specific service
	stats, exists := GetServiceStats(service)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service stats not found"})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

// RegisterStatsEndpoint registers the stats endpoint with the given router
func RegisterStatsEndpoint(router *gin.Engine) {
	router.GET("/stats", StatsHandler)
	router.GET("/stats/:service", StatsHandler)
}

// RecordRequestDuration records the duration of a request
func RecordRequestDuration(service, endpoint string, statusCode int, duration time.Duration) {
	RequestDuration.WithLabelValues(
		service,
		endpoint,
		http.StatusText(statusCode),
	).Observe(duration.Seconds())
}

// RecordTokenUsage records token usage for AI models
func RecordTokenUsage(service, model, tokenType string, count int) {
	TokenUsage.WithLabelValues(
		service,
		model,
		tokenType,
	).Add(float64(count))
}

// RecordCacheHit records a cache hit
func RecordCacheHit(service, cacheType string) {
	CacheHits.WithLabelValues(service, cacheType).Inc()
}

// RecordCacheMiss records a cache miss
func RecordCacheMiss(service, cacheType string) {
	CacheMisses.WithLabelValues(service, cacheType).Inc()
}

// RecordError records an error
func RecordError(service, errorType string) {
	ErrorCounter.WithLabelValues(service, errorType).Inc()
}

// SetCircuitBreakerState sets the circuit breaker state
func SetCircuitBreakerState(service string, isOpen bool) {
	value := 0.0
	if isOpen {
		value = 1.0
	}
	CircuitBreakerState.WithLabelValues(service).Set(value)
}

// RecordRateLimiterRejection records a rate limiter rejection
func RecordRateLimiterRejection(service, endpoint string) {
	RateLimiterRejections.WithLabelValues(service, endpoint).Inc()
}

// Helper functions

func getVersion() string {
	version := os.Getenv("ANA_VERSION")
	if version == "" {
		version = "dev"
	}
	return version
}

func getEnvironment() string {
	env := os.Getenv("ANA_ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	return env
}

func getGoVersion() string {
	// In a real implementation, we'd use runtime.Version()
	// For now, just return a placeholder
	return "go1.20"
}

func getDatabaseStatus() string {
	// In a real implementation, would check database connectivity
	// For now, just return ok
	return "ok"
}

func getAIServiceStatus() string {
	// In a real implementation, would check AI service availability
	// For now, we'll simulate checking the environment variable
	if os.Getenv("CEREBRAS_API_KEY") == "" {
		return "degraded"
	}
	return "ok"
}

