// Package main provides database migration tool for the GoEdu Control Testing Platform.
// It handles database schema creation, updates, and index management.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/config"
	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/migrations"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/database"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/logger"
)

// main is the entry point for the migration tool.
// It connects to the database and creates necessary indexes and collections.
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(&cfg.Logger)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	log.Info("Starting database migration...",
		logger.String("database", cfg.Database.Database),
	)

	// Connect to database
	dbClient, err := database.NewClient(&cfg.Database, log)
	if err != nil {
		log.Error(ctx, "Failed to connect to database", err)
		os.Exit(1)
	}
	defer dbClient.Close(ctx)

	// Run migrations
	if err := runMigrations(ctx, dbClient, log); err != nil {
		log.Error(ctx, "Migration failed", err)
		os.Exit(1)
	}

	log.Info("Database migration completed successfully")
}

// runMigrations executes all database migrations using the enhanced migration system.
// This includes creating indexes, collections, schema updates, and version tracking.
//
// Parameters:
//   - ctx: Context for migration operations
//   - db: Database client
//   - log: Logger instance
//
// Returns:
//   - error: Migration error if any step fails
func runMigrations(ctx context.Context, db *database.Client, log *logger.Logger) error {
	log.Info("Initializing migration manager...")
	
	// Create migration manager
	migrationManager := migrations.NewMigrationManager(db, log)
	
	// Get current database version
	currentVersion, err := migrationManager.GetVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current database version: %w", err)
	}
	
	log.Info("Current database version",
		logger.Int("version", currentVersion),
	)
	
	// Apply all pending migrations
	if err := migrationManager.Up(ctx); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	
	// Get updated version
	newVersion, err := migrationManager.GetVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get updated database version: %w", err)
	}
	
	log.Info("Database migration completed",
		logger.Int("previous_version", currentVersion),
		logger.Int("current_version", newVersion),
	)
	
	// Display migration history
	history, err := migrationManager.GetMigrationHistory(ctx)
	if err != nil {
		log.Error(ctx, "Failed to get migration history", err)
		// Don't fail the migration for this
	} else {
		log.Info("Migration history retrieved",
			logger.Int("total_migrations", len(history)),
		)
		
		for _, record := range history {
			log.Info("Applied migration",
				logger.Int("version", record.Version),
				logger.String("description", record.Description),
				logger.Time("applied_at", record.AppliedAt),
			)
		}
	}

	return nil
}

// Future migration functions would be added here:
// - createValidationRules()
// - updateExistingDocuments()
// - addNewFields()
// - etc.