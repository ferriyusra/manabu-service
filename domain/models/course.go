package models

import "time"

type Course struct {
	ID             uint       `gorm:"primaryKey;autoIncrement"`
	Title          string     `gorm:"type:varchar(200);not null;uniqueIndex:idx_course_title_jlpt"`
	Description    string     `gorm:"type:text;not null"`
	JlptLevelID    uint       `gorm:"not null;uniqueIndex:idx_course_title_jlpt;index"`
	ThumbnailURL   string     `gorm:"type:varchar(255)"`
	Difficulty     int        `gorm:"type:int;default:1;check:difficulty >= 1 AND difficulty <= 5"`
	EstimatedHours int        `gorm:"type:int"`
	IsPublished    bool       `gorm:"type:boolean;default:false"`
	PublishedAt    *time.Time `gorm:"type:timestamp"`
	JlptLevel      JlptLevel  `gorm:"foreignKey:JlptLevelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

// TableName specifies the table name for the Course model
func (Course) TableName() string {
	return "courses"
}
