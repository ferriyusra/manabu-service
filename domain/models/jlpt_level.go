package models

import "time"

type JlptLevel struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Code        string `gorm:"type:varchar(10);not null;uniqueIndex"`
	Name        string `gorm:"type:varchar(50);not null"`
	Description string `gorm:"type:varchar(255);not null"`
	LevelOrder  int    `gorm:"type:int;not null"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
