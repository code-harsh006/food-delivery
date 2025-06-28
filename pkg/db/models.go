package db

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Name      string         `json:"name" gorm:"not null"`
	Phone     string         `json:"phone"`
	Role      string         `json:"role" gorm:"default:'user'"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	Addresses []Address      `json:"addresses"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Address struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	UserID    uint   `json:"user_id"`
	Title     string `json:"title"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	IsDefault bool   `json:"is_default" gorm:"default:false"`
}

type Vendor struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	IsVerified  bool      `json:"is_verified" gorm:"default:false"`
	IsOnline    bool      `json:"is_online" gorm:"default:true"`
	Rating      float64   `json:"rating" gorm:"default:0"`
	Products    []Product `json:"products"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Category struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name" gorm:"unique;not null"`
	Products []Product `json:"products"`
}

type Product struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	Name        string   `json:"name" gorm:"not null"`
	Description string   `json:"description"`
	Price       float64  `json:"price" gorm:"not null"`
	Stock       int      `json:"stock" gorm:"default:0"`
	ImageURL    string   `json:"image_url"`
	VendorID    uint     `json:"vendor_id"`
	Vendor      Vendor   `json:"vendor"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `json:"category"`
	IsActive    bool     `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CartItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	UserID        uint        `json:"user_id"`
	User          User        `json:"user"`
	VendorID      uint        `json:"vendor_id"`
	Vendor        Vendor      `json:"vendor"`
	Items         []OrderItem `json:"items"`
	TotalAmount   float64     `json:"total_amount"`
	Status        string      `json:"status" gorm:"default:'pending'"`
	PaymentStatus string      `json:"payment_status" gorm:"default:'pending'"`
	DeliveryAddress string    `json:"delivery_address"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Payment struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	OrderID       uint      `json:"order_id"`
	Order         Order     `json:"order"`
	Amount        float64   `json:"amount"`
	Method        string    `json:"method"`
	Status        string    `json:"status" gorm:"default:'pending'"`
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type Delivery struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"order_id"`
	Order     Order     `json:"order"`
	AgentName string    `json:"agent_name"`
	AgentPhone string   `json:"agent_phone"`
	Status    string    `json:"status" gorm:"default:'assigned'"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

