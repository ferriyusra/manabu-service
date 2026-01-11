package models

import (
	"time"

	"github.com/google/uuid"
)

// UserVocabularyStatus represents the learning progress and spaced repetition data for a user learning a vocabulary word
type UserVocabularyStatus struct {
	ID             uint       `gorm:"primaryKey;autoIncrement"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_vocabulary"`
	VocabularyID   uint       `gorm:"not null;index;uniqueIndex:idx_user_vocabulary"`
	Status         string     `gorm:"type:varchar(20);not null;default:'learning';check:status IN ('learning', 'completed')"`
	Repetitions    int        `gorm:"type:int;not null;default:0"`
	LastReviewedAt *time.Time `gorm:"null"`
	Vocabulary     Vocabulary `gorm:"foreignKey:VocabularyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

// TableName specifies the table name for the UserVocabularyStatus model
func (UserVocabularyStatus) TableName() string {
	return "user_vocabulary_status"
}
