// Package middleware_test contains tests for the middleware package.
// This file tests the organization middleware and multi-tenancy functionality.
package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/models"
)

// MockOrganizationService is a mock implementation of OrganizationService for testing
type MockOrganizationService struct {
	mock.Mock
}

func (m *MockOrganizationService) GetOrganization(ctx context.Context, id string) (*models.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Organization), args.Error(1)
}

func (m *MockOrganizationService) GetOrganizationBySlug(ctx context.Context, slug string) (*models.Organization, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Organization), args.Error(1)
}

func (m *MockOrganizationService) GetFeatureFlags(ctx context.Context, orgID string) (map[string]bool, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]bool), args.Error(1)
}

func (m *MockOrganizationService) ValidateOrganizationAccess(ctx context.Context, userID, orgID string) error {
	args := m.Called(ctx, userID, orgID)
	return args.Error(0)
}

// MockUserService is a mock implementation of UserService for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUser(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// TestOrganizationMiddleware tests the organization middleware functionality
func TestOrganizationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name        string
		description string
		testFunc    func(t *testing.T)
	}{
		{
			name:        "middleware_creation",
			description: "Middleware should be created with required dependencies",
			testFunc:    testMiddlewareCreation,
		},
		{
			name:        "organization_context_extraction",
			description: "Middleware should extract organization context from request",
			testFunc:    testOrganizationContextExtraction,
		},
		{
			name:        "feature_flag_enforcement",
			description: "Middleware should enforce feature flag restrictions",
			testFunc:    testFeatureFlagEnforcement,
		},
		{
			name:        "helper_functions",
			description: "Helper functions should work correctly",
			testFunc:    testHelperFunctions,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.testFunc(t)
		})
	}
}

// testMiddlewareCreation validates middleware creation
func testMiddlewareCreation(t *testing.T) {
	orgService := &MockOrganizationService{}
	userService := &MockUserService{}
	logger := zap.NewNop()

	middleware := NewOrganizationMiddleware(orgService, userService, logger)

	assert.NotNil(t, middleware)
	assert.NotNil(t, middleware.orgService)
	assert.NotNil(t, middleware.userService)
	assert.NotNil(t, middleware.logger)
}

// testOrganizationContextExtraction tests context extraction functionality
func testOrganizationContextExtraction(t *testing.T) {
	orgService := &MockOrganizationService{}
	userService := &MockUserService{}
	logger := zap.NewNop()

	middleware := NewOrganizationMiddleware(orgService, userService, logger)

	// Test extracting organization ID from header
	t.Run("extract_from_header", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/api/test", nil)
		req.Header.Set("X-Organization-ID", "507f1f77bcf86cd799439011")
		c.Request = req

		orgID, err := middleware.extractOrganizationID(c)
		assert.NoError(t, err)
		assert.Equal(t, "507f1f77bcf86cd799439011", orgID.Hex())
	})

	// Test extracting user ID from context
	t.Run("extract_user_from_context", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("user_id", "507f1f77bcf86cd799439012")

		userID, err := middleware.extractUserID(c)
		assert.NoError(t, err)
		assert.Equal(t, "507f1f77bcf86cd799439012", userID.Hex())
	})

	// Test missing organization ID
	t.Run("missing_organization_id", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/api/test", nil)
		c.Request = req

		_, err := middleware.extractOrganizationID(c)
		assert.Error(t, err)
		assert.Equal(t, ErrOrganizationIDNotFound, err)
	})
}

// testFeatureFlagEnforcement tests feature flag middleware
func testFeatureFlagEnforcement(t *testing.T) {
	orgService := &MockOrganizationService{}
	userService := &MockUserService{}
	logger := zap.NewNop()

	middleware := NewOrganizationMiddleware(orgService, userService, logger)

	t.Run("feature_enabled", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// Set up organization context with feature enabled
		orgContext := &OrganizationContext{
			OrganizationID: primitive.NewObjectID(),
			FeatureFlags: map[string]bool{
				"advanced_reporting": true,
			},
		}

		// Set up route with middleware
		var nextCalled bool
		router.Use(func(c *gin.Context) {
			c.Set("organization_context", orgContext)
			c.Next()
		})
		router.Use(middleware.RequireFeature("advanced_reporting"))
		router.GET("/test", func(c *gin.Context) {
			nextCalled = true
			c.Status(http.StatusOK)
		})

		// Execute request
		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req
		router.ServeHTTP(w, req)

		// Verify next was called (feature is enabled)
		assert.True(t, nextCalled)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("feature_disabled", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// Set up organization context with feature disabled
		orgContext := &OrganizationContext{
			OrganizationID: primitive.NewObjectID(),
			FeatureFlags: map[string]bool{
				"advanced_reporting": false,
			},
		}

		// Set up route with middleware
		var nextCalled bool
		router.Use(func(c *gin.Context) {
			c.Set("organization_context", orgContext)
			c.Next()
		})
		router.Use(middleware.RequireFeature("advanced_reporting"))
		router.GET("/test", func(c *gin.Context) {
			nextCalled = true
			c.Status(http.StatusOK)
		})

		// Execute request
		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req
		router.ServeHTTP(w, req)

		// Verify next was not called (feature is disabled)
		assert.False(t, nextCalled)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("missing_organization_context", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// Set up route with middleware but no context
		var nextCalled bool
		router.Use(middleware.RequireFeature("advanced_reporting"))
		router.GET("/test", func(c *gin.Context) {
			nextCalled = true
			c.Status(http.StatusOK)
		})

		// Execute request
		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req
		router.ServeHTTP(w, req)

		// Verify next was not called (missing context)
		assert.False(t, nextCalled)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// testHelperFunctions tests the helper functions
func testHelperFunctions(t *testing.T) {
	t.Run("get_organization_context_success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		expectedContext := &OrganizationContext{
			OrganizationID: primitive.NewObjectID(),
			UserID:         primitive.NewObjectID(),
		}
		c.Set("organization_context", expectedContext)

		context, err := GetOrganizationContext(c)
		assert.NoError(t, err)
		assert.Equal(t, expectedContext, context)
	})

	t.Run("get_organization_context_missing", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		context, err := GetOrganizationContext(c)
		assert.Error(t, err)
		assert.Nil(t, context)
		assert.Equal(t, ErrOrganizationContextNotFound, err)
	})

	t.Run("get_organization_id_success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		expectedID := "507f1f77bcf86cd799439011"
		c.Set("organization_id", expectedID)

		id, err := GetOrganizationID(c)
		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
	})

	t.Run("get_organization_id_missing", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		id, err := GetOrganizationID(c)
		assert.Error(t, err)
		assert.Empty(t, id)
		assert.Equal(t, ErrOrganizationIDNotFound, err)
	})
}

// TestMiddlewareErrors tests custom middleware errors
func TestMiddlewareErrors(t *testing.T) {
	t.Run("middleware_error_implements_error", func(t *testing.T) {
		err := &MiddlewareError{
			Code:    "TEST_ERROR",
			Message: "Test error message",
		}

		assert.Equal(t, "Test error message", err.Error())
		assert.Equal(t, "TEST_ERROR", err.Code)
	})

	t.Run("predefined_errors", func(t *testing.T) {
		assert.Equal(t, "ORG_ID_NOT_FOUND", ErrOrganizationIDNotFound.Code)
		assert.Equal(t, "USER_ID_NOT_FOUND", ErrUserIDNotFound.Code)
		assert.Equal(t, "ORG_CONTEXT_NOT_FOUND", ErrOrganizationContextNotFound.Code)
		assert.Equal(t, "INVALID_ORG_CONTEXT", ErrInvalidOrganizationContext.Code)
		assert.Equal(t, "INVALID_ORG_ID", ErrInvalidOrganizationID.Code)
	})
}

// TestOrganizationContextStructure tests the OrganizationContext structure
func TestOrganizationContextStructure(t *testing.T) {
	orgID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	
	context := &OrganizationContext{
		OrganizationID:   orgID,
		UserID:           userID,
		UserRole:         "auditor",
		UserPermissions:  []string{"controls:read", "evidence:upload"},
		FeatureFlags:     map[string]bool{"basic_controls": true},
		SubscriptionTier: "professional",
		IsActive:         true,
	}

	assert.Equal(t, orgID, context.OrganizationID)
	assert.Equal(t, userID, context.UserID)
	assert.Equal(t, "auditor", context.UserRole)
	assert.Contains(t, context.UserPermissions, "controls:read")
	assert.True(t, context.FeatureFlags["basic_controls"])
	assert.Equal(t, "professional", context.SubscriptionTier)
	assert.True(t, context.IsActive)
}