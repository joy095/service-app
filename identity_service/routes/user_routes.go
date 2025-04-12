package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joy095/identity/controllers"
	"github.com/joy095/identity/controllers/relations"
	middleware "github.com/joy095/identity/middlewares"
	"github.com/joy095/identity/middlewares/auth"
	"github.com/joy095/identity/utils/mail"
)

func RegisterRoutes(router *gin.Engine) {
	userController := controllers.NewUserController()
	relationController := relations.NewRelationController()

	// Public routes
	router.POST("/register", middleware.NewRateLimiter("3-M"), userController.Register)
	router.POST("/login", middleware.NewRateLimiter("10-M"), userController.Login)
	router.POST("/refresh-token", middleware.NewRateLimiter("5-M"), userController.RefreshToken)

	router.POST("/request-otp", middleware.NewRateLimiter("2-M"), mail.RequestOTP)
	router.POST("/verify-otp", middleware.NewRateLimiter("5-M"), mail.VerifyOTP)

	// Protected routes
	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.POST("/logout", middleware.NewRateLimiter("10-M"), userController.Logout)
		protected.GET("/user/:username", middleware.NewRateLimiter("30-M"), userController.GetUserByUsername)

		// Relationship routes
		protected.POST("/relation/request", middleware.NewRateLimiter("10-M"), relationController.SendRequest)
		protected.POST("/relation/accept", middleware.NewRateLimiter("30-M"), relationController.AcceptRequest)
		protected.POST("/relation/reject", middleware.NewRateLimiter("30-M"), relationController.RejectRequest)
		protected.GET("/relation/pending", middleware.NewRateLimiter("20-M"), relationController.ListPendingRequests)
		protected.GET("/relation/connections", middleware.NewRateLimiter("30-M"), relationController.ListConnections)
		protected.GET("/relation/status/:user_id", middleware.NewRateLimiter("30-M"), relationController.CheckConnectionStatus)
	}
}
