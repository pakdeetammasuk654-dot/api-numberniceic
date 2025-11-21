package domain

import (
	"time"

	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	Title     string `gorm:"column:title;not null" json:"title"`
	Content   string `gorm:"column:content;type:text;not null" json:"content"` // เก็บ HTML จาก TinyMCE
	CoverURL  string `gorm:"column:cover_url" json:"cover_url"`
	AuthorID  uint   `gorm:"column:author_id" json:"author_id"`
	Author    User   `gorm:"foreignKey:AuthorID" json:"author"` // เชื่อมกับตาราง Users
	Published bool   `gorm:"default:true" json:"published"`
	Views     int    `gorm:"default:0" json:"views"`
	CreatedAt time.Time
}

func (Blog) TableName() string {
	return "blogs"
}
