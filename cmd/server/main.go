package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sebae/ana/internal/database"
	"github.com/sebae/ana/internal/models"
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

