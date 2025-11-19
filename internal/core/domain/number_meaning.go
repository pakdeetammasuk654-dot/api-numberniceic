package domain

// NumberMeaning ตรงกับตาราง 'numbers' ใน Database
type NumberMeaning struct {
	PairNumber    string `json:"pair_number"`    // pairnumber (char 2)
	PairType      string `json:"pair_type"`      // pairtype
	DetailVip     string `json:"detail_vip"`     // detail_vip
	MiracleDetail string `json:"miracle_detail"` // miracledetail
	MiracleDesc   string `json:"miracle_desc"`   // miracledesc
	PairPoint     int    `json:"pair_point"`     // pairpoint
}

// TableName บอก GORM ว่าให้ไปใช้ตารางชื่อ "numbers"
func (NumberMeaning) TableName() string {
	return "numbers"
}
