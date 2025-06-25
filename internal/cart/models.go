package cart

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;uniqueIndex;not null"`
	Items     []CartItem     `json:"items" gorm:"foreignKey:CartID"`
	Total     float64        `json:"total" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CartItem struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CartID    uuid.UUID      `json:"cart_id" gorm:"type:uuid;not null"`
	ProductID uuid.UUID      `json:"product_id" gorm:"type:uuid;not null"`
	VendorID  uuid.UUID      `json:"vendor_id" gorm:"type:uuid;not null"`
	Quantity  int            `json:"quantity" gorm:"not null;default:1"`
	Price     float64        `json:"price" gorm:"not null"`
	Total     float64        `json:"total" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	if ci.ID == uuid.Nil {
		ci.ID = uuid.New()
	}
	ci.Total = ci.Price * float64(ci.Quantity)
	return nil
}

func (ci *CartItem) BeforeUpdate(tx *gorm.DB) error {
	ci.Total = ci.Price * float64(ci.Quantity)
	return nil
}
