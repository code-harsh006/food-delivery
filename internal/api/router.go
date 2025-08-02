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
		"base_url":    "http://localhost:8080",
		"endpoints": gin.H{
			"root": gin.H{
				"method":       "GET",
				"path":         "/",
				"description":  "API information and available endpoints",
				"request_body": "None",
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"message": "Food Delivery API",
						"version": "1.0.0",
						"status":  "running",
						"endpoints": gin.H{
							"health": "/api/v1/health",
							"docs":   "/api/v1/docs",
							"status": "/api/v1/status",
						},
						"documentation": "Visit /api/v1/docs for API documentation",
					},
				},
			},
			"health": gin.H{
				"method":       "GET",
				"path":         "/api/v1/health",
				"description":  "Basic health check",
				"request_body": "None",
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"status":    "ok",
						"message":   "Food Delivery API is healthy",
						"timestamp": "2025-07-17T04:44:44.123Z",
						"service":   "food-delivery-api",
					},
				},
			},
			"health_detailed": gin.H{
				"method":       "GET",
				"path":         "/api/v1/health/detailed",
				"description":  "Detailed health check with system information",
				"request_body": "None",
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"status":    "ok",
						"timestamp": "2025-07-17T04:44:44.123Z",
						"service":   "food-delivery-api",
						"checks": gin.H{
							"api": gin.H{
								"status":  "ok",
								"message": "API is responsive",
								"latency": "0ms",
							},
							"memory": gin.H{
								"status":  "ok",
								"message": "Memory usage is normal",
							},
						},
					},
				},
			},
			"health_ready": gin.H{
				"method":       "GET",
				"path":         "/api/v1/health/ready",
				"description":  "Readiness check for load balancers",
				"request_body": "None",
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"ready":     true,
						"timestamp": "2025-07-17T04:44:44.123Z",
						"service":   "food-delivery-api",
					},
				},
			},
			"health_live": gin.H{
				"method":       "GET",
				"path":         "/api/v1/health/live",
				"description":  "Liveness check for container orchestration",
				"request_body": "None",
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"alive":     true,
						"timestamp": "2025-07-17T04:44:44.123Z",
						"service":   "food-delivery-api",
						"uptime":    "running",
					},
				},
			},
			"status": gin.H{
				"method":       "GET",
				"path":         "/api/v1/status",
				"description":  "API status and version information",
				"request_body": "None",
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"status":    "running",
						"message":   "Food Delivery API is operational",
						"timestamp": "2025-07-17T04:44:44.123Z",
						"version":   "1.0.0",
						"endpoints": gin.H{
							"health":          "/api/v1/health",
							"health_detailed": "/api/v1/health/detailed",
							"ready":           "/api/v1/health/ready",
							"live":            "/api/v1/health/live",
							"status":          "/api/v1/status",
						},
					},
				},
			},
			"docs": gin.H{
				"method":           "GET",
				"path":             "/api/v1/docs",
				"description":      "This documentation page",
				"request_body":     "None",
				"response_example": "This documentation",
			},
			"notifications_send": gin.H{
				"method":         "POST",
				"path":           "/api/v1/notifications/send",
				"description":    "Send a notification",
				"authentication": "Required (JWT)",
				"request_body": gin.H{
					"type":      "email|sms|push",
					"recipient": "user@example.com",
					"subject":   "Notification Subject",
					"message":   "Notification message content",
					"data": gin.H{
						"key": "value",
					},
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"message":         "Notification sent successfully",
						"notification_id": "notif_123456",
						"status":          "sent",
					},
				},
			},
			"notifications_subscribe": gin.H{
				"method":           "GET",
				"path":             "/api/v1/notifications/subscribe",
				"description":      "Subscribe to notifications (WebSocket)",
				"authentication":   "Required (JWT)",
				"request_body":     "None (WebSocket connection)",
				"response_example": "WebSocket connection established",
			},
			"mongo_health": gin.H{
				"method":       "GET",
				"path":         "/health",
				"description":  "MongoDB health check",
				"request_body": "None",
				"response_example": gin.H{
					"status":   "healthy",
					"database": "mongodb",
				},
			},
			"auth_register": gin.H{
				"method":         "POST",
				"path":           "/api/mongo/v1/auth/register",
				"description":    "Register a new user",
				"authentication": "Not required",
				"request_body": gin.H{
					"name":     "John Doe",
					"email":    "john@example.com",
					"phone":    "+1234567890",
					"password": "securepassword123",
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"message":               "User registered successfully",
						"user_id":               "user_123456",
						"email":                 "john@example.com",
						"verification_required": true,
					},
				},
			},
			"auth_login": gin.H{
				"method":         "POST",
				"path":           "/api/mongo/v1/auth/login",
				"description":    "User login",
				"authentication": "Not required",
				"request_body": gin.H{
					"email":    "john@example.com",
					"password": "securepassword123",
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"message": "Login successful",
						"token":   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
						"user": gin.H{
							"id":    "user_123456",
							"name":  "John Doe",
							"email": "john@example.com",
							"phone": "+1234567890",
						},
					},
				},
			},
			"auth_verify_otp": gin.H{
				"method":         "POST",
				"path":           "/api/mongo/v1/auth/verify-otp",
				"description":    "Verify OTP for email verification",
				"authentication": "Not required",
				"request_body": gin.H{
					"email": "john@example.com",
					"otp":   "123456",
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"message": "Email verified successfully",
						"user_id": "user_123456",
					},
				},
			},
			"auth_resend_otp": gin.H{
				"method":         "POST",
				"path":           "/api/mongo/v1/auth/resend-otp",
				"description":    "Resend OTP for email verification",
				"authentication": "Not required",
				"request_body": gin.H{
					"email": "john@example.com",
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"message": "OTP sent successfully",
						"email":   "john@example.com",
					},
				},
			},
			"services_list": gin.H{
				"method":         "GET",
				"path":           "/api/mongo/v1/services",
				"description":    "Get all available services",
				"authentication": "Not required",
				"request_body":   "None",
				"query_parameters": gin.H{
					"category": "food|cleaning|maintenance",
					"page":     "1",
					"limit":    "10",
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"services": []gin.H{
							gin.H{
								"id":          "service_123",
								"name":        "Food Delivery",
								"description": "Fast food delivery service",
								"price":       15.99,
								"category":    "food",
								"rating":      4.5,
								"image_url":   "https://example.com/food.jpg",
							},
						},
						"total": 1,
						"page":  1,
						"limit": 10,
					},
				},
			},
			"service_detail": gin.H{
				"method":         "GET",
				"path":           "/api/mongo/v1/services/:id",
				"description":    "Get service details by ID",
				"authentication": "Not required",
				"request_body":   "None",
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"id":          "service_123",
						"name":        "Food Delivery",
						"description": "Fast food delivery service",
						"price":       15.99,
						"category":    "food",
						"rating":      4.5,
						"image_url":   "https://example.com/food.jpg",
						"provider": gin.H{
							"id":     "provider_456",
							"name":   "Food Express",
							"rating": 4.8,
						},
					},
				},
			},
			"services_categories": gin.H{
				"method":         "GET",
				"path":           "/api/mongo/v1/services/categories",
				"description":    "Get all service categories",
				"authentication": "Not required",
				"request_body":   "None",
				"response_example": gin.H{
					"success": true,
					"data": []gin.H{
						gin.H{
							"id":          "cat_1",
							"name":        "Food",
							"description": "Food delivery services",
							"icon":        "🍕",
						},
						gin.H{
							"id":          "cat_2",
							"name":        "Cleaning",
							"description": "Cleaning services",
							"icon":        "🧹",
						},
					},
				},
			},
			"services_search": gin.H{
				"method":         "GET",
				"path":           "/api/mongo/v1/services/search",
				"description":    "Search services",
				"authentication": "Not required",
				"request_body":   "None",
				"query_parameters": gin.H{
					"q":         "food delivery",
					"category":  "food",
					"min_price": "10",
					"max_price": "50",
					"rating":    "4.0",
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"services": []gin.H{
							gin.H{
								"id":          "service_123",
								"name":        "Food Delivery",
								"description": "Fast food delivery service",
								"price":       15.99,
								"rating":      4.5,
							},
						},
						"total": 1,
					},
				},
			},
			"booking_create": gin.H{
				"method":         "POST",
				"path":           "/api/mongo/v1/bookings",
				"description":    "Create a new booking",
				"authentication": "Required (JWT)",
				"request_body": gin.H{
					"service_id":     "service_123",
					"scheduled_date": "2025-07-20T14:00:00Z",
					"address": gin.H{
						"street":   "123 Main St",
						"city":     "New York",
						"state":    "NY",
						"zip_code": "10001",
					},
					"special_instructions": "Please deliver to front door",
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"booking_id":     "booking_789",
						"message":        "Booking created successfully",
						"status":         "pending",
						"scheduled_date": "2025-07-20T14:00:00Z",
					},
				},
			},
			"bookings_list": gin.H{
				"method":         "GET",
				"path":           "/api/mongo/v1/bookings",
				"description":    "Get user's bookings",
				"authentication": "Required (JWT)",
				"request_body":   "None",
				"query_parameters": gin.H{
					"status": "pending|confirmed|completed|cancelled",
					"page":   "1",
					"limit":  "10",
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"bookings": []gin.H{
							gin.H{
								"id": "booking_789",
								"service": gin.H{
									"id":    "service_123",
									"name":  "Food Delivery",
									"price": 15.99,
								},
								"status":         "pending",
								"scheduled_date": "2025-07-20T14:00:00Z",
								"total_amount":   15.99,
							},
						},
						"total": 1,
						"page":  1,
						"limit": 10,
					},
				},
			},
			"booking_detail": gin.H{
				"method":         "GET",
				"path":           "/api/mongo/v1/bookings/:id",
				"description":    "Get booking details by ID",
				"authentication": "Required (JWT)",
				"request_body":   "None",
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"id": "booking_789",
						"service": gin.H{
							"id":          "service_123",
							"name":        "Food Delivery",
							"description": "Fast food delivery service",
							"price":       15.99,
						},
						"status":         "pending",
						"scheduled_date": "2025-07-20T14:00:00Z",
						"address": gin.H{
							"street":   "123 Main St",
							"city":     "New York",
							"state":    "NY",
							"zip_code": "10001",
						},
						"total_amount":         15.99,
						"special_instructions": "Please deliver to front door",
					},
				},
			},
			"booking_update": gin.H{
				"method":         "PUT",
				"path":           "/api/mongo/v1/bookings/:id",
				"description":    "Update booking details",
				"authentication": "Required (JWT)",
				"request_body": gin.H{
					"scheduled_date":       "2025-07-21T15:00:00Z",
					"special_instructions": "Updated delivery instructions",
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"message":    "Booking updated successfully",
						"booking_id": "booking_789",
					},
				},
			},
			"booking_cancel": gin.H{
				"method":         "DELETE",
				"path":           "/api/mongo/v1/bookings/:id",
				"description":    "Cancel a booking",
				"authentication": "Required (JWT)",
				"request_body":   "None",
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"message":    "Booking cancelled successfully",
						"booking_id": "booking_789",
					},
				},
			},
			"user_profile": gin.H{
				"method":         "GET",
				"path":           "/api/mongo/v1/users/profile",
				"description":    "Get user profile",
				"authentication": "Required (JWT)",
				"request_body":   "None",
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"id":    "user_123456",
						"name":  "John Doe",
						"email": "john@example.com",
						"phone": "+1234567890",
						"addresses": []gin.H{
							gin.H{
								"id":         "addr_1",
								"title":      "Home",
								"street":     "123 Main St",
								"city":       "New York",
								"state":      "NY",
								"zip_code":   "10001",
								"is_default": true,
							},
						},
					},
				},
			},
			"user_profile_update": gin.H{
				"method":         "PUT",
				"path":           "/api/mongo/v1/users/profile",
				"description":    "Update user profile",
				"authentication": "Required (JWT)",
				"request_body": gin.H{
					"name":  "John Smith",
					"phone": "+1234567890",
					"addresses": []gin.H{
						gin.H{
							"title":      "Home",
							"street":     "123 Main St",
							"city":       "New York",
							"state":      "NY",
							"zip_code":   "10001",
							"is_default": true,
						},
					},
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"message": "Profile updated successfully",
						"user_id": "user_123456",
					},
				},
			},
			"user_notifications": gin.H{
				"method":         "GET",
				"path":           "/api/mongo/v1/users/notifications",
				"description":    "Get user notifications",
				"authentication": "Required (JWT)",
				"request_body":   "None",
				"query_parameters": gin.H{
					"page":  "1",
					"limit": "10",
					"read":  "true|false",
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"notifications": []gin.H{
							gin.H{
								"id":         "notif_123",
								"type":       "booking_confirmation",
								"title":      "Booking Confirmed",
								"message":    "Your booking has been confirmed",
								"read":       false,
								"created_at": "2025-07-17T04:44:44.123Z",
							},
						},
						"total":        1,
						"unread_count": 1,
					},
				},
			},
			"notification_mark_read": gin.H{
				"method":         "PUT",
				"path":           "/api/mongo/v1/users/notifications/:id/read",
				"description":    "Mark notification as read",
				"authentication": "Required (JWT)",
				"request_body":   "None",
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"message":         "Notification marked as read",
						"notification_id": "notif_123",
					},
				},
			},
			"admin_bookings": gin.H{
				"method":         "GET",
				"path":           "/api/mongo/v1/admin/bookings",
				"description":    "Get all bookings (Admin only)",
				"authentication": "Required (JWT + Admin)",
				"request_body":   "None",
				"query_parameters": gin.H{
					"status":  "pending|confirmed|completed|cancelled",
					"page":    "1",
					"limit":   "10",
					"user_id": "user_123456",
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"bookings": []gin.H{
							gin.H{
								"id": "booking_789",
								"user": gin.H{
									"id":    "user_123456",
									"name":  "John Doe",
									"email": "john@example.com",
								},
								"service": gin.H{
									"id":   "service_123",
									"name": "Food Delivery",
								},
								"status":         "pending",
								"scheduled_date": "2025-07-20T14:00:00Z",
								"total_amount":   15.99,
							},
						},
						"total": 1,
						"page":  1,
						"limit": 10,
					},
				},
			},
			"admin_booking_status": gin.H{
				"method":         "PUT",
				"path":           "/api/mongo/v1/admin/bookings/:id/status",
				"description":    "Update booking status (Admin only)",
				"authentication": "Required (JWT + Admin)",
				"request_body": gin.H{
					"status": "confirmed|completed|cancelled",
					"notes":  "Admin notes about the booking",
				},
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"message":    "Booking status updated successfully",
						"booking_id": "booking_789",
						"new_status": "confirmed",
					},
				},
			},
			"admin_dashboard": gin.H{
				"method":         "GET",
				"path":           "/api/mongo/v1/admin/dashboard",
				"description":    "Get admin dashboard statistics",
				"authentication": "Required (JWT + Admin)",
				"request_body":   "None",
				"response_example": gin.H{
					"success": true,
					"data": gin.H{
						"total_users":        150,
						"total_bookings":     45,
						"pending_bookings":   12,
						"completed_bookings": 28,
						"cancelled_bookings": 5,
						"total_revenue":      1250.75,
						"recent_bookings": []gin.H{
							gin.H{
								"id":           "booking_789",
								"user_name":    "John Doe",
								"service_name": "Food Delivery",
								"status":       "pending",
								"amount":       15.99,
							},
						},
					},
				},
			},
		},
		"authentication": "Most endpoints require JWT authentication",
		"error_responses": gin.H{
			"400": gin.H{
				"success": false,
				"error":   "Bad Request - Invalid request data",
			},
			"401": gin.H{
				"success": false,
				"error":   "Unauthorized - Authentication required",
			},
			"403": gin.H{
				"success": false,
				"error":   "Forbidden - Insufficient permissions",
			},
			"404": gin.H{
				"success": false,
				"error":   "Not Found - Resource not found",
			},
			"500": gin.H{
				"success": false,
				"error":   "Internal Server Error - Something went wrong",
			},
		},
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
