// Package auth provides authentication utilities for the GoEdu Control Testing Platform.
// This package implements secure password hashing with bcrypt and JWT token management
// following banking security standards and compliance requirements.
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/models"
)

// Password hashing configuration following security best practices
const (
	// bcrypt cost factor - 12 provides good security/performance balance for financial services
	BcryptCost = 12
	
	// Minimum password requirements
	MinPasswordLength = 12
	MaxPasswordLength = 128
	
	// Maximum failed login attempts before lockout
	MaxFailedAttempts = 5
	
	// Account lockout duration
	LockoutDuration = 30 * time.Minute
)

// Errors for authentication operations
var (
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrAccountLocked        = errors.New("account temporarily locked due to too many failed attempts")
	ErrAccountInactive      = errors.New("account is inactive")
	ErrPasswordTooWeak      = errors.New("password does not meet security requirements")
	ErrInvalidToken         = errors.New("invalid or expired token")
	ErrTokenExpired         = errors.New("token has expired")
	ErrInsufficientPermissions = errors.New("insufficient permissions for this operation")
)

// PasswordHasher provides secure password hashing functionality using bcrypt.
type PasswordHasher struct {
	cost int
}

// NewPasswordHasher creates a new password hasher with the specified bcrypt cost.
// If cost is 0, uses the default BcryptCost value.
//
// Cost values:
//   - 10-12: Good for most applications
//   - 12-14: Recommended for financial/security applications
//   - 14+: High security, slower performance
func NewPasswordHasher(cost int) *PasswordHasher {
	if cost == 0 {
		cost = BcryptCost
	}
	return &PasswordHasher{cost: cost}
}

// HashPassword creates a bcrypt hash from a plaintext password.
// Uses a random salt and the configured cost factor for security.
//
// Parameters:
//   - password: plaintext password to hash (min 12, max 128 characters)
//
// Returns:
//   - Bcrypt hash string safe for database storage
//   - Error if password validation fails or hashing fails
//
// Security considerations:
//   - Uses cryptographically secure random salt
//   - Cost factor ensures adequate resistance to brute force
//   - Hash includes salt, cost, and algorithm version
func (ph *PasswordHasher) HashPassword(password string) (string, error) {
	// Validate password length
	if len(password) < MinPasswordLength {
		return "", fmt.Errorf("%w: minimum length %d characters", ErrPasswordTooWeak, MinPasswordLength)
	}
	if len(password) > MaxPasswordLength {
		return "", fmt.Errorf("%w: maximum length %d characters", ErrPasswordTooWeak, MaxPasswordLength)
	}
	
	// Generate bcrypt hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), ph.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	
	return string(hash), nil
}

// VerifyPassword compares a plaintext password against a bcrypt hash.
// Implements constant-time comparison to prevent timing attacks.
//
// Parameters:
//   - password: plaintext password to verify
//   - hash: bcrypt hash from database
//
// Returns:
//   - true if password matches hash
//   - error if comparison fails or hash is invalid
//
// Security considerations:
//   - Uses constant-time comparison
//   - Validates hash format before comparison
//   - Resistant to timing attacks
func (ph *PasswordHasher) VerifyPassword(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, fmt.Errorf("password verification failed: %w", err)
	}
	return true, nil
}

// GetHashCost extracts the cost factor from a bcrypt hash.
// Useful for password policy validation and hash migration.
func (ph *PasswordHasher) GetHashCost(hash string) (int, error) {
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		return 0, fmt.Errorf("failed to extract hash cost: %w", err)
	}
	return cost, nil
}

// JWTManager handles JWT token creation, validation, and management.
// Implements the JWT structure defined in SYSTEM_ARCHITECTURE.md
type JWTManager struct {
	secretKey       []byte
	issuer          string
	audience        string
	accessDuration  time.Duration
	refreshDuration time.Duration
}

// NewJWTManager creates a new JWT manager with the specified configuration.
//
// Parameters:
//   - secretKey: secret key for signing tokens (should be 256-bit for HS256)
//   - issuer: token issuer identifier (e.g., "goedu-platform")
//   - audience: token audience identifier (e.g., "goedu-api")
//
// Security considerations:
//   - Secret key should be cryptographically random
//   - Key should be rotated periodically
//   - Consider using asymmetric keys (RS256) for distributed systems
func NewJWTManager(secretKey []byte, issuer, audience string) *JWTManager {
	return &JWTManager{
		secretKey:       secretKey,
		issuer:          issuer,
		audience:        audience,
		accessDuration:  models.AccessTokenDuration,
		refreshDuration: models.RefreshTokenDuration,
	}
}

// GenerateAccessToken creates a JWT access token for the specified user.
// Includes user context, roles, and permissions in the token claims.
//
// Parameters:
//   - user: user profile for token generation
//   - sessionID: unique session identifier
//   - ipAddress: client IP address for security context
//
// Returns:
//   - Signed JWT token string
//   - Token expiration time
//   - Error if token generation fails
func (jm *JWTManager) GenerateAccessToken(user *models.UserProfileResponse, sessionID, ipAddress string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(jm.accessDuration)
	
	claims := models.JWTClaims{
		Issuer:         jm.issuer,
		Subject:        user.ID.Hex(),
		Audience:       jm.audience,
		ExpiresAt:      expiresAt.Unix(),
		IssuedAt:       now.Unix(),
		NotBefore:      now.Unix(),
		JWTID:          generateJTI(),
		UserID:         user.ID.Hex(),
		Email:          user.Email,
		Roles:          []string{user.Role},
		OrganizationID: user.OrganizationID.Hex(),
		Permissions:    user.Permissions,
		SessionID:      sessionID,
		TokenType:      models.TokenTypeAccess,
		IPAddress:      ipAddress,
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":             claims.Issuer,
		"sub":             claims.Subject,
		"aud":             claims.Audience,
		"exp":             claims.ExpiresAt,
		"iat":             claims.IssuedAt,
		"nbf":             claims.NotBefore,
		"jti":             claims.JWTID,
		"user_id":         claims.UserID,
		"email":           claims.Email,
		"roles":           claims.Roles,
		"organization_id": claims.OrganizationID,
		"permissions":     claims.Permissions,
		"session_id":      claims.SessionID,
		"token_type":      claims.TokenType,
		"ip_address":      claims.IPAddress,
	})
	
	tokenString, err := token.SignedString(jm.secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign access token: %w", err)
	}
	
	return tokenString, expiresAt, nil
}

// GenerateRefreshToken creates a JWT refresh token for session management.
// Refresh tokens have longer expiration and are used to obtain new access tokens.
//
// Parameters:
//   - userID: user identifier
//   - sessionID: unique session identifier
//
// Returns:
//   - Signed JWT refresh token string
//   - Token expiration time
//   - Error if token generation fails
func (jm *JWTManager) GenerateRefreshToken(userID, sessionID string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(jm.refreshDuration)
	
	claims := jwt.MapClaims{
		"iss":        jm.issuer,
		"sub":        userID,
		"aud":        jm.audience,
		"exp":        expiresAt.Unix(),
		"iat":        now.Unix(),
		"jti":        generateJTI(),
		"session_id": sessionID,
		"token_type": models.TokenTypeRefresh,
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jm.secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}
	
	return tokenString, expiresAt, nil
}

// ValidateToken validates a JWT token and extracts claims.
// Performs signature validation, expiration checking, and audience verification.
//
// Parameters:
//   - tokenString: JWT token string to validate
//
// Returns:
//   - JWT claims if token is valid
//   - Error if token is invalid, expired, or malformed
//
// Security considerations:
//   - Validates signature against secret key
//   - Checks token expiration
//   - Verifies issuer and audience claims
//   - Resistant to common JWT attacks
func (jm *JWTManager) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jm.secretKey, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	
	// Validate required claims
	if issuer, ok := claims["iss"].(string); !ok || issuer != jm.issuer {
		return nil, fmt.Errorf("%w: invalid issuer", ErrInvalidToken)
	}
	
	if audience, ok := claims["aud"].(string); !ok || audience != jm.audience {
		return nil, fmt.Errorf("%w: invalid audience", ErrInvalidToken)
	}
	
	// Extract claims
	jwtClaims := &models.JWTClaims{
		Issuer:    getStringClaim(claims, "iss"),
		Subject:   getStringClaim(claims, "sub"),
		Audience:  getStringClaim(claims, "aud"),
		ExpiresAt: getInt64Claim(claims, "exp"),
		IssuedAt:  getInt64Claim(claims, "iat"),
		NotBefore: getInt64Claim(claims, "nbf"),
		JWTID:     getStringClaim(claims, "jti"),
		
		UserID:         getStringClaim(claims, "user_id"),
		Email:          getStringClaim(claims, "email"),
		Roles:          getStringSliceClaim(claims, "roles"),
		OrganizationID: getStringClaim(claims, "organization_id"),
		Permissions:    getStringSliceClaim(claims, "permissions"),
		SessionID:      getStringClaim(claims, "session_id"),
		TokenType:      getStringClaim(claims, "token_type"),
		IPAddress:      getStringClaim(claims, "ip_address"),
	}
	
	// Check expiration
	if time.Unix(jwtClaims.ExpiresAt, 0).Before(time.Now()) {
		return nil, ErrTokenExpired
	}
	
	return jwtClaims, nil
}

// RefreshAccessToken generates a new access token using a valid refresh token.
// This implements the token refresh flow for seamless user experience.
func (jm *JWTManager) RefreshAccessToken(refreshToken string, user *models.UserProfileResponse, ipAddress string) (string, time.Time, error) {
	// Validate refresh token
	claims, err := jm.ValidateToken(refreshToken)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid refresh token: %w", err)
	}
	
	// Verify token type
	if claims.TokenType != models.TokenTypeRefresh {
		return "", time.Time{}, fmt.Errorf("invalid token type for refresh: %s", claims.TokenType)
	}
	
	// Generate new access token
	return jm.GenerateAccessToken(user, claims.SessionID, ipAddress)
}

// Permission checker for RBAC implementation
type PermissionChecker struct{}

// NewPermissionChecker creates a new permission checker instance.
func NewPermissionChecker() *PermissionChecker {
	return &PermissionChecker{}
}

// HasPermission checks if a user has a specific permission based on their roles.
// Implements the RBAC permission checking logic defined in SYSTEM_ARCHITECTURE.md
//
// Parameters:
//   - userRoles: list of user's role names
//   - requiredResource: resource being accessed (e.g., "controls")
//   - requiredAction: action being performed (e.g., "read", "write")
//   - requiredScope: scope of the operation (e.g., "own", "team", "organization")
//
// Returns:
//   - true if user has the required permission
//   - false if user lacks the required permission
//
// Permission Resolution:
//   - Checks each user role for matching permissions
//   - Supports wildcard permissions ("*" resource or action)
//   - Scope hierarchy: "organization" includes "team" and "own"
func (pc *PermissionChecker) HasPermission(userRoles []string, requiredResource, requiredAction, requiredScope string) bool {
	for _, roleName := range userRoles {
		role, exists := models.DefaultRoles[roleName]
		if !exists {
			continue
		}
		
		for _, permission := range role.Permissions {
			if pc.matchesPermission(permission, requiredResource, requiredAction, requiredScope) {
				return true
			}
		}
	}
	
	return false
}

// ValidatePermission validates if a user has the required permission, returning an error if not.
func (pc *PermissionChecker) ValidatePermission(userRoles []string, resource, action, scope string) error {
	if !pc.HasPermission(userRoles, resource, action, scope) {
		return fmt.Errorf("%w: %s:%s:%s", ErrInsufficientPermissions, resource, action, scope)
	}
	return nil
}

// matchesPermission checks if a permission matches the required access pattern.
func (pc *PermissionChecker) matchesPermission(permission models.Permission, resource, action, scope string) bool {
	// Check resource match (support wildcard)
	if permission.Resource != "*" && permission.Resource != resource {
		return false
	}
	
	// Check action match (support wildcard)
	if permission.Action != "*" && permission.Action != action {
		return false
	}
	
	// Check scope match with hierarchy
	return pc.scopeIncludes(permission.Scope, scope)
}

// scopeIncludes checks if the permission scope includes the required scope.
// Implements scope hierarchy: organization > team > own
func (pc *PermissionChecker) scopeIncludes(permissionScope, requiredScope string) bool {
	if permissionScope == requiredScope {
		return true
	}
	
	// Scope hierarchy
	switch permissionScope {
	case "organization":
		return requiredScope == "team" || requiredScope == "own"
	case "team":
		return requiredScope == "own"
	case "*":
		return true
	default:
		return false
	}
}

// Utility functions

// generateJTI creates a unique JWT ID for token tracking
func generateJTI() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("jti_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// Helper functions for extracting claims from JWT

func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getInt64Claim(claims jwt.MapClaims, key string) int64 {
	if val, ok := claims[key]; ok {
		switch v := val.(type) {
		case float64:
			return int64(v)
		case int64:
			return v
		case int:
			return int64(v)
		}
	}
	return 0
}

func getStringSliceClaim(claims jwt.MapClaims, key string) []string {
	if val, ok := claims[key]; ok {
		if slice, ok := val.([]interface{}); ok {
			result := make([]string, len(slice))
			for i, item := range slice {
				if str, ok := item.(string); ok {
					result[i] = str
				}
			}
			return result
		}
	}
	return nil
}