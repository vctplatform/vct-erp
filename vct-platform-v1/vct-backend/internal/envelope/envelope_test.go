package envelope

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOK(t *testing.T) {
	resp := OK(map[string]string{"name": "VCT"})
	if !resp.Success {
		t.Error("should be success")
	}
	if resp.Data == nil {
		t.Error("data should not be nil")
	}
	if resp.Error != nil {
		t.Error("error should be nil")
	}
	if resp.Timestamp == "" {
		t.Error("should have timestamp")
	}
}

func TestPaginated(t *testing.T) {
	resp := Paginated([]string{"a", "b"}, 2, 10, 25)
	if !resp.Success {
		t.Error("should be success")
	}
	if resp.Meta == nil {
		t.Fatal("meta should not be nil")
	}
	if resp.Meta.Page != 2 {
		t.Errorf("expected page 2, got %d", resp.Meta.Page)
	}
	if resp.Meta.TotalPages != 3 {
		t.Errorf("expected 3 total pages, got %d", resp.Meta.TotalPages)
	}
}

func TestPaginated_ExactDivision(t *testing.T) {
	resp := Paginated(nil, 1, 10, 20)
	if resp.Meta.TotalPages != 2 {
		t.Errorf("expected 2 pages for 20/10, got %d", resp.Meta.TotalPages)
	}
}

func TestErr(t *testing.T) {
	resp := Err("NOT_FOUND", "Athlete not found")
	if resp.Success {
		t.Error("should not be success")
	}
	if resp.Error == nil {
		t.Fatal("error should not be nil")
	}
	if resp.Error.Code != "NOT_FOUND" {
		t.Errorf("expected NOT_FOUND, got %s", resp.Error.Code)
	}
}

func TestValidationError(t *testing.T) {
	resp := ValidationError(map[string]string{
		"email": "invalid format",
		"name":  "required",
	})
	if resp.Error.Code != "VALIDATION_ERROR" {
		t.Error("should be VALIDATION_ERROR")
	}
	if len(resp.Error.Details) != 2 {
		t.Errorf("expected 2 details, got %d", len(resp.Error.Details))
	}
}

func TestWithRequestID(t *testing.T) {
	resp := OK("data").WithRequestID("req-123")
	if resp.Meta == nil || resp.Meta.RequestID != "req-123" {
		t.Error("request ID not set")
	}
}

func TestCommonErrors(t *testing.T) {
	tests := []struct {
		name string
		resp *Response
		code string
	}{
		{"NotFound", NotFound("athlete", "123"), "NOT_FOUND"},
		{"BadRequest", BadRequest("invalid"), "BAD_REQUEST"},
		{"Unauthorized", Unauthorized("no token"), "UNAUTHORIZED"},
		{"Forbidden", Forbidden("no access"), "FORBIDDEN"},
		{"InternalError", InternalError("oops"), "INTERNAL_ERROR"},
		{"Conflict", Conflict("duplicate"), "CONFLICT"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resp.Error.Code != tt.code {
				t.Errorf("expected %s, got %s", tt.code, tt.resp.Error.Code)
			}
		})
	}
}

func TestWriteOK(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteOK(rec, map[string]string{"id": "1"})

	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("content-type: %s", ct)
	}

	var resp Response
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp.Success {
		t.Error("should be success")
	}
}

func TestWriteCreated(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteCreated(rec, map[string]string{"id": "new"})

	if rec.Code != 201 {
		t.Errorf("expected 201, got %d", rec.Code)
	}
}

func TestWriteErr_StatusMapping(t *testing.T) {
	tests := []struct {
		resp   *Response
		status int
	}{
		{NotFound("x", "1"), http.StatusNotFound},
		{BadRequest("x"), http.StatusBadRequest},
		{Unauthorized("x"), http.StatusUnauthorized},
		{Forbidden("x"), http.StatusForbidden},
		{Conflict("x"), http.StatusConflict},
		{InternalError("x"), http.StatusInternalServerError},
		{ValidationError(nil), http.StatusBadRequest},
	}
	for _, tt := range tests {
		rec := httptest.NewRecorder()
		WriteErr(rec, tt.resp)
		if rec.Code != tt.status {
			t.Errorf("code %s: expected %d, got %d", tt.resp.Error.Code, tt.status, rec.Code)
		}
	}
}
