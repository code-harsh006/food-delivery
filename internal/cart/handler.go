package cart

import (
	"food-delivery-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler() *Handler {
	return &Handler{
		service: NewService(),
	}
}

func (h *Handler) GetCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	cart, err := h.service.GetCart(userID.(uuid.UUID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get cart", err)
		return
	}

	utils.SuccessResponse(c, "Cart retrieved successfully", cart)
}

func (h *Handler) AddToCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	cart, err := h.service.AddToCart(userID.(uuid.UUID), req)
	if err != nil {
		utils.ErrorResponse(c, 400, "Failed to add item to cart", err)
		return
	}

	utils.SuccessResponse(c, "Item added to cart successfully", cart)
}

func (h *Handler) UpdateCartItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid item ID")
		return
	}

	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	cart, err := h.service.UpdateCartItem(userID.(uuid.UUID), itemID, req)
	if err != nil {
		utils.ErrorResponse(c, 400, "Failed to update cart item", err)
		return
	}

	utils.SuccessResponse(c, "Cart item updated successfully", cart)
}

func (h *Handler) RemoveFromCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid item ID")
		return
	}

	cart, err := h.service.RemoveFromCart(userID.(uuid.UUID), itemID)
	if err != nil {
		utils.ErrorResponse(c, 400, "Failed to remove item from cart", err)
		return
	}

	utils.SuccessResponse(c, "Item removed from cart successfully", cart)
}

func (h *Handler) ClearCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	err := h.service.ClearCart(userID.(uuid.UUID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to clear cart", err)
		return
	}

	utils.SuccessResponse(c, "Cart cleared successfully", nil)
}
