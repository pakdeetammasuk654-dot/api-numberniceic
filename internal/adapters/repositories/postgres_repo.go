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

// ... (ฟังก์ชันเดิม GetSatValue, GetShaValue, GetNumberMeaning, GetKakisByDay คงเดิม) ...
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

// --- NEW: Implement SearchSimilarNames ---
func (r *postgresRepository) SearchSimilarNames(name string, day string, limit int) ([]domain.NamesMiracle, error) {
	var results []domain.NamesMiracle

	// Map วันเกิด -> ชื่อคอลัมน์ใน DB
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

	// เลือกคอลัมน์กาลกิณีให้ตรงกับวัน (Default เป็นวันอาทิตย์ถ้าหาไม่เจอ)
	targetCol, ok := columnMap[day]
	if !ok {
		targetCol = "k_sunday"
	}

	// SQL Query (ดัดแปลงจากโค้ดเดิมของคุณเพื่อให้เข้ากับ Gorm)
	// ใช้ levenshtein(...) / greatest(...) เพื่อหา % ความเหมือน
	// และกรองเฉพาะชื่อที่ไม่เป็นกาลกิณี (column = false)
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

	// Execute Raw Query
	// params: name (สำหรับ $1 ใน levenshtein), name (สำหรับ $1 ใน greatest), limit
	err := r.db.Raw(query, name, name, limit).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
