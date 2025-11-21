package main

import (
	"api-numberniceic/internal/adapters/handlers"
	"api-numberniceic/internal/adapters/middlewares"
	"api-numberniceic/internal/adapters/repositories"
	"api-numberniceic/internal/core/domain"
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

	// 2. ตั้งค่าเชื่อมต่อ Database
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

	// 🔥 Auto Migrate: เพิ่ม SavedName
	db.AutoMigrate(&domain.User{}, &domain.SavedName{})

	// 3. Setup Template Engine & Fiber
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/static", "./public")

	// 4. Init Layers
	repo := repositories.NewPostgresRepository(db)
	service := services.NewAnalyzerService(repo)
	handler := handlers.NewFiberHandler(service)

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// 5. Setup Routes
	app.Get("/", handler.ViewHome)
	app.Get("/about", handler.ViewAbout)
	app.Get("/articles", handler.ViewArticles)

	app.Get("/analysis", handler.ViewAnalysis)
	app.Post("/analysis", handler.HandleAnalysis)

	// Auth Routes
	app.Get("/login", authHandler.ViewLogin)
	app.Post("/login", authHandler.HandleLogin)
	app.Get("/register", authHandler.ViewRegister)
	app.Post("/register", authHandler.HandleRegister)
	app.Get("/logout", authHandler.HandleLogout)

	// Protected Routes
	app.Get("/dashboard", middlewares.IsAuthenticated, handler.ViewDashboard)

	// 🔥 API Routes (Updated)
	api := app.Group("/api")
	api.Post("/save-name", handler.ApiSaveName)                                     // API บันทึกชื่อ
	api.Get("/delete-name/:id", middlewares.IsAuthenticated, handler.ApiDeleteName) // ลบชื่อ
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
