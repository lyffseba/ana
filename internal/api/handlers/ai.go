// AI API Handlers for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "time"

    "github.com/lyffseba/ana/internal/ai/processors"
)

// AIHandler handles AI-related API endpoints
type AIHandler struct {
    manager    *processors.ProcessorManager
    monitoring *Monitoring
}

// NewAIHandler creates a new AI handler
func NewAIHandler(manager *processors.ProcessorManager) *AIHandler {
    return &AIHandler{
        manager: manager,
        monitoring: NewMonitoring(),
    }
}

// ProcessRequest handles AI processing requests
func (h *AIHandler) ProcessRequest(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    defer h.recordMetrics("process_request", start)

    // Parse request
    var request struct {
        Processor string          `json:"processor"`
        Input     string          `json:"input"`
        Options   map[string]interface{} `json:"options,omitempty"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        h.handleError(w, "Invalid request format", http.StatusBadRequest)
        return
    }

    // Validate request
    if request.Processor == "" || request.Input == "" {
        h.handleError(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    // Process request
    ctx := r.Context()
    result, err := h.processAI(ctx, request.Processor, request.Input, request.Options)
    if err != nil {
        h.handleError(w, "Processing error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Return result
    h.respondJSON(w, result)
}

// processAI handles AI processing
func (h *AIHandler) processAI(ctx context.Context, processorName, input string, options map[string]interface{}) (interface{}, error) {
    // Prepare input
    inputData := map[string]interface{}{
        "text": input,
        "options": options,
    }
    inputBytes, err := json.Marshal(inputData)
    if err != nil {
        return nil, err
    }

    // Process with AI
    result, err := h.manager.Process(ctx, processorName, inputBytes)
    if err != nil {
        return nil, err
    }

    // Parse result
    var response interface{}
    if err := json.Unmarshal(result, &response); err != nil {
        return nil, err
    }

    return response, nil
}

// GetMetrics returns AI processing metrics
func (h *AIHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
    metrics := h.manager.GetAllMetrics()
    h.respondJSON(w, metrics)
}

// ResetMetrics resets AI processing metrics
func (h *AIHandler) ResetMetrics(w http.ResponseWriter, r *http.Request) {
    h.manager.ResetAllMetrics()
    h.respondJSON(w, map[string]string{"status": "metrics reset"})
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
