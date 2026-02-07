package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config holds all configuration for the API Gateway
type Config struct {
	Environment string `mapstructure:"environment"`
	Port        int    `mapstructure:"port"`
	LogLevel    string `mapstructure:"log_level"`

	// Database
	Database DatabaseConfig `mapstructure:"database"`

	// Redis
	Redis RedisConfig `mapstructure:"redis"`

	// JWT
	JWT JWTConfig `mapstructure:"jwt"`

	// Services
	Services ServiceConfig `mapstructure:"services"`

	// Rate Limiting
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret    string `mapstructure:"secret"`
	ExpiresIn int    `mapstructure:"expires_in"`
}

// ServiceConfig holds downstream service URLs
type ServiceConfig struct {
	ModelRegistry string `mapstructure:"model_registry"`
	Inference     string `mapstructure:"inference"`
	UserCenter    string `mapstructure:"user_center"`
	Billing       string `mapstructure:"billing"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled bool `mapstructure:"enabled"`
	RPM     int  `mapstructure:"rpm"`
	Burst   int  `mapstructure:"burst"`
}

// Load returns the application configuration
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Read from environment variables
	v.SetEnvPrefix("MAAS")
	v.AutomaticEnv()

	// Read config file
	if err := loadConfigFile(v); err != nil {
		return nil, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// LoadWithWatch returns config and sets up file watching for hot reload
func LoadWithWatch(onChange func(*Config)) (*Config, error) {
	config, err := Load()
	if err != nil {
		return nil, err
	}

	// Setup file watcher
	setupWatcher(onChange)

	return config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	v.SetDefault("environment", "development")
	v.SetDefault("port", 8080)
	v.SetDefault("log_level", "info")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "maas_platform")
	v.SetDefault("database.ssl_mode", "disable")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("jwt.secret", "change-me-in-production")
	v.SetDefault("jwt.expires_in", 86400)

	v.SetDefault("services.model_registry", "http://localhost:8081")
	v.SetDefault("services.inference", "http://localhost:8082")
	v.SetDefault("services.user_center", "http://localhost:8083")
	v.SetDefault("services.billing", "http://localhost:8084")

	v.SetDefault("rate_limit.enabled", true)
	v.SetDefault("rate_limit.rpm", 1000)
	v.SetDefault("rate_limit.burst", 100)
}

// loadConfigFile attempts to load configuration from file
func loadConfigFile(v *viper.Viper) error {
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Add config paths
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/maas/")
	v.AddConfigPath("$HOME/.maas")

	// Try to read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error occurred
			return fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; use defaults and environment variables
		log.Println("Config file not found, using defaults and environment variables")
	} else {
		log.Printf("Loaded config file: %s", v.ConfigFileUsed())
	}

	return nil
}

// setupWatcher sets up file watching for configuration hot reload
func setupWatcher(onChange func(*Config)) {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create config watcher: %v", err)
		return
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Config file modified, reloading...")

					// Reload config
					if err := viper.ReadInConfig(); err != nil {
						log.Printf("Error reloading config: %v", err)
						continue
					}

					var newConfig Config
					if err := viper.Unmarshal(&newConfig); err != nil {
						log.Printf("Error unmarshaling config: %v", err)
						continue
					}

					if err := newConfig.Validate(); err != nil {
						log.Printf("Config validation failed: %v", err)
						continue
					}

					log.Println("Config reloaded successfully")
					if onChange != nil {
						onChange(&newConfig)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Config watcher error: %v", err)
			}
		}
	}()

	// Watch the config file directory
	configDir := filepath.Dir(configFile)
	if err := watcher.Add(configDir); err != nil {
		log.Printf("Failed to watch config directory: %v", err)
	}

	log.Printf("Config hot reload enabled, watching: %s", configFile)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate environment
	if c.Environment != "development" && c.Environment != "staging" && c.Environment != "production" {
		return fmt.Errorf("invalid environment: %s", c.Environment)
	}

	// Validate port
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port: %d", c.Port)
	}

	// Validate log level
	if c.LogLevel != "debug" && c.LogLevel != "info" && c.LogLevel != "warn" && c.LogLevel != "error" {
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	// Validate database
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Database.Port)
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}

	// Validate JWT secret in production
	if c.Environment == "production" {
		if c.JWT.Secret == "" || c.JWT.Secret == "change-me-in-production" {
			return fmt.Errorf("JWT secret must be set in production")
		}
		if len(c.JWT.Secret) < 32 {
			return fmt.Errorf("JWT secret must be at least 32 characters in production")
		}
	}

	// Validate rate limit
	if c.RateLimit.RPM <= 0 {
		return fmt.Errorf("rate limit RPM must be positive")
	}

	return nil
}

// GetEnvOrDefault returns environment variable value or default
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsDevelopment returns true if in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// DatabaseDSN returns PostgreSQL connection string
func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Database.Host,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.Port,
		c.Database.SSLMode,
	)
}

// RedisAddr returns Redis address
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}
