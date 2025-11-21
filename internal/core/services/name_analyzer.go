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

// ... (ฟังก์ชันเดิม AnalyzeName ฯลฯ คงเดิม) ...
// ... (Copy ส่วนเดิมมาวาง) ...

func (s *analyzerService) AnalyzeName(name, day string) (*domain.NameAnalysis, error) {
	return &domain.NameAnalysis{}, nil
}                                                                                  // Mock
func (s *analyzerService) GetNameLinguistics(n string) (string, error)             { return "", nil }
func (s *analyzerService) SaveNameForUser(uid uint, admin bool, n, d string) error { return nil }
func (s *analyzerService) GetSavedNames(uid uint) ([]domain.SavedName, error)      { return nil, nil }
func (s *analyzerService) RemoveSavedName(id, uid uint) error                      { return nil }
func (s *analyzerService) GetPairMeaning(p string) (*domain.NumberMeaning, error)  { return nil, nil }
func (s *analyzerService) GetKakisList(d string) ([]string, error)                 { return nil, nil }
func (s *analyzerService) GetEnrichedPairs(sum int) []domain.PairData              { return nil }

// --- Blog Service ---

// 🔥 อัปเดต: รับ shortTitle, typeID
func (s *analyzerService) CreateNewBlog(userID uint, isAdmin bool, title, shortTitle string, typeID uint, content, coverURL string) error {
	if !isAdmin {
		return fmt.Errorf("Unauthorized")
	}
	newBlog := &domain.Blog{
		Title:      title,
		ShortTitle: shortTitle,
		BlogTypeID: typeID,
		Content:    content,
		CoverURL:   coverURL,
		AuthorID:   userID,
	}
	return s.repo.CreateBlog(newBlog)
}

func (s *analyzerService) GetLatestBlogs() ([]domain.Blog, error) {
	return s.repo.GetAllBlogs()
}

func (s *analyzerService) GetBlogDetail(id uint) (*domain.Blog, error) {
	return s.repo.GetBlogByID(id)
}

// 🔥 อัปเดต: รับ shortTitle, typeID
func (s *analyzerService) UpdateExistingBlog(id uint, userID uint, isAdmin bool, title, shortTitle string, typeID uint, content, coverURL string) error {
	if !isAdmin {
		return fmt.Errorf("Unauthorized")
	}
	blog, err := s.repo.GetBlogByID(id)
	if err != nil {
		return err
	}

	blog.Title = title
	blog.ShortTitle = shortTitle
	blog.BlogTypeID = typeID
	blog.Content = content
	blog.CoverURL = coverURL

	return s.repo.UpdateBlog(blog)
}

func (s *analyzerService) RemoveBlog(id uint, userID uint, isAdmin bool) error {
	if !isAdmin {
		return fmt.Errorf("Unauthorized")
	}
	return s.repo.DeleteBlog(id)
}

// 🔥 เพิ่มใหม่: ดึงประเภท
func (s *analyzerService) GetBlogTypes() ([]domain.BlogType, error) {
	return s.repo.GetAllBlogTypes()
}
