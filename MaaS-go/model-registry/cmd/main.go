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

	"maas-platform/model-registry/internal/config"
	"maas-platform/model-registry/internal/handler"
	"maas-platform/model-registry/internal/middleware"
	"maas-platform/model-registry/internal/repository"
	"maas-platform/model-registry/internal/router"
	"maas-platform/model-registry/internal/service"
	"maas-platform/model-registry/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.LogLevel)

	log.Info("Starting Model Registry Service",
		"version", "1.0.0",
		"environment", cfg.Environment,
		"port", cfg.Port,
	)

	// Connect to database
	db, err := repository.NewDatabase(repository.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	// Run migrations
	log.Info("Running database migrations...")
	if err := repository.AutoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database", "error", err)
	}
	log.Info("Database migrations completed")

	// Initialize repository
	modelRepo := repository.NewGormModelRepository(db)

	// Initialize service
	modelService := service.NewModelService(modelRepo, log)

	// Initialize handler
	modelHandler := handler.NewModelHandler(modelService, log)

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
		// Check database connectivity
		sqlDB, err := db.DB()
		dbStatus := "up"
		if err != nil {
			dbStatus = "down"
		} else if err := sqlDB.Ping(); err != nil {
			dbStatus = "down"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"service":     "model-registry",
			"database":    dbStatus,
			"timestamp":   time.Now().Unix(),
			"environment": cfg.Environment,
		})
	})

	// API routes
	api := r.Group("/api/v1")
	router.RegisterRoutes(api, modelHandler)

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
