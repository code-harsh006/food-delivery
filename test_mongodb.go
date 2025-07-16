package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Test MongoDB connection
	uri := "mongodb+srv://madhavjadav638:GDuUTED803LIihgx@cluster0.jd56d.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0&tlsInsecure=true"

	fmt.Println("🔗 Testing MongoDB connection...")

	// Log the URI being used (without password for security)
	logURI := strings.Replace(uri, "GDuUTED803LIihgx", "***", -1)
	fmt.Printf("📡 Connecting to: %s\n", logURI)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetServerSelectionTimeout(15 * time.Second)
	clientOptions.SetSocketTimeout(15 * time.Second)
	clientOptions.SetConnectTimeout(15 * time.Second)

	// Configure TLS to be more permissive
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	clientOptions.SetTLSConfig(tlsConfig)

	fmt.Println("🔐 Attempting to connect...")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("❌ Failed to create MongoDB client: %v", err)
	}
	defer client.Disconnect(ctx)

	// Test ping with retry
	fmt.Println("🏓 Testing connection with ping...")
	var pingErr error
	for i := 0; i < 3; i++ {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
		pingErr = client.Ping(pingCtx, nil)
		pingCancel()

		if pingErr == nil {
			break
		}

		fmt.Printf("⚠️  Ping attempt %d failed: %v\n", i+1, pingErr)

		// Provide specific error guidance
		if strings.Contains(pingErr.Error(), "authentication failed") {
			fmt.Println("🔐 Authentication failed. Please check:")
			fmt.Println("   - Username: madhavjadav638")
			fmt.Println("   - Password: GDuUTED803LIihgx")
			fmt.Println("   - User exists in MongoDB Atlas Database Access")
			fmt.Println("   - User has 'Read and write to any database' privileges")
			fmt.Println("   - Wait 1-2 minutes after creating/resetting user")
		} else if strings.Contains(pingErr.Error(), "tls") {
			fmt.Println("🔒 TLS error. This might be a network/firewall issue.")
		} else if strings.Contains(pingErr.Error(), "timeout") {
			fmt.Println("⏰ Connection timeout. Check your internet connection.")
		}

		if i < 2 {
			fmt.Println("🔄 Retrying in 3 seconds...")
			time.Sleep(3 * time.Second)
		}
	}

	if pingErr != nil {
		log.Fatalf("❌ Failed to ping MongoDB after 3 attempts: %v", pingErr)
	}

	fmt.Println("✅ MongoDB connection successful!")

	// Test database access
	fmt.Println("📚 Testing database access...")
	db := client.Database("food")
	collections, err := db.ListCollectionNames(ctx, nil)
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not list collections: %v\n", err)
		fmt.Println("💡 This might mean the user doesn't have access to the 'food' database")
		fmt.Println("🔧 Trying to create a test collection to verify access...")

		// Try to create a test collection to verify write access
		testCollection := db.Collection("test_connection")
		_, insertErr := testCollection.InsertOne(ctx, map[string]interface{}{
			"test":      true,
			"timestamp": time.Now(),
		})
		if insertErr != nil {
			fmt.Printf("❌ Write access test failed: %v\n", insertErr)
		} else {
			fmt.Printf("✅ Write access test successful!\n")
			// Clean up test document
			testCollection.DeleteOne(ctx, map[string]interface{}{"test": true})
		}
	} else {
		fmt.Printf("📚 Collections in database: %v\n", collections)
	}

	fmt.Println("🎉 MongoDB connection test completed successfully!")
}
