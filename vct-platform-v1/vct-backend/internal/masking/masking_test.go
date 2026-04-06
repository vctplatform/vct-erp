package masking

import (
	"encoding/json"
	"testing"
)

func TestStrategy_Full(t *testing.T) {
	if got := Full("secret123"); got != "*********" {
		t.Errorf("expected 9 asterisks, got %q", got)
	}
	if got := Full(""); got != "" {
		t.Errorf("empty should stay empty, got %q", got)
	}
}

func TestStrategy_Email(t *testing.T) {
	tests := []struct{ in, want string }{
		{"user@example.com", "u***@example.com"},
		{"ab@test.vn", "a*@test.vn"},
		{"x@y.com", "*@y.com"},
		{"invalid", "*******"},
	}
	for _, tt := range tests {
		if got := Email(tt.in); got != tt.want {
			t.Errorf("Email(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestStrategy_Phone(t *testing.T) {
	got := Phone("0912-345-6789")
	if got != "*******6789" {
		t.Errorf("got %q", got)
	}
}

func TestStrategy_IDCard(t *testing.T) {
	got := IDCard("079123456789")
	if got != "079*******89" {
		t.Errorf("got %q", got)
	}
}

func TestStrategy_Partial(t *testing.T) {
	if got := Partial("Nguyen"); got != "N****n" {
		t.Errorf("got %q", got)
	}
	if got := Partial("AB"); got != "**" {
		t.Errorf("short string: got %q", got)
	}
}

func TestStrategy_Redact(t *testing.T) {
	if got := Redact("anything"); got != "[REDACTED]" {
		t.Errorf("got %q", got)
	}
}

func TestMasker_DefaultRules(t *testing.T) {
	m := NewMasker()

	if got := m.MaskValue("email", "user@test.com"); got == "user@test.com" {
		t.Error("email should be masked")
	}
	if got := m.MaskValue("password", "secret"); got != "[REDACTED]" {
		t.Errorf("password should be redacted, got %q", got)
	}
	if got := m.MaskValue("name", "Nguyễn Văn A"); got != "Nguyễn Văn A" {
		t.Error("name has no rule, should not be masked")
	}
}

func TestMasker_CustomRule(t *testing.T) {
	m := NewEmptyMasker()
	m.AddRule("bank_account", Partial, false)

	got := m.MaskValue("bank_account", "123456789")
	if got != "1*******9" {
		t.Errorf("got %q", got)
	}
}

func TestMasker_RegexRule(t *testing.T) {
	m := NewEmptyMasker()
	m.AddRule(".*_secret$", Redact, true)

	if got := m.MaskValue("api_secret", "key123"); got != "[REDACTED]" {
		t.Error("regex rule should match api_secret")
	}
	if got := m.MaskValue("name", "visible"); got != "visible" {
		t.Error("name should not match regex")
	}
}

func TestMaskMap(t *testing.T) {
	m := NewMasker()
	data := map[string]string{
		"email": "admin@vct.vn",
		"phone": "0901234567",
		"name":  "Trần B",
	}

	masked := m.MaskMap(data)
	if masked["email"] == "admin@vct.vn" {
		t.Error("email should be masked")
	}
	if masked["name"] != "Trần B" {
		t.Error("name should not be masked")
	}
}

func TestMaskJSON(t *testing.T) {
	m := NewMasker()
	input := `{"email":"user@test.com","password":"abc123","profile":{"phone":"0912345678","name":"Test"}}`

	masked, err := m.MaskJSON([]byte(input))
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	json.Unmarshal(masked, &result)

	if result["email"] == "user@test.com" {
		t.Error("email should be masked")
	}
	if result["password"] != "[REDACTED]" {
		t.Error("password should be redacted")
	}

	profile := result["profile"].(map[string]any)
	if profile["phone"] == "0912345678" {
		t.Error("nested phone should be masked")
	}
	if profile["name"] != "Test" {
		t.Error("name should not be masked")
	}
}

func TestMaskJSON_VietnamPII(t *testing.T) {
	m := NewMasker()
	input := `{"cccd":"079123456789","cmnd":"215678901"}`

	masked, err := m.MaskJSON([]byte(input))
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	json.Unmarshal(masked, &result)

	if result["cccd"] == "079123456789" {
		t.Error("CCCD should be masked")
	}
	if result["cmnd"] == "215678901" {
		t.Error("CMND should be masked")
	}
}
