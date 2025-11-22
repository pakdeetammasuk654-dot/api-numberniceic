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

// --- 1. ระบบวิเคราะห์ชื่อ (กู้คืน) ---

func (r *postgresRepository) GetSatValue(char string) (int, error) {
	var satNum domain.SatNum
	result := r.db.Table("sat_nums").Where("char_key = ?", char).First(&satNum)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, result.Error
	}
	return satNum.SatValue, nil
}

func (r *postgresRepository) GetShaValue(char string) (int, error) {
	var shaNum domain.ShaNum
	result := r.db.Table("sha_nums").Where("char_key = ?", char).First(&shaNum)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, result.Error
	}
	return shaNum.ShaValue, nil
}

func (r *postgresRepository) GetNumberMeaning(pair string) (*domain.NumberMeaning, error) {
	var meaning domain.NumberMeaning
	result := r.db.Table("numbers").Where("pairnumber = ?", pair).First(&meaning)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &meaning, nil
}

func (r *postgresRepository) GetKakisByDay(day string) ([]string, error) {
	var kakisList []string
	result := r.db.Table("kakis_day").Where("day = ?", day).Pluck("kakis", &kakisList)
	if result.Error != nil {
		return nil, result.Error
	}
	return kakisList, nil
}

func (r *postgresRepository) SearchSimilarNames(name string, day string, limit int) ([]domain.NamesMiracle, error) {
	var results []domain.NamesMiracle

	columnMap := map[string]string{
		"sunday":     "k_sunday",
		"monday":     "k_monday",
		"tuesday":    "k_tuesday",
		"wednesday1": "k_wednesday1",
		"wednesday2": "k_wednesday2",
		"thursday":   "k_thursday",
		"friday":     "k_friday",
		"saturday":   "k_saturday",
	}

	targetCol, ok := columnMap[day]
	if !ok {
		targetCol = "k_sunday"
	}

	query := fmt.Sprintf(`
		SELECT 
			name_id, thname, satnum, shanum, 
			k_sunday, k_monday, k_tuesday, k_wednesday1, k_wednesday2, k_thursday, k_friday, k_saturday,
			levenshtein(?, thname) / greatest(length(?), length(thname))::real as distance
		FROM names_miracle
		WHERE %s = false
		ORDER BY distance ASC
		LIMIT ?
	`, targetCol)

	err := r.db.Raw(query, name, name, limit).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// --- 2. ระบบบันทึกชื่อ (กู้คืน) ---

func (r *postgresRepository) SaveName(savedName *domain.SavedName) error {
	return r.db.Create(savedName).Error
}

func (r *postgresRepository) GetSavedNamesByUserID(userID uint) ([]domain.SavedName, error) {
	var names []domain.SavedName
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&names).Error
	return names, err
}

func (r *postgresRepository) DeleteSavedName(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&domain.SavedName{}).Error
}

// --- 3. ระบบ Blog (ใหม่) ---

func (r *postgresRepository) CreateBlog(blog *domain.Blog) error {
	return r.db.Create(blog).Error
}

func (r *postgresRepository) GetAllBlogs(limit int) ([]domain.Blog, error) {
	var blogs []domain.Blog
	query := r.db.Preload("Author").Preload("BlogType").Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&blogs).Error
	return blogs, err
}

func (r *postgresRepository) GetBlogByID(id uint) (*domain.Blog, error) {
	var blog domain.Blog
	err := r.db.Preload("Author").Preload("BlogType").First(&blog, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &blog, nil
}

func (r *postgresRepository) GetBlogBySlug(slug string) (*domain.Blog, error) {
	var blog domain.Blog
	err := r.db.Preload("Author").Preload("BlogType").Where("slug = ?", slug).First(&blog).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
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

// --- 4. ระบบ Blog Type (ใหม่) ---

func (r *postgresRepository) GetAllBlogTypes() ([]domain.BlogType, error) {
	var types []domain.BlogType
	err := r.db.Order("id asc").Find(&types).Error
	return types, err
}

func (r *postgresRepository) GetBlogTypeByID(id uint) (*domain.BlogType, error) {
	var blogType domain.BlogType
	err := r.db.First(&blogType, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("ไม่พบประเภท ID: %d", id)
		}
		return nil, err
	}
	return &blogType, nil
}

func (r *postgresRepository) GetBlogTypeByName(name string) (*domain.BlogType, error) {
	var blogType domain.BlogType
	err := r.db.Where("name = ?", name).First(&blogType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // คืนค่า nil เมื่อไม่พบ ถือเป็นเรื่องปกติ
		}
		return nil, err // คืนค่า error กรณีเกิดปัญหาอื่นๆ
	}
	return &blogType, nil
}


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

func (r *postgresRepository) CreateBlogType(blogType *domain.BlogType) error {
	return r.db.Create(blogType).Error
}

func (r *postgresRepository) UpdateBlogType(blogType *domain.BlogType) error {
	return r.db.Save(blogType).Error
}

func (r *postgresRepository) DeleteBlogType(id uint) error {
	return r.db.Delete(&domain.BlogType{}, id).Error
}
