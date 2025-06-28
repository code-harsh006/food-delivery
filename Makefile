.PHONY: build run test clean docker-build docker-run docker-stop migrate seed

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=food-delivery
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v cmd/main.go

# Run the application
run:
	$(GOCMD) run cmd/main.go

# Test the application
test:
	$(GOTEST) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v cmd/main.go

# Docker Development Commands
docker-dev-build:
	docker build --target development -t food-delivery-dev .

docker-dev-run:
	docker-compose -f doker-compose.yml -f docker-compose.override.yml up -d

docker-dev-stop:
	docker-compose -f doker-compose.yml -f docker-compose.override.yml down

docker-dev-logs:
	docker-compose -f doker-compose.yml -f docker-compose.override.yml logs -f

docker-dev-shell:
	docker-compose -f doker-compose.yml -f docker-compose.override.yml exec app sh

# Docker Production Commands
docker-prod-build:
	docker build --target production -t food-delivery-prod .

docker-prod-run:
	docker-compose -f docker-compose.prod.yml up -d

docker-prod-stop:
	docker-compose -f docker-compose.prod.yml down

docker-prod-logs:
	docker-compose -f docker-compose.prod.yml logs -f

# Docker General Commands
docker-build:
	docker build -t food-delivery-app .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker system prune -f
	docker volume prune -f

docker-rebuild:
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# Database commands
migrate:
	$(GOCMD) run cmd/main.go migrate

seed:
	$(GOCMD) run cmd/main.go seed

# Development commands
dev:
	air

install-air:
	go install github.com/cosmtrek/air@latest

# Linting
lint:
	golangci-lint run

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Format code
fmt:
	go fmt ./...

# Security check
security:
	gosec ./...

install-security:
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Generate swagger docs (optional)
swagger:
	swag init -g cmd/main.go

install-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest

# Production build
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o $(BINARY_NAME) cmd/main.go

# Render deployment
deploy-render:
	git add .
	git commit -m "Deploy to Render"
	git push origin main

# Environment setup
setup-env:
	@if [ ! -f .env ]; then \
		echo "Creating .env file from config.env.example..."; \
		cp config.env.example .env; \
		echo "Please update .env file with your actual values"; \
	else \
		echo ".env file already exists"; \
	fi

# Create necessary directories
setup-dirs:
	mkdir -p logs
	mkdir -p scripts
	mkdir -p monitoring/grafana/dashboards
	mkdir -p monitoring/grafana/datasources
	mkdir -p nginx/ssl

# Full setup
setup: setup-dirs setup-env deps

# Help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build files"
	@echo "  deps         - Download dependencies"
	@echo "  dev          - Run in development mode with air"
	@echo "  setup        - Full development setup"
	@echo ""
	@echo "Docker Development:"
	@echo "  docker-dev-build - Build development Docker image"
	@echo "  docker-dev-run   - Run development environment"
	@echo "  docker-dev-stop  - Stop development environment"
	@echo "  docker-dev-logs  - View development logs"
	@echo "  docker-dev-shell - Access development container shell"
	@echo ""
	@echo "Docker Production:"
	@echo "  docker-prod-build - Build production Docker image"
	@echo "  docker-prod-run   - Run production environment"
	@echo "  docker-prod-stop  - Stop production environment"
	@echo "  docker-prod-logs  - View production logs"
	@echo ""
	@echo "Docker General:"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker containers"
	@echo "  docker-logs   - View logs"
	@echo "  docker-clean  - Clean Docker resources"
	@echo "  docker-rebuild- Rebuild and restart containers"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  security     - Run security check"
	@echo ""
	@echo "Database:"
	@echo "  migrate      - Run database migrations"
	@echo "  seed         - Seed database"
	@echo ""
	@echo "Other:"
	@echo "  help         - Show this help"

