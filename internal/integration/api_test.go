// Integration tests for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package integration

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/lyffseba/ana/internal/api"
    "github.com/lyffseba/ana/internal/config"
    "github.com/lyffseba/ana/internal/logging"
)

func TestAPIIntegration(t *testing.T) {
    // Create test configuration
    cfg := &config.Config{
        Server: config.ServerConfig{
            Port: 8080,
        },
        AI: config.AIConfig{
            Cerebras: config.CerebrasConfig{
                Endpoint: "http://test-cerebras",
                APIKey:   "test-key",
            },
        },
        Logging: config.LoggingConfig{
            Level:  "debug",
            Format: "json",
            Output: "stdout",
        },
    }

    // Create logger
    logger, err := logging.NewLogger(&cfg.Logging)
    require.NoError(t, err)

    // Create API
    apiServer, err := api.NewAPI(cfg, logger)
    require.NoError(t, err)

    // Create test server
    server := httptest.NewServer(apiServer.Handler())
    defer server.Close()

    // Test cases
    tests := []struct {
        name           string
        method         string
        path           string
        body           interface{}
        expectedStatus int
        expectedBody   interface{}
    }{
        {
            name:           "health check",
            method:         "GET",
            path:          "/health",
            expectedStatus: http.StatusOK,
            expectedBody: map[string]string{
                "status": "healthy",
            },
        },
        {
            name:   "process AI request",
            method: "POST",
            path:   "/api/v1/ai/process",
            body: map[string]interface{}{
                "text": "test input",
            },
            expectedStatus: http.StatusOK,
            expectedBody: map[string]interface{}{
                "processed": true,
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create request
            var body []byte
            if tt.body != nil {
                body, err = json.Marshal(tt.body)
                require.NoError(t, err)
            }

            req, err := http.NewRequest(tt.method, server.URL+tt.path, bytes.NewReader(body))
            require.NoError(t, err)

            if tt.body != nil {
                req.Header.Set("Content-Type", "application/json")
            }

            // Send request
            client := &http.Client{Timeout: 5 * time.Second}
            resp, err := client.Do(req)
            require.NoError(t, err)
            defer resp.Body.Close()

            // Check status
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)

            // Check body
            if tt.expectedBody != nil {
                var result map[string]interface{}
                err = json.NewDecoder(resp.Body).Decode(&result)
                require.NoError(t, err)
                assert.Equal(t, tt.expectedBody, result)
            }
        })
    }
}

func TestWebSocketIntegration(t *testing.T) {
    // Create test configuration
    cfg := &config.Config{
        Server: config.ServerConfig{
            Port: 8080,
        },
        AI: config.AIConfig{
            Cerebras: config.CerebrasConfig{
                Endpoint: "http://test-cerebras",
                APIKey:   "test-key",
            },
        },
        Logging: config.LoggingConfig{
            Level:  "debug",
            Format: "json",
            Output: "stdout",
        },
    }

    // Create logger
    logger, err := logging.NewLogger(&cfg.Logging)
    require.NoError(t, err)

    // Create API
    apiServer, err := api.NewAPI(cfg, logger)
    require.NoError(t, err)

    // Create test server
    server := httptest.NewServer(apiServer.Handler())
    defer server.Close()

    // Create WebSocket connection
    // Note: Implementation depends on the WebSocket client library you're using
    // This is just a placeholder for the test structure
    t.Skip("WebSocket tests to be implemented")
}
