// AI API Handlers for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package api

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/lyffseba/ana/internal/ai"
    "github.com/lyffseba/ana/internal/metrics"
)

// AIHandler handles AI-related API endpoints
type AIHandler struct {
    aiService  *ai.AIService
    metrics *metrics.Metrics
}

// NewAIHandler creates a new AI handler
func NewAIHandler(aiService *ai.AIService, metrics *metrics.Metrics) *AIHandler {
    return &AIHandler{
        aiService:  aiService,
        metrics: metrics,
    }
}

// ProcessRequest handles AI processing requests
func (h *AIHandler) ProcessRequest(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    defer h.recordMetrics("process_request", start)

    // Parse request
    var task ai.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        h.handleError(w, "Invalid request format", http.StatusBadRequest)
        return
    }

    // Process task
    ctx := r.Context()
    result, err := h.aiService.ProcessTask(ctx, &task)
    if err != nil {
        h.handleError(w, "Failed to process task: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Return result
    h.respondJSON(w, result)
}

// GetModels returns available AI models
func (h *AIHandler) GetModels(w http.ResponseWriter, r *http.Request) {
    // Return a static list of models for now
    models := []string{"gpt-3.5-turbo", "gpt-4"}
    h.respondJSON(w, models)
}

// HealthCheck returns AI system health status
func (h *AIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    health := h.aiService.HealthCheck(r.Context())
    h.respondJSON(w, health)
}

func (h *AIHandler) handleError(w http.ResponseWriter, message string, status int) {
    h.metrics.RecordError("ai_api_error")
    http.Error(w, message, status)
}

func (h *AIHandler) respondJSON(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func (h *AIHandler) recordMetrics(operation string, start time.Time) {
    duration := time.Since(start).Seconds()
    h.metrics.RecordAIProcessing(operation, duration)
}

// RegisterRoutes registers AI endpoints
func (h *AIHandler) RegisterRoutes(router *http.ServeMux) {
    router.HandleFunc("/api/v1/ai/process", h.ProcessRequest)
    router.HandleFunc("/api/v1/ai/models", h.GetModels)
    router.HandleFunc("/api/v1/ai/health", h.HealthCheck)
}
