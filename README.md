# Food Delivery Backend

A comprehensive Go-based backend API for food delivery services with support for payments, maps, tracking, and real-time notifications.

## Features

- 🔐 **Authentication & Authorization** - JWT-based auth with role-based access
- 💳 **Payment Integration** - Stripe and PayPal support
- 🗺️ **Maps & Location** - Google Maps and Mapbox integration
- 📍 **Real-time Tracking** - Order tracking and delivery status
- 📧 **Email & SMS** - Twilio SMS and SMTP email notifications
- 🔔 **Push Notifications** - Firebase Cloud Messaging
- 📁 **File Upload** - AWS S3 compatible storage
- 🗄️ **Database** - PostgreSQL with Redis caching
- 📊 **Monitoring** - Prometheus and Grafana
- 🐳 **Docker** - Multi-stage builds for dev and production

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL
- **Cache**: Redis
- **ORM**: GORM
- **Container**: Docker & Docker Compose
- **Monitoring**: Prometheus + Grafana

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/code-harsh006/food-delivery.git
   cd food-delivery
   ```

2. **Setup environment**
   ```bash
   make setup
   # or manually:
   cp config.env.example .env
   # Edit .env with your actual values
   ```

3. **Run with Docker (Recommended)**
   ```bash
   # Development environment
   make docker-dev-run
   
   # Or production environment
   make docker-prod-run
   ```

4. **Run locally**
   ```bash
   # Install dependencies
   make deps
   
   # Run the application
   make run
   ```

## Environment Variables

Copy `config.env.example` to `.env` and configure the following:

### Required Variables
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `JWT_SECRET` - Secret key for JWT tokens

### Payment Configuration
- `STRIPE_KEY` - Stripe secret key
- `STRIPE_PUBLISHABLE_KEY` - Stripe publishable key
- `PAYPAL_CLIENT_ID` - PayPal client ID
- `PAYPAL_CLIENT_SECRET` - PayPal client secret

### Maps & Location
- `GOOGLE_MAPS_API_KEY` - Google Maps API key
- `MAPBOX_ACCESS_TOKEN` - Mapbox access token

### Tracking & Delivery
- `TRACKING_API_KEY` - Tracking service API key
- `DELIVERY_RADIUS_KM` - Maximum delivery radius

### Email & SMS
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`
- `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN`, `TWILIO_PHONE_NUMBER`

### File Upload
- `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`, `AWS_S3_BUCKET`

## API Endpoints

### Health & Status
- `GET /health` - Health check
- `GET /api` - API information

### Authentication
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout
- `GET /auth/profile` - Get user profile

### Orders
- `GET /orders` - List orders
- `POST /orders` - Create order
- `GET /orders/:id` - Get order details
- `PUT /orders/:id/status` - Update order status

### Products
- `GET /products` - List products
- `POST /products` - Create product (admin)
- `GET /products/:id` - Get product details

### Tracking
- `GET /tracking/:order_id` - Get order tracking
- `POST /tracking/:order_id/update` - Update tracking location

## Docker Commands

### Development
```bash
# Build development image
make docker-dev-build

# Run development environment
make docker-dev-run

# View logs
make docker-dev-logs

# Stop development environment
make docker-dev-stop

# Access container shell
make docker-dev-shell
```

### Production
```bash
# Build production image
make docker-prod-build

# Run production environment
make docker-prod-run

# View logs
make docker-prod-logs

# Stop production environment
make docker-prod-stop
```

### General Docker
```bash
# Clean Docker resources
make docker-clean

# Rebuild and restart
make docker-rebuild
```

## Development Tools

### Available Services (Development)
- **App**: http://localhost:8080
- **Adminer** (Database): http://localhost:8081
- **Redis Commander**: http://localhost:8082
- **Mailhog** (Email): http://localhost:8025
- **MinIO Console**: http://localhost:9001
- **Jaeger** (Tracing): http://localhost:16686

### Available Services (Production)
- **App**: http://localhost:8080
- **Nginx**: http://localhost:80, https://localhost:443
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000

## Project Structure

```
food-delivery/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── api/                 # HTTP handlers
│   ├── auth/                # Authentication module
│   ├── cart/                # Shopping cart module
│   ├── notification/        # Notifications module
│   ├── order/               # Orders module
│   ├── product/             # Products module
│   ├── search/              # Search module
│   ├── user/                # Users module
│   └── vendor/              # Vendors module
├── pkg/
│   ├── config/              # Configuration management
│   ├── db/                  # Database models and migrations
│   ├── logger/              # Logging utilities
│   ├── middleware/          # HTTP middleware
│   ├── response/            # Response utilities
│   └── utils/               # Utility functions
├── Dockerfile               # Multi-stage Docker build
├── docker-compose.yml       # Main Docker Compose
├── docker-compose.override.yml # Development overrides
├── docker-compose.prod.yml  # Production configuration
├── Makefile                 # Build and deployment commands
└── config.env.example       # Environment variables template
```

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run specific test
go test -v ./internal/auth
```

## Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Security check
make security
```

## Database

```bash
# Run migrations
make migrate

# Seed database
make seed
```

## Monitoring

The application includes monitoring with Prometheus and Grafana:

- **Prometheus**: Collects metrics from the application
- **Grafana**: Visualizes metrics and provides dashboards
- **Jaeger**: Distributed tracing for debugging

## Deployment

### Render
```bash
make deploy-render
```

### Docker
```bash
# Build production image
make docker-prod-build

# Run production stack
make docker-prod-run
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run tests and linting
6. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support and questions, please open an issue on GitHub. 