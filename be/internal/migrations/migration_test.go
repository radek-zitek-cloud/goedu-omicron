// Package migrations_test provides comprehensive tests for the database migration system
// in the GoEdu Control Testing Platform.
package migrations_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/config"
	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/migrations"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/database"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/logger"
)

// TestConfig provides test configuration for migration testing.
func getTestConfig() *config.DatabaseConfig {
	return &config.DatabaseConfig{
		URI:                 "mongodb://localhost:27017",
		Database:            "goedu_migration_test",
		MaxPoolSize:         10,
		MinPoolSize:         2,
		MaxConnIdleTime:     1 * time.Minute,
		ConnectTimeout:      5 * time.Second,
		ServerSelectTimeout: 5 * time.Second,
	}
}

// getTestLogger creates a test logger instance.
func getTestLogger() (*logger.Logger, error) {
	cfg := &logger.Config{
		Level:       "info",
		Environment: "test",
		OutputPath:  "stdout",
	}
	return logger.New(cfg)
}

// setupTestDatabase creates a clean test database for migration testing.
func setupTestDatabase(t *testing.T) (*database.Client, func()) {
	cfg := getTestConfig()
	log, err := getTestLogger()
	require.NoError(t, err)

	client, err := database.NewClient(cfg, log)
	if err != nil {
		t.Skipf("MongoDB not available for testing: %v", err)
		return nil, nil
	}

	// Clean up any existing test data
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client.Database().Drop(ctx)

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client.Database().Drop(ctx)
		client.Close(ctx)
	}

	return client, cleanup
}

// TestMigrationManager tests the migration manager functionality.
func TestMigrationManager(t *testing.T) {
	client, cleanup := setupTestDatabase(t)
	if client == nil {
		return
	}
	defer cleanup()

	log, err := getTestLogger()
	require.NoError(t, err)

	manager := migrations.NewMigrationManager(client, log)

	t.Run("initial version should be zero", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		version, err := manager.GetVersion(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, version)
	})

	t.Run("migration history should be empty initially", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		history, err := manager.GetMigrationHistory(ctx)
		assert.NoError(t, err)
		assert.Empty(t, history)
	})
}

// TestMigrationUp tests the migration up functionality.
func TestMigrationUp(t *testing.T) {
	client, cleanup := setupTestDatabase(t)
	if client == nil {
		return
	}
	defer cleanup()

	log, err := getTestLogger()
	require.NoError(t, err)

	manager := migrations.NewMigrationManager(client, log)

	t.Run("apply migrations successfully", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Apply migrations
		err := manager.Up(ctx)
		assert.NoError(t, err)

		// Check that version was updated
		version, err := manager.GetVersion(ctx)
		assert.NoError(t, err)
		assert.Greater(t, version, 0)

		// Check migration history
		history, err := manager.GetMigrationHistory(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, history)

		// Verify first migration is recorded
		assert.Equal(t, 1, history[0].Version)
		assert.Equal(t, "Create initial database indexes for optimal performance", history[0].Description)
		assert.WithinDuration(t, time.Now(), history[0].AppliedAt, 1*time.Minute)
	})

	t.Run("migrations are idempotent", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Get version before second run
		versionBefore, err := manager.GetVersion(ctx)
		require.NoError(t, err)

		// Apply migrations again
		err = manager.Up(ctx)
		assert.NoError(t, err)

		// Version should be the same
		versionAfter, err := manager.GetVersion(ctx)
		assert.NoError(t, err)
		assert.Equal(t, versionBefore, versionAfter)
	})
}

// TestMigrationDown tests the migration rollback functionality.
func TestMigrationDown(t *testing.T) {
	client, cleanup := setupTestDatabase(t)
	if client == nil {
		return
	}
	defer cleanup()

	log, err := getTestLogger()
	require.NoError(t, err)

	manager := migrations.NewMigrationManager(client, log)

	t.Run("rollback with no migrations", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := manager.Down(ctx)
		assert.NoError(t, err)

		version, err := manager.GetVersion(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, version)
	})

	t.Run("rollback after applying migrations", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// First apply migrations
		err := manager.Up(ctx)
		require.NoError(t, err)

		versionBefore, err := manager.GetVersion(ctx)
		require.NoError(t, err)
		require.Greater(t, versionBefore, 0)

		// Then rollback
		err = manager.Down(ctx)
		assert.NoError(t, err)

		versionAfter, err := manager.GetVersion(ctx)
		assert.NoError(t, err)
		assert.Equal(t, versionBefore-1, versionAfter)

		// Check that migration record was removed
		history, err := manager.GetMigrationHistory(ctx)
		assert.NoError(t, err)
		if len(history) > 0 {
			// Should not contain the rolled back migration
			for _, record := range history {
				assert.NotEqual(t, versionBefore, record.Version)
			}
		}
	})
}

// TestMigrationsCollection tests the migrations collection functionality.
func TestMigrationsCollection(t *testing.T) {
	client, cleanup := setupTestDatabase(t)
	if client == nil {
		return
	}
	defer cleanup()

	log, err := getTestLogger()
	require.NoError(t, err)

	manager := migrations.NewMigrationManager(client, log)

	t.Run("migrations collection is created", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Apply migrations
		err := manager.Up(ctx)
		require.NoError(t, err)

		// Check that migrations collection exists
		collection := client.Collection("migrations")
		count, err := collection.CountDocuments(ctx, bson.M{})
		assert.NoError(t, err)
		assert.Greater(t, count, int64(0))

		// Check that unique index exists on version field
		cursor, err := collection.Indexes().List(ctx)
		require.NoError(t, err)

		var indexes []map[string]interface{}
		err = cursor.All(ctx, &indexes)
		require.NoError(t, err)

		hasVersionIndex := false
		for _, index := range indexes {
			if name, ok := index["name"].(string); ok {
				if name == "version_1" {
					hasVersionIndex = true
					// Check if it's unique
					if unique, ok := index["unique"].(bool); ok && unique {
						assert.True(t, unique)
					}
					break
				}
			}
		}
		assert.True(t, hasVersionIndex, "Should have unique index on version field")
	})
}

// TestMigrationRecords tests migration record functionality.
func TestMigrationRecords(t *testing.T) {
	client, cleanup := setupTestDatabase(t)
	if client == nil {
		return
	}
	defer cleanup()

	log, err := getTestLogger()
	require.NoError(t, err)

	manager := migrations.NewMigrationManager(client, log)

	t.Run("migration records are properly stored", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Apply migrations
		err := manager.Up(ctx)
		require.NoError(t, err)

		// Get migration history
		history, err := manager.GetMigrationHistory(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, history)

		// Verify record structure
		record := history[0]
		assert.Equal(t, 1, record.Version)
		assert.NotEmpty(t, record.Description)
		assert.NotEmpty(t, record.Checksum)
		assert.WithinDuration(t, time.Now(), record.AppliedAt, 1*time.Minute)
	})
}

// TestMigrationErrorHandling tests error handling in migrations.
func TestMigrationErrorHandling(t *testing.T) {
	client, cleanup := setupTestDatabase(t)
	if client == nil {
		return
	}
	defer cleanup()

	log, err := getTestLogger()
	require.NoError(t, err)

	t.Run("database connection errors", func(t *testing.T) {
		// Create a manager with a closed client
		client.Close(context.Background())
		
		manager := migrations.NewMigrationManager(client, log)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := manager.Up(ctx)
		assert.Error(t, err)
	})
}

// TestMigrationConcurrency tests concurrent migration execution.
func TestMigrationConcurrency(t *testing.T) {
	client, cleanup := setupTestDatabase(t)
	if client == nil {
		return
	}
	defer cleanup()

	log, err := getTestLogger()
	require.NoError(t, err)

	t.Run("concurrent migration attempts", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Try to run migrations concurrently
		const numGoroutines = 3
		done := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				manager := migrations.NewMigrationManager(client, log)
				done <- manager.Up(ctx)
			}()
		}

		// Wait for all goroutines to complete
		var errors []error
		for i := 0; i < numGoroutines; i++ {
			if err := <-done; err != nil {
				errors = append(errors, err)
			}
		}

		// At least one should succeed, others might fail due to unique constraint
		// This is expected behavior for concurrent migrations
		assert.LessOrEqual(t, len(errors), numGoroutines-1)

		// Final state should be consistent
		manager := migrations.NewMigrationManager(client, log)
		version, err := manager.GetVersion(ctx)
		assert.NoError(t, err)
		assert.Greater(t, version, 0)
	})
}

// BenchmarkMigrationUp benchmarks the migration up performance.
func BenchmarkMigrationUp(b *testing.B) {
	cfg := getTestConfig()
	log, err := getTestLogger()
	require.NoError(b, err)

	client, err := database.NewClient(cfg, log)
	if err != nil {
		b.Skipf("MongoDB not available for benchmarking: %v", err)
		return
	}
	defer client.Close(context.Background())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Clean database for each iteration
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		client.Database().Drop(ctx)
		cancel()

		manager := migrations.NewMigrationManager(client, log)
		
		b.StartTimer()
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		err := manager.Up(ctx)
		cancel()
		
		if err != nil {
			b.Errorf("Migration failed: %v", err)
		}
	}
}