// Package migrations defines all database migrations for the GoEdu Control Testing Platform.
// This file contains the ordered list of all migrations that need to be applied to the database.
package migrations

import (
	"context"

	"github.com/radek-zitek-cloud/goedu-omicron/be/pkg/database"
)

// getAllMigrations returns all available migrations in the system.
// Migrations should be added to this list in version order.
//
// Returns:
//   - []Migration: Complete list of all available migrations
//
// Example:
//   migrations := getAllMigrations()
//   for _, migration := range migrations {
//       log.Info("Available migration", zap.Int("version", migration.Version))
//   }
func getAllMigrations() []Migration {
	return []Migration{
		migration001InitialIndexes(),
		migration002AuditIndexes(),
		migration003OptimizeQueries(),
		// Add new migrations here...
	}
}

// migration001InitialIndexes creates the initial database indexes.
// This migration creates all the basic indexes required for the application.
func migration001InitialIndexes() Migration {
	return Migration{
		Version:     1,
		Description: "Create initial database indexes for optimal performance",
		Up: func(ctx context.Context, db *database.Client) error {
			// This uses the existing CreateIndexes method from the database client
			return db.CreateIndexes(ctx)
		},
		Down: func(ctx context.Context, db *database.Client) error {
			// Rollback: drop all custom indexes (keep only _id indexes)
			collections := []string{"users", "organizations", "controls", "testing_cycles", "evidence_requests", "audit_logs"}
			
			for _, collectionName := range collections {
				collection := db.Collection(collectionName)
				
				// Get all indexes
				cursor, err := collection.Indexes().List(ctx)
				if err != nil {
					return err
				}
				
				// Drop all indexes except _id_
				var indexes []map[string]interface{}
				if err := cursor.All(ctx, &indexes); err != nil {
					return err
				}
				
				for _, index := range indexes {
					if name, ok := index["name"].(string); ok && name != "_id_" {
						if _, err := collection.Indexes().DropOne(ctx, name); err != nil {
							return err
						}
					}
				}
			}
			
			return nil
		},
	}
}

// migration002AuditIndexes creates additional indexes for audit trail optimization.
// This migration adds specialized indexes for compliance and audit queries.
func migration002AuditIndexes() Migration {
	return Migration{
		Version:     2,
		Description: "Add specialized indexes for audit trail and compliance queries",
		Up: func(ctx context.Context, db *database.Client) error {
			// Create additional audit-specific indexes
			// This would be implemented when audit requirements are finalized
			// For now, this is a placeholder
			return nil
		},
		Down: func(ctx context.Context, db *database.Client) error {
			// Rollback: remove audit-specific indexes
			// This would remove the indexes created in the Up function
			return nil
		},
	}
}

// migration003OptimizeQueries creates indexes for query optimization based on usage patterns.
// This migration adds indexes discovered through performance testing and monitoring.
func migration003OptimizeQueries() Migration {
	return Migration{
		Version:     3,
		Description: "Add performance optimization indexes based on query patterns",
		Up: func(ctx context.Context, db *database.Client) error {
			// Create performance optimization indexes
			// This would be implemented based on actual query patterns
			// For now, this is a placeholder
			return nil
		},
		Down: func(ctx context.Context, db *database.Client) error {
			// Rollback: remove performance optimization indexes
			return nil
		},
	}
}

// Future migration templates:
//
// func migration004ExampleMigration() Migration {
//     return Migration{
//         Version:     4,
//         Description: "Example migration description",
//         Up: func(ctx context.Context, db *database.Client) error {
//             // Forward migration logic
//             return nil
//         },
//         Down: func(ctx context.Context, db *database.Client) error {
//             // Rollback migration logic
//             return nil
//         },
//     }
// }