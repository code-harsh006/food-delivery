package product

import (
	"food-delivery-backend/pkg/config"
	"food-delivery-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, cfg *config.Config) {
	handler := NewHandler()
	
	// Public routes
	products := router.Group("/products")
	{
		products.GET("", handler.GetProducts)
		products.GET("/:id", handler.GetProductByID)
		products.GET("/search", handler.SearchProducts)
		products.GET("/featured", handler.GetFeaturedProducts)
	}

	categories := router.Group("/categories")
	{
		categories.GET("", handler.GetCategories)
		categories.GET("/:id", handler.GetCategoryByID)
	}

	// Admin/Vendor protected routes
	adminProducts := router.Group("/admin/products")
	adminProducts.Use(middleware.AuthMiddleware(&cfg.JWT))
	{
		adminProducts.POST("", handler.CreateProduct)
		adminProducts.PUT("/:id", handler.UpdateProduct)
		adminProducts.DELETE("/:id", handler.DeleteProduct)
	}

	adminCategories := router.Group("/admin/categories")
	adminCategories.Use(middleware.AuthMiddleware(&cfg.JWT), middleware.AdminMiddleware())
	{
		adminCategories.POST("", handler.CreateCategory)
	}
}
