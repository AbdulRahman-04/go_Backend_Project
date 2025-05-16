package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"Go_Backend/utils"
	"github.com/gin-gonic/gin"
)

const (
	RateLimitCount  = 5           // ✅ 5 requests allowed per minute
	RateLimitWindow = time.Minute // ✅ 1-minute window
)

func RateLimitMiddleware(group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()
		endpoint := c.FullPath() // ✅ Get the exact endpoint path
		key := "ratelimit:" + group + ":" + ip + ":" + endpoint

		// 🔥 Redis key increment
		count, err := utils.RedisClient.Incr(ctx, key).Result()
		if err != nil {
			log.Println("❌ Redis error:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			return
		}

		// 🔥 Fix: Expiry sirf tab set karo jab pehli request aaye (count == 1)
		exists, _ := utils.RedisClient.Exists(ctx, key).Result()
		if exists == 1 { // ✅ Agar key already exist karti hai toh expiry set NA karo!
			ttl, _ := utils.RedisClient.TTL(ctx, key).Result()
			log.Printf("🚀 RateLimit [%s] | IP: %s | Endpoint: %s | Count: %d | TTL Remaining: %v", group, ip, endpoint, count, ttl)
		} else { // ✅ Agar key nayi bani hai toh expiry set karo!
			utils.RedisClient.Expire(ctx, key, RateLimitWindow)
			log.Printf("🚀 RateLimit [%s] | IP: %s | Endpoint: %s | Count: %d | TTL SET to: %v", group, ip, endpoint, count, RateLimitWindow)
		}

		// ✅ Agar request limit exceed ho toh block karo
		if count > RateLimitCount {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			return
		}

		c.Next()
	}
}