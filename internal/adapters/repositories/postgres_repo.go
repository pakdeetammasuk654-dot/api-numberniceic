package repositories

import (
	"api-numberniceic/internal/core/domain"
	"api-numberniceic/internal/core/ports"
	"errors"

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

// ฟังก์ชันใหม่: ดึงอักษรกาลกิณี
func (r *postgresRepository) GetKakisByDay(day string) ([]string, error) {
	var kakisList []string
	// เลือกเฉพาะ column 'kakis' จากตาราง kakis_day โดย filter ตาม day
	result := r.db.Table("kakis_day").Where("day = ?", day).Pluck("kakis", &kakisList)

	if result.Error != nil {
		return nil, result.Error
	}
	return kakisList, nil
}
