package ports

import "api-numberniceic/internal/core/domain"

// NumberRepository (Output Port)
// สั่งให้ Database ไปหาค่าตัวเลข จากตัวอักษร (char) ที่ส่งไป
type NumberRepository interface {
	GetSatValue(char string) (int, error) // เช่น ส่ง "ก" ได้ 1
	GetShaValue(char string) (int, error) // เช่น ส่ง "ก" ได้ 15
}

// NumberService (Input Port)
// สั่งให้ระบบทำการวิเคราะห์ชื่อ
type NumberService interface {
	AnalyzeName(name string) (*domain.NameAnalysis, error)
}
