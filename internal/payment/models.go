package payment

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID         uuid.UUID      `json:"order_id" gorm:"type:uuid;not null"`
	UserID          uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Amount          float64        `json:"amount" gorm:"not null"`
	Currency        string         `json:"currency" gorm:"default:'INR'"`
	PaymentMethod   string         `json:"payment_method" gorm:"not null"` // card, upi, wallet, cod
	PaymentGateway  string         `json:"payment_gateway"` // stripe, razorpay, paytm
	GatewayPaymentID string        `json:"gateway_payment_id"`
	Status          string         `json:"status" gorm:"default:'pending'"` // pending, processing, completed, failed, cancelled, refunded
	FailureReason   string         `json:"failure_reason"`
	RefundAmount    float64        `json:"refund_amount" gorm:"default:0"`
	RefundReason    string         `json:"refund_reason"`
	Metadata        string         `json:"metadata" gorm:"type:jsonb"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type PaymentMethod struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Type      string         `json:"type" gorm:"not null"` // card, upi, wallet
	Details   string         `json:"details" gorm:"type:jsonb"`
	IsDefault bool           `json:"is_default" gorm:"default:false"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (pm *PaymentMethod) BeforeCreate(tx *gorm.DB) error {
	if pm.ID == uuid.Nil {
		pm.ID = uuid.New()
	}
	return nil
}
