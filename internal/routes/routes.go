package routes

import (
	"ai-mockup/internal/handlers"
	"ai-mockup/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine) {
	// Add logging middleware
	router.Use(middleware.LoggingMiddleware())

	// API routes
	api := router.Group("/api")
	{
		api.GET("/ping", handlers.PingHandler)
		api.GET("/health", handlers.HealthHandler)
		api.GET("/session", middleware.AuthMiddleware(), handlers.SessionHandler)
		api.GET("/usage", middleware.AuthMiddleware(), handlers.GetUsageHandler)
		api.POST("/generate", middleware.AuthMiddleware(), handlers.GenerateImageHandler)
	}
	auth := api.Group("/auth")
	{
		auth.POST("/logout", handlers.LogoutHandler)
		auth.GET("/google", handlers.OautHandler)
		auth.GET("/google/callback", handlers.OauthCallbackHandler)
	}

}
