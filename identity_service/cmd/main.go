package main

import (
	"log"
	"os"

	"github.com/joy095/identity/routes"

	"github.com/joy095/identity/config/db"
	"github.com/joy095/identity/logger"
	middleware "github.com/joy095/identity/middlewares/cors"
	logger_middleware "github.com/joy095/identity/middlewares/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Initialize loggers before using
	logger.InitLoggers()

	godotenv.Load(".env.local")
	db.Connect()
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	r := gin.Default()

	r.Use(middleware.CorsMiddleware())

	// Apply Logger Middleware
	r.Use(logger_middleware.GinLogger())

	routes.RegisterRoutes(r)

	r.GET("/health", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	logger.InfoLogger.Info("Server is started")

	log.Printf("Starting server on port %s...", port)

	r.Run(":" + port)
}
