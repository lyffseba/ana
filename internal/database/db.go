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
// DB is the global database connection instance
var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "ana_user"),
		getEnv("DB_PASSWORD", "your_secure_password"),
		getEnv("DB_NAME", "ana_world"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
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

