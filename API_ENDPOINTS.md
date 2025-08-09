# Food Delivery API Endpoints Documentation

This document provides a comprehensive list of all API endpoints available in the Food Delivery application, along with their current working status.

## Status Legend

- ✅ Working - Endpoint is fully functional
- ⚠️ Partially Working - Endpoint works but with limitations
- ❌ Not Working - Endpoint is not functional

## Core API Endpoints

### Root Endpoint

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GET | [/](http://localhost:8080/) | API root with general information | ✅ Working |

### Health Check Endpoints

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GET | [/api/v1/health](http://localhost:8080/api/v1/health) | Basic health check | ✅ Working |
| GET | [/api/v1/health/detailed](http://localhost:8080/api/v1/health/detailed) | Detailed health check with system information | ✅ Working |
| GET | [/api/v1/health/ready](http://localhost:8080/api/v1/health/ready) | Readiness check | ✅ Working |
| GET | [/api/v1/health/live](http://localhost:8080/api/v1/health/live) | Liveness check | ✅ Working |
| GET | [/api/v1/status](http://localhost:8080/api/v1/status) | API status information | ✅ Working |

### Documentation Endpoint

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GET | [/api/v1/docs](http://localhost:8080/api/v1/docs) | API documentation information | ✅ Working |

## Notification Endpoints

These endpoints require authentication (JWT token).

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| POST | [/api/v1/notifications/send](http://localhost:8080/api/v1/notifications/send) | Send a notification | ✅ Working |
| GET | [/api/v1/notifications/subscribe](http://localhost:8080/api/v1/notifications/subscribe) | Subscribe to notifications | ✅ Working |

## MongoDB API Endpoints

These endpoints use MongoDB for data storage. When MongoDB is not connected, these endpoints will return a 503 Service Unavailable error with the message "Database not available: MongoDB connection is not established".

### MongoDB Health Endpoint

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GET | [/health](http://localhost:8080/health) | MongoDB health check | ✅ Working |

### MongoDB API Root Endpoint

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GET | [/api/mongo/v1](http://localhost:8080/api/mongo/v1) | MongoDB API v1 root with endpoint information | ✅ Working |
| GET | [/api/mongo/v1/test](http://localhost:8080/api/mongo/v1/test) | Simple test endpoint to verify MongoDB routes | ✅ Working |

### Authentication Endpoints

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GET | [/api/mongo/v1/auth](http://localhost:8080/api/mongo/v1/auth) | Authentication endpoints information | ✅ Working |
| POST | [/api/mongo/v1/auth/register](http://localhost:8080/api/mongo/v1/auth/register) | User registration | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| POST | [/api/mongo/v1/auth/login](http://localhost:8080/api/mongo/v1/auth/login) | User login | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| POST | [/api/mongo/v1/auth/verify-otp](http://localhost:8080/api/mongo/v1/auth/verify-otp) | Verify OTP for login/registration | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| POST | [/api/mongo/v1/auth/resend-otp](http://localhost:8080/api/mongo/v1/auth/resend-otp) | Resend OTP | ⚠️ Partially Working (Returns database error when MongoDB not connected) |

### Services Endpoints

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GET | [/api/mongo/v1/services](http://localhost:8080/api/mongo/v1/services) | Get all services | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| GET | [/api/mongo/v1/services/categories](http://localhost:8080/api/mongo/v1/services/categories) | Get service categories | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| GET | [/api/mongo/v1/services/search](http://localhost:8080/api/mongo/v1/services/search) | Search services | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| GET | [/api/mongo/v1/services/:id](http://localhost:8080/api/mongo/v1/services/:id) | Get service by ID | ⚠️ Partially Working (Returns database error when MongoDB not connected) |

### Bookings Endpoints

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GET | [/api/mongo/v1/bookings](http://localhost:8080/api/mongo/v1/bookings) | Get user bookings | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| POST | [/api/mongo/v1/bookings](http://localhost:8080/api/mongo/v1/bookings) | Create a new booking | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| GET | [/api/mongo/v1/bookings/:id](http://localhost:8080/api/mongo/v1/bookings/:id) | Get booking by ID | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| PUT | [/api/mongo/v1/bookings/:id](http://localhost:8080/api/mongo/v1/bookings/:id) | Update booking | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| DELETE | [/api/mongo/v1/bookings/:id](http://localhost:8080/api/mongo/v1/bookings/:id) | Cancel booking | ⚠️ Partially Working (Returns database error when MongoDB not connected) |

### User Endpoints

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GET | [/api/mongo/v1/users](http://localhost:8080/api/mongo/v1/users) | User management endpoints information | ✅ Working |
| GET | [/api/mongo/v1/users/profile](http://localhost:8080/api/mongo/v1/users/profile) | Get user profile | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| PUT | [/api/mongo/v1/users/profile](http://localhost:8080/api/mongo/v1/users/profile) | Update user profile | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| GET | [/api/mongo/v1/users/notifications](http://localhost:8080/api/mongo/v1/users/notifications) | Get user notifications | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| PUT | [/api/mongo/v1/users/notifications/:id/read](http://localhost:8080/api/mongo/v1/users/notifications/:id/read) | Mark notification as read | ⚠️ Partially Working (Returns database error when MongoDB not connected) |

### Admin Endpoints

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GET | [/api/mongo/v1/admin](http://localhost:8080/api/mongo/v1/admin) | Admin management endpoints information | ✅ Working |
| GET | [/api/mongo/v1/admin/bookings](http://localhost:8080/api/mongo/v1/admin/bookings) | Get all bookings | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| PUT | [/api/mongo/v1/admin/bookings/:id/status](http://localhost:8080/api/mongo/v1/admin/bookings/:id/status) | Update booking status | ⚠️ Partially Working (Returns database error when MongoDB not connected) |
| GET | [/api/mongo/v1/admin/dashboard](http://localhost:8080/api/mongo/v1/admin/dashboard) | Get dashboard statistics | ⚠️ Partially Working (Returns database error when MongoDB not connected) |

## Summary

The Food Delivery API has two main categories of endpoints:

1. **Core API Endpoints** - These are always available and do not depend on external services like MongoDB. All core endpoints are ✅ Working.

2. **MongoDB API Endpoints** - These endpoints depend on a MongoDB database connection. When MongoDB is connected, these endpoints work normally. When MongoDB is not connected (as is currently the case), these endpoints ⚠️ Partially Work by returning a clear error message indicating that the database is not available.

This design allows the application to start and serve basic functionality even when MongoDB is not available, while clearly indicating which features are unavailable due to the database connection issue.

## How to Fix MongoDB Endpoints

To make the MongoDB-dependent endpoints fully functional:

1. Ensure you have a MongoDB instance running
2. Update the `MONGODB_URI` in your `.env` file with the correct connection string
3. Make sure the MongoDB user has the proper permissions
4. Restart the application

Example MongoDB URI format:
```
mongodb+srv://username:password@cluster0.example.mongodb.net/database_name?retryWrites=true&w=majority
```