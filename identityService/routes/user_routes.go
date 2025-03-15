package routes

import (
	"identity/controllers"
	"identity/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Create a new user controller
	userController := controllers.NewUserController()

	// Public routes
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.POST("/refresh-token", userController.RefreshToken)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/logout", userController.Logout)
		protected.GET("/user/:username", userController.GetUserByUsername)
	}
}
