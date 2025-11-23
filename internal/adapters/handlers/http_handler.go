package handlers

import (
	"api-numberniceic/internal/core/domain"
	"api-numberniceic/internal/core/ports"
	"html/template"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

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
    data["Path"] = c.Path() // Add current path for SEO

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

func truncate(s string, length int) string {
	if utf8.RuneCountInString(s) <= length {
		return s
	}
	runes := []rune(s)
	return string(runes[:length]) + "..."
}


// --- View Handlers (General) ---

func (h *FiberHandler) ViewHome(c *fiber.Ctx) error {
	blogs, _ := h.service.GetLatestBlogs(7)

	var mainArticle *domain.Blog
	var secondaryArticles []domain.Blog
	var linkListArticles []domain.Blog

	if len(blogs) > 0 {
		mainArticle = &blogs[0]
	}
	if len(blogs) > 1 {
		end := min(4, len(blogs))
		secondaryArticles = blogs[1:end]
	}
	if len(blogs) > 4 {
		linkListArticles = blogs[4:]
	}

	return h.RenderWithAuth(c, "landing_page", fiber.Map{
		"Title":             "หน้าแรก",
		"Description":       "วิเคราะห์ชื่อมงคลตามหลักเลขศาสตร์และพลังเงา ค้นพบความหมายที่ซ่อนอยู่เบื้องหลังชื่อของคุณ",
		"MainArticle":       mainArticle,
		"SecondaryArticles": secondaryArticles,
		"LinkListArticles":  linkListArticles,
	})
}

func (h *FiberHandler) ViewAbout(c *fiber.Ctx) error {
	return h.RenderWithAuth(c, "about", fiber.Map{
        "Title": "เกี่ยวกับเรา",
        "Description": "ชื่อดี.com ให้บริการวิเคราะห์ชื่อตามหลักเลขศาสตร์และพลังเงา เพื่อส่งเสริมสิริมงคลและความสำเร็จในชีวิต",
    })
}

func (h *FiberHandler) ViewSitemap(c *fiber.Ctx) error {
	blogs, _ := h.service.GetLatestBlogs(0) // Get all blogs
	return h.RenderWithAuth(c, "sitemap", fiber.Map{
		"Title": "แผนที่เว็บไซต์",
        "Description": "รวมลิงก์ทั้งหมดของเว็บไซต์ ชื่อดี.com เพื่อให้ง่ายต่อการเข้าถึง",
		"Blogs": blogs,
	})
}

func (h *FiberHandler) HandleSitemapXML(c *fiber.Ctx) error {
	blogs, _ := h.service.GetLatestBlogs(0) // Get all blogs
	c.Set("Content-Type", "application/xml")
	return c.Render("sitemap_xml", fiber.Map{
		"Blogs": blogs,
	})
}

// 🔥 กู้คืน: ViewAnalysis
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

	data := fiber.Map{
		"Title":       "วิเคราะห์ชื่อ",
        "Description": "วิเคราะห์ชื่อ " + name + " ตามหลักเลขศาสตร์และพลังเงา ค้นพบความหมายและพลังที่ซ่อนอยู่",
		"Name":     name,
		"BirthDay": birthDay,
	}
	if err == nil {
		data["Result"] = result
	} else {
		data["Error"] = "ไม่สามารถโหลดข้อมูลได้: " + err.Error()
	}

	return h.RenderWithAuth(c, "analysis", data)
}

// 🔥 กู้คืน: HandleAnalysis
func (h *FiberHandler) HandleAnalysis(c *fiber.Ctx) error {
	name := c.FormValue("name")
	birthDay := c.FormValue("birth_day")

	result, err := h.service.AnalyzeName(name, birthDay)

	data := fiber.Map{
		"Title":       "ผลการวิเคราะห์ชื่อ " + name,
        "Description": "ผลการวิเคราะห์ชื่อ " + name + " ตามหลักเลขศาสตร์และพลังเงา พร้อมคำทำนายและชื่อที่แนะนำ",
		"Name":     name,
		"BirthDay": birthDay,
	}

	if err != nil {
		data["Error"] = err.Error()
	} else {
		data["Result"] = result
	}

	return h.RenderWithAuth(c, "analysis", data)
}

// 🔥 กู้คืน: ViewDashboard
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
        "Title": "แดชบอร์ด",
        "Description": "ดูและจัดการชื่อที่คุณบันทึกไว้ทั้งหมด",
		"SavedNames":    viewModels,
		"TotalCount":    totalCount,
		"TotalScoreSum": totalScoreSum,
	})
}

// --- Blog Handlers ---

func (h *FiberHandler) ViewArticles(c *fiber.Ctx) error {
	blogs, _ := h.service.GetLatestBlogs(0) // 0 for all
	_, isAdmin := getUserInfoFromContext(c)
	return h.RenderWithAuth(c, "articles", fiber.Map{
        "Title": "บทความทั้งหมด",
        "Description": "รวมบทความน่ารู้เกี่ยวกับศาสตร์ตัวเลข พลังเงา และเคล็ดลับการตั้งชื่อ",
		"Blogs":   blogs,
		"IsAdmin": isAdmin,
	})
}

func (h *FiberHandler) ViewBlogDetail(c *fiber.Ctx) error {
	slug := c.Params("slug")
	blog, err := h.service.GetBlogDetail(slug)
	if err != nil || blog == nil {
		return c.Redirect("/articles")
	}
	return h.RenderWithAuth(c, "blog_detail", fiber.Map{
        "Title":       blog.Title,
        "Description": blog.Description,
        "Image":       blog.CoverURL,
		"Blog":        blog,
		"ContentHTML": template.HTML(blog.Content),
	})
}

// --- Admin Handlers ---

func (h *FiberHandler) ViewAdminPanel(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Redirect("/")
	}
	return h.RenderWithAuth(c, "admin/panel", fiber.Map{"Title": "แผงควบคุม"})
}

func (h *FiberHandler) ViewCreateBlog(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Redirect("/articles")
	}
	types, _ := h.service.GetBlogTypes()
	return h.RenderWithAuth(c, "admin/create_blog", fiber.Map{
        "Title": "เขียนบทความใหม่",
		"Types": types,
	})
}

func (h *FiberHandler) ViewAdminBlogs(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Redirect("/")
	}
	blogs, _ := h.service.GetLatestBlogs(0) // 0 for all

	// Create a view model to hold truncated values
	type BlogView struct {
		domain.Blog
		TruncatedTitle       string
		TruncatedDescription string
	}
	var blogViews []BlogView
	for _, b := range blogs {
		blogViews = append(blogViews, BlogView{
			Blog:                 b,
			TruncatedTitle:       truncate(b.Title, 30),
			TruncatedDescription: truncate(b.Description, 40),
		})
	}

	return h.RenderWithAuth(c, "admin/blogs", fiber.Map{
		"Title": "จัดการบทความ",
		"Blogs": blogViews,
	})
}

func (h *FiberHandler) ViewEditBlog(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Redirect("/")
	}
	id := c.Params("id")
	blog, err := h.service.GetBlogDetail(id)
	if err != nil {
		return c.Redirect("/admin/blogs")
	}
	types, _ := h.service.GetBlogTypes()
	return h.RenderWithAuth(c, "admin/edit_blog", fiber.Map{
        "Title": "แก้ไขบทความ",
		"Blog":  blog,
		"Types": types,
	})
}

func (h *FiberHandler) ViewAdminTypes(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Redirect("/")
	}
	types, _ := h.service.GetBlogTypes()
	return h.RenderWithAuth(c, "admin/types", fiber.Map{
        "Title": "จัดการประเภทบทความ",
		"Types": types,
	})
}

func (h *FiberHandler) ViewEditBlogType(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Redirect("/")
	}
	id, _ := c.ParamsInt("id")
	blogType, err := h.service.GetBlogTypeByID(uint(id))
	if err != nil {
		return c.Redirect("/admin/types")
	}
	return h.RenderWithAuth(c, "admin/edit_type", fiber.Map{
        "Title": "แก้ไขประเภทบทความ",
		"Type": blogType,
	})
}

// --- Admin Actions ---

func (h *FiberHandler) HandleCreateBlog(c *fiber.Ctx) error {
	userID, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Status(403).SendString("Unauthorized")
	}

	title := c.FormValue("title")
	shortTitle := c.FormValue("short_title")
	description := c.FormValue("description")
	content := c.FormValue("content")
	coverURL := c.FormValue("cover_url")
	typeID, _ := strconv.Atoi(c.FormValue("blog_type_id"))

	if err := h.service.CreateNewBlog(userID, isAdmin, title, shortTitle, description, uint(typeID), content, coverURL); err != nil {
		types, _ := h.service.GetBlogTypes()
		return h.RenderWithAuth(c, "admin/create_blog", fiber.Map{"Error": err.Error(), "Types": types})
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
	shortTitle := c.FormValue("short_title")
	description := c.FormValue("description")
	content := c.FormValue("content")
	coverURL := c.FormValue("cover_url")
	typeIDVal, _ := strconv.Atoi(c.FormValue("blog_type_id"))
	typeID := uint(typeIDVal)

	if err := h.service.UpdateExistingBlog(uint(id), userID, isAdmin, title, shortTitle, description, typeID, content, coverURL); err != nil {
		// หากเกิดข้อผิดพลาด ให้กลับไปที่หน้าแก้ไขพร้อมกับข้อความ Error
		blog, _ := h.service.GetBlogDetail(strconv.Itoa(id))
		types, _ := h.service.GetBlogTypes()
		return h.RenderWithAuth(c, "admin/edit_blog", fiber.Map{
			"Blog":  blog,
			"Types": types,
			"Error": err.Error(),
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

func (h *FiberHandler) HandleCreateBlogType(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Status(403).SendString("Unauthorized")
	}
	name := c.FormValue("name")
	if err := h.service.CreateNewBlogType(name); err != nil {
		// Consider showing an error message on the page
		return c.Redirect("/admin/types")
	}
	return c.Redirect("/admin/types")
}

func (h *FiberHandler) HandleEditBlogType(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Status(403).SendString("Unauthorized")
	}
	id, _ := c.ParamsInt("id")
	name := c.FormValue("name")

	if err := h.service.UpdateBlogType(uint(id), name); err != nil {
		// หากเกิดข้อผิดพลาด ให้กลับไปที่หน้าแก้ไขพร้อมกับข้อความ Error
		blogType, _ := h.service.GetBlogTypeByID(uint(id))
		return h.RenderWithAuth(c, "admin/edit_type", fiber.Map{
			"Type":  blogType,
			"Error": err.Error(),
		})
	}

	// หากสำเร็จ ให้กลับไปที่หน้ารายการ
	return c.Redirect("/admin/types")
}


func (h *FiberHandler) HandleDeleteBlogType(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Status(403).SendString("Unauthorized")
	}
	id, _ := c.ParamsInt("id")
	h.service.RemoveBlogType(uint(id))
	return c.Redirect("/admin/types")
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
