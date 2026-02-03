package config

import "os"

type AppConfig struct {
	Env      string
	HTTPAddr string
	RPCAddr  string
	LogLevel string
}

func Load() *AppConfig {
	return &AppConfig{
		Env:      getEnv("APP_ENV", "dev"),
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		RPCAddr:  getEnv("RPC_ADDR", ":9090"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
