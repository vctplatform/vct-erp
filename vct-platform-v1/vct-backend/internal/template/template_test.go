package template

import (
	"strings"
	"testing"
	"time"
)

func TestRegisterAndRender(t *testing.T) {
	eng := New()
	err := eng.Register("greeting", "Xin chào, {{.Name}}!")
	if err != nil {
		t.Fatal(err)
	}

	result, err := eng.Render("greeting", map[string]string{"Name": "Nguyễn Văn A"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "Xin chào, Nguyễn Văn A!" {
		t.Errorf("unexpected: %s", result)
	}
}

func TestRenderNotFound(t *testing.T) {
	eng := New()
	_, err := eng.Render("nonexistent", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestHelpers(t *testing.T) {
	eng := New()
	eng.Register("helpers", `{{upper .Name}} - {{year}}`)

	result, err := eng.Render("helpers", map[string]string{"Name": "test"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "TEST") {
		t.Error("upper helper failed")
	}
}

func TestDefaultHelper(t *testing.T) {
	eng := New()
	eng.Register("def", `{{default "N/A" .Value}}`)

	result, _ := eng.Render("def", map[string]string{"Value": ""})
	if result != "N/A" {
		t.Errorf("expected N/A, got %s", result)
	}

	result, _ = eng.Render("def", map[string]string{"Value": "Hello"})
	if result != "Hello" {
		t.Errorf("expected Hello, got %s", result)
	}
}

func TestFormatDate(t *testing.T) {
	eng := New()
	eng.Register("date", `{{formatDate .Date "02/01/2006"}}`)

	dt := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)
	result, err := eng.Render("date", map[string]any{"Date": dt})
	if err != nil {
		t.Fatal(err)
	}
	if result != "20/03/2026" {
		t.Errorf("expected 20/03/2026, got %s", result)
	}
}

func TestLayoutInheritance(t *testing.T) {
	eng := New()
	eng.AddLayout("email", `<html><body>{{block "content" .}}{{end}}<footer>VCT Platform</footer></body></html>`)

	err := eng.RegisterWithLayout("welcome", "email", `<h1>Chào mừng {{.Name}}</h1>`)
	if err != nil {
		t.Fatal(err)
	}

	result, err := eng.Render("welcome", map[string]string{"Name": "Athlete"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "<h1>Chào mừng Athlete</h1>") {
		t.Error("content not rendered")
	}
	if !strings.Contains(result, "VCT Platform") {
		t.Error("layout footer missing")
	}
}

func TestLayoutNotFound(t *testing.T) {
	eng := New()
	err := eng.RegisterWithLayout("test", "nonexistent", "content")
	if err == nil {
		t.Error("expected error for missing layout")
	}
}

func TestI18n(t *testing.T) {
	eng := New()
	eng.AddTranslations("vi", map[string]string{
		"welcome":    "Chào mừng",
		"tournament": "Giải đấu",
	})
	eng.AddTranslations("en", map[string]string{
		"welcome":    "Welcome",
		"tournament": "Tournament",
	})

	eng.Register("msg", `{{t .Locale "welcome"}} - {{t .Locale "tournament"}}`)

	result, _ := eng.Render("msg", map[string]string{"Locale": "vi"})
	if result != "Chào mừng - Giải đấu" {
		t.Errorf("vi: %s", result)
	}

	result, _ = eng.Render("msg", map[string]string{"Locale": "en"})
	if result != "Welcome - Tournament" {
		t.Errorf("en: %s", result)
	}
}

func TestI18n_FallbackToKey(t *testing.T) {
	eng := New()
	eng.AddTranslations("vi", map[string]string{})
	eng.Register("fb", `{{t "vi" "unknown_key"}}`)

	result, _ := eng.Render("fb", nil)
	if result != "unknown_key" {
		t.Errorf("should fall back to key, got %s", result)
	}
}

func TestRenderString(t *testing.T) {
	eng := New()
	result, err := eng.RenderString("Hello {{.Name}}!", map[string]string{"Name": "World"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "Hello World!" {
		t.Errorf("unexpected: %s", result)
	}
}

func TestHasAndList(t *testing.T) {
	eng := New()
	eng.Register("a", "template a")
	eng.Register("b", "template b")

	if !eng.Has("a") {
		t.Error("should have 'a'")
	}
	if eng.Has("c") {
		t.Error("should not have 'c'")
	}
	if len(eng.List()) != 2 {
		t.Errorf("expected 2 templates, got %d", len(eng.List()))
	}
}
