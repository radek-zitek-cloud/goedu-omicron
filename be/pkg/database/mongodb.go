// Package database provides MongoDB connection management and health checking
// for the GoEdu Control Testing Platform. It implements connection pooling,
// timeout handling, and monitoring capabilities required for financial applications.
package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/config"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/logger"
)

// Client wraps the MongoDB client with additional functionality for health checks
// and connection management. It provides a centralized point for database operations
// across all services in the application.
type Client struct {
	client   *mongo.Client
	database *mongo.Database
	config   *config.DatabaseConfig
	logger   *logger.Logger
}

// HealthStatus represents the health status of the database connection.
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency_ms"`
	Error     string    `json:"error,omitempty"`
}

// NewClient creates a new database client with the provided configuration.
// It establishes connection pooling, sets timeouts, and configures monitoring
// as required for high-availability financial applications.
//
// Parameters:
//   - cfg: Database configuration containing connection settings
//   - log: Logger instance for database operation logging
//
// Returns:
//   - *Client: Configured database client
//   - error: Connection establishment error
//
// Example:
//   client, err := NewClient(&config.Database, logger)
//   if err != nil {
//       return fmt.Errorf("failed to connect to database: %w", err)
//   }
//   defer client.Close(context.Background())
func NewClient(cfg *config.DatabaseConfig, log *logger.Logger) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	// Configure client options with connection pooling and timeouts
	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(uint64(cfg.MaxPoolSize)).
		SetMinPoolSize(uint64(cfg.MinPoolSize)).
		SetMaxConnIdleTime(cfg.MaxConnIdleTime).
		SetConnectTimeout(cfg.ConnectTimeout).
		SetServerSelectionTimeout(cfg.ServerSelectTimeout)

	// Create MongoDB client
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	// Test the connection
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Get database handle
	database := client.Database(cfg.Database)

	log.Info("Successfully connected to MongoDB",
		logger.String("database", cfg.Database),
		logger.Int("max_pool_size", cfg.MaxPoolSize),
		logger.Int("min_pool_size", cfg.MinPoolSize),
	)

	return &Client{
		client:   client,
		database: database,
		config:   cfg,
		logger:   log,
	}, nil
}

// Database returns the MongoDB database handle for performing operations.
// All repository implementations should use this method to get database access.
//
// Returns:
//   - *mongo.Database: MongoDB database handle
//
// Example:
//   db := client.Database()
//   collection := db.Collection("controls")
//   result, err := collection.InsertOne(ctx, document)
func (c *Client) Database() *mongo.Database {
	return c.database
}

// Collection provides direct access to a specific MongoDB collection.
// This is a convenience method for common collection access patterns.
//
// Parameters:
//   - name: Name of the collection to access
//
// Returns:
//   - *mongo.Collection: MongoDB collection handle
//
// Example:
//   controls := client.Collection("controls")
//   users := client.Collection("users")
func (c *Client) Collection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

// HealthCheck performs a comprehensive health check of the database connection.
// This includes connectivity testing, latency measurement, and basic query execution.
// Financial applications require detailed health monitoring for regulatory compliance.
//
// Parameters:
//   - ctx: Context for the health check operation with timeout
//
// Returns:
//   - *HealthStatus: Detailed health status information
//
// Example:
//   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//   defer cancel()
//   
//   health := client.HealthCheck(ctx)
//   if health.Status != "healthy" {
//       log.Error("Database health check failed", zap.String("error", health.Error))
//   }
func (c *Client) HealthCheck(ctx context.Context) *HealthStatus {
	start := time.Now()
	status := &HealthStatus{
		Timestamp: start,
	}

	// Test basic connectivity with ping
	if err := c.client.Ping(ctx, readpref.Primary()); err != nil {
		status.Status = "unhealthy"
		status.Error = fmt.Sprintf("ping failed: %v", err)
		status.Latency = time.Since(start).Milliseconds()
		return status
	}

	// Test basic query operation
	result := c.database.RunCommand(ctx, bson.D{{"ping", 1}})
	if result.Err() != nil {
		status.Status = "unhealthy"
		status.Error = fmt.Sprintf("command failed: %v", result.Err())
		status.Latency = time.Since(start).Milliseconds()
		return status
	}

	// All checks passed
	status.Status = "healthy"
	status.Latency = time.Since(start).Milliseconds()

	return status
}

// CreateIndexes creates necessary database indexes for optimal performance.
// Financial applications require well-designed indexes for query performance
// and to support audit trail requirements.
//
// Parameters:
//   - ctx: Context for index creation operations
//
// Returns:
//   - error: Index creation error
//
// Example:
//   if err := client.CreateIndexes(context.Background()); err != nil {
//       log.Error("Failed to create database indexes", zap.Error(err))
//   }
func (c *Client) CreateIndexes(ctx context.Context) error {
	c.logger.Info("Creating database indexes...")

	// Define index models for different collections
	indexModels := map[string][]mongo.IndexModel{
		"users": {
			{
				Keys: bson.D{{"email", 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.D{{"organization_id", 1}, {"role", 1}},
			},
			{
				Keys: bson.D{{"created_at", -1}},
			},
		},
		"organizations": {
			{
				Keys: bson.D{{"slug", 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.D{{"created_at", -1}},
			},
		},
		"controls": {
			{
				Keys: bson.D{{"organization_id", 1}, {"control_id", 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.D{{"organization_id", 1}, {"status", 1}},
			},
			{
				Keys: bson.D{{"framework", 1}, {"category", 1}},
			},
			{
				Keys: bson.D{{"created_at", -1}},
			},
		},
		"testing_cycles": {
			{
				Keys: bson.D{{"organization_id", 1}, {"cycle_id", 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.D{{"organization_id", 1}, {"status", 1}},
			},
			{
				Keys: bson.D{{"start_date", 1}, {"end_date", 1}},
			},
		},
		"evidence_requests": {
			{
				Keys: bson.D{{"organization_id", 1}, {"control_id", 1}},
			},
			{
				Keys: bson.D{{"assignee_id", 1}, {"status", 1}},
			},
			{
				Keys: bson.D{{"due_date", 1}},
			},
			{
				Keys: bson.D{{"created_at", -1}},
			},
		},
		"audit_logs": {
			{
				Keys: bson.D{{"user_id", 1}, {"timestamp", -1}},
			},
			{
				Keys: bson.D{{"organization_id", 1}, {"timestamp", -1}},
			},
			{
				Keys: bson.D{{"action", 1}, {"timestamp", -1}},
			},
			{
				Keys: bson.D{{"resource_id", 1}, {"timestamp", -1}},
			},
			{
				Keys: bson.D{{"correlation_id", 1}},
			},
		},
	}

	// Create indexes for each collection
	for collectionName, indexes := range indexModels {
		collection := c.database.Collection(collectionName)
		
		if len(indexes) > 0 {
			_, err := collection.Indexes().CreateMany(ctx, indexes)
			if err != nil {
				return fmt.Errorf("failed to create indexes for collection %s: %w", collectionName, err)
			}
			
			c.logger.Info("Created indexes for collection",
				logger.String("collection", collectionName),
				logger.Int("index_count", len(indexes)),
			)
		}
	}

	c.logger.Info("Successfully created all database indexes")
	return nil
}

// Close gracefully closes the database connection.
// This should be called during application shutdown to ensure clean resource cleanup.
//
// Parameters:
//   - ctx: Context for connection closure with timeout
//
// Returns:
//   - error: Connection closure error
//
// Example:
//   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//   defer cancel()
//   
//   if err := client.Close(ctx); err != nil {
//       log.Error("Failed to close database connection", zap.Error(err))
//   }
func (c *Client) Close(ctx context.Context) error {
	if c.client == nil {
		return nil
	}

	c.logger.Info("Closing database connection...")
	
	if err := c.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	c.logger.Info("Database connection closed successfully")
	return nil
}

// Stats returns connection statistics for monitoring and debugging.
// This information is useful for performance monitoring and capacity planning.
//
// Returns:
//   - map[string]interface{}: Connection statistics
//
// Example:
//   stats := client.Stats()
//   activeConnections := stats["active_connections"]
//   poolSize := stats["pool_size"]
func (c *Client) Stats() map[string]interface{} {
	return map[string]interface{}{
		"database_name":     c.config.Database,
		"max_pool_size":     c.config.MaxPoolSize,
		"min_pool_size":     c.config.MinPoolSize,
		"connect_timeout":   c.config.ConnectTimeout.String(),
		"max_conn_idle_time": c.config.MaxConnIdleTime.String(),
	}
}