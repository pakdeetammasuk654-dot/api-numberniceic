package repositories

import (
	"api-numberniceic/internal/core/domain"
	"api-numberniceic/internal/core/ports"

	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) ports.NumberRepository {
	return &postgresRepository{db: db}
}

// ... (ฟังก์ชันเดิม GetSatValue, SaveName, CreateBlog ฯลฯ คงเดิม) ...
// ... Copy โค้ดเดิมทั้งหมดมาวาง ...

// (ส่วนเดิม...)
func (r *postgresRepository) GetSatValue(char string) (int, error) { /*...*/ return 0, nil } // Mock for brevity, use original
func (r *postgresRepository) GetShaValue(char string) (int, error) { /*...*/ return 0, nil }
func (r *postgresRepository) GetNumberMeaning(pair string) (*domain.NumberMeaning, error) { /*...*/
	return nil, nil
}
func (r *postgresRepository) GetKakisByDay(day string) ([]string, error) { /*...*/ return nil, nil }
func (r *postgresRepository) SearchSimilarNames(name string, day string, limit int) ([]domain.NamesMiracle, error) { /*...*/
	return nil, nil
}
func (r *postgresRepository) SaveName(savedName *domain.SavedName) error {
	return r.db.Create(savedName).Error
}
func (r *postgresRepository) GetSavedNamesByUserID(userID uint) ([]domain.SavedName, error) { /*...*/
	return nil, nil
}
func (r *postgresRepository) DeleteSavedName(id uint, userID uint) error {
	return r.db.Delete(&domain.SavedName{}, id).Error
}
func (r *postgresRepository) CreateBlog(blog *domain.Blog) error        { return r.db.Create(blog).Error }
func (r *postgresRepository) GetAllBlogs() ([]domain.Blog, error)       { /*...*/ return nil, nil }
func (r *postgresRepository) GetBlogByID(id uint) (*domain.Blog, error) { /*...*/ return nil, nil }
func (r *postgresRepository) UpdateBlog(blog *domain.Blog) error        { return r.db.Save(blog).Error }
func (r *postgresRepository) DeleteBlog(id uint) error                  { return r.db.Delete(&domain.Blog{}, id).Error }

// --- Blog Type Repository (Updated) ---

func (r *postgresRepository) GetAllBlogTypes() ([]domain.BlogType, error) {
	var types []domain.BlogType
	err := r.db.Find(&types).Error
	return types, err
}

func (r *postgresRepository) SeedBlogTypes() error {
	var count int64
	r.db.Model(&domain.BlogType{}).Count(&count)
	if count == 0 {
		types := []domain.BlogType{
			{Name: "ข่าวสาร"}, {Name: "บทความทั่วไป"}, {Name: "เคล็ดลับ"}, {Name: "ดวงชะตา"},
		}
		return r.db.Create(&types).Error
	}
	return nil
}

// 🔥 เพิ่มใหม่: สร้างประเภท
func (r *postgresRepository) CreateBlogType(blogType *domain.BlogType) error {
	return r.db.Create(blogType).Error
}

// 🔥 เพิ่มใหม่: ลบประเภท
func (r *postgresRepository) DeleteBlogType(id uint) error {
	return r.db.Delete(&domain.BlogType{}, id).Error
}
