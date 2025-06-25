package cart

import (
	"food-delivery-backend/pkg/config"
	"food-delivery-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, cfg *config.Config) {
	handler := NewHandler()
	
	cart := router.Group("/cart")
	cart.Use(middleware.AuthMiddleware(&cfg.JWT))
	{
		cart.GET("", handler.GetCart)
		cart.POST("/items", handler.AddToCart)
		cart.PUT("/items/:id", handler.UpdateCartItem)
		cart.DELETE("/items/:id", handler.RemoveFromCart)
		cart.DELETE("", handler.ClearCart)
	}
}
