package handlers

import (
	"api-numberniceic/internal/core/domain" // 🔥 เพิ่ม import domain
	"api-numberniceic/internal/core/ports"
	"html/template"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type FiberHandler struct {
	service ports.NumberService
}

func NewFiberHandler(service ports.NumberService) *FiberHandler {
	return &FiberHandler{
		service: service,
	}
}

func getJwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "supersecretkey"
	}
	return []byte(secret)
}

func getUserIDFromContext(c *fiber.Ctx) uint {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return 0
	}
	token, _ := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return getJwtSecret(), nil
	})
	if token != nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if idFloat, ok := claims["user_id"].(float64); ok {
				return uint(idFloat)
			}
		}
	}
	return 0
}

func (h *FiberHandler) RenderWithAuth(c *fiber.Ctx, template string, data fiber.Map) error {
	cookie := c.Cookies("jwt")
	isLoggedIn := false
	displayName := ""

	if cookie != "" {
		token, _ := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			return getJwtSecret(), nil
		})
		if token != nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				isLoggedIn = true
				if name, ok := claims["display_name"].(string); ok {
					displayName = name
				}
			}
		}
	}

	if data == nil {
		data = fiber.Map{}
	}
	data["IsLoggedIn"] = isLoggedIn
	data["DisplayName"] = displayName

	return c.Render(template, data, "layouts/main")
}

func translateDay(day string) string {
	switch day {
	case "sunday":
		return "อาทิตย์"
	case "monday":
		return "จันทร์"
	case "tuesday":
		return "อังคาร"
	case "wednesday1":
		return "พุธ (กลางวัน)"
	case "wednesday2":
		return "พุธ (กลางคืน)"
	case "thursday":
		return "พฤหัสบดี"
	case "friday":
		return "ศุกร์"
	case "saturday":
		return "เสาร์"
	default:
		return day
	}
}

// --- View Handlers ---

func (h *FiberHandler) ViewHome(c *fiber.Ctx) error {
	return h.RenderWithAuth(c, "home", nil)
}

func (h *FiberHandler) ViewDashboard(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	savedNames, _ := h.service.GetSavedNames(userID)

	totalCount := len(savedNames)
	totalScoreSum := 0

	// ViewModel
	type SavedNameView struct {
		ID         uint
		Name       string
		NameHTML   template.HTML
		BirthDay   string
		BirthDayTH string
		TotalScore int
		SatSum     int
		// 🔥 เปลี่ยนจาก PairType ตัวเดียว เป็น Array ของ PairData
		SatPairs []domain.PairData
		ShaSum   int
		ShaPairs []domain.PairData
	}

	var viewModels []SavedNameView

	for _, n := range savedNames {
		totalScoreSum += n.TotalScore

		// 1. กาลกิณี
		kakis, _ := h.service.GetKakisList(n.BirthDay)
		kakisMap := make(map[string]bool)
		for _, k := range kakis {
			kakisMap[k] = true
		}

		var sb strings.Builder
		for _, r := range n.Name {
			s := string(r)
			if kakisMap[s] {
				sb.WriteString(`<span class="bad-char">` + s + `</span>`)
			} else {
				sb.WriteString(s)
			}
		}
		nameHTML := template.HTML(sb.String())

		// 2. 🔥 ใช้ GetEnrichedPairs เพื่อแตกผลรวมเป็นคู่เลข (เช่น 190 -> 19, 90)
		satPairs := h.service.GetEnrichedPairs(n.SatSum)
		shaPairs := h.service.GetEnrichedPairs(n.ShaSum)

		viewModels = append(viewModels, SavedNameView{
			ID:         n.ID,
			Name:       n.Name,
			NameHTML:   nameHTML,
			BirthDay:   n.BirthDay,
			BirthDayTH: translateDay(n.BirthDay),
			TotalScore: n.TotalScore,
			SatSum:     n.SatSum,
			SatPairs:   satPairs, // ส่ง Array ไป Loop
			ShaSum:     n.ShaSum,
			ShaPairs:   shaPairs, // ส่ง Array ไป Loop
		})
	}

	return h.RenderWithAuth(c, "dashboard", fiber.Map{
		"SavedNames":    viewModels,
		"TotalCount":    totalCount,
		"TotalScoreSum": totalScoreSum,
	})
}

func (h *FiberHandler) ViewArticles(c *fiber.Ctx) error {
	return h.RenderWithAuth(c, "articles", nil)
}

func (h *FiberHandler) ViewAbout(c *fiber.Ctx) error {
	return h.RenderWithAuth(c, "about", nil)
}

func (h *FiberHandler) ViewAnalysis(c *fiber.Ctx) error {
	name := c.Query("name")
	birthDay := c.Query("birth_day")
	if name == "" {
		name = "ณเดชน์"
	}
	if birthDay == "" {
		birthDay = "sunday"
	}
	result, err := h.service.AnalyzeName(name, birthDay)
	data := fiber.Map{"Name": name, "BirthDay": birthDay}
	if err == nil {
		data["Result"] = result
	} else {
		data["Error"] = "ไม่สามารถโหลดข้อมูลได้: " + err.Error()
	}
	return h.RenderWithAuth(c, "analysis", data)
}

func (h *FiberHandler) HandleAnalysis(c *fiber.Ctx) error {
	name := c.FormValue("name")
	birthDay := c.FormValue("birth_day")
	result, err := h.service.AnalyzeName(name, birthDay)
	data := fiber.Map{"Name": name, "BirthDay": birthDay}
	if err != nil {
		data["Error"] = err.Error()
	} else {
		data["Result"] = result
	}
	return h.RenderWithAuth(c, "analysis", data)
}

func (h *FiberHandler) ApiSaveName(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized", "redirect": "/login"})
	}
	type SaveRequest struct {
		Name     string `json:"name"`
		BirthDay string `json:"birth_day"`
	}
	var req SaveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid data"})
	}
	if err := h.service.SaveNameForUser(userID, req.Name, req.BirthDay); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Saved successfully"})
}

func (h *FiberHandler) ApiDeleteName(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	id, _ := c.ParamsInt("id")
	h.service.RemoveSavedName(uint(id), userID)
	return c.Redirect("/dashboard")
}

func (h *FiberHandler) ApiAnalyze(c *fiber.Ctx) error {
	name := c.Query("name")
	birthDay := c.Query("birth_day")
	if name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name is required"})
	}
	result, err := h.service.AnalyzeName(name, birthDay)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *FiberHandler) ApiGetLinguistics(c *fiber.Ctx) error {
	name := c.Query("name")
	if name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name is required"})
	}
	meaning, err := h.service.GetNameLinguistics(name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"text": meaning})
}
