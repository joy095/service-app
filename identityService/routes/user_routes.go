package routes

import (
	"identity/controllers"
	"identity/middlewares"
	"identity/utils/mail"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Create a new user controller
	userController := controllers.NewUserController()

	// Public routes
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.POST("/refresh-token", userController.RefreshToken)

	// OTP verification and password reset routes
	router.POST("/send-otp", mail.RequestOTP)
	router.POST("/verify-otp", mail.VerifyOTP)
	// router.POST("/reset-password", userController.ResetPassword)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/logout", userController.Logout)
		protected.GET("/user/:username", userController.GetUserByUsername)
	}
}
