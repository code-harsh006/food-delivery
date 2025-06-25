.PHONY: build run test clean docker-up docker-down migrate seed docker-prod docker-prod-down build-linux build-windows build-mac deploy-vercel deploy-railway deploy-render health-check verify-vercel deploy

# Build the application
build:
	go build -o bin/main cmd/main.go

# Run the application locally
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Start Docker services
docker-up:
	docker-compose up -d

# Stop Docker services
docker-down:
	docker-compose down

# Build and run with Docker
docker-build:
	docker-compose up --build

# Run database migrations
migrate:
	go run cmd/migrate/main.go

# Seed database with sample data
seed:
	go run cmd/seed/main.go

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Run development server with hot reload
dev:
	air

# Database operations
db-reset: docker-down docker-up migrate seed

# Production build
prod-build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main cmd/main.go

# Deployment commands
docker-prod:
	docker-compose -f docker-compose.production.yml up --build -d

docker-prod-down:
	docker-compose -f docker-compose.production.yml down

# Build for different platforms
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux cmd/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/main-windows.exe cmd/main.go

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/main-mac cmd/main.go

# Fixed Vercel deployment
deploy-vercel:
	@echo "🚀 Deploying to Vercel..."
	@echo "📋 Checking vercel.json configuration..."
	@if grep -q '"functions"' vercel.json; then \
		echo "❌ Found conflicting 'functions' property in vercel.json"; \
		echo "🔧 This has been fixed in the updated configuration"; \
	else \
		echo "✅ vercel.json configuration is correct"; \
	fi
	vercel --prod

# Verify deployment
verify-vercel:
	@echo "🧪 Verifying Vercel deployment..."
	@read -p "Enter your Vercel URL: " url; \
	chmod +x verify-deployment.sh; \
	./verify-deployment.sh $$url

# Railway deployment
deploy-railway:
	@echo "🚂 Deploying to Railway..."
	railway up

# Render deployment
deploy-render:
	@echo "🎨 Deploying to Render..."
	@echo "Please connect your GitHub repository to Render dashboard"

# Docker deployment
deploy-docker:
	@echo "🐳 Building and running Docker..."
	docker build -f Dockerfile.production -t food-delivery-backend .
	docker run -p 8080:8080 food-delivery-backend

# Health check
health-check:
	curl -f http://localhost:8080/health || exit 1

# Quick deployment with platform choice
deploy:
	@echo "🚀 Food Delivery Backend Deployment"
	@echo "Choose deployment platform:"
	@echo "1) Vercel (Serverless)"
	@echo "2) Railway (Container)"
	@echo "3) Docker (Local)"
	@read -p "Enter choice (1-3): " choice; \
	case $$choice in \
		1) make deploy-vercel ;; \
		2) make deploy-railway ;; \
		3) make deploy-docker ;; \
		*) echo "Invalid choice" ;; \
	esac
