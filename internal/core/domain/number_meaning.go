package domain

// NumberMeaning ตรงกับตาราง 'numbers' ใน Database
type NumberMeaning struct {
	// เพิ่ม tag gorm:"column:..." เพื่อระบุชื่อคอลัมน์ให้ถูกต้องตาม Database เป๊ะๆ
	PairNumber    string `json:"pair_number" gorm:"column:pairnumber"`
	PairType      string `json:"pair_type" gorm:"column:pairtype"`
	DetailVip     string `json:"detail_vip" gorm:"column:detail_vip"`
	MiracleDetail string `json:"miracle_detail" gorm:"column:miracledetail"`
	MiracleDesc   string `json:"miracle_desc" gorm:"column:miracledesc"`
	PairPoint     int    `json:"pair_point" gorm:"column:pairpoint"`
}

// TableName บอก GORM ว่าให้ไปใช้ตารางชื่อ "numbers"
func (NumberMeaning) TableName() string {
	return "numbers"
}
