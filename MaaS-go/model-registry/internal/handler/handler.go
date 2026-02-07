package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"maas-platform/model-registry/internal/model"
	"maas-platform/model-registry/internal/repository"
	"maas-platform/model-registry/internal/service"
	"maas-platform/model-registry/pkg/logger"
)

// ModelHandler handles model-related HTTP requests
type ModelHandler struct {
	service service.ModelService
	logger  *logger.Logger
}

// NewModelHandler creates a new model handler
func NewModelHandler(s service.ModelService, logger *logger.Logger) *ModelHandler {
	return &ModelHandler{
		service: s,
		logger:  logger,
	}
}

// CreateModelRequest represents a model creation request
type CreateModelRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Version     string            `json:"version" binding:"required"`
	Framework   string            `json:"framework" binding:"required"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
	IsPublic    bool              `json:"is_public"`
	OwnerID     string            `json:"owner_id" binding:"required"`
	TenantID    string            `json:"tenant_id" binding:"required"`
}

// CreateModel handles model creation
func (h *ModelHandler) CreateModel(c *gin.Context) {
	var req CreateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createReq := service.CreateModelRequest{
		Name:        req.Name,
		Description: req.Description,
		Version:     req.Version,
		Framework:   model.ModelFramework(req.Framework),
		Tags:        req.Tags,
		Metadata:    req.Metadata,
		OwnerID:     req.OwnerID,
		TenantID:    req.TenantID,
		IsPublic:    req.IsPublic,
	}

	m, err := h.service.CreateModel(c.Request.Context(), createReq)
	if err != nil {
		h.logger.Error("Failed to create model", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, m)
}

// GetModel handles getting a model by ID
func (h *ModelHandler) GetModel(c *gin.Context) {
	id := c.Param("id")

	m, err := h.service.GetModel(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrModelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "model not found"})
			return
		}
		h.logger.Error("Failed to get model", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, m)
}

// ListModels handles listing models
func (h *ModelHandler) ListModels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := service.ListModelsFilter{
		Name:  c.Query("name"),
		Page:  page,
		Limit: limit,
	}

	// Parse optional filters
	if framework := c.Query("framework"); framework != "" {
		filter.Framework = model.ModelFramework(framework)
	}
	if status := c.Query("status"); status != "" {
		filter.Status = model.ModelStatus(status)
	}

	response, err := h.service.ListModels(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to list models", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateModel handles updating a model
func (h *ModelHandler) UpdateModel(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m, err := h.service.UpdateModel(c.Request.Context(), id, req)
	if err != nil {
		if err == repository.ErrModelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "model not found"})
			return
		}
		h.logger.Error("Failed to update model", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, m)
}

// DeleteModel handles deleting a model
func (h *ModelHandler) DeleteModel(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteModel(c.Request.Context(), id); err != nil {
		if err == repository.ErrModelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "model not found"})
			return
		}
		h.logger.Error("Failed to delete model", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateModelStatus handles updating model status
func (h *ModelHandler) UpdateModelStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateModelStatus(c.Request.Context(), id, model.ModelStatus(req.Status)); err != nil {
		if err == repository.ErrModelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "model not found"})
			return
		}
		h.logger.Error("Failed to update model status", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
