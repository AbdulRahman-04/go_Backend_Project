package middleware

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

var fileProcessQueue chan string

// init creates the file processing queue and starts the background worker.
func init() {
	fileProcessQueue = make(chan string, 100)
	go processFileQueue()
}

// processFileQueue continuously processes file paths enqueued for post-upload tasks.
func processFileQueue() {
	for filePath := range fileProcessQueue {
		log.Println("Processing file asynchronously:", filePath)
		// Simulating a heavy processing task (e.g., image resizing, virus scan, etc.)
		time.Sleep(2 * time.Second)
		log.Println("Finished processing file:", filePath)
	}
}

// processFileSynchronously is the fallback if the queue is full.
func processFileSynchronously(filePath string) {
	log.Println("File queue full; processing file synchronously:", filePath)
	time.Sleep(2 * time.Second)
	log.Println("Finished processing file synchronously:", filePath)
}

// FileUploadMiddleware handles the file upload and enqueues post-upload processing.
func FileUploadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure uploads folder exists.
		if _, err := os.Stat("uploads"); os.IsNotExist(err) {
			if err := os.Mkdir("uploads", os.ModePerm); err != nil {
				log.Println("Error creating uploads folder:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Server error"})
				c.Abort()
				return
			}
		}

		// Retrieve file if provided.
		file, err := c.FormFile("file")
		if err == nil {
			filePath := filepath.Join("uploads", file.Filename)
			if err := c.SaveUploadedFile(file, filePath); err != nil {
				log.Println("❌ File upload failed:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "File upload failed ❌"})
				c.Abort()
				return
			}

			// Store file path in context for later use if needed.
			c.Set("filePath", filePath)
			log.Println("✅ File uploaded successfully:", filePath)

			// Enqueue file for asynchronous processing.
			select {
			case fileProcessQueue <- filePath:
				// Successfully enqueued.
			default:
				// If channel is full, process synchronously.
				processFileSynchronously(filePath)
			}
		}

		c.Next() // Proceed with the next handler.
	}
}