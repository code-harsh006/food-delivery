package user

import (
	"food-delivery-backend/pkg/config"
	"food-delivery-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, cfg *config.Config) {
	handler := NewHandler()
	
	user := router.Group("/user")
	user.Use(middleware.AuthMiddleware(&cfg.JWT))
	{
		user.PUT("/profile", handler.UpdateProfile)
		
		// Address management
		user.GET("/addresses", handler.GetAddresses)
		user.POST("/addresses", handler.AddAddress)
		user.PUT("/addresses/:id", handler.UpdateAddress)
		user.DELETE("/addresses/:id", handler.DeleteAddress)
		user.PUT("/addresses/:id/default", handler.SetDefaultAddress)
	}
}
