package models

import "time"

type Vocabulary struct {
	ID                     uint      `gorm:"primaryKey;autoIncrement"`
	Word                   string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_vocabulary_word_jlpt"`
	Reading                string    `gorm:"type:varchar(255)"`
	Meaning                string    `gorm:"type:varchar(500);not null"`
	PartOfSpeech           string    `gorm:"type:varchar(50)"`
	JlptLevelID            uint      `gorm:"not null;uniqueIndex:idx_vocabulary_word_jlpt;index"`
	CategoryID             uint      `gorm:"not null;index"`
	ExampleSentence        string    `gorm:"type:text"`
	ExampleSentenceReading string    `gorm:"type:text"`
	ExampleSentenceMeaning string    `gorm:"type:text"`
	AudioURL               string    `gorm:"type:varchar(255)"`
	ImageURL               string    `gorm:"type:varchar(255)"`
	Difficulty             int       `gorm:"type:int;default:1;check:difficulty >= 1 AND difficulty <= 5"`
	JlptLevel              JlptLevel `gorm:"foreignKey:JlptLevelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Category               Category  `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedAt              *time.Time
	UpdatedAt              *time.Time
}

// TableName specifies the table name for the Vocabulary model
func (Vocabulary) TableName() string {
	return "vocabularies"
}
