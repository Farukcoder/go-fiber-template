package database

import (
	"fmt"
	"go-fiber-template/helpers"
	"go-fiber-template/models"
	"time"

	"gorm.io/gorm"
)

// MigrationOperation represents a single migration operation
type MigrationOperation struct {
	Description string
	SQL         string
}

// DynamicMigrator handles database migrations
type DynamicMigrator struct {
	db *gorm.DB
}

// NewDynamicMigrator creates a new migrator instance
func NewDynamicMigrator(db *gorm.DB) *DynamicMigrator {
	return &DynamicMigrator{db: db}
}

// DetectChanges analyzes the database schema and returns needed migrations
func (m *DynamicMigrator) DetectChanges() ([]MigrationOperation, error) {
	var operations []MigrationOperation

	// Auto-migrate will handle the changes
	// Here we just return an empty slice since GORM handles the migrations
	return operations, nil
}

// Migration represents a migration record in the database
type Migration struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique"`
	CreatedAt time.Time
}

// ExecuteMigrations runs the migration operations
func (m *DynamicMigrator) ExecuteMigrations(operations []MigrationOperation) error {
	// Create migrations table if it doesn't exist
	if err := m.db.AutoMigrate(&Migration{}); err != nil {
		helpers.Error("Failed to create migrations table: %v", err)
		return fmt.Errorf("migrations table creation failed: %w", err)
	}

	// Auto-migrate the models
	if err := m.db.AutoMigrate(
		&models.User{},
		&models.Log{},
		// Add other models here as needed
	); err != nil {
		helpers.Error("Failed to execute migrations: %v", err)
		return fmt.Errorf("migration execution failed: %w", err)
	}

	// Record this migration
	migration := Migration{
		Name:      time.Now().UTC().Format("20060102150405"),
		CreatedAt: time.Now().UTC(),
	}
	if err := m.db.Create(&migration).Error; err != nil {
		helpers.Error("Failed to record migration: %v", err)
		return fmt.Errorf("failed to record migration: %w", err)
	}

	helpers.Info("Migrations executed successfully")
	return nil
}
