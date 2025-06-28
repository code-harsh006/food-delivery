package search

import (
	"net/http"
	"strings"

	"github.com/code-harsh006/food-delivery/pkg/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	db *gorm.DB
}

func NewModule(database *gorm.DB) *Module {
	return &Module{db: database}
}

func (m *Module) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	search := router.Group("/search")
	search.Use(authMiddleware)
	{
		search.GET("/products", m.searchProducts)
		search.GET("/vendors", m.searchVendors)
		search.GET("/recommendations", m.getRecommendations)
	}
}

func (m *Module) searchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var products []db.Product
	searchQuery := m.db.Preload("Vendor").Preload("Category").Where("is_active = ?", true)

	// Search in product name and description
	searchTerm := "%" + strings.ToLower(query) + "%"
	searchQuery = searchQuery.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)

	// Filter by category if provided
	if categoryID := c.Query("category_id"); categoryID != "" {
		searchQuery = searchQuery.Where("category_id = ?", categoryID)
	}

	// Filter by price range
	if minPrice := c.Query("min_price"); minPrice != "" {
		searchQuery = searchQuery.Where("price >= ?", minPrice)
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		searchQuery = searchQuery.Where("price <= ?", maxPrice)
	}

	// Sort by relevance (name match first, then description match)
	searchQuery = searchQuery.Order("CASE WHEN LOWER(name) LIKE '" + searchTerm + "' THEN 1 ELSE 2 END, name ASC")

	if err := searchQuery.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": products,
		"count":   len(products),
	})
}

func (m *Module) searchVendors(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var vendors []db.Vendor
	searchQuery := m.db.Where("is_verified = ? AND is_online = ?", true, true)

	// Search in vendor name and address
	searchTerm := "%" + strings.ToLower(query) + "%"
	searchQuery = searchQuery.Where("LOWER(name) LIKE ? OR LOWER(address) LIKE ?", searchTerm, searchTerm)

	// Sort by rating and name
	searchQuery = searchQuery.Order("rating DESC, name ASC")

	if err := searchQuery.Find(&vendors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search vendors"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": vendors,
		"count":   len(vendors),
	})
}

func (m *Module) getRecommendations(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Stub implementation for personalized recommendations
	// In a real implementation, you would use ML algorithms, user history, etc.

	var recommendations []db.Product

	// Get user's order history to find preferred categories
	var orderItems []db.OrderItem
	m.db.Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.user_id = ?", userID).
		Preload("Product.Category").
		Find(&orderItems)

	// Extract preferred categories
	categoryMap := make(map[uint]bool)
	for _, item := range orderItems {
		categoryMap[item.Product.CategoryID] = true
	}

	var preferredCategories []uint
	for categoryID := range categoryMap {
		preferredCategories = append(preferredCategories, categoryID)
	}

	// Get recommended products from preferred categories
	query := m.db.Preload("Vendor").Preload("Category").Where("is_active = ?", true)

	if len(preferredCategories) > 0 {
		query = query.Where("category_id IN ?", preferredCategories)
	}

	// Also include popular products (high stock, good vendor rating)
	query = query.Joins("JOIN vendors ON vendors.id = products.vendor_id").
		Where("vendors.is_verified = ? AND vendors.is_online = ?", true, true).
		Order("vendors.rating DESC, products.stock DESC").
		Limit(20)

	if err := query.Find(&recommendations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendations"})
		return
	}

	// If no personalized recommendations, get popular products
	if len(recommendations) == 0 {
		m.db.Preload("Vendor").Preload("Category").
			Where("is_active = ?", true).
			Joins("JOIN vendors ON vendors.id = products.vendor_id").
			Where("vendors.is_verified = ? AND vendors.is_online = ?", true, true).
			Order("vendors.rating DESC").
			Limit(10).
			Find(&recommendations)
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"count":           len(recommendations),
		"type":            "personalized",
	})
}
