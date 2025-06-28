package validator

import (
	"errors"

	"github.com/code-harsh006/food-delivery/internal/admin/model"
)

// AdminValidator provides methods to validate admin data
type AdminValidator struct{}

// NewAdminValidator creates a new instance of AdminValidator
func NewAdminValidator() *AdminValidator {
	return &AdminValidator{}
}

// ValidateUpdate checks if the admin details are valid for update
func (v *AdminValidator) ValidateUpdate(admin model.Admin) error {
	if admin.Name == "" {
		return errors.New("admin name cannot be empty")
	}
	if admin.Email == "" {
		return errors.New("admin email cannot be empty")
	}
	// Add more validation rules as needed
	return nil
}
