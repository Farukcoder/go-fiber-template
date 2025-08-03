package database

import (
	"fmt"
	"os"

	"go-fiber-template/helpers"
	"go-fiber-template/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection with auto migration and indexing
func InitDB() (*gorm.DB, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		helpers.Warning("Error loading .env file")
	}

	// Get database configuration from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	sslmode := os.Getenv("DB_SSLMODE") // Optional: "disable", "require", etc.
	if sslmode == "" {
		sslmode = "disable"
	}

	// Build PostgreSQL DSN string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, database, sslmode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		helpers.Error("Failed to connect to the database", err)
		return nil, err
	}
	helpers.Success("Successfully connected to the database")

	// Run model-wise migrations serially
	if err := RunSerialMigrations(DB); err != nil {
		helpers.Error("Failed to run migrations", err)
		return nil, err
	}

	// Handle foreign key constraints after migrations
	if err := createForeignKeyConstraints(); err != nil {
		helpers.Error("Failed to create foreign key constraints", err)
		return nil, err
	}

	// Create indexes for better performance
	if err := createIndexes(); err != nil {
		helpers.Error("Failed to create indexes", err)
		return nil, err
	}

	// Show final completion message
	helpers.Success("All migrations completed successfully!")

	return DB, nil
}

// autoMigrate runs auto migration for all models
func autoMigrate() error {
	// First, migrate models without foreign key constraints in stages

	// Stage 1: Core foundation models
	stage1Models := []interface{}{
		&models.User{},
	}

	for _, model := range stage1Models {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	// Stage 2: Remaining models
	remainingModels := []interface{}{
		// Logging
		&models.Log{},
		// Add other models here as needed
	}

	for _, model := range remainingModels {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	return nil
}

// createIndexes creates additional indexes for better performance
func createIndexes() error {
	// Define indexes to create
	indexes := []struct {
		Name        string
		Description string
		SQL         string
		CheckSQL    string
	}{
		{
			Name:        "idx_users_email",
			Description: "User email index",
			SQL:         "CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
			CheckSQL:    "SELECT COUNT(*) FROM pg_indexes WHERE indexname = 'idx_users_email'",
		},
		{
			Name:        "idx_users_user_type",
			Description: "User type index",
			SQL:         "CREATE INDEX IF NOT EXISTS idx_users_user_type ON users(user_type)",
			CheckSQL:    "SELECT COUNT(*) FROM pg_indexes WHERE indexname = 'idx_users_user_type'",
		},
		{
			Name:        "idx_users_is_active",
			Description: "User active status index",
			SQL:         "CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)",
			CheckSQL:    "SELECT COUNT(*) FROM pg_indexes WHERE indexname = 'idx_users_is_active'",
		},
		{
			Name:        "idx_logs_method",
			Description: "Log method index",
			SQL:         "CREATE INDEX IF NOT EXISTS idx_logs_method ON logs(method)",
			CheckSQL:    "SELECT COUNT(*) FROM pg_indexes WHERE indexname = 'idx_logs_method'",
		},
		{
			Name:        "idx_logs_status_code",
			Description: "Log status code index",
			SQL:         "CREATE INDEX IF NOT EXISTS idx_logs_status_code ON logs(status_code)",
			CheckSQL:    "SELECT COUNT(*) FROM pg_indexes WHERE indexname = 'idx_logs_status_code'",
		},
		{
			Name:        "idx_logs_created_at",
			Description: "Log created_at index",
			SQL:         "CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs(created_at)",
			CheckSQL:    "SELECT COUNT(*) FROM pg_indexes WHERE indexname = 'idx_logs_created_at'",
		},
	}

	// Check which indexes need to be created
	indexesToCreate := []struct {
		Name        string
		Description string
		SQL         string
		CheckSQL    string
	}{}

	for _, index := range indexes {
		var count int64
		err := DB.Raw(index.CheckSQL).Scan(&count).Error
		if err != nil || count == 0 {
			indexesToCreate = append(indexesToCreate, index)
		}
	}

	// If no indexes need to be created, return silently
	if len(indexesToCreate) == 0 {
		return nil
	}

	// Create indexes silently, only log errors
	errorCount := 0

	for _, index := range indexesToCreate {
		if err := DB.Exec(index.SQL).Error; err != nil {
			helpers.Error(fmt.Sprintf("Failed to create %s", index.Description), err)
			errorCount++
		}
	}

	if errorCount > 0 {
		return fmt.Errorf("index creation completed with %d errors", errorCount)
	}

	return nil
}

// createForeignKeyConstraints creates foreign key constraints after auto migration
func createForeignKeyConstraints() error {
	// Define constraints with their names for checking existence
	constraints := []struct {
		name        string
		description string
		sql         string
	}{
		// Add foreign key constraints here as needed
		// Example:
		// {
		//     name: "fk_orders_user_id",
		//     description: "Orders user foreign key constraint",
		//     sql:  "ALTER TABLE orders ADD CONSTRAINT fk_orders_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE",
		// },
	}

	if len(constraints) == 0 {
		return nil
	}

	errorCount := 0

	for _, constraint := range constraints {
		// Check if constraint already exists
		var exists bool
		checkSQL := `
            SELECT EXISTS (
                SELECT 1 FROM pg_constraint 
                WHERE conname = $1 
            )`

		err := DB.Raw(checkSQL, constraint.name).Scan(&exists).Error
		if err != nil {
			helpers.Error(fmt.Sprintf("Failed to check constraint existence: %s", constraint.name), err)
			errorCount++
			continue
		}

		if !exists {
			if err := DB.Exec(constraint.sql).Error; err != nil {
				helpers.Error(fmt.Sprintf("Failed to create constraint: %s", constraint.name), err)
				errorCount++
			}
		}
	}

	if errorCount > 0 {
		return fmt.Errorf("constraint creation completed with %d errors", errorCount)
	}

	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Legacy function for backward compatibility
func ConnectDB() (*gorm.DB, error) {
	return InitDB()
}
