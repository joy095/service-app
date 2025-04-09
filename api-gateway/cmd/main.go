package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/joy095/api-gateway/logger"
	middleware "github.com/joy095/api-gateway/middlewares/cors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	logger.InitLoggers()

	godotenv.Load(".env.local")
}

func main() {
	port := os.Getenv("PORT")

	router := gin.Default()

	router.Use(middleware.CorsMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	// Identity Service
	identityService := os.Getenv("IDENTITY_SERVICE_URL")
	if identityService == "" {
		log.Fatal("IDENTITY_SERVICE_URL environment variable is required")
	}

	logger.InfoLogger.Info("Identity Service Called")

	router.Any("/v1/auth/*proxyPath", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("proxyPath")
		reverseProxy(identityService)(c)
	})

	// Image Service
	imageService := os.Getenv("IMAGE_SERVICE_URL")
	if imageService == "" {
		log.Fatal("IMAGE_SERVICE_URL environment variable is required")
	}

	logger.InfoLogger.Info("Image Service Called")

	router.Any("/v1/image/*proxyPath", func(c *gin.Context) {

		c.Request.URL.Path = c.Param("proxyPath")
		reverseProxy(imageService)(c)
	})

	logger.InfoLogger.Info("Starting HTTP server on %s..." + port)

	log.Printf("Starting server on port %s...", port)

	router.Run(":" + port)
}

// reverseProxy creates a reverse proxy for the given target URL
func reverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL, err := url.Parse(target)
		if err != nil {
			logger.ErrorLogger.Error("Failed to parse target URL:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse target URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
