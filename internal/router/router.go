package router

import (
	"github.com/gin-gonic/gin"
	"quocbui.dev/m/internal/config"
	"quocbui.dev/m/internal/handlers"
	"quocbui.dev/m/internal/middleware"
)

func SetupRouter(
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	linkHandler *handlers.LinkHandler,
) *gin.Engine {
	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware(cfg.RateLimit.Requests, cfg.RateLimit.Window))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public shorten (anonymous)
		api.POST("/shorten", linkHandler.ShortenPublic)
	}

	// Protected routes
	protected := r.Group("/api/v1/me")
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		protected.GET("", userHandler.GetMe)
		protected.POST("/shorten", linkHandler.ShortenPrivate)
		protected.GET("/links", linkHandler.GetMyLinks)
		protected.GET("/links/:code", linkHandler.GetMyLinkDetail)
		protected.DELETE("/links/:code", linkHandler.DeleteMyLink)
	}

	// Redirect route (must be last to avoid conflicts)
	r.GET("/:code", linkHandler.Redirect)

	return r
}
