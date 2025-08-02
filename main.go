package main

import (
	"fmt"
	"go-fiber-template/database"
	"go-fiber-template/helpers"
	"go-fiber-template/middleware"
	"go-fiber-template/routes"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		helpers.Warning("⚠️  Could not load .env file, using system environment variables")
	}

	// Initialize logger
	if err := helpers.InitLogger(); err != nil {
		helpers.Error("❌ Could not initialize logger: %v", err)
		return
	}

	// Create Fiber app with config
	app := fiber.New(fiber.Config{
		ReadBufferSize:  32 * 1024,        // 32 KB
		WriteBufferSize: 32 * 1024,        // 32 KB
		ReadTimeout:     30 * time.Second, // Read timeout
		WriteTimeout:    30 * time.Second, // Write timeout
		BodyLimit:       50 * 1024 * 1024, // 50 MB
	})

	// Setup CORS
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
		helpers.Warning("FRONTEND_URL not set. Defaulting to: %s", frontendURL)
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     frontendURL,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true, // Required if using cookies or Authorization headers
	}))

	// Middleware
	app.Use(middleware.RequestLogger())

	// Connect to the database
	db, err := database.ConnectDB()
	if err != nil {
		helpers.Error("❌ Could not connect to database: %v", err)
		return
	}

	// Setup routes
	routes.SetupRoutes(app, db)

	// Server start logs
	appHost := os.Getenv("APP_HOST")
	appPort := os.Getenv("APP_PORT")
	serverAddress := appHost + ":" + appPort

	// Print to console directly for immediate visibility
	fmt.Printf("\n🚀 Server is running on http://%s\n", serverAddress)
	fmt.Println("📝 All HTTP requests will be logged below:")
	fmt.Println(strings.Repeat("=", 80))

	helpers.Success("🚀 Server is running on http://" + serverAddress)
	helpers.Success("\n\t******************************************************************************************\n")

	// Start server
	if err := app.Listen(serverAddress); err != nil {
		helpers.Error("❌ Failed to start server: %v", err)
	}
}
