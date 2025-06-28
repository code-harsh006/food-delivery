package order

import (
	"net/http"

	"github.com/code-harsh006/food-delivery/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Module struct{}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	notifications := router.Group("/notifications")
	notifications.Use(authMiddleware)
	{
		notifications.POST("/send", m.sendNotification)
		notifications.GET("/subscribe", m.subscribeToNotifications)
	}
}

type SendNotificationRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
	Type    string `json:"type"`
}

func (m *Module) sendNotification(c *gin.Context) {
	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Stub implementation for sending push notifications
	// In a real implementation, you would integrate with FCM, APNs, etc.

	logger.Info("Sending notification",
		zap.Uint("user_id", req.UserID),
		zap.String("title", req.Title),
		zap.String("message", req.Message),
		zap.String("type", req.Type),
	)

	// Simulate successful notification send
	c.JSON(http.StatusOK, gin.H{
		"message":         "Notification sent successfully",
		"notification_id": "notif_123456",
		"status":          "sent",
	})
}

func (m *Module) subscribeToNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// WebSocket connection would be established here
	// For now, we'll just return a success message

	logger.Info("User subscribed to notifications", zap.Any("user_id", userID))

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscribed to notifications successfully",
		"user_id": userID,
	})
}

// Stub methods for different notification types
func (m *Module) SendOrderStatusNotification(userID uint, orderID uint, status string) {
	logger.Info("Order status notification",
		zap.Uint("user_id", userID),
		zap.Uint("order_id", orderID),
		zap.String("status", status),
	)
}

func (m *Module) SendDeliveryNotification(userID uint, message string) {
	logger.Info("Delivery notification",
		zap.Uint("user_id", userID),
		zap.String("message", message),
	)
}

func (m *Module) SendPromotionalNotification(userID uint, title, message string) {
	logger.Info("Promotional notification",
		zap.Uint("user_id", userID),
		zap.String("title", title),
		zap.String("message", message),
	)
}
