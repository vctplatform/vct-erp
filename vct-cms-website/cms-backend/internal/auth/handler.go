package auth

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// ─── MOCK USER STORE ───
// TODO: Thay thế bằng GORM + PostgreSQL trong production
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Hash, không bao giờ trả về client
	Role     string `json:"role"`
}

// Tạo hash password cho admin mặc định
func hashPassword(pw string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(bytes)
}

// Mock database - 2 users mặc định
var users = []User{
	{ID: 1, Username: "admin", Password: hashPassword("admin123"), Role: "admin"},
	{ID: 2, Username: "editor", Password: hashPassword("editor123"), Role: "editor"},
}

// ─── ROUTES ───
func RegisterRoutes(router fiber.Router) {
	router.Post("/login", login)
	router.Get("/me", JWTMiddleware(), getMe)
}

// ─── LOGIN ───
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dữ liệu đăng nhập không hợp lệ"})
	}

	// Tìm user
	var found *User
	for _, u := range users {
		if u.Username == req.Username {
			found = &u
			break
		}
	}

	if found == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Tên đăng nhập hoặc mật khẩu không đúng"})
	}

	// So sánh mật khẩu bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Tên đăng nhập hoặc mật khẩu không đúng"})
	}

	// Tạo JWT Token
	token, err := GenerateToken(found.ID, found.Username, found.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Không thể tạo token"})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Đăng nhập thành công",
		"token":   token,
		"user": fiber.Map{
			"id":       found.ID,
			"username": found.Username,
			"role":     found.Role,
		},
	})
}

// ─── GET CURRENT USER ───
func getMe(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "success",
		"user": fiber.Map{
			"id":       c.Locals("userID"),
			"username": c.Locals("username"),
			"role":     c.Locals("role"),
		},
	})
}
