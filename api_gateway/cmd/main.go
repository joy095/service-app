package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/joy095/api-gateway/logger"
	middleware "github.com/joy095/api-gateway/middlewares/cors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	logger.InitLoggers()

	if err := godotenv.Load(".env.local"); err != nil {
		logger.ErrorLogger.Error("Error loading .env.local file")
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
		logger.InfoLogger.Info("PORT not set. Using default: 8080")
	}

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

	// Ensure the URL ends with a slash for proper path joining
	if !strings.HasSuffix(identityService, "/") {
		identityService += "/"
	}

	router.Any("/v1/auth/*proxyPath", func(c *gin.Context) {
		proxyPath := c.Param("proxyPath")
		logger.InfoLogger.Info("Identity Service Called: " + proxyPath)
		reverseProxy(identityService, proxyPath)(c)
	})

	// Image Service
	imageService := os.Getenv("IMAGE_SERVICE_URL")
	if imageService == "" {
		log.Fatal("IMAGE_SERVICE_URL environment variable is required")
	}

	// Ensure the URL ends with a slash for proper path joining
	if !strings.HasSuffix(imageService, "/") {
		imageService += "/"
	}

	router.Any("/v1/image/*proxyPath", func(c *gin.Context) {
		proxyPath := c.Param("proxyPath")
		logger.InfoLogger.Info("Image Service Called: " + proxyPath)
		reverseProxy(imageService, proxyPath)(c)
	})

	logger.InfoLogger.Info("Starting HTTP server on port " + port)
	log.Printf("Starting server on port %s...", port)

	if err := router.Run(":" + port); err != nil {
		logger.ErrorLogger.Error("Failed to start server: " + err.Error())
		log.Fatal(err)
	}
}

// reverseProxy creates a reverse proxy for the given target URL and path
func reverseProxy(target string, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL, err := url.Parse(target)
		if err != nil {
			logger.ErrorLogger.Error("Failed to parse target URL: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse target URL"})
			return
		}

		// Create a new director function to modify the request
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Store the original director
		originalDirector := proxy.Director

		// Create a new director that preserves query parameters and headers
		proxy.Director = func(req *http.Request) {
			originalDirector(req)

			// Remove the prefix from the path
			req.URL.Path = path

			// Preserve query parameters
			req.URL.RawQuery = c.Request.URL.RawQuery

			// Log the final target URL for debugging
			logger.InfoLogger.Info("Proxying to: " + req.URL.String())
		}

		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			logger.ErrorLogger.Error("Proxy error: " + err.Error())
			c.JSON(http.StatusBadGateway, gin.H{"error": "Proxy request failed"})
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
