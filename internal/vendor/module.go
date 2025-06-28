package vendor

import (
	"net/http"
	"strconv"

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
	vendors := router.Group("/vendors")
	vendors.Use(authMiddleware)
	{
		vendors.GET("", m.getVendors)
		vendors.GET("/:id", m.getVendor)
		vendors.POST("", m.createVendor)
		vendors.PUT("/:id", m.updateVendor)
		vendors.PUT("/:id/toggle-online", m.toggleOnlineStatus)
		vendors.GET("/:id/products", m.getVendorProducts)
	}
}

func (m *Module) getVendors(c *gin.Context) {
	var vendors []db.Vendor
	query := m.db.Preload("Products")

	// Filter by verification status
	if verified := c.Query("verified"); verified != "" {
		query = query.Where("is_verified = ?", verified == "true")
	}

	// Filter by online status
	if online := c.Query("online"); online != "" {
		query = query.Where("is_online = ?", online == "true")
	}

	if err := query.Find(&vendors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vendors"})
		return
	}

	c.JSON(http.StatusOK, vendors)
}

func (m *Module) getVendor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var vendor db.Vendor
	if err := m.db.Preload("Products").First(&vendor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
		return
	}

	c.JSON(http.StatusOK, vendor)
}

type CreateVendorRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func (m *Module) createVendor(c *gin.Context) {
	var req CreateVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vendor := db.Vendor{
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Address:    req.Address,
		IsVerified: false,
		IsOnline:   true,
	}

	if err := m.db.Create(&vendor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vendor"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Vendor created successfully", "vendor": vendor})
}

func (m *Module) updateVendor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req CreateVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var vendor db.Vendor
	if err := m.db.First(&vendor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
		return
	}

	vendor.Name = req.Name
	vendor.Email = req.Email
	vendor.Phone = req.Phone
	vendor.Address = req.Address

	if err := m.db.Save(&vendor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vendor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vendor updated successfully", "vendor": vendor})
}

func (m *Module) toggleOnlineStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var vendor db.Vendor
	if err := m.db.First(&vendor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
		return
	}

	vendor.IsOnline = !vendor.IsOnline

	if err := m.db.Save(&vendor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vendor status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vendor status updated", "vendor": vendor})
}

func (m *Module) getVendorProducts(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var products []db.Product
	if err := m.db.Where("vendor_id = ? AND is_active = ?", id, true).Preload("Category").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vendor products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
