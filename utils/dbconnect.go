package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"Go_Backend/config" // Apna module name use karo
)

// ConnectDB establishes a MongoDB connection
func ConnectDB() (*mongo.Client, error) {
	cfg := config.LoadConfig() // Load config
	clientOptions := options.Client().ApplyURI(cfg.DBUrl)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %v", err)
	}

	// Ping to confirm successful connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("Database ping failed: %v", err)
	}

	log.Println("DATABASE CONNECTED SUCCESSFULLY! ✅")
	return client, nil
}