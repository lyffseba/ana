// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
// Last Updated: Sat May 17 07:34:44 AM CEST 2025

// Package database provides MongoDB connection and management functionality
package database

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	clientOnce sync.Once
)

// InitMongo initializes the MongoDB client (singleton)
func InitMongo() *mongo.Client {
	clientOnce.Do(func() {
		uri := getEnv("MONGODB_URI", "mongodb://localhost:27017/ana_world")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var err error
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
		// Ping to ensure connection
		if err := client.Ping(ctx, nil); err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}
		log.Println("MongoDB connection established")
	})
	return client
}

// GetCollection returns a MongoDB collection handle
var defaultDB = "ana_world"
func GetCollection(dbName, collName string) *mongo.Collection {
	if dbName == "" {
		dbName = defaultDB
	}
	return InitMongo().Database(dbName).Collection(collName)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}


