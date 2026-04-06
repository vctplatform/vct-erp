// Package export provides CSV and JSON data export with streaming writers,
// configurable field selection, header mapping, and HTTP download handlers.
package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// Format
// ═══════════════════════════════════════════════════════════════

// Format defines the export output format.
type Format string

const (
	FormatCSV  Format = "csv"
	FormatJSON Format = "json"
)

// ParseFormat parses a format string.
func ParseFormat(s string) Format {
	switch strings.ToLower(s) {
	case "json":
		return FormatJSON
	default:
		return FormatCSV
	}
}

// ═══════════════════════════════════════════════════════════════
// Column Definition
// ═══════════════════════════════════════════════════════════════

// Column defines a field to export.
type Column struct {
	Key    string // Internal field key
	Header string // Display header for CSV
}

// ═══════════════════════════════════════════════════════════════
// Exporter
// ═══════════════════════════════════════════════════════════════

// Exporter streams rows to a writer in the specified format.
type Exporter struct {
	format  Format
	columns []Column
	w       io.Writer
	csvW    *csv.Writer
	jsonEnc *json.Encoder
	count   int
	started bool
}

// NewExporter creates a streaming exporter.
func NewExporter(w io.Writer, format Format, columns []Column) *Exporter {
	e := &Exporter{
		format:  format,
		columns: columns,
		w:       w,
	}

	switch format {
	case FormatCSV:
		e.csvW = csv.NewWriter(w)
	case FormatJSON:
		e.jsonEnc = json.NewEncoder(w)
	}

	return e
}

// WriteHeader writes the header row (CSV only).
func (e *Exporter) WriteHeader() error {
	if e.format != FormatCSV {
		return nil
	}
	headers := make([]string, len(e.columns))
	for i, col := range e.columns {
		headers[i] = col.Header
	}
	return e.csvW.Write(headers)
}

// WriteRow writes a single row of data.
// row maps column keys to string values.
func (e *Exporter) WriteRow(row map[string]string) error {
	e.count++

	switch e.format {
	case FormatCSV:
		record := make([]string, len(e.columns))
		for i, col := range e.columns {
			record[i] = row[col.Key]
		}
		return e.csvW.Write(record)

	case FormatJSON:
		// Write only selected columns
		filtered := make(map[string]string, len(e.columns))
		for _, col := range e.columns {
			filtered[col.Key] = row[col.Key]
		}
		return e.jsonEnc.Encode(filtered)
	}
	return nil
}

// Flush ensures all buffered data is written.
func (e *Exporter) Flush() {
	if e.csvW != nil {
		e.csvW.Flush()
	}
}

// Count returns the number of rows written.
func (e *Exporter) Count() int {
	return e.count
}

// ═══════════════════════════════════════════════════════════════
// HTTP Download Handler
// ═══════════════════════════════════════════════════════════════

// RowProvider is a function that yields rows to the exporter.
// Return io.EOF to signal end of data.
type RowProvider func() (map[string]string, error)

// DownloadHandler creates an HTTP handler that streams export data.
func DownloadHandler(filename string, format Format, columns []Column, provider RowProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Override format from query param
		if f := r.URL.Query().Get("format"); f != "" {
			format = ParseFormat(f)
		}

		ext := ".csv"
		contentType := "text/csv; charset=utf-8"
		if format == FormatJSON {
			ext = ".jsonl"
			contentType = "application/x-ndjson"
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition",
			fmt.Sprintf(`attachment; filename="%s_%s%s"`, filename, time.Now().Format("20060102_150405"), ext))
		w.Header().Set("Cache-Control", "no-cache")

		exp := NewExporter(w, format, columns)
		exp.WriteHeader()

		for {
			row, err := provider()
			if err == io.EOF {
				break
			}
			if err != nil {
				// Best-effort error in stream
				return
			}
			if writeErr := exp.WriteRow(row); writeErr != nil {
				return
			}
		}

		exp.Flush()
	})
}
