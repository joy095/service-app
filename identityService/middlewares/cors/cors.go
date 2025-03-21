package middleware

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// CorsMiddleware sets up CORS settings
func CorsMiddleware() gin.HandlerFunc {
	godotenv.Load()

	identityService := os.Getenv("IDENTITY_SERVICE_URL")
	imageService := os.Getenv("IMAGE_SERVICE_URL")

	return cors.New(cors.Config{
		AllowOrigins:     []string{identityService, imageService},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
