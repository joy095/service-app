package controllers

import (
	"identity/config/db"
	"identity/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	c.JSON(http.StatusCreated, gin.H{
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
