package database

import (
	"fmt"
	"os"
	"time"

	"garma_track/helpers"
)

// RunDynamicMigration runs and applies migrations
func RunDynamicMigration() error {
	db, err := InitDB()
	if err != nil {
		return fmt.Errorf("failed to initialize DB: %w", err)
	}

	migrator := NewDynamicMigrator(db)
	operations, err := migrator.DetectChanges()
	if err != nil {
		return fmt.Errorf("failed to detect schema changes: %w", err)
	}

	if len(operations) == 0 {
		helpers.Info("✅ No migrations needed - DB is up to date")
		return nil
	}

	helpers.Info("📋 Migration plan:")
	for i, op := range operations {
		helpers.Debug(fmt.Sprintf(" [%d] %s", i+1, op.Description))
	}

	return migrator.ExecuteMigrations(operations)
}

// GenerateMigrationFile generates migration SQL for manual inspection
func GenerateMigrationFile(filename string) error {
	db, err := InitDB()
	if err != nil {
		return fmt.Errorf("DB init failed: %w", err)
	}

	migrator := NewDynamicMigrator(db)
	operations, err := migrator.DetectChanges()
	if err != nil {
		return fmt.Errorf("schema diff failed: %w", err)
	}

	if len(operations) == 0 {
		helpers.Info("No migrations to generate")
		return nil
	}

	content := fmt.Sprintf("-- Migration generated on %s\n", time.Now().Format("2006-01-02 15:04:05"))
	for i, op := range operations {
		content += fmt.Sprintf("-- [%d] %s\n%s;\n\n", i+1, op.Description, op.SQL)
	}

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	helpers.Info(fmt.Sprintf("✅ Migration SQL written to: %s", filename))
	return nil
}
