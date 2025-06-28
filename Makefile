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

# Docker commands
docker-build:
	docker build -t food-delivery-app .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f

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

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build files"
	@echo "  deps         - Download dependencies"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker containers"
	@echo "  dev          - Run in development mode with air"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  security     - Run security check"
	@echo "  help         - Show this help"

