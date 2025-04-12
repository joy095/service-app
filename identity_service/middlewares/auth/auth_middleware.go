package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/joy095/identity/logger"
	"github.com/joy095/identity/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	logger.InfoLogger.Info("AuthMiddleware called")

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {

			logger.ErrorLogger.Error("Authorization header required")
			log.Println("Authorization header required")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		logger.ErrorLogger.Errorf("Token authHeader string: %s", authHeader)
		log.Printf("Token authHeader string: %s", authHeader)

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		logger.ErrorLogger.Errorf("Token string: %s", tokenString)
		log.Printf("Token string: %s", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// Ensure the token method is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.ErrorLogger.Errorf("Unexpected signing method: %v", token.Header["alg"])
				log.Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			jwtSecret := utils.GetJWTSecret()

			return []byte(jwtSecret), nil
		})

		if err != nil {
			logger.ErrorLogger.Errorf("Error passing token: %v", err)
			log.Printf("Error parsing token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			logger.ErrorLogger.Error("Token is not valid")
			log.Println("Token is not valid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.ErrorLogger.Error("Invalid token claims")
			log.Println("Invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			logger.ErrorLogger.Error("Invalid token claims: user_id not found")
			log.Println("Invalid token claims: user_id not found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		logger.InfoLogger.Infof("Authenticated user ID: %s", userID)
		log.Printf("Authenticated user ID: %s", userID)
		c.Set("user_id", userID)
		c.Next()
	}
}
