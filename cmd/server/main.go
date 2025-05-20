package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"Go_Backend/config"
	privateCtrl "Go_Backend/controllers/private" // aliased to avoid conflicts
	"Go_Backend/controllers/public"
	"Go_Backend/middleware"
	"Go_Backend/utils"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB (assumed to be implemented in utils.ConnectDB)
	if err := utils.ConnectDB(); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	fmt.Println("MongoDB connected")

	// (Optional) Initialize Redis if needed for other purposes.
	utils.InitRedis()

	// Create a new Gin router and add global middlewares
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Setup public routes with in-memory rate limiter for "public"
	publicGroup := router.Group("/api/public")
	publicGroup.Use(middleware.RateLimitMiddlewareInMemory("public"))
	public.SetupPublicRoutes(publicGroup)

	// Setup private routes with AuthMiddleware and in-memory rate limiter for "private"
	privateGroup := router.Group("/api/private")
	privateGroup.Use(middleware.AuthMiddleware())
	privateGroup.Use(middleware.RateLimitMiddlewareInMemory("private"))
	privateCtrl.SetupPrivateRoutes(privateGroup)

	// Determine port (default: 8080 if cfg.Port is 0)
	portStr := "8080" // Default fallback
    if cfg.Port != 0 {
    portStr = strconv.Itoa(cfg.Port)
    }
	log.Println("Server running on port", portStr)
	if err := router.Run(":" + portStr); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}