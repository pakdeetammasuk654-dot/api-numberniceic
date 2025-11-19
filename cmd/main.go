package main

import (
	"api-numberniceic/internal/adapters/handlers"
	"api-numberniceic/internal/adapters/repositories"
	"api-numberniceic/internal/core/services"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. ตั้งค่าเชื่อมต่อ Database (PostgreSQL)
	dsn := "host=localhost user=tayap password=IntelliP24.X dbname=tayap port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("✅ Connected to Database successfully")

	// 2. Setup Template Engine
	engine := html.New("./views", ".html")

	// 3. Setup Fiber App
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// บอกให้ /static ชี้ไปที่โฟลเดอร์ public (สำหรับ CSS)
	app.Static("/static", "./public")

	// 4. Init Layers
	repo := repositories.NewPostgresRepository(db)
	service := services.NewAnalyzerService(repo)
	handler := handlers.NewFiberHandler(service)

	// 5. Setup Routes (เส้นทางใหม่)

	// หน้าแรก (Home)
	app.Get("/", handler.ViewHome)

	// เมนูอื่นๆ
	app.Get("/dashboard", handler.ViewDashboard)
	app.Get("/articles", handler.ViewArticles)
	app.Get("/about", handler.ViewAbout)

	// หน้าวิเคราะห์ชื่อ (Analysis)
	app.Get("/analysis", handler.ViewAnalysis)    // แสดงฟอร์ม
	app.Post("/analysis", handler.HandleAnalysis) // กดส่งฟอร์ม

	// API (JSON)
	api := app.Group("/api")
	api.Get("/analyze", handler.ApiAnalyze)

	// 6. Start Server
	fmt.Println("🚀 Server running at http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
