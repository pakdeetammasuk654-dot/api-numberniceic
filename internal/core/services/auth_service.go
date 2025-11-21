package services

import (
	"api-numberniceic/internal/core/domain"
	"api-numberniceic/internal/core/ports"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo ports.UserRepository
}

func NewAuthService(userRepo ports.UserRepository) ports.AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(username, email, password, displayName string) error {
	// 1. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 2. Map ข้อมูลให้ตรงกับตาราง users
	newUser := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		DisplayName:  displayName,
		IsAdmin:      false, // Default เป็น User ธรรมดา
	}

	return s.userRepo.CreateUser(newUser)
}

func (s *authService) Login(email, password string) (string, error) {
	// 1. ค้นหา User
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// 2. ตรวจสอบรหัสผ่าน
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// 3. สร้าง JWT Token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	claims["display_name"] = user.DisplayName

	// 🔥 เพิ่ม: ใส่สถานะ Admin ลงใน Token
	claims["is_admin"] = user.IsAdmin

	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "supersecretkey"
	}

	return token.SignedString([]byte(secret))
}
