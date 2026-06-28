package model

import "time"

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex;size:128;not null" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}
