// Package repositories defines repository interfaces for data access layer
// in the GoEdu Control Testing Platform. These interfaces abstract database
// operations and provide consistent data access patterns.
package repositories

import (
	"context"
	"time"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/models"
)

// Repository defines common operations available on all repositories.
// This interface provides consistent CRUD operations across all entities.
type Repository[T any] interface {
	// Create inserts a new entity into the database
	Create(ctx context.Context, entity *T) error
	
	// GetByID retrieves an entity by its ID
	GetByID(ctx context.Context, id string) (*T, error)
	
	// Update updates an existing entity in the database
	Update(ctx context.Context, entity *T) error
	
	// Delete removes an entity from the database (soft delete)
	Delete(ctx context.Context, id string) error
	
	// List retrieves entities with optional filtering and pagination
	List(ctx context.Context, filter interface{}) ([]*T, error)
	
	// Count returns the total number of entities matching the filter
	Count(ctx context.Context, filter interface{}) (int64, error)
}

// OrganizationRepository handles data access for organizations.
// It provides organization-specific query methods and multi-tenancy support.
type OrganizationRepository interface {
	Repository[models.Organization]
	
	// GetBySlug retrieves an organization by its unique slug
	GetBySlug(ctx context.Context, slug string) (*models.Organization, error)
	
	// GetActiveOrganizations retrieves all active organizations
	GetActiveOrganizations(ctx context.Context, limit, offset int) ([]*models.Organization, error)
	
	// UpdateSettings updates organization-specific settings
	UpdateSettings(ctx context.Context, orgID string, settings map[string]interface{}) error
	
	// GetFeatureFlags returns feature flags for an organization
	GetFeatureFlags(ctx context.Context, orgID string) (map[string]bool, error)
	
	// UpdateFeatureFlag updates a specific feature flag
	UpdateFeatureFlag(ctx context.Context, orgID, flag string, enabled bool) error
}

// UserRepository handles data access for user accounts.
// It provides user authentication and authorization data operations.
type UserRepository interface {
	Repository[models.User]
	
	// GetByEmail retrieves a user by email address
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	
	// GetByOrganization retrieves users belonging to an organization
	GetByOrganization(ctx context.Context, orgID string, limit, offset int) ([]*models.User, error)
	
	// GetByRole retrieves users with a specific role in an organization
	GetByRole(ctx context.Context, orgID, role string) ([]*models.User, error)
	
	// UpdatePassword updates a user's password hash
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
	
	// UpdateLastLogin updates the user's last login timestamp
	UpdateLastLogin(ctx context.Context, userID string) error
	
	// IncrementFailedLogins increments the failed login counter
	IncrementFailedLogins(ctx context.Context, userID string) error
	
	// ResetFailedLogins resets the failed login counter
	ResetFailedLogins(ctx context.Context, userID string) error
	
	// LockUser temporarily locks a user account
	LockUser(ctx context.Context, userID string, until time.Time) error
	
	// UpdatePreferences updates user preferences
	UpdatePreferences(ctx context.Context, userID string, preferences map[string]interface{}) error
}

// ControlRepository handles data access for compliance controls.
// It provides control-specific queries and framework-based filtering.
type ControlRepository interface {
	Repository[models.Control]
	
	// GetByControlID retrieves a control by its business control ID within an organization
	GetByControlID(ctx context.Context, orgID, controlID string) (*models.Control, error)
	
	// GetByOrganization retrieves controls for a specific organization
	GetByOrganization(ctx context.Context, orgID string, filter *ControlFilter) ([]*models.Control, error)
	
	// GetByFramework retrieves controls for a specific compliance framework
	GetByFramework(ctx context.Context, orgID, framework string) ([]*models.Control, error)
	
	// GetByCategory retrieves controls by category within an organization
	GetByCategory(ctx context.Context, orgID, category string) ([]*models.Control, error)
	
	// GetByOwner retrieves controls assigned to a specific owner
	GetByOwner(ctx context.Context, orgID, owner string) ([]*models.Control, error)
	
	// GetByRiskLevel retrieves controls by risk level
	GetByRiskLevel(ctx context.Context, orgID, riskLevel string) ([]*models.Control, error)
	
	// Search performs full-text search on control title and description
	Search(ctx context.Context, orgID, query string, limit, offset int) ([]*models.Control, error)
	
	// GetControlStats returns statistics about controls in an organization
	GetControlStats(ctx context.Context, orgID string) (*ControlStats, error)
	
	// BulkUpdate updates multiple controls in a single operation
	BulkUpdate(ctx context.Context, updates []*ControlUpdate) error
}

// TestingCycleRepository handles data access for testing cycles.
// It manages testing cycle lifecycle and progress tracking.
type TestingCycleRepository interface {
	Repository[models.TestingCycle]
	
	// GetByCycleID retrieves a testing cycle by its business cycle ID
	GetByCycleID(ctx context.Context, orgID, cycleID string) (*models.TestingCycle, error)
	
	// GetByOrganization retrieves testing cycles for an organization
	GetByOrganization(ctx context.Context, orgID string, limit, offset int) ([]*models.TestingCycle, error)
	
	// GetActiveByOrganization retrieves active testing cycles
	GetActiveByOrganization(ctx context.Context, orgID string) ([]*models.TestingCycle, error)
	
	// GetByDateRange retrieves cycles within a date range
	GetByDateRange(ctx context.Context, orgID string, start, end time.Time) ([]*models.TestingCycle, error)
	
	// UpdateProgress updates the progress of a testing cycle
	UpdateProgress(ctx context.Context, cycleID string, progress *models.Progress) error
	
	// GetCyclesByControl returns cycles that include a specific control
	GetCyclesByControl(ctx context.Context, controlID string) ([]*models.TestingCycle, error)
}

// EvidenceRequestRepository handles data access for evidence requests.
// It manages evidence collection workflow and assignment tracking.
type EvidenceRequestRepository interface {
	Repository[models.EvidenceRequest]
	
	// GetByRequestID retrieves an evidence request by its business request ID
	GetByRequestID(ctx context.Context, orgID, requestID string) (*models.EvidenceRequest, error)
	
	// GetByAssignee retrieves evidence requests assigned to a user
	GetByAssignee(ctx context.Context, assigneeID string, status string) ([]*models.EvidenceRequest, error)
	
	// GetByControl retrieves evidence requests for a specific control
	GetByControl(ctx context.Context, controlID string) ([]*models.EvidenceRequest, error)
	
	// GetByCycle retrieves evidence requests for a testing cycle
	GetByCycle(ctx context.Context, cycleID string) ([]*models.EvidenceRequest, error)
	
	// GetOverdueRequests retrieves overdue evidence requests
	GetOverdueRequests(ctx context.Context, orgID string) ([]*models.EvidenceRequest, error)
	
	// GetPendingRequests retrieves pending evidence requests for an organization
	GetPendingRequests(ctx context.Context, orgID string, limit, offset int) ([]*models.EvidenceRequest, error)
	
	// UpdateStatus updates the status of an evidence request
	UpdateStatus(ctx context.Context, requestID, status string) error
	
	// AddEvidence adds evidence to an evidence request
	AddEvidence(ctx context.Context, requestID string, evidence *models.Evidence) error
	
	// AddComment adds a comment to an evidence request
	AddComment(ctx context.Context, requestID string, comment *models.Comment) error
	
	// GetRequestStats returns statistics about evidence requests
	GetRequestStats(ctx context.Context, orgID string) (*EvidenceRequestStats, error)
}

// AuditLogRepository handles data access for audit trail entries.
// It provides audit logging and compliance tracking capabilities.
type AuditLogRepository interface {
	// Create inserts a new audit log entry
	Create(ctx context.Context, entry *models.AuditLog) error
	
	// GetByUser retrieves audit logs for a specific user
	GetByUser(ctx context.Context, userID string, limit, offset int) ([]*models.AuditLog, error)
	
	// GetByOrganization retrieves audit logs for an organization
	GetByOrganization(ctx context.Context, orgID string, filter *AuditFilter) ([]*models.AuditLog, error)
	
	// GetByResource retrieves audit logs for a specific resource
	GetByResource(ctx context.Context, resourceType, resourceID string) ([]*models.AuditLog, error)
	
	// GetByAction retrieves audit logs for a specific action type
	GetByAction(ctx context.Context, orgID, action string, limit, offset int) ([]*models.AuditLog, error)
	
	// GetByTimeRange retrieves audit logs within a time range
	GetByTimeRange(ctx context.Context, orgID string, start, end time.Time) ([]*models.AuditLog, error)
	
	// GetByCorrelationID retrieves audit logs for a correlation ID
	GetByCorrelationID(ctx context.Context, correlationID string) ([]*models.AuditLog, error)
	
	// Search performs search across audit log entries
	Search(ctx context.Context, orgID, query string, filter *AuditFilter) ([]*models.AuditLog, error)
	
	// GetAuditStats returns audit statistics for an organization
	GetAuditStats(ctx context.Context, orgID string, timeRange *TimeRange) (*AuditStats, error)
	
	// Purge removes old audit logs based on retention policy
	Purge(ctx context.Context, retentionDays int) (int64, error)
}

// Filter and Stats structures

// ControlFilter defines filtering options for control queries
type ControlFilter struct {
	Framework   string   `json:"framework,omitempty"`
	Category    string   `json:"category,omitempty"`
	RiskLevel   string   `json:"risk_level,omitempty"`
	Status      string   `json:"status,omitempty"`
	Owner       string   `json:"owner,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	SearchQuery string   `json:"search_query,omitempty"`
	
	// Pagination
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	
	// Sorting
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

// ControlStats represents control statistics
type ControlStats struct {
	TotalControls    int            `json:"total_controls"`
	ByFramework      map[string]int `json:"by_framework"`
	ByCategory       map[string]int `json:"by_category"`
	ByRiskLevel      map[string]int `json:"by_risk_level"`
	ByStatus         map[string]int `json:"by_status"`
	RecentlyCreated  int            `json:"recently_created"`
	RecentlyModified int            `json:"recently_modified"`
}

// ControlUpdate represents a control update operation
type ControlUpdate struct {
	ID     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

// EvidenceRequestStats represents evidence request statistics
type EvidenceRequestStats struct {
	TotalRequests   int            `json:"total_requests"`
	ByStatus        map[string]int `json:"by_status"`
	OverdueCount    int            `json:"overdue_count"`
	CompletedToday  int            `json:"completed_today"`
	PendingCount    int            `json:"pending_count"`
	AverageTime     float64        `json:"average_completion_time_hours"`
}

// AuditFilter defines filtering options for audit log queries
type AuditFilter struct {
	UserID         string     `json:"user_id,omitempty"`
	Action         string     `json:"action,omitempty"`
	ResourceType   string     `json:"resource_type,omitempty"`
	ResourceID     string     `json:"resource_id,omitempty"`
	Success        *bool      `json:"success,omitempty"`
	CorrelationID  string     `json:"correlation_id,omitempty"`
	TimeRange      *TimeRange `json:"time_range,omitempty"`
	IPAddress      string     `json:"ip_address,omitempty"`
	
	// Pagination
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	
	// Sorting
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

// TimeRange represents a time range filter
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// AuditStats represents audit statistics
type AuditStats struct {
	TotalEvents      int            `json:"total_events"`
	ByAction         map[string]int `json:"by_action"`
	ByUser           map[string]int `json:"by_user"`
	ByResourceType   map[string]int `json:"by_resource_type"`
	SuccessfulEvents int            `json:"successful_events"`
	FailedEvents     int            `json:"failed_events"`
	EventsToday      int            `json:"events_today"`
	EventsThisWeek   int            `json:"events_this_week"`
}

// Cache-related interfaces for Redis operations

// CacheRepository provides caching operations for repositories.
// It helps improve performance by caching frequently accessed data.
type CacheRepository interface {
	// Set stores a value in the cache with expiration
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	
	// Get retrieves a value from the cache
	Get(ctx context.Context, key string, dest interface{}) error
	
	// Delete removes keys from the cache
	Delete(ctx context.Context, keys ...string) error
	
	// Exists checks if keys exist in the cache
	Exists(ctx context.Context, keys ...string) (int64, error)
	
	// Invalidate removes cache entries by pattern
	Invalidate(ctx context.Context, pattern string) error
	
	// GetStats returns cache statistics
	GetStats(ctx context.Context) (map[string]interface{}, error)
}

