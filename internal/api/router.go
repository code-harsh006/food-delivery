package api

import (
	"github.com/code-harsh006/food-delivery/internal/notification"
	"github.com/code-harsh006/food-delivery/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// APIRouter handles the main API routing and integration
type APIRouter struct {
	router        *gin.Engine
	healthHandler *HealthHandler
}

// NewAPIRouter creates a new APIRouter instance
func NewAPIRouter(router *gin.Engine) *APIRouter {
	return &APIRouter{
		router:        router,
		healthHandler: NewHealthHandler(),
	}
}

// SetupRoutes sets up all API routes
func (r *APIRouter) SetupRoutes() {
	// Health check routes (no authentication required)
	r.healthHandler.SetupHealthRoutes(r.router.Group("/api/v1"))

	// Main API routes
	api := r.router.Group("/api/v1")
	{
		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Notification routes
			notificationModule := notification.NewModule()
			notificationModule.SetupRoutes(protected, middleware.AuthMiddleware())
		}
	}

	// MongoDB routes (setup after main routes to avoid conflicts)
	r.SetupMongoDBRoutes()
}

// GetRouter returns the configured router
func (r *APIRouter) GetRouter() *gin.Engine {
	return r.router
}
