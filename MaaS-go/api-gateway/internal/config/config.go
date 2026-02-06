package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds all configuration for the API Gateway
type Config struct {
	Environment string `mapstructure:"environment"`
	Port        int    `mapstructure:"port"`
	LogLevel    string `mapstructure:"log_level"`

	// Database
	DBHost     string `mapstructure:"db_host"`
	DBPort     int    `mapstructure:"db_port"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
	DBName     string `mapstructure:"db_name"`
	DBSSLMode  string `mapstructure:"db_ssl_mode"`

	// Redis
	RedisHost     string `mapstructure:"redis_host"`
	RedisPort     int    `mapstructure:"redis_port"`
	RedisPassword string `mapstructure:"redis_password"`
	RedisDB       int    `mapstructure:"redis_db"`

	// JWT
	JWTSecret    string `mapstructure:"jwt_secret"`
	JWTExpiresIn int    `mapstructure:"jwt_expires_in"`

	// Services
	ModelRegistryURL string `mapstructure:"model_registry_url"`
	InferenceURL     string `mapstructure:"inference_url"`
	UserCenterURL    string `mapstructure:"user_center_url"`
	BillingURL       string `mapstructure:"billing_url"`

	// Rate Limiting
	RateLimitEnabled bool `mapstructure:"rate_limit_enabled"`
	RateLimitRPM     int  `mapstructure:"rate_limit_rpm"`
	RateLimitBurst   int  `mapstructure:"rate_limit_burst"`
}

// Load returns the application configuration
func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/maas/")

	// Set defaults
	viper.SetDefault("environment", "development")
	viper.SetDefault("port", 8080)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_port", 5432)
	viper.SetDefault("db_ssl_mode", "disable")
	viper.SetDefault("redis_host", "localhost")
	viper.SetDefault("redis_port", 6379)
	viper.SetDefault("redis_db", 0)
	viper.SetDefault("jwt_expires_in", 86400)
	viper.SetDefault("rate_limit_enabled", true)
	viper.SetDefault("rate_limit_rpm", 1000)
	viper.SetDefault("rate_limit_burst", 100)

	// Read from environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MAAS")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: config file not found: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	return &config
}
