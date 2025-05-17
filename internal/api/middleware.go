// API Middleware for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package api

import (
    "context"
    "net/http"
    "time"

    "github.com/lyffseba/ana/internal/monitoring"
)

// Middleware wraps http.Handler with additional functionality
type Middleware struct {
    monitoring *monitoring.Service
}

// NewMiddleware creates a new middleware instance
func NewMiddleware(monitoring *monitoring.Service) *Middleware {
    return &Middleware{
        monitoring: monitoring,
    }
}

// WithLogging adds logging to requests
func (m *Middleware) WithLogging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Process request
        next.ServeHTTP(w, r)

        // Log request
        duration := time.Since(start)
        m.monitoring.RecordHTTPRequest(r.Method, r.URL.Path, duration)
    })
}

// WithMetrics adds metrics collection
func (m *Middleware) WithMetrics(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Wrap response writer to capture status
        wrapped := NewResponseWriter(w)

        // Process request
        next.ServeHTTP(wrapped, r)

        // Record metrics
        duration := time.Since(start)
        m.monitoring.RecordHTTPMetrics(r.Method, r.URL.Path, wrapped.Status(), duration)
    })
}

// WithAuth adds authentication
func (m *Middleware) WithAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Validate token and get user
        user, err := validateToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Add user to context
        ctx := context.WithValue(r.Context(), "user", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// WithRateLimit adds rate limiting
func (m *Middleware) WithRateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get client IP
        ip := r.RemoteAddr

        // Check rate limit
        if isRateLimited(ip) {
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// ResponseWriter wraps http.ResponseWriter to capture status code
type ResponseWriter struct {
    http.ResponseWriter
    status int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
    return &ResponseWriter{ResponseWriter: w}
}

func (w *ResponseWriter) WriteHeader(status int) {
    w.status = status
    w.ResponseWriter.WriteHeader(status)
}

func (w *ResponseWriter) Status() int {
    if w.status == 0 {
        return http.StatusOK
    }
    return w.status
}

// Helper functions (to be implemented)
func validateToken(token string) (interface{}, error) {
    // TODO: Implement token validation
    return nil, nil
}

func isRateLimited(ip string) bool {
    // TODO: Implement rate limiting
    return false
}
