// Configuration management for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package config

import (
    "fmt"
    "os"
    "time"

    "gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    AI       AIConfig      `yaml:"ai"`
    Database DatabaseConfig `yaml:"database"`
    Logging  LoggingConfig `yaml:"logging"`
    Metrics  MetricsConfig `yaml:"metrics"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
    Host         string        `yaml:"host"`
    Port         int           `yaml:"port"`
    ReadTimeout  time.Duration `yaml:"read_timeout"`
    WriteTimeout time.Duration `yaml:"write_timeout"`
    IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// AIConfig holds AI configuration
type AIConfig struct {
    Cerebras CerebrasConfig `yaml:"cerebras"`
    Cache    CacheConfig    `yaml:"cache"`
}

// CerebrasConfig holds Cerebras-specific configuration
type CerebrasConfig struct {
    Endpoint    string  `yaml:"endpoint"`
    APIKey      string  `yaml:"api_key"`
    MaxTokens   int     `yaml:"max_tokens"`
    Temperature float64 `yaml:"temperature"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
    Enabled  bool          `yaml:"enabled"`
    Duration time.Duration `yaml:"duration"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    Database string `yaml:"database"`
    SSLMode  string `yaml:"ssl_mode"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
    Level  string `yaml:"level"`
    Format string `yaml:"format"`
    Output string `yaml:"output"`
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
    Enabled bool   `yaml:"enabled"`
    Port    int    `yaml:"port"`
    Path    string `yaml:"path"`
}

// LoadConfig loads configuration from file
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("reading config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("parsing config file: %w", err)
    }

    if err := config.validate(); err != nil {
        return nil, fmt.Errorf("validating config: %w", err)
    }

    config.setDefaults()
    return &config, nil
}

// validate validates the configuration
func (c *Config) validate() error {
    if c.Server.Port <= 0 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }

    if c.AI.Cerebras.Endpoint == "" {
        return fmt.Errorf("cerebras endpoint is required")
    }

    if c.AI.Cerebras.APIKey == "" {
        return fmt.Errorf("cerebras API key is required")
    }

    return nil
}

// setDefaults sets default values
func (c *Config) setDefaults() {
    if c.Server.ReadTimeout == 0 {
        c.Server.ReadTimeout = 30 * time.Second
    }

    if c.Server.WriteTimeout == 0 {
        c.Server.WriteTimeout = 30 * time.Second
    }

    if c.Server.IdleTimeout == 0 {
        c.Server.IdleTimeout = 60 * time.Second
    }

    if c.AI.Cerebras.MaxTokens == 0 {
        c.AI.Cerebras.MaxTokens = 1000
    }

    if c.AI.Cerebras.Temperature == 0 {
        c.AI.Cerebras.Temperature = 0.7
    }

    if c.AI.Cache.Duration == 0 {
        c.AI.Cache.Duration = 1 * time.Hour
    }

    if c.Logging.Level == "" {
        c.Logging.Level = "info"
    }

    if c.Logging.Format == "" {
        c.Logging.Format = "json"
    }

    if c.Metrics.Path == "" {
        c.Metrics.Path = "/metrics"
    }
}
