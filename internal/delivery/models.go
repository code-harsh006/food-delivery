package delivery

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeliveryAgent struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;uniqueIndex;not null"`
	LicenseNumber string        `json:"license_number" gorm:"uniqueIndex;not null"`
	VehicleType  string         `json:"vehicle_type" gorm:"not null"` // bike, car, bicycle
	VehicleNumber string        `json:"vehicle_number" gorm:"not null"`
	IsAvailable  bool           `json:"is_available" gorm:"default:true"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CurrentLat   float64        `json:"current_lat"`
	CurrentLng   float64        `json:"current_lng"`
	Rating       float64        `json:"rating" gorm:"default:0"`
	TotalDeliveries int         `json:"total_deliveries" gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Delivery struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID         uuid.UUID      `json:"order_id" gorm:"type:uuid;uniqueIndex;not null"`
	AgentID         *uuid.UUID     `json:"agent_id" gorm:"type:uuid"`
	Status          string         `json:"status" gorm:"default:'pending'"` // pending, assigned, picked_up, in_transit, delivered, cancelled
	PickupAddress   string         `json:"pickup_address" gorm:"type:jsonb"`
	DeliveryAddress string         `json:"delivery_address" gorm:"type:jsonb"`
	PickupTime      *time.Time     `json:"pickup_time"`
	DeliveryTime    *time.Time     `json:"delivery_time"`
	EstimatedTime   int            `json:"estimated_time"` // in minutes
	Distance        float64        `json:"distance"` // in km
	DeliveryFee     float64        `json:"delivery_fee"`
	Instructions    string         `json:"instructions"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	
	Agent *DeliveryAgent `json:"agent,omitempty" gorm:"foreignKey:AgentID"`
}

type DeliveryTracking struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	DeliveryID uuid.UUID `json:"delivery_id" gorm:"type:uuid;not null"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

func (da *DeliveryAgent) BeforeCreate(tx *gorm.DB) error {
	if da.ID == uuid.Nil {
		da.ID = uuid.New()
	}
	return nil
}

func (d *Delivery) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

func (dt *DeliveryTracking) BeforeCreate(tx *gorm.DB) error {
	if dt.ID == uuid.Nil {
		dt.ID = uuid.New()
	}
	return nil
}
