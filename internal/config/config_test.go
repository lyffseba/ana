package config

import (
    "os"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
    // Create temporary config file
    configContent := `
server:
  host: "localhost"
  port: 8080
ai:
  cerebras:
    endpoint: "http://cerebras.api"
    api_key: "test-key"
`
    tmpfile, err := os.CreateTemp("", "config-*.yaml")
    require.NoError(t, err)
    defer os.Remove(tmpfile.Name())

    _, err = tmpfile.WriteString(configContent)
    require.NoError(t, err)
    tmpfile.Close()

    // Test loading config
    config, err := LoadConfig(tmpfile.Name())
    require.NoError(t, err)

    // Verify values
    assert.Equal(t, "localhost", config.Server.Host)
    assert.Equal(t, 8080, config.Server.Port)
    assert.Equal(t, "http://cerebras.api", config.AI.Cerebras.Endpoint)
    assert.Equal(t, "test-key", config.AI.Cerebras.APIKey)

    // Verify defaults
    assert.Equal(t, 30*time.Second, config.Server.ReadTimeout)
    assert.Equal(t, 30*time.Second, config.Server.WriteTimeout)
    assert.Equal(t, 60*time.Second, config.Server.IdleTimeout)
    assert.Equal(t, 1000, config.AI.Cerebras.MaxTokens)
    assert.Equal(t, 0.7, config.AI.Cerebras.Temperature)
}

func TestConfigValidation(t *testing.T) {
    tests := []struct {
        name        string
        config      Config
        shouldError bool
    }{
        {
            name: "valid config",
            config: Config{
                Server: ServerConfig{Port: 8080},
                AI: AIConfig{
                    Cerebras: CerebrasConfig{
                        Endpoint: "http://cerebras.api",
                        APIKey:   "test-key",
                    },
                },
            },
            shouldError: false,
        },
        {
            name: "invalid port",
            config: Config{
                Server: ServerConfig{Port: 0},
                AI: AIConfig{
                    Cerebras: CerebrasConfig{
                        Endpoint: "http://cerebras.api",
                        APIKey:   "test-key",
                    },
                },
            },
            shouldError: true,
        },
        {
            name: "missing endpoint",
            config: Config{
                Server: ServerConfig{Port: 8080},
                AI: AIConfig{
                    Cerebras: CerebrasConfig{
                        APIKey: "test-key",
                    },
                },
            },
            shouldError: true,
        },
        {
            name: "missing API key",
            config: Config{
                Server: ServerConfig{Port: 8080},
                AI: AIConfig{
                    Cerebras: CerebrasConfig{
                        Endpoint: "http://cerebras.api",
                    },
                },
            },
            shouldError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.validate()
            if tt.shouldError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
