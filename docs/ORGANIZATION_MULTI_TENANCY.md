# Organization Entity and Multi-Tenancy Implementation

This document demonstrates the newly implemented Organization entity and multi-tenancy features for the GoEdu Control Testing Platform.

## 🏢 Enhanced Organization Entity

The Organization entity now supports comprehensive multi-tenant operations with the following key features:

### Core Organization Structure
```go
type Organization struct {
    BaseModel `bson:",inline"`
    
    // Basic organization information
    Name         string `bson:"name" json:"name"`
    DisplayName  string `bson:"display_name,omitempty" json:"display_name,omitempty"`
    Slug         string `bson:"slug" json:"slug"`
    Type         string `bson:"type" json:"type"` // commercial_bank, credit_union, etc.
    Industry     string `bson:"industry" json:"industry"`
    
    // Regulatory profile for compliance
    RegulatoryProfile RegulatoryProfile `bson:"regulatory_profile" json:"regulatory_profile"`
    
    // Subscription management
    Subscription OrganizationSubscription `bson:"subscription" json:"subscription"`
    
    // Feature flags per organization
    FeatureFlags map[string]bool `bson:"feature_flags,omitempty" json:"feature_flags,omitempty"`
    
    // Organization settings
    Settings OrganizationSettings `bson:"settings" json:"settings"`
}
```

### Regulatory Profile Support
```go
type RegulatoryProfile struct {
    Industry             string    `json:"industry"`
    PrimaryRegulator     string    `json:"primary_regulator"` // OCC, Federal Reserve, FDIC
    ApplicableFrameworks []string  `json:"applicable_frameworks"` // SOX, PCI-DSS, Basel III
    ExamCycle           string    `json:"exam_cycle"`
    RetentionPeriod     int       `json:"retention_period"` // 7 years for SOX
    RequiresSOX         bool      `json:"requires_sox"`
    RequiresPCIDSS      bool      `json:"requires_pci_dss"`
    RequiresFFIEC       bool      `json:"requires_ffiec"`
    RequiresBaselIII    bool      `json:"requires_basel_iii"`
}
```

### Subscription Management
```go
type OrganizationSubscription struct {
    Plan           string    `json:"plan"` // starter, professional, enterprise
    Status         string    `json:"status"` // active, trial, suspended
    UserLimit      int       `json:"user_limit"`
    StorageLimit   int64     `json:"storage_limit"`
    CurrentUsers   int       `json:"current_users"`
    CurrentStorage int64     `json:"current_storage"`
    IsInTrial      bool      `json:"is_in_trial"`
    TrialEnd       time.Time `json:"trial_end"`
}
```

## 🔒 Multi-Tenancy Data Isolation

### Organization Context Middleware

The middleware automatically enforces organization-level data isolation:

```go
// Usage in route setup
router.Use(organizationMiddleware.EnforceOrganizationContext())

// Access organization context in handlers
func (h *Handler) GetControls(c *gin.Context) {
    orgContext, err := middleware.GetOrganizationContext(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization context required"})
        return
    }
    
    // All data access is automatically scoped to orgContext.OrganizationID
    controls, err := h.controlService.GetControlsByOrganization(c, orgContext.OrganizationID.Hex())
    // ...
}
```

### Feature Flag Enforcement

Control access to features based on subscription tier:

```go
// Require specific feature for route access
router.Use(organizationMiddleware.RequireFeature("advanced_reporting"))

// Check feature programmatically
func (h *Handler) GenerateAdvancedReport(c *gin.Context) {
    orgContext, _ := middleware.GetOrganizationContext(c)
    
    if !orgContext.FeatureFlags["advanced_reporting"] {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "Advanced reporting not available in your subscription plan"
        })
        return
    }
    
    // Generate advanced report...
}
```

## 📊 Organization Service Layer

### Creating Organizations
```go
input := &services.CreateOrganizationInput{
    Name:         "First National Bank",
    Type:         "commercial_bank",
    Industry:     "financial_services",
    ContactEmail: "admin@firstnational.com",
    RegulatoryProfile: &services.CreateRegulatoryProfileInput{
        Industry:         "financial_services",
        PrimaryRegulator: "OCC",
        ApplicableFrameworks: []string{"SOX", "PCI-DSS", "Basel III"},
        RequiresSOX:      true,
        RequiresPCIDSS:   true,
        RetentionPeriod:  7,
    },
    SubscriptionPlan: "enterprise",
}

org, err := orgService.CreateOrganization(ctx, input)
```

### Managing Feature Flags
```go
// Enable a feature for an organization
err := orgService.UpdateFeatureFlag(ctx, orgID, "advanced_analytics", true)

// Check if feature is enabled
enabled, err := orgService.IsFeatureEnabled(ctx, orgID, "advanced_analytics")
```

## 🛡️ Data Isolation Examples

### Request Flow with Organization Context

1. **Request arrives with organization ID in header**:
   ```http
   GET /api/controls
   X-Organization-ID: 507f1f77bcf86cd799439011
   Authorization: Bearer <jwt_token>
   ```

2. **Middleware extracts and validates organization context**:
   - Validates user belongs to organization
   - Loads organization settings and feature flags
   - Injects context into request

3. **Handler automatically accesses organization-scoped data**:
   ```go
   func (h *Handler) GetControls(c *gin.Context) {
       // Organization ID automatically available
       orgID, _ := middleware.GetOrganizationID(c)
       
       // All queries automatically scoped to organization
       controls, err := h.controlRepo.GetByOrganization(c, orgID, filter)
   }
   ```

### Multi-Tenant Database Queries

All repository methods now include organization-based filtering:

```go
// Example: Get controls for specific organization only
func (r *controlRepository) GetByOrganization(ctx context.Context, orgID string, filter *ControlFilter) ([]*models.Control, error) {
    // MongoDB query automatically includes organization filter
    mongoFilter := bson.M{
        "organization_id": orgID,
        "status": "active",
    }
    
    // Additional filters...
    if filter.Framework != "" {
        mongoFilter["framework"] = filter.Framework
    }
    
    return r.find(ctx, mongoFilter)
}
```

## 🧪 Testing Multi-Tenancy

### Organization Entity Tests
```go
func TestOrganizationStructure(t *testing.T) {
    org := createValidOrganization()
    
    // Validate subscription structure
    assert.Equal(t, "professional", org.Subscription.Plan)
    assert.Equal(t, "active", org.Subscription.Status)
    assert.Greater(t, org.Subscription.UserLimit, 0)
    
    // Validate regulatory profile
    assert.True(t, org.RegulatoryProfile.RequiresSOX)
    assert.Contains(t, org.RegulatoryProfile.ApplicableFrameworks, "SOX")
    
    // Validate feature flags
    assert.True(t, org.FeatureFlags["basic_controls"])
    assert.True(t, org.FeatureFlags["advanced_reporting"])
}
```

### Middleware Tests
```go
func TestFeatureFlagEnforcement(t *testing.T) {
    // Test feature enabled
    orgContext := &OrganizationContext{
        FeatureFlags: map[string]bool{"advanced_reporting": true},
    }
    
    // Middleware should allow access
    handler := middleware.RequireFeature("advanced_reporting")
    // ... test passes through
    
    // Test feature disabled
    orgContext.FeatureFlags["advanced_reporting"] = false
    // ... test blocks access with 403 Forbidden
}
```

## 🔄 Migration from Existing Code

### Before (Single Tenant)
```go
// Old: Direct database access
controls, err := controlRepo.GetAll(ctx)
```

### After (Multi-Tenant)
```go
// New: Organization-scoped access
orgID, _ := middleware.GetOrganizationID(c)
controls, err := controlRepo.GetByOrganization(ctx, orgID, filter)
```

## 🚀 Benefits of Implementation

1. **Complete Data Isolation**: Each organization's data is strictly separated
2. **Feature Flag Control**: Different features available per subscription tier
3. **Regulatory Compliance**: Full support for banking regulations (SOX, PCI-DSS, etc.)
4. **Subscription Management**: Comprehensive billing and usage tracking
5. **Performance Optimization**: Caching layer for organization context
6. **Security**: User-to-organization access validation
7. **Audit Trail**: Complete logging of organization-level activities

## 📈 Subscription Tiers and Feature Matrix

| Feature | Starter | Professional | Enterprise |
|---------|---------|--------------|------------|
| Basic Controls | ✅ | ✅ | ✅ |
| Evidence Management | ✅ | ✅ | ✅ |
| Basic Reporting | ✅ | ✅ | ✅ |
| Advanced Reporting | ❌ | ✅ | ✅ |
| Custom Workflows | ❌ | ✅ | ✅ |
| API Access | ❌ | ✅ | ✅ |
| SSO Integration | ❌ | ✅ | ✅ |
| Advanced Analytics | ❌ | ❌ | ✅ |
| White Labeling | ❌ | ❌ | ✅ |
| User Limit | 10 | 50 | 1000 |
| Storage Limit | 5 GB | 50 GB | 500 GB |

This implementation provides a robust foundation for multi-tenant SaaS operations while maintaining strict data isolation and supporting complex banking regulatory requirements.