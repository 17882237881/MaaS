package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds all configuration for the Model Registry
type Config struct {
	Environment string         `mapstructure:"environment"`
	Port        int            `mapstructure:"port"`
	LogLevel    string         `mapstructure:"log_level"`
	Database    DatabaseConfig `mapstructure:"database"`
	Redis       RedisConfig    `mapstructure:"redis"`
	Services    ServiceConfig  `mapstructure:"services"`
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

// ServiceConfig holds service addresses
type ServiceConfig struct {
	APIGateway string `mapstructure:"api_gateway"`
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
	viper.SetDefault("port", 8081)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.name", "maas_registry")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

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
