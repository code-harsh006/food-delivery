package product

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	ParentID    *uuid.UUID     `json:"parent_id" gorm:"type:uuid"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Products []Product  `json:"products,omitempty" gorm:"many2many:product_categories;"`
}

type Product struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	VendorID    uuid.UUID      `json:"vendor_id" gorm:"type:uuid;not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Images      []string       `json:"images" gorm:"type:jsonb"`
	Price       float64        `json:"price" gorm:"not null"`
	DiscountPrice *float64     `json:"discount_price"`
	SKU         string         `json:"sku" gorm:"uniqueIndex"`
	Stock       int            `json:"stock" gorm:"default:0"`
	MinStock    int            `json:"min_stock" gorm:"default:0"`
	Unit        string         `json:"unit" gorm:"default:'piece'"`
	Weight      float64        `json:"weight"`
	Dimensions  string         `json:"dimensions" gorm:"type:jsonb"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	IsFeatured  bool           `json:"is_featured" gorm:"default:false"`
	Tags        []string       `json:"tags" gorm:"type:jsonb"`
	Attributes  string         `json:"attributes" gorm:"type:jsonb"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	Categories []Category `json:"categories,omitempty" gorm:"many2many:product_categories;"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
