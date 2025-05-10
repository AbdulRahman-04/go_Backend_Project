package utils

import (
    "context"
    "fmt"
    "log"

    "Go_Backend/config"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// ConnectDB establishes a connection to MongoDB
func ConnectDB() error {
    uri := config.LoadConfig().DBUrl               // ← DBUrl use karo
    clientOptions := options.Client().ApplyURI(uri)

    var err error
    client, err = mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        return fmt.Errorf("error connecting to MongoDB: %v", err)
    }

    // Ping MongoDB to ensure connection is successful
    if err = client.Ping(context.TODO(), nil); err != nil {
        return fmt.Errorf("error pinging MongoDB: %v", err)
    }

    log.Println("✅ Successfully connected to MongoDB")
    return nil
}

// GetCollection returns a collection from the MongoDB database
func GetCollection(collectionName string) *mongo.Collection {
    if client == nil {
        log.Fatal("MongoDB client is nil (did you call ConnectDB?)")
    }
    // "GO_BACKEND" ko apne actual DB name se replace kar sakte ho agar alag ho
    return client.Database("GO_BACKEND").Collection(collectionName)
}
