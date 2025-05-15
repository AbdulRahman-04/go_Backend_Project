package utils

import (
	"github.com/redis/go-redis/v9"
	"context"
	"log"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", 
		DB:       0,
	})

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	} else {
		log.Println("✅ Connected to Redis")
	}
}
