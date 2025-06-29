# Database Migration System

This document describes the enhanced database migration system for the GoEdu Control Testing Platform.

## Overview

The migration system provides version-controlled database schema evolution with rollback capabilities, audit trails, and production-safe migration execution. It ensures database consistency across development, staging, and production environments.

## Features

- **Version Control**: Sequential migration versioning with unique identifiers
- **Rollback Support**: Safe rollback of migrations with automatic cleanup
- **Audit Trail**: Complete history of applied migrations with timestamps
- **Idempotent Operations**: Safe to run multiple times without side effects
- **Concurrent Protection**: Prevents concurrent migration execution conflicts
- **Index Management**: Automated creation and maintenance of database indexes

## Architecture

### Components

1. **Migration Manager** (`internal/migrations/migration.go`)
   - Orchestrates migration execution
   - Tracks migration history
   - Handles rollback operations

2. **Migration Definitions** (`internal/migrations/migrations.go`)
   - Contains all migration definitions
   - Implements up/down functions for each migration

3. **Migration Tool** (`cmd/migrate/main.go`)
   - CLI tool for running migrations
   - Supports both up and down operations

### Database Schema

The system creates a `migrations` collection to track applied migrations:

```javascript
{
  version: 1,                    // Migration version number
  description: "Initial setup",  // Human-readable description
  applied_at: ISODate(),         // When migration was applied
  checksum: "1-Initial setup"    // Simple integrity check
}
```

## Usage

### Running Migrations

```bash
# Apply all pending migrations
make migrate

# Or run directly
./bin/migrate
```

### Creating New Migrations

1. Add a new migration function to `internal/migrations/migrations.go`:

```go
func migration004ExampleMigration() Migration {
    return Migration{
        Version:     4,
        Description: "Add new field to users collection",
        Up: func(ctx context.Context, db *database.Client) error {
            // Forward migration logic
            collection := db.Collection("users")
            _, err := collection.UpdateMany(
                ctx,
                bson.M{},
                bson.M{"$set": bson.M{"new_field": "default_value"}},
            )
            return err
        },
        Down: func(ctx context.Context, db *database.Client) error {
            // Rollback migration logic
            collection := db.Collection("users")
            _, err := collection.UpdateMany(
                ctx,
                bson.M{},
                bson.M{"$unset": bson.M{"new_field": ""}},
            )
            return err
        },
    }
}
```

2. Add the migration to the `getAllMigrations()` function:

```go
func getAllMigrations() []Migration {
    return []Migration{
        migration001InitialIndexes(),
        migration002AuditIndexes(),
        migration003OptimizeQueries(),
        migration004ExampleMigration(), // <- Add here
        // Add new migrations here...
    }
}
```

### Migration Guidelines

1. **Sequential Versioning**: Always increment version numbers sequentially
2. **Descriptive Names**: Use clear, descriptive migration descriptions
3. **Rollback Implementation**: Always implement the Down function
4. **Idempotent Operations**: Ensure migrations can be run multiple times safely
5. **Testing**: Test both up and down migrations thoroughly

## Current Migrations

### Migration 001: Initial Indexes
- **Version**: 1
- **Description**: Create initial database indexes for optimal performance
- **Purpose**: Sets up all basic indexes required for application functionality

**Created Indexes**:
- Users: email (unique), organization_id + role, created_at
- Organizations: slug (unique), created_at  
- Controls: organization_id + control_id (unique), organization_id + status, framework + category, created_at
- Testing Cycles: organization_id + cycle_id (unique), organization_id + status, start_date + end_date
- Evidence Requests: organization_id + control_id, assignee_id + status, due_date, created_at
- Audit Logs: user_id + timestamp, organization_id + timestamp, action + timestamp, resource_id + timestamp, correlation_id

### Migration 002: Audit Indexes (Placeholder)
- **Version**: 2
- **Description**: Add specialized indexes for audit trail and compliance queries
- **Purpose**: Performance optimization for compliance reporting
- **Status**: Placeholder for future implementation

### Migration 003: Query Optimization (Placeholder)
- **Version**: 3
- **Description**: Add performance optimization indexes based on query patterns
- **Purpose**: Performance improvements based on production usage patterns
- **Status**: Placeholder for future implementation

## API Reference

### MigrationManager Methods

#### `Up(ctx context.Context) error`
Applies all pending migrations to the database.

#### `Down(ctx context.Context) error`
Rolls back the last applied migration.

#### `GetVersion(ctx context.Context) (int, error)`
Returns the current database schema version.

#### `GetMigrationHistory(ctx context.Context) ([]MigrationRecord, error)`
Returns the complete migration history.

## Production Considerations

### Pre-Migration Checklist

1. **Backup Database**: Always backup before running migrations in production
2. **Test Environment**: Run migrations in staging environment first
3. **Rollback Plan**: Ensure rollback migration is tested and working
4. **Maintenance Window**: Schedule migrations during low-traffic periods
5. **Monitoring**: Monitor application health during and after migration

### Error Handling

- Migrations fail atomically - partial application is prevented
- Failed migrations must be fixed and re-run
- Migration history tracks both successful and failed attempts
- Rollback capability allows recovery from problematic migrations

### Performance Considerations

- Index creation can be slow on large collections
- Consider creating indexes with `background: true` for large datasets
- Monitor database performance during migration execution
- Plan for potential downtime during schema changes

## Troubleshooting

### Common Issues

1. **Migration Timeout**
   - Increase context timeout for large operations
   - Consider breaking large migrations into smaller chunks

2. **Duplicate Key Errors**
   - Check for unique constraint violations
   - Ensure data cleanup before index creation

3. **Version Conflicts**
   - Ensure sequential version numbering
   - Check for missing migrations in deployment

4. **Rollback Failures**
   - Verify rollback logic is correct
   - Check for data dependencies that prevent rollback

### Debugging

Enable debug logging for detailed migration information:

```bash
GOEDU_LOGGER_LEVEL=debug ./bin/migrate
```

Check migration collection directly in MongoDB:

```javascript
db.migrations.find().sort({version: 1})
```

## Security Considerations

- Migration tool requires database connection credentials
- Use environment variables for sensitive configuration
- Restrict migration tool access in production environments
- Audit migration execution with proper logging

## Future Enhancements

1. **Migration Validation**: Pre-migration validation checks
2. **Parallel Execution**: Safe parallel migration execution
3. **Schema Diffing**: Automatic schema difference detection
4. **Migration Templates**: Code generation for common migration patterns
5. **Integration Testing**: Automated migration testing in CI/CD pipeline