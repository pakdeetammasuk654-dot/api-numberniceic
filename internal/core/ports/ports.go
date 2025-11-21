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

	// --- Blog Repository ---
	CreateBlog(blog *domain.Blog) error
	GetAllBlogs() ([]domain.Blog, error)
	GetBlogByID(id uint) (*domain.Blog, error)
	UpdateBlog(blog *domain.Blog) error
	DeleteBlog(id uint) error

	// 🔥 เพิ่มใหม่: ดึงประเภทบทความทั้งหมด
	GetAllBlogTypes() ([]domain.BlogType, error)
	// (Optional) สร้างประเภทเริ่มต้น
	SeedBlogTypes() error
}

// NumberService (Input Port)
type NumberService interface {
	AnalyzeName(name string, birthDay string) (*domain.NameAnalysis, error)
	GetNameLinguistics(name string) (string, error)
	SaveNameForUser(userID uint, isAdmin bool, name, birthDay string) error
	GetSavedNames(userID uint) ([]domain.SavedName, error)
	RemoveSavedName(id uint, userID uint) error
	GetPairMeaning(pair string) (*domain.NumberMeaning, error)
	GetKakisList(day string) ([]string, error)
	GetEnrichedPairs(sum int) []domain.PairData

	// --- Blog Service ---
	// 🔥 อัปเดต: รับ shortTitle และ typeID เพิ่ม
	CreateNewBlog(userID uint, isAdmin bool, title, shortTitle string, typeID uint, content, coverURL string) error
	GetLatestBlogs() ([]domain.Blog, error)
	GetBlogDetail(id uint) (*domain.Blog, error)
	UpdateExistingBlog(id uint, userID uint, isAdmin bool, title, shortTitle string, typeID uint, content, coverURL string) error
	RemoveBlog(id uint, userID uint, isAdmin bool) error

	// 🔥 เพิ่มใหม่: ขอ list ประเภท
	GetBlogTypes() ([]domain.BlogType, error)
}
