# Garma Track - Go Fiber API Template

A robust and scalable Go API template built with Fiber framework, featuring authentication, database management, logging, and middleware support. This template provides a solid foundation for building production-ready REST APIs.

## 🚀 Features

- **Fast & Lightweight**: Built with Fiber framework for high performance
- **Authentication System**: JWT-based authentication with user registration and login
- **Database Integration**: PostgreSQL with GORM ORM and automatic migrations
- **Request Validation**: Input validation using go-playground/validator
- **Structured Logging**: Comprehensive logging system with different log levels
- **CORS Support**: Configurable CORS middleware
- **Middleware**: Request logging and authentication middleware
- **User Management**: Multi-role user system with different user types
- **Environment Configuration**: Environment variable management with .env support
- **Error Handling**: Standardized error responses and validation errors

## 🏗️ Project Structure

```
go_templete/
├── controllers/          # HTTP request handlers
│   └── auth_controller.go
├── database/            # Database configuration and migrations
│   ├── db.go
│   ├── migrations.go
│   ├── migrator.go
│   └── seeds/
├── helpers/             # Utility functions and helpers
│   ├── global_helper.go
│   └── logger.go
├── logs/                # Application logs
├── middleware/          # HTTP middleware
│   ├── auth.go
│   └── logger.go
├── models/              # Data models
│   ├── log.go
│   └── user.go
├── requests/            # Request validation structs
│   ├── auth_request.go
│   └── log_request.go
├── routes/              # Route definitions
│   └── routes.go
├── storage/             # File storage
├── main.go             # Application entry point
├── go.mod              # Go module file
└── go.sum              # Go dependencies checksum
```

## 🛠️ Prerequisites

- Go 1.23.0 or higher
- PostgreSQL database
- Git

## 📦 Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go_templete
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   Create a `.env` file in the root directory:
   ```env
   # Application
   APP_HOST=localhost
   APP_PORT=8080
   FRONTEND_URL=http://localhost:3000

   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=your_db_name

   # JWT
   JWT_SECRET=your_jwt_secret_key
   ```

4. **Set up PostgreSQL database**
   - Create a PostgreSQL database
   - Update the database credentials in your `.env` file

5. **Run the application**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080` (or the port specified in your `.env` file).

## 🔧 Configuration

### Using Docker (Recommended)

1. **Build the Docker image**
   ```bash
   docker build -t go-fiber-template .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 --env-file .env go-fiber-template
   ```

### Manual Deployment

1. **Build the binary**
   ```bash
   go build -o go-fiber-template-db main.go
   ```

2. **Run the application**
   ```bash
   ./go-fiber-template
   ```

## 📦 Dependencies

### Main Dependencies
- **Fiber v2**: Fast HTTP framework
- **GORM v1**: ORM for database operations
- **PostgreSQL Driver**: Database driver
- **JWT v5**: JWT token handling
- **Validator v10**: Request validation
- **Godotenv**: Environment variable management

### Development Dependencies
- **bcrypt**: Password hashing
- **crypto**: Cryptographic functions

**Happy Coding! 🎉**
