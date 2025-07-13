package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

// InitMongoDB initializes the MongoDB connection
func InitMongoDB(uri string) error {
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
	log.Println("✅ Connected to MongoDB Atlas cluster")
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
