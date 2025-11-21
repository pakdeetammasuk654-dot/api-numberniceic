package handlers

import (
	"api-numberniceic/internal/core/ports"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service ports.AuthService
}

func NewAuthHandler(service ports.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) ViewLogin(c *fiber.Ctx) error {
	// ใช้ Render ธรรมดาก็ได้ หรือ RenderWithAuth ก็ดี (แต่หน้า Login ปกติไม่ต้องโชว์ User menu)
	return c.Render("auth/login", fiber.Map{}, "layouts/main")
}

func (h *AuthHandler) ViewRegister(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{}, "layouts/main")
}

func (h *AuthHandler) HandleRegister(c *fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	displayName := c.FormValue("display_name") // รับค่าชื่อเล่น

	// ส่ง display_name ไป Service
	if err := h.service.Register(username, email, password, displayName); err != nil {
		return c.Render("auth/register", fiber.Map{"Error": "สมัครสมาชิกไม่สำเร็จ: " + err.Error()}, "layouts/main")
	}

	return c.Redirect("/login")
}

func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	token, err := h.service.Login(email, password)
	if err != nil {
		return c.Render("auth/login", fiber.Map{"Error": "เข้าสู่ระบบไม่สำเร็จ: " + err.Error()}, "layouts/main")
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HTTPOnly = true
	c.Cookie(cookie)

	return c.Redirect("/dashboard")
}

func (h *AuthHandler) HandleLogout(c *fiber.Ctx) error {
	c.ClearCookie("jwt")
	return c.Redirect("/")
}
