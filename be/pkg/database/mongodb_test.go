// Package database_test provides comprehensive tests for MongoDB connection management,
// health checking, and connection pooling functionality in the GoEdu Control Testing Platform.
package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/config"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/database"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/logger"
)

// TestConfig provides test configuration for database testing.
func getTestConfig() *config.DatabaseConfig {
	return &config.DatabaseConfig{
		URI:                 "mongodb://localhost:27017",
		Database:            "goedu_test",
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

// TestNewClient tests the MongoDB client creation and connection.
func TestNewClient(t *testing.T) {
	t.Run("successful connection", func(t *testing.T) {
		cfg := getTestConfig()
		log, err := getTestLogger()
		require.NoError(t, err)

		client, err := database.NewClient(cfg, log)
		if err != nil {
			t.Skipf("MongoDB not available for testing: %v", err)
			return
		}
		defer client.Close(context.Background())

		assert.NotNil(t, client)
		
		// Test database access
		db := client.Database()
		assert.NotNil(t, db)
		assert.Equal(t, "goedu_test", db.Name())
	})

	t.Run("invalid connection string", func(t *testing.T) {
		cfg := getTestConfig()
		cfg.URI = "invalid://connection"
		
		log, err := getTestLogger()
		require.NoError(t, err)

		_, err = database.NewClient(cfg, log)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create MongoDB client")
	})

	t.Run("connection timeout", func(t *testing.T) {
		cfg := getTestConfig()
		cfg.URI = "mongodb://nonexistent:27017"
		cfg.ConnectTimeout = 1 * time.Second
		cfg.ServerSelectTimeout = 1 * time.Second
		
		log, err := getTestLogger()
		require.NoError(t, err)

		_, err = database.NewClient(cfg, log)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to ping MongoDB")
	})
}

// TestHealthCheck tests the database health check functionality.
func TestHealthCheck(t *testing.T) {
	cfg := getTestConfig()
	log, err := getTestLogger()
	require.NoError(t, err)

	client, err := database.NewClient(cfg, log)
	if err != nil {
		t.Skipf("MongoDB not available for testing: %v", err)
		return
	}
	defer client.Close(context.Background())

	t.Run("healthy connection", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		health := client.HealthCheck(ctx)
		
		assert.NotNil(t, health)
		assert.Equal(t, "healthy", health.Status)
		assert.Empty(t, health.Error)
		assert.Greater(t, health.Latency, int64(0))
		assert.WithinDuration(t, time.Now(), health.Timestamp, 1*time.Second)
	})

	t.Run("health check with context timeout", func(t *testing.T) {
		// Create very short timeout to test timeout handling
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		health := client.HealthCheck(ctx)
		
		assert.NotNil(t, health)
		// Health check should handle timeout gracefully
		// The actual result depends on timing, but it should not panic
	})
}

// TestCreateIndexes tests the database index creation functionality.
func TestCreateIndexes(t *testing.T) {
	cfg := getTestConfig()
	log, err := getTestLogger()
	require.NoError(t, err)

	client, err := database.NewClient(cfg, log)
	if err != nil {
		t.Skipf("MongoDB not available for testing: %v", err)
		return
	}
	defer client.Close(context.Background())

	t.Run("create indexes successfully", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := client.CreateIndexes(ctx)
		assert.NoError(t, err)

		// Verify that indexes were created by checking a few key collections
		collections := []string{"users", "organizations", "controls"}
		
		for _, collectionName := range collections {
			collection := client.Collection(collectionName)
			
			cursor, err := collection.Indexes().List(ctx)
			require.NoError(t, err)
			
			var indexes []map[string]interface{}
			err = cursor.All(ctx, &indexes)
			require.NoError(t, err)
			
			// Should have at least the _id index plus our custom indexes
			assert.Greater(t, len(indexes), 1, "Collection %s should have custom indexes", collectionName)
		}
	})
}

// TestConnectionPooling tests the MongoDB connection pooling functionality.
func TestConnectionPooling(t *testing.T) {
	cfg := getTestConfig()
	cfg.MaxPoolSize = 5
	cfg.MinPoolSize = 2
	
	log, err := getTestLogger()
	require.NoError(t, err)

	client, err := database.NewClient(cfg, log)
	if err != nil {
		t.Skipf("MongoDB not available for testing: %v", err)
		return
	}
	defer client.Close(context.Background())

	t.Run("connection pool configuration", func(t *testing.T) {
		stats := client.Stats()
		
		assert.Equal(t, cfg.MaxPoolSize, stats["max_pool_size"])
		assert.Equal(t, cfg.MinPoolSize, stats["min_pool_size"])
		assert.Equal(t, cfg.Database, stats["database_name"])
		assert.Equal(t, cfg.ConnectTimeout.String(), stats["connect_timeout"])
	})

	t.Run("concurrent connections", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Test concurrent database operations
		const numGoroutines = 10
		done := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				collection := client.Collection("test_concurrent")
				_, err := collection.InsertOne(ctx, map[string]interface{}{
					"test":      true,
					"timestamp": time.Now(),
				})
				done <- err
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			err := <-done
			assert.NoError(t, err, "Concurrent operation %d should succeed", i)
		}

		// Clean up test data
		client.Collection("test_concurrent").Drop(ctx)
	})
}

// TestCollectionAccess tests direct collection access functionality.
func TestCollectionAccess(t *testing.T) {
	cfg := getTestConfig()
	log, err := getTestLogger()
	require.NoError(t, err)

	client, err := database.NewClient(cfg, log)
	if err != nil {
		t.Skipf("MongoDB not available for testing: %v", err)
		return
	}
	defer client.Close(context.Background())

	t.Run("collection access", func(t *testing.T) {
		collection := client.Collection("test_collection")
		assert.NotNil(t, collection)
		assert.Equal(t, "test_collection", collection.Name())
	})

	t.Run("database access", func(t *testing.T) {
		db := client.Database()
		assert.NotNil(t, db)
		assert.Equal(t, cfg.Database, db.Name())
	})
}

// TestGracefulShutdown tests the graceful shutdown functionality.
func TestGracefulShutdown(t *testing.T) {
	cfg := getTestConfig()
	log, err := getTestLogger()
	require.NoError(t, err)

	client, err := database.NewClient(cfg, log)
	if err != nil {
		t.Skipf("MongoDB not available for testing: %v", err)
		return
	}

	t.Run("graceful close", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := client.Close(ctx)
		assert.NoError(t, err)

		// Verify connection is closed by attempting a health check
		health := client.HealthCheck(ctx)
		assert.Equal(t, "unhealthy", health.Status)
		assert.NotEmpty(t, health.Error)
	})

	t.Run("double close", func(t *testing.T) {
		client2, err := database.NewClient(cfg, log)
		if err != nil {
			t.Skipf("MongoDB not available for testing: %v", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// First close should succeed
		err = client2.Close(ctx)
		assert.NoError(t, err)

		// Second close should not panic or error
		err = client2.Close(ctx)
		assert.NoError(t, err)
	})
}

// BenchmarkHealthCheck benchmarks the health check performance.
func BenchmarkHealthCheck(b *testing.B) {
	cfg := getTestConfig()
	log, err := getTestLogger()
	require.NoError(b, err)

	client, err := database.NewClient(cfg, log)
	if err != nil {
		b.Skipf("MongoDB not available for benchmarking: %v", err)
		return
	}
	defer client.Close(context.Background())

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			health := client.HealthCheck(ctx)
			if health.Status != "healthy" {
				b.Errorf("Health check failed: %s", health.Error)
			}
		}
	})
}

// BenchmarkConnectionPooling benchmarks connection pool performance.
func BenchmarkConnectionPooling(b *testing.B) {
	cfg := getTestConfig()
	cfg.MaxPoolSize = 20
	
	log, err := getTestLogger()
	require.NoError(b, err)

	client, err := database.NewClient(cfg, log)
	if err != nil {
		b.Skipf("MongoDB not available for benchmarking: %v", err)
		return
	}
	defer client.Close(context.Background())

	ctx := context.Background()
	collection := client.Collection("benchmark_test")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, err := collection.InsertOne(ctx, map[string]interface{}{
				"benchmark": true,
				"iteration": i,
				"timestamp": time.Now(),
			})
			if err != nil {
				b.Errorf("Insert failed: %v", err)
			}
			i++
		}
	})

	// Clean up
	collection.Drop(ctx)
}