package blog

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Khởi tạo Service
var service = NewBlogService(`d:\VCT PLATFORM\vct-website`)

func RegisterRoutes(router fiber.Router) {
	router.Get("/", getPosts)
	router.Post("/", createPost)
	router.Put("/:id", updatePost)
	router.Post("/upload", uploadImage)
}

// Lấy danh sách bài viết
func getPosts(c *fiber.Ctx) error {
	posts, err := service.GetPosts()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": posts})
}

// Thêm bài viết mới (sử dụng atomic CreatePost)
func createPost(c *fiber.Ctx) error {
	var newPost Post
	if err := c.BodyParser(&newPost); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dữ liệu không hợp lệ"})
	}

	if err := service.CreatePost(newPost); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Không thể lưu bài viết: " + err.Error()})
	}

	commitMsg := fmt.Sprintf("content: add blog [%s]", newPost.Title)
	if err := service.SyncToGit(commitMsg); err != nil {
		return c.Status(200).JSON(fiber.Map{"status": "warning", "message": "Lưu OK, lỗi Git: " + err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Đã thêm bài viết và push lên GitHub", "data": newPost})
}

// Cập nhật bài viết (sử dụng atomic UpdatePost)
func updatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedPost Post
	if err := c.BodyParser(&updatedPost); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dữ liệu không hợp lệ"})
	}

	if err := service.UpdatePost(id, updatedPost); err != nil {
		if err.Error() == "not_found" {
			return c.Status(404).JSON(fiber.Map{"error": "Không tìm thấy bài viết"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Lỗi cập nhật: " + err.Error()})
	}

	commitMsg := fmt.Sprintf("content: update blog [%s]", updatedPost.Title)
	if err := service.SyncToGit(commitMsg); err != nil {
		return c.Status(200).JSON(fiber.Map{"status": "warning", "message": "Lưu OK, lỗi Git: " + err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Cập nhật bài viết thành công", "data": updatedPost})
}

// allowedImageExts — danh sách extension ảnh hợp lệ
var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".webp": true, ".svg": true, ".avif": true,
}

// Upload hình ảnh — sanitized filename + extension whitelist
func uploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Không tìm thấy file upload"})
	}

	// ── SECURITY: Sanitize filename ──
	// 1. Chỉ lấy basename (loại bỏ path traversal ../../../)
	safeName := filepath.Base(file.Filename)

	// 2. Kiểm tra extension hợp lệ
	ext := strings.ToLower(filepath.Ext(safeName))
	if !allowedImageExts[ext] {
		return c.Status(400).JSON(fiber.Map{"error": "Định dạng file không được phép. Chỉ chấp nhận: jpg, png, gif, webp, svg, avif"})
	}

	// 3. Tạo tên file unique bằng timestamp
	filename := fmt.Sprintf("%d_%s", time.Now().UnixMilli(), safeName)

	// 4. Build path an toàn bằng filepath.Join (không dùng string concat)
	savePath := filepath.Join(service.ImagesDir, filename)

	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Không thể lưu tệp tin: " + err.Error()})
	}

	relativePath := "assets/images/" + filename

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Upload hình ảnh thành công",
		"url":     relativePath,
	})
}
