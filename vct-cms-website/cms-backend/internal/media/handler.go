package media

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router) {
	router.Post("/upload", uploadMedia)
}

func uploadMedia(c *fiber.Ctx) error {
	// Handle multipart form, save to assets/images
	return c.JSON(fiber.Map{"status": "success", "message": "Media upload handler"})
}
