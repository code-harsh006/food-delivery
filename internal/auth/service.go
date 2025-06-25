package auth

import (
	"errors"
	"food-delivery-backend/pkg/config"
	"food-delivery-backend/pkg/db"
	"food-delivery-backend/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	cfg *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Role      string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) Register(req RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	var existingUser User
	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists with this email")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Set default role
	if req.Role == "" {
		req.Role = "customer"
	}

	// Create user
	user := User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      req.Role,
		IsActive:  true,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role, &s.cfg.JWT)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role, &s.cfg.JWT)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	refreshTokenRecord := RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.cfg.JWT.RefreshExpiry),
	}

	if err := db.DB.Create(&refreshTokenRecord).Error; err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Login(req LoginRequest) (*AuthResponse, error) {
	var user User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Update last login
	user.LastLoginAt = &time.Time{}
	*user.LastLoginAt = time.Now()
	db.DB.Save(&user)

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role, &s.cfg.JWT)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role, &s.cfg.JWT)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	refreshTokenRecord := RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.cfg.JWT.RefreshExpiry),
	}

	if err := db.DB.Create(&refreshTokenRecord).Error; err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(refreshToken string) (*AuthResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(refreshToken, &s.cfg.JWT)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	// Check if refresh token exists in database
	var tokenRecord RefreshToken
	if err := db.DB.Where("token = ? AND expires_at > ?", refreshToken, time.Now()).First(&tokenRecord).Error; err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	// Get user
	var user User
	if err := db.DB.First(&user, tokenRecord.UserID).Error; err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Generate new tokens
	newAccessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role, &s.cfg.JWT)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role, &s.cfg.JWT)
	if err != nil {
		return nil, err
	}

	// Update refresh token in database
	tokenRecord.Token = newRefreshToken
	tokenRecord.ExpiresAt = time.Now().Add(s.cfg.JWT.RefreshExpiry)
	db.DB.Save(&tokenRecord)

	return &AuthResponse{
		User:         user,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) Logout(userID uuid.UUID) error {
	// Delete all refresh tokens for the user
	return db.DB.Where("user_id = ?", userID).Delete(&RefreshToken{}).Error
}

func (s *Service) GetUserByID(userID uuid.UUID) (*User, error) {
	var user User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
