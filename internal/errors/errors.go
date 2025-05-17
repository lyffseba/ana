// Error handling for ANA Project
// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

package errors

import (
    "fmt"
    "net/http"
)

// Error types
type ErrorType string

const (
    ErrorTypeValidation  ErrorType = "validation"
    ErrorTypeNotFound    ErrorType = "not_found"
    ErrorTypeInternal    ErrorType = "internal"
    ErrorTypeUnknown     ErrorType = "unknown"
)

// APIError represents an API error
type APIError struct {
    Type    ErrorType `json:"type"`
    Message string    `json:"message"`
    Code    int       `json:"code"`
    Details any       `json:"details,omitempty"`
}

// Error implements error interface
func (e *APIError) Error() string {
    return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// NewValidationError creates a validation error
func NewValidationError(message string, details any) *APIError {
    return &APIError{
        Type:    ErrorTypeValidation,
        Message: message,
        Code:    http.StatusBadRequest,
        Details: details,
    }
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string) *APIError {
    return &APIError{
        Type:    ErrorTypeNotFound,
        Message: message,
        Code:    http.StatusNotFound,
    }
}

// NewInternalError creates an internal error
func NewInternalError(message string) *APIError {
    return &APIError{
        Type:    ErrorTypeInternal,
        Message: message,
        Code:    http.StatusInternalServerError,
    }
}
