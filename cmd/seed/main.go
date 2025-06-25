package main

import (
	"food-delivery-backend/internal/auth"
	"food-delivery-backend/internal/product"
	"food-delivery-backend/pkg/config"
	"food-delivery-backend/pkg/db"
	"food-delivery-backend/pkg/logger"
	"food-delivery-backend/pkg/utils"

	"github.com/google/uuid"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.Logging.Level, cfg.Logging.Format)

	// Initialize database
	if err := db.InitPostgres(&cfg.Database); err != nil {
		logger.Fatal("Failed to connect to PostgreSQL: ", err)
	}

	logger.Info("Starting database seeding...")

	// Seed users
	if err := seedUsers(); err != nil {
		logger.Error("Failed to seed users: ", err)
	}

	// Seed categories
	if err := seedCategories(); err != nil {
		logger.Error("Failed to seed categories: ", err)
	}

	// Seed products
	if err := seedProducts(); err != nil {
		logger.Error("Failed to seed products: ", err)
	}

	logger.Info("Database seeding completed successfully")
}

func seedUsers() error {
	users := []auth.User{
		{
			Email:     "admin@fooddelivery.com",
			FirstName: "Admin",
			LastName:  "User",
			Phone:     "+1234567890",
			Role:      "admin",
			IsActive:  true,
		},
		{
			Email:     "vendor@fooddelivery.com",
			FirstName: "Vendor",
			LastName:  "User",
			Phone:     "+1234567891",
			Role:      "vendor",
			IsActive:  true,
		},
		{
			Email:     "customer@fooddelivery.com",
			FirstName: "Customer",
			LastName:  "User",
			Phone:     "+1234567892",
			Role:      "customer",
			IsActive:  true,
		},
	}

	for i := range users {
		hashedPassword, err := utils.HashPassword("password123")
		if err != nil {
			return err
		}
		users[i].Password = hashedPassword

		// Check if user already exists
		var existingUser auth.User
		if err := db.DB.Where("email = ?", users[i].Email).First(&existingUser).Error; err != nil {
			if err := db.DB.Create(&users[i]).Error; err != nil {
				return err
			}
			logger.Info("Created user: ", users[i].Email)
		}
	}

	return nil
}

func seedCategories() error {
	categories := []product.Category{
		{
			Name:        "Fruits & Vegetables",
			Description: "Fresh fruits and vegetables",
			IsActive:    true,
			SortOrder:   1,
		},
		{
			Name:        "Dairy & Eggs",
			Description: "Milk, cheese, eggs and dairy products",
			IsActive:    true,
			SortOrder:   2,
		},
		{
			Name:        "Meat & Seafood",
			Description: "Fresh meat and seafood",
			IsActive:    true,
			SortOrder:   3,
		},
		{
			Name:        "Beverages",
			Description: "Drinks and beverages",
			IsActive:    true,
			SortOrder:   4,
		},
		{
			Name:        "Snacks",
			Description: "Snacks and quick bites",
			IsActive:    true,
			SortOrder:   5,
		},
	}

	for _, category := range categories {
		var existingCategory product.Category
		if err := db.DB.Where("name = ?", category.Name).First(&existingCategory).Error; err != nil {
			if err := db.DB.Create(&category).Error; err != nil {
				return err
			}
			logger.Info("Created category: ", category.Name)
		}
	}

	return nil
}

func seedProducts() error {
	// Get a vendor user
	var vendor auth.User
	if err := db.DB.Where("role = ?", "vendor").First(&vendor).Error; err != nil {
		logger.Warn("No vendor found, skipping product seeding")
		return nil
	}

	// Get categories
	var categories []product.Category
	db.DB.Find(&categories)

	if len(categories) == 0 {
		logger.Warn("No categories found, skipping product seeding")
		return nil
	}

	products := []product.Product{
		{
			VendorID:    vendor.ID,
			Name:        "Fresh Bananas",
			Description: "Fresh organic bananas",
			Price:       2.99,
			Stock:       100,
			Unit:        "kg",
			Weight:      1.0,
			IsActive:    true,
			IsFeatured:  true,
			Tags:        []string{"organic", "fresh", "fruit"},
		},
		{
			VendorID:    vendor.ID,
			Name:        "Whole Milk",
			Description: "Fresh whole milk",
			Price:       3.49,
			Stock:       50,
			Unit:        "liter",
			Weight:      1.0,
			IsActive:    true,
			Tags:        []string{"dairy", "fresh"},
		},
		{
			VendorID:    vendor.ID,
			Name:        "Chicken Breast",
			Description: "Fresh chicken breast",
			Price:       8.99,
			Stock:       25,
			Unit:        "kg",
			Weight:      1.0,
			IsActive:    true,
			Tags:        []string{"meat", "protein", "fresh"},
		},
		{
			VendorID:    vendor.ID,
			Name:        "Orange Juice",
			Description: "Fresh orange juice",
			Price:       4.99,
			Stock:       30,
			Unit:        "liter",
			Weight:      1.0,
			IsActive:    true,
			Tags:        []string{"beverage", "juice", "fresh"},
		},
		{
			VendorID:    vendor.ID,
			Name:        "Potato Chips",
			Description: "Crispy potato chips",
			Price:       1.99,
			Stock:       75,
			Unit:        "pack",
			Weight:      0.2,
			IsActive:    true,
			Tags:        []string{"snack", "crispy"},
		},
	}

	for i, prod := range products {
		var existingProduct product.Product
		if err := db.DB.Where("name = ? AND vendor_id = ?", prod.Name, prod.VendorID).First(&existingProduct).Error; err != nil {
			if err := db.DB.Create(&products[i]).Error; err != nil {
				return err
			}

			// Associate with first category (simplified)
			if len(categories) > 0 {
				db.DB.Model(&products[i]).Association("Categories").Append(&categories[0])
			}

			logger.Info("Created product: ", prod.Name)
		}
	}

	return nil
}
