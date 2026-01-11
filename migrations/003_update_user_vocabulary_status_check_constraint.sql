-- Migration: Update user_vocabulary_status status check constraint
-- Description: Change status values from ('learning', 'reviewing', 'mastered') to ('learning', 'completed')
-- Reason: Simplified from SM-2 algorithm to simple boolean tracking system
-- Created: 2026-01-11

BEGIN;

-- Drop the old check constraint
ALTER TABLE user_vocabulary_status
DROP CONSTRAINT IF EXISTS user_vocabulary_status_status_check;

-- Add new check constraint with updated values
ALTER TABLE user_vocabulary_status
ADD CONSTRAINT user_vocabulary_status_status_check
CHECK (status IN ('learning', 'completed'));

-- Update any existing 'reviewing' or 'mastered' statuses to 'completed'
UPDATE user_vocabulary_status
SET status = 'completed'
WHERE status IN ('reviewing', 'mastered');

COMMIT;
