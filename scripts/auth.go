package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email             string `json:"email" binding:"required,email"`
	Phone             string `json:"phone" binding:"required"`
	Name              string `json:"name" binding:"required"`
	AccommodationType string `json:"accommodation_type"`
	Address           string `json:"address"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyOTPRequest represents OTP verification request
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// Register handles user registration
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser User
	collection := MongoDB.Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Create new user
	user := User{
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
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var user User
	collection := MongoDB.Collection("users")
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
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user User
	userCollection := MongoDB.Collection("users")
	err := userCollection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify OTP
	var otp OTP
	otpCollection := MongoDB.Collection("otps")
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

// generateAndSendOTP generates and sends OTP to user
func generateAndSendOTP(userID primitive.ObjectID, email, purpose string) error {
	// Generate 6-digit OTP
	code, err := generateOTP()
	if err != nil {
		return err
	}

	// Save OTP to database
	otp := OTP{
		UserID:    userID,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(10 * time.Minute), // 10 minutes expiry
		IsUsed:    false,
		CreatedAt: time.Now(),
	}

	collection := MongoDB.Collection("otps")
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

// ResendOTP handles OTP resend requests
func ResendOTP(c *gin.Context) {
	email := c.Query("email")
	purpose := c.Query("purpose")

	if email == "" || purpose == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and purpose are required"})
		return
	}

	// Find user
	var user User
	collection := MongoDB.Collection("users")
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

	c.JSON(http.StatusOK, gin.H{"message": "OTP resent successfully"})
}
