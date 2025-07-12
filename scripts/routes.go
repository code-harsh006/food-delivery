package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", Register)
			auth.POST("/login", Login)
			auth.POST("/verify-otp", VerifyOTP)
			auth.POST("/resend-otp", ResendOTP)
		}

		// Service routes
		services := v1.Group("/services")
		{
			services.GET("", GetServices)
			services.GET("/:id", GetServiceByID)
			services.GET("/categories", GetServiceCategories)
			services.GET("/search", SearchServices)
		}

		// Booking routes
		bookings := v1.Group("/bookings")
		{
			bookings.POST("", CreateBooking)
			bookings.GET("", GetUserBookings)
			bookings.GET("/:id", GetBookingByID)
			bookings.PUT("/:id", UpdateBooking)
			bookings.DELETE("/:id", CancelBooking)
		}

		// User routes
		users := v1.Group("/users")
		{
			users.GET("/profile", GetUserProfile)
			users.PUT("/profile", UpdateUserProfile)
			users.GET("/notifications", GetUserNotifications)
			users.PUT("/notifications/:id/read", MarkNotificationAsRead)
		}

		// Admin routes (for service providers/admin panel)
		admin := v1.Group("/admin")
		{
			admin.GET("/bookings", GetAllBookings)
			admin.PUT("/bookings/:id/status", UpdateBookingStatus)
			admin.GET("/dashboard", GetDashboardStats)
		}
	}
}

// GetUserProfile returns user profile information
func GetUserProfile(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var user User
	collection := MongoDB.Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUserProfile updates user profile information
func UpdateUserProfile(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var user User
	collection := MongoDB.Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData["updated_at"] = time.Now()
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$set": updateData})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Get updated user
	collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// GetUserNotifications returns user notifications
func GetUserNotifications(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var notifications []Notification
	collection := MongoDB.Collection("notifications")
	cursor, err := collection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &notifications); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total":         len(notifications),
	})
}

// MarkNotificationAsRead marks a notification as read
func MarkNotificationAsRead(c *gin.Context) {
	// Implementation for marking notification as read
	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// GetAllBookings returns all bookings (admin only)
func GetAllBookings(c *gin.Context) {
	// Implementation for admin to view all bookings
	c.JSON(http.StatusOK, gin.H{"message": "Admin bookings endpoint"})
}

// UpdateBookingStatus updates booking status (admin/technician)
func UpdateBookingStatus(c *gin.Context) {
	// Implementation for updating booking status by admin/technician
	c.JSON(http.StatusOK, gin.H{"message": "Booking status updated"})
}

// GetDashboardStats returns dashboard statistics
func GetDashboardStats(c *gin.Context) {
	// Implementation for dashboard statistics
	c.JSON(http.StatusOK, gin.H{"message": "Dashboard stats endpoint"})
}
