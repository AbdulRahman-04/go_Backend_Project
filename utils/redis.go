package utils

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// Global Redis client variable
var RedisClient *redis.Client

// InitRedis initializes Redis connection
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // If remote Redis, update this address
		Password: "",               // Set password if applicable
		DB:       0,                // Default Redis database
	})

	// Ping Redis to confirm connection
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	} else {
		log.Println("✅ Connected to Redis")
	}
}