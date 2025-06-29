// Package services provides service layer implementations for the GoEdu Control Testing Platform.
// This file contains the organization service implementation with multi-tenancy support.
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/models"
	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/repositories"
)

// organizationService implements the OrganizationService interface.
// It provides comprehensive organization management capabilities including
// subscription management, feature flags, and multi-tenancy support.
type organizationService struct {
	orgRepo     repositories.OrganizationRepository
	userRepo    repositories.UserRepository
	auditRepo   repositories.AuditLogRepository
	cacheRepo   repositories.CacheRepository
	logger      *zap.Logger
}

// NewOrganizationService creates a new organization service with required dependencies.
//
// Parameters:
//   - orgRepo: Repository for organization data operations
//   - userRepo: Repository for user data operations
//   - auditRepo: Repository for audit logging
//   - cacheRepo: Repository for caching operations
//   - logger: Logger for service operations
//
// Returns:
//   - OrganizationService: Configured organization service instance
func NewOrganizationService(
	orgRepo repositories.OrganizationRepository,
	userRepo repositories.UserRepository,
	auditRepo repositories.AuditLogRepository,
	cacheRepo repositories.CacheRepository,
	logger *zap.Logger,
) OrganizationService {
	return &organizationService{
		orgRepo:   orgRepo,
		userRepo:  userRepo,
		auditRepo: auditRepo,
		cacheRepo: cacheRepo,
		logger:    logger,
	}
}

// CreateOrganization creates a new organization with proper validation and setup.
// This method initializes all required organization components including subscription,
// settings, and default feature flags.
//
// Parameters:
//   - ctx: Request context
//   - input: Organization creation input data
//
// Returns:
//   - *models.Organization: Created organization
//   - error: Error if creation fails
func (s *organizationService) CreateOrganization(ctx context.Context, input *CreateOrganizationInput) (*models.Organization, error) {
	// Validate input
	if err := s.validateCreateOrganizationInput(input); err != nil {
		s.logger.Warn("Invalid organization creation input",
			zap.Error(err),
			zap.String("name", input.Name),
		)
		return nil, err
	}

	// Create slug from name if not provided
	slug := s.generateSlug(input.Name)

	// Check if slug is unique
	if exists, err := s.slugExists(ctx, slug); err != nil {
		return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
	} else if exists {
		return nil, ErrOrganizationSlugExists
	}

	// Create organization entity
	org := &models.Organization{
		BaseModel: models.BaseModel{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        input.Name,
		DisplayName: input.DisplayName,
		Slug:        slug,
		Description: input.Description,
		Type:        input.Type,
		Industry:    input.Industry,
		Size:        input.Size,
		Region:      input.Region,
		Country:     input.Country,
		Timezone:    input.Timezone,
		Currency:    input.Currency,
		ContactEmail: input.ContactEmail,
		ContactPhone: input.ContactPhone,
		Website:     input.Website,
		LogoURL:     input.LogoURL,
		Status:      models.OrganizationStatusActive,
		IsActive:    true,
		MemberCount: 0,
		MaxMembers:  s.getDefaultMemberLimit(input.SubscriptionPlan),
	}

	// Set address if provided
	if input.Address != nil {
		org.Address = models.Address{
			Street1:    input.Address.Street1,
			Street2:    input.Address.Street2,
			City:       input.Address.City,
			State:      input.Address.State,
			PostalCode: input.Address.PostalCode,
			Country:    input.Address.Country,
		}
	}

	// Set regulatory profile if provided
	if input.RegulatoryProfile != nil {
		org.RegulatoryProfile = s.createRegulatoryProfile(input.RegulatoryProfile)
	}

	// Initialize subscription
	org.Subscription = s.createDefaultSubscription(input.SubscriptionPlan)

	// Initialize settings
	org.Settings = s.createDefaultSettings(input.Settings)

	// Initialize default feature flags
	org.FeatureFlags = s.createDefaultFeatureFlags(input.SubscriptionPlan)

	// Create organization in database
	if err := s.orgRepo.Create(ctx, org); err != nil {
		s.logger.Error("Failed to create organization in database",
			zap.Error(err),
			zap.String("name", org.Name),
			zap.String("slug", org.Slug),
		)
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Log organization creation
	s.logOrganizationEvent(ctx, org.ID, "organization_created", map[string]interface{}{
		"name":              org.Name,
		"type":              org.Type,
		"industry":          org.Industry,
		"subscription_plan": org.Subscription.Plan,
	})

	s.logger.Info("Organization created successfully",
		zap.String("organization_id", org.ID.Hex()),
		zap.String("name", org.Name),
		zap.String("slug", org.Slug),
	)

	return org, nil
}

// GetOrganization retrieves an organization by ID with caching support.
//
// Parameters:
//   - ctx: Request context
//   - id: Organization ID
//
// Returns:
//   - *models.Organization: Organization entity
//   - error: Error if not found or retrieval fails
func (s *organizationService) GetOrganization(ctx context.Context, id string) (*models.Organization, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("org:%s", id)
	var org models.Organization
	if err := s.cacheRepo.Get(ctx, cacheKey, &org); err == nil {
		return &org, nil
	}

	// Get from database
	orgEntity, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Warn("Failed to get organization",
			zap.Error(err),
			zap.String("organization_id", id),
		)
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	// Cache for future requests
	if err := s.cacheRepo.Set(ctx, cacheKey, orgEntity, 15*time.Minute); err != nil {
		s.logger.Warn("Failed to cache organization",
			zap.Error(err),
			zap.String("organization_id", id),
		)
	}

	return orgEntity, nil
}

// GetOrganizationBySlug retrieves an organization by its unique slug.
//
// Parameters:
//   - ctx: Request context
//   - slug: Organization slug
//
// Returns:
//   - *models.Organization: Organization entity
//   - error: Error if not found or retrieval fails
func (s *organizationService) GetOrganizationBySlug(ctx context.Context, slug string) (*models.Organization, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("org:slug:%s", slug)
	var org models.Organization
	if err := s.cacheRepo.Get(ctx, cacheKey, &org); err == nil {
		return &org, nil
	}

	// Get from database
	orgEntity, err := s.orgRepo.GetBySlug(ctx, slug)
	if err != nil {
		s.logger.Warn("Failed to get organization by slug",
			zap.Error(err),
			zap.String("slug", slug),
		)
		return nil, fmt.Errorf("failed to get organization by slug: %w", err)
	}

	// Cache for future requests
	if err := s.cacheRepo.Set(ctx, cacheKey, orgEntity, 15*time.Minute); err != nil {
		s.logger.Warn("Failed to cache organization by slug",
			zap.Error(err),
			zap.String("slug", slug),
		)
	}

	return orgEntity, nil
}

// GetFeatureFlags retrieves feature flags for an organization with caching.
//
// Parameters:
//   - ctx: Request context
//   - orgID: Organization ID
//
// Returns:
//   - map[string]bool: Feature flags mapping
//   - error: Error if retrieval fails
func (s *organizationService) GetFeatureFlags(ctx context.Context, orgID string) (map[string]bool, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("org:flags:%s", orgID)
	var flags map[string]bool
	if err := s.cacheRepo.Get(ctx, cacheKey, &flags); err == nil {
		return flags, nil
	}

	// Get feature flags from database
	flags, err := s.orgRepo.GetFeatureFlags(ctx, orgID)
	if err != nil {
		s.logger.Warn("Failed to get feature flags",
			zap.Error(err),
			zap.String("organization_id", orgID),
		)
		return nil, fmt.Errorf("failed to get feature flags: %w", err)
	}

	// Cache for future requests
	if err := s.cacheRepo.Set(ctx, cacheKey, flags, 10*time.Minute); err != nil {
		s.logger.Warn("Failed to cache feature flags",
			zap.Error(err),
			zap.String("organization_id", orgID),
		)
	}

	return flags, nil
}

// UpdateFeatureFlag updates a single feature flag for an organization.
//
// Parameters:
//   - ctx: Request context
//   - orgID: Organization ID
//   - flag: Feature flag name
//   - enabled: Whether the feature is enabled
//
// Returns:
//   - error: Error if update fails
func (s *organizationService) UpdateFeatureFlag(ctx context.Context, orgID, flag string, enabled bool) error {
	// Update in database
	if err := s.orgRepo.UpdateFeatureFlag(ctx, orgID, flag, enabled); err != nil {
		s.logger.Error("Failed to update feature flag",
			zap.Error(err),
			zap.String("organization_id", orgID),
			zap.String("flag", flag),
			zap.Bool("enabled", enabled),
		)
		return fmt.Errorf("failed to update feature flag: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("org:flags:%s", orgID)
	if err := s.cacheRepo.Delete(ctx, cacheKey); err != nil {
		s.logger.Warn("Failed to invalidate feature flags cache",
			zap.Error(err),
			zap.String("organization_id", orgID),
		)
	}

	// Log feature flag change
	s.logOrganizationEvent(ctx, primitive.ObjectID{}, "feature_flag_updated", map[string]interface{}{
		"organization_id": orgID,
		"flag":           flag,
		"enabled":        enabled,
	})

	s.logger.Info("Feature flag updated",
		zap.String("organization_id", orgID),
		zap.String("flag", flag),
		zap.Bool("enabled", enabled),
	)

	return nil
}

// ValidateOrganizationAccess validates that a user has access to an organization.
//
// Parameters:
//   - ctx: Request context
//   - userID: User ID
//   - orgID: Organization ID
//
// Returns:
//   - error: Error if validation fails or access is denied
func (s *organizationService) ValidateOrganizationAccess(ctx context.Context, userID, orgID string) error {
	// Get user to check organization association
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user belongs to the organization
	if user.OrganizationID.Hex() != orgID {
		return ErrUserNotInOrganization
	}

	// Check if user is active
	if !user.IsActive || user.Status != models.UserStatusActive {
		return ErrUserNotActive
	}

	return nil
}

// IsFeatureEnabled checks if a specific feature is enabled for an organization.
//
// Parameters:
//   - ctx: Request context
//   - orgID: Organization ID
//   - feature: Feature name to check
//
// Returns:
//   - bool: Whether the feature is enabled
//   - error: Error if check fails
func (s *organizationService) IsFeatureEnabled(ctx context.Context, orgID, feature string) (bool, error) {
	flags, err := s.GetFeatureFlags(ctx, orgID)
	if err != nil {
		return false, err
	}

	enabled, exists := flags[feature]
	return exists && enabled, nil
}

// Helper methods for organization service

// validateCreateOrganizationInput validates the input for creating an organization.
func (s *organizationService) validateCreateOrganizationInput(input *CreateOrganizationInput) error {
	if input == nil {
		return ErrInvalidInput
	}

	if strings.TrimSpace(input.Name) == "" {
		return ErrOrganizationNameRequired
	}

	if strings.TrimSpace(input.ContactEmail) == "" {
		return ErrContactEmailRequired
	}

	if strings.TrimSpace(input.Type) == "" {
		return ErrOrganizationTypeRequired
	}

	return nil
}

// generateSlug creates a URL-friendly slug from the organization name.
func (s *organizationService) generateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// Remove any non-alphanumeric characters except hyphens
	result := ""
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result += string(char)
		}
	}
	return result
}

// slugExists checks if a slug already exists in the database.
func (s *organizationService) slugExists(ctx context.Context, slug string) (bool, error) {
	_, err := s.orgRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// getDefaultMemberLimit returns the default member limit for a subscription plan.
func (s *organizationService) getDefaultMemberLimit(plan string) int {
	switch plan {
	case models.SubscriptionPlanStarter:
		return 10
	case models.SubscriptionPlanProfessional:
		return 50
	case models.SubscriptionPlanEnterprise:
		return 1000
	default:
		return 5 // Free tier default
	}
}

// createRegulatoryProfile creates a regulatory profile from input data.
func (s *organizationService) createRegulatoryProfile(input *CreateRegulatoryProfileInput) models.RegulatoryProfile {
	return models.RegulatoryProfile{
		Industry:             input.Industry,
		PrimaryRegulator:     input.PrimaryRegulator,
		ApplicableFrameworks: input.ApplicableFrameworks,
		ExamCycle:           input.ExamCycle,
		RetentionPeriod:     input.RetentionPeriod,
		RiskTolerance:       input.RiskTolerance,
		RequiresSOX:         input.RequiresSOX,
		RequiresPCIDSS:      input.RequiresPCIDSS,
		RequiresFFIEC:       input.RequiresFFIEC,
		RequiresBaselIII:    input.RequiresBaselIII,
		ComplianceStatus:    "pending",
		CustomRequirements:  make(map[string]interface{}),
	}
}

// createDefaultSubscription creates a default subscription configuration.
func (s *organizationService) createDefaultSubscription(plan string) models.OrganizationSubscription {
	if plan == "" {
		plan = models.SubscriptionPlanStarter
	}

	now := time.Now()
	return models.OrganizationSubscription{
		Plan:               plan,
		Status:             models.SubscriptionStatusTrial,
		BillingPeriod:      "monthly",
		TrialStart:         now,
		TrialEnd:           now.AddDate(0, 0, 30), // 30-day trial
		IsInTrial:          true,
		UserLimit:          s.getDefaultMemberLimit(plan),
		StorageLimit:       s.getDefaultStorageLimit(plan),
		CurrentUsers:       0,
		CurrentStorage:     0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// createDefaultSettings creates default organization settings.
func (s *organizationService) createDefaultSettings(input *CreateOrganizationSettingsInput) models.OrganizationSettings {
	settings := models.OrganizationSettings{
		RequireMFA:            false,
		AllowInvitations:      true,
		SessionTimeoutMinutes: 480, // 8 hours
		EnableAuditLog:        true,
		DataRetentionDays:     2555, // ~7 years for compliance
		AllowDataExport:       true,
		CustomBranding:        false,
		Integrations: models.OrganizationIntegrations{
			SSO:             false,
			LDAP:            false,
			API:             true,
			APIKeyEnabled:   false,
			WebhooksEnabled: false,
		},
		Notifications: models.OrganizationNotifications{
			EmailEnabled:      true,
			SMSEnabled:        false,
			WebEnabled:        true,
			DigestFrequency:   "daily",
			EvidenceRequests:  true,
			DeadlineReminders: true,
			SystemAlerts:      true,
		},
		ComplianceSettings: models.OrganizationCompliance{
			AutoGenerateReports:      false,
			RequireDigitalSignature:  false,
			EnableChangeTracking:     true,
			RequireEvidenceApproval:  false,
			MinimumReviewers:         1,
			AutoSubmitReports:        false,
		},
		DefaultUserRole:         models.RoleViewer,
		AutoAssignControls:      false,
		AutoGenerateWorkpapers:  false,
		EnableWorkflowReminders: true,
		CustomSettings:          make(map[string]interface{}),
	}

	// Override with provided input
	if input != nil {
		if input.SessionTimeoutMinutes > 0 {
			settings.SessionTimeoutMinutes = input.SessionTimeoutMinutes
		}
		if input.DataRetentionDays > 0 {
			settings.DataRetentionDays = input.DataRetentionDays
		}
		settings.RequireMFA = input.RequireMFA
		settings.AllowInvitations = input.AllowInvitations
		settings.EnableAuditLog = input.EnableAuditLog
		settings.AllowDataExport = input.AllowDataExport
	}

	return settings
}

// createDefaultFeatureFlags creates default feature flags for a subscription plan.
func (s *organizationService) createDefaultFeatureFlags(plan string) map[string]bool {
	flags := map[string]bool{
		"basic_controls":        true,
		"evidence_management":   true,
		"basic_reporting":       true,
		"user_management":       true,
		"audit_trail":          true,
	}

	switch plan {
	case models.SubscriptionPlanProfessional:
		flags["advanced_reporting"] = true
		flags["custom_workflows"] = true
		flags["api_access"] = true
		flags["sso_integration"] = true

	case models.SubscriptionPlanEnterprise:
		flags["advanced_reporting"] = true
		flags["custom_workflows"] = true
		flags["api_access"] = true
		flags["sso_integration"] = true
		flags["advanced_analytics"] = true
		flags["custom_integrations"] = true
		flags["white_labeling"] = true
		flags["priority_support"] = true
	}

	return flags
}

// getDefaultStorageLimit returns the default storage limit for a subscription plan.
func (s *organizationService) getDefaultStorageLimit(plan string) int64 {
	switch plan {
	case models.SubscriptionPlanStarter:
		return 5 * 1024 * 1024 * 1024 // 5 GB
	case models.SubscriptionPlanProfessional:
		return 50 * 1024 * 1024 * 1024 // 50 GB
	case models.SubscriptionPlanEnterprise:
		return 500 * 1024 * 1024 * 1024 // 500 GB
	default:
		return 1 * 1024 * 1024 * 1024 // 1 GB for free tier
	}
}

// logOrganizationEvent logs an organization-related event for audit purposes.
func (s *organizationService) logOrganizationEvent(ctx context.Context, orgID primitive.ObjectID, action string, metadata map[string]interface{}) {
	auditEntry := &models.AuditLog{
		ID:             primitive.NewObjectID(),
		Timestamp:      time.Now(),
		OrganizationID: orgID,
		Action:         action,
		ResourceType:   "organization",
		Success:        true,
		Metadata:       metadata,
	}

	if err := s.auditRepo.Create(ctx, auditEntry); err != nil {
		s.logger.Warn("Failed to log organization event",
			zap.Error(err),
			zap.String("action", action),
			zap.String("organization_id", orgID.Hex()),
		)
	}
}

// Service errors
var (
	ErrInvalidInput              = errors.New("invalid input provided")
	ErrOrganizationNameRequired  = errors.New("organization name is required")
	ErrContactEmailRequired      = errors.New("contact email is required")
	ErrOrganizationTypeRequired  = errors.New("organization type is required")
	ErrOrganizationSlugExists    = errors.New("organization slug already exists")
	ErrUserNotInOrganization     = errors.New("user does not belong to organization")
	ErrUserNotActive             = errors.New("user is not active")
)

// Placeholder implementations for remaining OrganizationService methods
// These would be implemented based on specific business requirements

func (s *organizationService) UpdateOrganization(ctx context.Context, id string, input *UpdateOrganizationInput) (*models.Organization, error) {
	// Implementation would update organization and invalidate cache
	return nil, errors.New("not implemented")
}

func (s *organizationService) DeleteOrganization(ctx context.Context, id string) error {
	// Implementation would soft delete organization
	return errors.New("not implemented")
}

func (s *organizationService) ListOrganizations(ctx context.Context, filter *OrganizationFilter) (*OrganizationConnection, error) {
	// Implementation would list organizations with filtering and pagination
	return nil, errors.New("not implemented")
}

func (s *organizationService) GetActiveOrganizations(ctx context.Context, limit, offset int) ([]*models.Organization, error) {
	// Implementation would get active organizations
	return nil, errors.New("not implemented")
}

func (s *organizationService) UpdateSubscription(ctx context.Context, orgID string, subscription *models.OrganizationSubscription) error {
	// Implementation would update subscription
	return errors.New("not implemented")
}

func (s *organizationService) GetSubscriptionStatus(ctx context.Context, orgID string) (*models.OrganizationSubscription, error) {
	// Implementation would get subscription status
	return nil, errors.New("not implemented")
}

func (s *organizationService) UpgradeSubscription(ctx context.Context, orgID, newPlan string) error {
	// Implementation would upgrade subscription
	return errors.New("not implemented")
}

func (s *organizationService) DowngradeSubscription(ctx context.Context, orgID, newPlan string) error {
	// Implementation would downgrade subscription
	return errors.New("not implemented")
}

func (s *organizationService) CancelSubscription(ctx context.Context, orgID string) error {
	// Implementation would cancel subscription
	return errors.New("not implemented")
}

func (s *organizationService) RenewSubscription(ctx context.Context, orgID string) error {
	// Implementation would renew subscription
	return errors.New("not implemented")
}

func (s *organizationService) BulkUpdateFeatureFlags(ctx context.Context, orgID string, flags map[string]bool) error {
	// Implementation would bulk update feature flags
	return errors.New("not implemented")
}

func (s *organizationService) GetSettings(ctx context.Context, orgID string) (*models.OrganizationSettings, error) {
	// Implementation would get organization settings
	return nil, errors.New("not implemented")
}

func (s *organizationService) UpdateSettings(ctx context.Context, orgID string, settings *models.OrganizationSettings) error {
	// Implementation would update organization settings
	return errors.New("not implemented")
}

func (s *organizationService) UpdatePartialSettings(ctx context.Context, orgID string, updates map[string]interface{}) error {
	// Implementation would update partial settings
	return errors.New("not implemented")
}

func (s *organizationService) GetMemberCount(ctx context.Context, orgID string) (int, error) {
	// Implementation would get member count
	return 0, errors.New("not implemented")
}

func (s *organizationService) UpdateMemberCount(ctx context.Context, orgID string, count int) error {
	// Implementation would update member count
	return errors.New("not implemented")
}

func (s *organizationService) GetMemberLimits(ctx context.Context, orgID string) (current, max int, error) {
	// Implementation would get member limits
	return 0, 0, errors.New("not implemented")
}

func (s *organizationService) CanAddMember(ctx context.Context, orgID string) (bool, error) {
	// Implementation would check if can add member
	return false, errors.New("not implemented")
}

func (s *organizationService) ValidateOrganization(ctx context.Context, org *models.Organization) error {
	// Implementation would validate organization
	return errors.New("not implemented")
}

func (s *organizationService) CheckComplianceRequirements(ctx context.Context, orgID string) (*ComplianceStatus, error) {
	// Implementation would check compliance requirements
	return nil, errors.New("not implemented")
}

func (s *organizationService) UpdateRegulatoryProfile(ctx context.Context, orgID string, profile *models.RegulatoryProfile) error {
	// Implementation would update regulatory profile
	return errors.New("not implemented")
}

func (s *organizationService) GetOrganizationStats(ctx context.Context, orgID string) (*OrganizationStats, error) {
	// Implementation would get organization statistics
	return nil, errors.New("not implemented")
}

func (s *organizationService) GetUsageMetrics(ctx context.Context, orgID string, timeRange *TimeRange) (*UsageMetrics, error) {
	// Implementation would get usage metrics
	return nil, errors.New("not implemented")
}

func (s *organizationService) GetOrganizationContext(ctx context.Context, userID string) (*OrganizationContext, error) {
	// Implementation would get organization context for user
	return nil, errors.New("not implemented")
}

func (s *organizationService) GetUserOrganizations(ctx context.Context, userID string) ([]*models.Organization, error) {
	// Implementation would get user organizations
	return nil, errors.New("not implemented")
}