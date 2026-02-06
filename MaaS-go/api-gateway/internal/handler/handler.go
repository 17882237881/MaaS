package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"maas-platform/api-gateway/internal/config"
	"maas-platform/api-gateway/pkg/logger"
)

// Handler handles HTTP requests
type Handler struct {
	config *config.Config
	logger *logger.Logger
}

// New creates a new handler
func New(cfg *config.Config, log *logger.Logger) *Handler {
	return &Handler{
		config: cfg,
		logger: log,
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
