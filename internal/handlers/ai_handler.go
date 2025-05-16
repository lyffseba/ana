// Package handlers provides HTTP request handlers
package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juansgv/ana/internal/ai"
)

// AIRequest represents an incoming AI assistant request
type AIRequest struct {
	Query string `json:"query" binding:"required"`
}

// AIResponse represents the response from the AI assistant
type AIResponse struct {
	Response string `json:"response"`
}

// Global AI client instance
var aiClient = ai.NewCerebrasClient()

// GetAIAssistance handles requests to the AI assistant
func GetAIAssistance(c *gin.Context) {
	var request AIRequest
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
	response, err := aiClient.GenerateAssistantResponse(request.Query, systemContext)
	if err != nil {
		log.Printf("AI assistant error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate AI response"})
		return
	}

	c.JSON(http.StatusOK, AIResponse{
		Response: response,
	})
}

