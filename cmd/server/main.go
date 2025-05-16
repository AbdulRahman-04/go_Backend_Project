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
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	if err := utils.ConnectDB(); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	fmt.Println("MongoDB connected")

	// Initialize Redis
	utils.InitRedis()

	// Create Gin router
	router := gin.New()

	// Global Middlewares
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Setup public routes with "public" rate limiter middleware
	publicGroup := router.Group("/api/public")
	publicGroup.Use(middleware.RateLimitMiddleware("public"))
	public.SetupPublicRoutes(publicGroup)

	// Setup private routes with Auth and "private" rate limiter middleware
	privateGroup := router.Group("/api/private")
	privateGroup.Use(middleware.AuthMiddleware())
	privateGroup.Use(middleware.RateLimitMiddleware("private"))
	private.SetupPrivateRoutes(privateGroup)

	// Determine port (default is 8080 if not provided)
	portStr := "8080"
	if cfg.Port != 0 {
		portStr = strconv.Itoa(cfg.Port)
	}
	log.Println("Server running on port", portStr)
	if err := router.Run(":" + portStr); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}