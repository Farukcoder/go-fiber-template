package database

import (
	"fmt"
	"os"

	"go-fiber-template/helpers"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global DB instance
var DB *gorm.DB

// InitDB initializes the DB connection, runs migrations, constraints and indexes
func InitDB() (*gorm.DB, error) {
	// Load .env variables
	if err := godotenv.Load(); err != nil {
		helpers.Warning("Could not load .env file (skipping)")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	helpers.Info("Connecting to database...")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		helpers.Error("Failed to connect to database: %v", err)
		return nil, err
	}
	helpers.Info("Connected to database successfully")

	// Run migrations through the migrator
	migrator := NewDynamicMigrator(DB)
	if err := migrator.ExecuteMigrations(nil); err != nil {
		helpers.Error("Failed to run migrations: %v", err)
		return nil, fmt.Errorf("migration failed: %w", err)
	}
	helpers.Info("Database migrations completed successfully")

	return DB, nil
}

// ConnectDB is a legacy alias for InitDB
func ConnectDB() (*gorm.DB, error) {
	return InitDB()
}

// GetDB returns the global DB instance
func GetDB() *gorm.DB {
	return DB
}
