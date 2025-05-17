// AI API Handlers for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package api

import (
    "context"
    "encoding/json"
    "net/http"
    "time"

    "github.com/lyffseba/ana/internal/ai"
    "github.com/lyffseba/ana/internal/monitoring"
)

// AIHandler handles AI-related API endpoints
type AIHandler struct {
    aiService  *ai.AIService
    monitoring *monitoring.Service
}

// NewAIHandler creates a new AI handler
func NewAIHandler(aiService *ai.AIService, monitoring *monitoring.Service) *AIHandler {
    return &AIHandler{
        aiService:  aiService,
        monitoring: monitoring,
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
    models, err := h.aiService.GetAvailableModels(r.Context())
    if err != nil {
        h.handleError(w, "Failed to get models: "+err.Error(), http.StatusInternalServerError)
        return
    }

    h.respondJSON(w, models)
}

// HealthCheck returns AI system health status
func (h *AIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    health := h.aiService.HealthCheck(r.Context())
    h.respondJSON(w, health)
}

func (h *AIHandler) handleError(w http.ResponseWriter, message string, status int) {
    h.monitoring.RecordError("ai_api_error", message)
    http.Error(w, message, status)
}

func (h *AIHandler) respondJSON(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func (h *AIHandler) recordMetrics(operation string, start time.Time) {
    duration := time.Since(start)
    h.monitoring.RecordDuration("ai_api_duration", duration, map[string]string{
        "operation": operation,
    })
}

// RegisterRoutes registers AI endpoints
func (h *AIHandler) RegisterRoutes(router *http.ServeMux) {
    router.HandleFunc("/api/v1/ai/process", h.ProcessRequest)
    router.HandleFunc("/api/v1/ai/models", h.GetModels)
    router.HandleFunc("/api/v1/ai/health", h.HealthCheck)
}
