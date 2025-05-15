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
	// Load Config
	cfg := config.LoadConfig()
	portStr := strconv.Itoa(cfg.Port)

	// Connect to MongoDB
	err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}

	// Connect to Redis
	utils.InitRedis()

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Global Middlewares (excluding RateLimitMiddleware)
	router.Use(gin.Recovery())                 // Recover from panic
	router.Use(middleware.LoggerMiddleware())  // Logging
	router.Use(middleware.CORSMiddleware())    // CORS

	// Health Check
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "HELLO IN TS💙"})
	})

	// Public Routes
	public.SetupPublicRoutes(router)

	// Private Routes with JWT Auth
	privateGroup := router.Group("/api/private")
	privateGroup.Use(middleware.AuthMiddleware()) // JWT Auth
	private.SetupPrivateRoutes(privateGroup)       // Here we apply RateLimit per route inside

	// Server Start Log First
	log.Printf("🚀 SERVER IS LIVE AT PORT %s", portStr)
	log.Println("✅ DATABASE CONNECTED SUCCESSFULLY!")

	// Start Server
	if err := router.Run(fmt.Sprintf(":%s", portStr)); err != nil {
		log.Fatalf("❌ Failed to run server: %v", err)
	}
}
