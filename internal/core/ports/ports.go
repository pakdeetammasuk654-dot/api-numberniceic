package ports

import "api-numberniceic/internal/core/domain"

// NumberRepository (Output Port)
type NumberRepository interface {
	GetSatValue(char string) (int, error)
	GetShaValue(char string) (int, error)
	GetNumberMeaning(pair string) (*domain.NumberMeaning, error)
	GetKakisByDay(day string) ([]string, error)
	SearchSimilarNames(name string, day string, limit int) ([]domain.NamesMiracle, error)
	SaveName(savedName *domain.SavedName) error
	GetSavedNamesByUserID(userID uint) ([]domain.SavedName, error)
	DeleteSavedName(id uint, userID uint) error
}

// NumberService (Input Port)
type NumberService interface {
	AnalyzeName(name string, birthDay string) (*domain.NameAnalysis, error)
	GetNameLinguistics(name string) (string, error)
	SaveNameForUser(userID uint, name, birthDay string) error
	GetSavedNames(userID uint) ([]domain.SavedName, error)
	RemoveSavedName(id uint, userID uint) error
	GetPairMeaning(pair string) (*domain.NumberMeaning, error)
	GetKakisList(day string) ([]string, error)

	// 🔥 เพิ่มใหม่: ฟังก์ชันแปลงผลรวมเป็นคู่เลขพร้อมความหมาย (เช่น 190 -> 19(D), 90(R))
	GetEnrichedPairs(sum int) []domain.PairData
}
