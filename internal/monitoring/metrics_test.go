package monitoring

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

// Setup function to create a test Gin router
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// TestMetricsInit tests the initialization of metrics
func TestMetricsInit(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("ANA_VERSION", "test-version")
	os.Setenv("ANA_ENVIRONMENT", "test-env")
	defer func() {
		os.Unsetenv("ANA_VERSION")
		os.Unsetenv("ANA_ENVIRONMENT")
	}()

	// Initialize metrics system
	Init()

	// Can't directly test Prometheus metrics in a unit test,
	// but we can verify that initialization doesn't panic
	t.Log("Metrics initialized successfully")
}

// TestHealthCheckHandler tests the health check endpoint
func TestHealthCheckHandler(t *testing.T) {
	// Setup
	router := setupTestRouter()
	RegisterHealthEndpoint(router)
	
	// Mock environment variables
	os.Setenv("ANA_VERSION", "test-version")
	os.Setenv("ANA_ENVIRONMENT", "test-env")
	os.Setenv("CEREBRAS_API_KEY", "test-key")
	defer func() {
		os.Unsetenv("ANA_VERSION")
		os.Unsetenv("ANA_ENVIRONMENT")
		os.Unsetenv("CEREBRAS_API_KEY")
	}()
	
	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	
	// Check response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response body
	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	
	// Verify response
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, "test-version", response.Version)
	assert.Equal(t, "test-env", response.Environment)
	assert.NotNil(t, response.Services)
	assert.Equal(t, "ok", response.Services["ai_service"], "AI service should be ok")
	assert.NotZero(t, response.Timestamp)
}

// TestStatsHandler tests the stats handler
func TestStatsHandler(t *testing.T) {
	// Setup
	router := setupTestRouter()
	RegisterStatsEndpoint(router)
	
	// Initialize test stats
	testStats := ServiceStats{
		CacheStats: CacheStats{
			Size:         10,
			HitRate:      75.5,
			HitCount:     100,
			MissCount:    32,
			AvgValueSize: 1024,
		},
		CircuitStats: CircuitStats{
			State:        "closed",
			FailureCount: 2,
			LastFailure:  time.Now(),
			ResetTimeout: "60s",
		},
		PerformanceStats: PerformanceStats{
			AvgResponseTime: 150.75,
			RequestCount:    1000,
			P95ResponseTime: 250.0,
			P99ResponseTime: 450.0,
		},
		ErrorStats: ErrorStats{
			ErrorCount:          5,
			CircuitBreakerOpens: 1,
			RateLimitRejections: 10,
		},
	}
	
	UpdateServiceStats("test-service", testStats)
	
	// Test all stats endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stats", nil)
	router.ServeHTTP(w, req)
	
	// Check response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Test specific service stats
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/stats/test-service", nil)
	router.ServeHTTP(w, req)
	
	// Check response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response body
	var response ServiceStats
	err := json.Unmarshal(w.Body.Bytes(), &response)
	
	// Verify response
	assert.NoError(t, err)
	assert.Equal(t, testStats.CacheStats.Size, response.CacheStats.Size)
	assert.Equal(t, testStats.CacheStats.HitRate, response.CacheStats.HitRate)
	assert.Equal(t, testStats.CircuitStats.State, response.CircuitStats.State)
	assert.Equal(t, testStats.PerformanceStats.AvgResponseTime, response.PerformanceStats.AvgResponseTime)
	assert.Equal(t, testStats.ErrorStats.ErrorCount, response.ErrorStats.ErrorCount)
	
	// Test non-existent service
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/stats/non-existent-service", nil)
	router.ServeHTTP(w, req)
	
	// Check response - should be 404
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestRecordFunctions tests the record functions
func TestRecordFunctions(t *testing.T) {
	// Reset registry for clean test
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	
	// Test recording request duration
	RecordRequestDuration("test-service", "/test-endpoint", http.StatusOK, 150*time.Millisecond)
	
	// Test recording token usage
	RecordTokenUsage("test-service", "test-model", "prompt", 100)
	RecordTokenUsage("test-service", "test-model", "completion", 75)
	
	// Test recording cache operations
	RecordCacheHit("test-service", "response-cache")
	RecordCacheMiss("test-service", "response-cache")
	
	// Test recording errors
	RecordError("test-service", "timeout")
	
	// Test circuit breaker state
	SetCircuitBreakerState("test-service", true)
	SetCircuitBreakerState("test-service", false)
	
	// Test rate limiter rejection
	RecordRateLimiterRejection("test-service", "/test-endpoint")
	
	// Since we can't directly observe Prometheus metrics in tests,
	// we just verify that the functions don't panic
	t.Log("Record functions executed successfully")
}

// TestServiceStatsManagement tests the management of service stats
func TestServiceStatsManagement(t *testing.T) {
	// Test adding and retrieving stats
	testStats := ServiceStats{
		CacheStats: CacheStats{
			Size:     5,
			HitRate:  60.0,
			HitCount: 60,
			MissCount: 40,
		},
		PerformanceStats: PerformanceStats{
			AvgResponseTime: 200.0,
			RequestCount:    100,
		},
	}
	
	// Update and get stats
	UpdateServiceStats("test-service", testStats)
	stats, exists := GetServiceStats("test-service")
	
	// Verify results
	assert.True(t, exists)
	assert.Equal(t, testStats.CacheStats.Size, stats.CacheStats.Size)
	assert.Equal(t, testStats.CacheStats.HitRate, stats.CacheStats.HitRate)
	assert.Equal(t, testStats.PerformanceStats.AvgResponseTime, stats.PerformanceStats.AvgResponseTime)
	
	// Test non-existent service
	_, exists = GetServiceStats("non-existent")
	assert.False(t, exists)
}

// TestHelperFunctions tests the helper functions
func TestHelperFunctions(t *testing.T) {
	// Test getVersion
	os.Setenv("ANA_VERSION", "custom-version")
	assert.Equal(t, "custom-version", getVersion())
	os.Unsetenv("ANA_VERSION")
	assert.Equal(t, "dev", getVersion())
	
	// Test getEnvironment
	os.Setenv("ANA_ENVIRONMENT", "custom-env")
	assert.Equal(t, "custom-env", getEnvironment())
	os.Unsetenv("ANA_ENVIRONMENT")
	assert.Equal(t, "development", getEnvironment())
}

// BenchmarkRecordRequestDuration benchmarks the RecordRequestDuration function
func BenchmarkRecordRequestDuration(b *testing.B) {
	// Reset registry for clean benchmark
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RecordRequestDuration("benchmark-service", "/benchmark", http.StatusOK, 100*time.Millisecond)
	}
}

// BenchmarkRecordCacheOperations benchmarks the cache recording functions
func BenchmarkRecordCacheOperations(b *testing.B) {
	// Reset registry for clean benchmark
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			RecordCacheHit("benchmark-service", "response-cache")
		} else {
			RecordCacheMiss("benchmark-service", "response-cache")
		}
	}
}

// BenchmarkUpdateServiceStats benchmarks the UpdateServiceStats function
func BenchmarkUpdateServiceStats(b *testing.B) {
	testStats := ServiceStats{
		CacheStats: CacheStats{
			Size:     10,
			HitRate:  75.0,
			HitCount: 75,
			MissCount: 25,
		},
		PerformanceStats: PerformanceStats{
			AvgResponseTime: 150.0,
			RequestCount:    100,
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UpdateServiceStats("benchmark-service", testStats)
	}
}

