package handler

import (

	"github.com/gin-gonic/gin"
)

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token     string   `json:"token"`
	ExpiresIn int      `json:"expires_in"`
	User      UserInfo `json:"user"`
}

// UserInfo represents user information
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err.Error())
		return
	}

	// TODO: Call user center service for authentication
	h.Success(c, LoginResponse{
		Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		ExpiresIn: 86400,
		User: UserInfo{
			ID:       "user-123",
			Username: req.Username,
			Email:    "user@example.com",
			Role:     "developer",
		},
	})
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, err.Error())
		return
	}

	// TODO: Call user center service for registration
	h.Success(c, UserInfo{
		ID:       "user-456",
		Username: req.Username,
		Email:    req.Email,
		Role:     "developer",
	})
}

// GetCurrentUser returns the current user
func (h *Handler) GetCurrentUser(c *gin.Context) {
	// TODO: Get user from JWT token
	h.Success(c, UserInfo{
		ID:       "user-123",
		Username: "johndoe",
		Email:    "john@example.com",
		Role:     "developer",
	})
}
