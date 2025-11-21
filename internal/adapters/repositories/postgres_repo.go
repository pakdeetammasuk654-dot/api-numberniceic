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

// 🔥 เพิ่มส่วน SavedName Repository 🔥

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
