package ports

import "api-numberniceic/internal/core/domain"

// NumberRepository (Output Port)
type NumberRepository interface {
	GetSatValue(char string) (int, error)
	GetShaValue(char string) (int, error)
	GetNumberMeaning(pair string) (*domain.NumberMeaning, error)
	GetKakisByDay(day string) ([]string, error)

	// เพิ่ม: ค้นหาชื่อมงคลที่ใกล้เคียง
	SearchSimilarNames(name string, day string, limit int) ([]domain.NamesMiracle, error)
}

// NumberService (Input Port)
type NumberService interface {
	// อัปเดต: รับวันเกิดเพิ่มเข้ามา
	AnalyzeName(name string, birthDay string) (*domain.NameAnalysis, error)
}
