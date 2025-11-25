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
	if err := godotenv.Overload(); err != nil {
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

	// 3. Auto Migrate: สร้าง/อัปเดตตารางใน Database
	db.AutoMigrate(&domain.User{}, &domain.SavedName{}, &domain.Blog{}, &domain.BlogType{})

	// 4. Setup Template Engine & Fiber
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/static", "./public")

	// 5. Init Layers (Repository, Service, Handler)
	repo := repositories.NewPostgresRepository(db)

	if err := repo.SeedBlogTypes(); err != nil {
		log.Println("⚠️ Warning: Failed to seed blog types:", err)
	}

	service := services.NewAnalyzerService(repo)
	handler := handlers.NewFiberHandler(service)

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// 6. Setup Routes

	// --- SEO & Static Files ---
	app.Get("/robots.txt", func(c *fiber.Ctx) error {
		return c.SendFile("./public/robots.txt")
	})
	app.Get("/sitemap.xml", handler.HandleSitemapXML)

	// --- General Pages ---
	app.Get("/", handler.ViewHome)
	app.Get("/about", handler.ViewAbout)
	app.Get("/sitemap", handler.ViewSitemap)

	// --- Analysis Features ---
	app.Get("/analysis", handler.ViewAnalysis)
	app.Post("/analysis", handler.HandleAnalysis)

	// --- Auth Routes ---
	app.Get("/login", authHandler.ViewLogin)
	app.Post("/login", authHandler.HandleLogin)
	app.Get("/register", authHandler.ViewRegister)
	app.Post("/register", authHandler.HandleRegister)
	app.Get("/logout", authHandler.HandleLogout)

	// --- Protected User Routes ---
	app.Get("/dashboard", middlewares.IsAuthenticated, handler.ViewDashboard)

	// --- Blog Public Routes ---
	app.Get("/articles", handler.ViewArticles)
	app.Get("/articles/:slug", handler.ViewBlogDetail)

	// --- Admin Routes (Protected + Check Admin) ---
	admin := app.Group("/admin", middlewares.IsAuthenticated)

	admin.Get("/", handler.ViewAdminPanel)

	// Blog Management
	admin.Get("/blogs", handler.ViewAdminBlogs)
	admin.Get("/create-blog", handler.ViewCreateBlog)
	admin.Post("/create-blog", handler.HandleCreateBlog)
	admin.Get("/edit-blog/:id", handler.ViewEditBlog)
	admin.Post("/edit-blog/:id", handler.HandleEditBlog)
	admin.Get("/delete-blog/:id", handler.HandleDeleteBlog)

	// Blog Type Management
	admin.Get("/types", handler.ViewAdminTypes)
	admin.Post("/types", handler.HandleCreateBlogType)
	admin.Get("/edit-type/:id", handler.ViewEditBlogType)
	admin.Post("/edit-type/:id", handler.HandleEditBlogType)
	admin.Get("/delete-type/:id", handler.HandleDeleteBlogType)

	// --- API Routes ---
	api := app.Group("/api")
	api.Post("/save-name", handler.ApiSaveName)
	api.Get("/delete-name/:id", middlewares.IsAuthenticated, handler.ApiDeleteName)
	api.Get("/analyze", handler.ApiAnalyze)
	api.Get("/linguistics", handler.ApiGetLinguistics)

	// 7. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	fmt.Printf("🚀 Server running at http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
