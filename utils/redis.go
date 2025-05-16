package utils

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// RedisClient is the global Redis client.
var RedisClient *redis.Client

// InitRedis initializes the Redis connection.
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Change if Redis is remote.
		Password: "",               // Set password if needed.
		DB:       0,
	})

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	} else {
		log.Println("✅ Connected to Redis")
	}
}