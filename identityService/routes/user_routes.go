package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joy095/identity/controllers"
	middleware "github.com/joy095/identity/middlewares"
	"github.com/joy095/identity/middlewares/auth"
	"github.com/joy095/identity/utils/mail"
)

func RegisterRoutes(router *gin.Engine) {
	userController := controllers.NewUserController()

	// Use the rate limiter globally

	router.POST("/register", middleware.NewRateLimiter("3-M"), userController.Register)          // Limit registration abuse
	router.POST("/login", middleware.NewRateLimiter("10-M"), userController.Login)               // Allow retries but prevent brute force
	router.POST("/refresh-token", middleware.NewRateLimiter("5-M"), userController.RefreshToken) // Allow token refresh but not spam

	router.POST("/request-otp", middleware.NewRateLimiter("2-M"), mail.RequestOTP) // Prevent OTP spamming
	router.POST("/verify-otp", middleware.NewRateLimiter("5-M"), mail.VerifyOTP)   // Allow some retry margin

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.POST("/logout", middleware.NewRateLimiter("10-M"), userController.Logout)                   // Should be safe
		protected.GET("/user/:username", middleware.NewRateLimiter("30-M"), userController.GetUserByUsername) // Higher since it's a read-only action
	}

}
