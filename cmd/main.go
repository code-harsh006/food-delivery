package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/code-harsh006/food-delivery/internal/api"
	"github.com/code-harsh006/food-delivery/pkg/config"
	"github.com/code-harsh006/food-delivery/pkg/db"
	"github.com/code-harsh006/food-delivery/pkg/logger"
	"github.com/code-harsh006/food-delivery/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting Food Delivery Server...")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize logger
	logger.Init()
	fmt.Println("Logger initialized")

	// Load configuration
	cfg := config.Load()
	fmt.Println("Configuration loaded")

	// Check for PORT environment variable and override cfg.Port if set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback default
	}
	cfg.Port = port
	fmt.Printf("Using port: %s\n", cfg.Port)

	// Initialize MongoDB database
	if err := db.InitMongoDB(cfg.MongoDBURI, cfg.Environment == "production"); err != nil {
		log.Printf("⚠️  MongoDB connection failed: %v", err)
		log.Println("🚀 Starting server without MongoDB (limited functionality)")
	} else {
		log.Println("✅ MongoDB connected successfully")
	}

	// Initialize Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	fmt.Println("Initializing Gin router")

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	fmt.Println("Router middleware configured")

	// Initialize API router
	apiRouter := api.NewAPIRouter(router)
	fmt.Println("API router initialized")

	// Setup all API routes
	apiRouter.SetupRoutes()
	fmt.Println("API routes configured")

	// Start server
	srv := &http.Server{
		Addr:    "0.0.0.0:" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		fmt.Printf("🚀 Starting server on port %s\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	fmt.Printf("🚀 Food Delivery Server starting on port %s\n", cfg.Port)
	fmt.Printf("\n======= API Endpoints =======\n")
	fmt.Printf("📡 Health check:      http://0.0.0.0:%s/api/v1/health         [working]\n", cfg.Port)
	fmt.Printf("📡 Health detailed:   http://0.0.0.0:%s/api/v1/health/detailed [working]\n", cfg.Port)
	fmt.Printf("📡 Health ready:      http://0.0.0.0:%s/api/v1/health/ready     [working]\n", cfg.Port)
	fmt.Printf("📡 Health live:       http://0.0.0.0:%s/api/v1/health/live      [working]\n", cfg.Port)
	fmt.Printf("📋 API status:        http://0.0.0.0:%s/api/v1/status           [working]\n", cfg.Port)
	fmt.Printf("📖 API docs:          http://0.0.0.0:%s/api/v1/docs             [working]\n", cfg.Port)
	fmt.Printf("🔗 API v1 root:       http://0.0.0.0:%s/api/v1                  [working]\n", cfg.Port)
	fmt.Printf("\n--- Main API Groups ---\n")
	fmt.Printf("Auth:      /api/v1/auth/*         [working]\n")
	fmt.Printf("Users:     /api/v1/users/*        [working]\n")
	fmt.Printf("Products:  /api/v1/products/*     [working]\n")
	fmt.Printf("Cart:      /api/v1/cart/*         [working]\n")
	fmt.Printf("Orders:    /api/v1/orders/*       [working]\n")
	fmt.Printf("Vendors:   /api/v1/vendors/*      [working]\n")
	fmt.Printf("\n--- MongoDB API (v1) ---\n")
	fmt.Printf("MongoDB Auth:        /api/mongo/v1/auth/*         [working]\n")
	fmt.Printf("MongoDB Services:    /api/mongo/v1/services/*     [working]\n")
	fmt.Printf("MongoDB Bookings:    /api/mongo/v1/bookings/*     [working]\n")
	fmt.Printf("MongoDB Users:       /api/mongo/v1/users/*        [working]\n")
	fmt.Printf("MongoDB Admin:       /api/mongo/v1/admin/*        [working]\n")
	fmt.Printf("==============================\n\n")
	if db.GetMongoDB() == nil {
		fmt.Printf("⚠️  MongoDB: Not connected (limited functionality)\n")
	} else {
		fmt.Printf("✅ MongoDB: Connected\n")
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("✅ Server exited gracefully")
}
