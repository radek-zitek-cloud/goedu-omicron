// Package main provides database seeding tool for the GoEdu Control Testing Platform.
// It creates initial data for development and testing environments.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/config"
	"github.com/radek-zitek-cloud/goedu-omicron/be/internal/models"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/database"
	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/logger"
)

// main is the entry point for the seeding tool.
// It connects to the database and creates initial development data.
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Only run seeding in development environment
	if cfg.App.Environment == "production" {
		fmt.Println("Seeding is not allowed in production environment")
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(&cfg.Logger)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	log.Info("Starting database seeding...",
		logger.String("database", cfg.Database.Database),
		logger.String("environment", cfg.App.Environment),
	)

	// Connect to database
	dbClient, err := database.NewClient(&cfg.Database, log)
	if err != nil {
		log.Error(ctx, "Failed to connect to database", err)
		os.Exit(1)
	}
	defer dbClient.Close(ctx)

	// Run seeding
	if err := runSeeding(ctx, dbClient, log); err != nil {
		log.Error(ctx, "Seeding failed", err)
		os.Exit(1)
	}

	log.Info("Database seeding completed successfully")
}

// runSeeding creates initial development data.
// This includes sample organizations, users, controls, and testing cycles.
//
// Parameters:
//   - ctx: Context for seeding operations
//   - db: Database client
//   - log: Logger instance
//
// Returns:
//   - error: Seeding error if any step fails
func runSeeding(ctx context.Context, db *database.Client, log *logger.Logger) error {
	log.Info("Creating development data...")

	// Create sample organization
	if err := createSampleOrganization(ctx, db, log); err != nil {
		return fmt.Errorf("failed to create sample organization: %w", err)
	}

	// Create sample users
	if err := createSampleUsers(ctx, db, log); err != nil {
		return fmt.Errorf("failed to create sample users: %w", err)
	}

	// Create sample controls
	if err := createSampleControls(ctx, db, log); err != nil {
		return fmt.Errorf("failed to create sample controls: %w", err)
	}

	log.Info("All seeding completed successfully")
	return nil
}

// createSampleOrganization creates a sample organization for development.
func createSampleOrganization(ctx context.Context, db *database.Client, log *logger.Logger) error {
	log.Info("Creating sample organization...")

	org := &models.Organization{
		Name:         "Sample Financial Services Inc.",
		Slug:         "sample-finance",
		Description:  "A sample financial services organization for development",
		ContactEmail: "admin@samplefinance.com",
		ContactPhone: "+1-555-0123",
		Address: models.Address{
			Street1:    "123 Financial Street",
			City:       "New York",
			State:      "NY",
			PostalCode: "10001",
			Country:    "USA",
		},
		RegulatoryProfile: models.RegulatoryProfile{
			Industry:             "Financial Services",
			Regulations:          []string{"SOX", "GDPR", "PCI-DSS"},
			ApplicableFrameworks: []string{"COSO", "COBIT"},
			AuditFrequency:       "Annual",
			RetentionPeriod:      7,
			RiskTolerance:        "Low",
		},
		Subscription: models.OrganizationSubscription{
			Plan:   "Enterprise",
			Status: "active",
		},
		FeatureFlags: map[string]bool{
			"advanced_reporting": true,
			"api_access":        true,
			"custom_workflows":  true,
		},
		Status: models.OrganizationStatusActive,
	}

	org.UpdateTimestamps()

	// Insert organization
	collection := db.Collection("organizations")
	_, err := collection.InsertOne(ctx, org)
	if err != nil {
		return fmt.Errorf("failed to insert organization: %w", err)
	}

	log.Info("Sample organization created successfully",
		logger.String("name", org.Name),
		logger.String("slug", org.Slug),
	)

	return nil
}

// createSampleUsers creates sample users for development.
func createSampleUsers(ctx context.Context, db *database.Client, log *logger.Logger) error {
	log.Info("Creating sample users...")

	// Get the organization we just created
	orgCollection := db.Collection("organizations")
	var org models.Organization
	err := orgCollection.FindOne(ctx, map[string]string{"slug": "sample-finance"}).Decode(&org)
	if err != nil {
		return fmt.Errorf("failed to find sample organization: %w", err)
	}

	users := []*models.User{
		{
			Email: "admin@samplefinance.com",
			Profile: models.UserProfile{
				FirstName:  "John",
				LastName:   "Admin",
				Title:      "System Administrator",
				Department: "IT",
			},
			Authentication: models.AuthenticationDetails{
				PasswordHash: "$2a$12$placeholder", // In real implementation, hash properly
			},
			OrganizationID: org.ID,
			Roles:          []string{models.RoleAdmin},
			IsActive:       true,
			Status:         models.UserStatusActive,
			Permissions: models.UserPermissions{
				CanManageUsers:    true,
				CanManageSettings: true,
				CanViewReports:    true,
				CanEditControls:   true,
				CanViewControls:   true,
			},
		},
		{
			Email: "auditor@samplefinance.com",
			Profile: models.UserProfile{
				FirstName:  "Jane",
				LastName:   "Auditor",
				Title:      "Senior Auditor",
				Department: "Internal Audit",
			},
			Authentication: models.AuthenticationDetails{
				PasswordHash: "$2a$12$placeholder",
			},
			OrganizationID: org.ID,
			Roles:          []string{models.RoleAuditor},
			IsActive:       true,
			Status:         models.UserStatusActive,
			Permissions: models.UserPermissions{
				CanViewControls:   true,
				CanExecuteTests:   true,
				CanUploadEvidence: true,
				CanViewReports:    true,
			},
		},
		{
			Email: "owner@samplefinance.com",
			Profile: models.UserProfile{
				FirstName:  "Mike",
				LastName:   "Owner",
				Title:      "Control Owner",
				Department: "Operations",
			},
			Authentication: models.AuthenticationDetails{
				PasswordHash: "$2a$12$placeholder",
			},
			OrganizationID: org.ID,
			Roles:          []string{models.RoleOwner},
			IsActive:       true,
			Status:         models.UserStatusActive,
			Permissions: models.UserPermissions{
				CanViewControls:   true,
				CanUploadEvidence: true,
				CanViewReports:    true,
			},
		},
	}

	userCollection := db.Collection("users")
	for _, user := range users {
		user.UpdateTimestamps()
		_, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", user.Email, err)
		}

		log.Info("Sample user created",
			logger.String("email", user.Email),
			logger.String("role", strings.Join(user.Roles, ",")),
		)
	}

	return nil
}

// createSampleControls creates sample controls for development.
func createSampleControls(ctx context.Context, db *database.Client, log *logger.Logger) error {
	log.Info("Creating sample controls...")

	// Get the organization
	orgCollection := db.Collection("organizations")
	var org models.Organization
	err := orgCollection.FindOne(ctx, map[string]string{"slug": "sample-finance"}).Decode(&org)
	if err != nil {
		return fmt.Errorf("failed to find sample organization: %w", err)
	}

	controls := []*models.Control{
		{
			OrganizationID:   org.ID,
			ControlID:        "AC-001",
			Title:            "User Access Review",
			Description:      "Quarterly review of user access rights and permissions",
			Framework:        "COSO",
			Category:         "Access Control",
			RiskLevel:        "High",
			Importance:       "Critical",
			ControlType:      "Detective",
			ControlFrequency: "Quarterly",
			Owner:            "IT Security Manager",
			Process:          "Identity and Access Management",
			Systems:          []string{"Active Directory", "SAP", "Database"},
			TestingProcedure: "Review user access reports and validate segregation of duties",
			SampleSize:       25,
			EvidenceTypes:    []string{"User Access Report", "Approval Documentation"},
			Status:           models.ControlStatusActive,
		},
		{
			OrganizationID:   org.ID,
			ControlID:        "FIN-002",
			Title:            "Journal Entry Approval",
			Description:      "All journal entries above $10,000 require managerial approval",
			Framework:        "COSO",
			Category:         "Financial Reporting",
			RiskLevel:        "High",
			Importance:       "Critical",
			ControlType:      "Preventive",
			ControlFrequency: "Daily",
			Owner:            "Finance Manager",
			Process:          "Financial Close",
			Systems:          []string{"ERP System", "GL System"},
			TestingProcedure: "Test sample of journal entries for proper approval",
			SampleSize:       40,
			EvidenceTypes:    []string{"Journal Entry Report", "Approval Workflows"},
			Status:           models.ControlStatusActive,
		},
	}

	controlCollection := db.Collection("controls")
	for _, control := range controls {
		control.UpdateTimestamps()
		_, err := controlCollection.InsertOne(ctx, control)
		if err != nil {
			return fmt.Errorf("failed to insert control %s: %w", control.ControlID, err)
		}

		log.Info("Sample control created",
			logger.String("control_id", control.ControlID),
			logger.String("title", control.Title),
		)
	}

	return nil
}