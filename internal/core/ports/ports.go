package ports

import "api-numberniceic/internal/core/domain"

// NumberRepository (Output Port) - สำหรับติดต่อกับ Database
type NumberRepository interface {
	GetSatValue(char string) (int, error)
	GetShaValue(char string) (int, error)
	GetNumberMeaning(pair string) (*domain.NumberMeaning, error)
	GetKakisByDay(day string) ([]string, error)

	// ฟังก์ชันสำหรับค้นหาชื่อมงคลที่ใกล้เคียง (Levenshtein)
	SearchSimilarNames(name string, day string, limit int) ([]domain.NamesMiracle, error)
}

// NumberService (Input Port) - สำหรับ Business Logic
type NumberService interface {
	// ฟังก์ชันวิเคราะห์ชื่อ (รับชื่อ และ วันเกิด)
	AnalyzeName(name string, birthDay string) (*domain.NameAnalysis, error)
	GetNameLinguistics(name string) (string, error)
}
