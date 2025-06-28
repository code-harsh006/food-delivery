package service

import (
	"github.com/code-harsh006/food-delivery/internal/admin/model"
	"gorm.io/gorm"
)

// AdminService provides methods to interact with the admin database
type AdminService struct {
	db *gorm.DB
}

// NewAdminService creates a new instance of AdminService
func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{db: db}
}

// GetAdminDetails retrieves the details of an admin
func (s *AdminService) GetAdminDetails(adminID int) (*model.Admin, error) {
	var admin model.Admin
	if err := s.db.First(&admin, adminID).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// UpdateAdmin updates the details of an admin
func (s *AdminService) UpdateAdmin(admin model.Admin) error {
	if err := s.db.Save(&admin).Error; err != nil {
		return err
	}
	return nil
}
