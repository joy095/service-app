package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// CorsMiddleware sets up CORS settings
func CorsMiddleware() gin.HandlerFunc {
	godotenv.Load()

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	return cors.New(cors.Config{
		AllowOrigins:     strings.Split(allowedOrigins, ","), // Convert CSV to slice
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
