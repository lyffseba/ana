// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
// Last Updated: Sat May 17 07:34:44 AM CEST 2025

package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sebae/ana/internal/handlers"
	"github.com/sebae/ana/internal/monitoring"
)

// MetricsMiddleware adds request metrics collection
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip metrics collection for monitoring endpoints to avoid circular reporting
		path := c.Request.URL.Path
		if path == "/metrics" || path == "/health" || path == "/stats" {
			c.Next()
			return
		}

		// Start timer
		start := time.Now()
		
		// Process request
		c.Next()
		
		// Stop timer and collect metrics
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		endpoint := path
		
		// Record metrics
		monitoring.RecordRequestDuration("api", endpoint, statusCode, duration)
		
		// Record errors if any
		if statusCode >= 400 {
			monitoring.RecordError("api", http.StatusText(statusCode))
		}
	}
}

// SetupRouter configures all the routes for the application
func SetupRouter() *gin.Engine {
	r := gin.Default()
	
	// Add metrics middleware
	r.Use(MetricsMiddleware())

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API routes
	api := r.Group("/api")
	{
		// Task routes
		tasks := api.Group("/tasks")
		{
			tasks.GET("", handlers.GetTasks)
			tasks.GET("/:id", handlers.GetTaskByID)
			tasks.POST("", handlers.CreateTask)
			tasks.PUT("/:id", handlers.UpdateTask)
			tasks.DELETE("/:id", handlers.DeleteTask)
		}

		// Agenda routes
		agenda := api.Group("/agenda")
		{
			agenda.GET("/today", handlers.GetTasksDueToday)
		}
		
		// AI Assistant routes
		ai := api.Group("/ai")
		{
			// Cerebras AI assistant endpoint
			ai.POST("/cerebras", handlers.GetCerebrasAIAssistance)
		}
	}
	
	// Monitoring routes
	monitoring.RegisterHealthEndpoint(r)
	monitoring.RegisterMetricsEndpoint(r)
	monitoring.RegisterStatsEndpoint(r)
	
	// Register handlers for Cerebras monitoring endpoints
	handlers.RegisterCerebrasRoutes(r)

	// Handle frontend routes for development
	// Specifically handle index.html and other static assets
	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})
	
	// Serve individual static files explicitly
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./web/favicon.ico")
	})
	
	// For any other non-API routes, serve the SPA index.html
	// This approach mimics how Netlify would handle SPA routing
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		
		// Don't handle API routes here
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}
		
		// Check if the file exists in the web directory
		if fileExists("./web" + path) {
			c.File("./web" + path)
			return
		}
		
		// Default to serving index.html for client-side routing
		c.File("./web/index.html")
	})

	return r
}

// fileExists checks if a file exists at the given path
func fileExists(path string) bool {
	info, err := http.Dir("./").Open(path)
	if err != nil {
		return false
	}
	info.Close()
	return true
}

