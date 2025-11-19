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
	return &analyzerService{repo: repo}
}

func (s *analyzerService) AnalyzeName(name string) (*domain.NameAnalysis, error) {
	cleanName := strings.TrimSpace(name)
	satValues := []map[string]int{}
	shaValues := []map[string]int{}
	satSum := 0
	shaSum := 0

	// 1. คำนวณตัวเลข
	for _, charRune := range cleanName {
		charStr := string(charRune)
		if charStr == " " {
			continue
		}

		satVal, _ := s.repo.GetSatValue(charStr)
		satValues = append(satValues, map[string]int{charStr: satVal})
		satSum += satVal

		shaVal, _ := s.repo.GetShaValue(charStr)
		shaValues = append(shaValues, map[string]int{charStr: shaVal})
		shaSum += shaVal
	}

	// 2. สร้างคู่เลข (Strings)
	rawSatPairs := s.generatePairs(satSum)
	rawShaPairs := s.generatePairs(shaSum)

	// 3. ดึงความหมายของคู่เลข (Enrich Data)
	satPairData := s.enrichPairs(rawSatPairs)
	shaPairData := s.enrichPairs(rawShaPairs)

	return &domain.NameAnalysis{
		Name:      cleanName,
		SatValues: satValues,
		ShaValues: shaValues,
		SatSum:    satSum,
		SatPairs:  satPairData, // ส่งกลับแบบมีข้อมูลครบ
		ShaSum:    shaSum,
		ShaPairs:  shaPairData, // ส่งกลับแบบมีข้อมูลครบ
	}, nil
}

// ฟังก์ชันช่วย: วนลูปคู่เลขและดึงความหมายจาก DB
func (s *analyzerService) enrichPairs(pairs []string) []domain.PairData {
	var result []domain.PairData
	for _, p := range pairs {
		meaning, _ := s.repo.GetNumberMeaning(p) // ไปดึงจาก DB
		result = append(result, domain.PairData{
			Pair:    p,
			Meaning: meaning,
		})
	}
	return result
}

func (s *analyzerService) generatePairs(sum int) []string {
	strSum := strconv.Itoa(sum)
	length := len(strSum)
	if length == 1 {
		return []string{"0" + strSum}
	}
	if length == 2 {
		return []string{strSum}
	}
	if length == 3 {
		return []string{strSum[0:2], strSum[1:3]}
	}
	return []string{}
}
