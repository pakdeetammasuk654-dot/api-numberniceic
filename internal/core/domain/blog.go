package domain

import (
	"time"

	"gorm.io/gorm"
)

// BlogType: ตารางประเภทบทความ (post_types)
type BlogType struct {
	gorm.Model
	Name string `gorm:"column:name;unique;not null" json:"name"`
}

// กำหนดชื่อตารางใน DB เป็น 'post_types'
func (BlogType) TableName() string {
	return "post_types"
}

// Blog: ตารางบทความ (posts)
type Blog struct {
	gorm.Model
	Title      string `gorm:"column:title;not null" json:"title"`
	ShortTitle string `gorm:"column:short_title" json:"short_title"` // 🔥 เพิ่มใหม่
	Content    string `gorm:"column:content;type:text;not null" json:"content"`
	CoverURL   string `gorm:"column:cover_url" json:"cover_url"`

	AuthorID uint `gorm:"column:author_id" json:"author_id"`
	Author   User `gorm:"foreignKey:AuthorID" json:"author"`

	// 🔥 เชื่อมโยงกับ post_types (FK)
	BlogTypeID uint     `gorm:"column:blog_type_id" json:"blog_type_id"`
	BlogType   BlogType `gorm:"foreignKey:BlogTypeID" json:"blog_type"`

	Published bool `gorm:"default:true" json:"published"`
	Views     int  `gorm:"default:0" json:"views"`
	CreatedAt time.Time
}

// กำหนดชื่อตารางใน DB เป็น 'posts'
func (Blog) TableName() string {
	return "posts"
}
