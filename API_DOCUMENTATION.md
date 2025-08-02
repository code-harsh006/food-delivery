# Food Delivery API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## Response Format
All API responses follow this standard format:
```json
{
  "success": true,
  "data": { ... },
  "message": "Optional message"
}
```

## Error Responses
```json
{
  "success": false,
  "error": "Error description"
}
```

---

## 📋 Core Endpoints

### 1. Root Endpoint
**GET** `/`

**Description:** API information and available endpoints

**Request Body:** None

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Food Delivery API",
    "version": "1.0.0",
    "status": "running",
    "endpoints": {
      "health": "/api/v1/health",
      "docs": "/api/v1/docs",
      "status": "/api/v1/status"
    },
    "documentation": "Visit /api/v1/docs for API documentation"
  }
}
```

### 2. Health Check
**GET** `/api/v1/health`

**Description:** Basic health check

**Request Body:** None

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "ok",
    "message": "Food Delivery API is healthy",
    "timestamp": "2025-07-17T04:44:44.123Z",
    "service": "food-delivery-api"
  }
}
```

### 3. Detailed Health Check
**GET** `/api/v1/health/detailed`

**Description:** Detailed health check with system information

**Request Body:** None

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "ok",
    "timestamp": "2025-07-17T04:44:44.123Z",
    "service": "food-delivery-api",
    "checks": {
      "api": {
        "status": "ok",
        "message": "API is responsive",
        "latency": "0ms"
      },
      "memory": {
        "status": "ok",
        "message": "Memory usage is normal"
      }
    }
  }
}
```

### 4. Readiness Check
**GET** `/api/v1/health/ready`

**Description:** Readiness check for load balancers

**Request Body:** None

**Response:**
```json
{
  "success": true,
  "data": {
    "ready": true,
    "timestamp": "2025-07-17T04:44:44.123Z",
    "service": "food-delivery-api"
  }
}
```

### 5. Liveness Check
**GET** `/api/v1/health/live`

**Description:** Liveness check for container orchestration

**Request Body:** None

**Response:**
```json
{
  "success": true,
  "data": {
    "alive": true,
    "timestamp": "2025-07-17T04:44:44.123Z",
    "service": "food-delivery-api",
    "uptime": "running"
  }
}
```

### 6. API Status
**GET** `/api/v1/status`

**Description:** API status and version information

**Request Body:** None

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "running",
    "message": "Food Delivery API is operational",
    "timestamp": "2025-07-17T04:44:44.123Z",
    "version": "1.0.0",
    "endpoints": {
      "health": "/api/v1/health",
      "health_detailed": "/api/v1/health/detailed",
      "ready": "/api/v1/health/ready",
      "live": "/api/v1/health/live",
      "status": "/api/v1/status"
    }
  }
}
```

---

## 🔐 Authentication Endpoints

### 7. User Registration
**POST** `/api/mongo/v1/auth/register`

**Description:** Register a new user

**Authentication:** Not required

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "User registered successfully",
    "user_id": "user_123456",
    "email": "john@example.com",
    "verification_required": true
  }
}
```

### 8. User Login
**POST** `/api/mongo/v1/auth/login`

**Description:** User login

**Authentication:** Not required

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Login successful",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "user_123456",
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "+1234567890"
    }
  }
}
```

### 9. Verify OTP
**POST** `/api/mongo/v1/auth/verify-otp`

**Description:** Verify OTP for email verification

**Authentication:** Not required

**Request Body:**
```json
{
  "email": "john@example.com",
  "otp": "123456"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Email verified successfully",
    "user_id": "user_123456"
  }
}
```

### 10. Resend OTP
**POST** `/api/mongo/v1/auth/resend-otp`

**Description:** Resend OTP for email verification

**Authentication:** Not required

**Request Body:**
```json
{
  "email": "john@example.com"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "OTP sent successfully",
    "email": "john@example.com"
  }
}
```

---

## 🍕 Service Endpoints

### 11. List Services
**GET** `/api/mongo/v1/services`

**Description:** Get all available services

**Authentication:** Not required

**Request Body:** None

**Query Parameters:**
- `category` (optional): food|cleaning|maintenance
- `page` (optional): 1
- `limit` (optional): 10

**Response:**
```json
{
  "success": true,
  "data": {
    "services": [
      {
        "id": "service_123",
        "name": "Food Delivery",
        "description": "Fast food delivery service",
        "price": 15.99,
        "category": "food",
        "rating": 4.5,
        "image_url": "https://example.com/food.jpg"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10
  }
}
```

### 12. Get Service Details
**GET** `/api/mongo/v1/services/:id`

**Description:** Get service details by ID

**Authentication:** Not required

**Request Body:** None

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "service_123",
    "name": "Food Delivery",
    "description": "Fast food delivery service",
    "price": 15.99,
    "category": "food",
    "rating": 4.5,
    "image_url": "https://example.com/food.jpg",
    "provider": {
      "id": "provider_456",
      "name": "Food Express",
      "rating": 4.8
    }
  }
}
```

### 13. Get Service Categories
**GET** `/api/mongo/v1/services/categories`

**Description:** Get all service categories

**Authentication:** Not required

**Request Body:** None

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "cat_1",
      "name": "Food",
      "description": "Food delivery services",
      "icon": "🍕"
    },
    {
      "id": "cat_2",
      "name": "Cleaning",
      "description": "Cleaning services",
      "icon": "🧹"
    }
  ]
}
```

### 14. Search Services
**GET** `/api/mongo/v1/services/search`

**Description:** Search services

**Authentication:** Not required

**Request Body:** None

**Query Parameters:**
- `q` (optional): "food delivery"
- `category` (optional): food
- `min_price` (optional): 10
- `max_price` (optional): 50
- `rating` (optional): 4.0

**Response:**
```json
{
  "success": true,
  "data": {
    "services": [
      {
        "id": "service_123",
        "name": "Food Delivery",
        "description": "Fast food delivery service",
        "price": 15.99,
        "rating": 4.5
      }
    ],
    "total": 1
  }
}
```

---

## 📅 Booking Endpoints

### 15. Create Booking
**POST** `/api/mongo/v1/bookings`

**Description:** Create a new booking

**Authentication:** Required (JWT)

**Request Body:**
```json
{
  "service_id": "service_123",
  "scheduled_date": "2025-07-20T14:00:00Z",
  "address": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "zip_code": "10001"
  },
  "special_instructions": "Please deliver to front door"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "booking_id": "booking_789",
    "message": "Booking created successfully",
    "status": "pending",
    "scheduled_date": "2025-07-20T14:00:00Z"
  }
}
```

### 16. List User Bookings
**GET** `/api/mongo/v1/bookings`

**Description:** Get user's bookings

**Authentication:** Required (JWT)

**Request Body:** None

**Query Parameters:**
- `status` (optional): pending|confirmed|completed|cancelled
- `page` (optional): 1
- `limit` (optional): 10

**Response:**
```json
{
  "success": true,
  "data": {
    "bookings": [
      {
        "id": "booking_789",
        "service": {
          "id": "service_123",
          "name": "Food Delivery",
          "price": 15.99
        },
        "status": "pending",
        "scheduled_date": "2025-07-20T14:00:00Z",
        "total_amount": 15.99
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10
  }
}
```

### 17. Get Booking Details
**GET** `/api/mongo/v1/bookings/:id`

**Description:** Get booking details by ID

**Authentication:** Required (JWT)

**Request Body:** None

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "booking_789",
    "service": {
      "id": "service_123",
      "name": "Food Delivery",
      "description": "Fast food delivery service",
      "price": 15.99
    },
    "status": "pending",
    "scheduled_date": "2025-07-20T14:00:00Z",
    "address": {
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "zip_code": "10001"
    },
    "total_amount": 15.99,
    "special_instructions": "Please deliver to front door"
  }
}
```

### 18. Update Booking
**PUT** `/api/mongo/v1/bookings/:id`

**Description:** Update booking details

**Authentication:** Required (JWT)

**Request Body:**
```json
{
  "scheduled_date": "2025-07-21T15:00:00Z",
  "special_instructions": "Updated delivery instructions"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Booking updated successfully",
    "booking_id": "booking_789"
  }
}
```

### 19. Cancel Booking
**DELETE** `/api/mongo/v1/bookings/:id`

**Description:** Cancel a booking

**Authentication:** Required (JWT)

**Request Body:** None

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Booking cancelled successfully",
    "booking_id": "booking_789"
  }
}
```

---

## 👤 User Profile Endpoints

### 20. Get User Profile
**GET** `/api/mongo/v1/users/profile`

**Description:** Get user profile

**Authentication:** Required (JWT)

**Request Body:** None

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "user_123456",
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "addresses": [
      {
        "id": "addr_1",
        "title": "Home",
        "street": "123 Main St",
        "city": "New York",
        "state": "NY",
        "zip_code": "10001",
        "is_default": true
      }
    ]
  }
}
```

### 21. Update User Profile
**PUT** `/api/mongo/v1/users/profile`

**Description:** Update user profile

**Authentication:** Required (JWT)

**Request Body:**
```json
{
  "name": "John Smith",
  "phone": "+1234567890",
  "addresses": [
    {
      "title": "Home",
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "zip_code": "10001",
      "is_default": true
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Profile updated successfully",
    "user_id": "user_123456"
  }
}
```

---

## 🔔 Notification Endpoints

### 22. Send Notification
**POST** `/api/v1/notifications/send`

**Description:** Send a notification

**Authentication:** Required (JWT)

**Request Body:**
```json
{
  "type": "email|sms|push",
  "recipient": "user@example.com",
  "subject": "Notification Subject",
  "message": "Notification message content",
  "data": {
    "key": "value"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Notification sent successfully",
    "notification_id": "notif_123456",
    "status": "sent"
  }
}
```

### 23. Subscribe to Notifications
**GET** `/api/v1/notifications/subscribe`

**Description:** Subscribe to notifications (WebSocket)

**Authentication:** Required (JWT)

**Request Body:** None (WebSocket connection)

**Response:** WebSocket connection established

### 24. Get User Notifications
**GET** `/api/mongo/v1/users/notifications`

**Description:** Get user notifications

**Authentication:** Required (JWT)

**Request Body:** None

**Query Parameters:**
- `page` (optional): 1
- `limit` (optional): 10
- `read` (optional): true|false

**Response:**
```json
{
  "success": true,
  "data": {
    "notifications": [
      {
        "id": "notif_123",
        "type": "booking_confirmation",
        "title": "Booking Confirmed",
        "message": "Your booking has been confirmed",
        "read": false,
        "created_at": "2025-07-17T04:44:44.123Z"
      }
    ],
    "total": 1,
    "unread_count": 1
  }
}
```

### 25. Mark Notification as Read
**PUT** `/api/mongo/v1/users/notifications/:id/read`

**Description:** Mark notification as read

**Authentication:** Required (JWT)

**Request Body:** None

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Notification marked as read",
    "notification_id": "notif_123"
  }
}
```

---

## 👨‍💼 Admin Endpoints

### 26. Get All Bookings (Admin)
**GET** `/api/mongo/v1/admin/bookings`

**Description:** Get all bookings (Admin only)

**Authentication:** Required (JWT + Admin)

**Request Body:** None

**Query Parameters:**
- `status` (optional): pending|confirmed|completed|cancelled
- `page` (optional): 1
- `limit` (optional): 10
- `user_id` (optional): user_123456

**Response:**
```json
{
  "success": true,
  "data": {
    "bookings": [
      {
        "id": "booking_789",
        "user": {
          "id": "user_123456",
          "name": "John Doe",
          "email": "john@example.com"
        },
        "service": {
          "id": "service_123",
          "name": "Food Delivery"
        },
        "status": "pending",
        "scheduled_date": "2025-07-20T14:00:00Z",
        "total_amount": 15.99
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10
  }
}
```

### 27. Update Booking Status (Admin)
**PUT** `/api/mongo/v1/admin/bookings/:id/status`

**Description:** Update booking status (Admin only)

**Authentication:** Required (JWT + Admin)

**Request Body:**
```json
{
  "status": "confirmed|completed|cancelled",
  "notes": "Admin notes about the booking"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Booking status updated successfully",
    "booking_id": "booking_789",
    "new_status": "confirmed"
  }
}
```

### 28. Admin Dashboard
**GET** `/api/mongo/v1/admin/dashboard`

**Description:** Get admin dashboard statistics

**Authentication:** Required (JWT + Admin)

**Request Body:** None

**Response:**
```json
{
  "success": true,
  "data": {
    "total_users": 150,
    "total_bookings": 45,
    "pending_bookings": 12,
    "completed_bookings": 28,
    "cancelled_bookings": 5,
    "total_revenue": 1250.75,
    "recent_bookings": [
      {
        "id": "booking_789",
        "user_name": "John Doe",
        "service_name": "Food Delivery",
        "status": "pending",
        "amount": 15.99
      }
    ]
  }
}
```

---

## 🗄️ Database Health

### 29. MongoDB Health Check
**GET** `/health`

**Description:** MongoDB health check

**Request Body:** None

**Response:**
```json
{
  "status": "healthy",
  "database": "mongodb"
}
```

---

## 📝 Error Codes

| Status Code | Description |
|-------------|-------------|
| 400 | Bad Request - Invalid request data |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Something went wrong |

---

## 🔧 Testing the API

You can test the API using tools like:
- **Postman**
- **cURL**
- **Insomnia**
- **Thunder Client (VS Code extension)**

### Example cURL commands:

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Register user
curl -X POST http://localhost:8080/api/mongo/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "password": "securepassword123"
  }'

# Login
curl -X POST http://localhost:8080/api/mongo/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword123"
  }'

# Get services (with auth)
curl http://localhost:8080/api/mongo/v1/services \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 📚 Additional Resources

- **API Documentation:** `http://localhost:8080/api/v1/docs`
- **Health Check:** `http://localhost:8080/api/v1/health`
- **API Status:** `http://localhost:8080/api/v1/status`
- **Root Endpoint:** `http://localhost:8080/` 