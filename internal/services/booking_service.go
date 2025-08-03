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

// CreateBooking handles booking creation
func CreateBooking(c *gin.Context) {
	// Check if MongoDB is connected
	mongoDB := db.GetMongoDB()
	if mongoDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Database not available",
			"message": "MongoDB connection is not established",
		})
		return
	}

	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (you'll need to implement authentication middleware)
	userID := getUserIDFromContext(c)
	if userID.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse service ID
	serviceID, err := primitive.ObjectIDFromHex(req.ServiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	// Verify service exists
	var service models.Service
	serviceCollection := mongoDB.Collection("services")
	err = serviceCollection.FindOne(context.Background(), bson.M{"_id": serviceID, "is_active": true}).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Parse scheduled date
	scheduledDate, err := time.Parse("2006-01-02", req.ScheduledDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Check if the date is not in the past
	if scheduledDate.Before(time.Now().Truncate(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot book services for past dates"})
		return
	}

	// Create booking
	booking := models.Booking{
		UserID:          userID,
		ServiceID:       serviceID,
		ScheduledDate:   scheduledDate,
		ScheduledTime:   req.ScheduledTime,
		Status:          "pending",
		TotalAmount:     service.BasePrice,
		PaymentStatus:   "pending",
		SpecialRequests: req.SpecialRequests,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	bookingCollection := mongoDB.Collection("bookings")
	result, err := bookingCollection.InsertOne(context.Background(), booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	booking.ID = result.InsertedID.(primitive.ObjectID)

	// Create initial booking status
	status := models.BookingStatus{
		BookingID: booking.ID,
		Status:    "pending",
		Message:   "Booking created successfully",
		UpdatedBy: "system",
		CreatedAt: time.Now(),
	}
	statusCollection := mongoDB.Collection("booking_statuses")
	statusCollection.InsertOne(context.Background(), status)

	// Send notification to user
	notification := models.Notification{
		UserID:    userID,
		Title:     "Booking Confirmed",
		Message:   "Your booking for " + service.Name + " has been created successfully",
		Type:      "booking",
		CreatedAt: time.Now(),
	}
	notificationCollection := mongoDB.Collection("notifications")
	notificationCollection.InsertOne(context.Background(), notification)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Booking created successfully",
		"booking": booking,
	})
}

// GetUserBookings returns all bookings for a user
func GetUserBookings(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var bookings []models.Booking
	collection := db.GetMongoDB().Collection("bookings")
	cursor, err := collection.Find(context.Background(), bson.M{"user_id": userID})
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
	})
}

// GetBookingByID returns a specific booking
func GetBookingByID(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	userID := getUserIDFromContext(c)
	if userID.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var booking models.Booking
	collection := db.GetMongoDB().Collection("bookings")
	err = collection.FindOne(context.Background(), bson.M{"_id": bookingID, "user_id": userID}).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get booking status history
	var statusHistory []models.BookingStatus
	statusCollection := db.GetMongoDB().Collection("booking_statuses")
	statusCursor, err := statusCollection.Find(context.Background(), bson.M{"booking_id": booking.ID})
	if err == nil {
		defer statusCursor.Close(context.Background())
		statusCursor.All(context.Background(), &statusHistory)
	}

	c.JSON(http.StatusOK, gin.H{
		"booking":        booking,
		"status_history": statusHistory,
	})
}

// UpdateBooking handles booking updates
func UpdateBooking(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	userID := getUserIDFromContext(c)
	if userID.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var req models.UpdateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var booking models.Booking
	collection := db.GetMongoDB().Collection("bookings")
	err = collection.FindOne(context.Background(), bson.M{"_id": bookingID, "user_id": userID}).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Update booking fields
	updateData := bson.M{"updated_at": time.Now()}
	if req.Status != "" {
		updateData["status"] = req.Status
	}
	if req.TechnicianNotes != "" {
		updateData["technician_notes"] = req.TechnicianNotes
	}
	if req.Rating > 0 && req.Rating <= 5 {
		updateData["rating"] = req.Rating
	}
	if req.Review != "" {
		updateData["review"] = req.Review
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": bookingID}, bson.M{"$set": updateData})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}

	// Get updated booking
	collection.FindOne(context.Background(), bson.M{"_id": bookingID}).Decode(&booking)

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking updated successfully",
		"booking": booking,
	})
}

// CancelBooking handles booking cancellation
func CancelBooking(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	userID := getUserIDFromContext(c)
	if userID.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var booking models.Booking
	collection := db.GetMongoDB().Collection("bookings")
	err = collection.FindOne(context.Background(), bson.M{"_id": bookingID, "user_id": userID}).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if booking can be cancelled
	if booking.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking is already cancelled"})
		return
	}

	if booking.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot cancel completed booking"})
		return
	}

	// Update booking status to cancelled
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": bookingID}, bson.M{
		"$set": bson.M{
			"status":     "cancelled",
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking"})
		return
	}

	// Create cancellation status
	status := models.BookingStatus{
		BookingID: bookingID,
		Status:    "cancelled",
		Message:   "Booking cancelled by user",
		UpdatedBy: "user",
		CreatedAt: time.Now(),
	}
	statusCollection := db.GetMongoDB().Collection("booking_statuses")
	statusCollection.InsertOne(context.Background(), status)

	// Send notification
	notification := models.Notification{
		UserID:    userID,
		Title:     "Booking Cancelled",
		Message:   "Your booking has been cancelled successfully",
		Type:      "booking",
		CreatedAt: time.Now(),
	}
	notificationCollection := db.GetMongoDB().Collection("notifications")
	notificationCollection.InsertOne(context.Background(), notification)

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking cancelled successfully",
	})
}

// getUserIDFromContext extracts user ID from context (placeholder - implement proper auth)
func getUserIDFromContext(c *gin.Context) primitive.ObjectID {
	// This is a placeholder - implement proper authentication middleware
	// For now, we'll use a header or query parameter
	userIDStr := c.GetHeader("User-ID")
	if userIDStr == "" {
		userIDStr = c.Query("user_id")
	}

	if userIDStr != "" {
		if userID, err := primitive.ObjectIDFromHex(userIDStr); err == nil {
			return userID
		}
	}

	return primitive.NilObjectID
}
