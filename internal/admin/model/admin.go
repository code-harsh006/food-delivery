package model

import (
	"time"

	"gorm.io/gorm"
)

// Admin represents an administrator in the system
type Admin struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Admin) TableName() string {
	return "admins"
}
