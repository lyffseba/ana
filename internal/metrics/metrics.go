// Metrics collection for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics collects system metrics
type Metrics struct {
    requestCounter   *prometheus.CounterVec
    requestDuration  *prometheus.HistogramVec
    errorCounter     *prometheus.CounterVec
    aiProcessing     *prometheus.HistogramVec
    wsConnections    prometheus.Gauge
}

// NewMetrics creates a new metrics collector
func NewMetrics() *Metrics {
    m := &Metrics{}
    
    // Initialize request metrics
    m.requestCounter = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "ana_http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    m.requestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "ana_http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
    
    // Initialize error metrics
    m.errorCounter = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "ana_errors_total",
            Help: "Total number of errors",
        },
        []string{"type"},
    )
    
    // Initialize AI metrics
    m.aiProcessing = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "ana_ai_processing_duration_seconds",
            Help:    "AI processing duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"model"},
    )
    
    // Initialize WebSocket metrics
    m.wsConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "ana_websocket_connections",
            Help: "Current number of WebSocket connections",
        },
    )
    
    return m
}

// RecordRequest records an HTTP request
func (m *Metrics) RecordRequest(method, path string, status int, duration float64) {
    m.requestCounter.WithLabelValues(method, path, fmt.Sprintf("%d", status)).Inc()
    m.requestDuration.WithLabelValues(method, path).Observe(duration)
}

// RecordError records an error
func (m *Metrics) RecordError(errorType string) {
    m.errorCounter.WithLabelValues(errorType).Inc()
}

// RecordAIProcessing records AI processing duration
func (m *Metrics) RecordAIProcessing(model string, duration float64) {
    m.aiProcessing.WithLabelValues(model).Observe(duration)
}

// IncrementWSConnections increments WebSocket connections
func (m *Metrics) IncrementWSConnections() {
    m.wsConnections.Inc()
}

// DecrementWSConnections decrements WebSocket connections
func (m *Metrics) DecrementWSConnections() {
    m.wsConnections.Dec()
}
