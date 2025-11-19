package services

import (
	"api-numberniceic/internal/core/domain"
	"api-numberniceic/internal/core/ports"
	"strconv"
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

	satValues := []map[string]int{}
	shaValues := []map[string]int{}

	// ตัวแปรสะสมผลรวม (แยกกัน)
	satSum := 0
	shaSum := 0

	for _, charRune := range cleanName {
		charStr := string(charRune)

		if charStr == " " {
			continue
		}

		// 1. เลขศาสตร์
		satVal, err := s.repo.GetSatValue(charStr)
		if err != nil {
			satVal = 0
		}
		satValues = append(satValues, map[string]int{charStr: satVal})
		satSum += satVal // บวกสะสม

		// 2. เลขพลังเงา
		shaVal, err := s.repo.GetShaValue(charStr)
		if err != nil {
			shaVal = 0
		}
		shaValues = append(shaValues, map[string]int{charStr: shaVal})
		shaSum += shaVal // บวกสะสม
	}

	// สร้างคู่เลขของใครของมัน
	satPairs := s.generatePairs(satSum)
	shaPairs := s.generatePairs(shaSum)

	return &domain.NameAnalysis{
		Name:      cleanName,
		SatValues: satValues,
		ShaValues: shaValues,
		// ส่งค่ากลับแยกกัน
		SatSum:   satSum,
		SatPairs: satPairs,
		ShaSum:   shaSum,
		ShaPairs: shaPairs,
	}, nil
}

// generatePairs สร้างคู่เลขตามกฎที่กำหนด
func (s *analyzerService) generatePairs(sum int) []string {
	strSum := strconv.Itoa(sum)
	length := len(strSum)

	if length == 1 {
		return []string{"0" + strSum} // หลักหน่วย -> 09
	}
	if length == 2 {
		return []string{strSum} // หลักสิบ -> 12
	}
	if length == 3 {
		return []string{
			strSum[0:2], // คู่หน้า (11)
			strSum[1:3], // คู่หลัง (12)
		}
	}
	return []string{}
}
