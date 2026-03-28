package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"vct-cms/internal/auth"
	"vct-cms/internal/blog"
	"vct-cms/internal/media"
	"vct-cms/internal/settings"
)

// @title VCT Platform CMS API
// @version 1.0
// @description Backend CMS API for managing VCT Platform content.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	app := fiber.New(fiber.Config{
		AppName:   "VCT CMS Backend v1.0",
		BodyLimit: 10 * 1024 * 1024, // 10MB max upload
	})

	// ── MIDDLEWARE ──
	app.Use(recover.New())
	app.Use(logger.New())

	// CORS — whitelist chính xác
	allowedOrigins := os.Getenv("CORS_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// ── API V1 ──
	api := app.Group("/api/v1")

	// Auth routes — Rate limit 5 requests/phút cho login (chống brute-force)
	authContext := api.Group("/auth")
	authContext.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Quá nhiều lần đăng nhập. Vui lòng thử lại sau 1 phút.",
			})
		},
	}))
	auth.RegisterRoutes(authContext)

	// Protected routes — JWT Token bắt buộc
	blogContext := api.Group("/blog", auth.JWTMiddleware())
	blog.RegisterRoutes(blogContext)

	mediaContext := api.Group("/media", auth.JWTMiddleware())
	media.RegisterRoutes(mediaContext)

	settingsContext := api.Group("/settings", auth.JWTMiddleware())
	settings.RegisterRoutes(settingsContext)

	// ── START SERVER ──
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 VCT CMS Backend starting on :%s", port)
	log.Printf("   CORS: %s", allowedOrigins)
	log.Printf("   JWT Secret: %s", func() string {
		if os.Getenv("JWT_SECRET") != "" {
			return "✅ From env"
		}
		return "⚠️  Using dev fallback"
	}())

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
