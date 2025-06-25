package cart

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"food-delivery-backend/internal/product"
	"food-delivery-backend/pkg/db"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type AddToCartRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

type CartResponse struct {
	ID     uuid.UUID          `json:"id"`
	UserID uuid.UUID          `json:"user_id"`
	Items  []CartItemResponse `json:"items"`
	Total  float64            `json:"total"`
}

type CartItemResponse struct {
	ID        uuid.UUID       `json:"id"`
	ProductID uuid.UUID       `json:"product_id"`
	Product   product.Product `json:"product"`
	VendorID  uuid.UUID       `json:"vendor_id"`
	Quantity  int             `json:"quantity"`
	Price     float64         `json:"price"`
	Total     float64         `json:"total"`
}

func (s *Service) getCartKey(userID uuid.UUID) string {
	return fmt.Sprintf("cart:%s", userID.String())
}

func (s *Service) GetCart(userID uuid.UUID) (*CartResponse, error) {
	ctx := context.Background()
	cartKey := s.getCartKey(userID)

	// Try to get cart from Redis first
	cartData, err := db.RedisClient.Get(ctx, cartKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	var cart Cart
	if errors.Is(err, redis.Nil) {
		// Cart not in Redis, try database
		if err := db.DB.Preload("Items").Where("user_id = ?", userID).First(&cart).Error; err != nil {
			// Create new cart if not found
			cart = Cart{
				UserID: userID,
				Total:  0,
			}
			if err := db.DB.Create(&cart).Error; err != nil {
				return nil, err
			}
		}
	} else {
		// Parse cart from Redis
		if err := json.Unmarshal([]byte(cartData), &cart); err != nil {
			return nil, err
		}
	}

	// Get product details for cart items
	var cartItems []CartItemResponse
	total := 0.0

	for _, item := range cart.Items {
		var prod product.Product
		if err := db.DB.First(&prod, item.ProductID).Error; err != nil {
			continue // Skip if product not found
		}

		cartItem := CartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Product:   prod,
			VendorID:  item.VendorID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Total:     item.Total,
		}
		cartItems = append(cartItems, cartItem)
		total += item.Total
	}

	cartResponse := &CartResponse{
		ID:     cart.ID,
		UserID: cart.UserID,
		Items:  cartItems,
		Total:  total,
	}

	// Update cache
	s.updateCartCache(userID, cart)

	return cartResponse, nil
}

func (s *Service) AddToCart(userID uuid.UUID, req AddToCartRequest) (*CartResponse, error) {
	// Get product details
	var prod product.Product
	if err := db.DB.First(&prod, req.ProductID).Error; err != nil {
		return nil, errors.New("product not found")
	}

	if !prod.IsActive {
		return nil, errors.New("product is not available")
	}

	if prod.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// Get or create cart
	var cart Cart
	if err := db.DB.Preload("Items").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		cart = Cart{
			UserID: userID,
			Total:  0,
		}
		if err := db.DB.Create(&cart).Error; err != nil {
			return nil, err
		}
	}

	// Check if item already exists in cart
	var existingItem *CartItem
	for i := range cart.Items {
		if cart.Items[i].ProductID == req.ProductID {
			existingItem = &cart.Items[i]
			break
		}
	}

	price := prod.Price
	if prod.DiscountPrice != nil && *prod.DiscountPrice > 0 {
		price = *prod.DiscountPrice
	}

	if existingItem != nil {
		// Update existing item
		existingItem.Quantity += req.Quantity
		existingItem.Price = price
		existingItem.Total = price * float64(existingItem.Quantity)
		
		if err := db.DB.Save(existingItem).Error; err != nil {
			return nil, err
		}
	} else {
		// Add new item
		newItem := CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			VendorID:  prod.VendorID,
			Quantity:  req.Quantity,
			Price:     price,
			Total:     price * float64(req.Quantity),
		}
		
		if err := db.DB.Create(&newItem).Error; err != nil {
			return nil, err
		}
		cart.Items = append(cart.Items, newItem)
	}

	// Update cart total
	s.updateCartTotal(&cart)
	
	// Update cache
	s.updateCartCache(userID, cart)

	return s.GetCart(userID)
}

func (s *Service) UpdateCartItem(userID, itemID uuid.UUID, req UpdateCartItemRequest) (*CartResponse, error) {
	var cartItem CartItem
	if err := db.DB.Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		First(&cartItem).Error; err != nil {
		return nil, errors.New("cart item not found")
	}

	if req.Quantity == 0 {
		// Remove item from cart
		if err := db.DB.Delete(&cartItem).Error; err != nil {
			return nil, err
		}
	} else {
		// Update quantity
		cartItem.Quantity = req.Quantity
		cartItem.Total = cartItem.Price * float64(req.Quantity)
		
		if err := db.DB.Save(&cartItem).Error; err != nil {
			return nil, err
		}
	}

	// Clear cache to force refresh
	s.clearCartCache(userID)

	return s.GetCart(userID)
}

func (s *Service) RemoveFromCart(userID, itemID uuid.UUID) (*CartResponse, error) {
	if err := db.DB.Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		Delete(&CartItem{}).Error; err != nil {
		return nil, err
	}

	// Clear cache to force refresh
	s.clearCartCache(userID)

	return s.GetCart(userID)
}

func (s *Service) ClearCart(userID uuid.UUID) error {
	// Get cart
	var cart Cart
	if err := db.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return err
	}

	// Delete all cart items
	if err := db.DB.Where("cart_id = ?", cart.ID).Delete(&CartItem{}).Error; err != nil {
		return err
	}

	// Update cart total
	cart.Total = 0
	db.DB.Save(&cart)

	// Clear cache
	s.clearCartCache(userID)

	return nil
}

func (s *Service) updateCartTotal(cart *Cart) {
	var total float64
	for _, item := range cart.Items {
		total += item.Total
	}
	cart.Total = total
	db.DB.Save(cart)
}

func (s *Service) updateCartCache(userID uuid.UUID, cart Cart) {
	ctx := context.Background()
	cartKey := s.getCartKey(userID)
	
	cartData, _ := json.Marshal(cart)
	db.RedisClient.Set(ctx, cartKey, cartData, 24*time.Hour)
}

func (s *Service) clearCartCache(userID uuid.UUID) {
	ctx := context.Background()
	cartKey := s.getCartKey(userID)
	db.RedisClient.Del(ctx, cartKey)
}
