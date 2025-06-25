package auth

import (
	"food-delivery-backend/pkg/config"
	"food-delivery-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, cfg *config.Config) {
	handler := NewHandler(cfg)
	
	auth := router.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.RefreshToken)
		auth.POST("/logout", middleware.AuthMiddleware(&cfg.JWT), handler.Logout)
		auth.GET("/profile", middleware.AuthMiddleware(&cfg.JWT), handler.GetProfile)
	}
}
