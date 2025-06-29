// Package logger provides structured logging capabilities for the GoEdu application.
// It implements a centralized logging system with different log levels, correlation IDs,
// and audit trail support as required for financial compliance.
package logger

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with additional functionality for audit trails and correlation tracking.
// This is essential for financial institutions that require detailed system monitoring.
type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
}

// CorrelationIDKey is used to extract correlation IDs from context for request tracing.
// This enables tracking individual user actions across multiple system components.
const CorrelationIDKey = "correlation_id"

// LogLevel represents the different logging levels supported by the system.
type LogLevel string

const (
	// DebugLevel provides detailed information for debugging during development
	DebugLevel LogLevel = "debug"
	// InfoLevel provides general information about system operations
	InfoLevel LogLevel = "info"
	// WarnLevel indicates potentially harmful situations that should be monitored
	WarnLevel LogLevel = "warn"
	// ErrorLevel indicates error events that might still allow application to continue
	ErrorLevel LogLevel = "error"
	// FatalLevel indicates very severe error events that will lead to application termination
	FatalLevel LogLevel = "fatal"
)

// Config holds the logger configuration settings.
// This allows for different logging configurations across environments.
type Config struct {
	Level       LogLevel `mapstructure:"level"`
	Environment string   `mapstructure:"environment"`
	OutputPath  string   `mapstructure:"output_path"`
}

// DefaultConfig returns a default logger configuration suitable for development.
// Production environments should override these settings through environment variables.
func DefaultConfig() *Config {
	return &Config{
		Level:       InfoLevel,
		Environment: "development",
		OutputPath:  "stdout",
	}
}

// New creates a new logger instance with the provided configuration.
// The logger is configured with structured output, appropriate for both development
// and production environments with different output formats.
//
// Parameters:
//   - config: Logger configuration specifying level, environment, and output
//
// Returns:
//   - *Logger: Configured logger instance
//   - error: Configuration or initialization error
//
// Example:
//   config := &Config{Level: InfoLevel, Environment: "production", OutputPath: "stdout"}
//   logger, err := New(config)
//   if err != nil {
//       panic(err)
//   }
//   logger.Info("Application started", zap.String("version", "1.0.0"))
func New(config *Config) (*Logger, error) {
	// Convert string level to zapcore.Level
	var level zapcore.Level
	switch config.Level {
	case DebugLevel:
		level = zapcore.DebugLevel
	case InfoLevel:
		level = zapcore.InfoLevel
	case WarnLevel:
		level = zapcore.WarnLevel
	case ErrorLevel:
		level = zapcore.ErrorLevel
	case FatalLevel:
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	// Configure encoder based on environment
	var encoderConfig zapcore.EncoderConfig
	if config.Environment == "production" {
		encoderConfig = zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Create core with appropriate output
	var core zapcore.Core
	if config.OutputPath == "stdout" || config.OutputPath == "" {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
	} else {
		file, err := os.OpenFile(config.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(file),
			level,
		)
	}

	// Create logger with caller information for audit purposes
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}, nil
}

// WithContext creates a new logger with correlation ID from context.
// This is essential for tracking requests across microservices and maintaining audit trails.
//
// Parameters:
//   - ctx: Context containing correlation ID and other request metadata
//
// Returns:
//   - *Logger: Logger instance with context-specific fields
//
// Example:
//   correlationID := uuid.New().String()
//   ctx := context.WithValue(context.Background(), CorrelationIDKey, correlationID)
//   contextLogger := logger.WithContext(ctx)
//   contextLogger.Info("Processing user request", zap.String("user_id", userID))
func (l *Logger) WithContext(ctx context.Context) *Logger {
	fields := []zap.Field{}

	// Add correlation ID if present
	if correlationID, ok := ctx.Value(CorrelationIDKey).(string); ok {
		fields = append(fields, zap.String("correlation_id", correlationID))
	}

	// Add timestamp for audit trail
	fields = append(fields, zap.Time("request_time", time.Now()))

	return &Logger{
		Logger: l.Logger.With(fields...),
		sugar:  l.Logger.With(fields...).Sugar(),
	}
}

// Audit logs audit trail events with specific formatting for compliance requirements.
// Financial institutions require detailed audit logs for regulatory compliance.
//
// Parameters:
//   - ctx: Request context with correlation tracking
//   - action: The action being performed (e.g., "user_login", "control_update")
//   - userID: ID of the user performing the action
//   - resourceID: ID of the resource being acted upon
//   - fields: Additional structured fields for the audit log
//
// Example:
//   logger.Audit(ctx, "control_update", userID, controlID, 
//       zap.String("old_status", "draft"), 
//       zap.String("new_status", "active"))
func (l *Logger) Audit(ctx context.Context, action, userID, resourceID string, fields ...zap.Field) {
	auditFields := []zap.Field{
		zap.String("event_type", "audit"),
		zap.String("action", action),
		zap.String("user_id", userID),
		zap.String("resource_id", resourceID),
		zap.Time("audit_timestamp", time.Now()),
	}

	// Add correlation ID if present in context
	if correlationID, ok := ctx.Value(CorrelationIDKey).(string); ok {
		auditFields = append(auditFields, zap.String("correlation_id", correlationID))
	}

	// Append additional fields
	auditFields = append(auditFields, fields...)

	l.Logger.Info("Audit Event", auditFields...)
}

// Error logs error events with enhanced context for debugging and monitoring.
// Includes stack traces for debugging while maintaining structured format.
//
// Parameters:
//   - ctx: Request context for correlation tracking
//   - message: Human-readable error description
//   - err: The actual error that occurred
//   - fields: Additional structured fields for context
//
// Example:
//   logger.Error(ctx, "Failed to save control", err, 
//       zap.String("control_id", controlID),
//       zap.String("operation", "create"))
func (l *Logger) Error(ctx context.Context, message string, err error, fields ...zap.Field) {
	errorFields := []zap.Field{
		zap.Error(err),
		zap.Time("error_timestamp", time.Now()),
	}

	// Add correlation ID if present in context
	if correlationID, ok := ctx.Value(CorrelationIDKey).(string); ok {
		errorFields = append(errorFields, zap.String("correlation_id", correlationID))
	}

	// Append additional fields
	errorFields = append(errorFields, fields...)

	l.Logger.Error(message, errorFields...)
}

// Performance logs performance metrics for monitoring system performance.
// This is crucial for financial institutions that require performance monitoring.
//
// Parameters:
//   - ctx: Request context for correlation tracking
//   - operation: Name of the operation being measured
//   - duration: Time taken to complete the operation
//   - fields: Additional performance-related fields
//
// Example:
//   start := time.Now()
//   // ... perform operation
//   logger.Performance(ctx, "database_query", time.Since(start),
//       zap.String("table", "controls"),
//       zap.Int("rows_affected", rowCount))
func (l *Logger) Performance(ctx context.Context, operation string, duration time.Duration, fields ...zap.Field) {
	perfFields := []zap.Field{
		zap.String("event_type", "performance"),
		zap.String("operation", operation),
		zap.Duration("duration", duration),
		zap.Time("performance_timestamp", time.Now()),
	}

	// Add correlation ID if present in context
	if correlationID, ok := ctx.Value(CorrelationIDKey).(string); ok {
		perfFields = append(perfFields, zap.String("correlation_id", correlationID))
	}

	// Append additional fields
	perfFields = append(perfFields, fields...)

	l.Logger.Info("Performance Metric", perfFields...)
}

// Sugar returns the sugared logger for less structured logging scenarios.
// Use this for simple logging where structured fields are not required.
//
// Returns:
//   - *zap.SugaredLogger: Sugared logger instance for printf-style logging
//
// Example:
//   logger.Sugar().Infof("Processing %d controls for user %s", count, userID)
func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugar
}

// Sync flushes any buffered log entries.
// This should be called before application shutdown to ensure all logs are written.
//
// Returns:
//   - error: Sync error if any
//
// Example:
//   defer logger.Sync()
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Helper functions for creating zap fields with proper typing
// These functions provide a convenient way to create structured log fields

// String creates a string field for structured logging
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

// Int creates an integer field for structured logging
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Int64 creates an int64 field for structured logging
func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

// Duration creates a duration field for structured logging
func Duration(key string, value time.Duration) zap.Field {
	return zap.Duration(key, value)
}

// Strings creates a string slice field for structured logging
func Strings(key string, value []string) zap.Field {
	return zap.Strings(key, value)
}

// Error creates an error field for structured logging
func Error(err error) zap.Field {
	return zap.Error(err)
}

// Time creates a time field for structured logging
func Time(key string, value time.Time) zap.Field {
	return zap.Time(key, value)
}