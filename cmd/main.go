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
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize logger
	logger.Init()

	// Load configuration
	cfg := config.Load()

	// Initialize database
	database, err := db.Init(cfg.DatabaseURL)
	if err != nil {
		log.Printf("⚠️  Database connection failed: %v", err)
		log.Println("🚀 Starting server without database (limited functionality)")
		database = nil
	} else {
		log.Println("✅ Database connected successfully")

		// Run migrations
		if err := db.Migrate(database); err != nil {
			log.Printf("⚠️  Failed to run migrations: %v", err)
		}

		// Seed data
		if err := db.SeedData(database); err != nil {
			log.Printf("⚠️  Failed to seed data: %v", err)
		}
	}

	// Initialize Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Initialize API router
	apiRouter := api.NewAPIRouter(router, database)

	// Setup all API routes
	apiRouter.SetupRoutes()

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	fmt.Printf("🚀 Food Delivery Server starting on port %s\n", cfg.Port)
	fmt.Printf("📡 Health check: http://localhost:%s/api/v1/health\n", cfg.Port)
	fmt.Printf("🔗 API endpoint: http://localhost:%s/api/v1\n", cfg.Port)
	fmt.Printf("📋 API status: http://localhost:%s/api/v1/status\n", cfg.Port)
	fmt.Printf("📖 API docs: http://localhost:%s/api/v1/docs\n", cfg.Port)
	if database == nil {
		fmt.Printf("⚠️  Database: Not connected (limited functionality)\n")
	} else {
		fmt.Printf("✅ Database: Connected\n")
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
