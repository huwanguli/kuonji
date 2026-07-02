package model

import "time"

type Series struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:255;not null" json:"name"`
	Cover       string    `gorm:"size:512;default:''" json:"cover"`
	Description string    `gorm:"size:512;default:''" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
