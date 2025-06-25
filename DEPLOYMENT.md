# Deployment Guide

## Common Deployment Issues and Solutions

### 1. "Bun could not find package.json" Error

**Problem**: Deployment platform is trying to use Node.js/Bun instead of Go.

**Solutions**:

#### For Vercel:
- Use `vercel.json` configuration (already included)
- Ensure you're using `@vercel/go` builder

#### For Railway:
- Use `railway.json` or `Dockerfile`
- Set build command: `go build -o main cmd/main.go`

#### For Render:
- Use `render.yaml` configuration
- Set environment to `go`

#### For Heroku:
- Use `heroku.yml` configuration
- Ensure Dockerfile is properly configured

### 2. Platform-Specific Deployment

#### Vercel
\`\`\`bash
# Install Vercel CLI
npm i -g vercel

# Deploy
vercel --prod
\`\`\`

#### Railway
\`\`\`bash
# Install Railway CLI
npm i -g @railway/cli

# Login and deploy
railway login
railway up
\`\`\`

#### Render
\`\`\`bash
# Connect your GitHub repo to Render
# Use the render.yaml configuration
\`\`\`

#### Heroku
\`\`\`bash
# Install Heroku CLI
# Set stack to container
heroku stack:set container -a your-app-name
git push heroku main
\`\`\`

#### Docker Deployment
\`\`\`bash
# Production build
make docker-prod

# Or manually
docker-compose -f docker-compose.production.yml up -d
\`\`\`

### 3. Environment Variables

Ensure these are set in your deployment platform:

\`\`\`env
# Required
DB_HOST=your-db-host
DB_USER=your-db-user
DB_PASSWORD=your-db-password
DB_NAME=your-db-name
REDIS_HOST=your-redis-host
JWT_SECRET=your-jwt-secret

# Optional
PORT=8080
GIN_MODE=release
\`\`\`

### 4. Database Setup

#### PostgreSQL
- Create database and user
- Run migrations: `make migrate`
- Seed data: `make seed`

#### Redis
- Ensure Redis is accessible
- Test connection with health check

### 5. Health Check

Test your deployment:
\`\`\`bash
curl https://your-app-url/health
\`\`\`

Expected response:
\`\`\`json
{
  "status": "ok",
  "timestamp": 1234567890,
  "service": "food-delivery-backend"
}
\`\`\`

## Platform-Specific Instructions

### Vercel
1. Connect GitHub repository
2. Set build command: `go build -o api cmd/main.go`
3. Set output directory: `api`
4. Add environment variables

### Railway
1. Connect GitHub repository
2. Add PostgreSQL and Redis services
3. Set environment variables
4. Deploy automatically on push

### Render
1. Connect GitHub repository
2. Use Web Service with Go environment
3. Add PostgreSQL and Redis services
4. Set environment variables

### Docker
1. Build production image: `docker build -f Dockerfile.production -t food-delivery-backend .`
2. Run with docker-compose: `docker-compose -f docker-compose.production.yml up -d`
3. Access at http://localhost:8080

## Troubleshooting

### Build Fails
- Check Go version (requires 1.21+)
- Verify all dependencies in go.mod
- Check for syntax errors

### Database Connection Issues
- Verify database credentials
- Check network connectivity
- Ensure database exists

### Redis Connection Issues
- Verify Redis host and port
- Check Redis authentication
- Test Redis connectivity

### Performance Issues
- Enable connection pooling
- Configure Redis caching
- Monitor resource usage
- Scale horizontally if needed

## Monitoring

### Health Checks
- `/health` endpoint for basic health
- Monitor database connectivity
- Check Redis availability

### Logging
- Structured JSON logging
- Log levels: debug, info, warn, error
- Centralized log collection recommended

### Metrics
- Request/response times
- Error rates
- Database query performance
- Cache hit rates
