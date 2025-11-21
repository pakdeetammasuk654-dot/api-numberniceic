package ports

import "api-numberniceic/internal/core/domain"

type UserRepository interface {
	CreateUser(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}

type AuthService interface {
	// เพิ่ม displayName เข้าไปในพารามิเตอร์
	Register(username, email, password, displayName string) error
	Login(email, password string) (string, error)
}
