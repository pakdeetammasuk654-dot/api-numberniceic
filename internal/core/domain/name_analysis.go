package domain

type NameAnalysis struct {
	Name      string           `json:"name"`
	SatValues []map[string]int `json:"sat_values"` // เปลี่ยนชื่อให้สื่อความหมายว่าเป็น "ค่าเลขศาสตร์"
	ShaValues []map[string]int `json:"sha_values"` // เปลี่ยนชื่อให้สื่อความหมายว่าเป็น "ค่าเลขพลังเงา"

	// --- เลขศาสตร์ (Sat) ---
	SatSum   int      `json:"sat_sum"`   // ผลรวม
	SatPairs []string `json:"sat_pairs"` // คู่เลข

	// --- เลขพลังเงา (Sha) ---
	ShaSum   int      `json:"sha_sum"`   // ผลรวม
	ShaPairs []string `json:"sha_pairs"` // คู่เลข
}
