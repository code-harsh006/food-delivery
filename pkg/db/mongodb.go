package db

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

var MongoClient *mongo.Client
var MongoDB *mongo.Database

// InitMongoDB initializes the MongoDB connection
func InitMongoDB(uri string, isProduction bool) error {
	if uri == "" {
		uri = "mongodb+srv://madhavjadav638:GDuUTED803LIihgx@cluster0.jd56d.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0&tlsInsecure=true"
	}

	// Log the URI being used (without password for security)
	logURI := strings.Replace(uri, "GDuUTED803LIihgx", "***", -1)
	log.Printf("🔗 Attempting to connect to MongoDB: %s", logURI)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Configure client options with TLS settings
	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetServerSelectionTimeout(15 * time.Second)
	clientOptions.SetSocketTimeout(15 * time.Second)
	clientOptions.SetConnectTimeout(15 * time.Second)

	// Configure TLS to be more permissive
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	clientOptions.SetTLSConfig(tlsConfig)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		if isProduction {
			log.Fatal("❌ MongoDB connection failed:", err)
		} else {
			log.Printf("❌ Failed to create MongoDB client: %v", err)
			return fmt.Errorf("failed to create MongoDB client: %w", err)
		}
	}

	// Ping to verify connection with retry and better error handling
	var pingErr error
	for i := 0; i < 3; i++ {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
		pingErr = client.Ping(pingCtx, nil)
		pingCancel()

		if pingErr == nil {
			break
		}

		log.Printf("⚠️  MongoDB ping attempt %d failed: %v", i+1, pingErr)

		// Provide specific error guidance
		if strings.Contains(pingErr.Error(), "authentication failed") {
			log.Printf("🔐 Authentication failed. Please check:")
			log.Printf("   - Username: madhavjadav638")
			log.Printf("   - Password: GDuUTED803LIihgx")
			log.Printf("   - User exists in MongoDB Atlas Database Access")
			log.Printf("   - User has 'Read and write to any database' privileges")
			log.Printf("   - Wait 1-2 minutes after creating/resetting user")
		} else if strings.Contains(pingErr.Error(), "tls") {
			log.Printf("🔒 TLS error. This might be a network/firewall issue.")
		} else if strings.Contains(pingErr.Error(), "timeout") {
			log.Printf("⏰ Connection timeout. Check your internet connection.")
		}

		if i < 2 {
			log.Printf("🔄 Retrying in 3 seconds...")
			time.Sleep(3 * time.Second)
		}
	}

	if pingErr != nil {
		if isProduction {
			log.Fatal("❌ Failed to ping MongoDB after 3 attempts:", pingErr)
		} else {
			log.Printf("❌ Failed to ping MongoDB after 3 attempts: %v", pingErr)
			log.Println("⚠️  MongoDB connection failed, but continuing without it")
			log.Println("💡 To fix this:")
			log.Println("   1. Go to MongoDB Atlas → Database Access")
			log.Println("   2. Create/edit user 'madhavjadav638'")
			log.Println("   3. Set password to 'GDuUTED803LIihgx'")
			log.Println("   4. Set privileges to 'Read and write to any database'")
			log.Println("   5. Wait 1-2 minutes after changes")
			return nil
		}
	}

	MongoClient = client

	// Try to connect to 'food' database, if it doesn't exist it will be created
	MongoDB = client.Database("food")
	log.Println("✅ Connected to MongoDB Atlas cluster successfully!")

	// Test database access
	collections, err := MongoDB.ListCollectionNames(ctx, nil)
	if err != nil {
		log.Printf("⚠️  Warning: Could not list collections: %v", err)
		log.Printf("💡 This might mean the user doesn't have access to the 'food' database")
		log.Printf("🔧 Trying to create a test collection to verify access...")

		// Try to create a test collection to verify write access
		testCollection := MongoDB.Collection("test_connection")
		_, insertErr := testCollection.InsertOne(ctx, map[string]interface{}{
			"test":      true,
			"timestamp": time.Now(),
		})
		if insertErr != nil {
			log.Printf("❌ Write access test failed: %v", insertErr)
		} else {
			log.Printf("✅ Write access test successful!")
			// Clean up test document
			testCollection.DeleteOne(ctx, map[string]interface{}{"test": true})
		}
	} else {
		log.Printf("📚 Available collections: %v", collections)
	}

	return nil
}

// CloseMongoDB closes the MongoDB connection
func CloseMongoDB() error {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return MongoClient.Disconnect(ctx)
	}
	return nil
}

// GetMongoDB returns the MongoDB database instance
func GetMongoDB() *mongo.Database {
	return MongoDB
}

// GetMongoClient returns the MongoDB client instance
func GetMongoClient() *mongo.Client {
	return MongoClient
}
