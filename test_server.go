package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting simple test server...")

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Add a simple test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Test server is working"})
	})

	// Add MongoDB routes test
	mongoV1 := router.Group("/api/mongo/v1")
	{
		mongoV1.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "MongoDB API v1 root"})
		})

		auth := mongoV1.Group("/auth")
		{
			auth.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Auth endpoints available"})
			})
		}
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	fmt.Println("Server starting on port 8080...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait a bit for server to start
	time.Sleep(2 * time.Second)

	fmt.Println("Server should be running. Press Ctrl+C to stop.")

	// Keep the server running
	select {}
}
