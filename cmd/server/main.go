// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
// Last Updated: Sat May 17 07:34:44 AM CEST 2025

package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lyffseba/ana/internal/googleauth"
	"github.com/lyffseba/ana/internal/models"
	"github.com/lyffseba/ana/internal/monitoring"
	"github.com/lyffseba/ana/internal/repositories"
	"github.com/lyffseba/ana/internal/server"
	"go.uber.org/zap"
)

func main() {
	// Initialize Logger
	logger, err := zap.NewProduction() // Or zap.NewDevelopment() for more verbose logs
	if err != nil {
		zap.L().Fatal("can't initialize zap logger", zap.Error(err))
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		sugar.Warnw(".env file not found or could not be loaded. Using environment variables.", "error", err)
	}

	// Set default port if not specified
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// MongoDB: Connection is now handled internally by repositories/database package.
	sugar.Info("MongoDB connection will be established by repository/database package as needed.")

	// Create task repository and seed initial data if needed
	taskRepo := repositories.NewTaskRepository()
	if err := seedInitialData(taskRepo); err != nil {
		sugar.Warnf("Failed to seed initial data: %v", err)
	}
	
	// Initialize monitoring
	sugar.Info("Initializing monitoring system...")
	monitoring.Init()
	
	// Set up initial monitoring stats
	setupInitialMonitoringStats()

	// Initialize Google OAuth Service
	sugar.Info("Initializing Google OAuth Service...")
	// IMPORTANT: Ensure credentials.json is in the 'config' directory at the project root.
	// Add 'config/credentials.json' to your .gitignore file!
	credPath := "config/credentials.json"
	// This redirectURL must match exactly one of the Authorized redirect URIs in your Google Cloud Console
	redirectURL := "http://localhost:8080/api/auth/google/callback"
	scopes := []string{
		"https://www.googleapis.com/auth/calendar",       // Example: full calendar access
		"https://www.googleapis.com/auth/gmail.readonly", // Example: read-only gmail access
		// Add other scopes as needed, e.g., user profile:
		// "https://www.googleapis.com/auth/userinfo.email",
		// "https://www.googleapis.com/auth/userinfo.profile",
	}
	authService, err := googleauth.NewOAuthService(credPath, redirectURL, logger, scopes)
	if err != nil {
		sugar.Fatalf("Failed to initialize Google OAuth Service: %v", err)
	}

	// Initialize and start the server
	r := server.SetupRouter(authService) // Pass authService to SetupRouter
	sugar.Infof("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		sugar.Fatalf("Failed to start server: %v", err)
	}
}

// seedInitialData adds default data if the database is empty
func seedInitialData(taskRepo *repositories.TaskRepository) error {
	// Get a logger instance, assuming sugar is not accessible here or prefer direct logger
	logger, _ := zap.NewProduction() // Basic logger for this function
	defer logger.Sync()

	// Check if we have tasks, if not create sample tasks
	tasks, err := taskRepo.FindAll()
	if err != nil {
		return err
	}

	// If no tasks exist, seed with default data
	if len(tasks) == 0 {
		logger.Info("No tasks found, seeding initial data...")
		
		// Create sample tasks for development
		sampleTasks := []models.Task{
			{
				Title:       "Reunión con Cliente",
				Description: "Discutir requisitos del proyecto para el nuevo edificio residencial",
				DueDate:     time.Now().AddDate(0, 0, 1), // Tomorrow
				Priority:    "High",
				ProjectID:   1,
				Status:      "To-Do",
			},
			{
				Title:       "Finalizar Planos",
				Description: "Completar la versión final de los documentos del plano",
				DueDate:     time.Now().AddDate(0, 0, 3), // In 3 days
				Priority:    "Medium",
				ProjectID:   1,
				Status:      "To-Do",
			},
		}
		
		for _, task := range sampleTasks {
			if err := taskRepo.Create(&task); err != nil {
				return err
			}
		}
		
		logger.Info("Initial data seeding completed successfully")
	}
	
	return nil
}

// setupInitialMonitoringStats initializes monitoring statistics
func setupInitialMonitoringStats() {
	// Get a logger instance
	logger, _ := zap.NewProduction() // Basic logger for this function
	defer logger.Sync()

	// Initialize AI service stats
	aiStats := monitoring.ServiceStats{
		CacheStats: monitoring.CacheStats{
			Size:         0,
			HitRate:      0.0,
			HitCount:     0,
			MissCount:    0,
			AvgValueSize: 0,
		},
		CircuitStats: monitoring.CircuitStats{
			State:        "closed",
			FailureCount: 0,
			ResetTimeout: "60s",
		},
		PerformanceStats: monitoring.PerformanceStats{
			AvgResponseTime: 0.0,
			RequestCount:    0,
		},
		ErrorStats: monitoring.ErrorStats{
			ErrorCount:          0,
			CircuitBreakerOpens: 0,
			RateLimitRejections: 0,
		},
	}
	
	// Register AI service stats
	monitoring.UpdateServiceStats("cerebras", aiStats)
	
	// Set circuit breaker state
	monitoring.SetCircuitBreakerState("cerebras", false)
	
	logger.Info("Monitoring system initialized with initial stats")
}
