package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// Health is a simple health check handler
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "food-delivery-backend",
		"version":   "1.0.0",
	}
	
	json.NewEncoder(w).Encode(response)
}
