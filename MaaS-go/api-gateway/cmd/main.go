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
	"maas-platform/api-gateway/internal/service"
	rpc "maas-platform/api-gateway/pkg/grpc"
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
	// Load configuration with hot reload support
	cfg, err := config.LoadWithWatch(func(newCfg *config.Config) {
		// This callback is called when config file changes
		// You can update global settings here
		fmt.Println("Configuration reloaded")
	})
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Validate that required configs are present
	if err := validateEssentialConfig(cfg); err != nil {
		fmt.Printf("Configuration validation failed: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.LogLevel)

	log.Info("Starting API Gateway",
		"version", "1.0.0",
		"environment", cfg.Environment,
		"port", cfg.Port,
		"config_file", "config/config.yaml",
	)

	// Print configuration summary
	printConfigSummary(cfg, log)

	// Initialize gRPC client to Model Registry
	log.Info("Connecting to Model Registry gRPC service...", "address", cfg.Services.ModelRegistry)
	grpcClient, err := rpc.NewClient(cfg.Services.ModelRegistry)
	if err != nil {
		log.Fatal("Failed to connect to Model Registry", "error", err)
	}
	defer grpcClient.Close()
	log.Info("Connected to Model Registry gRPC service")

	// Initialize model service client
	modelServiceClient := service.NewModelServiceClient(grpcClient, log)

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
			"config":    cfg.Environment,
		})
	})

	// Config info endpoint (for debugging, disable in production)
	if cfg.IsDevelopment() {
		r.GET("/config", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"environment": cfg.Environment,
				"port":        cfg.Port,
				"log_level":   cfg.LogLevel,
				"database": gin.H{
					"host": cfg.Database.Host,
					"port": cfg.Database.Port,
					"name": cfg.Database.Name,
				},
				"redis": gin.H{
					"host": cfg.Redis.Host,
					"port": cfg.Redis.Port,
				},
			})
		})
	}

	// Register routes
	api := r.Group("/api/v1")
	h := handler.New(cfg, log, modelServiceClient)
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

// validateEssentialConfig validates that essential configuration is present
func validateEssentialConfig(cfg *config.Config) error {
	if cfg.Environment == "" {
		return fmt.Errorf("environment is required")
	}

	if cfg.Port <= 0 {
		return fmt.Errorf("invalid port: %d", cfg.Port)
	}

	// In production, ensure JWT secret is properly set
	if cfg.IsProduction() {
		if cfg.JWT.Secret == "" {
			return fmt.Errorf("JWT_SECRET is required in production")
		}
		if len(cfg.JWT.Secret) < 32 {
			return fmt.Errorf("JWT_SECRET must be at least 32 characters in production")
		}
	}

	return nil
}

// printConfigSummary prints a summary of the loaded configuration
func printConfigSummary(cfg *config.Config, log *logger.Logger) {
	log.Info("Configuration loaded",
		"environment", cfg.Environment,
		"port", cfg.Port,
		"log_level", cfg.LogLevel,
		"rate_limit_enabled", cfg.RateLimit.Enabled,
		"rate_limit_rpm", cfg.RateLimit.RPM,
	)

	log.Info("Database configuration",
		"host", cfg.Database.Host,
		"port", cfg.Database.Port,
		"name", cfg.Database.Name,
		"ssl_mode", cfg.Database.SSLMode,
	)

	log.Info("Redis configuration",
		"host", cfg.Redis.Host,
		"port", cfg.Redis.Port,
		"db", cfg.Redis.DB,
	)
}
