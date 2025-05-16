// Package handlers provides HTTP request handlers
package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sebae/ana/internal/ai"
)

// CerebrasAIRequest represents an incoming request to the Cerebras AI assistant
type CerebrasAIRequest struct {
	Query string `json:"query" binding:"required"`
}

// CerebrasAIResponse represents the response from the Cerebras AI assistant
type CerebrasAIResponse struct {
	Response string `json:"response"`
}

// Global Cerebras AI client instance
var cerebrasClient = ai.NewCerebrasClient()

// GetCerebrasAIAssistance handles requests to the Cerebras AI assistant
func GetCerebrasAIAssistance(c *gin.Context) {
	var request CerebrasAIRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// System message for context
	systemContext := []ai.Message{
		{
			Role:    "system",
			Content: "You are an AI assistant for the ana.world architectural project management platform. Help users with project management, task organization, and architectural design queries. Answer in Spanish where appropriate.",
		},
	}

	// Generate response using Cerebras API
	response, err := cerebrasClient.GenerateAssistantResponse(request.Query, systemContext)
	if err != nil {
		log.Printf("Cerebras AI assistant error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Cerebras AI response"})
		return
	}

	c.JSON(http.StatusOK, CerebrasAIResponse{
		Response: response,
	})
}

