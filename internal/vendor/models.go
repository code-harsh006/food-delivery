package vendor

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vendor struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID      `json:"user_id" gorm:"type:uuid;uniqueIndex;not null"`
	BusinessName    string         `json:"business_name" gorm:"not null"`
	BusinessType    string         `json:"business_type" gorm:"not null"` // restaurant, grocery, pharmacy
	Description     string         `json:"description"`
	Logo            string         `json:"logo"`
	Banner          string         `json:"banner"`
	Address         string         `json:"address" gorm:"not null"`
	City            string         `json:"city" gorm:"not null"`
	State           string         `json:"state" gorm:"not null"`
	PostalCode      string         `json:"postal_code" gorm:"not null"`
	Latitude        float64        `json:"latitude"`
	Longitude       float64        `json:"longitude"`
	Phone           string         `json:"phone" gorm:"not null"`
	Email           string         `json:"email" gorm:"not null"`
	Website         string         `json:"website"`
	LicenseNumber   string         `json:"license_number" gorm:"uniqueIndex"`
	TaxNumber       string         `json:"tax_number"`
	IsVerified      bool           `json:"is_verified" gorm:"default:false"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	IsOnline        bool           `json:"is_online" gorm:"default:false"`
	Rating          float64        `json:"rating" gorm:"default:0"`
	ReviewCount     int            `json:"review_count" gorm:"default:0"`
	DeliveryRadius  float64        `json:"delivery_radius" gorm:"default:5"` // in km
	MinOrderAmount  float64        `json:"min_order_amount" gorm:"default:0"`
	DeliveryFee     float64        `json:"delivery_fee" gorm:"default:0"`
	DeliveryTime    int            `json:"delivery_time" gorm:"default:30"` // in minutes
	OpeningHours    string         `json:"opening_hours" gorm:"type:jsonb"`
	PaymentMethods  []string       `json:"payment_methods" gorm:"type:jsonb"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type VendorReview struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	VendorID  uuid.UUID      `json:"vendor_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	OrderID   uuid.UUID      `json:"order_id" gorm:"type:uuid;not null"`
	Rating    int            `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment   string         `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (v *Vendor) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

func (vr *VendorReview) BeforeCreate(tx *gorm.DB) error {
	if vr.ID == uuid.Nil {
		vr.ID = uuid.New()
	}
	return nil
}
