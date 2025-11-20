package ports

import "api-numberniceic/internal/core/domain"

// NumberRepository (Output Port)
type NumberRepository interface {
	GetSatValue(char string) (int, error)
	GetShaValue(char string) (int, error)
	GetNumberMeaning(pair string) (*domain.NumberMeaning, error)

	// เพิ่ม: ดึงรายการอักษรกาลกิณีตามวันเกิด
	GetKakisByDay(day string) ([]string, error)
}

// NumberService (Input Port)
type NumberService interface {
	// อัปเดต: รับวันเกิดเพิ่มเข้ามา
	AnalyzeName(name string, birthDay string) (*domain.NameAnalysis, error)
}
