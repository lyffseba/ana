// Logging system for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package logging

import (
    "fmt"
    "os"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"

    "github.com/lyffseba/ana/internal/config"
)

// Logger wraps zap logger with additional functionality
type Logger struct {
    *zap.Logger
    config *config.LoggingConfig
}

// NewLogger creates a new logger
func NewLogger(config *config.LoggingConfig) (*Logger, error) {
    // Create encoder config
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "ts",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        FunctionKey:    zapcore.OmitKey,
        MessageKey:     "msg",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }

    // Create encoder based on format
    var encoder zapcore.Encoder
    switch config.Format {
    case "json":
        encoder = zapcore.NewJSONEncoder(encoderConfig)
    case "console":
        encoder = zapcore.NewConsoleEncoder(encoderConfig)
    default:
        return nil, fmt.Errorf("unsupported log format: %s", config.Format)
    }

    // Create output writer
    var writer zapcore.WriteSyncer
    switch config.Output {
    case "stdout":
        writer = zapcore.AddSync(os.Stdout)
    case "stderr":
        writer = zapcore.AddSync(os.Stderr)
    default:
        // Assume it's a file path
        file, err := os.OpenFile(config.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            return nil, fmt.Errorf("opening log file: %w", err)
        }
        writer = zapcore.AddSync(file)
    }

    // Parse log level
    level, err := zapcore.ParseLevel(config.Level)
    if err != nil {
        return nil, fmt.Errorf("parsing log level: %w", err)
    }

    // Create core
    core := zapcore.NewCore(
        encoder,
        writer,
        level,
    )

    // Create logger
    logger := zap.New(
        core,
        zap.AddCaller(),
        zap.AddStacktrace(zapcore.ErrorLevel),
    )

    return &Logger{
        Logger: logger,
        config: config,
    }, nil
}

// WithFields adds fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
    zapFields := make([]zap.Field, 0, len(fields))
    for k, v := range fields {
        zapFields = append(zapFields, zap.Any(k, v))
    }
    return &Logger{
        Logger: l.With(zapFields...),
        config: l.config,
    }
}

// WithError adds an error field to the logger
func (l *Logger) WithError(err error) *Logger {
    return &Logger{
        Logger: l.With(zap.Error(err)),
        config: l.config,
    }
}
