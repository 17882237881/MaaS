package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"maas-platform/api-gateway/internal/config"
	"maas-platform/api-gateway/internal/service"
	"maas-platform/api-gateway/pkg/logger"
	modelpb "maas-platform/shared/proto/model"
)

// Handler handles HTTP requests
type Handler struct {
	config      *config.Config
	logger      *logger.Logger
	modelClient *service.ModelServiceClient
}

// New creates a new handler
func New(cfg *config.Config, log *logger.Logger, modelClient *service.ModelServiceClient) *Handler {
	return &Handler{
		config:      cfg,
		logger:      log,
		modelClient: modelClient,
	}
}

// Response represents a standard API response
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// Success returns a successful response
func (h *Handler) Success(c *gin.Context, data interface{}) {
	requestID := c.GetString("request_id")
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		RequestID: requestID,
	})
}

// Error returns an error response
func (h *Handler) Error(c *gin.Context, code int, message string) {
	requestID := c.GetString("request_id")
	c.JSON(code, Response{
		Code:      code,
		Message:   message,
		RequestID: requestID,
	})
}

// BadRequest returns a 400 error
func (h *Handler) BadRequest(c *gin.Context, message string) {
	h.Error(c, http.StatusBadRequest, message)
}

// Unauthorized returns a 401 error
func (h *Handler) Unauthorized(c *gin.Context) {
	h.Error(c, http.StatusUnauthorized, "unauthorized")
}

// Forbidden returns a 403 error
func (h *Handler) Forbidden(c *gin.Context) {
	h.Error(c, http.StatusForbidden, "forbidden")
}

// NotFound returns a 404 error
func (h *Handler) NotFound(c *gin.Context, resource string) {
	h.Error(c, http.StatusNotFound, resource+" not found")
}

// InternalError returns a 500 error
func (h *Handler) InternalError(c *gin.Context, err error) {
	h.logger.Error("Internal error", "error", err)
	h.Error(c, http.StatusInternalServerError, "internal server error")
}

// ModelRequest represents a model upload request
type ModelRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Version     string            `json:"version" binding:"required"`
	Framework   string            `json:"framework" binding:"required"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
}

// ModelResponse represents a model response
type ModelResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	Framework   string            `json:"framework"`
	Status      string            `json:"status"`
	Size        int64             `json:"size"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

// CreateModel creates a new model via gRPC
func (h *Handler) CreateModel(c *gin.Context) {
	var req ModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err.Error())
		return
	}

	// Get user info from context
	ownerID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")

	ownerIDStr, _ := ownerID.(string)
	if ownerIDStr == "" {
		ownerIDStr = "default-user"
	}
	tenantIDStr, _ := tenantID.(string)
	if tenantIDStr == "" {
		tenantIDStr = "default-tenant"
	}

	// Call Model Registry via gRPC
	grpcReq := &modelpb.CreateModelRequest{
		Name:        req.Name,
		Description: req.Description,
		Version:     req.Version,
		Framework:   req.Framework,
		Tags:        req.Tags,
		Metadata:    req.Metadata,
		OwnerId:     ownerIDStr,
		TenantId:    tenantIDStr,
		IsPublic:    false,
	}

	model, err := h.modelClient.CreateModel(c.Request.Context(), grpcReq)
	if err != nil {
		h.InternalError(c, err)
		return
	}

	h.Success(c, convertProtoModelToResponse(model))
}

// ListModels lists all models via gRPC
func (h *Handler) ListModels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	ownerID, _ := c.Get("user_id")
	ownerIDStr, _ := ownerID.(string)

	grpcReq := &modelpb.ListModelsRequest{
		Page:    int32(page),
		Limit:   int32(limit),
		OwnerId: ownerIDStr,
	}

	if framework := c.Query("framework"); framework != "" {
		grpcReq.Framework = framework
	}
	if status := c.Query("status"); status != "" {
		grpcReq.Status = status
	}

	models, total, err := h.modelClient.ListModels(c.Request.Context(), grpcReq)
	if err != nil {
		h.InternalError(c, err)
		return
	}

	response := make([]ModelResponse, len(models))
	for i, m := range models {
		response[i] = convertProtoModelToResponse(m)
	}

	h.Success(c, gin.H{
		"models": response,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// GetModel gets a model by ID via gRPC
func (h *Handler) GetModel(c *gin.Context) {
	id := c.Param("id")

	model, err := h.modelClient.GetModel(c.Request.Context(), id)
	if err != nil {
		h.InternalError(c, err)
		return
	}

	h.Success(c, convertProtoModelToResponse(model))
}

// DeleteModel deletes a model via gRPC
func (h *Handler) DeleteModel(c *gin.Context) {
	id := c.Param("id")

	err := h.modelClient.DeleteModel(c.Request.Context(), id)
	if err != nil {
		h.InternalError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// convertProtoModelToResponse converts protobuf Model to HTTP response
func convertProtoModelToResponse(m *modelpb.Model) ModelResponse {
	return ModelResponse{
		ID:          m.Id,
		Name:        m.Name,
		Description: m.Description,
		Version:     m.Version,
		Framework:   m.Framework,
		Status:      m.Status,
		Size:        m.Size,
		Tags:        m.Tags,
		Metadata:    make(map[string]string),
		CreatedAt:   m.CreatedAt.AsTime().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   m.UpdatedAt.AsTime().Format("2006-01-02T15:04:05Z"),
	}
}
