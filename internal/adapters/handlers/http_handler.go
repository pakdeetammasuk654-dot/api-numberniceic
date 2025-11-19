package handlers

import (
	"api-numberniceic/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type FiberHandler struct {
	service ports.NumberService
}

func NewFiberHandler(service ports.NumberService) *FiberHandler {
	return &FiberHandler{
		service: service,
	}
}

// --- หน้า Home (Landing Page) ---
func (h *FiberHandler) ViewHome(c *fiber.Ctx) error {
	// ใช้ layout "layouts/main"
	return c.Render("home", fiber.Map{}, "layouts/main")
}

// --- หน้า Dashboard (ยังไม่มีเนื้อหา ใส่ placeholder ไว้ก่อน) ---
func (h *FiberHandler) ViewDashboard(c *fiber.Ctx) error {
	return c.Render("dashboard", fiber.Map{}, "layouts/main")
}

// --- หน้า บทความ ---
func (h *FiberHandler) ViewArticles(c *fiber.Ctx) error {
	return c.Render("articles", fiber.Map{}, "layouts/main")
}

// --- หน้า เกี่ยวกับเรา ---
func (h *FiberHandler) ViewAbout(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{}, "layouts/main")
}

// --- หน้า วิเคราะห์ชื่อ (Analysis) ---
func (h *FiberHandler) ViewAnalysis(c *fiber.Ctx) error {
	return c.Render("analysis", fiber.Map{}, "layouts/main")
}

// HandleAnalysis รับค่าจาก Form และแสดงผล
func (h *FiberHandler) HandleAnalysis(c *fiber.Ctx) error {
	name := c.FormValue("name") // รับค่าชื่อจาก input html

	result, err := h.service.AnalyzeName(name)
	if err != nil {
		return c.Render("analysis", fiber.Map{
			"Error": err.Error(),
			"Name":  name,
		}, "layouts/main") // อย่าลืมใส่ layout
	}

	return c.Render("analysis", fiber.Map{
		"Result": result,
		"Name":   name,
	}, "layouts/main") // อย่าลืมใส่ layout
}

// --- API (JSON) ---
func (h *FiberHandler) ApiAnalyze(c *fiber.Ctx) error {
	name := c.Query("name")
	if name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name is required"})
	}

	result, err := h.service.AnalyzeName(name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
