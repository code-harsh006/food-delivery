package config

import (
	"os"
)

type Config struct {
	DatabaseURL   string
	RedisURL      string
	JWTSecret     string
	Port          string
	Environment   string
	StripeKey     string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost/fooddelivery?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		StripeKey:   getEnv("STRIPE_KEY", "sk_test_dummy"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

