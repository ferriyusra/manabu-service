-- Migration: Create lessons table
-- Description: Creates the lessons table with proper constraints and indexes
-- Author: AI Assistant
-- Date: 2026-01-15

-- Create lessons table
CREATE TABLE IF NOT EXISTS lessons (
    id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    order_index INTEGER NOT NULL DEFAULT 0,
    estimated_minutes INTEGER DEFAULT 0,
    is_published BOOLEAN DEFAULT false,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key constraint with CASCADE on delete
    CONSTRAINT fk_lessons_course
        FOREIGN KEY (course_id)
        REFERENCES courses(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_lessons_course_id ON lessons(course_id);
CREATE INDEX IF NOT EXISTS idx_lessons_is_published ON lessons(is_published);

-- Create unique constraint on course_id and order_index combination
-- This ensures that each lesson within a course has a unique order_index
CREATE UNIQUE INDEX IF NOT EXISTS idx_lessons_course_order ON lessons(course_id, order_index);

-- Add comment to table
COMMENT ON TABLE lessons IS 'Stores lesson information for courses in the Japanese learning application';

-- Add comments to columns
COMMENT ON COLUMN lessons.id IS 'Primary key, auto-incrementing lesson ID';
COMMENT ON COLUMN lessons.course_id IS 'Foreign key reference to courses table';
COMMENT ON COLUMN lessons.title IS 'Title of the lesson';
COMMENT ON COLUMN lessons.content IS 'Full content/body of the lesson';
COMMENT ON COLUMN lessons.order_index IS 'Order position of lesson within the course (must be unique per course)';
COMMENT ON COLUMN lessons.estimated_minutes IS 'Estimated time to complete the lesson in minutes';
COMMENT ON COLUMN lessons.is_published IS 'Publication status of the lesson';
COMMENT ON COLUMN lessons.published_at IS 'Timestamp when the lesson was published';
COMMENT ON COLUMN lessons.created_at IS 'Timestamp when the lesson was created';
COMMENT ON COLUMN lessons.updated_at IS 'Timestamp when the lesson was last updated';
