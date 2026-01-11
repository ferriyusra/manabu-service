# Database Migrations

This directory contains SQL migration files for the Manabu Service database.

## Migration Tools

### Available Tools

1. **migrate.go** - Run SQL migration files
2. **check_constraints.go** - Check database constraints status

### How to Run Migrations

#### Method 1: Using Go Script (Recommended)

```bash
# Run a specific migration file
go run tools/migrate.go <migration_file>

# Example
go run tools/migrate.go 002_rename_users_uuid_constraint.sql
```

#### Method 2: Using Windows Batch Script

```bash
# Run migration using batch file
migrate.bat <migration_file>

# Example
migrate.bat 002_rename_users_uuid_constraint.sql
```

#### Method 3: Using psql Command Line

```bash
# Direct execution via psql
psql -U postgres -d manabu -f migrations/<migration_file>

# With password prompt
psql -U postgres -d manabu -W -f migrations/002_rename_users_uuid_constraint.sql
```

#### Method 4: Using Database GUI (DBeaver/pgAdmin)

1. Open DBeaver or pgAdmin
2. Connect to `manabu` database
3. Open SQL Editor
4. Copy-paste migration file content
5. Execute

### Check Migration Status

To verify if a migration has been applied:

```bash
# Check constraints on users table
go run tools/check_constraints.go
```

This will show:
- All constraints on the `users` table
- Whether `uni_users_uuid` constraint exists
- Migration status

## Migration History

### 001_alter_user_vocabulary_status_user_id_to_uuid.sql ✅

**Date:** 2026-01-08
**Status:** Applied
**Description:** Changes `user_id` column from `bigint` to `UUID` in `user_vocabulary_status` table

**Changes:**
1. Adds unique constraint to `users.uuid` column (constraint name: `uni_users_uuid`)
2. Drops existing foreign key constraint `fk_user_vocabulary_status_user`
3. Drops existing unique index `idx_user_vocabulary`
4. Truncates `user_vocabulary_status` table (⚠️ removes all data)
5. Changes `user_id` column type from `bigint` to `UUID`
6. Recreates unique index on `(user_id, vocabulary_id)`
7. Recreates foreign key constraint referencing `users.uuid`

**Impact:**
- ⚠️ **DATA LOSS**: All data in `user_vocabulary_status` table was deleted
- User references now use UUID for better security
- API responses expose UUID strings instead of integer IDs

**SQL:**
```sql
ALTER TABLE user_vocabulary_status
    ALTER COLUMN user_id TYPE UUID USING NULL;
```

---

### 002_rename_users_uuid_constraint.sql ✅

**Date:** 2026-01-08
**Status:** Applied (via fix_db.go)
**Description:** Renames UUID constraint to match GORM naming convention

**Changes:**
1. Renames constraint from `users_uuid_unique` to `uni_users_uuid`

**Reason:**
- GORM AutoMigrate expects constraint name pattern: `uni_{table}_{column}`
- Prevents "constraint not found" error on application startup
- Ensures compatibility with GORM's automatic schema management

**SQL:**
```sql
ALTER TABLE users RENAME CONSTRAINT users_uuid_unique TO uni_users_uuid;
```

---

## Creating New Migrations

### Naming Convention

Use the following format:
```
<number>_<description>.sql
```

Examples:
- `003_add_index_to_vocabulary.sql`
- `004_create_user_settings_table.sql`
- `005_alter_tags_add_color_field.sql`

### Migration Template

```sql
-- Migration: <Brief description>
-- Created: <YYYY-MM-DD>
-- Description: <Detailed explanation>

BEGIN;

-- Your SQL statements here
-- Example:
-- ALTER TABLE table_name ADD COLUMN new_column VARCHAR(255);

COMMIT;
```

### Best Practices

1. **Always use transactions** (`BEGIN` and `COMMIT`)
2. **Make migrations idempotent** when possible using:
   - `IF NOT EXISTS` for CREATE statements
   - `IF EXISTS` for DROP statements
   - Conditional logic with `DO $$` blocks
3. **Backup data** before running destructive migrations
4. **Test migrations** on development database first
5. **Document impact** especially for data loss or schema changes

### Example: Idempotent Migration

```sql
-- Migration: Add email_verified column to users
-- Created: 2026-01-08

BEGIN;

-- Add column only if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'users'
        AND column_name = 'email_verified'
    ) THEN
        ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;
    END IF;
END $$;

COMMIT;
```

## Current Database Schema Status

### users table
- ✅ `uuid` column has unique constraint (`uni_users_uuid`)
- ✅ Compatible with GORM AutoMigrate
- Model: `UUID uuid.UUID` with `gorm:"type:uuid;not null;unique"`

### user_vocabulary_status table
- ✅ `user_id` column type: `UUID`
- ✅ Foreign key references `users.uuid`
- ✅ Unique index on `(user_id, vocabulary_id)`
- Model: `UserID uuid.UUID` with `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_vocabulary"`

## Code Changes After Migration

When adding migrations that change schema, remember to update:

### 1. Models (`domain/models/*.go`)
Update GORM struct tags to match database schema:
```go
type User struct {
    UUID uuid.UUID `gorm:"type:uuid;not null;unique"`
}
```

### 2. Repositories (`repositories/*/*.go`)
Update queries to handle new column types:
```go
// For UUID columns, use ::uuid casting
Where("user_id = ?::uuid", userID)
```

### 3. DTOs (`domain/dto/*.go`)
Update request/response structures:
```go
type UserResponse struct {
    UserID string `json:"userId"` // Expose UUID as string
}
```

### 4. Services (`services/*/*.go`)
Update business logic to work with new types:
```go
// Pass UUID string to repository
userLogin.UUID.String()
```

## Rollback Strategy

⚠️ **Important:** These migrations don't have automatic rollback.

To rollback manually:

1. **Identify changes** from migration file
2. **Write reverse SQL** (example: `ALTER TABLE ... DROP COLUMN ...`)
3. **Backup data first** before rollback
4. **Test on dev** environment first

### Example Rollback for Migration 001

```sql
BEGIN;

-- Reverse: Change user_id back to bigint
ALTER TABLE user_vocabulary_status
    DROP CONSTRAINT fk_user_vocabulary_status_user;

ALTER TABLE user_vocabulary_status
    ALTER COLUMN user_id TYPE BIGINT USING NULL;

-- Recreate foreign key to users.id
ALTER TABLE user_vocabulary_status
    ADD CONSTRAINT fk_user_vocabulary_status_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE;

COMMIT;
```

⚠️ **Note:** Rollback may cause data loss. Always backup first!

## Troubleshooting

### Error: "constraint does not exist"

**Problem:** GORM can't find expected constraint name

**Solution:**
```bash
# Check actual constraint names
go run tools/check_constraints.go

# Rename constraint to match GORM convention
# Pattern: uni_{table}_{column} for unique constraints
```

### Error: "cannot cast type"

**Problem:** Type mismatch when changing column types

**Solution:** Use explicit USING clause:
```sql
-- For string to UUID
ALTER COLUMN column_name TYPE UUID USING column_name::uuid;

-- For resetting to NULL
ALTER COLUMN column_name TYPE UUID USING NULL;
```

### Error: "violates foreign key constraint"

**Problem:** Referenced data doesn't exist

**Solution:**
1. Drop foreign key constraint first
2. Modify column
3. Recreate constraint after data is fixed

## Safety Checklist

Before running any migration:

- [ ] Read migration file completely
- [ ] Understand what will change
- [ ] Check for data loss warnings
- [ ] Backup database (if production)
- [ ] Test on development database first
- [ ] Verify application still works after migration
- [ ] Update code (models, repos, services) if needed
- [ ] Rebuild application: `go build`
- [ ] Test all affected endpoints

## Questions?

If migration fails or you're unsure:
1. Check migration file syntax
2. Verify database connection (config.json)
3. Run `check_constraints.go` to see current state
4. Review error messages carefully
5. Test on development database first

---

**Last Updated:** 2026-01-08
**Maintainer:** Manabu Service Team
