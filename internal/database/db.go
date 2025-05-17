// Package database provides database connection and management functionality
package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sebae/ana/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
// DB is the global database instance
var DB *gorm.DB

// buildDSN builds a database connection string
func buildDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),  // Changed from ana_user to postgres
		getEnv("DB_PASSWORD", "postgres"), // Changed to a standard default
		getEnv("DB_NAME", "postgres"),  // Changed from ana_world to postgres
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)
}

// initializeDB is a var so it can be mocked in tests
var initializeDB = func() (*gorm.DB, error) {
	dsn := buildDSN()
	
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
}

// InitDB initializes the database connection
func InitDB() {
	var err error
	DB, err = initializeDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		panic(err)
	}

	log.Println("Database connection established")

	// Auto migrate the Task model (create table if it doesn't exist)
	// This is useful for development but in production we should use proper migrations
	err = DB.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

