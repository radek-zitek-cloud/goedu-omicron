// Package models contains authentication-related data structures for the GoEdu Control Testing Platform.
// These models handle user authentication, authorization, and security contexts
// required for financial compliance and audit environments.
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LoginRequest represents the request payload for user authentication.
// This structure handles both standard email/password login and MFA scenarios.
type LoginRequest struct {
	// Basic authentication credentials
	Email    string `json:"email" validate:"required,email" binding:"required,email"`
	Password string `json:"password" validate:"required,min=12,max=128" binding:"required,min=12,max=128"`
	
	// Multi-factor authentication
	MFACode string `json:"mfa_code,omitempty" binding:"omitempty,len=6,numeric"`
	
	// Session management
	RememberMe bool `json:"remember_me,omitempty"`
	
	// Security context
	IPAddress string `json:"-"` // Set by middleware, not from request body
	UserAgent string `json:"-"` // Set by middleware, not from request body
}

// LoginResponse represents the response payload after successful authentication.
// Contains user information, tokens, and security context required for the frontend.
type LoginResponse struct {
	// Authentication result
	Success      bool   `json:"success"`
	Message      string `json:"message,omitempty"`
	RequiresMFA  bool   `json:"requires_mfa,omitempty"`
	
	// User context
	User *UserProfileResponse `json:"user,omitempty"`
	
	// Token management
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	
	// Session information
	SessionID string `json:"session_id,omitempty"`
}

// RefreshTokenRequest represents a request to refresh an access token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" binding:"required"`
}

// Permission represents a granular permission in the RBAC system.
// Permissions follow the pattern: resource:action:scope
type Permission struct {
	// Permission identification
	ID          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name" validate:"required"`
	Description string `json:"description" bson:"description"`
	
	// Permission structure following resource:action:scope pattern
	Resource string `json:"resource" bson:"resource" validate:"required"` // e.g., "controls", "testing_cycles"
	Action   string `json:"action" bson:"action" validate:"required"`     // e.g., "read", "write", "delete"
	Scope    string `json:"scope" bson:"scope" validate:"required"`       // e.g., "own", "team", "organization"
	
	// Metadata
	Category string `json:"category,omitempty" bson:"category,omitempty"`
	Tags     []string `json:"tags,omitempty" bson:"tags,omitempty"`
	
	// Status
	IsActive  bool      `json:"is_active" bson:"is_active"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Role represents a collection of permissions for role-based access control.
// Roles are designed to match typical banking audit organizational structures.
type Role struct {
	// Role identification
	ID          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name" validate:"required"`
	Description string `json:"description" bson:"description"`
	
	// Permission management
	Permissions []Permission `json:"permissions" bson:"permissions"`
	
	// Role hierarchy and inheritance
	ParentRoleID string   `json:"parent_role_id,omitempty" bson:"parent_role_id,omitempty"`
	ChildRoles   []string `json:"child_roles,omitempty" bson:"child_roles,omitempty"`
	
	// Organization context
	OrganizationID primitive.ObjectID `json:"organization_id,omitempty" bson:"organization_id,omitempty"`
	IsSystemRole   bool               `json:"is_system_role" bson:"is_system_role"` // Built-in vs custom roles
	
	// Status and metadata
	IsActive  bool                   `json:"is_active" bson:"is_active"`
	Priority  int                    `json:"priority" bson:"priority"` // For role precedence
	Metadata  map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" bson:"updated_at"`
	CreatedBy string                 `json:"created_by,omitempty" bson:"created_by,omitempty"`
	UpdatedBy string                 `json:"updated_by,omitempty" bson:"updated_by,omitempty"`
}

// JWTClaims represents the claims structure for JWT tokens.
// Based on the JWT structure defined in SYSTEM_ARCHITECTURE.md
type JWTClaims struct {
	// Standard JWT claims
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	Audience  string `json:"aud"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	NotBefore int64  `json:"nbf,omitempty"`
	JWTID     string `json:"jti,omitempty"`
	
	// Custom claims for GoEdu platform
	UserID         string   `json:"user_id"`
	Email          string   `json:"email"`
	Roles          []string `json:"roles"`
	OrganizationID string   `json:"organization_id"`
	Permissions    []string `json:"permissions"`
	SessionID      string   `json:"session_id"`
	
	// Security context
	TokenType string `json:"token_type"` // "access" or "refresh"
	IPAddress string `json:"ip_address,omitempty"`
}

// Session represents an active user session for tracking and management.
type Session struct {
	BaseModel `bson:",inline"`
	
	// Session identification
	SessionID string             `bson:"session_id" json:"session_id" validate:"required"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id" validate:"required"`
	
	// Session security
	IPAddress     string    `bson:"ip_address" json:"ip_address"`
	UserAgent     string    `bson:"user_agent" json:"user_agent"`
	RefreshToken  string    `bson:"refresh_token" json:"-"` // Hashed refresh token
	LastActivity  time.Time `bson:"last_activity" json:"last_activity"`
	
	// Session management
	ExpiresAt time.Time `bson:"expires_at" json:"expires_at"`
	IsActive  bool      `bson:"is_active" json:"is_active"`
	
	// Tracking
	LoginMethod string                 `bson:"login_method" json:"login_method"` // "password", "sso", "mfa"
	DeviceInfo  map[string]interface{} `bson:"device_info,omitempty" json:"device_info,omitempty"`
}

// AuditEvent represents security and authentication audit events.
type AuditEvent struct {
	BaseModel `bson:",inline"`
	
	// Event identification
	EventID   string             `bson:"event_id" json:"event_id" validate:"required"`
	EventType string             `bson:"event_type" json:"event_type" validate:"required"`
	
	// Context
	UserID         primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	OrganizationID primitive.ObjectID `bson:"organization_id,omitempty" json:"organization_id,omitempty"`
	SessionID      string             `bson:"session_id,omitempty" json:"session_id,omitempty"`
	
	// Event details
	Action      string `bson:"action" json:"action" validate:"required"`
	Resource    string `bson:"resource,omitempty" json:"resource,omitempty"`
	ResourceID  string `bson:"resource_id,omitempty" json:"resource_id,omitempty"`
	Description string `bson:"description" json:"description"`
	
	// Security context
	IPAddress  string `bson:"ip_address,omitempty" json:"ip_address,omitempty"`
	UserAgent  string `bson:"user_agent,omitempty" json:"user_agent,omitempty"`
	
	// Result and metadata
	Success   bool                   `bson:"success" json:"success"`
	ErrorCode string                 `bson:"error_code,omitempty" json:"error_code,omitempty"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	
	// Risk assessment
	RiskLevel string `bson:"risk_level,omitempty" json:"risk_level,omitempty"`
	Severity  string `bson:"severity,omitempty" json:"severity,omitempty"`
}

// Built-in roles for banking audit workflows
// These match the role definitions in SYSTEM_ARCHITECTURE.md
var DefaultRoles = map[string]Role{
	RoleAuditor: {
		ID:          RoleAuditor,
		Name:        "Auditor",
		Description: "Standard auditor role with read access to controls and ability to manage own assignments",
		Permissions: []Permission{
			{Resource: "controls", Action: "read", Scope: "organization"},
			{Resource: "assignments", Action: "read", Scope: "own"},
			{Resource: "assignments", Action: "update", Scope: "own"},
			{Resource: "evidence_requests", Action: "create", Scope: "own"},
			{Resource: "test_executions", Action: "create", Scope: "own"},
			{Resource: "findings", Action: "create", Scope: "own"},
		},
		IsSystemRole: true,
		IsActive:     true,
		Priority:     100,
	},
	"audit_manager": {
		ID:          "audit_manager",
		Name:        "Audit Manager",
		Description: "Audit manager role with control management and team assignment capabilities",
		Permissions: []Permission{
			{Resource: "controls", Action: "read", Scope: "organization"},
			{Resource: "controls", Action: "write", Scope: "organization"},
			{Resource: "testing_cycles", Action: "create", Scope: "team"},
			{Resource: "assignments", Action: "create", Scope: "team"},
			{Resource: "assignments", Action: "read", Scope: "team"},
			{Resource: "assignments", Action: "update", Scope: "team"},
			{Resource: "reports", Action: "read", Scope: "team"},
			{Resource: "findings", Action: "approve", Scope: "team"},
		},
		IsSystemRole: true,
		IsActive:     true,
		Priority:     200,
	},
	RoleAdmin: {
		ID:          RoleAdmin,
		Name:        "Administrator",
		Description: "System administrator with full access to all resources",
		Permissions: []Permission{
			{Resource: "*", Action: "*", Scope: "*"},
		},
		IsSystemRole: true,
		IsActive:     true,
		Priority:     1000,
	},
}

// Authentication constants
const (
	// Token types
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
	
	// Login methods
	LoginMethodPassword = "password"
	LoginMethodSSO      = "sso"
	LoginMethodMFA      = "mfa"
	
	// Event types for audit logging
	EventTypeLogin           = "login"
	EventTypeLogout          = "logout"
	EventTypeTokenRefresh    = "token_refresh"
	EventTypePasswordChange  = "password_change"
	EventTypeAccountLocked   = "account_locked"
	EventTypePermissionDenied = "permission_denied"
	
	// Risk levels
	RiskLevelLow      = "low"
	RiskLevelMedium   = "medium"
	RiskLevelHigh     = "high"
	RiskLevelCritical = "critical"
	
	// Default token expiration times
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour // 7 days
	SessionDuration      = 8 * time.Hour      // 8 hours
)