package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	if cookie == "" {
		return c.Redirect("/login")
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "supersecretkey"
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Redirect("/login")
	}

	// ถ้าผ่าน ให้ไปต่อได้
	return c.Next()
}
