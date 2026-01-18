package models

import (
	"time"

	"github.com/google/uuid"
)

// Status constants for user course progress
const (
	ProgressStatusNotStarted = "not_started"
	ProgressStatusInProgress = "in_progress"
	ProgressStatusCompleted  = "completed"
)

type UserCourseProgress struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID             uint       `gorm:"not null;uniqueIndex:idx_user_course_progress_user_course;index"`
	CourseID           uint       `gorm:"not null;uniqueIndex:idx_user_course_progress_user_course;index"`
	Status             string     `gorm:"type:varchar(50);not null;default:'not_started';index;check:status IN ('not_started', 'in_progress', 'completed')"`
	ProgressPercentage float64    `gorm:"type:decimal(5,2);not null;default:0.00;check:progress_percentage >= 0 AND progress_percentage <= 100"`
	CompletedLessons   int        `gorm:"type:int;not null;default:0"`
	TotalLessons       int        `gorm:"type:int;not null;default:0"`
	StartedAt          *time.Time `gorm:"type:timestamp"`
	CompletedAt        *time.Time `gorm:"type:timestamp"`
	LastAccessedAt     *time.Time `gorm:"type:timestamp;index"`
	User               User       `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Course             Course     `gorm:"foreignKey:CourseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}

// TableName specifies the table name for the UserCourseProgress model
func (UserCourseProgress) TableName() string {
	return "user_course_progress"
}
