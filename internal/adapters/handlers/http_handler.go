package handlers

import (
	"api-numberniceic/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type FiberHandler struct {
	service ports.NumberService
}

// NewFiberHandler สร้าง Handler ใหม่ (ต้องชื่อนี้เพื่อให้ตรงกับ main.go)
func NewFiberHandler(service ports.NumberService) *FiberHandler {
	return &FiberHandler{
		service: service,
	}
}

// ViewIndex แสดงหน้าแรก (HTML Form)
func (h *FiberHandler) ViewIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

// ViewResult รับค่าจาก Form มาคำนวณและแสดงผล
func (h *FiberHandler) ViewResult(c *fiber.Ctx) error {
	name := c.FormValue("name") // รับค่าชื่อจาก input html

	// เรียกใช้ Service AnalyzeName ที่เราสร้างไว้
	result, err := h.service.AnalyzeName(name)
	if err != nil {
		return c.Render("index", fiber.Map{
			"Error": err.Error(),
			"Name":  name,
		})
	}

	// ส่งผลลัพธ์กลับไปแสดงที่หน้า index.html
	return c.Render("index", fiber.Map{
		"Result": result,
		"Name":   name,
	})
}

// ApiAnalyze สำหรับเรียกผ่าน API (JSON)
// ตัวอย่าง: GET /api/analyze?name=ทดสอบ
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
