# 📋 Quick Migration Guide

Quick reference for running database migrations in Manabu Service.

## 🚀 Quick Start

### Run a Migration

```bash
# Recommended way
go run tools/migrate.go <migration_file>

# Example
go run tools/migrate.go 002_rename_users_uuid_constraint.sql
```

### Check Migration Status

```bash
go run tools/check_constraints.go
```

## 📂 Available Migrations

| File | Status | Description |
|------|--------|-------------|
| `001_alter_user_vocabulary_status_user_id_to_uuid.sql` | ✅ Applied | Changes user_id from bigint to UUID |
| `002_rename_users_uuid_constraint.sql` | ✅ Applied | Renames constraint for GORM compatibility |
| `003_update_user_vocabulary_status_check_constraint.sql` | 🔄 Optional | Updates status constraint from ('learning', 'reviewing', 'mastered') to ('learning', 'completed') |

## 🛠️ Tools

### tools/migrate.go
Runs SQL migration files with confirmation prompt.

**Usage:**
```bash
go run tools/migrate.go <migration_file.sql>
```

**Example:**
```bash
go run tools/migrate.go 003_add_new_feature.sql
```

### tools/check_constraints.go
Checks database constraints status.

**Usage:**
```bash
go run tools/check_constraints.go
```

**Output:**
- Lists all constraints on `users` table
- Shows constraint types (PRIMARY KEY, UNIQUE, FOREIGN KEY)
- Verifies migration status

### migrate.bat (Windows)
Batch script wrapper for easier migration execution.

**Usage:**
```bash
migrate.bat <migration_file.sql>
```

## ⚠️ When to Run Migrations

### First Time Setup
**Don't run migrations!** GORM AutoMigrate will create all tables automatically on first run.

### After Schema Changes
Run migrations when:
- ✅ You pull code with new migration files
- ✅ Deploying to production/staging
- ✅ Schema needs modification after initial setup
- ✅ Constraint names need fixing

### Before Running Migration
Always:
1. ✅ **Backup database** (production)
2. ✅ **Test on dev** database first
3. ✅ **Read migration file** completely
4. ✅ **Check for data loss** warnings
5. ✅ **Verify app still works** after migration

## 📝 Creating New Migrations

### 1. Create Migration File

```bash
# Navigate to migrations folder
cd migrations

# Create new file with naming convention:
# <number>_<description>.sql
```

**Example:** `003_add_user_settings_table.sql`

### 2. Write Migration

```sql
-- Migration: Add user settings table
-- Created: 2026-01-08

BEGIN;

-- Your SQL here
CREATE TABLE IF NOT EXISTS user_settings (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(uuid),
    theme VARCHAR(20) DEFAULT 'light',
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

COMMIT;
```

### 3. Test Migration

```bash
# Test on dev database
go run tools/migrate.go 003_add_user_settings_table.sql
```

### 4. Update Model (if needed)

```go
// domain/models/user_settings.go
type UserSettings struct {
    ID        uint      `gorm:"primaryKey;autoIncrement"`
    UserID    uuid.UUID `gorm:"type:uuid;not null"`
    Theme     string    `gorm:"type:varchar(20);default:light"`
    Language  string    `gorm:"type:varchar(10);default:en"`
    CreatedAt *time.Time
    UpdatedAt *time.Time
}
```

### 5. Add to AutoMigrate

```go
// cmd/main.go
db.AutoMigrate(
    // ... existing models
    &models.UserSettings{},
)
```

## 🔧 Alternative Methods

### Using psql

```bash
# Direct execution
psql -U postgres -d manabu -f migrations/002_rename_users_uuid_constraint.sql

# With password prompt
psql -U postgres -d manabu -W -f migrations/003_new_migration.sql
```

### Using Database GUI

1. Open **DBeaver** or **pgAdmin**
2. Connect to `manabu` database
3. Open SQL Editor
4. Copy migration file content
5. Execute

## ❓ Troubleshooting

### Migration Already Applied

**Error:** "relation already exists" or "column already exists"

**Solution:** Migration was already applied, skip it.

```bash
# Check what's in database
go run tools/check_constraints.go
```

### Constraint Name Mismatch

**Error:** "constraint does not exist"

**Cause:** GORM expects different constraint name

**Solution:**
1. Check actual constraint name: `go run tools/check_constraints.go`
2. Rename constraint in migration or model tag

### Type Cast Error

**Error:** "cannot cast type X to Y"

**Solution:** Use explicit USING clause:
```sql
ALTER COLUMN column_name TYPE uuid USING column_name::text::uuid;
```

### Foreign Key Violation

**Error:** "violates foreign key constraint"

**Solution:**
1. Drop foreign key first
2. Modify data/column
3. Recreate foreign key

## 📚 Full Documentation

For comprehensive guide, see:
- **[migrations/README.md](migrations/README.md)** - Complete migration documentation
- **[README.md#database-migrations](README.md#database-migrations)** - Setup guide

## 🆘 Need Help?

1. Read migration file completely
2. Check [migrations/README.md](migrations/README.md) troubleshooting section
3. Run `check_constraints.go` to see current state
4. Test on dev database first
5. Ask team if unsure

---

**Remember:** Always backup production database before running migrations! 🔒
