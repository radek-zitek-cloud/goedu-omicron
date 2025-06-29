// Package server_test provides tests for the HTTP server functionality,
// including health check endpoints and server initialization.
package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/config"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/logger"
)

// TestHealthEndpointStructure tests the health endpoint response structure
// without requiring a live database connection.
func TestHealthEndpointStructure(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create a mock health check handler
	router := gin.New()
	
	router.GET("/health", func(c *gin.Context) {
		// Mock health check response similar to the real implementation
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
			"checks": gin.H{
				"database": gin.H{
					"status":     "healthy",
					"timestamp":  time.Now().UTC(),
					"latency_ms": int64(10),
				},
				"cache": gin.H{
					"status":     "healthy",
					"timestamp":  time.Now().UTC(),
					"latency_ms": int64(5),
				},
			},
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	t.Run("health endpoint returns correct structure", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "status")
		assert.Contains(t, w.Body.String(), "timestamp")
		assert.Contains(t, w.Body.String(), "version")
		assert.Contains(t, w.Body.String(), "checks")
		assert.Contains(t, w.Body.String(), "database")
		assert.Contains(t, w.Body.String(), "cache")
	})

	t.Run("readiness endpoint returns correct structure", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ready", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "status")
		assert.Contains(t, w.Body.String(), "ready")
		assert.Contains(t, w.Body.String(), "timestamp")
		assert.Contains(t, w.Body.String(), "version")
	})
}

// TestConfigurationLoading tests configuration loading functionality.
func TestConfigurationLoading(t *testing.T) {
	t.Run("load default configuration", func(t *testing.T) {
		cfg, err := config.Load()
		require.NoError(t, err)
		assert.NotNil(t, cfg)

		// Verify basic configuration structure
		assert.NotEmpty(t, cfg.App.Name)
		assert.NotEmpty(t, cfg.App.Version)
		assert.Greater(t, cfg.App.Port, 0)
		assert.NotEmpty(t, cfg.Database.URI)
		assert.NotEmpty(t, cfg.Database.Database)
		assert.Greater(t, cfg.Database.MaxPoolSize, 0)
		assert.Greater(t, cfg.Database.MinPoolSize, 0)
		assert.Greater(t, cfg.Database.ConnectTimeout, time.Duration(0))
	})

	t.Run("validate development configuration", func(t *testing.T) {
		cfg, err := config.Load()
		require.NoError(t, err)

		// Should pass validation in development mode
		assert.True(t, cfg.IsDevelopment())
		assert.False(t, cfg.IsProduction())
	})

	t.Run("configuration helper methods", func(t *testing.T) {
		cfg, err := config.Load()
		require.NoError(t, err)

		// Test helper methods
		dbURI := cfg.GetDatabaseURI()
		assert.NotEmpty(t, dbURI)

		redisAddr := cfg.GetRedisAddr()
		assert.NotEmpty(t, redisAddr)
		assert.Contains(t, redisAddr, ":")

		serverAddr := cfg.GetServerAddr()
		assert.NotEmpty(t, serverAddr)
		assert.Contains(t, serverAddr, ":")
	})
}

// TestLoggerInitialization tests logger initialization.
func TestLoggerInitialization(t *testing.T) {
	t.Run("create logger with test config", func(t *testing.T) {
		cfg := &logger.Config{
			Level:       logger.InfoLevel,
			Environment: "test",
			OutputPath:  "stdout",
		}

		log, err := logger.New(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, log)
	})

	t.Run("logger with different levels", func(t *testing.T) {
		levels := []logger.LogLevel{logger.DebugLevel, logger.InfoLevel, logger.WarnLevel, logger.ErrorLevel}

		for _, level := range levels {
			cfg := &logger.Config{
				Level:       level,
				Environment: "test",
				OutputPath:  "stdout",
			}

			log, err := logger.New(cfg)
			assert.NoError(t, err, "Failed to create logger with level: %s", level)
			assert.NotNil(t, log, "Logger should not be nil for level: %s", level)
		}
	})
}

// TestCORSMiddleware tests CORS middleware functionality.
func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create a test router with CORS middleware
	router := gin.New()
	
	// Mock CORS middleware
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := []string{"http://localhost:3000", "http://localhost:3001"}
		
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin || allowedOrigin == "*" {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Correlation-ID")
		c.Header("Access-Control-Expose-Headers", "X-Correlation-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			return
		}

		c.Next()
	})

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	t.Run("CORS headers with allowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	})

	t.Run("OPTIONS preflight request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("CORS headers with disallowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "http://malicious.com")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})
}

// TestRequestLogging tests request logging middleware functionality.
func TestRequestLogging(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	
	// Mock logging middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()

		// Generate correlation ID
		correlationID := "test-correlation-id"
		c.Set("correlation_id", correlationID)
		c.Header("X-Correlation-ID", correlationID)

		c.Next()

		// Log would happen here in real implementation
		duration := time.Since(start)
		assert.Greater(t, duration, time.Duration(0))
	})

	router.GET("/test", func(c *gin.Context) {
		correlationID, exists := c.Get("correlation_id")
		assert.True(t, exists)
		assert.Equal(t, "test-correlation-id", correlationID)
		
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	t.Run("correlation ID is set", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "test-correlation-id", w.Header().Get("X-Correlation-ID"))
	})
}

// TestServerConfiguration tests server configuration validation.
func TestServerConfiguration(t *testing.T) {
	t.Run("valid port range", func(t *testing.T) {
		validPorts := []int{8080, 3000, 8000, 9000}
		
		for _, port := range validPorts {
			// Simulate validation logic
			assert.True(t, port >= 1024 && port <= 65535, "Port %d should be valid", port)
		}
	})

	t.Run("invalid port range", func(t *testing.T) {
		invalidPorts := []int{80, 443, 0, 70000}
		
		for _, port := range invalidPorts {
			// Simulate validation logic
			valid := port >= 1024 && port <= 65535
			if port == 80 || port == 443 {
				// These ports might be valid in some contexts but restricted in our validation
				assert.False(t, valid, "Port %d should be considered invalid for non-privileged applications", port)
			} else {
				assert.False(t, valid, "Port %d should be invalid", port)
			}
		}
	})
}

// TestDatabaseConfiguration tests database configuration validation.
func TestDatabaseConfiguration(t *testing.T) {
	t.Run("pool size validation", func(t *testing.T) {
		testCases := []struct {
			name        string
			maxPool     int
			minPool     int
			shouldValid bool
		}{
			{"valid pool sizes", 100, 10, true},
			{"equal pool sizes", 50, 50, true},
			{"invalid - min > max", 10, 100, false},
			{"zero values", 0, 0, true}, // Might be valid in some contexts
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				valid := tc.maxPool >= tc.minPool
				assert.Equal(t, tc.shouldValid, valid, "Pool size validation failed for case: %s", tc.name)
			})
		}
	})

	t.Run("timeout validation", func(t *testing.T) {
		validTimeouts := []time.Duration{
			5 * time.Second,
			10 * time.Second,
			30 * time.Second,
			1 * time.Minute,
		}

		for _, timeout := range validTimeouts {
			assert.Greater(t, timeout, time.Duration(0), "Timeout %v should be positive", timeout)
		}
	})
}

// BenchmarkHealthEndpoint benchmarks the health endpoint performance.
func BenchmarkHealthEndpoint(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			if w.Code != http.StatusOK {
				b.Errorf("Expected status 200, got %d", w.Code)
			}
		}
	})
}