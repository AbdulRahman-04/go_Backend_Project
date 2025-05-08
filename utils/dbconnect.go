package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"Go_Backend/config"
)

// Global MongoDB Client
var client *mongo.Client

// ConnectDB initializes a connection to MongoDB
func ConnectDB() error {
	cfg := config.LoadConfig()
	log.Printf("🔄 Connecting to MongoDB at %s...", cfg.DBUrl)

	clientOptions := options.Client().ApplyURI(cfg.DBUrl)

	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Ping to confirm connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
		return fmt.Errorf("database ping failed: %v", err)
	}

	log.Println("✅ DATABASE CONNECTED SUCCESSFULLY!")
	return nil
}

// GetCollection returns the specified MongoDB collection
func GetCollection(name string) *mongo.Collection {
	if client == nil {
		log.Fatalf("❌ MongoDB client is nil! Check if ConnectDB() failed.")
	}
	return client.Database("your_database_name").Collection(name)
}