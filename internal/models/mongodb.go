package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email             string             `bson:"email" json:"email"`
	Phone             string             `bson:"phone" json:"phone"`
	Name              string             `bson:"name" json:"name"`
	AccommodationType string             `bson:"accommodation_type" json:"accommodation_type"`
	Address           string             `bson:"address" json:"address"`
	IsVerified        bool               `bson:"is_verified" json:"is_verified"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}

// OTP represents one-time passwords for verification
type OTP struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Code      string             `bson:"code" json:"code"`
	Purpose   string             `bson:"purpose" json:"purpose"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expires_at"`
	IsUsed    bool               `bson:"is_used" json:"is_used"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// Service represents available services
type Service struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"`
	BasePrice   float64            `bson:"base_price" json:"base_price"`
	Duration    int                `bson:"duration" json:"duration"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Booking represents service bookings
type Booking struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	ServiceID       primitive.ObjectID `bson:"service_id" json:"service_id"`
	ScheduledDate   time.Time          `bson:"scheduled_date" json:"scheduled_date"`
	ScheduledTime   string             `bson:"scheduled_time" json:"scheduled_time"`
	Status          string             `bson:"status" json:"status"`
	TotalAmount     float64            `bson:"total_amount" json:"total_amount"`
	PaymentStatus   string             `bson:"payment_status" json:"payment_status"`
	SpecialRequests string             `bson:"special_requests" json:"special_requests"`
	TechnicianNotes string             `bson:"technician_notes" json:"technician_notes"`
	Rating          int                `bson:"rating" json:"rating"`
	Review          string             `bson:"review" json:"review"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// BookingStatus represents possible booking statuses
type BookingStatus struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookingID primitive.ObjectID `bson:"booking_id" json:"booking_id"`
	Status    string             `bson:"status" json:"status"`
	Message   string             `bson:"message" json:"message"`
	UpdatedBy string             `bson:"updated_by" json:"updated_by"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// Notification represents user notifications
type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title     string             `bson:"title" json:"title"`
	Message   string             `bson:"message" json:"message"`
	Type      string             `bson:"type" json:"type"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// Request/Response structures
type RegisterRequest struct {
	Email             string `json:"email" binding:"required,email"`
	Phone             string `json:"phone" binding:"required"`
	Name              string `json:"name" binding:"required"`
	AccommodationType string `json:"accommodation_type"`
	Address           string `json:"address"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type CreateBookingRequest struct {
	ServiceID       string `json:"service_id" binding:"required"`
	ScheduledDate   string `json:"scheduled_date" binding:"required"`
	ScheduledTime   string `json:"scheduled_time" binding:"required"`
	SpecialRequests string `json:"special_requests"`
}

type UpdateBookingRequest struct {
	Status          string `json:"status"`
	TechnicianNotes string `json:"technician_notes"`
	Rating          int    `json:"rating"`
	Review          string `json:"review"`
}
