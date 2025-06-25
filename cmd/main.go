package main

import (
	"context"
	"fmt"
	"food-delivery-backend/internal/auth"
	"food-delivery-backend/internal/cart"
	"food-delivery-backend/internal/product"
	"food-delivery-backend/internal/user"
	"food-delivery-backend/pkg/config"
	"food-delivery-backend/pkg/db"
	"food-delivery-backend/pkg/logger"
	"food-delivery-backend/pkg/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.Logging.Level, cfg.Logging.Format)

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize databases
	if err := db.InitPostgres(&cfg.Database); err != nil {
		logger.Fatal("Failed to connect to PostgreSQL: ", err)
	}

	if err := db.InitRedis(&cfg.Redis); err != nil {
		logger.Fatal("Failed to connect to Redis: ", err)
	}

	// Auto-migrate database
	if err := autoMigrate(); err != nil {
		logger.Fatal("Failed to migrate database: ", err)
	}

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.ErrorHandlingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "food-delivery-backend",
		})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Food Delivery Backend API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Register module routes
		auth.RegisterRoutes(api, cfg)
		user.RegisterRoutes(api, cfg)
		product.RegisterRoutes(api, cfg)
		cart.RegisterRoutes(api, cfg)
	}

	// Get port from environment or config
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}
	if port == "" {
		port = "8080"
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: ", err)
		}
	}()

	logger.Info("Server started on port ", port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}

	logger.Info("Server exited")
}

func autoMigrate() error {
	// Import all models here for auto-migration
	models := []interface{}{
		&auth.User{},
		&auth.RefreshToken{},
		&user.Address{},
		&user.UserProfile{},
		&product.Category{},
		&product.Product{},
		&cart.Cart{},
		&cart.CartItem{},
	}

	for _, model := range models {
		if err := db.DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	logger.Info("Database migration completed successfully")
	return nil
}
