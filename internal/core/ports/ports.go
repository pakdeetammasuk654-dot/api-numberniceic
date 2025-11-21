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

	// --- Blog Type Repository ---
	GetAllBlogTypes() ([]domain.BlogType, error)
	SeedBlogTypes() error
	CreateBlogType(blogType *domain.BlogType) error // 🔥 เพิ่มใหม่
	DeleteBlogType(id uint) error                   // 🔥 เพิ่มใหม่
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
	CreateNewBlog(userID uint, isAdmin bool, title, shortTitle string, typeID uint, content, coverURL string) error
	GetLatestBlogs() ([]domain.Blog, error)
	GetBlogDetail(id uint) (*domain.Blog, error)
	UpdateExistingBlog(id uint, userID uint, isAdmin bool, title, shortTitle string, typeID uint, content, coverURL string) error
	RemoveBlog(id uint, userID uint, isAdmin bool) error

	// --- Blog Type Service ---
	GetBlogTypes() ([]domain.BlogType, error)
	CreateNewBlogType(name string) error // 🔥 เพิ่มใหม่
	RemoveBlogType(id uint) error        // 🔥 เพิ่มใหม่
}
