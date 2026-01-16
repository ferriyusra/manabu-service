package models

import "time"

type Exercise struct {
	ID               uint       `gorm:"primaryKey;autoIncrement"`
	LessonID         uint       `gorm:"not null;uniqueIndex:idx_exercise_lesson_order;index"`
	Title            string     `gorm:"type:varchar(200);not null"`
	Description      string     `gorm:"type:varchar(1000)"`
	ExerciseType     string     `gorm:"type:varchar(50);not null;index"`
	OrderIndex       int        `gorm:"type:int;not null;default:0;uniqueIndex:idx_exercise_lesson_order"`
	DifficultyLevel  int        `gorm:"type:int;default:1"`
	EstimatedMinutes int        `gorm:"type:int;default:0"`
	IsPublished      bool       `gorm:"type:boolean;default:false;index"`
	PublishedAt      *time.Time `gorm:"type:timestamp"`
	Lesson           Lesson     `gorm:"foreignKey:LessonID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}

// TableName specifies the table name for the Exercise model
func (Exercise) TableName() string {
	return "exercises"
}
