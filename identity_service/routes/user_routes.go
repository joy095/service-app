package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joy095/identity/controllers"
	"github.com/joy095/identity/controllers/relations"
	"github.com/joy095/identity/middlewares/auth"
	"github.com/joy095/identity/utils/mail"
)

func RegisterRoutes(router *gin.Engine) {
	userController := controllers.NewUserController()
	relationController := relations.NewRelationController()

	// Public routes
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.POST("/refresh-token", userController.RefreshToken)

	router.POST("/request-otp", mail.RequestOTP)
	router.POST("/verify-otp", mail.VerifyOTP)

	// Protected routes
	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.POST("/logout", userController.Logout)

		protected.GET("/user/:username", userController.GetUserByUsername)

	}

	// Relationship routes using JWT by default to send accept and reject requests
	router.POST("/relation/request", relationController.SendRequest)
	router.POST("/relation/accept", relationController.AcceptRequest)
	router.POST("/relation/reject", relationController.RejectRequest)
	router.GET("/relation/pending", relationController.ListPendingRequests)
	router.GET("/relation/connections", relationController.ListConnections)
	router.GET("/relation/status/:username", relationController.CheckConnectionStatus)

	// // Public routes
	// router.POST("/register", middleware.CombinedRateLimiter("5-5m", "20-2h"), userController.Register)
	// router.POST("/login", middleware.CombinedRateLimiter("25-5m", "20-2h"), userController.Login)
	// router.POST("/refresh-token", middleware.CombinedRateLimiter("10-15m", "30-2h"), userController.RefreshToken)

	// router.POST("/request-otp", middleware.CombinedRateLimiter("5-5m", "20-2h"), mail.RequestOTP)
	// router.POST("/verify-otp", middleware.CombinedRateLimiter("25-5m", "20-2h"), mail.VerifyOTP)

	// // Protected routes
	// protected := router.Group("/")
	// protected.Use(auth.AuthMiddleware())
	// {
	// 	protected.POST("/logout", middleware.NewRateLimiter("10-5m"), userController.Logout)
	// 	protected.GET("/user/:username", middleware.NewRateLimiter("30-1m"), userController.GetUserByUsername)

	// 	// Relationship routes
	// 	protected.POST("/relation/request", middleware.NewRateLimiter("30-3m"), relationController.SendRequest)
	// 	protected.POST("/relation/accept", middleware.NewRateLimiter("30-1m"), relationController.AcceptRequest)
	// 	protected.POST("/relation/reject", middleware.NewRateLimiter("30-1m"), relationController.RejectRequest)
	// 	protected.GET("/relation/pending", middleware.NewRateLimiter("30-1m"), relationController.ListPendingRequests)
	// 	protected.GET("/relation/connections", middleware.NewRateLimiter("30-1m"), relationController.ListConnections)
	// 	protected.GET("/relation/status/:user_id", middleware.NewRateLimiter("30-1m"), relationController.CheckConnectionStatus)
	// }
}
