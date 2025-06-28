package api

import (
	"github.com/code-harsh006/food-delivery/internal/admin"
	"github.com/code-harsh006/food-delivery/internal/auth"
	"github.com/code-harsh006/food-delivery/internal/cart"
	"github.com/code-harsh006/food-delivery/internal/notification"
	"github.com/code-harsh006/food-delivery/internal/product"
	"github.com/code-harsh006/food-delivery/internal/search"
	"github.com/code-harsh006/food-delivery/internal/user"
	"github.com/code-harsh006/food-delivery/internal/vendor"
	"github.com/code-harsh006/food-delivery/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// APIRouter handles the main API routing and integration
type APIRouter struct {
	router        *gin.Engine
	db            *gorm.DB
	authHandler   *AuthHandler
	indexHandler  *IndexHandler
	healthHandler *HealthHandler
}

// NewAPIRouter creates a new APIRouter instance
func NewAPIRouter(router *gin.Engine, db *gorm.DB) *APIRouter {
	return &APIRouter{
		router:        router,
		db:            db,
		authHandler:   NewAuthHandler(auth.NewModule(db)),
		indexHandler:  NewIndexHandler(),
		healthHandler: NewHealthHandler(db),
	}
}

// SetupRoutes sets up all API routes
func (r *APIRouter) SetupRoutes() {
	// Health check routes (no authentication required)
	r.healthHandler.SetupHealthRoutes(r.router.Group("/api/v1"))

	// Main API routes
	api := r.router.Group("/api/v1")
	{
		// Index routes (no authentication required)
		r.indexHandler.SetupIndexRoutes(api)

		// Auth routes (no authentication required for login/register)
		r.authHandler.SetupAuthRoutes(api)

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			if r.db != nil {
				userModule := user.NewModule(r.db)
				userModule.SetupRoutes(protected, middleware.AuthMiddleware())
			}

			// Product routes
			if r.db != nil {
				productModule := product.NewModule(r.db)
				productModule.SetupRoutes(protected, middleware.AuthMiddleware())
			}

			// Vendor routes
			if r.db != nil {
				vendorModule := vendor.NewModule(r.db)
				vendorModule.SetupRoutes(protected, middleware.AuthMiddleware())
			}

			// Cart routes
			if r.db != nil {
				cartModule := cart.NewModule(r.db, nil) // Redis client would be passed here
				cartModule.SetupRoutes(protected, middleware.AuthMiddleware())
			}

			// Notification routes
			notificationModule := notification.NewModule()
			notificationModule.SetupRoutes(protected, middleware.AuthMiddleware())

			// Search routes
			if r.db != nil {
				searchModule := search.NewModule(r.db)
				searchModule.SetupRoutes(protected, middleware.AuthMiddleware())
			}
		}

		// Admin routes (require admin authentication)
		if r.db != nil {
			admin.Module(r.router, r.db)
		}
	}
}

// GetRouter returns the configured router
func (r *APIRouter) GetRouter() *gin.Engine {
	return r.router
}
