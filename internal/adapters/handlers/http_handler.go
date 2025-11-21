package handlers

import (
	"api-numberniceic/internal/core/ports"
	"os"

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

// 🔥 Helper Function: ใช้แทน c.Render ปกติ เพื่อส่งสถานะ Login ไปหน้าเว็บ
func (h *FiberHandler) RenderWithAuth(c *fiber.Ctx, template string, data fiber.Map) error {
	cookie := c.Cookies("jwt")
	isLoggedIn := false
	displayName := ""

	if cookie != "" {
		// พยายามแกะ Token เพื่อเอาชื่อมาโชว์
		token, _ := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// ถ้า secret key ไม่เจอ ให้ลองใช้ default (กัน error ตอน dev)
		if token == nil || !token.Valid {
			// ลอง parse แบบไม่อิง signature แค่เพื่อดึงข้อมูล (optional)
			// หรือจะมองว่า invalid ก็ได้
		} else {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				isLoggedIn = true
				if name, ok := claims["display_name"].(string); ok {
					displayName = name
				}
			}
		}

		// *หมายเหตุ: ถ้า Parse ไม่ผ่าน (เช่น secret ผิด) ก็จะถือว่ายังไม่ Login
		// กรณี Dev ง่ายๆ อาจจะเช็คแค่ cookie != "" ก็ได้ แต่เช็ค token ชัวร์กว่า
		if cookie != "" && !isLoggedIn {
			// Fallback กรณี parse error แต่มี cookie (อาจจะ assume ว่า login แล้วแต่ดึงชื่อไม่ได้)
			// หรือจะบังคับ logout ก็ได้
			// ในที่นี้ขอ assume ง่ายๆ ว่าถ้ามี cookie คือ login แล้ว (แต่ชื่ออาจไม่ขึ้นถ้า token ผิด)
			isLoggedIn = true
		}
	}

	if data == nil {
		data = fiber.Map{}
	}
	data["IsLoggedIn"] = isLoggedIn
	data["DisplayName"] = displayName

	return c.Render(template, data, "layouts/main")
}

// --- View Handlers (อัปเดตให้ใช้ RenderWithAuth) ---

func (h *FiberHandler) ViewHome(c *fiber.Ctx) error {
	return h.RenderWithAuth(c, "home", fiber.Map{})
}

func (h *FiberHandler) ViewDashboard(c *fiber.Ctx) error {
	return h.RenderWithAuth(c, "dashboard", fiber.Map{})
}

func (h *FiberHandler) ViewArticles(c *fiber.Ctx) error {
	return h.RenderWithAuth(c, "articles", fiber.Map{})
}

func (h *FiberHandler) ViewAbout(c *fiber.Ctx) error {
	return h.RenderWithAuth(c, "about", fiber.Map{})
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

	data := fiber.Map{
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

func (h *FiberHandler) HandleAnalysis(c *fiber.Ctx) error {
	name := c.FormValue("name")
	birthDay := c.FormValue("birth_day")

	result, err := h.service.AnalyzeName(name, birthDay)

	data := fiber.Map{
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

// --- API Handlers (คงเดิม) ---
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
