package admin

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"food-delivery/pkg/db"
	"food-delivery/pkg/response"
	"food-delivery/internal/admin/model"
	"food-delivery/internal/admin/service"
	"food-delivery/internal/admin/validator"
	"food-delivery/internal/admin/middleware"
)
// Module initializes the admin module
func Module(r *gin.Engine, db *gorm.DB) {
	// Initialize the database connection
	dbConn := db.NewDB(db)

	// Initialize the service layer
	adminService := service.NewAdminService(dbConn)



	// Initialize the validator
	adminValidator := validator.NewAdminValidator()


	// Group routes under /admin
	adminGroup := r.Group("/admin")
	{
		// Middleware for authentication
		adminGroup.Use(middleware.Authenticate())

		// Route to get admin details
		adminGroup.GET("/details", func(c *gin.Context) {
			adminID, _ := strconv.Atoi(c.Param("id"))
			admin, err := adminService.GetAdminDetails(adminID)
			if err != nil {
				response.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
response.Success(c, admin)
		})

		// Route to update admin details
		adminGroup.PUT("/update", func(c *gin.Context) {
			var admin model.Admin
			if err := c.ShouldBindJSON(&admin); err != nil {
				response.Error(c, http.StatusBadRequest, "Invalid input")
				return
			}
			if err := adminValidator.ValidateUpdate(admin); err != nil {
				response.Error(c, http.StatusBadRequest, err.Error())
				return
			}
			if err := adminService.UpdateAdmin(admin); err != nil {
				response.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
			response.Success(c, "Admin updated successfully")
		})
	}
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

// AdminService provides methods to interact with the admin database
package service
import (
	"errors"
	"gorm.io/gorm"
	"food-delivery/internal/admin/model"
	"food-delivery/pkg/db"
)
// AdminService struct
type AdminService struct {
	db *db.DB
}

// NewAdminService creates a new instance of AdminService

func NewAdminService(db *db.DB) *AdminService {
	return &AdminService{db: db}
}

// AdminValidator provides methods to validate admin data

package validator
import (
	"errors"
	"food-delivery/internal/admin/model"
	"food-delivery/pkg/response"
)

// AdminValidator struct

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

// Admin model represents an admin in the system

package model
import (
	"gorm.io/gorm"
	"time"
	"food-delivery/pkg/db"
)
// Admin struct
// Admin represents an administrator in the system
// It contains fields for ID, Name, Email, and CreatedAt
// ID is the primary key, Name is the administrator's name, Email is their email address, and CreatedAt is the timestamp of when the admin was created
// The gorm.Model is embedded to provide default fields like ID, CreatedAt, UpdatedAt, and DeletedAt
// The table name is set to "admins" using the TableName method
type Admin struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
func (Admin) TableName() string {
	return "admins"
}
// Middleware for authentication
// This middleware checks if the user is authenticated before allowing access to admin routes
// It retrieves the user from the context and checks if they are an admin
// If not, it returns a 403 Forbidden response
// package middleware
package middleware
import (
	"net/http"
	"github.com/gin-gonic/gin"
	"food-delivery/pkg/response"
	"food-delivery/internal/admin/model"
	"food-delivery/pkg/auth"
)
// Authenticate middleware checks if the user is authenticated
// It retrieves the user from the context and checks if they are an admin
// If not, it returns a 403 Forbidden response
// Authenticate checks if the user is authenticated
// It retrieves the user from the context and checks if they are an admin
// If not, it returns a 403 Forbidden response
// package middleware
// Authenticate checks if the user is authenticated
// It retrieves the user from the context and checks if they are an admin
// If not, it returns a 403 Forbidden response
// Authenticate checks if the user is authenticated
// It retrieves the user from the context and checks if they are an admin
// If not, it returns a 403 Forbidden response
// Authenticate checks if the user is authenticated
// It retrieves the user from the context and checks if they are an admin
// If not, it returns a 403 Forbidden response
// Authenticate checks if the user is authenticated
// It retrieves the user from the context and checks if they are an admin
// If not, it returns a 403 Forbidden response
// Authenticate checks if the user is authenticated
// It retrieves the user from the context and checks if they are an admin
// If not, it returns a 403 Forbidden response

func Authenticate() gin.HandlerFunc {
l	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			response.Error(c, http.StatusForbidden, "Access denied")
			c.Abort()
			return
		}

		admin, ok := user.(*model.Admin)
		if !ok || admin.ID == 0 {
			response.Error(c, http.StatusForbidden, "Access denied")
			c.Abort()
			return
		}

		c.Next()
	}
}
// Authenticate checks if the user is authenticated

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			response.Error(c, http.StatusForbidden, "Access denied")
			c.Abort()
			return
		}

		admin, ok := user.(*model.Admin)
		if !ok || admin.ID == 0 {
			response.Error(c, http.StatusForbidden, "Access denied")
			c.Abort()
			return
		}

		c.Next()
	}
}



