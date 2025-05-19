// ANA API Package
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package api

import (
    "context"
    "fmt"
    "net/http"
    "time"
    "encoding/json"

    "go.uber.org/zap"
    "github.com/lyffseba/ana/internal/metrics"
)

// API represents the ANA API service
type API struct {
    server  *http.Server
    logger  *zap.Logger
    config  *Config
    metrics *metrics.Metrics
}

// Config holds API configuration
type Config struct {
    Port         int           `yaml:"port"`
    ReadTimeout  time.Duration `yaml:"read_timeout"`
    WriteTimeout time.Duration `yaml:"write_timeout"`
    IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// NewAPI creates a new API instance
func NewAPI(config *Config, logger *zap.Logger) (*API, error) {
    if config == nil {
        return nil, fmt.Errorf("config is required")
    }
    if logger == nil {
        return nil, fmt.Errorf("logger is required")
    }

    api := &API{
        config:  config,
        logger:  logger,
        metrics: metrics.NewMetrics(),
    }

    api.setupServer()
    return api, nil
}

// Start starts the API server
func (a *API) Start() error {
    a.logger.Info("Starting API server", zap.Int("port", a.config.Port))
    return a.server.ListenAndServe()
}

// Stop stops the API server
func (a *API) Stop(ctx context.Context) error {
    a.logger.Info("Stopping API server")
    return a.server.Shutdown(ctx)
}

// setupServer configures the HTTP server
func (a *API) setupServer() {
    mux := http.NewServeMux()

    // Add routes
    a.addRoutes(mux)

    // Create server
    a.server = &http.Server{
        Addr:         fmt.Sprintf(":%d", a.config.Port),
        Handler:      a.middleware(mux),
        ReadTimeout:  a.config.ReadTimeout,
        WriteTimeout: a.config.WriteTimeout,
        IdleTimeout:  a.config.IdleTimeout,
    }
}

// addRoutes adds API routes
func (a *API) addRoutes(mux *http.ServeMux) {
    // API version 1
    mux.Handle("/api/v1/health", a.handleHealth())
    mux.Handle("/api/v1/metrics", a.handleMetrics())
    
    // AI endpoints
    mux.Handle("/api/v1/ai/process", a.handleAIProcess())
    mux.Handle("/api/v1/ai/models", a.handleAIModels())
    mux.Handle("/api/v1/ai/train", a.handleAITrain())
    
    // WebSocket endpoint
    mux.Handle("/api/v1/ws", a.handleWebSocket())
}

// middleware adds common middleware
func (a *API) middleware(handler http.Handler) http.Handler {
    // Add logging
    handler = a.loggingMiddleware(handler)
    
    // Add metrics
    handler = a.metricsMiddleware(handler)
    
    // Add recovery
    handler = a.recoveryMiddleware(handler)
    
    // Add auth
    handler = a.authMiddleware(handler)
    
    // Add CORS
    handler = a.corsMiddleware(handler)
    
    return handler
}

// Middleware implementations...

func (a *API) loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        
        a.logger.Info("Request processed",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.Duration("duration", time.Since(start)),
        )
    })
}

func (a *API) metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        
        a.metrics.RecordRequest(r.Method, r.URL.Path, 200, time.Since(start).Seconds())
    })
}

func (a *API) recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                a.logger.Error("Panic recovered",
                    zap.Any("error", err),
                    zap.Stack("stack"),
                )
                http.Error(w, "Internal server error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func (a *API) authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // TODO: Implement proper token validation
        next.ServeHTTP(w, r)
    })
}

func (a *API) corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// Route handlers...

func (a *API) handleHealth() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        resp := map[string]string{
            "status": "healthy",
            "version": "1.0.0",
        }
        writeJSON(w, resp)
    })
}

func (a *API) handleMetrics() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // TODO: Implement metrics export
        writeJSON(w, map[string]string{"status": "metrics not implemented"})
    })
}

func (a *API) handleAIProcess() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // TODO: Implement AI processing
        http.Error(w, "Not implemented", http.StatusNotImplemented)
    })
}

func (a *API) handleAIModels() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // TODO: Implement model listing
        http.Error(w, "Not implemented", http.StatusNotImplemented)
    })
}

func (a *API) handleAITrain() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // TODO: Implement model training
        http.Error(w, "Not implemented", http.StatusNotImplemented)
    })
}

func (a *API) handleWebSocket() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // TODO: Implement WebSocket handling
        http.Error(w, "Not implemented", http.StatusNotImplemented)
    })
}

// Helper functions...

func writeJSON(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
