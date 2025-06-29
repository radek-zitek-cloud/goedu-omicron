// Package models contains the core domain models for the GoEdu Control Testing Platform.
// These models represent the fundamental entities in the control testing domain,
// designed for financial compliance and audit requirements.
package models

import (
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
type User struct {
	BaseModel `bson:",inline"`
	
	// Personal information
	FirstName string `bson:"first_name" json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `bson:"last_name" json:"last_name" validate:"required,min=1,max=100"`
	Email     string `bson:"email" json:"email" validate:"required,email"`
	
	// Authentication
	PasswordHash string    `bson:"password_hash" json:"-"`
	LastLogin    time.Time `bson:"last_login,omitempty" json:"last_login,omitempty"`
	
	// Organization and role
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" validate:"required"`
	Role           string             `bson:"role" json:"role" validate:"required"`
	Permissions    []string           `bson:"permissions,omitempty" json:"permissions,omitempty"`
	
	// Profile information
	Title      string `bson:"title,omitempty" json:"title,omitempty"`
	Department string `bson:"department,omitempty" json:"department,omitempty"`
	Phone      string `bson:"phone,omitempty" json:"phone,omitempty"`
	
	// Status and preferences
	Status      string                 `bson:"status" json:"status"`
	Preferences map[string]interface{} `bson:"preferences,omitempty" json:"preferences,omitempty"`
	
	// Security settings
	MFAEnabled    bool      `bson:"mfa_enabled" json:"mfa_enabled"`
	MFASecret     string    `bson:"mfa_secret,omitempty" json:"-"`
	FailedLogins  int       `bson:"failed_logins" json:"failed_logins"`
	LockedUntil   time.Time `bson:"locked_until,omitempty" json:"locked_until,omitempty"`
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