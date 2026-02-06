package config

// Config holds all configuration for the Model Registry
type Config struct {
	Environment string
	Port        int
	LogLevel    string
}

// Load returns the application configuration
func Load() *Config {
	return &Config{
		Environment: "development",
		Port:        8081,
		LogLevel:    "info",
	}
}
