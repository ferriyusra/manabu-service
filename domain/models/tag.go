package models

import "time"

type Tag struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(50);not null;uniqueIndex:idx_tag_name"`
	Description string `gorm:"type:varchar(255)"`
	Color       string `gorm:"type:varchar(7)"` // Hex color code, e.g., "#FF5733"
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

// TableName specifies the table name for the Tag model
func (Tag) TableName() string {
	return "tags"
}
