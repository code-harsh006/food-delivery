# Deployment Guide

This guide explains how to deploy the Food Delivery API to Render.

## Prerequisites

1. **GitHub Repository**: Your code should be in a GitHub repository
2. **Render Account**: Sign up at [render.com](https://render.com)
3. **MongoDB Atlas**: Set up your MongoDB Atlas cluster

## Method 1: Deploy with render.yaml (Recommended)

### 1. Push your code to GitHub
```bash
git add .
git commit -m "Add deployment configuration"
git push origin main
```

### 2. Connect to Render
1. Go to [render.com](https://render.com)
2. Click "New +" → "Web Service"
3. Connect your GitHub repository
4. Render will automatically detect the `render.yaml` file

### 3. Configure Environment Variables
In Render dashboard, add these environment variables:

```
PORT=8080
ENVIRONMENT=production
MONGODB_URI=mongodb+srv://madhavinternship2024:GDuUTED803LIihgx@cluster0.zpn8u9a.mongodb.net/food?retryWrites=true&w=majority
JWT_SECRET=your-secret-key-here
CORS_ORIGIN=*
```

### 4. Deploy
Click "Create Web Service" and wait for deployment.

## Method 2: Manual Deployment

### 1. Create Web Service
1. Go to Render dashboard
2. Click "New +" → "Web Service"
3. Connect your GitHub repository

### 2. Configure Build Settings
- **Environment**: Go
- **Build Command**: `go build -o food-delivery cmd/main.go`
- **Start Command**: `./food-delivery`

### 3. Add Environment Variables
Add the same environment variables as above.

### 4. Deploy
Click "Create Web Service"

## Method 3: Docker Deployment

### 1. Enable Docker
1. In Render dashboard, select your service
2. Go to "Settings" → "Build & Deploy"
3. Set **Docker** as the environment

### 2. Deploy
Render will automatically use the `Dockerfile` for deployment.

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `PORT` | Port to bind to (Render sets this) | Yes |
| `ENVIRONMENT` | Set to "production" | Yes |
| `MONGODB_URI` | MongoDB connection string | Yes |
| `JWT_SECRET` | Secret for JWT tokens | Yes |
| `CORS_ORIGIN` | CORS allowed origins | No |
| `DATABASE_URL` | PostgreSQL connection (optional) | No |

## Health Check

The application provides a health check endpoint:
- **URL**: `/health`
- **Method**: GET
- **Response**: `{"status":"healthy","database":"mongodb"}`

## API Endpoints

### PostgreSQL API (if available)
- Base URL: `https://your-app.onrender.com/api/v1`
- Health: `GET /api/v1/health`

### MongoDB API
- Base URL: `https://your-app.onrender.com/api/mongo/v1`
- Health: `GET /health`

## Troubleshooting

### 1. Port Binding Issues
- Make sure your app binds to `0.0.0.0:PORT` (not just `:PORT`)
- The `PORT` environment variable is set by Render

### 2. MongoDB Connection Issues
- Check your MongoDB Atlas IP whitelist
- Verify your connection string
- Ensure your MongoDB user has proper permissions

### 3. Build Failures
- Check the build logs in Render dashboard
- Ensure all dependencies are in `go.mod`
- Verify the build command is correct

### 4. Runtime Errors
- Check the logs in Render dashboard
- Verify all environment variables are set
- Test locally with the same environment variables

## Local Testing

Test your deployment configuration locally:

```bash
# Set environment variables
export PORT=8080
export ENVIRONMENT=production
export MONGODB_URI="your-mongodb-uri"
export JWT_SECRET="your-secret"

# Build and run
go build -o food-delivery cmd/main.go
./food-delivery
```

## Monitoring

- **Logs**: Available in Render dashboard
- **Metrics**: Basic metrics in Render dashboard
- **Health Checks**: Automatic health checks by Render

## Security Notes

1. **Environment Variables**: Never commit secrets to Git
2. **CORS**: Configure `CORS_ORIGIN` properly for production
3. **MongoDB**: Use proper authentication and IP restrictions
4. **HTTPS**: Render provides automatic HTTPS

## Support

If you encounter issues:
1. Check Render logs
2. Verify environment variables
3. Test locally with the same configuration
4. Check MongoDB Atlas connection 

---

## Step-by-Step Fixes

### 1. PostgreSQL: Connection Refused

- **If you do NOT need PostgreSQL:**  
  Remove or comment out all PostgreSQL-related code and configuration. Make sure your app does not try to connect to PostgreSQL if you only want to use MongoDB.

- **If you DO need PostgreSQL:**  
  - Make sure PostgreSQL is running locally or update your `DATABASE_URL` to point to a running instance (local or cloud).
  - Check your credentials and that the user/database exist.
  - If running locally, start PostgreSQL and ensure it listens on the correct port.

### 2. MongoDB: TLS Internal Error

This is a common issue with MongoDB Atlas and Go drivers. Here’s how to fix it:

- **Check your MongoDB URI:**  
  Make sure it uses the correct format and credentials.  
  Example:  
  ```
  mongodb+srv://<username>:<password>@cluster0.zpn8u9a.mongodb.net/food?retryWrites=true&w=majority
  ```

- **IP Whitelist:**  
  In MongoDB Atlas, make sure your current IP (or 0.0.0.0/0 for testing) is whitelisted.

- **TLS/SSL Support:**  
  - The Go MongoDB driver requires valid certificates. If you are behind a proxy or have custom CA, you may need to set `tlsInsecure=true` for testing (not for production).
  - Try adding `&tlsInsecure=true` to your connection string for local testing:
    ```
    mongodb+srv://<username>:<password>@cluster0.zpn8u9a.mongodb.net/food?retryWrites=true&w=majority&tlsInsecure=true
    ```
  - If this works, the issue is with certificate validation (e.g., missing CA certs on your machine).

- **Firewall/Network:**  
  - Ensure your network allows outbound connections to MongoDB Atlas (ports 27017+).
  - Some corporate or university networks block these ports.

---

## Code Improvements

- **Graceful Fallback:**  
  Your code already starts with "limited functionality" if a DB is not connected. Consider logging more details and providing a health endpoint that reports which DBs are up.

- **Environment Variable Checks:**  
  Before starting, check if all required environment variables are set. If not, log a clear error and exit.

- **Configuration File:**  
  Use a config file or `.env` for local development, and document all required variables.

---

## Next Steps

1. **Decide if you need PostgreSQL.**  
   - If not, remove its usage from your codebase and configs.

2. **Fix MongoDB Atlas connection:**  
   - Try `&tlsInsecure=true` in your URI for local testing.
   - Check your IP whitelist in Atlas.
   - Make sure your credentials are correct.

3. **Test Locally:**  
   - Use the same environment variables as in production.
   - Run your app and check logs for improved error messages.

---

Would you like me to:
- Update your code to skip PostgreSQL if not configured?
- Add better error messages and health reporting?
- Show you exactly where to add `tlsInsecure=true` in your MongoDB connection code?
- Refactor your config loading for better local/production separation?

Let me know which improvements you want, and I’ll make the code changes for you! 

---

## 1. Production-Ready Go Server Example

```go
// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get PORT from environment, fallback to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := "0.0.0.0:" + port

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Example root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Food Delivery API is running!")
	})

	log.Printf("🚀 Server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
```

---

## 2. Render.com Requirements Explained

- **PORT Environment Variable:**  
  Render automatically sets the `PORT` environment variable for your service. Your app **must** listen on this port, or Render will not detect an open port and will fail deployment.

- **Bind to 0.0.0.0:**  
  Binding to `0.0.0.0` ensures your app is accessible from outside the container (i.e., by Render’s load balancer). Binding to `localhost` or `127.0.0.1` will make your app unreachable.

- **render.yaml Example (optional):**
  ```yaml
  services:
    - type: web
      name: food-delivery-api
      env: go
      buildCommand: go build -o food-delivery cmd/main.go
      startCommand: ./food-delivery
      envVars:
        - key: ENVIRONMENT
          value: production
        - key: MONGODB_URI
          sync: false
  ```

---

## 3. Production Dockerfile (Multi-Stage)

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o food-delivery ./cmd/main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/food-delivery .
EXPOSE 8080
CMD ["./food-delivery"]
```
- **Note:** Render will still inject the `PORT` variable. Your Go code must use it (see above).

---

## 4. Troubleshooting Checklist

- **No open ports detected:**  
  - Ensure your Go app uses `os.Getenv("PORT")` and binds to `0.0.0.0:<PORT>`.
  - Check logs in Render dashboard for startup errors.

- **Check if port is open:**  
  - Log the address your server is listening on.
  - Use `/health` endpoint to verify service is running.

- **Verify environment variables:**  
  - In Render dashboard, check "Environment" tab for all required variables.
  - Use `log.Printf` to print out critical env vars at startup (avoid printing secrets).

- **Debugging connection timeouts:**  
  - Ensure your database/network dependencies are accessible from Render.
  - Check security groups, firewalls, and connection strings.
  - Use Render’s shell (if available) to test connectivity.

---

## Copy-Paste Summary

- Use the Go server code above.
- Ensure you use the `PORT` env variable and bind to `0.0.0.0`.
- Use the Dockerfile if deploying with Docker.
- Check Render logs for errors and confirm `/health` is reachable.

Let me know if you want a more advanced example (e.g., with MongoDB integration or graceful shutdown)! 