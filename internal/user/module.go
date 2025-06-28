package user

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
	users := router.Group("/users")
	users.Use(authMiddleware)
	{
		users.GET("/profile", m.getProfile)
		users.PUT("/profile", m.updateProfile)
		users.GET("/addresses", m.getAddresses)
		users.POST("/addresses", m.createAddress)
		users.PUT("/addresses/:id", m.updateAddress)
		users.DELETE("/addresses/:id", m.deleteAddress)
	}
}

func (m *Module) getProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var user db.User
	if err := m.db.Preload("Addresses").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func (m *Module) updateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user db.User
	if err := m.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Name = req.Name
	user.Phone = req.Phone

	if err := m.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "user": user})
}

func (m *Module) getAddresses(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var addresses []db.Address
	if err := m.db.Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch addresses"})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

type CreateAddressRequest struct {
	Title     string `json:"title" binding:"required"`
	Street    string `json:"street" binding:"required"`
	City      string `json:"city" binding:"required"`
	State     string `json:"state" binding:"required"`
	ZipCode   string `json:"zip_code" binding:"required"`
	IsDefault bool   `json:"is_default"`
}

func (m *Module) createAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address := db.Address{
		UserID:    userID.(uint),
		Title:     req.Title,
		Street:    req.Street,
		City:      req.City,
		State:     req.State,
		ZipCode:   req.ZipCode,
		IsDefault: req.IsDefault,
	}

	if err := m.db.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Address created successfully", "address": address})
}

func (m *Module) updateAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	addressID, _ := strconv.Atoi(c.Param("id"))
	
	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var address db.Address
	if err := m.db.Where("id = ? AND user_id = ?", addressID, userID).First(&address).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	address.Title = req.Title
	address.Street = req.Street
	address.City = req.City
	address.State = req.State
	address.ZipCode = req.ZipCode
	address.IsDefault = req.IsDefault

	if err := m.db.Save(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address updated successfully", "address": address})
}

func (m *Module) deleteAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	addressID, _ := strconv.Atoi(c.Param("id"))
	
	if err := m.db.Where("id = ? AND user_id = ?", addressID, userID).Delete(&db.Address{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}

