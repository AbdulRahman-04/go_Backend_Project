package middleware

import (
	"log"
	"time"

	"github.com/fatih/color"  // Terminal colors ke liye
	"github.com/gin-gonic/gin"
)

// LoggerMiddleware returns a gin.HandlerFunc middleware for logging requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()                // Request start time capture karte hain
		c.Next()                          // Request aage next middleware / handler ko jaye
		duration := time.Since(start)     // Request processing ka total time calculate karte hain

		method := c.Request.Method        // HTTP method (GET, POST, etc.)
		path := c.Request.URL.Path        // Request URL path
		statusCode := c.Writer.Status()   // Response ka HTTP status code
		clientIP := c.ClientIP()          // Client ka IP address
		userAgent := c.Request.UserAgent()// Client ka User-Agent header (browser info)

		// Status code ke hisaab se color set kar rahe hain
		var statusColor func(format string, a ...interface{}) string

		switch {
		case statusCode >= 500:
			statusColor = color.New(color.FgRed).SprintfFunc()    // Server errors red color
		case statusCode >= 400:
			statusColor = color.New(color.FgYellow).SprintfFunc() // Client errors yellow color
		case statusCode >= 300:
			statusColor = color.New(color.FgCyan).SprintfFunc()   // Redirects cyan color
		default:
			statusColor = color.New(color.FgGreen).SprintfFunc()  // Success green color
		}

		// Log print kar rahe hain request method, path, colored status code, duration in ms, IP, User-Agent
		log.Printf("📥 %s %s | %s | %d | %s | IP: %s | UA: %s",
			method,
			path,
			statusColor("%d", statusCode),
			duration.Milliseconds(),  // Time in milliseconds for precision
			"ms",
			clientIP,
			userAgent,
		)
	}
}
