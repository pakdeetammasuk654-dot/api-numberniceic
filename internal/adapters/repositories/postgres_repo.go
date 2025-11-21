package repositories

import (
	"api-numberniceic/internal/core/domain"
	"api-numberniceic/internal/core/ports"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) ports.NumberRepository {
	return &postgresRepository{db: db}
}

// --- Existing Methods (คงเดิม) ---
func (r *postgresRepository) GetSatValue(char string) (int, error) {
	var satNum domain.SatNum
	if err := r.db.Table("sat_nums").Where("char_key = ?", char).First(&satNum).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return satNum.SatValue, nil
}
func (r *postgresRepository) GetShaValue(char string) (int, error) {
	var shaNum domain.ShaNum
	if err := r.db.Table("sha_nums").Where("char_key = ?", char).First(&shaNum).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return shaNum.ShaValue, nil
}
func (r *postgresRepository) GetNumberMeaning(pair string) (*domain.NumberMeaning, error) {
	var meaning domain.NumberMeaning
	if err := r.db.Table("numbers").Where("pairnumber = ?", pair).First(&meaning).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &meaning, nil
}
func (r *postgresRepository) GetKakisByDay(day string) ([]string, error) {
	var kakis []string
	err := r.db.Table("kakis_day").Where("day = ?", day).Pluck("kakis", &kakis).Error
	return kakis, err
}
func (r *postgresRepository) SearchSimilarNames(name string, day string, limit int) ([]domain.NamesMiracle, error) {
	var results []domain.NamesMiracle
	col := "k_sunday" // Mock logic
	query := fmt.Sprintf(`SELECT * FROM (SELECT *, levenshtein($1, thname) / greatest(length($2), length(thname))::real as distance FROM names_miracle WHERE %s = false) as sub ORDER BY distance ASC LIMIT $3`, col)
	r.db.Raw(query, name, name, limit).Scan(&results)
	return results, nil
}
func (r *postgresRepository) SaveName(s *domain.SavedName) error { return r.db.Create(s).Error }
func (r *postgresRepository) GetSavedNamesByUserID(uid uint) ([]domain.SavedName, error) {
	var names []domain.SavedName
	err := r.db.Where("user_id = ?", uid).Order("created_at desc").Find(&names).Error
	return names, err
}
func (r *postgresRepository) DeleteSavedName(id, uid uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, uid).Delete(&domain.SavedName{}).Error
}

// --- Blog Repository ---

func (r *postgresRepository) CreateBlog(blog *domain.Blog) error {
	return r.db.Create(blog).Error
}

func (r *postgresRepository) GetAllBlogs() ([]domain.Blog, error) {
	var blogs []domain.Blog
	// 🔥 Preload BlogType เพิ่มเข้ามา
	err := r.db.Preload("Author").Preload("BlogType").Order("created_at desc").Find(&blogs).Error
	return blogs, err
}

func (r *postgresRepository) GetBlogByID(id uint) (*domain.Blog, error) {
	var blog domain.Blog
	// 🔥 Preload BlogType เพิ่มเข้ามา
	err := r.db.Preload("Author").Preload("BlogType").First(&blog, id).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *postgresRepository) UpdateBlog(blog *domain.Blog) error {
	return r.db.Save(blog).Error
}

func (r *postgresRepository) DeleteBlog(id uint) error {
	return r.db.Delete(&domain.Blog{}, id).Error
}

// 🔥 เพิ่มใหม่: ดึงประเภทบทความทั้งหมด
func (r *postgresRepository) GetAllBlogTypes() ([]domain.BlogType, error) {
	var types []domain.BlogType
	err := r.db.Find(&types).Error
	return types, err
}

// 🔥 เพิ่มใหม่: สร้างประเภทเริ่มต้น (ถ้ายังไม่มี)
func (r *postgresRepository) SeedBlogTypes() error {
	var count int64
	r.db.Model(&domain.BlogType{}).Count(&count)
	if count == 0 {
		types := []domain.BlogType{
			{Name: "ข่าวสาร"},
			{Name: "บทความทั่วไป"},
			{Name: "เคล็ดลับ"},
			{Name: "ดวงชะตา"},
		}
		return r.db.Create(&types).Error
	}
	return nil
}
