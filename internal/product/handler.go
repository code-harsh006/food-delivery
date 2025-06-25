package product

import (
	"food-delivery-backend/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler() *Handler {
	return &Handler{
		service: NewService(),
	}
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	category, err := h.service.CreateCategory(req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create category", err)
		return
	}

	utils.SuccessResponse(c, "Category created successfully", category)
}

func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get categories", err)
		return
	}

	utils.SuccessResponse(c, "Categories retrieved successfully", categories)
}

func (h *Handler) GetCategoryByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid category ID")
		return
	}

	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		utils.NotFoundResponse(c, "Category not found")
		return
	}

	utils.SuccessResponse(c, "Category retrieved successfully", category)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	product, err := h.service.CreateProduct(req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create product", err)
		return
	}

	utils.SuccessResponse(c, "Product created successfully", product)
}

func (h *Handler) GetProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var vendorID *uuid.UUID
	if vendorIDStr := c.Query("vendor_id"); vendorIDStr != "" {
		if id, err := uuid.Parse(vendorIDStr); err == nil {
			vendorID = &id
		}
	}

	var categoryID *uuid.UUID
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if id, err := uuid.Parse(categoryIDStr); err == nil {
			categoryID = &id
		}
	}

	products, err := h.service.GetProducts(vendorID, categoryID, limit, offset)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get products", err)
		return
	}

	utils.SuccessResponse(c, "Products retrieved successfully", products)
}

func (h *Handler) GetProductByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid product ID")
		return
	}

	product, err := h.service.GetProductByID(id)
	if err != nil {
		utils.NotFoundResponse(c, "Product not found")
		return
	}

	utils.SuccessResponse(c, "Product retrieved successfully", product)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid product ID")
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	product, err := h.service.UpdateProduct(id, req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update product", err)
		return
	}

	utils.SuccessResponse(c, "Product updated successfully", product)
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid product ID")
		return
	}

	err = h.service.DeleteProduct(id)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete product", err)
		return
	}

	utils.SuccessResponse(c, "Product deleted successfully", nil)
}

func (h *Handler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ValidationErrorResponse(c, "Search query is required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	products, err := h.service.SearchProducts(query, limit, offset)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to search products", err)
		return
	}

	utils.SuccessResponse(c, "Products found", products)
}

func (h *Handler) GetFeaturedProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.service.GetFeaturedProducts(limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get featured products", err)
		return
	}

	utils.SuccessResponse(c, "Featured products retrieved successfully", products)
}
