package domain

import "gorm.io/gorm"

type SavedName struct {
	gorm.Model
	UserID     uint   `gorm:"column:user_id;index" json:"user_id"` // ผูกกับ User คนไหน
	Name       string `gorm:"column:name" json:"name"`
	BirthDay   string `gorm:"column:birth_day" json:"birth_day"`
	TotalScore int    `gorm:"column:total_score" json:"total_score"` // เก็บผลคะแนนไว้โชว์เลย ไม่ต้องคำนวณใหม่
	SatSum     int    `gorm:"column:sat_sum" json:"sat_sum"`
	ShaSum     int    `gorm:"column:sha_sum" json:"sha_sum"`
}

func (SavedName) TableName() string {
	return "saved_names"
}
