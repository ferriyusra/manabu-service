-- Migration: Rename users UUID constraint to match GORM naming convention
-- Created: 2026-01-08

BEGIN;

-- Rename constraint from users_uuid_unique to uni_users_uuid (GORM default naming)
ALTER TABLE users RENAME CONSTRAINT users_uuid_unique TO uni_users_uuid;

COMMIT;
