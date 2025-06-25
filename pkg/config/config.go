package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Payment  PaymentConfig
	External ExternalConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret        string
	Expiry        time.Duration
	RefreshExpiry time.Duration
}

type PaymentConfig struct {
	StripeSecretKey    string
	StripeWebhookSecret string
}

type ExternalConfig struct {
	NotificationServiceURL string
	MapsAPIKey            string
}

type LoggingConfig struct {
	Level  string
	Format string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using environment variables: %v", err)
	}

	// Handle Heroku/Railway DATABASE_URL
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		// Parse DATABASE_URL for Heroku/Railway
		viper.Set("DATABASE_URL", databaseURL)
	}

	// Handle Redis URL
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		viper.Set("REDIS_URL", redisURL)
	}

	jwtExpiry, _ := time.ParseDuration(getEnvOrDefault("JWT_EXPIRY", "24h"))
	jwtRefreshExpiry, _ := time.ParseDuration(getEnvOrDefault("JWT_REFRESH_EXPIRY", "168h"))

	return &Config{
		Server: ServerConfig{
			Port:    getEnvOrDefault("PORT", "8080"),
			GinMode: getEnvOrDefault("GIN_MODE", "release"),
		},
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "password"),
			Name:     getEnvOrDefault("DB_NAME", "food_delivery"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnvOrDefault("REDIS_HOST", "localhost"),
			Port:     getEnvOrDefault("REDIS_PORT", "6379"),
			Password: getEnvOrDefault("REDIS_PASSWORD", ""),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:        getEnvOrDefault("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
			Expiry:        jwtExpiry,
			RefreshExpiry: jwtRefreshExpiry,
		},
		Payment: PaymentConfig{
			StripeSecretKey:     getEnvOrDefault("STRIPE_SECRET_KEY", ""),
			StripeWebhookSecret: getEnvOrDefault("STRIPE_WEBHOOK_SECRET", ""),
		},
		External: ExternalConfig{
			NotificationServiceURL: getEnvOrDefault("NOTIFICATION_SERVICE_URL", ""),
			MapsAPIKey:            getEnvOrDefault("MAPS_API_KEY", ""),
		},
		Logging: LoggingConfig{
			Level:  getEnvOrDefault("LOG_LEVEL", "info"),
			Format: getEnvOrDefault("LOG_FORMAT", "json"),
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return viper.GetString(key)
}
