package domain

type KakisDay struct {
	KakisID int    `gorm:"primaryKey;column:kakis_id"`
	Day     string `gorm:"column:day"`
	Kakis   string `gorm:"column:kakis"` // ตัวอักษรกาลกิณี 1 ตัว
}

// TableName ระบุชื่อตารางให้ตรงกับ database
func (KakisDay) TableName() string {
	return "kakis_day"
}
