package model

import "time"

type Article struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Slug        string    `gorm:"uniqueIndex;size:255;not null" json:"slug"`
	ContentMD   string    `gorm:"type:longtext;not null" json:"content_md"`
	ContentHTML string    `gorm:"type:longtext;not null" json:"content_html"`
	Excerpt     string    `gorm:"type:text" json:"excerpt"`
	Cover       string    `gorm:"size:512;default:''" json:"cover"`
	Status      int       `gorm:"default:0;index" json:"status"`
	ViewCount   int       `gorm:"default:0" json:"view_count"`
	IsTop        int       `gorm:"default:0" json:"is_top"`
	IsAnnouncement int     `gorm:"default:0" json:"is_announcement"`
	CategoryID  *uint     `gorm:"index" json:"category_id"`
	Category    *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags        []Tag     `gorm:"many2many:article_tags" json:"tags,omitempty"`
	Series      string    `gorm:"size:255;index;default:''" json:"series"`
	SeriesOrder int       `gorm:"default:0" json:"series_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Article) TableName() string {
	return "articles"
}
