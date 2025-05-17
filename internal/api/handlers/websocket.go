// WebSocket Handler for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package handlers

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/websocket"
    "github.com/lyffseba/ana/internal/ai/processors"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
    manager    *processors.ProcessorManager
    monitoring *Monitoring
    upgrader   websocket.Upgrader
}

// Message represents a WebSocket message
type Message struct {
    Type      string          `json:"type"`
    Processor string          `json:"processor"`
    Input     string          `json:"input"`
    Options   map[string]interface{} `json:"options,omitempty"`
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(manager *processors.ProcessorManager) *WebSocketHandler {
    return &WebSocketHandler{
        manager: manager,
        monitoring: NewMonitoring(),
        upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                // TODO: Implement proper origin check
                return true
            },
        },
    }
}

// HandleConnection handles WebSocket connections
func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
    // Upgrade connection
    conn, err := h.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade error: %v", err)
        return
    }
    defer conn.Close()

    // Handle connection
    h.handleClient(r.Context(), conn)
}

// handleClient handles a WebSocket client
func (h *WebSocketHandler) handleClient(ctx context.Context, conn *websocket.Conn) {
    for {
        // Read message
        _, data, err := conn.ReadMessage()
        if err != nil {
            log.Printf("WebSocket read error: %v", err)
            return
        }

        // Parse message
        var msg Message
        if err := json.Unmarshal(data, &msg); err != nil {
            h.sendError(conn, "Invalid message format")
            continue
        }

        // Process message
        go h.processMessage(ctx, conn, msg)
    }
}

// processMessage processes a WebSocket message
func (h *WebSocketHandler) processMessage(ctx context.Context, conn *websocket.Conn, msg Message) {
    start := time.Now()
    defer h.recordMetrics("process_message", start)

    // Process with AI
    inputData := map[string]interface{}{
        "text": msg.Input,
        "options": msg.Options,
    }
    inputBytes, err := json.Marshal(inputData)
    if err != nil {
        h.sendError(conn, "Invalid input format")
        return
    }

    result, err := h.manager.Process(ctx, msg.Processor, inputBytes)
    if err != nil {
        h.sendError(conn, "Processing error: "+err.Error())
        return
    }

    // Send result
    if err := conn.WriteMessage(websocket.TextMessage, result); err != nil {
        log.Printf("WebSocket write error: %v", err)
        return
    }
}

func (h *WebSocketHandler) sendError(conn *websocket.Conn, message string) {
    response := map[string]string{"error": message}
    data, _ := json.Marshal(response)
    conn.WriteMessage(websocket.TextMessage, data)
}

func (h *WebSocketHandler) recordMetrics(operation string, start time.Time) {
    duration := time.Since(start)
    h.monitoring.RecordDuration("websocket_duration", duration, map[string]string{
        "operation": operation,
    })
}
