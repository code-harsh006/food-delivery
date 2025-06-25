package user

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

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := h.service.UpdateProfile(userID.(uuid.UUID), req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update profile", err)
		return
	}

	utils.SuccessResponse(c, "Profile updated successfully", user)
}

func (h *Handler) GetAddresses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	addresses, err := h.service.GetAddresses(userID.(uuid.UUID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get addresses", err)
		return
	}

	utils.SuccessResponse(c, "Addresses retrieved successfully", addresses)
}

func (h *Handler) AddAddress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req AddAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	address, err := h.service.AddAddress(userID.(uuid.UUID), req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to add address", err)
		return
	}

	utils.SuccessResponse(c, "Address added successfully", address)
}

func (h *Handler) UpdateAddress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid address ID")
		return
	}

	var req AddAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	address, err := h.service.UpdateAddress(userID.(uuid.UUID), addressID, req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update address", err)
		return
	}

	utils.SuccessResponse(c, "Address updated successfully", address)
}

func (h *Handler) DeleteAddress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid address ID")
		return
	}

	err = h.service.DeleteAddress(userID.(uuid.UUID), addressID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete address", err)
		return
	}

	utils.SuccessResponse(c, "Address deleted successfully", nil)
}

func (h *Handler) SetDefaultAddress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid address ID")
		return
	}

	err = h.service.SetDefaultAddress(userID.(uuid.UUID), addressID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to set default address", err)
		return
	}

	utils.SuccessResponse(c, "Default address set successfully", nil)
}
