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
	// **แก้ข้อมูลตรงนี้ให้ตรงกับเครื่องของคุณนะครับ**
	dsn := "host=localhost user=tayap password=IntelliP24.X dbname=tayap port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("✅ Connected to Database successfully")

	// 2. Setup Template Engine (สำหรับหน้าเว็บ)
	engine := html.New("./views", ".html")

	// 3. Setup Fiber App
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./public")

	// 4. Init Layers (ต่อจิกซอว์เข้าด้วยกัน)
	// สร้าง Repository (ต่อ DB)
	repo := repositories.NewPostgresRepository(db)

	// สร้าง Service (ต่อ Repository)
	service := services.NewAnalyzerService(repo)

	// สร้าง Handler (ต่อ Service)
	handler := handlers.NewFiberHandler(service)

	// 5. Setup Routes
	// หน้าเว็บ
	app.Get("/", handler.ViewIndex)
	app.Post("/", handler.ViewResult)

	// API (JSON)
	api := app.Group("/api")
	api.Get("/analyze", handler.ApiAnalyze)

	// 6. Start Server
	fmt.Println("🚀 Server running at http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
