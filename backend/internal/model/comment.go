package model

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `gorm:"index;not null" json:"article_id"`
	ParentID  *uint     `gorm:"index" json:"parent_id"`
	Author    string    `gorm:"size:64;not null" json:"author"`
	Email     string    `gorm:"size:128;default:''" json:"email,omitempty"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Status    int       `gorm:"default:1;index" json:"status"`
	IP        string    `gorm:"size:64;default:''" json:"ip,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Article   *Article  `gorm:"foreignKey:ArticleID" json:"-"`
	Parent    *Comment  `gorm:"foreignKey:ParentID" json:"-"`
	Replies   []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}
