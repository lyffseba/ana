// API Server for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package api

import (
    "context"
    "log"
    "net/http"
    "time"

    "github.com/lyffseba/ana/internal/ai/processors"
    "github.com/lyffseba/ana/internal/api/handlers"
    "github.com/lyffseba/ana/internal/metrics"
)

// Server represents the API server
type Server struct {
    server  *http.Server
    manager *processors.ProcessorManager
    metrics *metrics.Metrics
}

// NewServer creates a new API server
func NewServer(addr string, manager *processors.ProcessorManager) *Server {
    m := metrics.NewMetrics()
    s := &Server{
        manager: manager,
        metrics: m,
    }

    // Create handlers
    aiHandler := handlers.NewAIHandler(manager, m)
    wsHandler := handlers.NewWebSocketHandler(manager, m)

    // Create router
    mux := http.NewServeMux()

    // Register routes
    mux.HandleFunc("/api/v1/ai/process", aiHandler.ProcessRequest)
    mux.HandleFunc("/api/v1/ai/metrics", aiHandler.GetMetrics)
    mux.HandleFunc("/api/v1/ai/metrics/reset", aiHandler.ResetMetrics)
    mux.HandleFunc("/api/v1/ws", wsHandler.HandleConnection)

    // Create server
    s.server = &http.Server{
        Addr:         addr,
        Handler:      s.withMiddleware(mux),
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    return s
}

// Start starts the server
func (s *Server) Start() error {
    log.Printf("Starting server on %s", s.server.Addr)
    return s.server.ListenAndServe()
}

// Stop stops the server
func (s *Server) Stop(ctx context.Context) error {
    return s.server.Shutdown(ctx)
}

// withMiddleware adds middleware to the handler
func (s *Server) withMiddleware(handler http.Handler) http.Handler {
    // Add logging
    handler = s.logMiddleware(handler)

    // Add metrics
    handler = s.metricsMiddleware(handler)

    // Add recovery
    handler = s.recoveryMiddleware(handler)

    return handler
}

// logMiddleware adds request logging
func (s *Server) logMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}

// metricsMiddleware adds metrics collection
func (s *Server) metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        s.metrics.RecordRequest(r.Method, r.URL.Path, 200, time.Since(start).Seconds()) // 200 as placeholder, use actual status if available
    })
}

// recoveryMiddleware adds panic recovery
func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic: %v", err)
                http.Error(w, "Internal server error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
