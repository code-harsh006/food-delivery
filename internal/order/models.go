package order

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	OrderNumber     string         `json:"order_number" gorm:"uniqueIndex;not null"`
	Status          string         `json:"status" gorm:"default:'pending'"` // pending, confirmed, preparing, ready, picked_up, delivered, cancelled
	PaymentStatus   string         `json:"payment_status" gorm:"default:'pending'"` // pending, paid, failed, refunded
	PaymentMethod   string         `json:"payment_method"`
	PaymentID       string         `json:"payment_id"`
	Items           []OrderItem    `json:"items" gorm:"foreignKey:OrderID"`
	SubTotal        float64        `json:"sub_total" gorm:"not null"`
	DeliveryFee     float64        `json:"delivery_fee" gorm:"default:0"`
	Tax             float64        `json:"tax" gorm:"default:0"`
	Discount        float64        `json:"discount" gorm:"default:0"`
	Total           float64        `json:"total" gorm:"not null"`
	DeliveryAddress string         `json:"delivery_address" gorm:"type:jsonb"`
	DeliveryTime    *time.Time     `json:"delivery_time"`
	EstimatedTime   int            `json:"estimated_time"` // in minutes
	SpecialInstructions string     `json:"special_instructions"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type OrderItem struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID   uuid.UUID      `json:"order_id" gorm:"type:uuid;not null"`
	ProductID uuid.UUID      `json:"product_id" gorm:"type:uuid;not null"`
	VendorID  uuid.UUID      `json:"vendor_id" gorm:"type:uuid;not null"`
	Name      string         `json:"name" gorm:"not null"`
	Price     float64        `json:"price" gorm:"not null"`
	Quantity  int            `json:"quantity" gorm:"not null"`
	Total     float64        `json:"total" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type OrderTracking struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID   uuid.UUID      `json:"order_id" gorm:"type:uuid;not null"`
	Status    string         `json:"status" gorm:"not null"`
	Message   string         `json:"message"`
	Location  string         `json:"location" gorm:"type:jsonb"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	if o.OrderNumber == "" {
		o.OrderNumber = generateOrderNumber()
	}
	return nil
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	oi.Total = oi.Price * float64(oi.Quantity)
	return nil
}

func (ot *OrderTracking) BeforeCreate(tx *gorm.DB) error {
	if ot.ID == uuid.Nil {
		ot.ID = uuid.New()
	}
	return nil
}

func generateOrderNumber() string {
	return "ORD" + time.Now().Format("20060102150405") + uuid.New().String()[:8]
}
