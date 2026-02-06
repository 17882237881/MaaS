package handler

import (
	"github.com/gin-gonic/gin"
)

// InferenceRequest represents an inference request
type InferenceRequest struct {
	ModelID string                 `json:"model_id" binding:"required"`
	Input   map[string]interface{} `json:"input" binding:"required"`
}

// InferenceResponse represents an inference response
type InferenceResponse struct {
	ModelID   string                 `json:"model_id"`
	Output    map[string]interface{} `json:"output"`
	Latency   int64                  `json:"latency_ms"`
	RequestID string                 `json:"request_id"`
}

// RunInference runs model inference
func (h *Handler) RunInference(c *gin.Context) {
	var req InferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err.Error())
		return
	}

	// TODO: Call inference service
	h.Success(c, InferenceResponse{
		ModelID:   req.ModelID,
		Output:    map[string]interface{}{"prediction": "result"},
		Latency:   100,
		RequestID: c.GetString("request_id"),
	})
}
