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