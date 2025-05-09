package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"Go_Backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global MongoDB Client
var client *mongo.Client

// ConnectDB initializes a connection to MongoDB
func ConnectDB() error {
	cfg := config.LoadConfig()
	log.Printf("DEBUG: Loaded DB_URL: '%s'", cfg.DBUrl)
	log.Printf("🔄 Connecting to MongoDB at %s...", cfg.DBUrl)

	// Extended timeouts for remote connections
	clientOptions := options.Client().
		ApplyURI(cfg.DBUrl).
		SetConnectTimeout(20 * time.Second).
		SetServerSelectionTimeout(20 * time.Second)

	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Ping the database to ensure connection is valid
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
		return fmt.Errorf("database ping failed: %v", err)
	}

	log.Println("✅ DATABASE CONNECTED SUCCESSFULLY!")
	return nil
}

// GetCollection returns the specified MongoDB collection.
// **Important:** Ensure that this is not called during package initialization!
func GetCollection(name string) *mongo.Collection {
	if client == nil {
		log.Fatalf("❌ MongoDB client is nil! Check if ConnectDB() failed.")
		// This code is unreachable; return nil to satisfy the compiler.
		return nil
	}
	return client.Database("GO_BACKEND").Collection(name)
}