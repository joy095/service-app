package controllers

import (
	"context"
	"errors"
	"identity/config/db"
	"identity/models"
	"net/http"
	"time"

	"identity/utils/mail"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// UserController handles user-related requests
type UserController struct{}

// NewUserController creates a new UserController
func NewUserController() *UserController {
	return &UserController{}
}

// Register handles user registration
func (uc *UserController) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, accessToken, refreshToken, err := models.CreateUser(db.DB, req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	otp := mail.GenerateSecureOTP()
	mail.SendOTP(req.Email, otp)

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"otp":      otp,
		},
		"tokens": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// Login handles user login
func (uc *UserController) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, accessToken, refreshToken, err := models.LoginUser(db.DB, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
		"tokens": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (uc *UserController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Query the database to find the user with this refresh token
	var user models.User
	query := `SELECT id, username, email, refresh_token FROM users WHERE refresh_token = $1`
	err := db.DB.QueryRow(context.Background(), query, req.RefreshToken).Scan(
		&user.ID, &user.Username, &user.Email, &user.RefreshToken,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Token not found in database
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		} else {
			// Database error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Generate a new access token
	accessToken, err := models.GenerateAccessToken(user.ID, time.Minute*15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Generate a new refresh token (optional, for token rotation)
	newRefreshToken, err := models.GenerateRefreshToken(user.ID, time.Hour*24*7) // Stronger Refresh Token for 7 days
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Update the refresh token in the database
	_, err = db.DB.Exec(context.Background(), `UPDATE users SET refresh_token = $1 WHERE id = $2`, newRefreshToken, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
		return
	}

	// Return the new tokens
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Logout handles user logout
func (uc *UserController) Logout(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required,uuid"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format or missing fields", "error-message": err.Error()})
		return
	}

	// Get the user ID from the context
	userIDFromToken, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Ensure the user can only log out their own account
	if userIDFromToken != req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only log out your own account"})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	if err := models.LogoutUser(db.DB, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// GetUserByUsername retrieves a user by username
func (uc *UserController) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := models.GetUserByUsername(db.DB, username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
