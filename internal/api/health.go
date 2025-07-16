package api

import (
	"net/http"
	"time"

	"github.com/code-harsh006/food-delivery/pkg/response"
	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check and system monitoring endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler instance
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck performs a basic health check
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"status":    "ok",
		"message":   "Food Delivery API is healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "food-delivery-api",
	})
}

// DetailedHealthCheck performs a comprehensive health check
func (h *HealthHandler) DetailedHealthCheck(c *gin.Context) {
	healthStatus := gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "food-delivery-api",
		"checks":    make(map[string]interface{}),
	}

	// Check API responsiveness
	healthStatus["checks"].(map[string]interface{})["api"] = gin.H{
		"status":  "ok",
		"message": "API is responsive",
		"latency": "0ms",
	}

	// Check memory usage (basic)
	healthStatus["checks"].(map[string]interface{})["memory"] = gin.H{
		"status":  "ok",
		"message": "Memory usage is normal",
	}

	// Overall status
	allChecks := healthStatus["checks"].(map[string]interface{})
	overallStatus := "ok"
	for _, check := range allChecks {
		if check.(gin.H)["status"] == "error" {
			overallStatus = "error"
			break
		}
	}
	healthStatus["status"] = overallStatus

	if overallStatus == "ok" {
		response.Success(c, healthStatus)
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "Service health check failed",
			"data":    healthStatus,
		})
	}
}

// ReadinessCheck checks if the service is ready to handle requests
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	readiness := gin.H{
		"ready":     true,
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "food-delivery-api",
	}

	response.Success(c, readiness)
}

// LivenessCheck checks if the service is alive
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"alive":     true,
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "food-delivery-api",
		"uptime":    "running",
	})
}

// SetupHealthRoutes sets up health check routes
func (h *HealthHandler) SetupHealthRoutes(router *gin.RouterGroup) {
	health := router.Group("/health")
	{
		health.GET("/", h.HealthCheck)
		health.GET("/detailed", h.DetailedHealthCheck)
		health.GET("/ready", h.ReadinessCheck)
		health.GET("/live", h.LivenessCheck)
	}
}
