package router

import (
	"github.com/gin-gonic/gin"

	"maas-platform/api-gateway/internal/handler"
)

// RegisterRoutes registers all routes
func RegisterRoutes(r *gin.RouterGroup, h *handler.Handler) {
	// Auth routes (no authentication required)
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}

	// Protected routes
	protected := r.Group("")
	{
		// User routes
		users := protected.Group("/users")
		{
			users.GET("/me", h.GetCurrentUser)
		}

		// Model routes
		models := protected.Group("/models")
		{
			models.POST("", h.CreateModel)
			models.GET("", h.ListModels)
			models.GET("/:id", h.GetModel)
			models.DELETE("/:id", h.DeleteModel)
		}

		// Inference routes
		protected.POST("/inference", h.RunInference)
	}
}
