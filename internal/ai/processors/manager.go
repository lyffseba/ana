// AI Processor Manager for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package processors

import (
    "context"
    "fmt"
    "sync"
)

// ProcessorManager manages multiple AI processors
type ProcessorManager struct {
    processors map[string]Processor
    mu         sync.RWMutex
}

// Processor interface for AI processors
type Processor interface {
    Process(ctx context.Context, input []byte) ([]byte, error)
    Train(ctx context.Context, data []byte) error
    Evaluate(ctx context.Context) (*EvaluationResult, error)
    GetMetrics() *Metrics
    Reset()
}

// NewProcessorManager creates a new processor manager
func NewProcessorManager() *ProcessorManager {
    return &ProcessorManager{
        processors: make(map[string]Processor),
    }
}

// RegisterProcessor registers a new processor
func (m *ProcessorManager) RegisterProcessor(name string, processor Processor) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.processors[name] = processor
}

// GetProcessor retrieves a processor by name
func (m *ProcessorManager) GetProcessor(name string) (Processor, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    processor, exists := m.processors[name]
    if !exists {
        return nil, fmt.Errorf("processor not found: %s", name)
    }
    return processor, nil
}

// Process processes input with specified processor
func (m *ProcessorManager) Process(ctx context.Context, processorName string, input []byte) ([]byte, error) {
    processor, err := m.GetProcessor(processorName)
    if err != nil {
        return nil, err
    }
    return processor.Process(ctx, input)
}

// GetAllMetrics returns metrics for all processors
func (m *ProcessorManager) GetAllMetrics() map[string]*Metrics {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    metrics := make(map[string]*Metrics)
    for name, processor := range m.processors {
        metrics[name] = processor.GetMetrics()
    }
    return metrics
}

// ResetAllMetrics resets metrics for all processors
func (m *ProcessorManager) ResetAllMetrics() {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    for _, processor := range m.processors {
        processor.Reset()
    }
}
