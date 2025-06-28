package admin

import (
	"net/http"
	"strconv"

	"github.com/code-harsh006/food-delivery/internal/admin/model"
	"github.com/code-harsh006/food-delivery/internal/admin/service"
	"github.com/code-harsh006/food-delivery/internal/admin/validator"
	"github.com/code-harsh006/food-delivery/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module initializes the admin module
func Module(r *gin.Engine, db *gorm.DB) {
	// Initialize the database connection
	dbConn := db

	// Initialize the service layer
	adminService := service.NewAdminService(dbConn)

	// Initialize the validator
	adminValidator := validator.NewAdminValidator()

	// Group routes under /admin
	adminGroup := r.Group("/admin")
	{
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
