package model

import "time"

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	Slug        string    `gorm:"uniqueIndex;size:128;not null" json:"slug"`
	Description string    `gorm:"size:255;default:''" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Articles    []Article `gorm:"foreignKey:CategoryID" json:"-"`
}
