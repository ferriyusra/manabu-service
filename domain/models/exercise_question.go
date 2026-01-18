package models

import "time"

type ExerciseQuestion struct {
	ID            uint       `gorm:"primaryKey;autoIncrement"`
	ExerciseID    uint       `gorm:"not null;uniqueIndex:idx_question_exercise_order;index"`
	QuestionText  string     `gorm:"type:text;not null"`
	QuestionType  string     `gorm:"type:varchar(50);not null;index"`
	Options       string     `gorm:"type:text"`
	CorrectAnswer string     `gorm:"type:text;not null"`
	Explanation   string     `gorm:"type:text"`
	AudioURL      string     `gorm:"type:varchar(500)"`
	ImageURL      string     `gorm:"type:varchar(500)"`
	OrderIndex    int        `gorm:"type:int;not null;default:0;uniqueIndex:idx_question_exercise_order"`
	Points        int        `gorm:"type:int;not null;default:10"`
	IsPublished   bool       `gorm:"type:boolean;default:false;index"`
	PublishedAt   *time.Time `gorm:"type:timestamp"`
	Exercise      Exercise   `gorm:"foreignKey:ExerciseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

// TableName specifies the table name for the ExerciseQuestion model
func (ExerciseQuestion) TableName() string {
	return "exercise_questions"
}
