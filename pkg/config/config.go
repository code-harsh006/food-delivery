package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Database Configuration
	DatabaseURL string
	RedisURL    string

	// JWT Configuration
	JWTSecret string

	// Server Configuration
	Port        string
	Environment string

	// Payment Configuration
	StripeKey            string
	StripePublishableKey string
	StripeWebhookSecret  string
	PayPalClientID       string
	PayPalClientSecret   string
	PayPalMode           string

	// Map & Location Services
	GoogleMapsAPIKey  string
	MapboxAccessToken string
	DefaultLatitude   float64
	DefaultLongitude  float64
	DefaultRadius     int

	// Tracking & Delivery
	TrackingAPIKey           string
	DeliveryRadiusKM         int
	EstimatedDeliveryTimeMin int
	MaxDeliveryTimeMin       int

	// Email Configuration
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	EmailFrom    string

	// SMS Configuration
	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioPhoneNumber string

	// Push Notifications
	FirebaseProjectID    string
	FirebasePrivateKeyID string
	FirebasePrivateKey   string
	FirebaseClientEmail  string
	FirebaseClientID     string

	// File Upload
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSS3Bucket        string
	MaxFileSize        int64

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   int

	// Logging
	LogLevel string
	LogFile  string

	// Security
	CORSOrigin    string
	SessionSecret string
	BcryptCost    int

	// Cache Configuration
	CacheTTL      int
	RedisCacheTTL int

	// Monitoring
	SentryDSN          string
	NewRelicLicenseKey string
}

func Load() *Config {
	return &Config{
		// Database Configuration
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost/fooddelivery?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),

		// JWT Configuration
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),

		// Server Configuration
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Payment Configuration
		StripeKey:            getEnv("STRIPE_KEY", "sk_test_dummy"),
		StripePublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", "pk_test_dummy"),
		StripeWebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", "whsec_dummy"),
		PayPalClientID:       getEnv("PAYPAL_CLIENT_ID", ""),
		PayPalClientSecret:   getEnv("PAYPAL_CLIENT_SECRET", ""),
		PayPalMode:           getEnv("PAYPAL_MODE", "sandbox"),

		// Map & Location Services
		GoogleMapsAPIKey:  getEnv("GOOGLE_MAPS_API_KEY", ""),
		MapboxAccessToken: getEnv("MAPBOX_ACCESS_TOKEN", ""),
		DefaultLatitude:   getEnvAsFloat("DEFAULT_LATITUDE", 40.7128),
		DefaultLongitude:  getEnvAsFloat("DEFAULT_LONGITUDE", -74.0060),
		DefaultRadius:     getEnvAsInt("DEFAULT_RADIUS", 5000),

		// Tracking & Delivery
		TrackingAPIKey:           getEnv("TRACKING_API_KEY", ""),
		DeliveryRadiusKM:         getEnvAsInt("DELIVERY_RADIUS_KM", 10),
		EstimatedDeliveryTimeMin: getEnvAsInt("ESTIMATED_DELIVERY_TIME_MIN", 30),
		MaxDeliveryTimeMin:       getEnvAsInt("MAX_DELIVERY_TIME_MIN", 60),

		// Email Configuration
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		EmailFrom:    getEnv("EMAIL_FROM", "noreply@fooddelivery.com"),

		// SMS Configuration
		TwilioAccountSID:  getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:   getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioPhoneNumber: getEnv("TWILIO_PHONE_NUMBER", ""),

		// Push Notifications
		FirebaseProjectID:    getEnv("FIREBASE_PROJECT_ID", ""),
		FirebasePrivateKeyID: getEnv("FIREBASE_PRIVATE_KEY_ID", ""),
		FirebasePrivateKey:   getEnv("FIREBASE_PRIVATE_KEY", ""),
		FirebaseClientEmail:  getEnv("FIREBASE_CLIENT_EMAIL", ""),
		FirebaseClientID:     getEnv("FIREBASE_CLIENT_ID", ""),

		// File Upload
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSRegion:          getEnv("AWS_REGION", "us-east-1"),
		AWSS3Bucket:        getEnv("AWS_S3_BUCKET", ""),
		MaxFileSize:        getEnvAsInt64("MAX_FILE_SIZE", 5242880), // 5MB

		// Rate Limiting
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvAsInt("RATE_LIMIT_WINDOW", 900),

		// Logging
		LogLevel: getEnv("LOG_LEVEL", "info"),
		LogFile:  getEnv("LOG_FILE", "logs/app.log"),

		// Security
		CORSOrigin:    getEnv("CORS_ORIGIN", "http://localhost:3000"),
		SessionSecret: getEnv("SESSION_SECRET", "your_session_secret"),
		BcryptCost:    getEnvAsInt("BCRYPT_COST", 12),

		// Cache Configuration
		CacheTTL:      getEnvAsInt("CACHE_TTL", 3600),
		RedisCacheTTL: getEnvAsInt("REDIS_CACHE_TTL", 1800),

		// Monitoring
		SentryDSN:          getEnv("SENTRY_DSN", ""),
		NewRelicLicenseKey: getEnv("NEW_RELIC_LICENSE_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}
