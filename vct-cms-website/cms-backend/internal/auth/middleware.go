package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware bảo vệ các route yêu cầu xác thực
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Lấy Header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Thiếu token xác thực. Vui lòng đăng nhập.",
			})
		}

		// Kiểm tra định dạng "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Định dạng token không hợp lệ. Sử dụng: Bearer <token>",
			})
		}

		// Giải mã và xác thực Token
		claims, err := ValidateToken(parts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token đã hết hạn hoặc không hợp lệ. Vui lòng đăng nhập lại.",
			})
		}

		// Lưu thông tin user vào context để các handler sau sử dụng
		c.Locals("userID", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RoleGuard chỉ cho phép user có role nhất định truy cập
func RoleGuard(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		for _, r := range allowedRoles {
			if role == r {
				return c.Next()
			}
		}
		return c.Status(403).JSON(fiber.Map{
			"error": "Bạn không có quyền thực hiện thao tác này (yêu cầu: " + strings.Join(allowedRoles, "/") + ")",
		})
	}
}
