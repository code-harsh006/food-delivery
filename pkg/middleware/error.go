package middleware

import (
	"food-delivery-backend/pkg/logger"
	"food-delivery-backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered: ", err)
				utils.InternalServerErrorResponse(c, "Internal server error", nil)
				c.Abort()
			}
		}()

		c.Next()

		// Handle errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			logger.Error("Request error: ", err.Error())
			
			switch err.Type {
			case gin.ErrorTypeBind:
				utils.ValidationErrorResponse(c, "Invalid request data")
			case gin.ErrorTypePublic:
				utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
			default:
				utils.InternalServerErrorResponse(c, "Internal server error", err.Err)
			}
		}
	}
}
