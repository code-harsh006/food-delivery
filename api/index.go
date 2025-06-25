package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Handler is the main serverless function entry point
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route requests
	path := r.URL.Path
	method := r.Method

	switch {
	case path == "/" || path == "/api":
		handleRoot(w, r)
	case path == "/health" || path == "/api/health":
		handleHealth(w, r)
	case strings.HasPrefix(path, "/api/v1/auth"):
		handleAuth(w, r)
	case strings.HasPrefix(path, "/api/v1/products"):
		handleProducts(w, r)
	case strings.HasPrefix(path, "/api/v1/user"):
		handleUser(w, r)
	case strings.HasPrefix(path, "/api/v1/cart"):
		handleCart(w, r)
	default:
		handleNotFound(w, r)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Food Delivery Backend API",
		"version": "1.0.0",
		"status":  "running",
		"docs":    "https://github.com/your-repo/food-delivery-backend",
		"endpoints": map[string]string{
			"health":   "/health",
			"auth":     "/api/v1/auth/*",
			"products": "/api/v1/products",
			"user":     "/api/v1/user/*",
			"cart":     "/api/v1/cart",
		},
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ok",
		"service":   "food-delivery-backend",
		"version":   "1.0.0",
		"timestamp": time.Now().Unix(),
		"platform":  "vercel-serverless",
		"go_version": os.Getenv("GO_VERSION"),
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/auth")
	
	switch {
	case path == "/register" && r.Method == "POST":
		handleAuthRegister(w, r)
	case path == "/login" && r.Method == "POST":
		handleAuthLogin(w, r)
	case path == "/profile" && r.Method == "GET":
		handleAuthProfile(w, r)
	case path == "/refresh" && r.Method == "POST":
		handleAuthRefresh(w, r)
	case path == "/logout" && r.Method == "POST":
		handleAuthLogout(w, r)
	default:
		handleNotFound(w, r)
	}
}

func handleAuthRegister(w http.ResponseWriter, r *http.Request) {
	// For now, return a mock response
	response := map[string]interface{}{
		"success": true,
		"message": "Registration endpoint - Connect your database to enable full functionality",
		"data": map[string]interface{}{
			"note": "This is a demo response. Configure DATABASE_URL environment variable to enable full auth functionality.",
		},
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	// For now, return a mock response
	response := map[string]interface{}{
		"success": true,
		"message": "Login endpoint - Connect your database to enable full functionality",
		"data": map[string]interface{}{
			"access_token": "demo-token-configure-database-for-real-auth",
			"user": map[string]interface{}{
				"id":    "demo-user-id",
				"email": "demo@example.com",
				"role":  "customer",
			},
		},
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleAuthProfile(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"success": true,
		"message": "Profile endpoint - Connect your database to enable full functionality",
		"data": map[string]interface{}{
			"id":         "demo-user-id",
			"email":      "demo@example.com",
			"first_name": "Demo",
			"last_name":  "User",
			"role":       "customer",
		},
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleAuthRefresh(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"success": true,
		"message": "Token refresh endpoint",
		"data": map[string]interface{}{
			"access_token": "new-demo-token",
		},
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleAuthLogout(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"success": true,
		"message": "Logout successful",
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleProducts(w http.ResponseWriter, r *http.Request) {
	// Mock products data
	products := []map[string]interface{}{
		{
			"id":          "1",
			"name":        "Fresh Bananas",
			"description": "Organic fresh bananas",
			"price":       2.99,
			"category":    "Fruits",
			"in_stock":    true,
		},
		{
			"id":          "2",
			"name":        "Whole Milk",
			"description": "Fresh whole milk 1L",
			"price":       3.49,
			"category":    "Dairy",
			"in_stock":    true,
		},
		{
			"id":          "3",
			"name":        "Chicken Breast",
			"description": "Fresh chicken breast 1kg",
			"price":       8.99,
			"category":    "Meat",
			"in_stock":    true,
		},
	}
	
	response := map[string]interface{}{
		"success": true,
		"message": "Products retrieved successfully",
		"data":    products,
		"note":    "This is demo data. Configure DATABASE_URL to enable full product management.",
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"success": true,
		"message": "User management endpoint",
		"data": map[string]interface{}{
			"note": "Configure DATABASE_URL environment variable to enable full user management functionality.",
		},
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleCart(w http.ResponseWriter, r *http.Request) {
	// Mock cart data
	cart := map[string]interface{}{
		"id":    "demo-cart-id",
		"items": []map[string]interface{}{},
		"total": 0.0,
		"note":  "Configure REDIS_URL to enable full cart functionality.",
	}
	
	response := map[string]interface{}{
		"success": true,
		"message": "Cart retrieved successfully",
		"data":    cart,
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"success": false,
		"message": fmt.Sprintf("Endpoint not found: %s %s", r.Method, r.URL.Path),
		"available_endpoints": []string{
			"GET /",
			"GET /health",
			"POST /api/v1/auth/register",
			"POST /api/v1/auth/login",
			"GET /api/v1/auth/profile",
			"GET /api/v1/products",
			"GET /api/v1/cart",
		},
	}
	
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(response)
}
