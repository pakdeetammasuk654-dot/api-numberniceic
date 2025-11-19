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

// NewPostgresRepository สร้าง Repo ใหม่ที่เชื่อมต่อกับ DB
func NewPostgresRepository(db *gorm.DB) ports.NumberRepository {
	return &postgresRepository{
		db: db,
	}
}

// GetSatValue ดึงค่าเลขศาสตร์จากตาราง sat_nums
func (r *postgresRepository) GetSatValue(char string) (int, error) {
	var satNum domain.SatNum

	// คำสั่ง SQL: SELECT * FROM sat_nums WHERE char_key = 'char' LIMIT 1
	// GORM จะแปลงชื่อ Field 'CharKey' เป็น 'char_key' ให้อัตโนมัติ
	result := r.db.Table("sat_nums").Where("char_key = ?", char).First(&satNum)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil // ถ้าไม่เจอตัวอักษร ให้ค่าเป็น 0 (ไม่ถือว่า Error ร้ายแรง)
		}
		return 0, result.Error
	}

	return satNum.SatValue, nil
}

// GetShaValue ดึงค่าเลขพลังเงาจากตาราง sha_nums
func (r *postgresRepository) GetShaValue(char string) (int, error) {
	var shaNum domain.ShaNum

	// คำสั่ง SQL: SELECT * FROM sha_nums WHERE char_key = 'char' LIMIT 1
	result := r.db.Table("sha_nums").Where("char_key = ?", char).First(&shaNum)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil // ถ้าไม่เจอตัวอักษร ให้ค่าเป็น 0
		}
		return 0, result.Error
	}

	return shaNum.ShaValue, nil
}
