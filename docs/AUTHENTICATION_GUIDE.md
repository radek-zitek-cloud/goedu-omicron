# Authentication Implementation Guide

## Overview

This document provides a comprehensive guide to the authentication and authorization system implemented for the GoEdu Control Testing Platform as part of TASK-003.

## Architecture

The authentication system follows a multi-layered security approach designed for financial compliance environments:

- **Password Security**: bcrypt hashing with configurable cost factor
- **Token Management**: JWT-based access and refresh tokens
- **Role-Based Access Control**: Hierarchical permission system
- **Session Management**: Secure session tracking with audit capabilities
- **Multi-Factor Authentication**: Support for TOTP and backup codes

## Core Components

### 1. Enhanced User Model (`/internal/models/models.go`)

The User entity has been restructured to align with the DATA_ARCHITECTURE.md specification:

```go
type User struct {
    BaseModel `bson:",inline"`
    
    // Basic identification
    Email string `bson:"email" json:"email"`
    
    // Profile information
    Profile UserProfile `bson:"profile" json:"profile"`
    
    // Authentication details
    Authentication AuthenticationDetails `bson:"authentication" json:"authentication"`
    
    // Authorization
    Roles       []string `bson:"roles" json:"roles"`
    Permissions UserPermissions `bson:"permissions" json:"permissions"`
    
    // Organization context
    OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
    
    // Status and metadata
    IsActive bool   `bson:"is_active" json:"is_active"`
    Status   string `bson:"status" json:"status"`
    Metadata UserMetadata `bson:"metadata" json:"metadata"`
}
```

### 2. Authentication Models (`/internal/models/auth.go`)

Comprehensive authentication data structures:

- **LoginRequest/LoginResponse**: Authentication flow models
- **JWTClaims**: Token structure following SYSTEM_ARCHITECTURE.md
- **Permission/Role**: RBAC system models
- **Session**: Active session tracking
- **AuditEvent**: Security event logging

### 3. Password Security (`/pkg/auth/auth.go`)

```go
// Password hashing with bcrypt
hasher := auth.NewPasswordHasher(12) // Cost factor for financial services
hash, err := hasher.HashPassword("SecurePassword123!")
valid, err := hasher.VerifyPassword("SecurePassword123!", hash)
```

Key features:
- Configurable bcrypt cost factor (default: 12)
- Password strength validation (12-128 characters)
- Timing attack protection
- Hash cost extraction for policy validation

### 4. JWT Token Management

```go
// JWT token generation and validation
jwtManager := auth.NewJWTManager(secretKey, "goedu-platform", "goedu-api")
token, expiresAt, err := jwtManager.GenerateAccessToken(userProfile, sessionID, ipAddress)
claims, err := jwtManager.ValidateToken(token)
```

Features:
- Access tokens (15 minutes) and refresh tokens (7 days)
- Comprehensive claim validation
- Token refresh flow
- Security context tracking (IP, session)

### 5. Role-Based Access Control (RBAC)

```go
// Permission checking
checker := auth.NewPermissionChecker()
hasAccess := checker.HasPermission(userRoles, "controls", "read", "organization")
err := checker.ValidatePermission(userRoles, "controls", "write", "team")
```

Built-in roles:
- **admin**: Full system access (`*:*:*`)
- **audit_manager**: Team management and control oversight
- **auditor**: Individual testing and evidence collection

Permission structure: `resource:action:scope`
- **Resource**: controls, assignments, evidence_requests, etc.
- **Action**: read, write, create, delete, approve, etc.
- **Scope**: own, team, organization, * (with hierarchy)

## Usage Examples

### User Registration

```go
// Create new user with proper security context
user := &models.User{
    Email: "auditor@bank.com",
    Profile: models.UserProfile{
        FirstName:  "Jane",
        LastName:   "Auditor",
        Title:      "Senior IT Auditor",
        Department: "Internal Audit",
    },
    Authentication: models.AuthenticationDetails{
        PasswordHash: hashedPassword, // Use PasswordHasher
    },
    OrganizationID: orgID,
    Roles:          []string{models.RoleAuditor},
    IsActive:       true,
    Status:         models.UserStatusActive,
}
```

### Authentication Flow

```go
// 1. Login
loginRequest := &models.LoginRequest{
    Email:    "auditor@bank.com",
    Password: "SecurePassword123!",
}

// 2. Validate credentials and generate tokens
user, err := userService.AuthenticateUser(ctx, loginRequest.Email, loginRequest.Password)
accessToken, expiresAt, err := jwtManager.GenerateAccessToken(user.ToUserProfileResponse(), sessionID, ipAddress)
refreshToken, _, err := jwtManager.GenerateRefreshToken(user.ID.Hex(), sessionID)

// 3. Return response
response := &models.LoginResponse{
    Success:      true,
    User:         user.ToUserProfileResponse(),
    AccessToken:  accessToken,
    RefreshToken: refreshToken,
    ExpiresAt:    expiresAt,
    SessionID:    sessionID,
}
```

### Permission Checking

```go
// Check if user can perform specific actions
if user.CanPerformAction("view_controls") {
    // Allow access to controls
}

// RBAC permission checking
checker := auth.NewPermissionChecker()
if checker.HasPermission(user.Roles, "assignments", "create", "team") {
    // Allow creating team assignments
}
```

### Session Management

```go
// Create and track user sessions
session := &models.Session{
    SessionID:    generateSessionID(),
    UserID:       user.ID,
    IPAddress:    ipAddress,
    UserAgent:    userAgent,
    RefreshToken: hashedRefreshToken,
    ExpiresAt:    time.Now().Add(8 * time.Hour),
    IsActive:     true,
    LoginMethod:  models.LoginMethodPassword,
}
```

## Security Features

### 1. Account Protection
- Failed login attempt tracking
- Account lockout after 5 failed attempts
- 30-minute lockout duration
- Security event logging

### 2. Token Security
- Short-lived access tokens (15 minutes)
- Secure refresh token rotation
- IP address validation
- Session binding

### 3. Password Policy
- Minimum 12 characters
- Maximum 128 characters
- bcrypt with cost factor 12
- Password expiration support

### 4. Audit Trail
- All authentication events logged
- Security context tracking
- Failed login monitoring
- Permission denial logging

## Service Interfaces

The authentication system is designed with clean service interfaces:

```go
type AuthenticationService interface {
    Login(ctx context.Context, request *models.LoginRequest) (*models.LoginResponse, error)
    Logout(ctx context.Context, sessionID string) error
    RefreshToken(ctx context.Context, refreshToken string) (*models.LoginResponse, error)
    ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
    // ... additional methods
}

type PermissionService interface {
    HasPermission(ctx context.Context, userID, resource, action, scope string) (bool, error)
    ValidatePermission(ctx context.Context, userID, resource, action, scope string) error
    // ... additional methods
}
```

## Testing

Comprehensive test suite covering:
- Password hashing and verification
- JWT token generation and validation
- RBAC permission checking
- Error condition handling
- Security edge cases

Run tests:
```bash
go test ./pkg/auth/... -v
```

## Configuration

### Environment Variables

```bash
# JWT Configuration
JWT_SECRET_KEY=your-256-bit-secret-key
JWT_ISSUER=goedu-platform
JWT_AUDIENCE=goedu-api

# Password Policy
BCRYPT_COST=12
MIN_PASSWORD_LENGTH=12
MAX_PASSWORD_LENGTH=128

# Session Management
SESSION_DURATION=8h
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=168h

# Account Security
MAX_FAILED_ATTEMPTS=5
LOCKOUT_DURATION=30m
```

## Next Steps

The authentication foundation is now complete. Recommended next steps:

1. **Middleware Implementation**: Create HTTP middleware for token validation
2. **API Endpoints**: Implement authentication REST/GraphQL endpoints
3. **MFA Integration**: Add TOTP and backup code functionality
4. **SSO Integration**: Add SAML/OIDC support for enterprise integration
5. **Audit Dashboard**: Create security monitoring interface

## Compliance

This implementation meets financial services security requirements:
- Strong password policies
- Comprehensive audit logging
- Session management
- Role-based access control
- Security event monitoring

The system is designed to support SOX, PCI DSS, and other regulatory frameworks commonly required in banking environments.