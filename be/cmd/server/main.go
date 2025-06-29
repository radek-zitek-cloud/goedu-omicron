// Package main is the entry point for the GoEdu Control Testing Platform backend server.
// It initializes all services, establishes database connections, and starts the HTTP server
// with proper graceful shutdown handling for production deployment.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/config"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/cache"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/database"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/logger"
)

// Application holds all application dependencies and services.
// This structure provides dependency injection and service management.
type Application struct {
	config   *config.Config
	logger   *logger.Logger
	database *database.Client
	cache    *cache.Client
	server   *http.Server
}

// main is the application entry point.
// It initializes the application, starts services, and handles graceful shutdown.
//
// The application follows these initialization steps:
// 1. Load configuration from environment and files
// 2. Initialize structured logging
// 3. Connect to MongoDB database
// 4. Connect to Redis cache
// 5. Set up HTTP routes and middleware
// 6. Start the HTTP server
// 7. Wait for shutdown signals
// 8. Perform graceful shutdown
func main() {
	// Create application context for initialization
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize application
	app, err := NewApplication(ctx)
	if err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// Start the application
	if err := app.Start(ctx); err != nil {
		app.logger.Error(ctx, "Failed to start application", err)
		os.Exit(1)
	}

	// Wait for shutdown signal
	app.WaitForShutdown()

	// Perform graceful shutdown
	if err := app.Shutdown(ctx); err != nil {
		app.logger.Error(ctx, "Error during shutdown", err)
		os.Exit(1)
	}

	app.logger.Info("Application shutdown complete")
}

// NewApplication creates and initializes a new application instance.
// It loads configuration, establishes database connections, and sets up services.
//
// Parameters:
//   - ctx: Context for initialization operations with timeout
//
// Returns:
//   - *Application: Initialized application instance
//   - error: Initialization error if any step fails
//
// Example:
//   app, err := NewApplication(context.Background())
//   if err != nil {
//       log.Fatal("Application initialization failed:", err)
//   }
func NewApplication(ctx context.Context) (*Application, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger
	log, err := logger.New(&cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	log.Info("Application initialization started",
		logger.String("name", cfg.App.Name),
		logger.String("version", cfg.App.Version),
		logger.String("environment", cfg.App.Environment),
	)

	// Connect to MongoDB
	log.Info("Connecting to MongoDB...")
	dbClient, err := database.NewClient(&cfg.Database, log)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create database indexes
	log.Info("Creating database indexes...")
	if err := dbClient.CreateIndexes(ctx); err != nil {
		log.Error(ctx, "Failed to create database indexes", err)
		// Don't fail startup, just log the error
	}

	// Connect to Redis
	log.Info("Connecting to Redis...")
	cacheClient, err := cache.NewClient(&cfg.Cache, log)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cache: %w", err)
	}

	// Create application instance
	app := &Application{
		config:   cfg,
		logger:   log,
		database: dbClient,
		cache:    cacheClient,
	}

	// Set up HTTP server
	if err := app.setupServer(); err != nil {
		return nil, fmt.Errorf("failed to setup HTTP server: %w", err)
	}

	log.Info("Application initialized successfully")
	return app, nil
}

// setupServer configures the HTTP server with routes, middleware, and settings.
// It sets up the Gin router with CORS, logging, recovery, and health check endpoints.
//
// Returns:
//   - error: Server setup error
func (app *Application) setupServer() error {
	// Set Gin mode based on environment
	if app.config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(app.loggingMiddleware())
	router.Use(app.corsMiddleware())

	// Health check endpoint
	router.GET("/health", app.healthCheckHandler)
	router.GET("/ready", app.readinessHandler)

	// API version group (placeholder for future endpoints)
	v1 := router.Group("/api/v1")
	_ = v1 // Placeholder for routes to be added
	{
		// Authentication routes would go here
		// v1.POST("/auth/login", app.loginHandler)
		// v1.POST("/auth/logout", app.logoutHandler)

		// GraphQL endpoint would go here
		// v1.POST("/graphql", app.graphqlHandler)
		// v1.GET("/graphql", app.graphqlPlaygroundHandler)

		// REST API endpoints would go here
		// controls := v1.Group("/controls")
		// controls.Use(app.authMiddleware())
		// controls.GET("", app.listControlsHandler)
		// controls.POST("", app.createControlHandler)
	}

	// Create HTTP server
	app.server = &http.Server{
		Addr:         app.config.GetServerAddr(),
		Handler:      router,
		ReadTimeout:  app.config.App.Timeout,
		WriteTimeout: app.config.App.Timeout,
		IdleTimeout:  2 * app.config.App.Timeout,
	}

	return nil
}

// Start begins serving HTTP requests on the configured port.
// It starts the server in a goroutine to allow for graceful shutdown handling.
//
// Parameters:
//   - ctx: Context for startup operations
//
// Returns:
//   - error: Server startup error
func (app *Application) Start(ctx context.Context) error {
	app.logger.Info("Starting HTTP server",
		logger.String("address", app.server.Addr),
		logger.String("environment", app.config.App.Environment),
	)

	// Start server in a goroutine
	go func() {
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error(ctx, "HTTP server error", err)
		}
	}()

	app.logger.Info("HTTP server started successfully",
		logger.String("address", app.server.Addr),
	)

	return nil
}

// WaitForShutdown waits for termination signals and begins graceful shutdown.
// It listens for SIGINT and SIGTERM signals commonly used in containerized environments.
func (app *Application) WaitForShutdown() {
	// Create channel to receive OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	sig := <-quit
	app.logger.Info("Received shutdown signal",
		logger.String("signal", sig.String()),
	)
}

// Shutdown performs graceful shutdown of all application services.
// It shuts down the HTTP server, closes database connections, and cleans up resources.
//
// Parameters:
//   - ctx: Context for shutdown operations with timeout
//
// Returns:
//   - error: Shutdown error if any operation fails
func (app *Application) Shutdown(ctx context.Context) error {
	app.logger.Info("Starting graceful shutdown...")

	// Shutdown HTTP server
	app.logger.Info("Shutting down HTTP server...")
	if err := app.server.Shutdown(ctx); err != nil {
		app.logger.Error(ctx, "HTTP server shutdown error", err)
		return fmt.Errorf("HTTP server shutdown failed: %w", err)
	}

	// Close cache connection
	app.logger.Info("Closing cache connection...")
	if err := app.cache.Close(); err != nil {
		app.logger.Error(ctx, "Cache connection close error", err)
		// Don't return error, continue with other cleanup
	}

	// Close database connection
	app.logger.Info("Closing database connection...")
	if err := app.database.Close(ctx); err != nil {
		app.logger.Error(ctx, "Database connection close error", err)
		return fmt.Errorf("database connection close failed: %w", err)
	}

	// Sync logger (flush any pending logs)
	if err := app.logger.Sync(); err != nil {
		// Ignore sync errors during shutdown
	}

	return nil
}

// HTTP Handlers

// healthCheckHandler provides a basic health check endpoint.
// It returns the application status and basic system information.
func (app *Application) healthCheckHandler(c *gin.Context) {
	ctx := c.Request.Context()

	// Check database health
	dbHealth := app.database.HealthCheck(ctx)

	// Check cache health
	cacheHealth := app.cache.HealthCheck(ctx)

	// Determine overall health
	status := "healthy"
	if dbHealth.Status != "healthy" || cacheHealth.Status != "healthy" {
		status = "unhealthy"
		c.Status(http.StatusServiceUnavailable)
	} else {
		c.Status(http.StatusOK)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    status,
		"timestamp": time.Now().UTC(),
		"version":   app.config.App.Version,
		"checks": gin.H{
			"database": dbHealth,
			"cache":    cacheHealth,
		},
	})
}

// readinessHandler provides a readiness check for Kubernetes deployments.
// It indicates when the application is ready to receive traffic.
func (app *Application) readinessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now().UTC(),
		"version":   app.config.App.Version,
	})
}

// Middleware

// loggingMiddleware provides structured request logging for audit trails.
// It logs all HTTP requests with correlation IDs and performance metrics.
func (app *Application) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate correlation ID for request tracing
		correlationID := fmt.Sprintf("%d", start.UnixNano())
		c.Set("correlation_id", correlationID)

		// Add correlation ID to response headers
		c.Header("X-Correlation-ID", correlationID)

		// Process request
		c.Next()

		// Log request completion
		duration := time.Since(start)
		app.logger.Performance(c.Request.Context(), "http_request", duration,
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.String("correlation_id", correlationID),
			logger.Int("status", c.Writer.Status()),
			logger.String("client_ip", c.ClientIP()),
			logger.String("user_agent", c.Request.UserAgent()),
		)
	}
}

// corsMiddleware configures Cross-Origin Resource Sharing (CORS) settings.
// It allows frontend applications to make requests to the API server.
func (app *Application) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range app.config.App.CORS.AllowedOrigins {
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

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			return
		}

		c.Next()
	}
}