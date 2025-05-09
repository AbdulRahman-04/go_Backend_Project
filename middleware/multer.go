package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// File Upload Middleware for Gin
func FileUploadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ✅ Ensure `uploads/` directory exists
		if _, err := os.Stat("uploads"); os.IsNotExist(err) {
			os.Mkdir("uploads", os.ModePerm)
		}

		// ✅ Retrieve file from request
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "No file uploaded ❌"})
			c.Abort()
			return
		}

		// ✅ Define file storage path
		filePath := "uploads/" + file.Filename

		// ✅ Save file to local storage
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			log.Println("❌ File upload failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "File upload failed ❌"})
			c.Abort()
			return
		}

		// ✅ Store file path in request context
		c.Set("filePath", filePath)

		log.Println("✅ File uploaded successfully:", filePath)
		c.Next() // Proceed to next handler
	}
}