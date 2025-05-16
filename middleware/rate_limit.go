package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiters sync.Map // global map to store rate limiters per key

// getLimiter returns the rate limiter for the given key.
// If none exists, it creates a new one with the specified rate (5 req/min) and burst 5.
func getLimiter(key string) *rate.Limiter {
	if lim, ok := limiters.Load(key); ok {
		return lim.(*rate.Limiter)
	}
	// 5 requests per minute => 5/60 per second
	limiter := rate.NewLimiter(rate.Limit(5)/60, 5)
	limiters.Store(key, limiter)
	return limiter
}

// RateLimitMiddlewareInMemory provides in-memory rate limiting per endpoint, per IP.
func RateLimitMiddlewareInMemory(group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		endpoint := c.FullPath() // unique per endpoint
		key := "ratelimit:" + group + ":" + ip + ":" + endpoint

		limiter := getLimiter(key)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			return
		}
		c.Next()
	}
}