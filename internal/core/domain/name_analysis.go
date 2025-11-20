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

	// --- เพิ่มส่วนสรุปคะแนน (Summary) ---
	TotalScore int `json:"total_score"` // คะแนนรวมสุทธิ
	GoodScore  int `json:"good_score"`  // ผลรวมคะแนนดี (D)
	BadScore   int `json:"bad_score"`   // ผลรวมคะแนนร้าย (R)
}
