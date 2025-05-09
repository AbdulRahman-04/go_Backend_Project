package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"Go_Backend/config"
	"Go_Backend/utils"
	"Go_Backend/controllers/public"
	"Go_Backend/controllers/private"
	"Go_Backend/middleware"
)

func main() {
	// Initialize Gin Router
	router := gin.Default()

	// Set Release Mode (Production)
	gin.SetMode(gin.ReleaseMode)

	// Load Configuration
	cfg := config.LoadConfig()
	portStr := strconv.Itoa(cfg.Port)

	// Connect to MongoDB first
	log.Println("Connecting to MongoDB... 🔄")
	err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	log.Println("✅ DATABASE CONNECTED SUCCESSFULLY!")

	// Health Check Route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "HELLO IN TS💙"})
	})

	// Public API Routes
	public.SetupPublicRoutes(router)

	// Private Routes with JWT Middleware
	privateRoutes := router.Group("/api/private")
	privateRoutes.Use(middleware.AuthMiddleware()) // JWT middleware applied here
	{
		privateRoutes.POST("/addtodo", private.AddTodo)
		privateRoutes.GET("/alltodos", private.GetAllTodos)
		privateRoutes.GET("/getone/:id", private.GetOneTodo)
		privateRoutes.PUT("/editone/:id", private.EditTodo)
		privateRoutes.DELETE("/deleteone/:id", private.DeleteTodo)
		privateRoutes.DELETE("/deleteall", private.DeleteAllTodos)
	}

	// Start the Server
	log.Printf("YOUR SERVER IS LIVE AT PORT %s", portStr)
	if err := router.Run(fmt.Sprintf(":%s", portStr)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}