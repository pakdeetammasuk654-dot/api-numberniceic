package main

import (
	"api-numberniceic/internal/adapters/handlers"
	"api-numberniceic/internal/adapters/middlewares" // ✅ เพิ่ม Middleware
	"api-numberniceic/internal/adapters/repositories"
	"api-numberniceic/internal/core/domain" // ✅ เพิ่ม Domain เพื่อใช้ Migrate User
	"api-numberniceic/internal/core/services"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. โหลด Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: No .env file found")
	}

	// 2. ตั้งค่าเชื่อมต่อ Database (อ่านจาก env)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("✅ Connected to Database successfully")

	// 🔥 Auto Migrate User Table (สร้างตาราง users อัตโนมัติ ถ้ายังไม่มี)
	db.AutoMigrate(&domain.User{})

	// 3. Setup Template Engine & Fiber
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/static", "./public")

	// 4. Init Layers (ประกาศตัวแปร Service & Handler)

	// --- Existing (ระบบวิเคราะห์ชื่อ) ---
	repo := repositories.NewPostgresRepository(db)
	service := services.NewAnalyzerService(repo)
	handler := handlers.NewFiberHandler(service)

	// --- 🔥 New Auth Layers (ระบบสมาชิก) ---
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// 5. Setup Routes

	// --- General Pages ---
	app.Get("/", handler.ViewHome)
	app.Get("/about", handler.ViewAbout)
	app.Get("/articles", handler.ViewArticles)

	// --- Analysis Feature ---
	app.Get("/analysis", handler.ViewAnalysis)
	app.Post("/analysis", handler.HandleAnalysis)

	// --- 🔥 Auth Routes (Login/Register) ---
	app.Get("/login", authHandler.ViewLogin)
	app.Post("/login", authHandler.HandleLogin)
	app.Get("/register", authHandler.ViewRegister)
	app.Post("/register", authHandler.HandleRegister)
	app.Get("/logout", authHandler.HandleLogout)

	// --- 🔥 Protected Routes (ต้อง Login ก่อนถึงเข้าได้) ---
	// ใช้ middlewares.IsAuthenticated คั่นไว้หน้า ViewDashboard
	app.Get("/dashboard", middlewares.IsAuthenticated, handler.ViewDashboard)

	// --- API Routes ---
	api := app.Group("/api")
	api.Get("/analyze", handler.ApiAnalyze)
	api.Get("/linguistics", handler.ApiGetLinguistics)

	// 6. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Printf("🚀 Server running at http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
