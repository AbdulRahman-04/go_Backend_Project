package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"Go_Backend/config"
	"Go_Backend/controllers/private"
	"Go_Backend/controllers/public"
	"Go_Backend/middleware"
	"Go_Backend/utils"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()
	portStr := strconv.Itoa(cfg.Port)

	// Connect to MongoDB first
	log.Println("Connecting to MongoDB... 🔄")
	err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	// First log server live, then DB connected
log.Printf("YOUR SERVER IS LIVE AT PORT %s", portStr)
log.Println("✅ DATABASE CONNECTED SUCCESSFULLY!")

	// Initialize Gin Router
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// Health Check Route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "HELLO IN TS💙"})
	})

	// Public API Routes
	public.SetupPublicRoutes(router)

	// Private Routes with JWT Middleware
	privateGroup := router.Group("/api/private")
	privateGroup.Use(middleware.AuthMiddleware())
	private.SetupPrivateRoutes(privateGroup)

	// Start the Server
	log.Printf("YOUR SERVER IS LIVE AT PORT %s", portStr)
	if err := router.Run(fmt.Sprintf(":%s", portStr)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
