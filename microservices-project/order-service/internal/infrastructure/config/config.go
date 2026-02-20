package config

import "os"

type Config struct {
	ServiceName    string
	Port           string
	DBDSN          string
	UserServiceURL string
	LogLevel       string
}

func Load() Config {
	return Config{
		ServiceName:    getEnv("SERVICE_NAME", "order-service"),
		Port:           getEnv("PORT", "8081"),
		DBDSN:          getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/orderdb?sslmode=disable"),
		UserServiceURL: getEnv("USER_SERVICE_URL", "http://user-service:8080"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
