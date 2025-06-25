package user

import (
	"errors"
	"food-delivery-backend/internal/auth"
	"food-delivery-backend/pkg/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type UpdateProfileRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Avatar      string `json:"avatar"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
}

type AddAddressRequest struct {
	Title        string  `json:"title" binding:"required"`
	AddressLine1 string  `json:"address_line1" binding:"required"`
	AddressLine2 string  `json:"address_line2"`
	City         string  `json:"city" binding:"required"`
	State        string  `json:"state" binding:"required"`
	PostalCode   string  `json:"postal_code" binding:"required"`
	Country      string  `json:"country"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	IsDefault    bool    `json:"is_default"`
}

func (s *Service) UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*auth.User, error) {
	var user auth.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}

	// Update user fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := db.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	// Update or create user profile
	var profile UserProfile
	err := db.DB.Where("user_id = ?", userID).First(&profile).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		profile = UserProfile{
			UserID: userID,
		}
	}

	if req.Avatar != "" {
		profile.Avatar = req.Avatar
	}
	if req.Gender != "" {
		profile.Gender = req.Gender
	}

	if err := db.DB.Save(&profile).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) GetAddresses(userID uuid.UUID) ([]Address, error) {
	var addresses []Address
	if err := db.DB.Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

func (s *Service) AddAddress(userID uuid.UUID, req AddAddressRequest) (*Address, error) {
	// If this is set as default, unset other default addresses
	if req.IsDefault {
		db.DB.Model(&Address{}).Where("user_id = ?", userID).Update("is_default", false)
	}

	address := Address{
		UserID:       userID,
		Title:        req.Title,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		Country:      req.Country,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		IsDefault:    req.IsDefault,
	}

	if address.Country == "" {
		address.Country = "India"
	}

	if err := db.DB.Create(&address).Error; err != nil {
		return nil, err
	}

	return &address, nil
}

func (s *Service) UpdateAddress(userID, addressID uuid.UUID, req AddAddressRequest) (*Address, error) {
	var address Address
	if err := db.DB.Where("id = ? AND user_id = ?", addressID, userID).First(&address).Error; err != nil {
		return nil, err
	}

	// If this is set as default, unset other default addresses
	if req.IsDefault && !address.IsDefault {
		db.DB.Model(&Address{}).Where("user_id = ? AND id != ?", userID, addressID).Update("is_default", false)
	}

	// Update fields
	address.Title = req.Title
	address.AddressLine1 = req.AddressLine1
	address.AddressLine2 = req.AddressLine2
	address.City = req.City
	address.State = req.State
	address.PostalCode = req.PostalCode
	address.Country = req.Country
	address.Latitude = req.Latitude
	address.Longitude = req.Longitude
	address.IsDefault = req.IsDefault

	if err := db.DB.Save(&address).Error; err != nil {
		return nil, err
	}

	return &address, nil
}

func (s *Service) DeleteAddress(userID, addressID uuid.UUID) error {
	return db.DB.Where("id = ? AND user_id = ?", addressID, userID).Delete(&Address{}).Error
}

func (s *Service) SetDefaultAddress(userID, addressID uuid.UUID) error {
	// First, unset all default addresses for the user
	if err := db.DB.Model(&Address{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		return err
	}

	// Then set the specified address as default
	return db.DB.Model(&Address{}).Where("id = ? AND user_id = ?", addressID, userID).Update("is_default", true).Error
}
