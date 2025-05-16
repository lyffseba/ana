package ai

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

// MockRoundTripper is a mock implementation of http.RoundTripper for testing
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

// RoundTrip implements the http.RoundTripper interface
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

// TestGenerateTextResponse tests the text model response generation
func TestGenerateTextResponse(t *testing.T) {
	// Create a mock round tripper
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			// Check if the request is properly formed
			if req.Header.Get("Authorization") != "Bearer test-api-key" {
				t.Error("Authorization header not set correctly")
			}
			
			if req.Header.Get("Content-Type") != "application/json" {
				t.Error("Content-Type header not set correctly")
			}
			
			// Check request body for correct model
			body, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(body)) // Reset the body for further reading
			
			if !strings.Contains(string(body), "\"model\":\"test-model\"") {
				t.Error("Request body doesn't contain the correct model")
			}
			
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"id": "test-id",
					"object": "chat.completion",
					"created": 1620000000,
					"model": "test-model",
					"choices": [
						{
							"message": {
								"role": "assistant",
								"content": "This is a test response"
							},
							"finish_reason": "stop"
						}
					],
					"usage": {
						"prompt_tokens": 10,
						"completion_tokens": 20,
						"total_tokens": 30
					}
				}`)),
			}, nil
		},
	}
	
	// Create a client with the mock transport
	httpClient := &http.Client{Transport: mockTransport}
	
	client := &CerebrasClient{
		apiKey:     "test-api-key",
		apiURL:     "https://test-url.com",
		httpClient: httpClient,
	}
	
	// Test with a simple query
	response, err := client.GenerateTextResponse("Test query", "test-model", []Message{
		{Role: "system", Content: "You are a test assistant"},
	})
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if response != "This is a test response" {
		t.Errorf("Expected 'This is a test response', got '%s'", response)
	}
}

// TestGenerateVisionResponse tests the vision model response generation
func TestGenerateVisionResponse(t *testing.T) {
	// Create a mock round tripper for vision requests
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			// Check if the request is properly formed
			if req.Header.Get("Authorization") != "Bearer test-api-key" {
				t.Error("Authorization header not set correctly")
			}
			
			// Check request body for image data and vision model
			body, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(body)) // Reset the body for further reading
			
			if !strings.Contains(string(body), "base64,test-image-data") {
				t.Error("Request body doesn't contain the image data")
			}
			
			if !strings.Contains(string(body), "\"model\":\"cerebras/QWen-2.5-Vision\"") {
				t.Error("Request body doesn't use the vision model")
			}
			
			// Mock successful response
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"id": "test-id",
					"object": "chat.completion",
					"created": 1620000000,
					"model": "cerebras/QWen-2.5-Vision",
					"choices": [
						{
							"message": {
								"role": "assistant",
								"content": "This is an analysis of the image"
							},
							"finish_reason": "stop"
						}
					],
					"usage": {
						"prompt_tokens": 10,
						"completion_tokens": 20,
						"total_tokens": 30
					}
				}`)),
			}, nil
		},
	}
	
	// Create a client with the mock transport
	httpClient := &http.Client{Transport: mockTransport}
	
	client := &CerebrasClient{
		apiKey:     "test-api-key",
		apiURL:     "https://test-url.com",
		httpClient: httpClient,
	}
	
	// Test with a query and image
	response, err := client.GenerateVisionResponse("Analyze this image", "test-image-data", []Message{
		{Role: "system", Content: "You are a vision assistant"},
	})
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if response != "This is an analysis of the image" {
		t.Errorf("Expected 'This is an analysis of the image', got '%s'", response)
	}
}

// TestErrorHandling tests the error handling in the client
func TestErrorHandling(t *testing.T) {
	// Create a mock round tripper for error testing
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			// Mock error response
			return &http.Response{
				StatusCode: 401,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"error": {
						"message": "Invalid API key",
						"type": "invalid_request_error"
					}
				}`)),
			}, nil
		},
	}
	
	// Create a client with the mock transport
	httpClient := &http.Client{Transport: mockTransport}
	
	client := &CerebrasClient{
		apiKey:     "test-api-key",
		apiURL:     "https://test-url.com",
		httpClient: httpClient,
	}
	
	// Test with a simple query
	response, err := client.GenerateTextResponse("Test query", "test-model", []Message{
		{Role: "system", Content: "You are a test assistant"},
	})
	
	// Should return a user-friendly error message without throwing an error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if !strings.Contains(response, "No se pudo autenticar") {
		t.Errorf("Expected authentication error message, got '%s'", response)
	}
}

// TestNoAPIKey tests the behavior when no API key is provided
func TestNoAPIKey(t *testing.T) {
	// Create a client with no API key
	client := &CerebrasClient{
		apiKey:     "",
		apiURL:     "https://test-url.com",
		httpClient: &http.Client{},
	}
	
	// Test text response
	textResponse, err := client.GenerateTextResponse("Test query", "test-model", nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !strings.Contains(textResponse, "Lo sentimos, el asistente de arquitectura no está disponible") {
		t.Errorf("Expected service unavailable message, got '%s'", textResponse)
	}
	
	// Test vision response
	visionResponse, err := client.GenerateVisionResponse("Test query", "test-image", nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !strings.Contains(visionResponse, "Lo sentimos, el asistente de visión arquitectónica no está disponible") {
		t.Errorf("Expected vision service unavailable message, got '%s'", visionResponse)
	}
}

// TestVisionNoImage tests the vision model with no image
func TestVisionNoImage(t *testing.T) {
	client := &CerebrasClient{
		apiKey: "test-api-key",
		apiURL: "https://test-url.com",
	}
	
	// Test vision response with empty image
	_, err := client.GenerateVisionResponse("Test query", "", nil)
	if err == nil {
		t.Error("Expected error for missing image, got nil")
	}
	if !strings.Contains(err.Error(), "image data is required") {
		t.Errorf("Expected 'image data is required' error, got '%v'", err)
	}
}

