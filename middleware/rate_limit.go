package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"log"

	"Go_Backend/utils"
	"github.com/gin-gonic/gin"
)

const (
	RateLimitCount  = 5
	RateLimitWindow = time.Minute
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()
		key := "rate:" + ip + ":" + c.Request.Method + ":" + c.Request.URL.Path

		// Increment the request count
		count, err := utils.RedisClient.Incr(ctx, key).Result()
		if err != nil {
			log.Printf("Redis INCR error: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// If first request, set expiration for key
		if count == 1 {
			err := utils.RedisClient.Expire(ctx, key, RateLimitWindow).Err()
			if err != nil {
				log.Printf("Redis EXPIRE error: %v\n", err)
			}
		}

		// Get TTL to send in header
		ttl, err := utils.RedisClient.TTL(ctx, key).Result()
		if err != nil || ttl < 0 {
			ttl = RateLimitWindow
		}

		remaining := RateLimitCount - int(count)
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(RateLimitCount))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))

		if count > int64(RateLimitCount) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}

		c.Next()
	}
}
