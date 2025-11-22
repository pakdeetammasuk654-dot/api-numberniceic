package services

import (
	"api-numberniceic/internal/core/domain"
	"api-numberniceic/internal/core/ports"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
	"google.golang.org/genai"
)

type analyzerService struct {
	repo ports.NumberRepository
}

func NewAnalyzerService(repo ports.NumberRepository) ports.NumberService {
	return &analyzerService{repo: repo}
}

// --- 1. Logic วิเคราะห์ชื่อ (กู้คืน) ---

func (s *analyzerService) AnalyzeName(name string, birthDay string) (*domain.NameAnalysis, error) {
	cleanName := strings.TrimSpace(name)
	satValues := []map[string]int{}
	shaValues := []map[string]int{}
	satSum := 0
	shaSum := 0

	kakisList, _ := s.repo.GetKakisByDay(birthDay)
	foundKakis := []string{}
	kakisMap := make(map[string]bool)
	for _, k := range kakisList {
		kakisMap[k] = true
	}

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

	rawSatPairs := s.generatePairs(satSum)
	rawShaPairs := s.generatePairs(shaSum)
	satPairData := s.enrichPairs(rawSatPairs)
	shaPairData := s.enrichPairs(rawShaPairs)

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

	similarNames, _ := s.repo.SearchSimilarNames(cleanName, birthDay, 12)
	for i := range similarNames {
		similarNames[i].SatPairs = s.enrichPairs(similarNames[i].SatNum)
		similarNames[i].ShaPairs = s.enrichPairs(similarNames[i].ShaNum)
	}

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
		SimilarNames: similarNames,
	}, nil
}

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

func (s *analyzerService) GetNameLinguistics(name string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API Key configuration error")
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("GenAI Client Error: %v", err)
	}
	prompt := fmt.Sprintf("อธิบายความหมายและรากศัพท์ของชื่อ '%s' แบบสั้นๆ กระชับ เข้าใจง่าย ในเชิงภาษาศาสตร์และสิริมงคล", name)
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-flash-latest",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("GenAI Generate Error: %v", err)
	}
	return result.Text(), nil
}

func (s *analyzerService) SaveNameForUser(userID uint, isAdmin bool, name, birthDay string) error {
	if !isAdmin {
		existingNames, err := s.repo.GetSavedNamesByUserID(userID)
		if err != nil {
			return err
		}
		if len(existingNames) >= 12 {
			return fmt.Errorf("สมาชิกทั่วไปบันทึกได้สูงสุด 12 รายชื่อ (กรุณาลบชื่อเก่าออกก่อน)")
		}
	}

	analysis, err := s.AnalyzeName(name, birthDay)
	if err != nil {
		return err
	}
	newSave := &domain.SavedName{
		UserID: userID, Name: name, BirthDay: birthDay,
		TotalScore: analysis.TotalScore, SatSum: analysis.SatSum, ShaSum: analysis.ShaSum,
	}
	return s.repo.SaveName(newSave)
}

func (s *analyzerService) GetSavedNames(userID uint) ([]domain.SavedName, error) {
	return s.repo.GetSavedNamesByUserID(userID)
}

func (s *analyzerService) RemoveSavedName(id uint, userID uint) error {
	return s.repo.DeleteSavedName(id, userID)
}

func (s *analyzerService) GetPairMeaning(pair string) (*domain.NumberMeaning, error) {
	return s.repo.GetNumberMeaning(pair)
}

func (s *analyzerService) GetKakisList(day string) ([]string, error) {
	return s.repo.GetKakisByDay(day)
}

func (s *analyzerService) GetEnrichedPairs(sum int) []domain.PairData {
	pairs := s.generatePairs(sum)
	return s.enrichPairs(pairs)
}

// --- 2. Logic Blog Service (ใหม่) ---

func (s *analyzerService) generateUniqueSlug(title string) (string, error) {
	baseSlug := slug.Make(title)
	uniqueSlug := baseSlug
	counter := 1
	for {
		existing, err := s.repo.GetBlogBySlug(uniqueSlug)
		if err != nil {
			return "", err // DB error
		}
		if existing == nil {
			break // Slug is unique
		}
		uniqueSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}
	return uniqueSlug, nil
}

func (s *analyzerService) CreateNewBlog(userID uint, isAdmin bool, title, shortTitle string, typeID uint, content, coverURL string) error {
	if !isAdmin {
		return fmt.Errorf("Unauthorized: Only admin can create blogs")
	}

	generatedSlug, err := s.generateUniqueSlug(title)
	if err != nil {
		return fmt.Errorf("failed to generate slug: %w", err)
	}

	newBlog := &domain.Blog{
		Title:      title,
		ShortTitle: shortTitle,
		Slug:       generatedSlug,
		BlogTypeID: typeID,
		Content:    content,
		CoverURL:   coverURL,
		AuthorID:   userID,
	}
	return s.repo.CreateBlog(newBlog)
}

func (s *analyzerService) GetLatestBlogs() ([]domain.Blog, error) {
	return s.repo.GetAllBlogs()
}

func (s *analyzerService) GetBlogDetail(identifier string) (*domain.Blog, error) {
	// พยายามหาจาก ID ก่อน
	if id, err := strconv.Atoi(identifier); err == nil {
		blog, err := s.repo.GetBlogByID(uint(id))
		if err == nil && blog != nil {
			return blog, nil
		}
	}
	// ถ้าไม่เจอ ID หรือ identifier ไม่ใช่ตัวเลข ให้หาจาก Slug
	return s.repo.GetBlogBySlug(identifier)
}

func (s *analyzerService) UpdateExistingBlog(id uint, userID uint, isAdmin bool, title, shortTitle string, typeID uint, content, coverURL string) error {
	if !isAdmin {
		return fmt.Errorf("Unauthorized")
	}
	blog, err := s.repo.GetBlogByID(id)
	if err != nil {
		return err
	}

	// ถ้ามีการเปลี่ยน Title ให้สร้าง Slug ใหม่
	if blog.Title != title {
		newSlug, err := s.generateUniqueSlug(title)
		if err != nil {
			return fmt.Errorf("failed to generate new slug: %w", err)
		}
		blog.Slug = newSlug
	}

	blog.Title = title
	blog.ShortTitle = shortTitle
	blog.BlogTypeID = typeID
	blog.Content = content
	blog.CoverURL = coverURL

	return s.repo.UpdateBlog(blog)
}

func (s *analyzerService) RemoveBlog(id uint, userID uint, isAdmin bool) error {
	if !isAdmin {
		return fmt.Errorf("Unauthorized")
	}
	return s.repo.DeleteBlog(id)
}

func (s *analyzerService) GetBlogTypes() ([]domain.BlogType, error) {
	return s.repo.GetAllBlogTypes()
}

func (s *analyzerService) CreateNewBlogType(name string) error {
	if name == "" {
		return fmt.Errorf("Category name cannot be empty")
	}
	return s.repo.CreateBlogType(&domain.BlogType{Name: name})
}

func (s *analyzerService) RemoveBlogType(id uint) error {
	return s.repo.DeleteBlogType(id)
}
