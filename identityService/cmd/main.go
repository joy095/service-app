package main

import (
	"identity/config/db"
	"identity/logger"
	logger_middleware "identity/middlewares/logger"
	"identity/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	db.Connect()
	logger.InitLoggers()
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	r := gin.Default()

	// Apply Logger Middleware
	r.Use(logger_middleware.GinLogger())

	routes.RegisterRoutes(r)

	r.GET("/health", func(c *gin.Context) {
		logger.InfoLogger.Info("Server is healthy")
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	log.Printf("Starting server on port %s...", port)
	r.Run(":" + port)
}
