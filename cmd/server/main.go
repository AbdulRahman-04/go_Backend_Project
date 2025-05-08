package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"Go_Backend/config"
	"Go_Backend/utils"
)

func main() {
	// ✅ Step 1: Initialize Gin Router
	router := gin.Default()

	// ✅ Step 2: Set Release Mode (Production)
	gin.SetMode(gin.ReleaseMode)

	// ✅ Step 3: Print Server Start Message FIRST
	cfg := config.LoadConfig()
	portStr := strconv.Itoa(cfg.Port)
	log.Printf("YOUR SERVER IS LIVE AT PORT %s", portStr)

	// ✅ Step 4: Start Database Connection AFTER Server Message
	_, err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	log.Println("DATABASE CONNECTED SUCCESSFULLY! ✅") // ✅ Now this prints after server message

	// ✅ Step 5: Health Check Route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "HELLO IN TS💙"})
	})

	// ✅ Step 6: Start the Server
	if err := router.Run(fmt.Sprintf(":%s", portStr)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}