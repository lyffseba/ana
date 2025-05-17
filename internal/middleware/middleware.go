// Middleware for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package middleware

import (
    "net/http"
    "time"

    "go.uber.org/zap"
    
    "github.com/lyffseba/ana/internal/errors"
    "github.com/lyffseba/ana/internal/metrics"
)

// Middleware manages HTTP middleware
type Middleware struct {
    logger  *zap.Logger
    metrics *metrics.Metrics
}

// NewMiddleware creates a new middleware manager
func NewMiddleware(logger *zap.Logger, metrics *metrics.Metrics) *Middleware {
    return &Middleware{
        logger:  logger,
        metrics: metrics,
    }
}

// Logging adds request logging
func (m *Middleware) Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Create response wrapper
        ww := NewResponseWriter(w)
        
        // Process request
        next.ServeHTTP(ww, r)
        
        // Log request
        m.logger.Info("HTTP request",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.Int("status", ww.Status()),
            zap.Duration("duration", time.Since(start)),
        )
    })
}

// Metrics adds metrics collection
func (m *Middleware) Metrics(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Create response wrapper
        ww := NewResponseWriter(w)
        
        // Process request
        next.ServeHTTP(ww, r)
        
        // Record metrics
        m.metrics.RecordRequest(
            r.Method,
            r.URL.Path,
            ww.Status(),
            time.Since(start).Seconds(),
        )
    })
}

// Recovery adds panic recovery
func (m *Middleware) Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // Log error
                m.logger.Error("Panic recovered",
                    zap.Any("error", err),
                    zap.Stack("stack"),
                )
                
                // Record error
                m.metrics.RecordError("panic")
                
                // Return error
                apiError := errors.NewInternalError("Internal server error")
                WriteError(w, apiError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

// Auth adds authentication
func (m *Middleware) Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            apiError := errors.NewValidationError("Missing authorization", nil)
            WriteError(w, apiError)
            return
        }
        
        // TODO: Implement proper token validation
        next.ServeHTTP(w, r)
    })
}

// CORS adds CORS headers
func (m *Middleware) CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// ResponseWriter wraps http.ResponseWriter
type ResponseWriter struct {
    http.ResponseWriter
    status int
    size   int64
}

// NewResponseWriter creates a new response writer
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
    return &ResponseWriter{ResponseWriter: w}
}

// WriteHeader captures the status code
func (w *ResponseWriter) WriteHeader(status int) {
    w.status = status
    w.ResponseWriter.WriteHeader(status)
}

// Write captures the response size
func (w *ResponseWriter) Write(b []byte) (int, error) {
    if w.status == 0 {
        w.status = http.StatusOK
    }
    n, err := w.ResponseWriter.Write(b)
    w.size += int64(n)
    return n, err
}

// Status returns the response status code
func (w *ResponseWriter) Status() int {
    if w.status == 0 {
        return http.StatusOK
    }
    return w.status
}

// Size returns the response size
func (w *ResponseWriter) Size() int64 {
    return w.size
}

// WriteError writes an API error response
func WriteError(w http.ResponseWriter, err *errors.APIError) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(err.Code)
    json.NewEncoder(w).Encode(err)
}
