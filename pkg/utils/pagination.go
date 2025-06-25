package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationParams struct {
	Page     int `json:"page"`
	Limit    int `json:"limit"`
	Offset   int `json:"offset"`
	SortBy   string `json:"sort_by"`
	SortDesc bool `json:"sort_desc"`
}

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
}

func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDesc := c.DefaultQuery("sort_desc", "true") == "true"

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	return PaginationParams{
		Page:     page,
		Limit:    limit,
		Offset:   offset,
		SortBy:   sortBy,
		SortDesc: sortDesc,
	}
}

func CreatePaginationResponse(data interface{}, params PaginationParams, total int64) PaginationResponse {
	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))
	
	return PaginationResponse{
		Data:       data,
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}
}
