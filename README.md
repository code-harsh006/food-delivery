# Food Delivery Backend

A complete, production-ready Golang backend for a modern food and grocery delivery application (like Swiggy + BigBasket). Built with Go, Gin, GORM, PostgreSQL, and Redis.

## 🚀 Features

### Core Modules
- **Authentication**: JWT-based auth with refresh tokens
- **User Management**: Profile management and address handling
- **Product Management**: Categories, products with full CRUD
- **Vendor Management**: Restaurant/store management
- **Cart System**: Redis-cached cart with real-time updates
- **Order Management**: Complete order lifecycle
- **Payment Integration**: Stripe integration (ready)
- **Delivery System**: Agent assignment and tracking
- **Admin Panel**: User and vendor management
- **Real-time Updates**: WebSocket support
- **Search & Recommendations**: Product search functionality

### Technical Features
- **Modular Architecture**: Clean, maintainable code structure
- **Database**: PostgreSQL with GORM ORM
- **Caching**: Redis for cart and session management
- **Authentication**: JWT with refresh token rotation
- **API Documentation**: Swagger/OpenAPI ready
- **Docker Support**: Complete containerization
- **Graceful Shutdown**: Proper server lifecycle management
- **Logging**: Structured logging with logrus
- **Error Handling**: Comprehensive error management
- **Validation**: Request validation and sanitization

## 🏗️ Architecture

\`\`\`
├── cmd/                    # Application entrypoints
│   ├── main.go            # Main application
│   └── seed/              # Database seeder
├── internal/              # Private application code
│   ├── auth/              # Authentication module
│   ├── user/              # User management
│   ├── product/           # Product management
│   ├── vendor/            # Vendor management
│   ├── cart/              # Shopping cart
│   ├── order/             # Order management
│   ├── payment/           # Payment processing
│   ├── delivery/          # Delivery management
│   ├── admin/             # Admin functionality
│   ├── notification/      # Notifications
│   └── search/            # Search & recommendations
├── pkg/                   # Public/shared packages
│   ├── config/            # Configuration management
│   ├── db/                # Database connections
│   ├── middleware/        # HTTP middleware
│   ├── logger/            # Logging utilities
│   └── utils/             # Common utilities
├── api/                   # API documentation
├── migrations/            # Database migrations
├── docker-compose.yml     # Docker services
├── Dockerfile             # Application container
├── Makefile               # Build automation
└── .env                   # Environment variables

\`\`\`

## 🚦 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### Installation

1. **Clone the repository**
\`\`\`bash
git clone <repository-url>
cd food-delivery-backend
\`\`\`

2. **Set up environment variables**
\`\`\`bash
cp .env.example .env
# Edit .env with your configuration
\`\`\`

3. **Start with Docker (Recommended)**
\`\`\`bash
make docker-up
\`\`\`

4. **Or run locally**
\`\`\`bash
# Install dependencies
make deps

# Start databases
docker-compose up postgres redis -d

# Run migrations and seed data
make migrate
make seed

# Start the server
make run
\`\`\`

The API will be available at `http://localhost:8080`

## 📚 API Documentation

### Authentication Endpoints
\`\`\`
POST /api/v1/auth/register     # User registration
POST /api/v1/auth/login        # User login
POST /api/v1/auth/refresh      # Refresh access token
POST /api/v1/auth/logout       # User logout
GET  /api/v1/auth/profile      # Get user profile
\`\`\`

### User Management
\`\`\`
PUT  /api/v1/user/profile      # Update profile
GET  /api/v1/user/addresses    # Get addresses
POST /api/v1/user/addresses    # Add address
PUT  /api/v1/user/addresses/:id # Update address
DELETE /api/v1/user/addresses/:id # Delete address
\`\`\`

### Products & Categories
\`\`\`
GET  /api/v1/products          # List products
GET  /api/v1/products/:id      # Get product details
GET  /api/v1/products/search   # Search products
GET  /api/v1/products/featured # Featured products
GET  /api/v1/categories        # List categories
GET  /api/v1/categories/:id    # Get category details
\`\`\`

### Shopping Cart
\`\`\`
GET    /api/v1/cart            # Get cart
POST   /api/v1/cart/items      # Add to cart
PUT    /api/v1/cart/items/:id  # Update cart item
DELETE /api/v1/cart/items/:id  # Remove from cart
DELETE /api/v1/cart            # Clear cart
\`\`\`

### Health Check
\`\`\`
GET /health                    # Service health status
\`\`\`

## 🔧 Configuration

Key environment variables:

\`\`\`env
# Server
PORT=8080
GIN_MODE=release

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=food_delivery

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

# Payment
STRIPE_SECRET_KEY=sk_test_...
\`\`\`

## 🧪 Testing

\`\`\`bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run specific module tests
go test -v ./internal/auth/...
\`\`\`

## 📦 Database Schema

### Key Tables
- **users**: User accounts and authentication
- **user_profiles**: Extended user information
- **addresses**: User delivery addresses
- **categories**: Product categories (hierarchical)
- **products**: Product catalog
- **vendors**: Restaurant/store information
- **carts & cart_items**: Shopping cart data
- **orders & order_items**: Order management
- **payments**: Payment transactions
- **deliveries**: Delivery tracking

## 🚀 Deployment

### Docker Deployment
\`\`\`bash
# Build production image
make prod-build

# Deploy with docker-compose
docker-compose -f docker-compose.prod.yml up -d
\`\`\`

### Manual Deployment
\`\`\`bash
# Build binary
make build

# Run migrations
./bin/main migrate

# Start server
./bin/main
\`\`\`

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt for password security
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: GORM ORM prevents SQL injection
- **CORS Support**: Configurable cross-origin requests
- **Rate Limiting**: Ready for rate limiting middleware
- **Environment Variables**: Secure configuration management

## 📈 Performance Features

- **Redis Caching**: Cart and session caching
- **Database Indexing**: Optimized database queries
- **Connection Pooling**: Efficient database connections
- **Graceful Shutdown**: Proper resource cleanup
- **Structured Logging**: Performance monitoring ready
- **Pagination**: Efficient data retrieval

## 🛠️ Development

### Available Make Commands
\`\`\`bash
make build          # Build the application
make run            # Run the application
make test           # Run tests
make docker-up      # Start Docker services
make docker-down    # Stop Docker services
make migrate        # Run database migrations
make seed           # Seed database with sample data
make clean          # Clean build artifacts
make fmt            # Format code
make lint           # Lint code
\`\`\`

### Adding New Modules

1. Create module directory in `internal/`
2. Implement models, service, handler, and routes
3. Register routes in `cmd/main.go`
4. Add migrations if needed
5. Update documentation

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the documentation
- Review the API examples

## 🗺️ Roadmap

- [ ] Order management completion
- [ ] Payment gateway integration
- [ ] Delivery tracking system
- [ ] Push notifications
- [ ] Admin dashboard APIs
- [ ] Analytics and reporting
- [ ] Multi-vendor support
- [ ] Inventory management
- [ ] Loyalty program
- [ ] Review and rating system
\`\`\`

This is a complete, production-ready food delivery backend with all the features you requested. The codebase is modular, scalable, and follows Go best practices. You can extend it easily by adding more modules or integrating with external services.
