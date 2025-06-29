// Package config provides environment-based configuration management for the GoEdu application.
// It supports multiple environments (development, staging, production) with secure handling
// of sensitive data like database credentials and API keys.
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/logger"
)

// Config holds all configuration settings for the application.
// It includes database connections, external services, security settings,
// and operational parameters required for the control testing platform.
type Config struct {
	// Application settings
	App AppConfig `mapstructure:"app"`
	
	// Database configurations
	Database DatabaseConfig `mapstructure:"database"`
	
	// Cache configuration (Redis)
	Cache CacheConfig `mapstructure:"cache"`
	
	// Object storage configuration (MinIO)
	Storage StorageConfig `mapstructure:"storage"`
	
	// Authentication and security
	Auth AuthConfig `mapstructure:"auth"`
	
	// Logging configuration
	Logger logger.Config `mapstructure:"logger"`
	
	// External services
	Email    EmailConfig    `mapstructure:"email"`
	Webhook  WebhookConfig  `mapstructure:"webhook"`
	
	// Monitoring and observability
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
}

// AppConfig contains basic application settings.
type AppConfig struct {
	Name        string        `mapstructure:"name"`
	Version     string        `mapstructure:"version"`
	Environment string        `mapstructure:"environment"`
	Port        int           `mapstructure:"port"`
	Host        string        `mapstructure:"host"`
	Timeout     time.Duration `mapstructure:"timeout"`
	CORS        CORSConfig    `mapstructure:"cors"`
}

// CORSConfig defines Cross-Origin Resource Sharing settings.
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

// DatabaseConfig contains MongoDB connection settings.
// Connection pooling and timeout settings are crucial for performance
// in financial applications where reliability is paramount.
type DatabaseConfig struct {
	URI               string        `mapstructure:"uri"`
	Database          string        `mapstructure:"database"`
	MaxPoolSize       int           `mapstructure:"max_pool_size"`
	MinPoolSize       int           `mapstructure:"min_pool_size"`
	MaxConnIdleTime   time.Duration `mapstructure:"max_conn_idle_time"`
	ConnectTimeout    time.Duration `mapstructure:"connect_timeout"`
	ServerSelectTimeout time.Duration `mapstructure:"server_select_timeout"`
}

// CacheConfig contains Redis connection and caching settings.
type CacheConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	Database     int           `mapstructure:"database"`
	MaxRetries   int           `mapstructure:"max_retries"`
	PoolSize     int           `mapstructure:"pool_size"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// StorageConfig contains MinIO object storage settings for evidence files.
type StorageConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	BucketName      string `mapstructure:"bucket_name"`
	Region          string `mapstructure:"region"`
	UseSSL          bool   `mapstructure:"use_ssl"`
}

// AuthConfig contains authentication and JWT settings.
// Security is critical for financial applications handling sensitive audit data.
type AuthConfig struct {
	JWTSecret     string        `mapstructure:"jwt_secret"`
	JWTExpiration time.Duration `mapstructure:"jwt_expiration"`
	BCryptCost    int           `mapstructure:"bcrypt_cost"`
	
	// OAuth/OIDC settings for enterprise authentication
	OAuthProvider   string `mapstructure:"oauth_provider"`
	OAuthClientID   string `mapstructure:"oauth_client_id"`
	OAuthClientSecret string `mapstructure:"oauth_client_secret"`
	OAuthRedirectURL  string `mapstructure:"oauth_redirect_url"`
}

// EmailConfig contains email service settings for notifications.
type EmailConfig struct {
	Provider string `mapstructure:"provider"`
	APIKey   string `mapstructure:"api_key"`
	From     string `mapstructure:"from"`
	
	// SMTP settings (alternative to API)
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUsername string `mapstructure:"smtp_username"`
	SMTPPassword string `mapstructure:"smtp_password"`
}

// WebhookConfig contains webhook settings for external integrations.
type WebhookConfig struct {
	Secret      string        `mapstructure:"secret"`
	Timeout     time.Duration `mapstructure:"timeout"`
	RetryCount  int           `mapstructure:"retry_count"`
	RetryDelay  time.Duration `mapstructure:"retry_delay"`
}

// MonitoringConfig contains settings for application monitoring and metrics.
type MonitoringConfig struct {
	Enabled           bool   `mapstructure:"enabled"`
	MetricsPath       string `mapstructure:"metrics_path"`
	HealthCheckPath   string `mapstructure:"health_check_path"`
	PrometheusEnabled bool   `mapstructure:"prometheus_enabled"`
}

// Load reads configuration from environment variables, config files, and defaults.
// It follows the 12-factor app methodology for configuration management.
//
// Configuration precedence (highest to lowest):
// 1. Environment variables
// 2. Configuration file (config.yaml, config.json)
// 3. Default values
//
// Returns:
//   - *Config: Loaded configuration instance
//   - error: Configuration loading or validation error
//
// Example:
//   config, err := Load()
//   if err != nil {
//       log.Fatal("Failed to load configuration", err)
//   }
//   
//   server := gin.New()
//   server.Run(fmt.Sprintf(":%d", config.App.Port))
func Load() (*Config, error) {
	// Initialize viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/goedu")
	
	// Set environment variable prefix
	viper.SetEnvPrefix("GOEDU")
	viper.AutomaticEnv()
	
	// Set default values
	setDefaults()
	
	// Read configuration file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is acceptable, we'll use defaults and env vars
	}
	
	// Unmarshal configuration
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	
	// Validate configuration
	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	
	return &config, nil
}

// setDefaults configures default values for all configuration options.
// These defaults are suitable for development and should be overridden in production.
func setDefaults() {
	// Application defaults
	viper.SetDefault("app.name", "GoEdu Control Testing Platform")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.host", "0.0.0.0")
	viper.SetDefault("app.timeout", "30s")
	
	// CORS defaults
	viper.SetDefault("app.cors.allowed_origins", []string{"http://localhost:3000"})
	viper.SetDefault("app.cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("app.cors.allowed_headers", []string{"Authorization", "Content-Type"})
	
	// Database defaults
	viper.SetDefault("database.uri", "mongodb://localhost:27017")
	viper.SetDefault("database.database", "goedu")
	viper.SetDefault("database.max_pool_size", 100)
	viper.SetDefault("database.min_pool_size", 10)
	viper.SetDefault("database.max_conn_idle_time", "5m")
	viper.SetDefault("database.connect_timeout", "10s")
	viper.SetDefault("database.server_select_timeout", "10s")
	
	// Cache defaults
	viper.SetDefault("cache.host", "localhost")
	viper.SetDefault("cache.port", 6379)
	viper.SetDefault("cache.password", "")
	viper.SetDefault("cache.database", 0)
	viper.SetDefault("cache.max_retries", 3)
	viper.SetDefault("cache.pool_size", 10)
	viper.SetDefault("cache.dial_timeout", "5s")
	viper.SetDefault("cache.read_timeout", "3s")
	viper.SetDefault("cache.write_timeout", "3s")
	viper.SetDefault("cache.idle_timeout", "5m")
	
	// Storage defaults
	viper.SetDefault("storage.endpoint", "localhost:9000")
	viper.SetDefault("storage.access_key_id", "minioadmin")
	viper.SetDefault("storage.secret_access_key", "minioadmin")
	viper.SetDefault("storage.bucket_name", "goedu-evidence")
	viper.SetDefault("storage.region", "us-east-1")
	viper.SetDefault("storage.use_ssl", false)
	
	// Auth defaults
	viper.SetDefault("auth.jwt_secret", "your-secret-key-change-in-production")
	viper.SetDefault("auth.jwt_expiration", "24h")
	viper.SetDefault("auth.bcrypt_cost", 12)
	
	// Email defaults
	viper.SetDefault("email.provider", "smtp")
	viper.SetDefault("email.from", "noreply@goedu.com")
	viper.SetDefault("email.smtp_host", "localhost")
	viper.SetDefault("email.smtp_port", 587)
	
	// Webhook defaults
	viper.SetDefault("webhook.timeout", "30s")
	viper.SetDefault("webhook.retry_count", 3)
	viper.SetDefault("webhook.retry_delay", "5s")
	
	// Monitoring defaults
	viper.SetDefault("monitoring.enabled", true)
	viper.SetDefault("monitoring.metrics_path", "/metrics")
	viper.SetDefault("monitoring.health_check_path", "/health")
	viper.SetDefault("monitoring.prometheus_enabled", true)
	
	// Logger defaults
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.environment", "development")
	viper.SetDefault("logger.output_path", "stdout")
}

// validate performs configuration validation to ensure required fields are set
// and values are within acceptable ranges for production deployment.
func validate(config *Config) error {
	// Validate required fields for production
	if config.App.Environment == "production" {
		if config.Auth.JWTSecret == "your-secret-key-change-in-production" {
			return fmt.Errorf("JWT secret must be changed in production")
		}
		
		if config.Database.URI == "mongodb://localhost:27017" {
			return fmt.Errorf("database URI must be configured for production")
		}
		
		if config.Storage.AccessKeyID == "minioadmin" {
			return fmt.Errorf("storage credentials must be configured for production")
		}
	}
	
	// Validate port range
	if config.App.Port < 1024 || config.App.Port > 65535 {
		return fmt.Errorf("app port must be between 1024 and 65535, got %d", config.App.Port)
	}
	
	// Validate database pool sizes
	if config.Database.MaxPoolSize < config.Database.MinPoolSize {
		return fmt.Errorf("database max_pool_size must be >= min_pool_size")
	}
	
	// Validate BCrypt cost
	if config.Auth.BCryptCost < 10 || config.Auth.BCryptCost > 15 {
		return fmt.Errorf("bcrypt cost must be between 10 and 15, got %d", config.Auth.BCryptCost)
	}
	
	return nil
}

// IsDevelopment returns true if the application is running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if the application is running in production mode.
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// GetDatabaseURI returns the complete database connection URI.
// It handles both development and production database configurations.
func (c *Config) GetDatabaseURI() string {
	return c.Database.URI
}

// GetRedisAddr returns the Redis server address in host:port format.
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Cache.Host, c.Cache.Port)
}

// GetServerAddr returns the server address in host:port format.
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.App.Host, c.App.Port)
}