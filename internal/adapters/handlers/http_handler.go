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

// ViewAnalysis: แก้ไขให้รับค่าจาก Query Params (สำหรับการคลิกจากตาราง)
func (h *FiberHandler) ViewAnalysis(c *fiber.Ctx) error {
	// 1. รับค่าจาก Link (Query Params) ก่อน
	// เช่น /analysis?name=สมชาย&birth_day=monday
	name := c.Query("name")
	birthDay := c.Query("birth_day")

	// 2. ถ้าไม่มีค่าส่งมา (เปิดหน้าเว็บครั้งแรก) ให้ใช้ค่า Default "ณเดชน์"
	if name == "" {
		name = "ณเดชน์"
	}
	if birthDay == "" {
		birthDay = "sunday"
	}

	// 3. สั่ง Service ให้วิเคราะห์
	result, err := h.service.AnalyzeName(name, birthDay)

	// 4. เตรียมข้อมูลสำหรับส่งไปแสดงผล
	data := fiber.Map{
		"Name":     name,
		"BirthDay": birthDay,
	}

	// ถ้าไม่มี Error ให้ส่งผลลัพธ์ (Result) ไปด้วย
	if err == nil {
		data["Result"] = result
	} else {
		data["Error"] = "ไม่สามารถโหลดข้อมูลได้: " + err.Error()
	}

	return c.Render("analysis", data, "layouts/main")
}

// HandleAnalysis (POST Form): สำหรับกรณีที่ User พิมพ์ชื่อแล้วกด Enter (ระบบ Auto Submit)
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

func (h *FiberHandler) ApiGetLinguistics(c *fiber.Ctx) error {
	name := c.Query("name")
	if name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name is required"})
	}

	// เรียกใช้ Service
	meaning, err := h.service.GetNameLinguistics(name)
	if err != nil {
		// กรณี AI มีปัญหา หรือไม่ได้ใส่ Key
		return c.Status(500).JSON(fiber.Map{"error": "AI Service Unavailable: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"text": meaning,
	})
}
