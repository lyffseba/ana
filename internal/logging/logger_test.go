package logging

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/lyffseba/ana/internal/config"
)

func TestNewLogger(t *testing.T) {
    tests := []struct {
        name        string
        config      *config.LoggingConfig
        shouldError bool
    }{
        {
            name: "valid json logger",
            config: &config.LoggingConfig{
                Level:  "info",
                Format: "json",
                Output: "stdout",
            },
            shouldError: false,
        },
        {
            name: "valid console logger",
            config: &config.LoggingConfig{
                Level:  "debug",
                Format: "console",
                Output: "stderr",
            },
            shouldError: false,
        },
        {
            name: "invalid format",
            config: &config.LoggingConfig{
                Level:  "info",
                Format: "invalid",
                Output: "stdout",
            },
            shouldError: true,
        },
        {
            name: "invalid level",
            config: &config.LoggingConfig{
                Level:  "invalid",
                Format: "json",
                Output: "stdout",
            },
            shouldError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            logger, err := NewLogger(tt.config)
            if tt.shouldError {
                assert.Error(t, err)
                assert.Nil(t, logger)
            } else {
                require.NoError(t, err)
                assert.NotNil(t, logger)
            }
        })
    }
}

func TestLoggerWithFields(t *testing.T) {
    config := &config.LoggingConfig{
        Level:  "info",
        Format: "json",
        Output: "stdout",
    }

    logger, err := NewLogger(config)
    require.NoError(t, err)

    fields := map[string]interface{}{
        "key1": "value1",
        "key2": 123,
    }

    loggerWithFields := logger.WithFields(fields)
    assert.NotNil(t, loggerWithFields)
    assert.NotEqual(t, logger, loggerWithFields)
}

func TestLoggerWithError(t *testing.T) {
    config := &config.LoggingConfig{
        Level:  "info",
        Format: "json",
        Output: "stdout",
    }

    logger, err := NewLogger(config)
    require.NoError(t, err)

    testError := assert.AnError
    loggerWithError := logger.WithError(testError)
    assert.NotNil(t, loggerWithError)
    assert.NotEqual(t, logger, loggerWithError)
}
