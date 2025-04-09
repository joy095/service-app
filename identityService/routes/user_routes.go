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
	router.Use(middleware.RateLimiterMiddleware())

	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.POST("/refresh-token", userController.RefreshToken)
	router.POST("/request-otp", mail.RequestOTP)
	router.POST("/verify-otp", mail.VerifyOTP)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.POST("/logout", userController.Logout)
		protected.GET("/user/:username", userController.GetUserByUsername)
	}
}
