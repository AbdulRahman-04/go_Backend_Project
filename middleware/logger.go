package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

var logChannel chan string

func init() {
	// Buffered channel with capacity 100 (adjust capacity as needed)
	logChannel = make(chan string, 100)
	go processLogs()
}

// processLogs continuously listens on logChannel and logs messages.
func processLogs() {
	for msg := range logChannel {
		log.Println(msg)
	}
}

// asyncLog pushes log messages into the channel without blocking.
func asyncLog(msg string) {
	select {
	case logChannel <- msg:
		// message queued successfully
	default:
		// If channel is full, log synchronously as a fallback.
		log.Println("Log Channel full, logging synchronously:", msg)
	}
}

// LoggerMiddleware returns a Gin middleware function for logging requests asynchronously.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Capture the request start time.
		c.Next()           // Process request in next middleware / handler.
		duration := time.Since(start)

		method := c.Request.Method         // HTTP method (GET, POST, etc.)
		path := c.Request.URL.Path         // URL path requested.
		statusCode := c.Writer.Status()    // Response HTTP status code.
		clientIP := c.ClientIP()           // Client IP address.
		userAgent := c.Request.UserAgent() // User-Agent header.

		// Choose the log color based on HTTP status.
		var statusColor func(format string, a ...interface{}) string
		switch {
		case statusCode >= 500:
			statusColor = color.New(color.FgRed).SprintfFunc() // 5xx: Red.
		case statusCode >= 400:
			statusColor = color.New(color.FgYellow).SprintfFunc() // 4xx: Yellow.
		case statusCode >= 300:
			statusColor = color.New(color.FgCyan).SprintfFunc() // 3xx: Cyan.
		default:
			statusColor = color.New(color.FgGreen).SprintfFunc() // 2xx: Green.
		}

		// Format the log message.
		msg := fmt.Sprintf("📥 %s %s | %s | %d ms | IP: %s | UA: %s",
			method,
			path,
			statusColor("%d", statusCode),
			duration.Milliseconds(),
			clientIP,
			userAgent,
		)

		// Send the log message asynchronously.
		asyncLog(msg)
	}
}