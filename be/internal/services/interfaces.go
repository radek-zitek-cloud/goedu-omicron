// Package services defines the service interfaces for the GoEdu Control Testing Platform.
// These interfaces define the business logic layer contracts following clean architecture principles.
package services

import (
	"context"
	"time"

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

// AuthenticationService handles user authentication, authorization, and security operations.
// It provides comprehensive authentication features including JWT token management,
// password hashing, session management, and role-based access control.
type AuthenticationService interface {
	// Authentication operations
	Login(ctx context.Context, request *models.LoginRequest) (*models.LoginResponse, error)
	Logout(ctx context.Context, sessionID string) error
	RefreshToken(ctx context.Context, refreshToken string) (*models.LoginResponse, error)
	
	// Password management
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, email string) error
	ValidatePasswordReset(ctx context.Context, token, newPassword string) error
	
	// Multi-factor authentication
	EnableMFA(ctx context.Context, userID string) (*MFASetupResponse, error)
	DisableMFA(ctx context.Context, userID, mfaCode string) error
	ValidateMFA(ctx context.Context, userID, mfaCode string) error
	GenerateBackupCodes(ctx context.Context, userID string) ([]string, error)
	
	// Session management
	CreateSession(ctx context.Context, userID, ipAddress, userAgent string) (*models.Session, error)
	GetSession(ctx context.Context, sessionID string) (*models.Session, error)
	UpdateSessionActivity(ctx context.Context, sessionID string) error
	TerminateSession(ctx context.Context, sessionID string) error
	TerminateAllSessions(ctx context.Context, userID string) error
	
	// Token operations
	ValidateAccessToken(ctx context.Context, token string) (*models.JWTClaims, error)
	GenerateAccessToken(ctx context.Context, user *models.User, sessionID, ipAddress string) (string, time.Time, error)
	GenerateRefreshToken(ctx context.Context, userID, sessionID string) (string, time.Time, error)
	
	// Account security
	LockAccount(ctx context.Context, userID string, reason string) error
	UnlockAccount(ctx context.Context, userID string) error
	RecordFailedLogin(ctx context.Context, userID, ipAddress string) error
	RecordSuccessfulLogin(ctx context.Context, userID, ipAddress string) error
	
	// Security questions
	SetSecurityQuestions(ctx context.Context, userID string, questions []models.SecurityQuestion) error
	ValidateSecurityAnswer(ctx context.Context, userID, question, answer string) error
	
	// Audit and compliance
	LogSecurityEvent(ctx context.Context, event *models.AuditEvent) error
	GetSecurityEvents(ctx context.Context, userID string, timeRange *TimeRange) ([]*models.AuditEvent, error)
}

// PermissionService handles role-based access control and permission management.
// It provides comprehensive RBAC functionality with hierarchical permissions.
type PermissionService interface {
	// Permission checking
	HasPermission(ctx context.Context, userID, resource, action, scope string) (bool, error)
	ValidatePermission(ctx context.Context, userID, resource, action, scope string) error
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	
	// Role management
	CreateRole(ctx context.Context, role *models.Role) error
	GetRole(ctx context.Context, roleID string) (*models.Role, error)
	UpdateRole(ctx context.Context, roleID string, updates *UpdateRoleInput) error
	DeleteRole(ctx context.Context, roleID string) error
	ListRoles(ctx context.Context, organizationID string) ([]*models.Role, error)
	
	// User role assignment
	AssignRole(ctx context.Context, userID, roleID string) error
	RevokeRole(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]*models.Role, error)
	
	// Permission management
	CreatePermission(ctx context.Context, permission *models.Permission) error
	GetPermission(ctx context.Context, permissionID string) (*models.Permission, error)
	UpdatePermission(ctx context.Context, permissionID string, updates *UpdatePermissionInput) error
	DeletePermission(ctx context.Context, permissionID string) error
	ListPermissions(ctx context.Context, filter *PermissionFilter) ([]*models.Permission, error)
	
	// Role-permission association
	GrantPermissionToRole(ctx context.Context, roleID, permissionID string) error
	RevokePermissionFromRole(ctx context.Context, roleID, permissionID string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]*models.Permission, error)
	
	// Bulk operations
	BulkAssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error
	BulkRevokePermissions(ctx context.Context, roleID string, permissionIDs []string) error
	SyncUserPermissions(ctx context.Context, userID string) error
}

// OrganizationService handles organization management and multi-tenancy operations.
// It provides comprehensive organization lifecycle management, subscription handling,
// and feature flag management for the multi-tenant platform.
type OrganizationService interface {
	// Organization CRUD operations
	CreateOrganization(ctx context.Context, input *CreateOrganizationInput) (*models.Organization, error)
	GetOrganization(ctx context.Context, id string) (*models.Organization, error)
	GetOrganizationBySlug(ctx context.Context, slug string) (*models.Organization, error)
	UpdateOrganization(ctx context.Context, id string, input *UpdateOrganizationInput) (*models.Organization, error)
	DeleteOrganization(ctx context.Context, id string) error
	
	// Organization listing and filtering
	ListOrganizations(ctx context.Context, filter *OrganizationFilter) (*OrganizationConnection, error)
	GetActiveOrganizations(ctx context.Context, limit, offset int) ([]*models.Organization, error)
	
	// Subscription management
	UpdateSubscription(ctx context.Context, orgID string, subscription *models.OrganizationSubscription) error
	GetSubscriptionStatus(ctx context.Context, orgID string) (*models.OrganizationSubscription, error)
	UpgradeSubscription(ctx context.Context, orgID, newPlan string) error
	DowngradeSubscription(ctx context.Context, orgID, newPlan string) error
	CancelSubscription(ctx context.Context, orgID string) error
	RenewSubscription(ctx context.Context, orgID string) error
	
	// Feature flag management
	GetFeatureFlags(ctx context.Context, orgID string) (map[string]bool, error)
	UpdateFeatureFlag(ctx context.Context, orgID, flag string, enabled bool) error
	BulkUpdateFeatureFlags(ctx context.Context, orgID string, flags map[string]bool) error
	IsFeatureEnabled(ctx context.Context, orgID, feature string) (bool, error)
	
	// Organization settings management
	GetSettings(ctx context.Context, orgID string) (*models.OrganizationSettings, error)
	UpdateSettings(ctx context.Context, orgID string, settings *models.OrganizationSettings) error
	UpdatePartialSettings(ctx context.Context, orgID string, updates map[string]interface{}) error
	
	// Member management
	GetMemberCount(ctx context.Context, orgID string) (int, error)
	UpdateMemberCount(ctx context.Context, orgID string, count int) error
	GetMemberLimits(ctx context.Context, orgID string) (current, max int, error)
	CanAddMember(ctx context.Context, orgID string) (bool, error)
	
	// Organization validation and compliance
	ValidateOrganization(ctx context.Context, org *models.Organization) error
	CheckComplianceRequirements(ctx context.Context, orgID string) (*ComplianceStatus, error)
	UpdateRegulatoryProfile(ctx context.Context, orgID string, profile *models.RegulatoryProfile) error
	
	// Organization statistics and analytics
	GetOrganizationStats(ctx context.Context, orgID string) (*OrganizationStats, error)
	GetUsageMetrics(ctx context.Context, orgID string, timeRange *TimeRange) (*UsageMetrics, error)
	
	// Multi-tenancy support methods
	GetOrganizationContext(ctx context.Context, userID string) (*OrganizationContext, error)
	ValidateOrganizationAccess(ctx context.Context, userID, orgID string) error
	GetUserOrganizations(ctx context.Context, userID string) ([]*models.Organization, error)
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

// Authentication service input/output structures

// MFASetupResponse contains MFA setup information
type MFASetupResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

// UpdateRoleInput contains data for updating roles
type UpdateRoleInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
}

// UpdatePermissionInput contains data for updating permissions
type UpdatePermissionInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	Category    *string `json:"category,omitempty"`
}

// PermissionFilter defines filtering options for permission queries
type PermissionFilter struct {
	Resource     string `json:"resource,omitempty"`
	Action       string `json:"action,omitempty"`
	Scope        string `json:"scope,omitempty"`
	Category     string `json:"category,omitempty"`
	IsActive     *bool  `json:"is_active,omitempty"`
	
	// Pagination
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Organization service input/output structures

// CreateOrganizationInput contains data for creating a new organization
type CreateOrganizationInput struct {
	// Basic information
	Name        string `json:"name" validate:"required,min=1,max=255"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
	
	// Organization classification
	Type     string `json:"type" validate:"required"` // commercial_bank, credit_union, etc.
	Industry string `json:"industry" validate:"required"`
	Size     string `json:"size,omitempty"`
	Region   string `json:"region,omitempty"`
	Country  string `json:"country,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	Currency string `json:"currency,omitempty"`
	
	// Contact information
	ContactEmail string  `json:"contact_email" validate:"required,email"`
	ContactPhone string  `json:"contact_phone,omitempty"`
	Website      string  `json:"website,omitempty"`
	LogoURL      string  `json:"logo_url,omitempty"`
	Address      *CreateAddressInput `json:"address,omitempty"`
	
	// Regulatory profile
	RegulatoryProfile *CreateRegulatoryProfileInput `json:"regulatory_profile,omitempty"`
	
	// Subscription plan
	SubscriptionPlan string `json:"subscription_plan,omitempty"` // starter, professional, enterprise
	
	// Initial settings
	Settings *CreateOrganizationSettingsInput `json:"settings,omitempty"`
}

// UpdateOrganizationInput contains data for updating an organization
type UpdateOrganizationInput struct {
	// Basic information
	Name        *string `json:"name,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
	
	// Organization classification
	Type     *string `json:"type,omitempty"`
	Industry *string `json:"industry,omitempty"`
	Size     *string `json:"size,omitempty"`
	Region   *string `json:"region,omitempty"`
	Country  *string `json:"country,omitempty"`
	Timezone *string `json:"timezone,omitempty"`
	Currency *string `json:"currency,omitempty"`
	
	// Contact information
	ContactEmail *string `json:"contact_email,omitempty"`
	ContactPhone *string `json:"contact_phone,omitempty"`
	Website      *string `json:"website,omitempty"`
	LogoURL      *string `json:"logo_url,omitempty"`
	
	// Status and limits
	Status     *string `json:"status,omitempty"`
	MaxMembers *int    `json:"max_members,omitempty"`
}

// CreateAddressInput contains address information for organization creation
type CreateAddressInput struct {
	Street1    string `json:"street1,omitempty"`
	Street2    string `json:"street2,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country,omitempty"`
}

// CreateRegulatoryProfileInput contains regulatory profile information
type CreateRegulatoryProfileInput struct {
	Industry             string   `json:"industry"`
	PrimaryRegulator     string   `json:"primary_regulator"`
	ApplicableFrameworks []string `json:"applicable_frameworks"`
	ExamCycle            string   `json:"exam_cycle"`
	RetentionPeriod      int      `json:"retention_period"`
	RiskTolerance        string   `json:"risk_tolerance"`
	RequiresSOX          bool     `json:"requires_sox"`
	RequiresPCIDSS       bool     `json:"requires_pci_dss"`
	RequiresFFIEC        bool     `json:"requires_ffiec"`
	RequiresBaselIII     bool     `json:"requires_basel_iii"`
}

// CreateOrganizationSettingsInput contains initial organization settings
type CreateOrganizationSettingsInput struct {
	RequireMFA            bool `json:"require_mfa"`
	AllowInvitations      bool `json:"allow_invitations"`
	SessionTimeoutMinutes int  `json:"session_timeout_minutes"`
	EnableAuditLog        bool `json:"enable_audit_log"`
	DataRetentionDays     int  `json:"data_retention_days"`
	AllowDataExport       bool `json:"allow_data_export"`
}

// OrganizationFilter defines filtering options for organization queries
type OrganizationFilter struct {
	Type     string `json:"type,omitempty"`
	Industry string `json:"industry,omitempty"`
	Status   string `json:"status,omitempty"`
	Region   string `json:"region,omitempty"`
	Country  string `json:"country,omitempty"`
	Plan     string `json:"plan,omitempty"`
	
	// Search
	Search string `json:"search,omitempty"`
	
	// Pagination
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	
	// Sorting
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

// OrganizationConnection represents a paginated list of organizations
type OrganizationConnection struct {
	Nodes      []*models.Organization `json:"nodes"`
	TotalCount int                    `json:"total_count"`
	HasMore    bool                   `json:"has_more"`
}

// ComplianceStatus represents organization compliance status
type ComplianceStatus struct {
	OverallStatus        string            `json:"overall_status"`
	RequiredFrameworks   []string          `json:"required_frameworks"`
	CompletedFrameworks  []string          `json:"completed_frameworks"`
	PendingRequirements  []string          `json:"pending_requirements"`
	ComplianceScore      float64           `json:"compliance_score"`
	LastAssessmentDate   time.Time         `json:"last_assessment_date"`
	NextAssessmentDate   time.Time         `json:"next_assessment_date"`
	FrameworkStatus      map[string]string `json:"framework_status"`
}

// OrganizationStats represents organization statistics and metrics
type OrganizationStats struct {
	TotalUsers        int                `json:"total_users"`
	ActiveUsers       int                `json:"active_users"`
	TotalControls     int                `json:"total_controls"`
	ActiveCycles      int                `json:"active_cycles"`
	PendingEvidence   int                `json:"pending_evidence"`
	StorageUsed       int64              `json:"storage_used"`
	LastActivity      time.Time          `json:"last_activity"`
	UsersByRole       map[string]int     `json:"users_by_role"`
	ControlsByFramework map[string]int   `json:"controls_by_framework"`
}

// UsageMetrics represents organization usage metrics over time
type UsageMetrics struct {
	Period           string                   `json:"period"`
	UserActivity     []UserActivityMetric     `json:"user_activity"`
	SystemUsage      []SystemUsageMetric      `json:"system_usage"`
	FeatureUsage     map[string]int           `json:"feature_usage"`
	StorageGrowth    []StorageMetric          `json:"storage_growth"`
	ComplianceMetrics []ComplianceMetric      `json:"compliance_metrics"`
}

// UserActivityMetric represents user activity over time
type UserActivityMetric struct {
	Date       time.Time `json:"date"`
	ActiveUsers int      `json:"active_users"`
	Logins     int      `json:"logins"`
	Sessions   int      `json:"sessions"`
}

// SystemUsageMetric represents system usage metrics
type SystemUsageMetric struct {
	Date      time.Time `json:"date"`
	Requests  int       `json:"requests"`
	Errors    int       `json:"errors"`
	AvgResponseTime float64 `json:"avg_response_time"`
}

// StorageMetric represents storage usage over time
type StorageMetric struct {
	Date        time.Time `json:"date"`
	TotalBytes  int64     `json:"total_bytes"`
	FileCount   int       `json:"file_count"`
	GrowthRate  float64   `json:"growth_rate"`
}

// ComplianceMetric represents compliance metrics over time
type ComplianceMetric struct {
	Date            time.Time `json:"date"`
	ComplianceScore float64   `json:"compliance_score"`
	OpenFindings    int       `json:"open_findings"`
	ClosedFindings  int       `json:"closed_findings"`
}

// OrganizationContext represents the organization context for a user
type OrganizationContext struct {
	OrganizationID   string                     `json:"organization_id"`
	Organization     *models.Organization       `json:"organization"`
	UserRole         string                     `json:"user_role"`
	UserPermissions  []string                   `json:"user_permissions"`
	FeatureFlags     map[string]bool            `json:"feature_flags"`
	Settings         *models.OrganizationSettings `json:"settings"`
}