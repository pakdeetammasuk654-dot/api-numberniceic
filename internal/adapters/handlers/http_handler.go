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

// --- View Handlers ---

func (h *FiberHandler) ViewHome(c *fiber.Ctx) error {
	return c.Render("home", fiber.Map{}, "layouts/main")
}

func (h *FiberHandler) ViewDashboard(c *fiber.Ctx) error {
	return c.Render("dashboard", fiber.Map{}, "layouts/main")
}

func (h *FiberHandler) ViewArticles(c *fiber.Ctx) error {
	return c.Render("articles", fiber.Map{}, "layouts/main")
}

func (h *FiberHandler) ViewAbout(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{}, "layouts/main")
}

func (h *FiberHandler) ViewAnalysis(c *fiber.Ctx) error {
	return c.Render("analysis", fiber.Map{}, "layouts/main")
}

// HandleAnalysis (POST Form)
func (h *FiberHandler) HandleAnalysis(c *fiber.Ctx) error {
	name := c.FormValue("name")
	birthDay := c.FormValue("birth_day") // รับค่าวันเกิด

	result, err := h.service.AnalyzeName(name, birthDay)
	if err != nil {
		return c.Render("analysis", fiber.Map{
			"Error": err.Error(),
			"Name":  name,
		}, "layouts/main")
	}

	return c.Render("analysis", fiber.Map{
		"Result":   result,
		"Name":     name,
		"BirthDay": birthDay, // ส่งกลับไปให้ UI แสดงค่าที่เลือกไว้
	}, "layouts/main")
}

// --- API Handlers ---

func (h *FiberHandler) ApiAnalyze(c *fiber.Ctx) error {
	name := c.Query("name")
	birthDay := c.Query("birth_day") // รับค่าทาง query params

	if name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name is required"})
	}

	result, err := h.service.AnalyzeName(name, birthDay)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
