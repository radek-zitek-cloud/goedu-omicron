// Package migrations provides database migration functionality for the GoEdu Control Testing Platform.
// It implements a version-controlled migration system that supports schema evolution,
// rollback capabilities, and audit trail for all database changes.
package migrations

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/database"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/logger"
)

// Migration represents a single database migration with version control.
// Each migration includes both forward (Up) and rollback (Down) operations
// for safe schema evolution in production environments.
type Migration struct {
	// Version is the unique migration version number (must be sequential)
	Version int `json:"version" bson:"version"`
	
	// Description provides a human-readable description of the migration
	Description string `json:"description" bson:"description"`
	
	// Up performs the forward migration operation
	Up func(ctx context.Context, db *database.Client) error `json:"-" bson:"-"`
	
	// Down performs the rollback migration operation  
	Down func(ctx context.Context, db *database.Client) error `json:"-" bson:"-"`
}

// MigrationRecord represents a migration that has been applied to the database.
// This is stored in the migrations collection for version tracking.
type MigrationRecord struct {
	Version     int       `json:"version" bson:"version"`
	Description string    `json:"description" bson:"description"`
	AppliedAt   time.Time `json:"applied_at" bson:"applied_at"`
	Checksum    string    `json:"checksum" bson:"checksum"`
}

// MigrationManager handles the execution and tracking of database migrations.
// It ensures migrations are applied in order and maintains version history.
type MigrationManager struct {
	db         *database.Client
	logger     *logger.Logger
	migrations []Migration
}

// NewMigrationManager creates a new migration manager with the provided database client.
//
// Parameters:
//   - db: Database client for migration operations
//   - log: Logger instance for migration logging
//
// Returns:
//   - *MigrationManager: Configured migration manager
//
// Example:
//   manager := NewMigrationManager(dbClient, logger)
//   if err := manager.Up(ctx); err != nil {
//       log.Fatal("Migration failed", err)
//   }
func NewMigrationManager(db *database.Client, log *logger.Logger) *MigrationManager {
	return &MigrationManager{
		db:         db,
		logger:     log,
		migrations: getAllMigrations(),
	}
}

// Up applies all pending migrations to the database.
// It checks the current database version and applies any missing migrations.
//
// Parameters:
//   - ctx: Context for migration operations with timeout
//
// Returns:
//   - error: Migration execution error
//
// Example:
//   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
//   defer cancel()
//   
//   if err := manager.Up(ctx); err != nil {
//       return fmt.Errorf("migration failed: %w", err)
//   }
func (m *MigrationManager) Up(ctx context.Context) error {
	m.logger.Info("Starting database migrations...")

	// Ensure migrations collection exists and is indexed
	if err := m.ensureMigrationsCollection(ctx); err != nil {
		return fmt.Errorf("failed to ensure migrations collection: %w", err)
	}

	// Get current database version
	currentVersion, err := m.getCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	m.logger.Info("Current database version",
		logger.Int("version", currentVersion),
	)

	// Sort migrations by version
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	// Apply pending migrations
	applied := 0
	for _, migration := range m.migrations {
		if migration.Version <= currentVersion {
			continue // Already applied
		}

		m.logger.Info("Applying migration",
			logger.Int("version", migration.Version),
			logger.String("description", migration.Description),
		)

		start := time.Now()

		// Apply the migration
		if err := migration.Up(ctx, m.db); err != nil {
			return fmt.Errorf("migration %d failed: %w", migration.Version, err)
		}

		// Record the migration
		if err := m.recordMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		applied++
		duration := time.Since(start)

		m.logger.Info("Migration applied successfully",
			logger.Int("version", migration.Version),
			logger.Duration("duration", duration),
		)
	}

	if applied == 0 {
		m.logger.Info("No pending migrations to apply")
	} else {
		m.logger.Info("Database migrations completed",
			logger.Int("applied_count", applied),
		)
	}

	return nil
}

// Down rolls back the last applied migration.
// This should be used with caution in production environments.
//
// Parameters:
//   - ctx: Context for rollback operations with timeout
//
// Returns:
//   - error: Rollback execution error
//
// Example:
//   if err := manager.Down(ctx); err != nil {
//       return fmt.Errorf("rollback failed: %w", err)
//   }
func (m *MigrationManager) Down(ctx context.Context) error {
	m.logger.Warn("Starting migration rollback...")

	// Get current version
	currentVersion, err := m.getCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if currentVersion == 0 {
		m.logger.Info("No migrations to rollback")
		return nil
	}

	// Find the migration to rollback
	var targetMigration *Migration
	for i := range m.migrations {
		if m.migrations[i].Version == currentVersion {
			targetMigration = &m.migrations[i]
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration version %d not found", currentVersion)
	}

	m.logger.Warn("Rolling back migration",
		logger.Int("version", targetMigration.Version),
		logger.String("description", targetMigration.Description),
	)

	start := time.Now()

	// Execute rollback
	if err := targetMigration.Down(ctx, m.db); err != nil {
		return fmt.Errorf("rollback of migration %d failed: %w", targetMigration.Version, err)
	}

	// Remove migration record
	if err := m.removeMigrationRecord(ctx, targetMigration.Version); err != nil {
		return fmt.Errorf("failed to remove migration record %d: %w", targetMigration.Version, err)
	}

	duration := time.Since(start)
	m.logger.Warn("Migration rolled back successfully",
		logger.Int("version", targetMigration.Version),
		logger.Duration("duration", duration),
	)

	return nil
}

// GetVersion returns the current database schema version.
//
// Parameters:
//   - ctx: Context for database operations
//
// Returns:
//   - int: Current database version
//   - error: Version retrieval error
//
// Example:
//   version, err := manager.GetVersion(ctx)
//   if err != nil {
//       return fmt.Errorf("failed to get version: %w", err)
//   }
//   log.Info("Database version", zap.Int("version", version))
func (m *MigrationManager) GetVersion(ctx context.Context) (int, error) {
	return m.getCurrentVersion(ctx)
}

// GetMigrationHistory returns the history of applied migrations.
//
// Parameters:
//   - ctx: Context for database operations
//
// Returns:
//   - []MigrationRecord: List of applied migrations
//   - error: History retrieval error
//
// Example:
//   history, err := manager.GetMigrationHistory(ctx)
//   if err != nil {
//       return fmt.Errorf("failed to get history: %w", err)
//   }
//   
//   for _, record := range history {
//       log.Info("Applied migration", zap.Int("version", record.Version))
//   }
func (m *MigrationManager) GetMigrationHistory(ctx context.Context) ([]MigrationRecord, error) {
	collection := m.db.Collection("migrations")

	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"version": 1}))
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}
	defer cursor.Close(ctx)

	var records []MigrationRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("failed to decode migrations: %w", err)
	}

	return records, nil
}

// Private helper methods

// ensureMigrationsCollection creates the migrations collection and necessary indexes.
func (m *MigrationManager) ensureMigrationsCollection(ctx context.Context) error {
	collection := m.db.Collection("migrations")

	// Create unique index on version
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"version", 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create migrations index: %w", err)
	}

	return nil
}

// getCurrentVersion retrieves the current database schema version.
func (m *MigrationManager) getCurrentVersion(ctx context.Context) (int, error) {
	collection := m.db.Collection("migrations")

	// Find the highest version number
	opts := options.FindOne().SetSort(bson.M{"version": -1})
	var record MigrationRecord

	err := collection.FindOne(ctx, bson.M{}, opts).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil // No migrations applied yet
		}
		return 0, fmt.Errorf("failed to query current version: %w", err)
	}

	return record.Version, nil
}

// recordMigration stores a migration record in the database.
func (m *MigrationManager) recordMigration(ctx context.Context, migration Migration) error {
	collection := m.db.Collection("migrations")

	record := MigrationRecord{
		Version:     migration.Version,
		Description: migration.Description,
		AppliedAt:   time.Now().UTC(),
		Checksum:    fmt.Sprintf("%d-%s", migration.Version, migration.Description), // Simple checksum
	}

	_, err := collection.InsertOne(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to insert migration record: %w", err)
	}

	return nil
}

// removeMigrationRecord removes a migration record from the database.
func (m *MigrationManager) removeMigrationRecord(ctx context.Context, version int) error {
	collection := m.db.Collection("migrations")

	result, err := collection.DeleteOne(ctx, bson.M{"version": version})
	if err != nil {
		return fmt.Errorf("failed to delete migration record: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("migration record %d not found", version)
	}

	return nil
}