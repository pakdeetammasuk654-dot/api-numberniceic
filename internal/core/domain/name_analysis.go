package domain

type PairData struct {
	Pair    string         `json:"pair"`
	Meaning *NumberMeaning `json:"meaning"`
}

type NameAnalysis struct {
	Name      string           `json:"name"`
	SatValues []map[string]int `json:"sat_values"`
	ShaValues []map[string]int `json:"sha_values"`

	SatSum   int        `json:"sat_sum"`
	SatPairs []PairData `json:"sat_pairs"`

	ShaSum   int        `json:"sha_sum"`
	ShaPairs []PairData `json:"sha_pairs"`

	// --- ส่วนสรุปคะแนน (Summary) ---
	TotalScore int `json:"total_score"` // คะแนนรวมสุทธิ
	GoodScore  int `json:"good_score"`  // ผลรวมคะแนนดี (D)
	BadScore   int `json:"bad_score"`   // ผลรวมคะแนนร้าย (R)

	// --- ส่วนกาลกิณี (New) ---
	BirthDay   string   `json:"birth_day"`   // วันเกิดที่เลือก
	KakisFound []string `json:"kakis_found"` // รายการอักษรกาลกิณีที่เจอในชื่อ
	HasKakis   bool     `json:"has_kakis"`   // เจอหรือไม่ (true/false)
}
