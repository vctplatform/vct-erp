package export

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testColumns = []Column{
	{Key: "id", Header: "ID"},
	{Key: "name", Header: "Tên VĐV"},
	{Key: "belt", Header: "Đai"},
}

func TestCSVExport(t *testing.T) {
	var buf bytes.Buffer
	exp := NewExporter(&buf, FormatCSV, testColumns)

	exp.WriteHeader()
	exp.WriteRow(map[string]string{"id": "1", "name": "Nguyễn Văn A", "belt": "Đai Đen"})
	exp.WriteRow(map[string]string{"id": "2", "name": "Trần Thị B", "belt": "Đai Đỏ"})
	exp.Flush()

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
	if !strings.Contains(lines[0], "Tên VĐV") {
		t.Error("header should contain Vietnamese column name")
	}
	if exp.Count() != 2 {
		t.Errorf("expected count 2, got %d", exp.Count())
	}
}

func TestJSONExport(t *testing.T) {
	var buf bytes.Buffer
	exp := NewExporter(&buf, FormatJSON, testColumns)

	exp.WriteHeader() // no-op for JSON
	exp.WriteRow(map[string]string{"id": "1", "name": "Nguyễn Văn A", "belt": "Đai Đen", "extra": "ignored"})
	exp.Flush()

	output := buf.String()
	if !strings.Contains(output, `"name":"Nguyễn Văn A"`) {
		t.Error("JSON should contain name field")
	}
	if strings.Contains(output, `"extra"`) {
		t.Error("extra fields outside column list should not be exported")
	}
}

func TestFieldSelection(t *testing.T) {
	cols := []Column{{Key: "name", Header: "Name"}}
	var buf bytes.Buffer
	exp := NewExporter(&buf, FormatCSV, cols)

	exp.WriteHeader()
	exp.WriteRow(map[string]string{"id": "1", "name": "Test", "belt": "Black"})
	exp.Flush()

	output := buf.String()
	if strings.Contains(output, "Black") {
		t.Error("non-selected fields should not appear")
	}
}

func TestParseFormat(t *testing.T) {
	if ParseFormat("JSON") != FormatJSON {
		t.Error("should parse JSON")
	}
	if ParseFormat("csv") != FormatCSV {
		t.Error("should parse csv")
	}
	if ParseFormat("unknown") != FormatCSV {
		t.Error("unknown should default to CSV")
	}
}

func TestDownloadHandler_CSV(t *testing.T) {
	rows := []map[string]string{
		{"id": "1", "name": "Athlete A", "belt": "Black"},
		{"id": "2", "name": "Athlete B", "belt": "Red"},
	}
	idx := 0
	provider := func() (map[string]string, error) {
		if idx >= len(rows) {
			return nil, io.EOF
		}
		row := rows[idx]
		idx++
		return row, nil
	}

	handler := DownloadHandler("athletes", FormatCSV, testColumns, provider)
	req := httptest.NewRequest("GET", "/export/athletes", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Header().Get("Content-Disposition"), "athletes_") {
		t.Error("should have filename in Content-Disposition")
	}
	if rec.Header().Get("Content-Type") != "text/csv; charset=utf-8" {
		t.Errorf("unexpected content-type: %s", rec.Header().Get("Content-Type"))
	}

	body := rec.Body.String()
	lines := strings.Split(strings.TrimSpace(body), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestDownloadHandler_JSONOverride(t *testing.T) {
	provider := func() (map[string]string, error) {
		return nil, io.EOF
	}

	handler := DownloadHandler("athletes", FormatCSV, testColumns, provider)
	req := httptest.NewRequest("GET", "/export/athletes?format=json", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Type") != "application/x-ndjson" {
		t.Errorf("should override to NDJSON, got %s", rec.Header().Get("Content-Type"))
	}
}

func TestDownloadHandler_CorrectHeaders(t *testing.T) {
	provider := func() (map[string]string, error) { return nil, io.EOF }
	handler := DownloadHandler("test", FormatCSV, testColumns, provider)

	req := httptest.NewRequest("GET", "/export", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Cache-Control") != "no-cache" {
		t.Error("should set Cache-Control: no-cache")
	}
	disp := rec.Header().Get("Content-Disposition")
	if !strings.HasPrefix(disp, `attachment; filename="test_`) {
		t.Errorf("unexpected Content-Disposition: %s", disp)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/csv; charset=utf-8" {
		t.Errorf("unexpected Content-Type: %s", ct)
	}

	// Status should work
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
