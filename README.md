# Go Fiber Template

A robust and scalable Go API template built with the Fiber framework, featuring JWT authentication, role-based user management, database integration, structured logging, and comprehensive middleware support. This template provides a production-ready foundation for building high-performance REST APIs.

## ✨ Features

- **⚡ High Performance**: Built with Fiber v2 framework for exceptional speed
- **🔐 JWT Authentication**: Secure token-based authentication system
- **👥 Role-Based Access**: Multi-tier user system (System Admin, Garments Admin, Department Admin, Employee)
- **🗄️ Database Integration**: PostgreSQL with GORM ORM and automatic migrations
- **✅ Input Validation**: Comprehensive request validation using go-playground/validator
- **📝 Structured Logging**: Multi-level logging system with file rotation
- **🌐 CORS Support**: Configurable Cross-Origin Resource Sharing
- **🛡️ Security Middleware**: Authentication and request logging middleware
- **⚙️ Environment Config**: Flexible configuration management with .env support
- **🎯 Error Handling**: Standardized error responses and validation messages
- **🔒 Password Security**: Bcrypt password hashing
- **📁 File Storage**: Organized storage system for uploads

## 🏗️ Project Structure

```
go-fiber-template/
├── controllers/          # HTTP request handlers
│   └── auth_controller.go   # Authentication endpoints
├── database/            # Database configuration and migrations
│   ├── db.go               # Database connection setup
│   └── migrations.go       # Database schema migrations
├── helpers/             # Utility functions and helpers
│   ├── global_helper.go    # Common utility functions
│   └── logger.go           # Logging configuration and methods
├── logs/                # Application logs (auto-generated)
│   ├── 2025-08-01.log     # Daily log files
│   └── ...
├── middleware/          # HTTP middleware
│   ├── auth.go             # JWT authentication middleware
│   └── logger.go           # Request logging middleware
├── models/              # Data models and database schemas
│   ├── log.go              # Log model for system logging
│   └── user.go             # User model with role-based access
├── requests/            # Request validation structures
│   ├── auth_request.go     # Authentication request validators
│   └── log_request.go      # Log request validators
├── routes/              # Route definitions and grouping
│   └── routes.go           # API route configurations
├── storage/             # File storage directory
├── main.go             # Application entry point
├── go.mod              # Go module dependencies
├── go.sum              # Dependency checksums
├── .env                # Environment configuration
└── README.md           # Project documentation
```

## 🛠️ Prerequisites

- **Go**: Version 1.23.0 or higher
- **PostgreSQL**: Version 12 or higher
- **Git**: For version control

## � Quick Start

### 1. Clone the Repository
```bash
git clone https://github.com/Farukcoder/go-fiber-template.git
cd go-fiber-template
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Environment Configuration
Create a `.env` file in the root directory with the following configuration:

```env
# Application Configuration
APP_NAME=go-fiber-template
APP_ENV=local
APP_DEBUG=true
APP_TIMEZONE=UTC
APP_PORT=8001
APP_HOST=localhost
FRONTEND_URL=http://localhost:3001

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=go-fiber-template-db
DB_USER=your_db_user
DB_PASSWORD=your_db_password

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_here
```

### 4. Database Setup
1. Create a PostgreSQL database:
   ```sql
   CREATE DATABASE "go-fiber-template-db";
   CREATE USER "your_db_user" WITH PASSWORD 'your_db_password';
   GRANT ALL PRIVILEGES ON DATABASE "go-fiber-template-db" TO "your_db_user";
   ```

2. The application will automatically run migrations on startup.

### 5. Run the Application
```bash
go run main.go
```

The API will be available at `http://localhost:8001`

## 📡 API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration

### Protected Routes
All protected routes require a valid JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## 👥 User Roles

The system supports four user types with different access levels:

- **System Admin**: Full system access and management
- **Garments Admin**: Garments-specific administrative functions
- **Department Admin**: Department-level administrative access
- **Employee**: Basic user access

## 📦 Dependencies

### Core Dependencies
```go
github.com/gofiber/fiber/v2 v2.52.9          // High-performance HTTP framework
github.com/golang-jwt/jwt/v5 v5.3.0          // JWT token handling
gorm.io/gorm v1.30.1                         // ORM for database operations
gorm.io/driver/postgres v1.6.0               // PostgreSQL driver
github.com/go-playground/validator/v10 v10.27.0  // Request validation
github.com/joho/godotenv v1.5.1              // Environment variable management
golang.org/x/crypto v0.40.0                  // Cryptographic functions
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | go-fiber-template |
| `APP_ENV` | Environment (local/staging/production) | local |
| `APP_DEBUG` | Debug mode | true |
| `APP_PORT` | Server port | 8001 |
| `APP_HOST` | Server host | localhost |
| `FRONTEND_URL` | Frontend application URL | http://localhost:3001 |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_NAME` | Database name | - |
| `DB_USER` | Database username | - |
| `DB_PASSWORD` | Database password | - |
| `JWT_SECRET` | JWT signing secret | - |

## 📝 Logging

The application features comprehensive logging with:
- Daily log file rotation
- Multiple log levels (Info, Warning, Error, Debug)
- Structured log format
- Automatic log directory creation

Log files are stored in the `logs/` directory with the format: `YYYY-MM-DD.log`

## 🧪 Development

### Project Commands
```bash
# Install dependencies
go mod download

# Run the application
go run main.go

# Build the application
go build -o app main.go

# Run with live reload (using air)
air

# Format code
go fmt ./...

# Run tests
go test ./...
```

### Adding New Features
1. Create model in `models/`
2. Add request validation in `requests/`
3. Implement controller in `controllers/`
4. Add routes in `routes/routes.go`
5. Add migrations if needed in `database/migrations.go`

## 🐳 Docker Support

You can run this application using Docker:

```bash
# Build the image
docker build -t go-fiber-template .

# Run the container
docker run -p 8001:8001 --env-file .env go-fiber-template
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Express inspired web framework
- [GORM](https://gorm.io/) - The fantastic ORM library for Golang
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation for Go

---

**Happy Coding! 🚀**

Built with ❤️ using Go and Fiber
