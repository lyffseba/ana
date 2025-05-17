// API Routes for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package api

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/lyffseba/ana/internal/ai"
    "github.com/lyffseba/ana/internal/monitoring"
)

// Router handles API routing
type Router struct {
    mux        *http.ServeMux
    middleware *Middleware
    ai         *AIHandler
    monitoring *monitoring.Service
}

// NewRouter creates a new router instance
func NewRouter(aiService *ai.AIService, monitoring *monitoring.Service) *Router {
    router := &Router{
        mux:        http.NewServeMux(),
        middleware: NewMiddleware(monitoring),
        ai:         NewAIHandler(aiService, monitoring),
        monitoring: monitoring,
    }

    router.setupRoutes()
    return router
}

// setupRoutes configures all API routes
func (r *Router) setupRoutes() {
    // AI routes
    r.handle("/api/v1/ai/process", r.ai.ProcessRequest)
    r.handle("/api/v1/ai/models", r.ai.GetModels)
    r.handle("/api/v1/ai/health", r.ai.HealthCheck)

    // Health check
    r.handle("/health", r.healthCheck)

    // Metrics
    r.handle("/metrics", r.metrics)
}

// handle wraps handlers with middleware
func (r *Router) handle(pattern string, handler http.HandlerFunc) {
    wrapped := handler
    wrapped = r.middleware.WithLogging(http.HandlerFunc(wrapped)).ServeHTTP
    wrapped = r.middleware.WithMetrics(http.HandlerFunc(wrapped)).ServeHTTP
    wrapped = r.middleware.WithRateLimit(http.HandlerFunc(wrapped)).ServeHTTP

    // Add authentication for protected routes
    if isProtectedRoute(pattern) {
        wrapped = r.middleware.WithAuth(http.HandlerFunc(wrapped)).ServeHTTP
    }

    r.mux.HandleFunc(pattern, wrapped)
}

// healthCheck handles system health check
func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
    health := struct {
        Status    string `json:"status"`
        Version   string `json:"version"`
        Timestamp string `json:"timestamp"`
    }{
        Status:    "healthy",
        Version:   "1.0.0",
        Timestamp: time.Now().UTC().Format(time.RFC3339),
    }

    respondJSON(w, health)
}

// metrics handles metrics endpoint
func (r *Router) metrics(w http.ResponseWriter, req *http.Request) {
    metrics := r.monitoring.GetMetrics()
    respondJSON(w, metrics)
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    r.mux.ServeHTTP(w, req)
}

// Helper functions
func isProtectedRoute(pattern string) bool {
    // Add logic to determine if route needs authentication
    return true
}

func respondJSON(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
