package handler

import (
	"encoding/json"
	"food-delivery-backend/internal/auth"
	"food-delivery-backend/pkg/config"
	"food-delivery-backend/pkg/db"
	"food-delivery-backend/pkg/logger"
	"food-delivery-backend/pkg/utils"
	"net/http"
	"strings"
	"sync"
)

var (
	authService *auth.Service
	authOnce    sync.Once
)

func initAuth() {
	cfg := config.Load()
	logger.Init(cfg.Logging.Level, cfg.Logging.Format)
	
	if err := db.InitPostgres(&cfg.Database); err != nil {
		logger.Error("Failed to connect to PostgreSQL: ", err)
		return
	}
	
	if err := db.InitRedis(&cfg.Redis); err != nil {
		logger.Error("Failed to connect to Redis: ", err)
		return
	}
	
	// Auto-migrate
	models := []interface{}{
		&auth.User{},
		&auth.RefreshToken{},
	}
	
	for _, model := range models {
		db.DB.AutoMigrate(model)
	}
	
	authService = auth.NewService(cfg)
}

// AuthRegister handles user registration
func AuthRegister(w http.ResponseWriter, r *http.Request) {
	authOnce.Do(initAuth)
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	response, err := authService.Register(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.Response{
		Success: true,
		Message: "User registered successfully",
		Data:    response,
	})
}

// AuthLogin handles user login
func AuthLogin(w http.ResponseWriter, r *http.Request) {
	authOnce.Do(initAuth)
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	response, err := authService.Login(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.Response{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// AuthProfile handles getting user profile
func AuthProfile(w http.ResponseWriter, r *http.Request) {
	authOnce.Do(initAuth)
	
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Extract JWT token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}
	
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}
	
	cfg := config.Load()
	claims, err := utils.ValidateToken(tokenParts[1], &cfg.JWT)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	
	user, err := authService.GetUserByID(claims.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.Response{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    user,
	})
}
