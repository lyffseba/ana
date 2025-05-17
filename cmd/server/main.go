package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sebae/ana/internal/database"
	"github.com/sebae/ana/internal/models"
	"github.com/sebae/ana/internal/monitoring"
	"github.com/sebae/ana/internal/repositories"
	"github.com/sebae/ana/internal/server"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or could not be loaded. Using environment variables.")
	}

	// Set default port if not specified
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database connection
	log.Println("Initializing database connection...")
	database.InitDB()

	// Create task repository and seed initial data if needed
	taskRepo := repositories.NewTaskRepository()
	if err := seedInitialData(taskRepo); err != nil {
		log.Printf("Warning: Failed to seed initial data: %v", err)
	}
	
	// Initialize monitoring
	log.Println("Initializing monitoring system...")
	monitoring.Init()
	
	// Set up initial monitoring stats
	setupInitialMonitoringStats()

	// Initialize and start the server
	r := server.SetupRouter()
	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// seedInitialData adds default data if the database is empty
func seedInitialData(taskRepo *repositories.TaskRepository) error {
	// Check if we have tasks, if not create sample tasks
	tasks, err := taskRepo.FindAll()
	if err != nil {
		return err
	}

	// If no tasks exist, seed with default data
	if len(tasks) == 0 {
		log.Println("No tasks found, seeding initial data...")
		
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
		
		log.Println("Initial data seeding completed successfully")
	}
	
	return nil
}

// setupInitialMonitoringStats initializes monitoring statistics
func setupInitialMonitoringStats() {
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
	
	log.Println("Monitoring system initialized")
}
