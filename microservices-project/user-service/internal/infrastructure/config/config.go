package config

import (
	"os"
)

type Config struct {
	ServiceName string
	Port        string
	DBDSN       string
	LogLevel    string
}

func Load() Config {
	return Config{
		ServiceName: getEnv("SERVICE_NAME", "user-service"),
		Port:        getEnv("PORT", "8080"),
		DBDSN:       getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
