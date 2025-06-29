// Package services defines the service interfaces for the GoEdu Control Testing Platform.
// These interfaces define the business logic layer contracts following clean architecture principles.
package services

import (
	"context"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/models"
)

// ControlService handles all business logic related to compliance controls.
// It manages control lifecycle, validation, and business rules enforcement.
type ControlService interface {
	// CreateControl creates a new control with validation and audit logging
	CreateControl(ctx context.Context, input *CreateControlInput) (*models.Control, error)
	
	// GetControl retrieves a control by ID with proper authorization checks
	GetControl(ctx context.Context, id string) (*models.Control, error)
	
	// UpdateControl updates an existing control with change tracking
	UpdateControl(ctx context.Context, id string, input *UpdateControlInput) (*models.Control, error)
	
	// DeleteControl soft deletes a control with audit trail
	DeleteControl(ctx context.Context, id string) error
	
	// ListControls retrieves controls with filtering, sorting, and pagination
	ListControls(ctx context.Context, filter *ControlFilter) (*ControlConnection, error)
	
	// ValidateControl performs business rule validation on control data
	ValidateControl(ctx context.Context, control *models.Control) error
	
	// GetControlsByFramework retrieves controls filtered by compliance framework
	GetControlsByFramework(ctx context.Context, framework string, orgID string) ([]*models.Control, error)
}

// TestingService manages testing cycles and control assignments.
// It orchestrates the testing workflow and manages test execution.
type TestingService interface {
	// CreateTestingCycle creates a new testing cycle with control assignments
	CreateTestingCycle(ctx context.Context, input *CreateCycleInput) (*models.TestingCycle, error)
	
	// GetTestingCycle retrieves a testing cycle by ID
	GetTestingCycle(ctx context.Context, id string) (*models.TestingCycle, error)
	
	// UpdateTestingCycle updates cycle information and status
	UpdateTestingCycle(ctx context.Context, id string, input *UpdateCycleInput) (*models.TestingCycle, error)
	
	// AssignControlToAuditor assigns a control to an auditor for testing
	AssignControlToAuditor(ctx context.Context, input *AssignmentInput) (*Assignment, error)
	
	// UpdateTestProgress updates the progress of control testing
	UpdateTestProgress(ctx context.Context, testID string, progress *TestProgress) error
	
	// CompleteTestingCycle marks a testing cycle as complete
	CompleteTestingCycle(ctx context.Context, cycleID string) error
	
	// GenerateWorkpaper creates testing workpapers for controls
	GenerateWorkpaper(ctx context.Context, testID string) (*Document, error)
	
	// GetCycleProgress retrieves progress information for a testing cycle
	GetCycleProgress(ctx context.Context, cycleID string) (*models.Progress, error)
}

// EvidenceService handles evidence collection and management.
// It manages evidence requests, file uploads, and evidence validation.
type EvidenceService interface {
	// CreateEvidenceRequest creates a new evidence request for a control
	CreateEvidenceRequest(ctx context.Context, input *EvidenceRequestInput) (*models.EvidenceRequest, error)
	
	// GetEvidenceRequest retrieves an evidence request by ID
	GetEvidenceRequest(ctx context.Context, id string) (*models.EvidenceRequest, error)
	
	// UpdateEvidenceRequest updates request status and information
	UpdateEvidenceRequest(ctx context.Context, id string, input *UpdateEvidenceRequestInput) (*models.EvidenceRequest, error)
	
	// UploadEvidence handles evidence file uploads with validation
	UploadEvidence(ctx context.Context, requestID string, files []*FileUpload) error
	
	// ProcessEvidenceFile processes uploaded evidence files for text extraction
	ProcessEvidenceFile(ctx context.Context, fileID string) error
	
	// GetEvidenceByRequest retrieves all evidence for a request
	GetEvidenceByRequest(ctx context.Context, requestID string) ([]*models.Evidence, error)
	
	// ValidateEvidence performs business validation on evidence
	ValidateEvidence(ctx context.Context, evidence *models.Evidence) error
	
	// GetPendingRequests retrieves pending evidence requests for a user
	GetPendingRequests(ctx context.Context, userID string) ([]*models.EvidenceRequest, error)
}

// UserService manages user accounts, authentication, and authorization.
// It handles user lifecycle and security operations.
type UserService interface {
	// CreateUser creates a new user account with proper validation
	CreateUser(ctx context.Context, input *CreateUserInput) (*models.User, error)
	
	// GetUser retrieves a user by ID
	GetUser(ctx context.Context, id string) (*models.User, error)
	
	// GetUserByEmail retrieves a user by email address
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	
	// UpdateUser updates user information
	UpdateUser(ctx context.Context, id string, input *UpdateUserInput) (*models.User, error)
	
	// UpdatePassword updates user password with proper hashing
	UpdatePassword(ctx context.Context, userID string, oldPassword, newPassword string) error
	
	// AuthenticateUser validates user credentials
	AuthenticateUser(ctx context.Context, email, password string) (*models.User, error)
	
	// DeactivateUser deactivates a user account
	DeactivateUser(ctx context.Context, id string) error
	
	// ListUsers retrieves users with filtering and pagination
	ListUsers(ctx context.Context, filter *UserFilter) (*UserConnection, error)
	
	// ValidatePermissions checks if user has required permissions
	ValidatePermissions(ctx context.Context, userID string, permissions []string) error
}

// AuditService handles audit logging and compliance tracking.
// It maintains detailed audit trails for regulatory compliance.
type AuditService interface {
	// LogAction records an audit trail entry
	LogAction(ctx context.Context, input *AuditInput) error
	
	// GetAuditTrail retrieves audit log entries with filtering
	GetAuditTrail(ctx context.Context, filter *AuditFilter) ([]*models.AuditLog, error)
	
	// GetUserActivity retrieves activity log for a specific user
	GetUserActivity(ctx context.Context, userID string, timeRange *TimeRange) ([]*models.AuditLog, error)
	
	// GetResourceActivity retrieves activity log for a specific resource
	GetResourceActivity(ctx context.Context, resourceType, resourceID string) ([]*models.AuditLog, error)
	
	// ExportAuditLog exports audit logs for compliance reporting
	ExportAuditLog(ctx context.Context, filter *AuditFilter, format string) (*Document, error)
}

// NotificationService handles system notifications and communications.
// It manages email notifications and system alerts.
type NotificationService interface {
	// SendEvidenceRequest sends notification for new evidence request
	SendEvidenceRequest(ctx context.Context, request *models.EvidenceRequest) error
	
	// SendReminderNotification sends reminder for overdue evidence
	SendReminderNotification(ctx context.Context, request *models.EvidenceRequest) error
	
	// SendTestingCycleNotification notifies about testing cycle updates
	SendTestingCycleNotification(ctx context.Context, cycle *models.TestingCycle, eventType string) error
	
	// SendSystemAlert sends system alerts to administrators
	SendSystemAlert(ctx context.Context, alert *SystemAlert) error
	
	// GetNotificationPreferences retrieves user notification preferences
	GetNotificationPreferences(ctx context.Context, userID string) (*NotificationPreferences, error)
	
	// UpdateNotificationPreferences updates user notification settings
	UpdateNotificationPreferences(ctx context.Context, userID string, prefs *NotificationPreferences) error
}

// Input/Output structures for service operations

// CreateControlInput contains the data needed to create a new control
type CreateControlInput struct {
	OrganizationID   string   `json:"organization_id" validate:"required"`
	ControlID        string   `json:"control_id" validate:"required"`
	Title            string   `json:"title" validate:"required,min=1,max=500"`
	Description      string   `json:"description" validate:"required"`
	Framework        string   `json:"framework" validate:"required"`
	Category         string   `json:"category" validate:"required"`
	SubCategory      string   `json:"sub_category,omitempty"`
	RiskLevel        string   `json:"risk_level" validate:"required"`
	Importance       string   `json:"importance" validate:"required"`
	ControlType      string   `json:"control_type" validate:"required"`
	ControlFrequency string   `json:"control_frequency" validate:"required"`
	Owner            string   `json:"owner" validate:"required"`
	Process          string   `json:"process"`
	Systems          []string `json:"systems,omitempty"`
	TestingProcedure string   `json:"testing_procedure" validate:"required"`
	SampleSize       int      `json:"sample_size"`
	EvidenceTypes    []string `json:"evidence_types"`
	TestingNotes     string   `json:"testing_notes,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

// UpdateControlInput contains the data for updating an existing control
type UpdateControlInput struct {
	Title            *string  `json:"title,omitempty"`
	Description      *string  `json:"description,omitempty"`
	RiskLevel        *string  `json:"risk_level,omitempty"`
	Importance       *string  `json:"importance,omitempty"`
	ControlFrequency *string  `json:"control_frequency,omitempty"`
	Owner            *string  `json:"owner,omitempty"`
	Process          *string  `json:"process,omitempty"`
	Systems          []string `json:"systems,omitempty"`
	TestingProcedure *string  `json:"testing_procedure,omitempty"`
	SampleSize       *int     `json:"sample_size,omitempty"`
	EvidenceTypes    []string `json:"evidence_types,omitempty"`
	TestingNotes     *string  `json:"testing_notes,omitempty"`
	Status           *string  `json:"status,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

// ControlFilter defines filtering options for control queries
type ControlFilter struct {
	OrganizationID string   `json:"organization_id"`
	Framework      string   `json:"framework,omitempty"`
	Category       string   `json:"category,omitempty"`
	RiskLevel      string   `json:"risk_level,omitempty"`
	Status         string   `json:"status,omitempty"`
	Owner          string   `json:"owner,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	
	// Pagination
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	
	// Sorting
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

// ControlConnection represents a paginated list of controls
type ControlConnection struct {
	Nodes      []*models.Control `json:"nodes"`
	TotalCount int               `json:"total_count"`
	HasMore    bool              `json:"has_more"`
}

// Additional service input/output structures...

// CreateCycleInput contains data for creating a testing cycle
type CreateCycleInput struct {
	OrganizationID string   `json:"organization_id" validate:"required"`
	CycleID        string   `json:"cycle_id" validate:"required"`
	Name           string   `json:"name" validate:"required"`
	Description    string   `json:"description,omitempty"`
	StartDate      string   `json:"start_date" validate:"required"`
	EndDate        string   `json:"end_date" validate:"required"`
	ControlScope   []string `json:"control_scope"`
	TestingType    string   `json:"testing_type" validate:"required"`
	Framework      string   `json:"framework"`
}

// UpdateCycleInput contains data for updating a testing cycle
type UpdateCycleInput struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	EndDate     *string  `json:"end_date,omitempty"`
	Status      *string  `json:"status,omitempty"`
}

// AssignmentInput contains data for control assignments
type AssignmentInput struct {
	CycleID    string `json:"cycle_id" validate:"required"`
	ControlID  string `json:"control_id" validate:"required"`
	AuditorID  string `json:"auditor_id" validate:"required"`
	DueDate    string `json:"due_date" validate:"required"`
	Priority   string `json:"priority,omitempty"`
	Instructions string `json:"instructions,omitempty"`
}

// Assignment represents a control testing assignment
type Assignment struct {
	ID           string `json:"id"`
	CycleID      string `json:"cycle_id"`
	ControlID    string `json:"control_id"`
	AuditorID    string `json:"auditor_id"`
	Status       string `json:"status"`
	AssignedDate string `json:"assigned_date"`
	DueDate      string `json:"due_date"`
}

// TestProgress represents testing progress information
type TestProgress struct {
	Status      string `json:"status"`
	Progress    int    `json:"progress"`
	Notes       string `json:"notes,omitempty"`
	CompletedAt string `json:"completed_at,omitempty"`
}

// Document represents a generated document
type Document struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  []byte `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// EvidenceRequestInput contains data for creating evidence requests
type EvidenceRequestInput struct {
	OrganizationID string   `json:"organization_id" validate:"required"`
	ControlID      string   `json:"control_id" validate:"required"`
	CycleID        string   `json:"cycle_id" validate:"required"`
	AssigneeID     string   `json:"assignee_id" validate:"required"`
	Title          string   `json:"title" validate:"required"`
	Description    string   `json:"description" validate:"required"`
	DueDate        string   `json:"due_date" validate:"required"`
	EvidenceTypes  []string `json:"evidence_types"`
	SampleSize     int      `json:"sample_size"`
	Instructions   string   `json:"instructions,omitempty"`
}

// UpdateEvidenceRequestInput contains data for updating evidence requests
type UpdateEvidenceRequestInput struct {
	Status       *string `json:"status,omitempty"`
	Response     *string `json:"response,omitempty"`
	CompletedAt  *string `json:"completed_at,omitempty"`
}

// FileUpload represents an uploaded file
type FileUpload struct {
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	FileType    string `json:"file_type"`
	Content     []byte `json:"content"`
	Description string `json:"description,omitempty"`
}

// CreateUserInput contains data for creating users
type CreateUserInput struct {
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	OrganizationID string `json:"organization_id" validate:"required"`
	Role           string `json:"role" validate:"required"`
	Title          string `json:"title,omitempty"`
	Department     string `json:"department,omitempty"`
	Phone          string `json:"phone,omitempty"`
}

// UpdateUserInput contains data for updating users
type UpdateUserInput struct {
	FirstName  *string  `json:"first_name,omitempty"`
	LastName   *string  `json:"last_name,omitempty"`
	Title      *string  `json:"title,omitempty"`
	Department *string  `json:"department,omitempty"`
	Phone      *string  `json:"phone,omitempty"`
	Status     *string  `json:"status,omitempty"`
	Role       *string  `json:"role,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// UserFilter defines filtering options for user queries
type UserFilter struct {
	OrganizationID string `json:"organization_id"`
	Role           string `json:"role,omitempty"`
	Status         string `json:"status,omitempty"`
	Department     string `json:"department,omitempty"`
	
	// Pagination
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// UserConnection represents a paginated list of users
type UserConnection struct {
	Nodes      []*models.User `json:"nodes"`
	TotalCount int            `json:"total_count"`
	HasMore    bool           `json:"has_more"`
}

// AuditInput contains data for audit logging
type AuditInput struct {
	UserID         string                 `json:"user_id" validate:"required"`
	OrganizationID string                 `json:"organization_id" validate:"required"`
	Action         string                 `json:"action" validate:"required"`
	ResourceType   string                 `json:"resource_type" validate:"required"`
	ResourceID     string                 `json:"resource_id" validate:"required"`
	OldValues      map[string]interface{} `json:"old_values,omitempty"`
	NewValues      map[string]interface{} `json:"new_values,omitempty"`
	IPAddress      string                 `json:"ip_address,omitempty"`
	UserAgent      string                 `json:"user_agent,omitempty"`
	Success        bool                   `json:"success"`
	ErrorMessage   string                 `json:"error_message,omitempty"`
}

// AuditFilter defines filtering options for audit queries
type AuditFilter struct {
	OrganizationID string     `json:"organization_id"`
	UserID         string     `json:"user_id,omitempty"`
	Action         string     `json:"action,omitempty"`
	ResourceType   string     `json:"resource_type,omitempty"`
	ResourceID     string     `json:"resource_id,omitempty"`
	TimeRange      *TimeRange `json:"time_range,omitempty"`
	
	// Pagination
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// TimeRange represents a time range filter
type TimeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// SystemAlert represents a system alert
type SystemAlert struct {
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Recipients  []string               `json:"recipients"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NotificationPreferences represents user notification settings
type NotificationPreferences struct {
	Email           bool `json:"email"`
	SMS             bool `json:"sms"`
	InApp           bool `json:"in_app"`
	EvidenceRequest bool `json:"evidence_request"`
	Reminders       bool `json:"reminders"`
	SystemAlerts    bool `json:"system_alerts"`
}