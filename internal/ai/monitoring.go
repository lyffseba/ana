// AI Monitoring for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package ai

import (
    "context"
    "time"

    "github.com/prometheus/client_golang/prometheus"
)

var (
    modelExecutionTime = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "ai_model_execution_time_seconds",
            Help: "Time taken to execute AI models",
        },
        []string{"model", "type"},
    )

    modelExecutionErrors = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "ai_model_execution_errors_total",
            Help: "Total number of AI model execution errors",
        },
        []string{"model", "type", "error"},
    )

    modelUsageCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "ai_model_usage_total",
            Help: "Total number of AI model executions",
        },
        []string{"model", "type"},
    )
)

func init() {
    // Register metrics
    prometheus.MustRegister(modelExecutionTime)
    prometheus.MustRegister(modelExecutionErrors)
    prometheus.MustRegister(modelUsageCount)
}

// MonitoredModel wraps a model executor with monitoring
type MonitoredModel struct {
    executor ModelExecutor
    name     string
    typ      string
}

func NewMonitoredModel(executor ModelExecutor, name, typ string) *MonitoredModel {
    return &MonitoredModel{
        executor: executor,
        name:     name,
        typ:      typ,
    }
}

func (m *MonitoredModel) Execute(ctx context.Context, input interface{}) (*ModelResponse, error) {
    start := time.Now()
    
    // Record model usage
    modelUsageCount.WithLabelValues(m.name, m.typ).Inc()

    // Execute model
    response, err := m.executor.Execute(ctx, input)
    
    // Record execution time
    duration := time.Since(start)
    modelExecutionTime.WithLabelValues(m.name, m.typ).Observe(duration.Seconds())

    // Record errors if any
    if err != nil {
        modelExecutionErrors.WithLabelValues(m.name, m.typ, err.Error()).Inc()
    }

    return response, err
}

// Health check for AI system
type HealthCheck struct {
    Status    string
    Message   string
    Timestamp time.Time
    Models    map[string]ModelHealth
}

type ModelHealth struct {
    Status     string
    LastUsed   time.Time
    ErrorCount int
}

func (s *AIService) HealthCheck(ctx context.Context) *HealthCheck {
    health := &HealthCheck{
        Status:    "healthy",
        Message:   "AI system operational",
        Timestamp: time.Now(),
        Models:    make(map[string]ModelHealth),
    }

    // Check each model's health
    s.mu.RLock()
    defer s.mu.RUnlock()

    for name, model := range s.models {
        modelHealth := ModelHealth{
            Status:   "healthy",
            LastUsed: time.Now(), // Replace with actual last used time
        }

        // Get error count from metrics
        if counter, err := modelExecutionErrors.GetMetricWithLabelValues(name, string(model.Version)); err == nil {
            modelHealth.ErrorCount = int(counter.Value())
        }

        health.Models[name] = modelHealth
    }

    return health
}
