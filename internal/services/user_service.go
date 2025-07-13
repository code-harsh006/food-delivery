package services

import (
	"context"
	"net/http"
	"time"

	"github.com/code-harsh006/food-delivery/internal/models"
	"github.com/code-harsh006/food-delivery/pkg/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUserProfile returns user profile information
func GetUserProfile(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var user models.User
	collection := db.GetMongoDB().Collection("users")
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

	var user models.User
	collection := db.GetMongoDB().Collection("users")
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

	var notifications []models.Notification
	collection := db.GetMongoDB().Collection("notifications")
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
	notificationID := c.Param("id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID is required"})
		return
	}

	userID := getUserIDFromContext(c)
	if userID.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse notification ID
	notifID, err := primitive.ObjectIDFromHex(notificationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Update notification as read
	collection := db.GetMongoDB().Collection("notifications")
	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": notifID, "user_id": userID},
		bson.M{"$set": bson.M{"read": true}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// GetAllBookings returns all bookings (admin only)
func GetAllBookings(c *gin.Context) {
	// This would typically check for admin role
	// For now, we'll return all bookings with pagination
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	var bookings []models.Booking
	collection := db.GetMongoDB().Collection("bookings")

	// Add pagination logic here
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &bookings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bookings": bookings,
		"total":    len(bookings),
		"page":     page,
		"limit":    limit,
	})
}

// UpdateBookingStatus updates booking status (admin/technician)
func UpdateBookingStatus(c *gin.Context) {
	bookingID := c.Param("id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking ID is required"})
		return
	}

	// Parse booking ID
	bookID, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var req struct {
		Status  string `json:"status" binding:"required"`
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update booking status
	bookingCollection := db.GetMongoDB().Collection("bookings")
	result, err := bookingCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": bookID},
		bson.M{"$set": bson.M{
			"status":     req.Status,
			"updated_at": time.Now(),
		}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	// Create status history entry
	status := models.BookingStatus{
		BookingID: bookID,
		Status:    req.Status,
		Message:   req.Message,
		UpdatedBy: "admin", // This should come from auth context
		CreatedAt: time.Now(),
	}
	statusCollection := db.GetMongoDB().Collection("booking_statuses")
	statusCollection.InsertOne(context.Background(), status)

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking status updated successfully",
		"status":  req.Status,
	})
}

// GetDashboardStats returns dashboard statistics
func GetDashboardStats(c *gin.Context) {
	// This would typically check for admin role
	collection := db.GetMongoDB().Collection("bookings")

	// Get total bookings
	totalBookings, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get booking count"})
		return
	}

	// Get pending bookings
	pendingBookings, err := collection.CountDocuments(context.Background(), bson.M{"status": "pending"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending bookings"})
		return
	}

	// Get completed bookings
	completedBookings, err := collection.CountDocuments(context.Background(), bson.M{"status": "completed"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get completed bookings"})
		return
	}

	// Get cancelled bookings
	cancelledBookings, err := collection.CountDocuments(context.Background(), bson.M{"status": "cancelled"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cancelled bookings"})
		return
	}

	// Get total users
	userCollection := db.GetMongoDB().Collection("users")
	totalUsers, err := userCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"total_bookings":     totalBookings,
			"pending_bookings":   pendingBookings,
			"completed_bookings": completedBookings,
			"cancelled_bookings": cancelledBookings,
			"total_users":        totalUsers,
		},
	})
}
