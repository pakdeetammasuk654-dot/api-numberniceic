package domain

import "github.com/lib/pq"

type NamesMiracle struct {
	NameID int            `gorm:"primaryKey;column:name_id" json:"name_id"`
	ThName string         `gorm:"column:thname" json:"th_name"`
	SatNum pq.StringArray `gorm:"column:satnum;type:text[]" json:"sat_num"`
	ShaNum pq.StringArray `gorm:"column:shanum;type:text[]" json:"sha_num"`

	// Flag กาลกิณี
	KSunday     bool `gorm:"column:k_sunday"`
	KMonday     bool `gorm:"column:k_monday"`
	KTuesday    bool `gorm:"column:k_tuesday"`
	KWednesday1 bool `gorm:"column:k_wednesday1"`
	KWednesday2 bool `gorm:"column:k_wednesday2"`
	KThursday   bool `gorm:"column:k_thursday"`
	KFriday     bool `gorm:"column:k_friday"`
	KSaturday   bool `gorm:"column:k_saturday"`

	Distance float64 `gorm:"column:distance" json:"distance"`

	// --- NEW: ฟิลด์สำหรับเก็บข้อมูลคู่เลขพร้อมความหมาย (ใช้แสดงผลหน้าเว็บ) ---
	SatPairs []PairData `gorm:"-" json:"sat_pairs"`
	ShaPairs []PairData `gorm:"-" json:"sha_pairs"`
}

func (NamesMiracle) TableName() string {
	return "names_miracle"
}
