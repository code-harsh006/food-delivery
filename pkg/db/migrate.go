package db

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Address{},
		&Vendor{},
		&Category{},
		&Product{},
		&CartItem{},
		&Order{},
		&OrderItem{},
		&Payment{},
		&Delivery{},
	)
}

func SeedData(db *gorm.DB) error {
	// Seed categories
	categories := []Category{
		{Name: "Fruits"},
		{Name: "Vegetables"},
		{Name: "Snacks"},
		{Name: "Beverages"},
		{Name: "Fast Food"},
		{Name: "Desserts"},
	}

	for _, category := range categories {
		db.FirstOrCreate(&category, Category{Name: category.Name})
	}

	// Seed vendors
	vendors := []Vendor{
		{Name: "Fresh Mart", Email: "freshmart@example.com", Phone: "1234567890", Address: "123 Market St", IsVerified: true},
		{Name: "Pizza Palace", Email: "pizza@example.com", Phone: "0987654321", Address: "456 Food Ave", IsVerified: true},
		{Name: "Healthy Bites", Email: "healthy@example.com", Phone: "1122334455", Address: "789 Green St", IsVerified: true},
	}

	for _, vendor := range vendors {
		db.FirstOrCreate(&vendor, Vendor{Email: vendor.Email})
	}

	// Seed products
	products := []Product{
		{Name: "Apple", Description: "Fresh red apples", Price: 2.99, Stock: 100, VendorID: 1, CategoryID: 1},
		{Name: "Banana", Description: "Ripe bananas", Price: 1.99, Stock: 150, VendorID: 1, CategoryID: 1},
		{Name: "Margherita Pizza", Description: "Classic pizza with tomato and mozzarella", Price: 12.99, Stock: 50, VendorID: 2, CategoryID: 5},
		{Name: "Pepperoni Pizza", Description: "Pizza with pepperoni", Price: 15.99, Stock: 30, VendorID: 2, CategoryID: 5},
		{Name: "Green Salad", Description: "Fresh mixed greens", Price: 8.99, Stock: 25, VendorID: 3, CategoryID: 2},
	}

	for _, product := range products {
		db.FirstOrCreate(&product, Product{Name: product.Name, VendorID: product.VendorID})
	}

	return nil
}

