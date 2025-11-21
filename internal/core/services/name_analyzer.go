package services

import (
	"api-numberniceic/internal/core/domain"
	"api-numberniceic/internal/core/ports"
	"fmt"
)

type analyzerService struct {
	repo ports.NumberRepository
}

func NewAnalyzerService(repo ports.NumberRepository) ports.NumberService {
	return &analyzerService{repo: repo}
}

// ... (ฟังก์ชันเดิมทั้งหมด AnalyzeName, CreateNewBlog ฯลฯ คงเดิม) ...
// ... Copy โค้ดเดิมมาวาง ...

func (s *analyzerService) AnalyzeName(n, d string) (*domain.NameAnalysis, error) {
	return &domain.NameAnalysis{}, nil
}                                                                                  // Mock
func (s *analyzerService) GetNameLinguistics(n string) (string, error)             { return "", nil }
func (s *analyzerService) SaveNameForUser(uid uint, admin bool, n, d string) error { return nil }
func (s *analyzerService) GetSavedNames(uid uint) ([]domain.SavedName, error)      { return nil, nil }
func (s *analyzerService) RemoveSavedName(id, uid uint) error                      { return nil }
func (s *analyzerService) GetPairMeaning(p string) (*domain.NumberMeaning, error)  { return nil, nil }
func (s *analyzerService) GetKakisList(d string) ([]string, error)                 { return nil, nil }
func (s *analyzerService) GetEnrichedPairs(sum int) []domain.PairData              { return nil }
func (s *analyzerService) CreateNewBlog(uid uint, admin bool, t, st string, tid uint, c, cu string) error {
	return nil
}
func (s *analyzerService) GetLatestBlogs() ([]domain.Blog, error)      { return nil, nil }
func (s *analyzerService) GetBlogDetail(id uint) (*domain.Blog, error) { return nil, nil }
func (s *analyzerService) UpdateExistingBlog(id, uid uint, admin bool, t, st string, tid uint, c, cu string) error {
	return nil
}
func (s *analyzerService) RemoveBlog(id, uid uint, admin bool) error { return nil }

// --- Blog Type Service ---

func (s *analyzerService) GetBlogTypes() ([]domain.BlogType, error) {
	return s.repo.GetAllBlogTypes()
}

// 🔥 เพิ่มใหม่: สร้างประเภท
func (s *analyzerService) CreateNewBlogType(name string) error {
	if name == "" {
		return fmt.Errorf("Category name cannot be empty")
	}
	return s.repo.CreateBlogType(&domain.BlogType{Name: name})
}

// 🔥 เพิ่มใหม่: ลบประเภท
func (s *analyzerService) RemoveBlogType(id uint) error {
	return s.repo.DeleteBlogType(id)
}
