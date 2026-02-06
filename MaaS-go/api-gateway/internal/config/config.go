package config

// Config holds all configuration for the API Gateway
type Config struct {
	Environment string
	Port        int
	LogLevel    string
}

// Load returns the application configuration
func Load() *Config {
	return &Config{
		Environment: "development",
		Port:        8080,
		LogLevel:    "info",
	}
}
