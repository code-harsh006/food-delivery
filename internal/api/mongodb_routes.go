package api

import (
	"log"
	"net/http"

	"github.com/code-harsh006/food-delivery/internal/services"
	"github.com/code-harsh006/food-delivery/pkg/response"
	"github.com/gin-gonic/gin"
)

// SetupMongoDBRoutes configures MongoDB-specific API routes
func (r *APIRouter) SetupMongoDBRoutes() {
	log.Println("Setting up MongoDB routes...")

	// Health check
	r.router.GET("/health", func(c *gin.Context) {
		log.Println("Health check endpoint called")
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "database": "mongodb"})
	})

	// MongoDB API v1 routes (using different prefix to avoid conflicts)
	mongoV1 := r.router.Group("/api/mongo/v1")
	log.Println("Created MongoDB v1 group: /api/mongo/v1")

	{
		// Add a handler for the /api/mongo/v1 root path
		mongoV1.GET("", func(c *gin.Context) {
			log.Println("MongoDB API v1 root endpoint called")
			response.Success(c, gin.H{
				"message": "MongoDB API v1",
				"version": "1.0.0",
				"endpoints": gin.H{
					"auth":     "/api/mongo/v1/auth",
					"services": "/api/mongo/v1/services",
					"bookings": "/api/mongo/v1/bookings",
					"users":    "/api/mongo/v1/users",
					"admin":    "/api/mongo/v1/admin",
				},
				"description": "Food Delivery API with MongoDB backend",
				"health":      "/health",
			})
		})

		// Add a simple test route to verify registration
		mongoV1.GET("/test", func(c *gin.Context) {
			log.Println("MongoDB test endpoint called")
			c.JSON(http.StatusOK, gin.H{"message": "MongoDB routes are working"})
		})

		// Authentication routes
		auth := mongoV1.Group("/auth")
		log.Println("Created auth group: /api/mongo/v1/auth")

		{
			// Add a handler for the /auth root path
			auth.GET("", func(c *gin.Context) {
				log.Println("Auth root endpoint called")
				response.Success(c, gin.H{
					"message": "Authentication endpoints",
					"endpoints": gin.H{
						"register":   "POST /api/mongo/v1/auth/register",
						"login":      "POST /api/mongo/v1/auth/login",
						"verify_otp": "POST /api/mongo/v1/auth/verify-otp",
						"resend_otp": "POST /api/mongo/v1/auth/resend-otp",
					},
					"description": "Use these endpoints for user authentication and management",
				})
			})
			auth.POST("/register", services.Register)
			auth.POST("/login", services.Login)
			auth.POST("/verify-otp", services.VerifyOTP)
			auth.POST("/resend-otp", services.ResendOTP)
			log.Println("Registered auth endpoints")
		}

		// Service routes
		serviceRoutes := mongoV1.Group("/services")
		log.Println("Created services group: /api/mongo/v1/services")

		{
			// Add a handler for the /services root path
			serviceRoutes.GET("", services.GetServices)
			serviceRoutes.GET("/categories", services.GetServiceCategories)
			serviceRoutes.GET("/search", services.SearchServices)
			serviceRoutes.GET("/:id", services.GetServiceByID)
			log.Println("Registered service endpoints")
		}

		// Booking routes
		bookings := mongoV1.Group("/bookings")
		log.Println("Created bookings group: /api/mongo/v1/bookings")

		{
			// Add a handler for the /bookings root path
			bookings.GET("", services.GetUserBookings)
			bookings.POST("", services.CreateBooking)
			bookings.GET("/:id", services.GetBookingByID)
			bookings.PUT("/:id", services.UpdateBooking)
			bookings.DELETE("/:id", services.CancelBooking)
			log.Println("Registered booking endpoints")
		}

		// User routes
		users := mongoV1.Group("/users")
		log.Println("Created users group: /api/mongo/v1/users")

		{
			// Add a handler for the /users root path
			users.GET("", func(c *gin.Context) {
				log.Println("Users root endpoint called")
				response.Success(c, gin.H{
					"message": "User management endpoints",
					"endpoints": gin.H{
						"profile":        "GET /api/mongo/v1/users/profile",
						"update_profile": "PUT /api/mongo/v1/users/profile",
						"notifications":  "GET /api/mongo/v1/users/notifications",
						"mark_read":      "PUT /api/mongo/v1/users/notifications/:id/read",
					},
					"description": "Use these endpoints for user profile and notification management",
				})
			})
			users.GET("/profile", services.GetUserProfile)
			users.PUT("/profile", services.UpdateUserProfile)
			users.GET("/notifications", services.GetUserNotifications)
			users.PUT("/notifications/:id/read", services.MarkNotificationAsRead)
			log.Println("Registered user endpoints")
		}

		// Admin routes (for service providers/admin panel)
		admin := mongoV1.Group("/admin")
		log.Println("Created admin group: /api/mongo/v1/admin")

		{
			// Add a handler for the /admin root path
			admin.GET("", func(c *gin.Context) {
				log.Println("Admin root endpoint called")
				response.Success(c, gin.H{
					"message": "Admin management endpoints",
					"endpoints": gin.H{
						"bookings":      "GET /api/mongo/v1/admin/bookings",
						"update_status": "PUT /api/mongo/v1/admin/bookings/:id/status",
						"dashboard":     "GET /api/mongo/v1/admin/dashboard",
					},
					"description": "Use these endpoints for admin panel functionality (requires admin privileges)",
				})
			})
			admin.GET("/bookings", services.GetAllBookings)
			admin.PUT("/bookings/:id/status", services.UpdateBookingStatus)
			admin.GET("/dashboard", services.GetDashboardStats)
			log.Println("Registered admin endpoints")
		}
	}

	log.Println("MongoDB routes setup completed")
}
