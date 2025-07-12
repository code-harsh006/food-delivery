package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetServices returns all active services
func GetServices(c *gin.Context) {
	var services []Service

	// Get query parameters for filtering
	category := c.Query("category")

	collection := MongoDB.Collection("services")
	filter := bson.M{"is_active": true}
	if category != "" {
		filter["category"] = category
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &services); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode services"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
		"total":    len(services),
	})
}

// GetServiceByID returns a specific service by ID
func GetServiceByID(c *gin.Context) {
	idParam := c.Param("id")
	serviceID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	var service Service
	collection := MongoDB.Collection("services")
	err = collection.FindOne(context.Background(), bson.M{"_id": serviceID, "is_active": true}).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"service": service})
}

// GetServiceCategories returns all service categories
func GetServiceCategories(c *gin.Context) {
	collection := MongoDB.Collection("services")

	// Use aggregation to get distinct categories
	pipeline := []bson.M{
		{"$match": bson.M{"is_active": true}},
		{"$group": bson.M{"_id": "$category"}},
		{"$project": bson.M{"category": "$_id", "_id": 0}},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode categories"})
		return
	}

	var categories []string
	for _, result := range results {
		if category, ok := result["category"].(string); ok {
			categories = append(categories, category)
		}
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// SearchServices searches services by name or description
func SearchServices(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var services []Service
	collection := MongoDB.Collection("services")

	// Create a text search filter
	filter := bson.M{
		"is_active": true,
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search services"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &services); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode services"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
		"total":    len(services),
		"query":    query,
	})
}
