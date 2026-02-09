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
			models.PUT("/:id", h.UpdateModel)
			models.DELETE("/:id", h.DeleteModel)
			models.PATCH("/:id/status", h.UpdateModelStatus)
			models.POST("/:id/tags", h.AddModelTags)
			models.DELETE("/:id/tags", h.RemoveModelTags)
			models.GET("/:id/metadata", h.GetModelMetadata)
			models.PUT("/:id/metadata", h.SetModelMetadata)
		}

		// Inference routes
		protected.POST("/inference", h.RunInference)
	}
}
