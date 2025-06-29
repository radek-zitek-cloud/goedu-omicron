// Package middleware provides HTTP middleware functions for the GoEdu Control Testing Platform.
// This file contains organization-level multi-tenancy and data isolation middleware.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/models"
)

// OrganizationContextKey is the key used to store organization context in the request context
type OrganizationContextKey string

const (
	// OrganizationIDKey is the context key for organization ID
	OrganizationIDKey OrganizationContextKey = "organization_id"
	// OrganizationKey is the context key for the full organization object
	OrganizationKey OrganizationContextKey = "organization"
	// UserOrganizationKey is the context key for user's organization association
	UserOrganizationKey OrganizationContextKey = "user_organization"
)

// OrganizationContext contains organization-related information for the current request.
// This provides comprehensive context for multi-tenant operations and data isolation.
type OrganizationContext struct {
	OrganizationID   primitive.ObjectID     `json:"organization_id"`
	Organization     *models.Organization   `json:"organization"`
	UserID           primitive.ObjectID     `json:"user_id"`
	UserRole         string                 `json:"user_role"`
	UserPermissions  []string               `json:"user_permissions"`
	FeatureFlags     map[string]bool        `json:"feature_flags"`
	Settings         *models.OrganizationSettings `json:"settings"`
	SubscriptionTier string                 `json:"subscription_tier"`
	IsActive         bool                   `json:"is_active"`
}

// OrganizationService interface for middleware dependencies
type OrganizationService interface {
	GetOrganization(ctx context.Context, id string) (*models.Organization, error)
	GetOrganizationBySlug(ctx context.Context, slug string) (*models.Organization, error)
	GetFeatureFlags(ctx context.Context, orgID string) (map[string]bool, error)
	ValidateOrganizationAccess(ctx context.Context, userID, orgID string) error
}

// UserService interface for middleware dependencies
type UserService interface {
	GetUser(ctx context.Context, id string) (*models.User, error)
}

// OrganizationMiddleware provides multi-tenancy and data isolation functionality.
// It ensures that all requests are properly scoped to the correct organization
// and that users can only access data belonging to their organization.
type OrganizationMiddleware struct {
	orgService  OrganizationService
	userService UserService
	logger      *zap.Logger
}

// NewOrganizationMiddleware creates a new organization middleware with required dependencies.
//
// Parameters:
//   - orgService: Service for organization operations
//   - userService: Service for user operations  
//   - logger: Logger for audit and debugging
//
// Returns:
//   - *OrganizationMiddleware: Configured middleware instance
func NewOrganizationMiddleware(
	orgService OrganizationService,
	userService UserService,
	logger *zap.Logger,
) *OrganizationMiddleware {
	return &OrganizationMiddleware{
		orgService:  orgService,
		userService: userService,
		logger:      logger,
	}
}

// EnforceOrganizationContext is the main middleware function that enforces organization-level
// data isolation. It extracts the organization context from the request and validates
// that the user has appropriate access to the organization's data.
//
// This middleware:
// 1. Extracts organization ID from request (header, path, or user context)
// 2. Validates user access to the organization
// 3. Loads organization settings and feature flags
// 4. Injects organization context into the request context
// 5. Enforces data isolation boundaries
//
// Usage:
//   router.Use(orgMiddleware.EnforceOrganizationContext())
func (m *OrganizationMiddleware) EnforceOrganizationContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		
		// Extract organization ID from various sources
		orgID, err := m.extractOrganizationID(c)
		if err != nil {
			m.logger.Warn("Failed to extract organization ID",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid or missing organization context",
				"code":  "INVALID_ORGANIZATION_CONTEXT",
			})
			return
		}

		// Extract user ID from JWT token (assuming it's already validated by auth middleware)
		userID, err := m.extractUserID(c)
		if err != nil {
			m.logger.Warn("Failed to extract user ID",
				zap.Error(err),
				zap.String("organization_id", orgID.Hex()),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User authentication required",
				"code":  "USER_NOT_AUTHENTICATED",
			})
			return
		}

		// Validate user access to organization
		if err := m.orgService.ValidateOrganizationAccess(ctx, userID.Hex(), orgID.Hex()); err != nil {
			m.logger.Warn("User access denied to organization",
				zap.Error(err),
				zap.String("user_id", userID.Hex()),
				zap.String("organization_id", orgID.Hex()),
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Access denied to organization",
				"code":  "ORGANIZATION_ACCESS_DENIED",
			})
			return
		}

		// Load organization details
		org, err := m.orgService.GetOrganization(ctx, orgID.Hex())
		if err != nil {
			m.logger.Error("Failed to load organization",
				zap.Error(err),
				zap.String("organization_id", orgID.Hex()),
			)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to load organization context",
				"code":  "ORGANIZATION_LOAD_ERROR",
			})
			return
		}

		// Check if organization is active
		if !org.IsActive || org.Status != models.OrganizationStatusActive {
			m.logger.Warn("Access attempted to inactive organization",
				zap.String("organization_id", orgID.Hex()),
				zap.String("status", org.Status),
				zap.Bool("is_active", org.IsActive),
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Organization is not active",
				"code":  "ORGANIZATION_INACTIVE",
			})
			return
		}

		// Load user details
		user, err := m.userService.GetUser(ctx, userID.Hex())
		if err != nil {
			m.logger.Error("Failed to load user",
				zap.Error(err),
				zap.String("user_id", userID.Hex()),
			)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to load user context",
				"code":  "USER_LOAD_ERROR",
			})
			return
		}

		// Load organization feature flags
		featureFlags, err := m.orgService.GetFeatureFlags(ctx, orgID.Hex())
		if err != nil {
			m.logger.Warn("Failed to load feature flags, using defaults",
				zap.Error(err),
				zap.String("organization_id", orgID.Hex()),
			)
			featureFlags = make(map[string]bool) // Default to empty flags
		}

		// Create organization context
		orgContext := &OrganizationContext{
			OrganizationID:   orgID,
			Organization:     org,
			UserID:           userID,
			UserRole:         strings.Join(user.Roles, ","), // Primary role for simplicity
			UserPermissions:  user.GetPermissionsList(),
			FeatureFlags:     featureFlags,
			Settings:         &org.Settings,
			SubscriptionTier: org.Subscription.Plan,
			IsActive:         org.IsActive,
		}

		// Inject organization context into request context
		ctx = context.WithValue(ctx, OrganizationIDKey, orgID)
		ctx = context.WithValue(ctx, OrganizationKey, org)
		ctx = context.WithValue(ctx, UserOrganizationKey, orgContext)
		c.Request = c.Request.WithContext(ctx)

		// Set organization context in Gin context for easy access
		c.Set("organization_context", orgContext)
		c.Set("organization_id", orgID.Hex())
		c.Set("organization", org)

		// Log successful organization context establishment
		m.logger.Debug("Organization context established",
			zap.String("organization_id", orgID.Hex()),
			zap.String("organization_name", org.Name),
			zap.String("user_id", userID.Hex()),
			zap.String("user_role", orgContext.UserRole),
			zap.String("subscription_tier", orgContext.SubscriptionTier),
		)

		c.Next()
	}
}

// RequireFeature creates middleware that checks if a specific feature is enabled
// for the organization. This enforces feature flag-based access control.
//
// Parameters:
//   - feature: The feature flag name to check
//
// Returns:
//   - gin.HandlerFunc: Middleware function that validates feature access
//
// Usage:
//   router.Use(orgMiddleware.RequireFeature("advanced_reporting"))
func (m *OrganizationMiddleware) RequireFeature(feature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgContext, exists := c.Get("organization_context")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Organization context not found",
				"code":  "MISSING_ORGANIZATION_CONTEXT",
			})
			return
		}

		ctx, ok := orgContext.(*OrganizationContext)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid organization context",
				"code":  "INVALID_ORGANIZATION_CONTEXT",
			})
			return
		}

		// Check if feature is enabled
		if enabled, exists := ctx.FeatureFlags[feature]; !exists || !enabled {
			m.logger.Warn("Feature access denied",
				zap.String("feature", feature),
				zap.String("organization_id", ctx.OrganizationID.Hex()),
				zap.Bool("feature_enabled", enabled),
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Feature not available for your organization",
				"code":    "FEATURE_NOT_ENABLED",
				"feature": feature,
			})
			return
		}

		c.Next()
	}
}

// extractOrganizationID extracts the organization ID from various request sources.
// It checks in the following order:
// 1. X-Organization-ID header
// 2. organization_id path parameter
// 3. User's organization ID from JWT token
func (m *OrganizationMiddleware) extractOrganizationID(c *gin.Context) (primitive.ObjectID, error) {
	// Try header first
	if orgIDHeader := c.GetHeader("X-Organization-ID"); orgIDHeader != "" {
		return primitive.ObjectIDFromHex(orgIDHeader)
	}

	// Try path parameter
	if orgIDParam := c.Param("organization_id"); orgIDParam != "" {
		return primitive.ObjectIDFromHex(orgIDParam)
	}
	if orgIDParam := c.Param("orgId"); orgIDParam != "" {
		return primitive.ObjectIDFromHex(orgIDParam)
	}

	// Try organization slug in path
	if orgSlug := c.Param("organization_slug"); orgSlug != "" {
		return m.getOrganizationIDBySlug(c.Request.Context(), orgSlug)
	}

	// Try to extract from user context (JWT token should contain organization_id)
	if userOrgID, exists := c.Get("user_organization_id"); exists {
		if orgID, ok := userOrgID.(string); ok {
			return primitive.ObjectIDFromHex(orgID)
		}
	}

	return primitive.NilObjectID, ErrOrganizationIDNotFound
}

// extractUserID extracts the user ID from the JWT token context.
// This assumes the authentication middleware has already validated the token
// and injected the user ID into the context.
func (m *OrganizationMiddleware) extractUserID(c *gin.Context) (primitive.ObjectID, error) {
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			return primitive.ObjectIDFromHex(uid)
		}
	}
	return primitive.NilObjectID, ErrUserIDNotFound
}

// getOrganizationIDBySlug resolves an organization slug to an organization ID.
func (m *OrganizationMiddleware) getOrganizationIDBySlug(ctx context.Context, slug string) (primitive.ObjectID, error) {
	org, err := m.orgService.GetOrganizationBySlug(ctx, slug)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return org.ID, nil
}

// GetOrganizationContext extracts the organization context from the Gin context.
// This is a helper function for handlers to easily access organization information.
//
// Parameters:
//   - c: Gin context containing the organization context
//
// Returns:
//   - *OrganizationContext: Organization context or nil if not found
//   - error: Error if context is invalid or missing
func GetOrganizationContext(c *gin.Context) (*OrganizationContext, error) {
	orgContext, exists := c.Get("organization_context")
	if !exists {
		return nil, ErrOrganizationContextNotFound
	}

	ctx, ok := orgContext.(*OrganizationContext)
	if !ok {
		return nil, ErrInvalidOrganizationContext
	}

	return ctx, nil
}

// GetOrganizationID extracts just the organization ID from the context.
// This is a convenience function for handlers that only need the ID.
//
// Parameters:
//   - c: Gin context containing the organization ID
//
// Returns:
//   - string: Organization ID as a string
//   - error: Error if ID is not found or invalid
func GetOrganizationID(c *gin.Context) (string, error) {
	orgID, exists := c.Get("organization_id")
	if !exists {
		return "", ErrOrganizationIDNotFound
	}

	if id, ok := orgID.(string); ok {
		return id, nil
	}

	return "", ErrInvalidOrganizationID
}

// Custom errors for organization middleware
var (
	ErrOrganizationIDNotFound      = &MiddlewareError{Code: "ORG_ID_NOT_FOUND", Message: "Organization ID not found in request"}
	ErrUserIDNotFound              = &MiddlewareError{Code: "USER_ID_NOT_FOUND", Message: "User ID not found in request context"}
	ErrOrganizationContextNotFound = &MiddlewareError{Code: "ORG_CONTEXT_NOT_FOUND", Message: "Organization context not found"}
	ErrInvalidOrganizationContext  = &MiddlewareError{Code: "INVALID_ORG_CONTEXT", Message: "Invalid organization context type"}
	ErrInvalidOrganizationID       = &MiddlewareError{Code: "INVALID_ORG_ID", Message: "Invalid organization ID type"}
)

// MiddlewareError represents an error that occurred in middleware processing.
type MiddlewareError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *MiddlewareError) Error() string {
	return e.Message
}