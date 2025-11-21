package handlers

import (
	"api-numberniceic/internal/core/domain"
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

func getUserInfoFromContext(c *fiber.Ctx) (uint, bool) {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return 0, false
	}
	token, _ := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return getJwtSecret(), nil
	})
	if token != nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := uint(0)
			isAdmin := false
			if idFloat, ok := claims["user_id"].(float64); ok {
				userID = uint(idFloat)
			}
			if adminVal, ok := claims["is_admin"].(bool); ok {
				isAdmin = adminVal
			}
			return userID, isAdmin
		}
	}
	return 0, false
}

func (h *FiberHandler) RenderWithAuth(c *fiber.Ctx, template string, data fiber.Map) error {
	cookie := c.Cookies("jwt")
	isLoggedIn := false
	displayName := ""
	isAdmin := false

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
				if adminVal, ok := claims["is_admin"].(bool); ok {
					isAdmin = adminVal
				}
			}
		}
	}

	if data == nil {
		data = fiber.Map{}
	}
	data["IsLoggedIn"] = isLoggedIn
	data["DisplayName"] = displayName
	data["IsAdmin"] = isAdmin

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

func (h *FiberHandler) ViewAbout(c *fiber.Ctx) error {
	return h.RenderWithAuth(c, "about", nil)
}

func (h *FiberHandler) ViewDashboard(c *fiber.Ctx) error {
	userID, _ := getUserInfoFromContext(c)
	savedNames, _ := h.service.GetSavedNames(userID)

	totalCount := len(savedNames)
	totalScoreSum := 0

	type SavedNameView struct {
		ID         uint
		Name       string
		NameHTML   template.HTML
		BirthDay   string
		BirthDayTH string
		TotalScore int
		SatSum     int
		SatPairs   []domain.PairData
		ShaSum     int
		ShaPairs   []domain.PairData
	}

	var viewModels []SavedNameView

	for _, n := range savedNames {
		totalScoreSum += n.TotalScore

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
			SatPairs:   satPairs,
			ShaSum:     n.ShaSum,
			ShaPairs:   shaPairs,
		})
	}

	return h.RenderWithAuth(c, "dashboard", fiber.Map{
		"SavedNames":    viewModels,
		"TotalCount":    totalCount,
		"TotalScoreSum": totalScoreSum,
	})
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

// --- Blog View Handlers ---

func (h *FiberHandler) ViewArticles(c *fiber.Ctx) error {
	blogs, _ := h.service.GetLatestBlogs()
	type BlogView struct {
		domain.Blog
		SummaryHTML template.HTML
	}
	var blogViews []BlogView
	for _, b := range blogs {
		summary := b.Content
		if len(summary) > 200 {
			summary = summary[:200] + "..."
		}
		blogViews = append(blogViews, BlogView{
			Blog:        b,
			SummaryHTML: template.HTML(summary),
		})
	}
	_, isAdmin := getUserInfoFromContext(c)
	return h.RenderWithAuth(c, "articles", fiber.Map{
		"Blogs":   blogViews,
		"IsAdmin": isAdmin,
	})
}

func (h *FiberHandler) ViewBlogDetail(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	blog, err := h.service.GetBlogDetail(uint(id))
	if err != nil {
		return c.Redirect("/articles")
	}
	return h.RenderWithAuth(c, "blog_detail", fiber.Map{
		"Blog":        blog,
		"ContentHTML": template.HTML(blog.Content),
	})
}

// --- Admin Handlers ---

// 🔥 เพิ่มใหม่: หน้า Admin Panel รวมเมนู
func (h *FiberHandler) ViewAdminPanel(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Redirect("/")
	}
	return h.RenderWithAuth(c, "admin/panel", nil)
}

func (h *FiberHandler) ViewCreateBlog(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Redirect("/articles")
	}
	return h.RenderWithAuth(c, "admin/create_blog", nil)
}

func (h *FiberHandler) ViewAdminBlogs(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Redirect("/")
	}
	blogs, _ := h.service.GetLatestBlogs()
	return h.RenderWithAuth(c, "admin/blogs", fiber.Map{"Blogs": blogs})
}

func (h *FiberHandler) ViewEditBlog(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Redirect("/")
	}
	id, _ := c.ParamsInt("id")
	blog, err := h.service.GetBlogDetail(uint(id))
	if err != nil {
		return c.Redirect("/admin/blogs")
	}
	return h.RenderWithAuth(c, "admin/edit_blog", fiber.Map{"Blog": blog})
}

// --- Blog Action Handlers ---

func (h *FiberHandler) HandleCreateBlog(c *fiber.Ctx) error {
	userID, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Status(403).SendString("Unauthorized")
	}
	title := c.FormValue("title")
	content := c.FormValue("content")
	coverURL := c.FormValue("cover_url")

	if err := h.service.CreateNewBlog(userID, isAdmin, title, content, coverURL); err != nil {
		return h.RenderWithAuth(c, "admin/create_blog", fiber.Map{"Error": err.Error()})
	}
	return c.Redirect("/articles")
}

func (h *FiberHandler) HandleEditBlog(c *fiber.Ctx) error {
	userID, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Status(403).SendString("Unauthorized")
	}
	id, _ := c.ParamsInt("id")
	title := c.FormValue("title")
	content := c.FormValue("content")
	coverURL := c.FormValue("cover_url")

	if err := h.service.UpdateExistingBlog(uint(id), userID, isAdmin, title, content, coverURL); err != nil {
		return h.RenderWithAuth(c, "admin/edit_blog", fiber.Map{
			"Error": err.Error(),
			"Blog":  fiber.Map{"ID": id, "Title": title, "Content": content, "CoverURL": coverURL},
		})
	}
	return c.Redirect("/admin/blogs")
}

func (h *FiberHandler) HandleDeleteBlog(c *fiber.Ctx) error {
	userID, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Status(403).SendString("Unauthorized")
	}
	id, _ := c.ParamsInt("id")
	h.service.RemoveBlog(uint(id), userID, isAdmin)
	return c.Redirect("/admin/blogs")
}

// --- API Handlers ---

func (h *FiberHandler) ApiSaveName(c *fiber.Ctx) error {
	userID, isAdmin := getUserInfoFromContext(c)
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
	if err := h.service.SaveNameForUser(userID, isAdmin, req.Name, req.BirthDay); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Saved successfully"})
}

func (h *FiberHandler) ApiDeleteName(c *fiber.Ctx) error {
	userID, _ := getUserInfoFromContext(c)
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
