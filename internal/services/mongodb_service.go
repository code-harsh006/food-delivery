package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/code-harsh006/food-delivery/internal/models"
	"github.com/code-harsh006/food-delivery/pkg/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetServices returns all active services
func GetServices(c *gin.Context) {
	// Check if MongoDB is connected
	mongoDB := db.GetMongoDB()
	if mongoDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Database not available",
			"message": "MongoDB connection is not established",
		})
		return
	}

	var services []models.Service

	// Get query parameters for filtering
	category := c.Query("category")

	collection := mongoDB.Collection("services")
	filter := bson.M{"is_active": true}
	if category != "" {
		filter["category"] = category
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &services); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode services"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
		"total":    len(services),
	})
}

// GetServiceByID returns a specific service by ID
func GetServiceByID(c *gin.Context) {
	// Check if MongoDB is connected
	mongoDB := db.GetMongoDB()
	if mongoDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Database not available",
			"message": "MongoDB connection is not established",
		})
		return
	}

	idParam := c.Param("id")
	serviceID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	var service models.Service
	collection := mongoDB.Collection("services")
	err = collection.FindOne(context.Background(), bson.M{"_id": serviceID, "is_active": true}).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"service": service})
}

// GetServiceCategories returns all service categories
func GetServiceCategories(c *gin.Context) {
	// Check if MongoDB is connected
	mongoDB := db.GetMongoDB()
	if mongoDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Database not available",
			"message": "MongoDB connection is not established",
		})
		return
	}

	collection := mongoDB.Collection("services")

	// Use aggregation to get distinct categories
	pipeline := []bson.M{
		{"$match": bson.M{"is_active": true}},
		{"$group": bson.M{"_id": "$category"}},
		{"$project": bson.M{"category": "$_id", "_id": 0}},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode categories"})
		return
	}

	var categories []string
	for _, result := range results {
		if category, ok := result["category"].(string); ok {
			categories = append(categories, category)
		}
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// SearchServices searches services by name or description
func SearchServices(c *gin.Context) {
	// Check if MongoDB is connected
	mongoDB := db.GetMongoDB()
	if mongoDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Database not available",
			"message": "MongoDB connection is not established",
		})
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var services []models.Service
	collection := mongoDB.Collection("services")

	// Create a text search filter
	filter := bson.M{
		"is_active": true,
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search services"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &services); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode services"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
		"total":    len(services),
		"query":    query,
	})
}

// Register handles user registration
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	collection := db.GetMongoDB().Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Create new user
	user := models.User{
		Email:             req.Email,
		Phone:             req.Phone,
		Name:              req.Name,
		AccommodationType: req.AccommodationType,
		Address:           req.Address,
		IsVerified:        false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	// Generate and send OTP
	if err := generateAndSendOTP(user.ID, user.Email, "registration"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please verify your email with the OTP sent.",
		"user_id": user.ID.Hex(),
	})
}

// Login handles user login
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var user models.User
	collection := db.GetMongoDB().Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate and send OTP for login
	if err := generateAndSendOTP(user.ID, user.Email, "login"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to your email for login verification",
		"user_id": user.ID.Hex(),
	})
}

// VerifyOTP handles OTP verification
func VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	userCollection := db.GetMongoDB().Collection("users")
	err := userCollection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify OTP
	var otp models.OTP
	otpCollection := db.GetMongoDB().Collection("otps")
	err = otpCollection.FindOne(context.Background(), bson.M{
		"user_id":    user.ID,
		"code":       req.Code,
		"is_used":    false,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&otp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Mark OTP as used
	otpCollection.UpdateOne(context.Background(), bson.M{"_id": otp.ID}, bson.M{"$set": bson.M{"is_used": true}})

	// Mark user as verified if it's registration OTP
	if otp.Purpose == "registration" {
		userCollection.UpdateOne(context.Background(), bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"is_verified": true}})
		user.IsVerified = true
	}

	// Generate session token (simplified - in production use JWT)
	sessionToken := generateSessionToken()

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verified successfully",
		"user": gin.H{
			"id":                 user.ID.Hex(),
			"email":              user.Email,
			"name":               user.Name,
			"accommodation_type": user.AccommodationType,
			"is_verified":        user.IsVerified,
		},
		"session_token": sessionToken,
	})
}

// ResendOTP handles OTP resend requests
func ResendOTP(c *gin.Context) {
	email := c.Query("email")
	purpose := c.Query("purpose")

	if email == "" || purpose == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and purpose are required"})
		return
	}

	// Find user
	var user models.User
	collection := db.GetMongoDB().Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate and send new OTP
	if err := generateAndSendOTP(user.ID, user.Email, purpose); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP resent successfully",
	})
}

// generateAndSendOTP generates and sends OTP to user
func generateAndSendOTP(userID primitive.ObjectID, email, purpose string) error {
	// Generate 6-digit OTP
	code, err := generateOTP()
	if err != nil {
		return err
	}

	// Save OTP to database
	otp := models.OTP{
		UserID:    userID,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(10 * time.Minute), // 10 minutes expiry
		IsUsed:    false,
		CreatedAt: time.Now(),
	}

	collection := db.GetMongoDB().Collection("otps")
	_, err = collection.InsertOne(context.Background(), otp)
	if err != nil {
		return err
	}

	// Send OTP via email (implement your email service here)
	fmt.Printf("OTP for %s (%s): %s\n", email, purpose, code)

	return nil
}

// generateOTP generates a 6-digit OTP
func generateOTP() (string, error) {
	max := big.NewInt(999999)
	min := big.NewInt(100000)

	n, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, big.NewInt(1)))
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(n.Int64()) + 100000), nil
}

// generateSessionToken generates a simple session token
func generateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
