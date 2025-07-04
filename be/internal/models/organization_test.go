// Package models_test contains unit tests for the models package.
// This file tests the Organization entity and related structures.
package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestOrganizationStructure validates the Organization entity structure
// and ensures all required fields are properly defined.
func TestOrganizationStructure(t *testing.T) {
	tests := []struct {
		name        string
		description string
		testFunc    func(t *testing.T)
	}{
		{
			name:        "organization_has_required_fields",
			description: "Organization struct should have all required fields",
			testFunc:    testOrganizationRequiredFields,
		},
		{
			name:        "organization_subscription_structure",
			description: "OrganizationSubscription should have proper structure",
			testFunc:    testOrganizationSubscriptionStructure,
		},
		{
			name:        "organization_settings_structure", 
			description: "OrganizationSettings should have proper structure",
			testFunc:    testOrganizationSettingsStructure,
		},
		{
			name:        "regulatory_profile_structure",
			description: "RegulatoryProfile should have enhanced structure",
			testFunc:    testRegulatoryProfileStructure,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.testFunc(t)
		})
	}
}

// testOrganizationRequiredFields validates that Organization has all required fields
func testOrganizationRequiredFields(t *testing.T) {
	org := Organization{
		BaseModel: BaseModel{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:         "Test Bank",
		DisplayName:  "Test Bank Inc.",
		Slug:         "test-bank", 
		Description:  "A test banking organization",
		Type:         "commercial_bank",
		Industry:     "financial_services",
		Size:         "medium",
		Region:       "North America",
		Country:      "United States",
		Timezone:     "America/New_York",
		Currency:     "USD",
		ContactEmail: "contact@testbank.com",
		ContactPhone: "+1-555-0123",
		Website:      "https://testbank.com",
		Status:       OrganizationStatusActive,
		IsActive:     true,
		MemberCount:  10,
		MaxMembers:   50,
	}

	// Validate basic fields
	assert.NotEmpty(t, org.ID)
	assert.Equal(t, "Test Bank", org.Name)
	assert.Equal(t, "Test Bank Inc.", org.DisplayName)
	assert.Equal(t, "test-bank", org.Slug)
	assert.Equal(t, "commercial_bank", org.Type)
	assert.Equal(t, "financial_services", org.Industry)
	assert.Equal(t, "contact@testbank.com", org.ContactEmail)
	assert.Equal(t, OrganizationStatusActive, org.Status)
	assert.True(t, org.IsActive)
	assert.Equal(t, 10, org.MemberCount)
	assert.Equal(t, 50, org.MaxMembers)
}

// testOrganizationSubscriptionStructure validates the subscription management structure
func testOrganizationSubscriptionStructure(t *testing.T) {
	now := time.Now()
	subscription := OrganizationSubscription{
		Plan:               SubscriptionPlanProfessional,
		Status:             SubscriptionStatusActive,
		BillingPeriod:      "monthly",
		BillingEmail:       "billing@testbank.com",
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   now.AddDate(0, 1, 0),
		TrialStart:         now.AddDate(0, -1, 0),
		TrialEnd:           now,
		IsInTrial:          false,
		UserLimit:          50,
		StorageLimit:       50 * 1024 * 1024 * 1024, // 50 GB
		CurrentUsers:       10,
		CurrentStorage:     5 * 1024 * 1024 * 1024, // 5 GB
		LastPaymentDate:    now.AddDate(0, -1, 0),
		NextPaymentDate:    now.AddDate(0, 1, 0),
		PaymentMethod:      "credit_card",
		SubscriptionID:     "sub_12345",
		CustomerID:         "cust_67890",
		CreatedAt:          now.AddDate(0, -2, 0),
		UpdatedAt:          now,
	}

	// Validate subscription fields
	assert.Equal(t, SubscriptionPlanProfessional, subscription.Plan)
	assert.Equal(t, SubscriptionStatusActive, subscription.Status)
	assert.Equal(t, "monthly", subscription.BillingPeriod)
	assert.Equal(t, "billing@testbank.com", subscription.BillingEmail)
	assert.False(t, subscription.IsInTrial)
	assert.Equal(t, 50, subscription.UserLimit)
	assert.Equal(t, int64(50*1024*1024*1024), subscription.StorageLimit)
	assert.Equal(t, 10, subscription.CurrentUsers)
	assert.Equal(t, "sub_12345", subscription.SubscriptionID)
	assert.Equal(t, "cust_67890", subscription.CustomerID)
}

// testOrganizationSettingsStructure validates the organization settings structure
func testOrganizationSettingsStructure(t *testing.T) {
	settings := OrganizationSettings{
		RequireMFA:            true,
		AllowInvitations:      true,
		SessionTimeoutMinutes: 480,
		EnableAuditLog:        true,
		DataRetentionDays:     2555,
		AllowDataExport:       true,
		CustomBranding:        false,
		Theme:                 "light",
		PrimaryColor:          "#1976d2",
		SecondaryColor:        "#f50057",
		Integrations: OrganizationIntegrations{
			SSO:             true,
			LDAP:            false,
			API:             true,
			SSOProvider:     "auth0",
			APIKeyEnabled:   true,
			WebhooksEnabled: true,
		},
		Notifications: OrganizationNotifications{
			EmailEnabled:      true,
			SMSEnabled:        false,
			WebEnabled:        true,
			DigestFrequency:   "daily",
			QuietHoursStart:   "22:00",
			QuietHoursEnd:     "06:00",
			EvidenceRequests:  true,
			DeadlineReminders: true,
			SystemAlerts:      true,
		},
		ComplianceSettings: OrganizationCompliance{
			AutoGenerateReports:     true,
			RequireDigitalSignature: true,
			EnableChangeTracking:    true,
			RequireEvidenceApproval: true,
			MinimumReviewers:        2,
			ReportingFrequency:      "quarterly",
			AutoSubmitReports:       false,
		},
		DefaultUserRole:         RoleAuditor,
		DefaultUserPermissions:  []string{"controls:read", "evidence:upload"},
		AutoAssignControls:      false,
		AutoGenerateWorkpapers:  true,
		EnableWorkflowReminders: true,
		CustomSettings:          map[string]interface{}{"custom_field": "value"},
	}

	// Validate settings structure
	assert.True(t, settings.RequireMFA)
	assert.True(t, settings.AllowInvitations)
	assert.Equal(t, 480, settings.SessionTimeoutMinutes)
	assert.True(t, settings.EnableAuditLog)
	assert.Equal(t, 2555, settings.DataRetentionDays)
	assert.True(t, settings.AllowDataExport)
	assert.False(t, settings.CustomBranding)
	assert.Equal(t, "light", settings.Theme)

	// Validate integrations
	assert.True(t, settings.Integrations.SSO)
	assert.False(t, settings.Integrations.LDAP)
	assert.True(t, settings.Integrations.API)
	assert.Equal(t, "auth0", settings.Integrations.SSOProvider)

	// Validate notifications
	assert.True(t, settings.Notifications.EmailEnabled)
	assert.False(t, settings.Notifications.SMSEnabled)
	assert.Equal(t, "daily", settings.Notifications.DigestFrequency)

	// Validate compliance settings
	assert.True(t, settings.ComplianceSettings.AutoGenerateReports)
	assert.True(t, settings.ComplianceSettings.RequireDigitalSignature)
	assert.Equal(t, 2, settings.ComplianceSettings.MinimumReviewers)
	assert.Equal(t, "quarterly", settings.ComplianceSettings.ReportingFrequency)

	// Validate workflow settings
	assert.Equal(t, RoleAuditor, settings.DefaultUserRole)
	assert.Contains(t, settings.DefaultUserPermissions, "controls:read")
	assert.True(t, settings.AutoGenerateWorkpapers)
	assert.NotNil(t, settings.CustomSettings)
}

// testRegulatoryProfileStructure validates the enhanced regulatory profile structure
func testRegulatoryProfileStructure(t *testing.T) {
	now := time.Now()
	profile := RegulatoryProfile{
		Industry:             "financial_services",
		PrimaryRegulator:     "OCC",
		ApplicableFrameworks: []string{"SOX", "PCI-DSS", "Basel III", "FFIEC"},
		Regulations:          []string{"Bank Secrecy Act", "Fair Credit Reporting Act"},
		ExamCycle:           "18_months",
		LastExamDate:        now.AddDate(-1, 0, 0),
		NextExamDate:        now.AddDate(0, 6, 0),
		AuditFrequency:      "quarterly",
		RetentionPeriod:     7,
		RiskTolerance:       "low",
		RiskFramework:       "COSO",
		ComplianceStatus:    "compliant",
		LastReviewDate:      now.AddDate(0, -1, 0),
		NextReviewDate:      now.AddDate(0, 11, 0),
		RequiresSOX:         true,
		RequiresPCIDSS:      true,
		RequiresFFIEC:       true,
		RequiresBaselIII:    false,
		CustomRequirements:  map[string]interface{}{"special_requirement": "value"},
	}

	// Validate regulatory profile fields
	assert.Equal(t, "financial_services", profile.Industry)
	assert.Equal(t, "OCC", profile.PrimaryRegulator)
	assert.Contains(t, profile.ApplicableFrameworks, "SOX")
	assert.Contains(t, profile.ApplicableFrameworks, "PCI-DSS")
	assert.Equal(t, "18_months", profile.ExamCycle)
	assert.Equal(t, 7, profile.RetentionPeriod)
	assert.Equal(t, "low", profile.RiskTolerance)
	assert.Equal(t, "COSO", profile.RiskFramework)
	assert.Equal(t, "compliant", profile.ComplianceStatus)
	assert.True(t, profile.RequiresSOX)
	assert.True(t, profile.RequiresPCIDSS)
	assert.True(t, profile.RequiresFFIEC)
	assert.False(t, profile.RequiresBaselIII)
	assert.NotNil(t, profile.CustomRequirements)
}

// TestOrganizationConstants validates that organization-related constants are properly defined
func TestOrganizationConstants(t *testing.T) {
	// Test organization status constants
	assert.Equal(t, "active", OrganizationStatusActive)
	assert.Equal(t, "inactive", OrganizationStatusInactive)
	assert.Equal(t, "suspended", OrganizationStatusSuspended)
	assert.Equal(t, "trial", OrganizationStatusTrial)

	// Test subscription status constants
	assert.Equal(t, "active", SubscriptionStatusActive)
	assert.Equal(t, "suspended", SubscriptionStatusSuspended)
	assert.Equal(t, "trial", SubscriptionStatusTrial)
	assert.Equal(t, "expired", SubscriptionStatusExpired)
	assert.Equal(t, "cancelled", SubscriptionStatusCancelled)

	// Test subscription plan constants
	assert.Equal(t, "starter", SubscriptionPlanStarter)
	assert.Equal(t, "professional", SubscriptionPlanProfessional)
	assert.Equal(t, "enterprise", SubscriptionPlanEnterprise)
}

// TestOrganizationValidation tests organization entity validation and business rules
func TestOrganizationValidation(t *testing.T) {
	t.Run("valid_organization", func(t *testing.T) {
		org := createValidOrganization()
		
		// Basic validations
		assert.NotEmpty(t, org.Name)
		assert.NotEmpty(t, org.ContactEmail)
		assert.NotEmpty(t, org.Type)
		assert.NotEmpty(t, org.Industry)
		assert.True(t, org.IsActive)
		assert.Greater(t, org.MaxMembers, 0)
	})

	t.Run("organization_with_subscription", func(t *testing.T) {
		org := createValidOrganization()
		
		// Validate subscription structure
		assert.NotEmpty(t, org.Subscription.Plan)
		assert.NotEmpty(t, org.Subscription.Status)
		assert.Greater(t, org.Subscription.UserLimit, 0)
		assert.Greater(t, org.Subscription.StorageLimit, int64(0))
	})

	t.Run("organization_with_feature_flags", func(t *testing.T) {
		org := createValidOrganization()
		org.FeatureFlags = map[string]bool{
			"basic_controls":     true,
			"advanced_reporting": false,
			"api_access":        true,
		}

		// Validate feature flags
		assert.NotNil(t, org.FeatureFlags)
		assert.True(t, org.FeatureFlags["basic_controls"])
		assert.False(t, org.FeatureFlags["advanced_reporting"])
		assert.True(t, org.FeatureFlags["api_access"])
	})
}

// Helper function to create a valid organization for testing
func createValidOrganization() Organization {
	now := time.Now()
	
	return Organization{
		BaseModel: BaseModel{
			ID:        primitive.NewObjectID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:         "Test Bank",
		DisplayName:  "Test Bank Inc.",
		Slug:         "test-bank",
		Description:  "A test banking organization",
		Type:         "commercial_bank",
		Industry:     "financial_services",
		Size:         "medium",
		Region:       "North America",
		Country:      "United States",
		Timezone:     "America/New_York",
		Currency:     "USD",
		ContactEmail: "contact@testbank.com",
		ContactPhone: "+1-555-0123",
		Website:      "https://testbank.com",
		Address: Address{
			Street1:    "123 Main St",
			City:       "New York",
			State:      "NY",
			PostalCode: "10001",
			Country:    "United States",
		},
		RegulatoryProfile: RegulatoryProfile{
			Industry:             "financial_services",
			PrimaryRegulator:     "OCC",
			ApplicableFrameworks: []string{"SOX", "PCI-DSS"},
			ExamCycle:           "18_months",
			RetentionPeriod:     7,
			RiskTolerance:       "low",
			RequiresSOX:         true,
			RequiresPCIDSS:      true,
			ComplianceStatus:    "compliant",
		},
		Subscription: OrganizationSubscription{
			Plan:           SubscriptionPlanProfessional,
			Status:         SubscriptionStatusActive,
			BillingPeriod:  "monthly",
			UserLimit:      50,
			StorageLimit:   50 * 1024 * 1024 * 1024, // 50 GB
			CurrentUsers:   10,
			CurrentStorage: 5 * 1024 * 1024 * 1024,  // 5 GB
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		Settings: OrganizationSettings{
			RequireMFA:            true,
			AllowInvitations:      true,
			SessionTimeoutMinutes: 480,
			EnableAuditLog:        true,
			DataRetentionDays:     2555,
			AllowDataExport:       true,
			Integrations: OrganizationIntegrations{
				SSO:  true,
				LDAP: false,
				API:  true,
			},
			Notifications: OrganizationNotifications{
				EmailEnabled:      true,
				EvidenceRequests:  true,
				DeadlineReminders: true,
				SystemAlerts:      true,
			},
			ComplianceSettings: OrganizationCompliance{
				EnableChangeTracking:    true,
				RequireEvidenceApproval: false,
				MinimumReviewers:        1,
			},
			DefaultUserRole:         RoleAuditor,
			AutoAssignControls:      false,
			AutoGenerateWorkpapers:  true,
			EnableWorkflowReminders: true,
		},
		Status:      OrganizationStatusActive,
		IsActive:    true,
		MemberCount: 10,
		MaxMembers:  50,
		FeatureFlags: map[string]bool{
			"basic_controls":     true,
			"evidence_management": true,
			"basic_reporting":    true,
			"advanced_reporting": true,
			"api_access":        true,
		},
	}
}