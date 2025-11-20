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

func (s *analyzerService) AnalyzeName(name string, birthDay string) (*domain.NameAnalysis, error) {
	cleanName := strings.TrimSpace(name)
	satValues := []map[string]int{}
	shaValues := []map[string]int{}
	satSum := 0
	shaSum := 0

	// 1. เตรียมข้อมูลกาลกิณี
	kakisList, _ := s.repo.GetKakisByDay(birthDay)
	foundKakis := []string{}
	kakisMap := make(map[string]bool)
	for _, k := range kakisList {
		kakisMap[k] = true
	}

	// 2. คำนวณค่าพลัง
	for _, charRune := range cleanName {
		charStr := string(charRune)
		if charStr == " " {
			continue
		}
		if kakisMap[charStr] {
			foundKakis = append(foundKakis, charStr)
		}
		satVal, _ := s.repo.GetSatValue(charStr)
		satValues = append(satValues, map[string]int{charStr: satVal})
		satSum += satVal

		shaVal, _ := s.repo.GetShaValue(charStr)
		shaValues = append(shaValues, map[string]int{charStr: shaVal})
		shaSum += shaVal
	}

	// 3. สร้างคู่เลข & ดึงความหมาย
	rawSatPairs := s.generatePairs(satSum)
	rawShaPairs := s.generatePairs(shaSum)
	satPairData := s.enrichPairs(rawSatPairs)
	shaPairData := s.enrichPairs(rawShaPairs)

	// 4. คำนวณคะแนนรวม
	totalScore := 0
	goodScore := 0
	badScore := 0
	calculatePoints := func(pairs []domain.PairData) {
		for _, p := range pairs {
			if p.Meaning != nil {
				score := p.Meaning.PairPoint
				totalScore += score
				pType := strings.ToUpper(p.Meaning.PairType)
				if strings.HasPrefix(pType, "D") {
					goodScore += score
				} else if strings.HasPrefix(pType, "R") {
					badScore += score
				}
			}
		}
	}
	calculatePoints(satPairData)
	calculatePoints(shaPairData)

	// --- 5. (NEW) ค้นหาชื่อที่คล้ายกันจาก Database ---
	// ดึงมา 12 ชื่อ ตามโจทย์ โดยใช้ชื่อที่ User พิมพ์เข้ามาเป็นตัวตั้งต้นค้นหา
	similarNames, _ := s.repo.SearchSimilarNames(cleanName, birthDay, 12)

	return &domain.NameAnalysis{
		Name:         cleanName,
		BirthDay:     birthDay,
		KakisFound:   foundKakis,
		HasKakis:     len(foundKakis) > 0,
		SatValues:    satValues,
		ShaValues:    shaValues,
		SatSum:       satSum,
		SatPairs:     satPairData,
		ShaSum:       shaSum,
		ShaPairs:     shaPairData,
		TotalScore:   totalScore,
		GoodScore:    goodScore,
		BadScore:     badScore,
		SimilarNames: similarNames, // ส่งรายชื่อกลับไปที่ View
	}, nil
}

// ... (func enrichPairs, generatePairs คงเดิม) ...
func (s *analyzerService) enrichPairs(pairs []string) []domain.PairData {
	var result []domain.PairData
	for _, p := range pairs {
		meaning, _ := s.repo.GetNumberMeaning(p)
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
