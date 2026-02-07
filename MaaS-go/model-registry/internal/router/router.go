package router

import (
	"github.com/gin-gonic/gin"

	"maas-platform/model-registry/internal/handler"
)

// RegisterRoutes registers all routes
func RegisterRoutes(r *gin.RouterGroup, h *handler.ModelHandler) {
	// Model routes
	models := r.Group("/models")
	{
		models.POST("", h.CreateModel)
		models.GET("", h.ListModels)
		models.GET("/:id", h.GetModel)
		models.PUT("/:id", h.UpdateModel)
		models.DELETE("/:id", h.DeleteModel)
		models.PATCH("/:id/status", h.UpdateModelStatus)
	}
}
