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
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: No .env file found")
	}

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

	db.AutoMigrate(&domain.User{}, &domain.SavedName{}, &domain.Blog{})

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/static", "./public")

	repo := repositories.NewPostgresRepository(db)
	service := services.NewAnalyzerService(repo)
	handler := handlers.NewFiberHandler(service)

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	app.Get("/", handler.ViewHome)
	app.Get("/about", handler.ViewAbout)

	app.Get("/analysis", handler.ViewAnalysis)
	app.Post("/analysis", handler.HandleAnalysis)

	app.Get("/login", authHandler.ViewLogin)
	app.Post("/login", authHandler.HandleLogin)
	app.Get("/register", authHandler.ViewRegister)
	app.Post("/register", authHandler.HandleRegister)
	app.Get("/logout", authHandler.HandleLogout)

	app.Get("/dashboard", middlewares.IsAuthenticated, handler.ViewDashboard)

	app.Get("/articles", handler.ViewArticles)
	app.Get("/articles/:id", handler.ViewBlogDetail)

	// 🔥 Admin Routes
	admin := app.Group("/admin", middlewares.IsAuthenticated)

	admin.Get("/", handler.ViewAdminPanel) // 🔥 หน้าหลัก Admin

	admin.Get("/blogs", handler.ViewAdminBlogs)
	admin.Get("/create-blog", handler.ViewCreateBlog)
	admin.Post("/create-blog", handler.HandleCreateBlog)
	admin.Get("/edit-blog/:id", handler.ViewEditBlog)
	admin.Post("/edit-blog/:id", handler.HandleEditBlog)
	admin.Get("/delete-blog/:id", handler.HandleDeleteBlog)

	api := app.Group("/api")
	api.Post("/save-name", handler.ApiSaveName)
	api.Get("/delete-name/:id", middlewares.IsAuthenticated, handler.ApiDeleteName)
	api.Get("/analyze", handler.ApiAnalyze)
	api.Get("/linguistics", handler.ApiGetLinguistics)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	fmt.Printf("🚀 Server running at http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
