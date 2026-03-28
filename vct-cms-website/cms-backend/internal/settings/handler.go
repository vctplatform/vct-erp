package settings

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"vct-cms/internal/auth"
)

// Khởi tạo Services — trỏ tới thư mục chứa website
var i18nService = NewI18nService(`d:\VCT PLATFORM\vct-website`)
var siteService = NewSiteService(`d:\VCT PLATFORM\vct-website`)

func RegisterRoutes(router fiber.Router) {
	// i18n
	router.Get("/languages", getLanguages)
	router.Put("/languages", auth.RoleGuard("admin"), saveLanguages)

	// Site Config
	router.Get("/site", getSiteConfig)
	router.Put("/site", auth.RoleGuard("admin"), saveSiteConfig)

	// Website Status Check
	router.Get("/status", checkWebsiteStatus)
}

// ─── i18n ───

func getLanguages(c *fiber.Ctx) error {
	entries, err := i18nService.GetMergedTranslations()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": entries, "totalKeys": len(entries)})
}

func saveLanguages(c *fiber.Ctx) error {
	var body struct {
		Entries []TranslationEntry `json:"entries"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dữ liệu không hợp lệ"})
	}
	if err := i18nService.SaveTranslations(body.Entries); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Lỗi ghi file: " + err.Error()})
	}
	if err := i18nService.SyncToGit("i18n: update translations"); err != nil {
		return c.Status(200).JSON(fiber.Map{"status": "warning", "message": "Lưu OK, lỗi Git: " + err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Đã cập nhật bản dịch và sync lên GitHub"})
}

// ─── SITE CONFIG ───

func getSiteConfig(c *fiber.Ctx) error {
	config, err := siteService.GetSiteConfig()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": config})
}

func saveSiteConfig(c *fiber.Ctx) error {
	var config SiteConfig
	if err := c.BodyParser(&config); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dữ liệu không hợp lệ"})
	}
	if err := siteService.SaveSiteConfig(&config); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Lỗi ghi site.json: " + err.Error()})
	}
	if err := siteService.SyncToGit("settings: update site config"); err != nil {
		return c.Status(200).JSON(fiber.Map{"status": "warning", "message": "Lưu OK, lỗi Git: " + err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Đã cập nhật cấu hình và sync lên GitHub"})
}

// ─── WEBSITE STATUS ───

func checkWebsiteStatus(c *fiber.Ctx) error {
	url := "https://vct-platform.github.io/vct-website/"
	resp, err := http.Head(url)
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "Không thể kết nối: " + err.Error()})
	}
	defer resp.Body.Close()

	return c.JSON(fiber.Map{
		"status":     "success",
		"httpStatus": resp.StatusCode,
		"online":     resp.StatusCode == 200,
		"url":        url,
	})
}
