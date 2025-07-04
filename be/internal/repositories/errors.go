// Package repositories provides common repository errors and utilities.
package repositories

import "errors"

// Common repository errors
var (
	ErrNotFound = errors.New("entity not found")
	ErrDuplicate = errors.New("duplicate entity")
	ErrInvalidInput = errors.New("invalid input")
	ErrDatabaseConnection = errors.New("database connection error")
)