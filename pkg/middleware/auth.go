package middleware

import (
	"food-delivery-backend/pkg/config"
	"food-delivery-backend/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header required")
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1], cfg)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid token")
			c.Abort()
			return
		}

		if claims.TokenType != "access" {
			utils.UnauthorizedResponse(c, "Invalid token type")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != "admin" {
			utils.ErrorResponse(c, 403, "Admin access required", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func VendorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || (role != "vendor" && role != "admin") {
			utils.ErrorResponse(c, 403, "Vendor access required", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
