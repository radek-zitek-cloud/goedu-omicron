// Package auth_test provides comprehensive tests for the GoEdu authentication package.
// Tests cover password hashing, JWT token management, and RBAC permission checking
// following security best practices for financial compliance systems.
package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/models"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/auth"
)

// Test password hashing functionality
func TestPasswordHasher(t *testing.T) {
	hasher := auth.NewPasswordHasher(0) // Use default cost

	t.Run("hash and verify valid password", func(t *testing.T) {
		password := "SecurePassword123!"
		
		hash, err := hasher.HashPassword(password)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
		
		// Verify the password
		valid, err := hasher.VerifyPassword(password, hash)
		require.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("verify incorrect password", func(t *testing.T) {
		password := "SecurePassword123!"
		wrongPassword := "WrongPassword123!"
		
		hash, err := hasher.HashPassword(password)
		require.NoError(t, err)
		
		// Verify wrong password
		valid, err := hasher.VerifyPassword(wrongPassword, hash)
		require.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("reject password too short", func(t *testing.T) {
		shortPassword := "short"
		
		_, err := hasher.HashPassword(shortPassword)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "minimum length")
	})

	t.Run("reject password too long", func(t *testing.T) {
		// Create a password longer than 128 characters
		longPassword := string(make([]byte, 129))
		for i := range longPassword {
			longPassword = longPassword[:i] + "a"
		}
		
		_, err := hasher.HashPassword(longPassword)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "maximum length")
	})

	t.Run("get hash cost", func(t *testing.T) {
		password := "SecurePassword123!"
		
		hash, err := hasher.HashPassword(password)
		require.NoError(t, err)
		
		cost, err := hasher.GetHashCost(hash)
		require.NoError(t, err)
		assert.Equal(t, auth.BcryptCost, cost)
	})
}

// Test JWT token management
func TestJWTManager(t *testing.T) {
	secretKey := []byte("test-secret-key-256-bits-long-enough-for-hs256-algorithm")
	issuer := "goedu-test"
	audience := "goedu-api-test"
	
	jwtManager := auth.NewJWTManager(secretKey, issuer, audience)

	// Create test user profile
	userProfile := &models.UserProfileResponse{
		ID:             primitive.NewObjectID(),
		Email:          "test@example.com",
		FirstName:      "Test",
		LastName:       "User",
		OrganizationID: primitive.NewObjectID(),
		Role:           "auditor",
		Permissions:    []string{"controls:read", "assignments:read"},
		Status:         "active",
		MFAEnabled:     false,
	}

	t.Run("generate and validate access token", func(t *testing.T) {
		sessionID := "test-session-123"
		ipAddress := "192.168.1.1"
		
		// Generate access token
		token, expiresAt, err := jwtManager.GenerateAccessToken(userProfile, sessionID, ipAddress)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.True(t, expiresAt.After(time.Now()))
		
		// Validate token
		claims, err := jwtManager.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, issuer, claims.Issuer)
		assert.Equal(t, audience, claims.Audience)
		assert.Equal(t, userProfile.ID.Hex(), claims.UserID)
		assert.Equal(t, userProfile.Email, claims.Email)
		assert.Equal(t, sessionID, claims.SessionID)
		assert.Equal(t, models.TokenTypeAccess, claims.TokenType)
		assert.Equal(t, ipAddress, claims.IPAddress)
	})

	t.Run("generate and validate refresh token", func(t *testing.T) {
		userID := userProfile.ID.Hex()
		sessionID := "test-session-456"
		
		// Generate refresh token
		token, expiresAt, err := jwtManager.GenerateRefreshToken(userID, sessionID)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.True(t, expiresAt.After(time.Now()))
		
		// Validate token
		claims, err := jwtManager.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, issuer, claims.Issuer)
		assert.Equal(t, audience, claims.Audience)
		assert.Equal(t, userID, claims.Subject)
		assert.Equal(t, sessionID, claims.SessionID)
		assert.Equal(t, models.TokenTypeRefresh, claims.TokenType)
	})

	t.Run("reject invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.string"
		
		_, err := jwtManager.ValidateToken(invalidToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid or expired token")
	})

	t.Run("reject token with wrong secret", func(t *testing.T) {
		// Create token with different secret
		wrongSecretManager := auth.NewJWTManager([]byte("wrong-secret"), issuer, audience)
		token, _, err := wrongSecretManager.GenerateAccessToken(userProfile, "session", "127.0.0.1")
		require.NoError(t, err)
		
		// Try to validate with correct manager
		_, err = jwtManager.ValidateToken(token)
		assert.Error(t, err)
	})

	t.Run("refresh access token", func(t *testing.T) {
		userID := userProfile.ID.Hex()
		sessionID := "test-session-789"
		ipAddress := "192.168.1.2"
		
		// Generate refresh token
		refreshToken, _, err := jwtManager.GenerateRefreshToken(userID, sessionID)
		require.NoError(t, err)
		
		// Refresh access token
		newAccessToken, newExpiresAt, err := jwtManager.RefreshAccessToken(refreshToken, userProfile, ipAddress)
		require.NoError(t, err)
		assert.NotEmpty(t, newAccessToken)
		assert.True(t, newExpiresAt.After(time.Now()))
		
		// Validate new access token
		claims, err := jwtManager.ValidateToken(newAccessToken)
		require.NoError(t, err)
		assert.Equal(t, sessionID, claims.SessionID)
		assert.Equal(t, models.TokenTypeAccess, claims.TokenType)
	})

	t.Run("reject wrong token type for refresh", func(t *testing.T) {
		// Try to refresh using an access token instead of refresh token
		accessToken, _, err := jwtManager.GenerateAccessToken(userProfile, "session", "127.0.0.1")
		require.NoError(t, err)
		
		_, _, err = jwtManager.RefreshAccessToken(accessToken, userProfile, "127.0.0.1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token type")
	})
}

// Test permission checking functionality
func TestPermissionChecker(t *testing.T) {
	checker := auth.NewPermissionChecker()

	t.Run("admin has all permissions", func(t *testing.T) {
		userRoles := []string{models.RoleAdmin}
		
		// Admin should have access to everything
		assert.True(t, checker.HasPermission(userRoles, "controls", "read", "organization"))
		assert.True(t, checker.HasPermission(userRoles, "controls", "write", "organization"))
		assert.True(t, checker.HasPermission(userRoles, "users", "manage", "organization"))
		assert.True(t, checker.HasPermission(userRoles, "any_resource", "any_action", "any_scope"))
	})

	t.Run("auditor has limited permissions", func(t *testing.T) {
		userRoles := []string{"auditor"}
		
		// Auditor should have read access to controls
		assert.True(t, checker.HasPermission(userRoles, "controls", "read", "organization"))
		assert.True(t, checker.HasPermission(userRoles, "controls", "read", "team"))
		assert.True(t, checker.HasPermission(userRoles, "controls", "read", "own"))
		
		// But not write access
		assert.False(t, checker.HasPermission(userRoles, "controls", "write", "organization"))
		
		// Should have access to own assignments
		assert.True(t, checker.HasPermission(userRoles, "assignments", "read", "own"))
		assert.True(t, checker.HasPermission(userRoles, "assignments", "update", "own"))
		
		// But not team assignments
		assert.False(t, checker.HasPermission(userRoles, "assignments", "read", "team"))
	})

	t.Run("audit_manager has team permissions", func(t *testing.T) {
		userRoles := []string{"audit_manager"}
		
		// Should have control management permissions
		assert.True(t, checker.HasPermission(userRoles, "controls", "read", "organization"))
		assert.True(t, checker.HasPermission(userRoles, "controls", "write", "organization"))
		
		// Should have team assignment permissions
		assert.True(t, checker.HasPermission(userRoles, "assignments", "create", "team"))
		assert.True(t, checker.HasPermission(userRoles, "assignments", "read", "team"))
		assert.True(t, checker.HasPermission(userRoles, "assignments", "update", "team"))
		
		// Should also have own permissions (scope hierarchy)
		assert.True(t, checker.HasPermission(userRoles, "assignments", "read", "own"))
	})

	t.Run("unknown role has no permissions", func(t *testing.T) {
		userRoles := []string{"unknown_role"}
		
		assert.False(t, checker.HasPermission(userRoles, "controls", "read", "organization"))
		assert.False(t, checker.HasPermission(userRoles, "assignments", "read", "own"))
	})

	t.Run("validate permission returns error for insufficient access", func(t *testing.T) {
		userRoles := []string{"auditor"}
		
		// Should not error for valid permission
		err := checker.ValidatePermission(userRoles, "controls", "read", "organization")
		assert.NoError(t, err)
		
		// Should error for invalid permission
		err = checker.ValidatePermission(userRoles, "controls", "write", "organization")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient permissions")
	})

	t.Run("scope hierarchy works correctly", func(t *testing.T) {
		userRoles := []string{"audit_manager"}
		
		// Organization scope should include team and own
		assert.True(t, checker.HasPermission(userRoles, "controls", "read", "organization"))
		assert.True(t, checker.HasPermission(userRoles, "controls", "read", "team"))
		assert.True(t, checker.HasPermission(userRoles, "controls", "read", "own"))
		
		// Team scope should include own but not organization
		// (audit_manager has team scope for assignments)
		assert.True(t, checker.HasPermission(userRoles, "assignments", "create", "team"))
		assert.True(t, checker.HasPermission(userRoles, "assignments", "create", "own"))
	})
}

// Test authentication error constants
func TestAuthenticationErrors(t *testing.T) {
	t.Run("error constants are defined", func(t *testing.T) {
		assert.NotEmpty(t, auth.ErrInvalidCredentials.Error())
		assert.NotEmpty(t, auth.ErrAccountLocked.Error())
		assert.NotEmpty(t, auth.ErrAccountInactive.Error())
		assert.NotEmpty(t, auth.ErrPasswordTooWeak.Error())
		assert.NotEmpty(t, auth.ErrInvalidToken.Error())
		assert.NotEmpty(t, auth.ErrTokenExpired.Error())
		assert.NotEmpty(t, auth.ErrInsufficientPermissions.Error())
	})
}