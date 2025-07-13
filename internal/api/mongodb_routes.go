package api

import (
	"net/http"

	"github.com/code-harsh006/food-delivery/internal/services"
	"github.com/gin-gonic/gin"
)

// SetupMongoDBRoutes configures MongoDB-specific API routes
func (r *APIRouter) SetupMongoDBRoutes() {
	// Health check
	r.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "database": "mongodb"})
	})

	// MongoDB API v1 routes (using different prefix to avoid conflicts)
	mongoV1 := r.router.Group("/api/mongo/v1")
	{
		// Authentication routes
		auth := mongoV1.Group("/auth")
		{
			auth.POST("/register", services.Register)
			auth.POST("/login", services.Login)
			auth.POST("/verify-otp", services.VerifyOTP)
			auth.POST("/resend-otp", services.ResendOTP)
		}

		// Service routes
		serviceRoutes := mongoV1.Group("/services")
		{
			serviceRoutes.GET("", services.GetServices)
			serviceRoutes.GET("/:id", services.GetServiceByID)
			serviceRoutes.GET("/categories", services.GetServiceCategories)
			serviceRoutes.GET("/search", services.SearchServices)
		}

		// Booking routes
		bookings := mongoV1.Group("/bookings")
		{
			bookings.POST("", services.CreateBooking)
			bookings.GET("", services.GetUserBookings)
			bookings.GET("/:id", services.GetBookingByID)
			bookings.PUT("/:id", services.UpdateBooking)
			bookings.DELETE("/:id", services.CancelBooking)
		}

		// User routes
		users := mongoV1.Group("/users")
		{
			users.GET("/profile", services.GetUserProfile)
			users.PUT("/profile", services.UpdateUserProfile)
			users.GET("/notifications", services.GetUserNotifications)
			users.PUT("/notifications/:id/read", services.MarkNotificationAsRead)
		}

		// Admin routes (for service providers/admin panel)
		admin := mongoV1.Group("/admin")
		{
			admin.GET("/bookings", services.GetAllBookings)
			admin.PUT("/bookings/:id/status", services.UpdateBookingStatus)
			admin.GET("/dashboard", services.GetDashboardStats)
		}
	}
}
