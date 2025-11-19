package domain

type NameAnalysis struct {
	Name      string           `json:"name"`
	SatValues []map[string]int `json:"sat_values"` // เปลี่ยนชื่อให้สื่อความหมายว่าเป็น "ค่าเลขศาสตร์"
	ShaValues []map[string]int `json:"sha_values"` // เปลี่ยนชื่อให้สื่อความหมายว่าเป็น "ค่าเลขพลังเงา"
}
