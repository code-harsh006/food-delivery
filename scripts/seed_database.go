package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User model for food delivery
type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Phone     string             `bson:"phone"`
	Address   string             `bson:"address"`
	City      string             `bson:"city"`
	State     string             `bson:"state"`
	ZipCode   string             `bson:"zip_code"`
	IsActive  bool               `bson:"is_active"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// Service model for food delivery services
type Service struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Category    string             `bson:"category"`
	BasePrice   float64            `bson:"base_price"`
	Duration    int                `bson:"duration"` // in minutes
	IsActive    bool               `bson:"is_active"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

// Vendor model for service providers
type Vendor struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Category    string             `bson:"category"`
	Location    string             `bson:"location"`
	City        string             `bson:"city"`
	State       string             `bson:"state"`
	Phone       string             `bson:"phone"`
	Email       string             `bson:"email"`
	Rating      float64            `bson:"rating"`
	IsActive    bool               `bson:"is_active"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

// Booking model for service bookings
type Booking struct {
	ID              primitive.ObjectID `bson:"_id"`
	UserID          primitive.ObjectID `bson:"user_id"`
	ServiceID       primitive.ObjectID `bson:"service_id"`
	VendorID        primitive.ObjectID `bson:"vendor_id"`
	ScheduledDate   time.Time          `bson:"scheduled_date"`
	ScheduledTime   string             `bson:"scheduled_time"`
	Status          string             `bson:"status"`
	TotalAmount     float64            `bson:"total_amount"`
	PaymentStatus   string             `bson:"payment_status"`
	SpecialRequests string             `bson:"special_requests"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
}

// BookingStatus model for tracking booking status changes
type BookingStatus struct {
	ID        primitive.ObjectID `bson:"_id"`
	BookingID primitive.ObjectID `bson:"booking_id"`
	Status    string             `bson:"status"`
	Message   string             `bson:"message"`
	UpdatedBy string             `bson:"updated_by"`
	CreatedAt time.Time          `bson:"created_at"`
}

// Notification model for user notifications
type Notification struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Title     string             `bson:"title"`
	Message   string             `bson:"message"`
	Type      string             `bson:"type"`
	IsRead    bool               `bson:"is_read"`
	CreatedAt time.Time          `bson:"created_at"`
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// MongoDB connection
	client, err := mongo.Connect(context.TODO(),
		options.Client().ApplyURI("mongodb+srv://madhavjadav638:GDuUTED803LIihgx@cluster0.jd56d.mongodb.net/food?retryWrites=true&w=majority&appName=Cluster0&tls=true&tlsAllowInvalidCertificates=true"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	db := client.Database("food")

	// Collections
	usersColl := db.Collection("users")
	servicesColl := db.Collection("services")
	vendorsColl := db.Collection("vendors")
	bookingsColl := db.Collection("bookings")
	bookingStatusesColl := db.Collection("booking_statuses")
	notificationsColl := db.Collection("notifications")

	fmt.Println("🚀 Starting database seeding...")

	// Clear existing data
	fmt.Println("🧹 Clearing existing data...")
	usersColl.DeleteMany(ctx, bson.M{})
	servicesColl.DeleteMany(ctx, bson.M{})
	vendorsColl.DeleteMany(ctx, bson.M{})
	bookingsColl.DeleteMany(ctx, bson.M{})
	bookingStatusesColl.DeleteMany(ctx, bson.M{})
	notificationsColl.DeleteMany(ctx, bson.M{})

	// Generate Users
	fmt.Println("👥 Generating users...")
	users := []User{}
	for i := 0; i < 50; i++ {
		user := User{
			ID:        primitive.NewObjectID(),
			Name:      faker.Name(),
			Email:     faker.Email(),
			Password:  hashPassword("password123"),
			Phone:     faker.Phonenumber(),
			Address:   faker.GetAddress().Street(),
			City:      faker.GetAddress().City(),
			State:     faker.GetAddress().State(),
			ZipCode:   faker.GetAddress().ZipCode(),
			IsActive:  true,
			CreatedAt: time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour),
			UpdatedAt: time.Now(),
		}
		users = append(users, user)
	}

	if _, err := usersColl.InsertMany(ctx, toDocs(users)); err != nil {
		log.Fatal("Failed to insert users:", err)
	}
	fmt.Printf("✅ Created %d users\n", len(users))

	// Generate Services
	fmt.Println("🛠️  Generating services...")
	serviceCategories := []string{"Food Delivery", "Restaurant Booking", "Catering", "Food Preparation", "Kitchen Cleaning", "Menu Planning", "Food Safety Training"}
	services := []Service{}

	for i := 0; i < 30; i++ {
		service := Service{
			ID:          primitive.NewObjectID(),
			Name:        faker.Word() + " " + faker.Word(),
			Description: faker.Sentence(),
			Category:    serviceCategories[rand.Intn(len(serviceCategories))],
			BasePrice:   float64(rand.Intn(200) + 50),
			Duration:    rand.Intn(120) + 30, // 30-150 minutes
			IsActive:    true,
			CreatedAt:   time.Now().Add(-time.Duration(rand.Intn(180)) * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		}
		services = append(services, service)
	}

	if _, err := servicesColl.InsertMany(ctx, toDocs(services)); err != nil {
		log.Fatal("Failed to insert services:", err)
	}
	fmt.Printf("✅ Created %d services\n", len(services))

	// Generate Vendors
	fmt.Println("🏪 Generating vendors...")
	vendorCategories := []string{"Restaurant", "Catering Service", "Food Truck", "Cafe", "Bakery", "Food Consultant", "Kitchen Equipment"}
	vendors := []Vendor{}

	for i := 0; i < 25; i++ {
		vendor := Vendor{
			ID:          primitive.NewObjectID(),
			Name:        faker.Company(),
			Description: faker.Sentence(),
			Category:    vendorCategories[rand.Intn(len(vendorCategories))],
			Location:    faker.GetAddress().Street(),
			City:        faker.GetAddress().City(),
			State:       faker.GetAddress().State(),
			Phone:       faker.Phonenumber(),
			Email:       faker.Email(),
			Rating:      float64(rand.Intn(20)+80) / 20.0, // 4.0-5.0
			IsActive:    true,
			CreatedAt:   time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		}
		vendors = append(vendors, vendor)
	}

	if _, err := vendorsColl.InsertMany(ctx, toDocs(vendors)); err != nil {
		log.Fatal("Failed to insert vendors:", err)
	}
	fmt.Printf("✅ Created %d vendors\n", len(vendors))

	// Generate Bookings
	fmt.Println("📅 Generating bookings...")
	bookingStatuses := []string{"pending", "confirmed", "completed", "cancelled"}
	paymentStatuses := []string{"pending", "paid", "failed"}
	timeSlots := []string{"09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00", "16:00", "17:00", "18:00"}
	bookings := []Booking{}
	bookingStatusesList := []BookingStatus{}

	for i := 0; i < 100; i++ {
		user := users[rand.Intn(len(users))]
		service := services[rand.Intn(len(services))]
		vendor := vendors[rand.Intn(len(vendors))]

		// Random date within next 30 days
		scheduledDate := time.Now().AddDate(0, 0, rand.Intn(30))

		booking := Booking{
			ID:              primitive.NewObjectID(),
			UserID:          user.ID,
			ServiceID:       service.ID,
			VendorID:        vendor.ID,
			ScheduledDate:   scheduledDate,
			ScheduledTime:   timeSlots[rand.Intn(len(timeSlots))],
			Status:          bookingStatuses[rand.Intn(len(bookingStatuses))],
			TotalAmount:     service.BasePrice,
			PaymentStatus:   paymentStatuses[rand.Intn(len(paymentStatuses))],
			SpecialRequests: faker.Sentence(),
			CreatedAt:       time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
			UpdatedAt:       time.Now(),
		}
		bookings = append(bookings, booking)

		// Create booking status
		status := BookingStatus{
			ID:        primitive.NewObjectID(),
			BookingID: booking.ID,
			Status:    booking.Status,
			Message:   "Booking " + booking.Status,
			UpdatedBy: "system",
			CreatedAt: booking.CreatedAt,
		}
		bookingStatusesList = append(bookingStatusesList, status)
	}

	if _, err := bookingsColl.InsertMany(ctx, toDocs(bookings)); err != nil {
		log.Fatal("Failed to insert bookings:", err)
	}
	if _, err := bookingStatusesColl.InsertMany(ctx, toDocs(bookingStatusesList)); err != nil {
		log.Fatal("Failed to insert booking statuses:", err)
	}
	fmt.Printf("✅ Created %d bookings\n", len(bookings))

	// Generate Notifications
	fmt.Println("🔔 Generating notifications...")
	notificationTypes := []string{"booking_confirmed", "booking_reminder", "payment_received", "service_completed"}
	notifications := []Notification{}

	for i := 0; i < 150; i++ {
		user := users[rand.Intn(len(users))]
		notification := Notification{
			ID:        primitive.NewObjectID(),
			UserID:    user.ID,
			Title:     faker.Word() + " " + faker.Word(),
			Message:   faker.Sentence(),
			Type:      notificationTypes[rand.Intn(len(notificationTypes))],
			IsRead:    rand.Float32() > 0.3, // 70% read
			CreatedAt: time.Now().Add(-time.Duration(rand.Intn(7)) * 24 * time.Hour),
		}
		notifications = append(notifications, notification)
	}

	if _, err := notificationsColl.InsertMany(ctx, toDocs(notifications)); err != nil {
		log.Fatal("Failed to insert notifications:", err)
	}
	fmt.Printf("✅ Created %d notifications\n", len(notifications))

	fmt.Println("🎉 Database seeded successfully!")
	fmt.Println("📊 Summary:")
	fmt.Printf("   👥 Users: %d\n", len(users))
	fmt.Printf("   🛠️  Services: %d\n", len(services))
	fmt.Printf("   🏪 Vendors: %d\n", len(vendors))
	fmt.Printf("   📅 Bookings: %d\n", len(bookings))
	fmt.Printf("   🔔 Notifications: %d\n", len(notifications))
	fmt.Println("🚀 Your food delivery API is ready with realistic data!")
}

func toDocs[T any](items []T) []interface{} {
	docs := make([]interface{}, len(items))
	for i, v := range items {
		docs[i] = v
	}
	return docs
}
