package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"word-filter/badwords"
)

func main() {
	// Set up Gin router
	router := gin.Default()

	// Step 1: Load bad words from a text file
	success, err := badwords.LoadBadWords("badwords/en.txt")
	if !success {
		fmt.Println("Failed to load bad words:", err)
		return
	}
	fmt.Println("Bad words loaded successfully!")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Keep only the essential endpoint for checking bad words
	router.POST("/check", func(c *gin.Context) {
		var request badwords.BadWordRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Use the updated CheckText function to check for bad words
		response := badwords.CheckText(request.Text)
		c.JSON(http.StatusOK, response)
	})

	// Health check endpoint (keeping this as it's a good practice)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Service is running"})
	})

	// Start the Gin server directly
	serverAddr := ":" + port
	log.Println("Starting HTTP server on " + serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
