# Fixed Deployment Guide

## The Issue
The error "Could not find an exported function" occurs when deployment platforms expect a serverless function but find a traditional Go server.

## Solutions by Platform

### 1. Vercel (Serverless)
✅ **Fixed with `api/index.go`**

\`\`\`bash
# Deploy to Vercel
npm i -g vercel
vercel --prod
\`\`\`

**Configuration**: Uses `api/index.go` as serverless function entry point.

### 2. Railway (Container)
✅ **Fixed with `Dockerfile.railway`**

\`\`\`bash
# Deploy to Railway
npm i -g @railway/cli
railway login
railway up
\`\`\`

**Configuration**: Uses Docker container with proper Go build.

### 3. Render (Native Go)
✅ **Fixed with `render.yaml`**

\`\`\`bash
# Connect GitHub repo to Render
# Automatic deployment on push
\`\`\`

**Configuration**: Native Go runtime with proper build commands.

### 4. Heroku (Container)
✅ **Fixed with `heroku.yml` and `Dockerfile.heroku`**

\`\`\`bash
# Deploy to Heroku
heroku create your-app-name
heroku stack:set container
git push heroku main
\`\`\`

**Configuration**: Container deployment with PostgreSQL and Redis addons.

### 5. Docker (Universal)
✅ **Works everywhere**

\`\`\`bash
# Local deployment
make docker-prod

# Or manual
docker-compose -f docker-compose.production.yml up -d
\`\`\`

## Quick Deploy Commands

### Option 1: Vercel (Recommended for Serverless)
\`\`\`bash
vercel --prod
\`\`\`

### Option 2: Railway (Recommended for Full Apps)
\`\`\`bash
railway up
\`\`\`

### Option 3: Docker (Universal)
\`\`\`bash
docker build -t food-delivery-backend .
docker run -p 8080:8080 food-delivery-backend
\`\`\`

## Environment Variables Required

\`\`\`env
# Database (Required)
DB_HOST=your-postgres-host
DB_USER=your-db-user  
DB_PASSWORD=your-db-password
DB_NAME=food_delivery

# Redis (Required)
REDIS_HOST=your-redis-host

# JWT (Required)
JWT_SECRET=your-super-secret-key

# Optional
PORT=8080
GIN_MODE=release
\`\`\`

## Testing Deployment

After deployment, test these endpoints:

\`\`\`bash
# Health check
curl https://your-app-url/health

# API test
curl https://your-app-url/api/v1/products

# Register user
curl -X POST https://your-app-url/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","first_name":"John","last_name":"Doe","phone":"+1234567890"}'
\`\`\`

## Platform-Specific Notes

### Vercel
- ✅ Serverless function in `api/index.go`
- ✅ Automatic HTTPS
- ✅ Global CDN
- ⚠️ 10-second timeout limit
- ⚠️ Cold starts

### Railway
- ✅ Full container support
- ✅ Built-in PostgreSQL/Redis
- ✅ Automatic deployments
- ✅ No timeout limits
- ✅ Always warm

### Render
- ✅ Native Go support
- ✅ Built-in PostgreSQL/Redis
- ✅ Automatic SSL
- ✅ Health checks
- ⚠️ Slower cold starts

### Heroku
- ✅ Container support
- ✅ Add-on ecosystem
- ✅ Mature platform
- ⚠️ Dyno sleeping (free tier)
- ⚠️ More expensive

## Troubleshooting

### Build Fails
1. Check Go version (1.21+ required)
2. Verify go.mod dependencies
3. Check platform-specific Dockerfile

### Database Issues
1. Verify connection strings
2. Check firewall rules
3. Ensure database exists

### Function Export Error
1. Use `api/index.go` for serverless
2. Use `cmd/main.go` for containers
3. Check platform configuration

## Recommended Deployment Strategy

1. **Development**: Local Docker
2. **Staging**: Railway (easy setup)
3. **Production**: Railway or Render (full features)
4. **Serverless**: Vercel (if timeout limits work for you)
\`\`\`

Let's also update the package.json to be more deployment-friendly:
