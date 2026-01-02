package models

import "time"

type Category struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_category_name_jlpt"`
	Description string    `gorm:"type:varchar(255)"`
	JlptLevelID uint      `gorm:"not null;uniqueIndex:idx_category_name_jlpt"`
	JlptLevel   JlptLevel `gorm:"foreignKey:JlptLevelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

// TableName specifies the table name for the Category model
func (Category) TableName() string {
	return "categories"
}
