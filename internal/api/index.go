package api

import (
	"net/http"
	"time"

	"github.com/code-harsh006/food-delivery/pkg/response"
	"github.com/gin-gonic/gin"
)

// IndexHandler handles the main API index and information endpoints
type IndexHandler struct {
	version   string
	startTime time.Time
}

// NewIndexHandler creates a new IndexHandler instance
func NewIndexHandler() *IndexHandler {
	return &IndexHandler{
		version:   "1.0.0",
		startTime: time.Now(),
	}
}

// GetAPIInfo returns general API information
func (h *IndexHandler) GetAPIInfo(c *gin.Context) {
	uptime := time.Since(h.startTime)

	response.Success(c, gin.H{
		"name":        "Food Delivery API",
		"version":     h.version,
		"description": "A comprehensive food delivery platform API",
		"uptime":      uptime.String(),
		"started_at":  h.startTime.Format(time.RFC3339),
		"endpoints": gin.H{
			"auth": gin.H{
				"register": "POST /api/v1/auth/register",
				"login":    "POST /api/v1/auth/login",
				"logout":   "POST /api/v1/auth/logout",
				"refresh":  "POST /api/v1/auth/refresh",
				"profile":  "GET /api/v1/auth/profile",
			},
			"users": gin.H{
				"profile":   "GET /api/v1/users/profile",
				"update":    "PUT /api/v1/users/profile",
				"addresses": "GET /api/v1/users/addresses",
			},
			"products": gin.H{
				"list":   "GET /api/v1/products",
				"detail": "GET /api/v1/products/:id",
				"search": "GET /api/v1/search/products",
			},
			"cart": gin.H{
				"view":   "GET /api/v1/cart",
				"add":    "POST /api/v1/cart/items",
				"update": "PUT /api/v1/cart/items/:id",
				"remove": "DELETE /api/v1/cart/items/:id",
				"clear":  "DELETE /api/v1/cart",
			},
			"orders": gin.H{
				"create": "POST /api/v1/orders",
				"list":   "GET /api/v1/orders",
				"detail": "GET /api/v1/orders/:id",
				"cancel": "PUT /api/v1/orders/:id/cancel",
			},
			"vendors": gin.H{
				"list":   "GET /api/v1/vendors",
				"detail": "GET /api/v1/vendors/:id",
				"search": "GET /api/v1/search/vendors",
			},
		},
		"status": "operational",
	})
}

// GetAPIVersion returns the API version information
func (h *IndexHandler) GetAPIVersion(c *gin.Context) {
	response.Success(c, gin.H{
		"version":     h.version,
		"build_date":  h.startTime.Format(time.RFC3339),
		"environment": "development",
	})
}

// GetAPIStatus returns the current API status
func (h *IndexHandler) GetAPIStatus(c *gin.Context) {
	uptime := time.Since(h.startTime)

	response.Success(c, gin.H{
		"status":    "operational",
		"uptime":    uptime.String(),
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   h.version,
		"services": gin.H{
			"api":      "operational",
			"database": "unknown", // This would be checked dynamically
			"redis":    "unknown", // This would be checked dynamically
		},
	})
}

// SetupIndexRoutes sets up the index routes
func (h *IndexHandler) SetupIndexRoutes(router *gin.RouterGroup) {
	// API information routes
	router.GET("/", h.GetAPIInfo)
	router.GET("/version", h.GetAPIVersion)
	router.GET("/status", h.GetAPIStatus)

	// API documentation route
	router.GET("/docs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Documentation",
			"swagger": "/api/v1/swagger/index.html",
			"postman": "https://documenter.getpostman.com/view/your-collection",
			"github":  "https://github.com/code-harsh006/food-delivery",
		})
	})
}
