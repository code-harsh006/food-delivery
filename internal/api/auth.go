package api

import (
	"net/http"

	"github.com/code-harsh006/food-delivery/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related API endpoints
type AuthHandler struct{}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	response.Success(c, gin.H{
		"user_id": userID,
		"message": "Profile retrieved successfully",
	})
}

// SetupAuthRoutes sets up authentication routes
func (h *AuthHandler) SetupAuthRoutes(router *gin.RouterGroup) {
	// Add auth-related routes
	auth := router.Group("/auth")
	{
		auth.GET("/profile", h.GetProfile)
	}
}
