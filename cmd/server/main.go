package main

import (
	"log"
	"os"

	"github.com/sebae/ana/internal/server"
)

func main() {
	// Set default port if not specified
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize and start the server
	r := server.SetupRouter()
	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

