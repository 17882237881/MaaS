package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"maas-platform/api-gateway/internal/config"
	"maas-platform/api-gateway/internal/handler"
	"maas-platform/api-gateway/internal/middleware"
	"maas-platform/api-gateway/internal/router"
	"maas-platform/api-gateway/pkg/logger"
)

// @title MaaS Platform API
// @version 1.0
// @description Model-as-a-Service Platform API Gateway
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@maas-platform.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.LogLevel)

	log.Info("Starting API Gateway",
		"version", "1.0.0",
		"environment", cfg.Environment,
		"port", cfg.Port,
	)

	// Set gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	r := gin.New()

	// Global middleware
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "api-gateway",
			"timestamp": time.Now().Unix(),
		})
	})

	// Register routes
	api := r.Group("/api/v1")
	h := handler.New(cfg, log)
	router.RegisterRoutes(api, h)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("Server forced to shutdown", "error", err)
		}
	}()

	// Start server
	log.Info("Server starting", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server", "error", err)
	}

	log.Info("Server exited")
}
