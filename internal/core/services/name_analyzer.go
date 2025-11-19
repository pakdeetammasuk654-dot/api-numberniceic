package services

import (
	"api-numberniceic/internal/core/domain"
	"api-numberniceic/internal/core/ports"
	"strings"
)

type analyzerService struct {
	repo ports.NumberRepository
}

func NewAnalyzerService(repo ports.NumberRepository) ports.NumberService {
	return &analyzerService{
		repo: repo,
	}
}

func (s *analyzerService) AnalyzeName(name string) (*domain.NameAnalysis, error) {
	cleanName := strings.TrimSpace(name)

	// เตรียมตัวแปรสำหรับเก็บผลลัพธ์ (เปลี่ยนชื่อตัวแปรให้สอดคล้องกัน)
	satValues := []map[string]int{}
	shaValues := []map[string]int{}

	// 1. Logic: แยกสระและอักษร
	for _, charRune := range cleanName {
		charStr := string(charRune)

		if charStr == " " {
			continue
		}

		// 2.1 หาค่าเลขศาสตร์ (Sat)
		satVal, err := s.repo.GetSatValue(charStr)
		if err != nil {
			satVal = 0
		}
		// เก็บผลลัพธ์คู่ ["ก": 1]
		satValues = append(satValues, map[string]int{charStr: satVal})

		// 2.2 หาค่าเลขพลังเงา (Sha)
		shaVal, err := s.repo.GetShaValue(charStr)
		if err != nil {
			shaVal = 0
		}
		// เก็บผลลัพธ์คู่ ["ก": 15]
		shaValues = append(shaValues, map[string]int{charStr: shaVal})
	}

	// 3. ส่งคืนผลลัพธ์ (ใช้ชื่อฟิลด์ใหม่)
	return &domain.NameAnalysis{
		Name:      cleanName,
		SatValues: satValues,
		ShaValues: shaValues,
	}, nil
}
