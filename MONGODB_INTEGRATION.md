# MongoDB Integration

This project now includes MongoDB integration alongside the existing PostgreSQL database. The MongoDB integration provides additional functionality for user management, service bookings, and notifications.

## Features Added

### 1. User Management
- User registration with email verification
- OTP-based authentication
- User profile management
- Email verification system

### 2. Service Management
- Service catalog with categories
- Service search functionality
- Service details and pricing

### 3. Booking System
- Create service bookings
- View booking history
- Update booking status
- Cancel bookings
- Booking status tracking

### 4. Notification System
- User notifications
- Booking confirmations
- Status updates

### 5. Admin Dashboard
- View all bookings
- Update booking status
- Dashboard statistics

## Database Configuration

### MongoDB Connection
The MongoDB connection is configured using the `MONGODB_URI` environment variable:

```env
MONGODB_URI=mongodb+srv://madhavinternship2024:GDuUTED803LIihgx@cluster0.zpn8u9a.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
```

### Database Collections
The following collections are used in MongoDB:

- `users` - User profiles and authentication
- `otps` - One-time passwords for verification
- `services` - Available services
- `bookings` - Service bookings
- `booking_statuses` - Booking status history
- `notifications` - User notifications

## API Endpoints

### Authentication
- `POST /api/mongo/v1/auth/register` - User registration
- `POST /api/mongo/v1/auth/login` - User login
- `POST /api/mongo/v1/auth/verify-otp` - OTP verification
- `POST /api/mongo/v1/auth/resend-otp` - Resend OTP

### Services
- `GET /api/mongo/v1/services` - Get all services
- `GET /api/mongo/v1/services/:id` - Get service by ID
- `GET /api/mongo/v1/services/categories` - Get service categories
- `GET /api/mongo/v1/services/search` - Search services

### Bookings
- `POST /api/mongo/v1/bookings` - Create booking
- `GET /api/mongo/v1/bookings` - Get user bookings
- `GET /api/mongo/v1/bookings/:id` - Get booking by ID
- `PUT /api/mongo/v1/bookings/:id` - Update booking
- `DELETE /api/mongo/v1/bookings/:id` - Cancel booking

### User Profile
- `GET /api/mongo/v1/users/profile` - Get user profile
- `PUT /api/mongo/v1/users/profile` - Update user profile
- `GET /api/mongo/v1/users/notifications` - Get user notifications
- `PUT /api/mongo/v1/users/notifications/:id/read` - Mark notification as read

### Admin (Service Provider)
- `GET /api/mongo/v1/admin/bookings` - Get all bookings
- `PUT /api/mongo/v1/admin/bookings/:id/status` - Update booking status
- `GET /api/mongo/v1/admin/dashboard` - Get dashboard stats

## Usage Examples

### 1. User Registration
```bash
curl -X POST http://localhost:8080/api/mongo/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "phone": "+1234567890",
    "name": "John Doe",
    "accommodation_type": "apartment",
    "address": "123 Main St"
  }'
```

### 2. OTP Verification
```bash
curl -X POST http://localhost:8080/api/mongo/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "code": "123456"
  }'
```

### 3. Create Booking
```bash
curl -X POST http://localhost:8080/api/mongo/v1/bookings \
  -H "Content-Type: application/json" \
  -H "User-ID: 507f1f77bcf86cd799439011" \
  -d '{
    "service_id": "507f1f77bcf86cd799439012",
    "scheduled_date": "2024-01-15",
    "scheduled_time": "14:00",
    "special_requests": "Please arrive 10 minutes early"
  }'
```

### 4. Get Services
```bash
curl http://localhost:8080/api/mongo/v1/services
```

### 5. Search Services
```bash
curl "http://localhost:8080/api/mongo/v1/services/search?q=cleaning"
```

## Authentication

Currently, the system uses a simple header-based authentication for demonstration purposes. In production, you should implement proper JWT-based authentication.

To use protected endpoints, include the `User-ID` header:
```bash
curl -H "User-ID: 507f1f77bcf86cd799439011" \
  http://localhost:8080/api/mongo/v1/users/profile
```

## Environment Variables

Make sure to set the following environment variables:

```env
# MongoDB Configuration
MONGODB_URI=mongodb+srv://madhavinternship2024:GDuUTED803LIihgx@cluster0.zpn8u9a.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0

# Server Configuration
PORT=8080
ENVIRONMENT=development

# JWT Configuration
JWT_SECRET=your-secret-key
```

## Running the Application

1. Copy the environment file:
```bash
cp config.env.example .env
```

2. Update the `.env` file with your configuration

3. Run the application:
```bash
go run cmd/main.go
```

4. The server will start and show the status of both PostgreSQL and MongoDB connections

## Health Check

You can check the health of the application:
```bash
curl http://localhost:8080/api/v1/health
```

This will return:
```json
{
  "status": "healthy",
  "database": "mongodb"
}
```

## Notes

- The MongoDB integration runs alongside the existing PostgreSQL database
- Both databases can be used simultaneously
- The system gracefully handles cases where either database is unavailable
- OTP codes are printed to the console for development purposes
- In production, implement proper email/SMS services for OTP delivery 