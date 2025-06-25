package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Title        string         `json:"title" gorm:"not null"`
	AddressLine1 string         `json:"address_line1" gorm:"not null"`
	AddressLine2 string         `json:"address_line2"`
	City         string         `json:"city" gorm:"not null"`
	State        string         `json:"state" gorm:"not null"`
	PostalCode   string         `json:"postal_code" gorm:"not null"`
	Country      string         `json:"country" gorm:"default:'India'"`
	Latitude     float64        `json:"latitude"`
	Longitude    float64        `json:"longitude"`
	IsDefault    bool           `json:"is_default" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserProfile struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;uniqueIndex;not null"`
	Avatar      string         `json:"avatar"`
	DateOfBirth *time.Time     `json:"date_of_birth"`
	Gender      string         `json:"gender"`
	Preferences string         `json:"preferences" gorm:"type:jsonb"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (up *UserProfile) BeforeCreate(tx *gorm.DB) error {
	if up.ID == uuid.Nil {
		up.ID = uuid.New()
	}
	return nil
}
