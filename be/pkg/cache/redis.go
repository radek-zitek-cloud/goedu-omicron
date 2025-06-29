// Package cache provides Redis connection management and caching utilities
// for the GoEdu Control Testing Platform. It implements connection pooling,
// health checking, and common caching patterns for improved performance.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/config"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/logger"
)

// Client wraps the Redis client with additional functionality for caching,
// health checks, and common operations. It provides a centralized caching
// layer for the application with monitoring and error handling.
type Client struct {
	client *redis.Client
	config *config.CacheConfig
	logger *logger.Logger
}

// HealthStatus represents the health status of the Redis connection.
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency_ms"`
	Error     string    `json:"error,omitempty"`
	Info      map[string]string `json:"info,omitempty"`
}

// NewClient creates a new Redis cache client with the provided configuration.
// It establishes connection pooling, sets timeouts, and configures retry logic
// for reliable caching operations in financial applications.
//
// Parameters:
//   - cfg: Cache configuration containing Redis connection settings
//   - log: Logger instance for cache operation logging
//
// Returns:
//   - *Client: Configured Redis cache client
//   - error: Connection establishment error
//
// Example:
//   client, err := NewClient(&config.Cache, logger)
//   if err != nil {
//       return fmt.Errorf("failed to connect to Redis: %w", err)
//   }
//   defer client.Close()
func NewClient(cfg *config.CacheConfig, log *logger.Logger) (*Client, error) {
	// Configure Redis options with connection pooling and timeouts
	options := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.Database,
		MaxRetries:   cfg.MaxRetries,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Create Redis client
	client := redis.NewClient(options)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Info("Successfully connected to Redis",
		logger.String("host", cfg.Host),
		logger.Int("port", cfg.Port),
		logger.Int("database", cfg.Database),
		logger.Int("pool_size", cfg.PoolSize),
	)

	return &Client{
		client: client,
		config: cfg,
		logger: log,
	}, nil
}

// Set stores a value in the cache with the specified key and expiration time.
// This is the primary method for caching data with automatic serialization.
//
// Parameters:
//   - ctx: Context for the cache operation with timeout
//   - key: Cache key for the value
//   - value: Value to cache (will be JSON serialized)
//   - expiration: TTL for the cached value (0 for no expiration)
//
// Returns:
//   - error: Cache operation error
//
// Example:
//   user := &User{ID: "123", Name: "John Doe"}
//   err := client.Set(ctx, "user:123", user, 1*time.Hour)
//   if err != nil {
//       log.Error("Failed to cache user", zap.Error(err))
//   }
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Serialize value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}

	// Store in Redis
	if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
		c.logger.Error(ctx, "Failed to set cache value", err,
			logger.String("key", key),
			logger.Duration("expiration", expiration),
		)
		return fmt.Errorf("failed to set cache value for key %s: %w", key, err)
	}

	c.logger.Info("Cached value successfully",
		logger.String("key", key),
		logger.Duration("expiration", expiration),
	)

	return nil
}

// Get retrieves a value from the cache and deserializes it into the provided destination.
// Returns redis.Nil error if the key doesn't exist.
//
// Parameters:
//   - ctx: Context for the cache operation with timeout
//   - key: Cache key to retrieve
//   - dest: Destination to unmarshal the cached value into
//
// Returns:
//   - error: Cache operation or deserialization error (redis.Nil if key not found)
//
// Example:
//   var user User
//   err := client.Get(ctx, "user:123", &user)
//   if err == redis.Nil {
//       // Key not found, load from database
//   } else if err != nil {
//       return fmt.Errorf("cache error: %w", err)
//   }
func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
	// Get from Redis
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			c.logger.Info("Cache miss",
				logger.String("key", key),
			)
		} else {
			c.logger.Error(ctx, "Failed to get cache value", err,
				logger.String("key", key),
			)
		}
		return err
	}

	// Deserialize JSON
	if err := json.Unmarshal([]byte(data), dest); err != nil {
		c.logger.Error(ctx, "Failed to unmarshal cached value", err,
			logger.String("key", key),
		)
		return fmt.Errorf("failed to unmarshal cached value for key %s: %w", key, err)
	}

	c.logger.Info("Cache hit",
		logger.String("key", key),
	)

	return nil
}

// Delete removes a value from the cache.
//
// Parameters:
//   - ctx: Context for the cache operation with timeout
//   - keys: Cache keys to delete
//
// Returns:
//   - error: Cache operation error
//
// Example:
//   err := client.Delete(ctx, "user:123", "session:abc")
//   if err != nil {
//       log.Error("Failed to delete cache keys", zap.Error(err))
//   }
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	deleted, err := c.client.Del(ctx, keys...).Result()
	if err != nil {
		c.logger.Error(ctx, "Failed to delete cache keys", err,
			logger.Strings("keys", keys),
		)
		return fmt.Errorf("failed to delete cache keys: %w", err)
	}

	c.logger.Info("Deleted cache keys",
		logger.Strings("keys", keys),
		logger.Int64("deleted_count", deleted),
	)

	return nil
}

// Exists checks if one or more keys exist in the cache.
//
// Parameters:
//   - ctx: Context for the cache operation with timeout
//   - keys: Cache keys to check
//
// Returns:
//   - int64: Number of keys that exist
//   - error: Cache operation error
//
// Example:
//   count, err := client.Exists(ctx, "user:123", "user:456")
//   if err != nil {
//       return err
//   }
//   log.Info("Found keys", zap.Int64("count", count))
func (c *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	count, err := c.client.Exists(ctx, keys...).Result()
	if err != nil {
		c.logger.Error(ctx, "Failed to check key existence", err,
			logger.Strings("keys", keys),
		)
		return 0, fmt.Errorf("failed to check key existence: %w", err)
	}

	return count, nil
}

// Expire sets a timeout on a key.
//
// Parameters:
//   - ctx: Context for the cache operation with timeout
//   - key: Key to set expiration on
//   - expiration: TTL for the key
//
// Returns:
//   - error: Cache operation error
//
// Example:
//   err := client.Expire(ctx, "user:123", 30*time.Minute)
//   if err != nil {
//       log.Error("Failed to set key expiration", zap.Error(err))
//   }
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	success, err := c.client.Expire(ctx, key, expiration).Result()
	if err != nil {
		c.logger.Error(ctx, "Failed to set key expiration", err,
			logger.String("key", key),
			logger.Duration("expiration", expiration),
		)
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}

	if !success {
		return fmt.Errorf("key %s does not exist", key)
	}

	return nil
}

// HealthCheck performs a comprehensive health check of the Redis connection.
// This includes connectivity testing, latency measurement, and basic operations.
//
// Parameters:
//   - ctx: Context for the health check operation with timeout
//
// Returns:
//   - *HealthStatus: Detailed health status information
//
// Example:
//   health := client.HealthCheck(ctx)
//   if health.Status != "healthy" {
//       log.Error("Redis health check failed", zap.String("error", health.Error))
//   }
func (c *Client) HealthCheck(ctx context.Context) *HealthStatus {
	start := time.Now()
	status := &HealthStatus{
		Timestamp: start,
		Info:      make(map[string]string),
	}

	// Test basic connectivity with ping
	pong, err := c.client.Ping(ctx).Result()
	if err != nil {
		status.Status = "unhealthy"
		status.Error = fmt.Sprintf("ping failed: %v", err)
		status.Latency = time.Since(start).Milliseconds()
		return status
	}

	status.Info["ping_response"] = pong

	// Test basic operations
	testKey := fmt.Sprintf("health_check_%d", time.Now().UnixNano())
	testValue := "test"

	// Test SET operation
	if err := c.client.Set(ctx, testKey, testValue, 10*time.Second).Err(); err != nil {
		status.Status = "unhealthy"
		status.Error = fmt.Sprintf("set operation failed: %v", err)
		status.Latency = time.Since(start).Milliseconds()
		return status
	}

	// Test GET operation
	value, err := c.client.Get(ctx, testKey).Result()
	if err != nil {
		status.Status = "unhealthy"
		status.Error = fmt.Sprintf("get operation failed: %v", err)
		status.Latency = time.Since(start).Milliseconds()
		return status
	}

	if value != testValue {
		status.Status = "unhealthy"
		status.Error = "get operation returned incorrect value"
		status.Latency = time.Since(start).Milliseconds()
		return status
	}

	// Clean up test key
	c.client.Del(ctx, testKey)

	// Get server info
	info, err := c.client.Info(ctx, "server").Result()
	if err == nil {
		status.Info["server_info"] = info
	}

	// All checks passed
	status.Status = "healthy"
	status.Latency = time.Since(start).Milliseconds()

	return status
}

// FlushDB removes all keys from the current database.
// This should only be used in development or testing environments.
//
// Parameters:
//   - ctx: Context for the cache operation with timeout
//
// Returns:
//   - error: Cache operation error
//
// Example:
//   if config.Environment == "development" {
//       err := client.FlushDB(ctx)
//       if err != nil {
//           log.Error("Failed to flush cache", zap.Error(err))
//       }
//   }
func (c *Client) FlushDB(ctx context.Context) error {
	if err := c.client.FlushDB(ctx).Err(); err != nil {
		c.logger.Error(ctx, "Failed to flush database", err)
		return fmt.Errorf("failed to flush database: %w", err)
	}

	c.logger.Warn("Flushed Redis database")
	return nil
}

// Close gracefully closes the Redis connection.
// This should be called during application shutdown.
//
// Returns:
//   - error: Connection closure error
//
// Example:
//   if err := client.Close(); err != nil {
//       log.Error("Failed to close Redis connection", zap.Error(err))
//   }
func (c *Client) Close() error {
	c.logger.Info("Closing Redis connection...")
	
	if err := c.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	c.logger.Info("Redis connection closed successfully")
	return nil
}

// Stats returns connection statistics for monitoring and debugging.
//
// Returns:
//   - map[string]interface{}: Connection statistics
//
// Example:
//   stats := client.Stats()
//   poolSize := stats["pool_size"]
//   database := stats["database"]
func (c *Client) Stats() map[string]interface{} {
	return map[string]interface{}{
		"host":         c.config.Host,
		"port":         c.config.Port,
		"database":     c.config.Database,
		"pool_size":    c.config.PoolSize,
		"max_retries":  c.config.MaxRetries,
		"dial_timeout": c.config.DialTimeout.String(),
		"read_timeout": c.config.ReadTimeout.String(),
		"write_timeout": c.config.WriteTimeout.String(),
		"idle_timeout": c.config.IdleTimeout.String(),
	}
}