package domain

import (
	"time"

	"gorm.io/gorm"
)

// User Struct: ตรงกับตาราง public.users
type User struct {
	ID           uint           `gorm:"primaryKey;column:id" json:"id"`
	Username     string         `gorm:"unique;not null;column:username" json:"username"`
	Email        string         `gorm:"unique;not null;column:email" json:"email"`
	PasswordHash string         `gorm:"not null;column:password_hash" json:"-"` // ใช้ password_hash ตาม DB
	DisplayName  string         `gorm:"not null;column:display_name" json:"display_name"`
	IsAdmin      bool           `gorm:"default:false;column:is_admin" json:"is_admin"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

func (User) TableName() string {
	return "users"
}
