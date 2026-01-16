package models

import "time"

type Lesson struct {
	ID               uint       `gorm:"primaryKey;autoIncrement"`
	CourseID         uint       `gorm:"not null;uniqueIndex:idx_lesson_course_order;index"`
	Title            string     `gorm:"type:varchar(255);not null"`
	Content          string     `gorm:"type:text"`
	OrderIndex       int        `gorm:"type:int;not null;default:0;uniqueIndex:idx_lesson_course_order"`
	EstimatedMinutes int        `gorm:"type:int;default:0"`
	IsPublished      bool       `gorm:"type:boolean;default:false;index"`
	PublishedAt      *time.Time `gorm:"type:timestamp"`
	Course           Course     `gorm:"foreignKey:CourseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}

// TableName specifies the table name for the Lesson model
func (Lesson) TableName() string {
	return "lessons"
}
