package api

import (
	"net/http"

	"github.com/code-harsh006/food-delivery/internal/notification"
	"github.com/code-harsh006/food-delivery/pkg/middleware"
	"github.com/code-harsh006/food-delivery/pkg/response"
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
	// Root route
	r.router.GET("/", r.rootHandler)

	// Health check routes (no authentication required)
	r.healthHandler.SetupHealthRoutes(r.router.Group("/api/v1"))

	// Main API routes
	api := r.router.Group("/api/v1")
	{
		// API documentation
		api.GET("/docs", r.docsHandler)

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

	// 404 handler for unmatched routes
	r.router.NoRoute(r.notFoundHandler)
}

// rootHandler handles the root path
func (r *APIRouter) rootHandler(c *gin.Context) {
	response.Success(c, gin.H{
		"message": "Food Delivery API",
		"version": "1.0.0",
		"status":  "running",
		"endpoints": gin.H{
			"health": "/api/v1/health",
			"docs":   "/api/v1/docs",
			"status": "/api/v1/status",
		},
		"documentation": "Visit /api/v1/docs for API documentation",
	})
}

// docsHandler provides API documentation
func (r *APIRouter) docsHandler(c *gin.Context) {
	response.Success(c, gin.H{
		"title":       "Food Delivery API Documentation",
		"version":     "1.0.0",
		"description": "A comprehensive Go-based backend API for food delivery services",
		"endpoints": gin.H{
			"root": gin.H{
				"method":      "GET",
				"path":        "/",
				"description": "API information and available endpoints",
			},
			"health": gin.H{
				"method":      "GET",
				"path":        "/api/v1/health",
				"description": "Basic health check",
			},
			"health_detailed": gin.H{
				"method":      "GET",
				"path":        "/api/v1/health/detailed",
				"description": "Detailed health check with system information",
			},
			"health_ready": gin.H{
				"method":      "GET",
				"path":        "/api/v1/health/ready",
				"description": "Readiness check for load balancers",
			},
			"health_live": gin.H{
				"method":      "GET",
				"path":        "/api/v1/health/live",
				"description": "Liveness check for container orchestration",
			},
			"status": gin.H{
				"method":      "GET",
				"path":        "/api/v1/status",
				"description": "API status and version information",
			},
			"docs": gin.H{
				"method":      "GET",
				"path":        "/api/v1/docs",
				"description": "This documentation page",
			},
		},
		"authentication": "Most endpoints require JWT authentication",
		"base_url":       "http://localhost:8080",
	})
}

// notFoundHandler handles 404 errors
func (r *APIRouter) notFoundHandler(c *gin.Context) {
	response.Error(c, http.StatusNotFound, "Endpoint not found. Visit / for API information.")
}

// GetRouter returns the configured router
func (r *APIRouter) GetRouter() *gin.Engine {
	return r.router
}
