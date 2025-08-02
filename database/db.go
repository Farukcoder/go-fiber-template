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

	helpers.Debug("DSN: %s", dsn)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		helpers.Error("Failed to connect to the database", err)
		return nil, err
	}
	helpers.Success("Successfully connected to the database")

	// Use dynamic migration system instead of simple AutoMigrate
	migrator := NewDynamicMigrator(DB)

	// Detect schema changes
	operations, err := migrator.DetectChanges()
	if err != nil {
		helpers.Error("Failed to detect schema changes", err)
		return nil, err
	}

	// Execute migrations
	if err := migrator.ExecuteMigrations(operations); err != nil {
		helpers.Error("Failed to execute migrations", err)
		return nil, err
	}
	helpers.Success("All dynamic migrations completed successfully")

	// Handle foreign key constraints after migrations
	if err := createForeignKeyConstraints(); err != nil {
		helpers.Error("Failed to create foreign key constraints", err)
		return nil, err
	}
	helpers.Success("All foreign key constraints created successfully")

	// Create indexes for better performance
	if err := createIndexes(); err != nil {
		helpers.Error("Failed to create indexes", err)
		return nil, err
	}
	helpers.Success("All indexes created successfully")

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
	// User indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
		return fmt.Errorf("failed to create user email index: %w", err)
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_user_type ON users(user_type)").Error; err != nil {
		return fmt.Errorf("failed to create user user_type index: %w", err)
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)").Error; err != nil {
		return fmt.Errorf("failed to create user is_active index: %w", err)
	}

	// Log indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_logs_method ON logs(method)").Error; err != nil {
		return fmt.Errorf("failed to create log method index: %w", err)
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_logs_status_code ON logs(status_code)").Error; err != nil {
		return fmt.Errorf("failed to create log status_code index: %w", err)
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs(created_at)").Error; err != nil {
		return fmt.Errorf("failed to create log created_at index: %w", err)
	}

	return nil
}

// createForeignKeyConstraints creates foreign key constraints after auto migration
func createForeignKeyConstraints() error {
	// Define constraints with their names for checking existence
	constraints := []struct {
		name string
		sql  string
	}{
		// Add foreign key constraints here as needed
		// Example:
		// {
		//     name: "fk_orders_user_id",
		//     sql:  "ALTER TABLE orders ADD CONSTRAINT fk_orders_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE",
		// },
	}

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
			helpers.Warning("Failed to check constraint existence: %s - Error: %v", constraint.name, err)
			continue
		}

		if !exists {
			if err := DB.Exec(constraint.sql).Error; err != nil {
				helpers.Warning("Failed to create constraint: %s - Error: %v", constraint.name, err)
			} else {
				helpers.Success("Successfully created constraint: %s", constraint.name)
			}
		} else {
			helpers.Debug("Constraint already exists: %s", constraint.name)
		}
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
