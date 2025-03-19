package auth

import (
	"fmt"
	"identity/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Authorization header required")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		log.Printf("Token authHeader string: %s", authHeader)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("Token string: %s", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// Ensure the token method is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			jwtSecret := utils.GetJWTSecret()

			return []byte(jwtSecret), nil
		})

		if err != nil {
			log.Printf("Error parsing token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			log.Println("Token is not valid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			log.Println("Invalid token claims: user_id not found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		log.Printf("Authenticated user ID: %s", userID)
		c.Set("userID", userID)
		c.Next()
	}
}
