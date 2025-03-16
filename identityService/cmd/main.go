package main

import (
	"fmt"
	"identity/config/db"
	"identity/routes"
	"log"
	"os"

	"identity/badwords"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	db.Connect()
}

func main() {
	// Step 1: Load bad words from a text file
	success, err := badwords.LoadBadWords("badwords/en.txt")
	if !success {
		fmt.Println("Failed to load bad words:", err)
		return
	}
	fmt.Println("Bad words loaded successfully!")

	// Test the ContainsBadWords function
	testInput := "This is a test message with fuck."
	containsBadWords := badwords.ContainsBadWords(testInput)

	if containsBadWords {
		fmt.Println("Bad words detected in the input.")
	} else {
		fmt.Println("No bad words detected in the input.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	r := gin.Default()

	routes.RegisterRoutes(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	log.Printf("Starting server on port %s...", port)
	r.Run(":" + port)
}
