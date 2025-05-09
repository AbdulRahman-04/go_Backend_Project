package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	dbURI := "mongodb+srv://abdrahman:abdrahman@rahmann18.hy9zl.mongodb.net/GO_BACKEND?retryWrites=true&w=majority"

	clientOptions := options.Client().
		ApplyURI(dbURI).
		SetConnectTimeout(15 * time.Second).
		SetServerSelectionTimeout(15 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("❌ mongo.Connect failed: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("❌ Ping to MongoDB failed: %v", err)
	}

	fmt.Println("✅ Connected to MongoDB Successfully!")
}
