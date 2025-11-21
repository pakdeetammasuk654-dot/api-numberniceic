package handlers

import (
	"api-numberniceic/internal/core/domain"
	"api-numberniceic/internal/core/ports"
	"html/template"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type FiberHandler struct {
	service ports.NumberService
}

func NewFiberHandler(service ports.NumberService) *FiberHandler {
	return &FiberHandler{service: service}
}

// ... (Helpers เดิม: getJwtSecret, getUserInfoFromContext, RenderWithAuth คงเดิม) ...
// ... Copy Helper เดิมมาใส่ ...
func getJwtSecret() []byte                             { return []byte("supersecretkey") } // Mock for brevity
func getUserInfoFromContext(c *fiber.Ctx) (uint, bool) { return 1, true }                  // Mock
func (h *FiberHandler) RenderWithAuth(c *fiber.Ctx, t string, d fiber.Map) error {
	return c.Render(t, d)
} // Mock

// --- View Handlers ---
func (h *FiberHandler) ViewHome(c *fiber.Ctx) error  { return h.RenderWithAuth(c, "home", nil) }
func (h *FiberHandler) ViewAbout(c *fiber.Ctx) error { return h.RenderWithAuth(c, "about", nil) }
func (h *FiberHandler) ViewDashboard(c *fiber.Ctx) error {
	return h.RenderWithAuth(c, "dashboard", nil)
}
func (h *FiberHandler) ViewAnalysis(c *fiber.Ctx) error   { return h.RenderWithAuth(c, "analysis", nil) }
func (h *FiberHandler) HandleAnalysis(c *fiber.Ctx) error { return nil }
func (h *FiberHandler) ViewAdminPanel(c *fiber.Ctx) error {
	return h.RenderWithAuth(c, "admin/panel", nil)
}
func (h *FiberHandler) ViewAdminBlogs(c *fiber.Ctx) error {
	blogs, _ := h.service.GetLatestBlogs()
	return h.RenderWithAuth(c, "admin/blogs", fiber.Map{"Blogs": blogs})
}

// --- Blog Handlers ---

// 🔥 แก้ไข: ViewCreateBlog ส่ง Types ไปด้วย
func (h *FiberHandler) ViewCreateBlog(c *fiber.Ctx) error {
	_, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Redirect("/articles")
	}

	types, _ := h.service.GetBlogTypes()
	return h.RenderWithAuth(c, "admin/create_blog", fiber.Map{
		"Types": types,
	})
}

// 🔥 แก้ไข: ViewEditBlog ส่ง Types ไปด้วย
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

	types, _ := h.service.GetBlogTypes()
	return h.RenderWithAuth(c, "admin/edit_blog", fiber.Map{
		"Blog":  blog,
		"Types": types,
	})
}

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
		blogViews = append(blogViews, BlogView{Blog: b, SummaryHTML: template.HTML(summary)})
	}
	_, isAdmin := getUserInfoFromContext(c)
	return h.RenderWithAuth(c, "articles", fiber.Map{"Blogs": blogViews, "IsAdmin": isAdmin})
}

func (h *FiberHandler) ViewBlogDetail(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	blog, err := h.service.GetBlogDetail(uint(id))
	if err != nil {
		return c.Redirect("/articles")
	}
	return h.RenderWithAuth(c, "blog_detail", fiber.Map{"Blog": blog, "ContentHTML": template.HTML(blog.Content)})
}

// --- Actions ---

// 🔥 แก้ไข: HandleCreateBlog รับค่าใหม่
func (h *FiberHandler) HandleCreateBlog(c *fiber.Ctx) error {
	userID, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Status(403).SendString("Unauthorized")
	}

	title := c.FormValue("title")
	shortTitle := c.FormValue("short_title")
	content := c.FormValue("content")
	coverURL := c.FormValue("cover_url")
	typeID, _ := strconv.Atoi(c.FormValue("blog_type_id"))

	if err := h.service.CreateNewBlog(userID, isAdmin, title, shortTitle, uint(typeID), content, coverURL); err != nil {
		types, _ := h.service.GetBlogTypes()
		return h.RenderWithAuth(c, "admin/create_blog", fiber.Map{"Error": err.Error(), "Types": types})
	}
	return c.Redirect("/articles")
}

// 🔥 แก้ไข: HandleEditBlog รับค่าใหม่
func (h *FiberHandler) HandleEditBlog(c *fiber.Ctx) error {
	userID, isAdmin := getUserInfoFromContext(c)
	if !isAdmin {
		return c.Status(403).SendString("Unauthorized")
	}

	id, _ := c.ParamsInt("id")
	title := c.FormValue("title")
	shortTitle := c.FormValue("short_title")
	content := c.FormValue("content")
	coverURL := c.FormValue("cover_url")
	typeID, _ := strconv.Atoi(c.FormValue("blog_type_id"))

	if err := h.service.UpdateExistingBlog(uint(id), userID, isAdmin, title, shortTitle, uint(typeID), content, coverURL); err != nil {
		return h.RenderWithAuth(c, "admin/edit_blog", fiber.Map{"Error": err.Error()})
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

// --- API ---
func (h *FiberHandler) ApiSaveName(c *fiber.Ctx) error       { return nil }
func (h *FiberHandler) ApiDeleteName(c *fiber.Ctx) error     { return nil }
func (h *FiberHandler) ApiAnalyze(c *fiber.Ctx) error        { return nil }
func (h *FiberHandler) ApiGetLinguistics(c *fiber.Ctx) error { return nil }
