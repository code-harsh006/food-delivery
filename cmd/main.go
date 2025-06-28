package main 

import (
	"context"
	"fmt"
	"log"
	"os"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"food-delivery/internal/auth"
	"food-delivery/internal/user"
	"food-delivery/internal/produck"
	"food-delivery/internal/vendor"
	"food-devivery/internal/cart"
	"food-devivery/internal/order"
	"food-devivery/internal/payment"
	"food-delivery/internal/delivery"
	"food-delibery/internal/admin"
	"food-delivery/internal/notification"
	"food-delivery/internal/search"
	"food-delivery/internal/config"
	"food-delivery/pkg/config"
	"food-delivery/pkg/db"
	"food-delivery/pkg/logger"
	"food-delivery/pkg/middleware"
	"food-delivery/pkg/metrics"
	"food-delivery/pkg/heathcheck"
	"food-devivery/pkg/radinst"
	"food-delivery/pkg/tracking"
	"food-devivery/pkg/redis"
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
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Redis
	redisClient := db.InitRedis(cfg.RedisURL)

	// Run migrations
	if err := db.Migrate(database); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Seed data
	if err := db.SeedData(database); err != nil {
		log.Println("Warning: Failed to seed data:", err)
	}

	// Initialize Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Initialize modules
	authModule := auth.NewModule(database)
	userModule := user.NewModule(database)
	productModule := product.NewModule(database)
	vendorModule := vendor.NewModule(database)
	cartModule := cart.NewModule(database, redisClient)
	orderModule := order.NewModule(database)
	paymentModule := payment.NewModule(database)
	deliveryModule := delivery.NewModule(database)
	adminModule := admin.NewModule(database)
	notificationModule := notification.NewModule()
	searchModule := search.NewModule(database)

	// Setup routes
	api := router.Group("/api/v1")
	{
		authModule.SetupRoutes(api)
		userModule.SetupRoutes(api, middleware.AuthMiddleware())
		productModule.SetupRoutes(api, middleware.AuthMiddleware())
		vendorModule.SetupRoutes(api, middleware.AuthMiddleware())
		cartModule.SetupRoutes(api, middleware.AuthMiddleware())
		orderModule.SetupRoutes(api, middleware.AuthMiddleware())
		paymentModule.SetupRoutes(api, middleware.AuthMiddleware())
		deliveryModule.SetupRoutes(api, middleware.AuthMiddleware())
		adminModule.SetupRoutes(api, middleware.AuthMiddleware(), middleware.AdminMiddleware())
		notificationModule.SetupRoutes(api, middleware.AuthMiddleware())
		searchModule.SetupRoutes(api, middleware.AuthMiddleware())
	}

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

	fmt.Printf("Server starting on port %s\n", cfg.Port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("Server exited")
}

