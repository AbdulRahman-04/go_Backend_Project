package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"Go_Backend/config"
	"Go_Backend/utils"
	"Go_Backend/controllers/public"
)

func main() {
	// ✅ Step 1: Initialize Gin Router
	router := gin.Default()

	// ✅ Step 2: Set Release Mode (Production)
	gin.SetMode(gin.ReleaseMode)

	// ✅ Step 3: Load Configuration
	cfg := config.LoadConfig()
	portStr := strconv.Itoa(cfg.Port)

	// ✅ Step 4: Start Database Connection FIRST
	log.Println("Connecting to MongoDB... 🔄") // Debug log
	err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err) // Stop execution if connection fails
	}
	log.Println("✅ DATABASE CONNECTED SUCCESSFULLY!")

	// ✅ Step 5: Health Check Route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "HELLO IN TS💙"})
	})

	// ✅ Step 6: Load Public API Routes
	public.SetupPublicRoutes(router)

	// ✅ Step 7: Start the Server
	log.Printf("YOUR SERVER IS LIVE AT PORT %s", portStr)
	if err := router.Run(fmt.Sprintf(":%s", portStr)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}