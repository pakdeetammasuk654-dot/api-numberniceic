package domain

// PairData โครงสร้างใหม่สำหรับเก็บคู่เลขพร้อมความหมาย
type PairData struct {
	Pair    string         `json:"pair"`
	Meaning *NumberMeaning `json:"meaning"` // เก็บความหมายที่ดึงจาก DB
}

type NameAnalysis struct {
	Name      string           `json:"name"`
	SatValues []map[string]int `json:"sat_values"`
	ShaValues []map[string]int `json:"sha_values"`

	// --- เลขศาสตร์ (Sat) ---
	SatSum   int        `json:"sat_sum"`
	SatPairs []PairData `json:"sat_pairs"` // เปลี่ยนเป็น []PairData

	// --- เลขพลังเงา (Sha) ---
	ShaSum   int        `json:"sha_sum"`
	ShaPairs []PairData `json:"sha_pairs"` // เปลี่ยนเป็น []PairData
}
