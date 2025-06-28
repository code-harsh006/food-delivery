package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"food-delivery/pkg/db"
)

type Module struct {
	db *gorm.DB
}

func NewModule(database *gorm.DB) *Module {
	return &Module{db: database}
}

func (m *Module) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	products := router.Group("/products")
	products.Use(authMiddleware)
	{
		products.GET("", m.getProducts)
		products.GET("/:id", m.getProduct)
		products.POST("", m.createProduct)
		products.PUT("/:id", m.updateProduct)
		products.DELETE("/:id", m.deleteProduct)
	}

	categories := router.Group("/categories")
	categories.Use(authMiddleware)
	{
		categories.GET("", m.getCategories)
		categories.POST("", m.createCategory)
	}
}

func (m *Module) getProducts(c *gin.Context) {
	var products []db.Product
	query := m.db.Preload("Vendor").Preload("Category")

	// Filter by category
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// Filter by vendor
	if vendorID := c.Query("vendor_id"); vendorID != "" {
		query = query.Where("vendor_id = ?", vendorID)
	}

	// Filter active products
	query = query.Where("is_active = ?", true)

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (m *Module) getProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var product db.Product
	if err := m.db.Preload("Vendor").Preload("Category").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
	VendorID    uint    `json:"vendor_id" binding:"required"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

func (m *Module) createProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := db.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		VendorID:    req.VendorID,
		CategoryID:  req.CategoryID,
		IsActive:    true,
	}

	if err := m.db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product": product})
}

func (m *Module) updateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product db.Product
	if err := m.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.ImageURL = req.ImageURL
	product.CategoryID = req.CategoryID

	if err := m.db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product": product})
}

func (m *Module) deleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := m.db.Delete(&db.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (m *Module) getCategories(c *gin.Context) {
	var categories []db.Category
	if err := m.db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

func (m *Module) createCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := db.Category{
		Name: req.Name,
	}

	if err := m.db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "category": category})
}

