package auth

import (
	"food-delivery-backend/pkg/config"
	"food-delivery-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		service: NewService(cfg),
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	response, err := h.service.Register(req)
	if err != nil {
		utils.ErrorResponse(c, 400, "Registration failed", err)
		return
	}

	utils.SuccessResponse(c, "User registered successfully", response)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	response, err := h.service.Login(req)
	if err != nil {
		utils.ErrorResponse(c, 401, "Login failed", err)
		return
	}

	utils.SuccessResponse(c, "Login successful", response)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	response, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, 401, "Token refresh failed", err)
		return
	}

	utils.SuccessResponse(c, "Token refreshed successfully", response)
}

func (h *Handler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	err := h.service.Logout(userID.(uuid.UUID))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Logout failed", err)
		return
	}

	utils.SuccessResponse(c, "Logged out successfully", nil)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	user, err := h.service.GetUserByID(userID.(uuid.UUID))
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "Profile retrieved successfully", user)
}
