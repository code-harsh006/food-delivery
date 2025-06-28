package cart

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/code-harsh006/food-delivery/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Module struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewModule(database *gorm.DB, redisClient *redis.Client) *Module {
	return &Module{db: database, redis: redisClient}
}

func (m *Module) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	cart := router.Group("/cart")
	cart.Use(authMiddleware)
	{
		cart.GET("", m.getCart)
		cart.POST("/add", m.addToCart)
		cart.PUT("/update", m.updateCartItem)
		cart.DELETE("/remove/:product_id", m.removeFromCart)
		cart.DELETE("/clear", m.clearCart)
	}
}

type CartResponse struct {
	Items       []CartItemResponse `json:"items"`
	TotalAmount float64            `json:"total_amount"`
	ItemCount   int                `json:"item_count"`
}

type CartItemResponse struct {
	ProductID uint       `json:"product_id"`
	Product   db.Product `json:"product"`
	Quantity  int        `json:"quantity"`
	Subtotal  float64    `json:"subtotal"`
}

func (m *Module) getCart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Try to get cart from Redis first
	cartKey := fmt.Sprintf("cart:%d", userID)
	cartData, err := m.redis.Get(context.Background(), cartKey).Result()

	var cartItems []db.CartItem

	if err == redis.Nil {
		// Cart not in Redis, get from database
		if err := m.db.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
			return
		}

		// Cache in Redis
		if len(cartItems) > 0 {
			cartJSON, _ := json.Marshal(cartItems)
			m.redis.Set(context.Background(), cartKey, cartJSON, time.Hour)
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart from cache"})
		return
	} else {
		// Parse cart from Redis
		json.Unmarshal([]byte(cartData), &cartItems)

		// Load product details
		for i := range cartItems {
			m.db.First(&cartItems[i].Product, cartItems[i].ProductID)
		}
	}

	// Build response
	var items []CartItemResponse
	var totalAmount float64
	var itemCount int

	for _, item := range cartItems {
		subtotal := item.Product.Price * float64(item.Quantity)
		items = append(items, CartItemResponse{
			ProductID: item.ProductID,
			Product:   item.Product,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
		})
		totalAmount += subtotal
		itemCount += item.Quantity
	}

	response := CartResponse{
		Items:       items,
		TotalAmount: totalAmount,
		ItemCount:   itemCount,
	}

	c.JSON(http.StatusOK, response)
}

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

func (m *Module) addToCart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if product exists and is active
	var product db.Product
	if err := m.db.Where("id = ? AND is_active = ?", req.ProductID, true).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check stock
	if product.Stock < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	// Check if item already exists in cart
	var existingItem db.CartItem
	if err := m.db.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&existingItem).Error; err == nil {
		// Update quantity
		existingItem.Quantity += req.Quantity
		if err := m.db.Save(&existingItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
			return
		}
	} else {
		// Create new cart item
		cartItem := db.CartItem{
			UserID:    userID.(uint),
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}

		if err := m.db.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
			return
		}
	}

	// Clear Redis cache
	cartKey := fmt.Sprintf("cart:%d", userID)
	m.redis.Del(context.Background(), cartKey)

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
}

type UpdateCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

func (m *Module) updateCartItem(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cartItem db.CartItem
	if err := m.db.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	cartItem.Quantity = req.Quantity
	if err := m.db.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
		return
	}

	// Clear Redis cache
	cartKey := fmt.Sprintf("cart:%d", userID)
	m.redis.Del(context.Background(), cartKey)

	c.JSON(http.StatusOK, gin.H{"message": "Cart item updated successfully"})
}

func (m *Module) removeFromCart(c *gin.Context) {
	userID, _ := c.Get("user_id")
	productID, _ := strconv.Atoi(c.Param("product_id"))

	if err := m.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&db.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	// Clear Redis cache
	cartKey := fmt.Sprintf("cart:%d", userID)
	m.redis.Del(context.Background(), cartKey)

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
}

func (m *Module) clearCart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	if err := m.db.Where("user_id = ?", userID).Delete(&db.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	// Clear Redis cache
	cartKey := fmt.Sprintf("cart:%d", userID)
	m.redis.Del(context.Background(), cartKey)

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}
