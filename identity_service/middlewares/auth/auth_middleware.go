package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/joy095/identity/config/db"
	"github.com/joy095/identity/logger"
	"github.com/joy095/identity/models"
	"github.com/joy095/identity/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.ErrorLogger.Error("Authorization header required")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(utils.GetJWTSecret()), nil
		})
		if err != nil || !token.Valid {
			logger.ErrorLogger.Errorf("Invalid token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.ErrorLogger.Error("Invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userIDFromToken, ok := claims["user_id"].(string)
		if !ok {
			logger.ErrorLogger.Error("Token does not contain user_id")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		usernameParam := c.Param("username")
		// Check if the user_id from the token matches the username parameter
		rawBody, _ := c.GetRawData()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody)) // allow re-reading

		var body struct {
			UserID string `json:"user_id"`
		}
		json.Unmarshal(rawBody, &body)

		if usernameParam == "" && body.UserID == "" {
			logger.ErrorLogger.Error("Either 'username' param or 'user_id' in body is required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Either 'username' param or 'user_id' in body is required"})
			c.Abort()
			return
		}

		// If username is provided, fetch user and match token user_id
		if usernameParam != "" {
			user, err := models.GetUserByUsername(db.DB, usernameParam)
			if err != nil {
				logger.ErrorLogger.Errorf("User not found by username: %s", usernameParam)
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				c.Abort()
				return
			}

			if user.ID.String() != userIDFromToken {
				logger.ErrorLogger.Errorf("User ID mismatch: token(%s) vs db(%s)", userIDFromToken, user.ID)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
				c.Abort()
				return
			}

			logger.InfoLogger.Infof("Authenticated via username: %s", user.Username)
		}

		// If user_id from body is provided, ensure it matches token
		if body.UserID != "" && body.UserID != userIDFromToken {
			logger.ErrorLogger.Errorf("User ID mismatch: token(%s) vs body(%s)", userIDFromToken, body.UserID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			c.Abort()
			return
		}

		// Attach user_id to context for downstream handlers
		c.Set("user_id", userIDFromToken)
		logger.InfoLogger.Infof("Authenticated user_id: %s", userIDFromToken)
		c.Next()
	}
}
