package product

import (
	"food-delivery-backend/pkg/db"

	"github.com/google/uuid"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type CreateCategoryRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	ParentID    *uuid.UUID `json:"parent_id"`
	SortOrder   int        `json:"sort_order"`
}

type CreateProductRequest struct {
	VendorID      uuid.UUID   `json:"vendor_id" binding:"required"`
	Name          string      `json:"name" binding:"required"`
	Description   string      `json:"description"`
	Images        []string    `json:"images"`
	Price         float64     `json:"price" binding:"required,gt=0"`
	DiscountPrice *float64    `json:"discount_price"`
	SKU           string      `json:"sku"`
	Stock         int         `json:"stock"`
	MinStock      int         `json:"min_stock"`
	Unit          string      `json:"unit"`
	Weight        float64     `json:"weight"`
	CategoryIDs   []uuid.UUID `json:"category_ids"`
	Tags          []string    `json:"tags"`
	IsFeatured    bool        `json:"is_featured"`
}

type UpdateProductRequest struct {
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	Images        []string    `json:"images"`
	Price         float64     `json:"price"`
	DiscountPrice *float64    `json:"discount_price"`
	Stock         int         `json:"stock"`
	MinStock      int         `json:"min_stock"`
	Unit          string      `json:"unit"`
	Weight        float64     `json:"weight"`
	CategoryIDs   []uuid.UUID `json:"category_ids"`
	Tags          []string    `json:"tags"`
	IsFeatured    bool        `json:"is_featured"`
	IsActive      bool        `json:"is_active"`
}

func (s *Service) CreateCategory(req CreateCategoryRequest) (*Category, error) {
	category := Category{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
		IsActive:    true,
	}

	if err := db.DB.Create(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *Service) GetCategories() ([]Category, error) {
	var categories []Category
	if err := db.DB.Where("is_active = ?", true).Order("sort_order ASC, name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *Service) GetCategoryByID(id uuid.UUID) (*Category, error) {
	var category Category
	if err := db.DB.Preload("Children").Preload("Products").First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *Service) CreateProduct(req CreateProductRequest) (*Product, error) {
	product := Product{
		VendorID:      req.VendorID,
		Name:          req.Name,
		Description:   req.Description,
		Images:        req.Images,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		SKU:           req.SKU,
		Stock:         req.Stock,
		MinStock:      req.MinStock,
		Unit:          req.Unit,
		Weight:        req.Weight,
		Tags:          req.Tags,
		IsFeatured:    req.IsFeatured,
		IsActive:      true,
	}

	if err := db.DB.Create(&product).Error; err != nil {
		return nil, err
	}

	// Associate with categories
	if len(req.CategoryIDs) > 0 {
		var categories []Category
		db.DB.Where("id IN ?", req.CategoryIDs).Find(&categories)
		db.DB.Model(&product).Association("Categories").Append(categories)
	}

	return &product, nil
}

func (s *Service) GetProducts(vendorID *uuid.UUID, categoryID *uuid.UUID, limit, offset int) ([]Product, error) {
	query := db.DB.Preload("Categories").Where("is_active = ?", true)

	if vendorID != nil {
		query = query.Where("vendor_id = ?", *vendorID)
	}

	if categoryID != nil {
		query = query.Joins("JOIN product_categories ON products.id = product_categories.product_id").
			Where("product_categories.category_id = ?", *categoryID)
	}

	var products []Product
	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Service) GetProductByID(id uuid.UUID) (*Product, error) {
	var product Product
	if err := db.DB.Preload("Categories").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *Service) UpdateProduct(id uuid.UUID, req UpdateProductRequest) (*Product, error) {
	var product Product
	if err := db.DB.First(&product, id).Error; err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if len(req.Images) > 0 {
		product.Images = req.Images
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.DiscountPrice != nil {
		product.DiscountPrice = req.DiscountPrice
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.MinStock >= 0 {
		product.MinStock = req.MinStock
	}
	if req.Unit != "" {
		product.Unit = req.Unit
	}
	if req.Weight > 0 {
		product.Weight = req.Weight
	}
	if len(req.Tags) > 0 {
		product.Tags = req.Tags
	}
	product.IsFeatured = req.IsFeatured
	product.IsActive = req.IsActive

	if err := db.DB.Save(&product).Error; err != nil {
		return nil, err
	}

	// Update categories
	if len(req.CategoryIDs) > 0 {
		var categories []Category
		db.DB.Where("id IN ?", req.CategoryIDs).Find(&categories)
		db.DB.Model(&product).Association("Categories").Replace(categories)
	}

	return &product, nil
}

func (s *Service) DeleteProduct(id uuid.UUID) error {
	return db.DB.Delete(&Product{}, id).Error
}

func (s *Service) SearchProducts(query string, limit, offset int) ([]Product, error) {
	var products []Product
	if err := db.DB.Preload("Categories").
		Where("is_active = ? AND (name ILIKE ? OR description ILIKE ?)", true, "%"+query+"%", "%"+query+"%").
		Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (s *Service) GetFeaturedProducts(limit int) ([]Product, error) {
	var products []Product
	if err := db.DB.Preload("Categories").
		Where("is_active = ? AND is_featured = ?", true, true).
		Limit(limit).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
