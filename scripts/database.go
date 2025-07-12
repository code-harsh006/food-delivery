package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

// InitDB initializes the MongoDB connection
func InitDB() error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb+srv://madhavinternship2024:GDuUTED803LIihgx@cluster0.zpn8u9a.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}
	MongoClient = client
	MongoDB = client.Database("food")
	log.Println("Connected to MongoDB Atlas cluster")
	return nil
}

// RunMigrations is a stub for MongoDB (no-op)
func RunMigrations() error {
	// MongoDB is schemaless; implement index creation here if needed
	return nil
}

// SeedServices is a stub for MongoDB
func SeedServices() error {
	// Implement initial data seeding for MongoDB if needed
	return nil
}
