package domain

import "github.com/lib/pq"

type NamesMiracle struct {
	NameID int            `gorm:"primaryKey;column:name_id" json:"name_id"`
	ThName string         `gorm:"column:thname" json:"th_name"`
	SatNum pq.StringArray `gorm:"column:satnum;type:text[]" json:"sat_num"`
	ShaNum pq.StringArray `gorm:"column:shanum;type:text[]" json:"sha_num"`

	// Flag กาลกิณี (true = เป็นกาลกิณี/ห้ามใช้)
	KSunday     bool `gorm:"column:k_sunday"`
	KMonday     bool `gorm:"column:k_monday"`
	KTuesday    bool `gorm:"column:k_tuesday"`
	KWednesday1 bool `gorm:"column:k_wednesday1"`
	KWednesday2 bool `gorm:"column:k_wednesday2"`
	KThursday   bool `gorm:"column:k_thursday"`
	KFriday     bool `gorm:"column:k_friday"`
	KSaturday   bool `gorm:"column:k_saturday"`

	// Field พิเศษสำหรับเก็บค่าความห่างของชื่อ (Distance) ที่คำนวณได้จาก SQL
	Distance float64 `gorm:"column:distance" json:"distance"`
}

func (NamesMiracle) TableName() string {
	return "names_miracle"
}
