// Package models contains the core domain models for the GoEdu Control Testing Platform.
// These models represent the fundamental entities in the control testing domain,
// designed for financial compliance and audit requirements.
package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseModel contains common fields for all domain entities.
// This provides consistent auditing and tracking across all models.
type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy string             `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy string             `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
}

// Organization represents a client organization using the platform.
// Organizations provide multi-tenancy and data isolation for different clients.
type Organization struct {
	BaseModel `bson:",inline"`
	
	// Basic organization information
	Name        string `bson:"name" json:"name" validate:"required,min=1,max=255"`
	Slug        string `bson:"slug" json:"slug" validate:"required,min=1,max=100"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	
	// Contact information
	ContactEmail string `bson:"contact_email" json:"contact_email" validate:"required,email"`
	ContactPhone string `bson:"contact_phone,omitempty" json:"contact_phone,omitempty"`
	
	// Address information
	Address Address `bson:"address,omitempty" json:"address,omitempty"`
	
	// Regulatory profile for compliance requirements
	RegulatoryProfile RegulatoryProfile `bson:"regulatory_profile" json:"regulatory_profile"`
	
	// Subscription and feature flags
	SubscriptionTier string            `bson:"subscription_tier" json:"subscription_tier"`
	FeatureFlags     map[string]bool   `bson:"feature_flags,omitempty" json:"feature_flags,omitempty"`
	
	// Status and metadata
	Status   string                 `bson:"status" json:"status"`
	Settings map[string]interface{} `bson:"settings,omitempty" json:"settings,omitempty"`
}

// Address represents a physical address for organizations.
type Address struct {
	Street1    string `bson:"street1,omitempty" json:"street1,omitempty"`
	Street2    string `bson:"street2,omitempty" json:"street2,omitempty"`
	City       string `bson:"city,omitempty" json:"city,omitempty"`
	State      string `bson:"state,omitempty" json:"state,omitempty"`
	PostalCode string `bson:"postal_code,omitempty" json:"postal_code,omitempty"`
	Country    string `bson:"country,omitempty" json:"country,omitempty"`
}

// RegulatoryProfile defines compliance requirements for an organization.
type RegulatoryProfile struct {
	Industry    string   `bson:"industry" json:"industry"`
	Regulations []string `bson:"regulations" json:"regulations"`
	
	// Compliance requirements
	RequiredFrameworks []string `bson:"required_frameworks" json:"required_frameworks"`
	AuditFrequency     string   `bson:"audit_frequency" json:"audit_frequency"`
	RetentionPeriod    int      `bson:"retention_period" json:"retention_period"` // in years
	
	// Risk profile
	RiskTolerance string `bson:"risk_tolerance" json:"risk_tolerance"`
	
	// Certification requirements
	Certifications []Certification `bson:"certifications,omitempty" json:"certifications,omitempty"`
}

// Certification represents regulatory certifications held by the organization.
type Certification struct {
	Name       string    `bson:"name" json:"name"`
	Number     string    `bson:"number,omitempty" json:"number,omitempty"`
	Issuer     string    `bson:"issuer" json:"issuer"`
	IssuedDate time.Time `bson:"issued_date" json:"issued_date"`
	ExpiryDate time.Time `bson:"expiry_date" json:"expiry_date"`
	Status     string    `bson:"status" json:"status"`
}

// User represents a user in the system with role-based access control.
// Enhanced to match DATA_ARCHITECTURE.md specification with comprehensive security context.
type User struct {
	BaseModel `bson:",inline"`
	
	// Basic identification - matches email field from DATA_ARCHITECTURE.md
	Email string `bson:"email" json:"email" validate:"required,email"`
	
	// Profile information - structured to match DATA_ARCHITECTURE.md profile section
	Profile UserProfile `bson:"profile" json:"profile"`
	
	// Authentication details - enhanced security structure from DATA_ARCHITECTURE.md
	Authentication AuthenticationDetails `bson:"authentication" json:"authentication"`
	
	// Authorization - role-based access control
	Roles       []string `bson:"roles" json:"roles" validate:"required,min=1"`
	Permissions UserPermissions `bson:"permissions" json:"permissions"`
	
	// Organization context
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" validate:"required"`
	
	// Status and lifecycle
	IsActive bool   `bson:"is_active" json:"is_active"`
	Status   string `bson:"status" json:"status"` // active, inactive, suspended, locked
	
	// Metadata and preferences - matches DATA_ARCHITECTURE.md metadata section
	Metadata UserMetadata `bson:"metadata" json:"metadata"`
}

// UserProfile contains the user's personal and professional information.
// This structure matches the profile section in DATA_ARCHITECTURE.md
type UserProfile struct {
	FirstName   string `bson:"first_name" json:"first_name" validate:"required,min=1,max=100"`
	LastName    string `bson:"last_name" json:"last_name" validate:"required,min=1,max=100"`
	Title       string `bson:"title,omitempty" json:"title,omitempty"`
	Department  string `bson:"department,omitempty" json:"department,omitempty"`
	PhoneNumber string `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
}

// GetFullName returns the user's full name for display purposes.
func (up *UserProfile) GetFullName() string {
	return fmt.Sprintf("%s %s", up.FirstName, up.LastName)
}

// GetInitials returns the user's initials for avatar display.
func (up *UserProfile) GetInitials() string {
	if len(up.FirstName) == 0 || len(up.LastName) == 0 {
		return ""
	}
	return fmt.Sprintf("%c%c", up.FirstName[0], up.LastName[0])
}

// AuthenticationDetails contains comprehensive authentication and security information.
// Enhanced structure based on DATA_ARCHITECTURE.md authentication section
type AuthenticationDetails struct {
	// Password management
	PasswordHash         string    `bson:"password_hash" json:"-"`
	LastPasswordChange   time.Time `bson:"last_password_change,omitempty" json:"last_password_change,omitempty"`
	PasswordExpiresAt    time.Time `bson:"password_expires_at,omitempty" json:"password_expires_at,omitempty"`
	RequirePasswordReset bool      `bson:"require_password_reset" json:"require_password_reset"`
	
	// Session tracking
	LastLoginAt     time.Time `bson:"last_login_at,omitempty" json:"last_login_at,omitempty"`
	LastLoginIP     string    `bson:"last_login_ip,omitempty" json:"last_login_ip,omitempty"`
	CurrentSessionCount int   `bson:"current_session_count" json:"current_session_count"`
	
	// Multi-factor authentication
	MFAEnabled    bool   `bson:"mfa_enabled" json:"mfa_enabled"`
	MFASecret     string `bson:"mfa_secret,omitempty" json:"-"` // Encrypted TOTP secret
	MFABackupCodes []string `bson:"mfa_backup_codes,omitempty" json:"-"` // Encrypted backup codes
	MFAMethod     string `bson:"mfa_method,omitempty" json:"mfa_method,omitempty"` // totp, sms, email
	
	// Account security
	FailedLoginAttempts int       `bson:"failed_login_attempts" json:"failed_login_attempts"`
	LockoutUntil        time.Time `bson:"lockout_until,omitempty" json:"lockout_until,omitempty"`
	AccountLockedAt     time.Time `bson:"account_locked_at,omitempty" json:"account_locked_at,omitempty"`
	SecurityQuestions   []SecurityQuestion `bson:"security_questions,omitempty" json:"-"`
	
	// Compliance and audit
	LastSecurityReview  time.Time `bson:"last_security_review,omitempty" json:"last_security_review,omitempty"`
	ComplianceFlags     []string  `bson:"compliance_flags,omitempty" json:"compliance_flags,omitempty"`
}

// SecurityQuestion represents a user security question for account recovery.
type SecurityQuestion struct {
	Question     string `bson:"question" json:"question"`
	AnswerHash   string `bson:"answer_hash" json:"-"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}

// UserPermissions represents granular permissions for the user.
// This structure supports the detailed permission model from DATA_ARCHITECTURE.md
type UserPermissions struct {
	// Control management permissions
	CanViewControls    bool `bson:"can_view_controls" json:"can_view_controls"`
	CanEditControls    bool `bson:"can_edit_controls" json:"can_edit_controls"`
	CanDeleteControls  bool `bson:"can_delete_controls" json:"can_delete_controls"`
	CanCreateControls  bool `bson:"can_create_controls" json:"can_create_controls"`
	
	// Testing and assignment permissions
	CanAssignTests     bool `bson:"can_assign_tests" json:"can_assign_tests"`
	CanExecuteTests    bool `bson:"can_execute_tests" json:"can_execute_tests"`
	CanReviewTests     bool `bson:"can_review_tests" json:"can_review_tests"`
	
	// Evidence and findings permissions
	CanViewEvidence    bool `bson:"can_view_evidence" json:"can_view_evidence"`
	CanUploadEvidence  bool `bson:"can_upload_evidence" json:"can_upload_evidence"`
	CanApproveFindings bool `bson:"can_approve_findings" json:"can_approve_findings"`
	
	// Reporting and analytics permissions
	CanViewReports     bool `bson:"can_view_reports" json:"can_view_reports"`
	CanCreateReports   bool `bson:"can_create_reports" json:"can_create_reports"`
	CanExportData      bool `bson:"can_export_data" json:"can_export_data"`
	
	// User and organization management
	CanManageUsers     bool `bson:"can_manage_users" json:"can_manage_users"`
	CanManageSettings  bool `bson:"can_manage_settings" json:"can_manage_settings"`
	CanViewAuditLogs   bool `bson:"can_view_audit_logs" json:"can_view_audit_logs"`
	
	// Custom permissions for specific organizational needs
	CustomPermissions  map[string]bool `bson:"custom_permissions,omitempty" json:"custom_permissions,omitempty"`
}

// UserMetadata contains additional user information and preferences.
// Enhanced structure based on DATA_ARCHITECTURE.md metadata section
type UserMetadata struct {
	// System preferences
	Timezone    string `bson:"timezone" json:"timezone"`
	Language    string `bson:"language,omitempty" json:"language,omitempty"`
	DateFormat  string `bson:"date_format,omitempty" json:"date_format,omitempty"`
	
	// Application preferences
	Preferences UserPreferences `bson:"preferences" json:"preferences"`
	
	// Professional information
	EmployeeID      string    `bson:"employee_id,omitempty" json:"employee_id,omitempty"`
	HireDate        time.Time `bson:"hire_date,omitempty" json:"hire_date,omitempty"`
	ManagerID       primitive.ObjectID `bson:"manager_id,omitempty" json:"manager_id,omitempty"`
	CostCenter      string    `bson:"cost_center,omitempty" json:"cost_center,omitempty"`
	
	// Certification and training
	Certifications  []UserCertification `bson:"certifications,omitempty" json:"certifications,omitempty"`
	TrainingRecords []TrainingRecord    `bson:"training_records,omitempty" json:"training_records,omitempty"`
	
	// Contact and emergency information
	EmergencyContact EmergencyContact `bson:"emergency_contact,omitempty" json:"emergency_contact,omitempty"`
	WorkLocation     string           `bson:"work_location,omitempty" json:"work_location,omitempty"`
}

// UserPreferences contains user interface and notification preferences.
type UserPreferences struct {
	// Notification preferences
	EmailNotifications    bool `bson:"email_notifications" json:"email_notifications"`
	SMSNotifications      bool `bson:"sms_notifications" json:"sms_notifications"`
	BrowserNotifications  bool `bson:"browser_notifications" json:"browser_notifications"`
	
	// Interface preferences
	DashboardLayout       string `bson:"dashboard_layout" json:"dashboard_layout"` // compact, detailed, custom
	Theme                 string `bson:"theme,omitempty" json:"theme,omitempty"`   // light, dark, auto
	
	// Workflow preferences
	DefaultAssignmentView string `bson:"default_assignment_view,omitempty" json:"default_assignment_view,omitempty"`
	AutoSaveDrafts        bool   `bson:"auto_save_drafts" json:"auto_save_drafts"`
	ShowAdvancedFeatures  bool   `bson:"show_advanced_features" json:"show_advanced_features"`
}

// UserCertification represents professional certifications held by the user.
type UserCertification struct {
	Name           string    `bson:"name" json:"name"`
	IssuingBody    string    `bson:"issuing_body" json:"issuing_body"`
	CertificationID string   `bson:"certification_id,omitempty" json:"certification_id,omitempty"`
	IssuedDate     time.Time `bson:"issued_date" json:"issued_date"`
	ExpiryDate     time.Time `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	Status         string    `bson:"status" json:"status"` // active, expired, suspended
}

// TrainingRecord represents training completed by the user.
type TrainingRecord struct {
	Title          string    `bson:"title" json:"title"`
	Provider       string    `bson:"provider" json:"provider"`
	CompletedDate  time.Time `bson:"completed_date" json:"completed_date"`
	ExpiryDate     time.Time `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	CertificateURL string    `bson:"certificate_url,omitempty" json:"certificate_url,omitempty"`
	Credits        float64   `bson:"credits,omitempty" json:"credits,omitempty"`
}

// EmergencyContact represents emergency contact information for the user.
type EmergencyContact struct {
	Name         string `bson:"name" json:"name"`
	Relationship string `bson:"relationship" json:"relationship"`
	PhoneNumber  string `bson:"phone_number" json:"phone_number"`
	Email        string `bson:"email,omitempty" json:"email,omitempty"`
}

// ToUserProfileResponse converts a User to a UserProfileResponse for API responses.
// This provides a safe, sanitized version of user data for client consumption.
func (u *User) ToUserProfileResponse() *UserProfileResponse {
	return &UserProfileResponse{
		ID:             u.ID,
		Email:          u.Email,
		FirstName:      u.Profile.FirstName,
		LastName:       u.Profile.LastName,
		Title:          u.Profile.Title,
		Department:     u.Profile.Department,
		OrganizationID: u.OrganizationID,
		Role:           strings.Join(u.Roles, ","), // Primary role for backwards compatibility
		Permissions:    u.getPermissionsList(),
		Status:         u.Status,
		LastLogin:      u.Authentication.LastLoginAt,
		MFAEnabled:     u.Authentication.MFAEnabled,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

// getPermissionsList converts the user's permissions to a string slice for API compatibility.
func (u *User) getPermissionsList() []string {
	var permissions []string
	
	// Add permissions based on the UserPermissions structure
	if u.Permissions.CanViewControls {
		permissions = append(permissions, "controls:read")
	}
	if u.Permissions.CanEditControls {
		permissions = append(permissions, "controls:write")
	}
	if u.Permissions.CanAssignTests {
		permissions = append(permissions, "assignments:create")
	}
	if u.Permissions.CanApproveFindings {
		permissions = append(permissions, "findings:approve")
	}
	if u.Permissions.CanViewReports {
		permissions = append(permissions, "reports:read")
	}
	if u.Permissions.CanManageUsers {
		permissions = append(permissions, "users:manage")
	}
	
	// Add custom permissions
	for permission, granted := range u.Permissions.CustomPermissions {
		if granted {
			permissions = append(permissions, permission)
		}
	}
	
	return permissions
}

// IsLocked checks if the user account is currently locked.
func (u *User) IsLocked() bool {
	return !u.Authentication.LockoutUntil.IsZero() && u.Authentication.LockoutUntil.After(time.Now())
}

// IsPasswordExpired checks if the user's password has expired.
func (u *User) IsPasswordExpired() bool {
	if u.Authentication.PasswordExpiresAt.IsZero() {
		return false
	}
	return u.Authentication.PasswordExpiresAt.Before(time.Now())
}

// RequiresMFA checks if the user requires multi-factor authentication.
func (u *User) RequiresMFA() bool {
	return u.Authentication.MFAEnabled
}

// CanPerformAction checks if the user has permission to perform a specific action.
// This is a helper method for common permission checks.
func (u *User) CanPerformAction(action string) bool {
	switch action {
	case "view_controls":
		return u.Permissions.CanViewControls
	case "edit_controls":
		return u.Permissions.CanEditControls
	case "assign_tests":
		return u.Permissions.CanAssignTests
	case "approve_findings":
		return u.Permissions.CanApproveFindings
	case "manage_users":
		return u.Permissions.CanManageUsers
	case "view_reports":
		return u.Permissions.CanViewReports
	default:
		// Check custom permissions
		if granted, exists := u.Permissions.CustomPermissions[action]; exists {
			return granted
		}
		return false
	}
}

// UserProfileResponse represents the user information included in authentication responses.
// This is a sanitized version of the User model safe for client consumption.
type UserProfileResponse struct {
	ID             primitive.ObjectID `json:"id"`
	Email          string             `json:"email"`
	FirstName      string             `json:"first_name"`
	LastName       string             `json:"last_name"`
	Title          string             `json:"title,omitempty"`
	Department     string             `json:"department,omitempty"`
	OrganizationID primitive.ObjectID `json:"organization_id"`
	Role           string             `json:"role"`
	Permissions    []string           `json:"permissions,omitempty"`
	Status         string             `json:"status"`
	LastLogin      time.Time          `json:"last_login,omitempty"`
	MFAEnabled     bool               `json:"mfa_enabled"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

// Control represents a compliance control that needs to be tested.
type Control struct {
	BaseModel `bson:",inline"`
	
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" validate:"required"`
	
	// Control identification
	ControlID   string `bson:"control_id" json:"control_id" validate:"required"`
	Title       string `bson:"title" json:"title" validate:"required,min=1,max=500"`
	Description string `bson:"description" json:"description" validate:"required"`
	
	// Framework and classification
	Framework string `bson:"framework" json:"framework" validate:"required"`
	Category  string `bson:"category" json:"category" validate:"required"`
	SubCategory string `bson:"sub_category,omitempty" json:"sub_category,omitempty"`
	
	// Risk and importance
	RiskLevel  string `bson:"risk_level" json:"risk_level" validate:"required"`
	Importance string `bson:"importance" json:"importance" validate:"required"`
	
	// Control details
	ControlType      string   `bson:"control_type" json:"control_type" validate:"required"`
	ControlFrequency string   `bson:"control_frequency" json:"control_frequency" validate:"required"`
	Owner            string   `bson:"owner" json:"owner" validate:"required"`
	Process          string   `bson:"process" json:"process"`
	Systems          []string `bson:"systems,omitempty" json:"systems,omitempty"`
	
	// Testing requirements
	TestingProcedure string                 `bson:"testing_procedure" json:"testing_procedure" validate:"required"`
	SampleSize       int                    `bson:"sample_size" json:"sample_size"`
	EvidenceTypes    []string               `bson:"evidence_types" json:"evidence_types"`
	TestingNotes     string                 `bson:"testing_notes,omitempty" json:"testing_notes,omitempty"`
	
	// Status and metadata
	Status      string                 `bson:"status" json:"status"`
	Tags        []string               `bson:"tags,omitempty" json:"tags,omitempty"`
	CustomFields map[string]interface{} `bson:"custom_fields,omitempty" json:"custom_fields,omitempty"`
}

// TestingCycle represents a period during which controls are tested.
type TestingCycle struct {
	BaseModel `bson:",inline"`
	
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" validate:"required"`
	
	// Cycle identification
	CycleID     string `bson:"cycle_id" json:"cycle_id" validate:"required"`
	Name        string `bson:"name" json:"name" validate:"required,min=1,max=255"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	
	// Timing
	StartDate time.Time `bson:"start_date" json:"start_date" validate:"required"`
	EndDate   time.Time `bson:"end_date" json:"end_date" validate:"required"`
	
	// Scope and configuration
	ControlScope    []primitive.ObjectID `bson:"control_scope" json:"control_scope"`
	TestingType     string               `bson:"testing_type" json:"testing_type" validate:"required"`
	Framework       string               `bson:"framework" json:"framework"`
	
	// Status and progress
	Status       string    `bson:"status" json:"status"`
	Progress     Progress  `bson:"progress" json:"progress"`
	CompletedAt  time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	
	// Settings
	Settings map[string]interface{} `bson:"settings,omitempty" json:"settings,omitempty"`
}

// Progress tracks the completion status of a testing cycle.
type Progress struct {
	TotalControls     int `bson:"total_controls" json:"total_controls"`
	CompletedControls int `bson:"completed_controls" json:"completed_controls"`
	FailedControls    int `bson:"failed_controls" json:"failed_controls"`
	PercentComplete   int `bson:"percent_complete" json:"percent_complete"`
}

// EvidenceRequest represents a request for evidence from a control owner.
type EvidenceRequest struct {
	BaseModel `bson:",inline"`
	
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" validate:"required"`
	ControlID      primitive.ObjectID `bson:"control_id" json:"control_id" validate:"required"`
	CycleID        primitive.ObjectID `bson:"cycle_id" json:"cycle_id" validate:"required"`
	
	// Request details
	RequestID   string `bson:"request_id" json:"request_id" validate:"required"`
	Title       string `bson:"title" json:"title" validate:"required,min=1,max=255"`
	Description string `bson:"description" json:"description" validate:"required"`
	
	// Assignment
	AssigneeID   primitive.ObjectID `bson:"assignee_id" json:"assignee_id" validate:"required"`
	AssignerID   primitive.ObjectID `bson:"assigner_id" json:"assigner_id" validate:"required"`
	AssignedDate time.Time          `bson:"assigned_date" json:"assigned_date"`
	
	// Timing
	DueDate     time.Time `bson:"due_date" json:"due_date" validate:"required"`
	CompletedAt time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	
	// Requirements
	EvidenceTypes []string `bson:"evidence_types" json:"evidence_types"`
	SampleSize    int      `bson:"sample_size" json:"sample_size"`
	Instructions  string   `bson:"instructions,omitempty" json:"instructions,omitempty"`
	
	// Status and responses
	Status   string     `bson:"status" json:"status"`
	Response string     `bson:"response,omitempty" json:"response,omitempty"`
	Evidence []Evidence `bson:"evidence,omitempty" json:"evidence,omitempty"`
	
	// Communication
	Comments []Comment `bson:"comments,omitempty" json:"comments,omitempty"`
}

// Evidence represents uploaded evidence files and metadata.
type Evidence struct {
	ID          string    `bson:"id" json:"id"`
	FileName    string    `bson:"file_name" json:"file_name"`
	FileSize    int64     `bson:"file_size" json:"file_size"`
	FileType    string    `bson:"file_type" json:"file_type"`
	StoragePath string    `bson:"storage_path" json:"storage_path"`
	UploadedAt  time.Time `bson:"uploaded_at" json:"uploaded_at"`
	UploadedBy  string    `bson:"uploaded_by" json:"uploaded_by"`
	
	// Metadata
	Description string                 `bson:"description,omitempty" json:"description,omitempty"`
	Tags        []string               `bson:"tags,omitempty" json:"tags,omitempty"`
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	
	// Processing status
	ProcessingStatus string `bson:"processing_status" json:"processing_status"`
	TextExtracted    string `bson:"text_extracted,omitempty" json:"text_extracted,omitempty"`
}

// Comment represents a comment on an evidence request or other entity.
type Comment struct {
	ID        string    `bson:"id" json:"id"`
	AuthorID  string    `bson:"author_id" json:"author_id"`
	Content   string    `bson:"content" json:"content"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	
	// Threading
	ParentID string `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	
	// Status
	Status string `bson:"status" json:"status"`
}

// AuditLog represents an audit trail entry for compliance tracking.
type AuditLog struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp    time.Time          `bson:"timestamp" json:"timestamp"`
	
	// Context
	CorrelationID  string             `bson:"correlation_id,omitempty" json:"correlation_id,omitempty"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	
	// Action details
	Action       string `bson:"action" json:"action"`
	ResourceType string `bson:"resource_type" json:"resource_type"`
	ResourceID   string `bson:"resource_id" json:"resource_id"`
	
	// Change tracking
	OldValues map[string]interface{} `bson:"old_values,omitempty" json:"old_values,omitempty"`
	NewValues map[string]interface{} `bson:"new_values,omitempty" json:"new_values,omitempty"`
	
	// Request metadata
	IPAddress string                 `bson:"ip_address,omitempty" json:"ip_address,omitempty"`
	UserAgent string                 `bson:"user_agent,omitempty" json:"user_agent,omitempty"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	
	// Result
	Success     bool   `bson:"success" json:"success"`
	ErrorMessage string `bson:"error_message,omitempty" json:"error_message,omitempty"`
}

// Common status constants
const (
	// User statuses
	UserStatusActive    = "active"
	UserStatusInactive  = "inactive"
	UserStatusSuspended = "suspended"
	
	// Organization statuses
	OrganizationStatusActive    = "active"
	OrganizationStatusInactive  = "inactive"
	OrganizationStatusSuspended = "suspended"
	
	// Control statuses
	ControlStatusDraft    = "draft"
	ControlStatusActive   = "active"
	ControlStatusInactive = "inactive"
	ControlStatusArchived = "archived"
	
	// Testing cycle statuses
	CycleStatusPlanning   = "planning"
	CycleStatusActive     = "active"
	CycleStatusCompleted  = "completed"
	CycleStatusCancelled  = "cancelled"
	
	// Evidence request statuses
	EvidenceRequestStatusPending    = "pending"
	EvidenceRequestStatusInProgress = "in_progress"
	EvidenceRequestStatusCompleted  = "completed"
	EvidenceRequestStatusOverdue    = "overdue"
	EvidenceRequestStatusCancelled  = "cancelled"
	
	// Common roles
	RoleAdmin     = "admin"
	RoleManager   = "manager"
	RoleAuditor   = "auditor"
	RoleOwner     = "owner"
	RoleViewer    = "viewer"
)

// NewID generates a new UUID string for entity IDs
func NewID() string {
	return uuid.New().String()
}

// UpdateTimestamps updates the created_at and updated_at timestamps
func (b *BaseModel) UpdateTimestamps() {
	now := time.Now()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	b.UpdatedAt = now
}