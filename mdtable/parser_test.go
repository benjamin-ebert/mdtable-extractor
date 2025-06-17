package mdtable

import (
	"testing"
)

const sampleMarkdown = `
| Name  | Age | City     |
|-------|-----|----------|
| Alice | 30  | New York |
| Bob   | 25  | London   |
`

func TestExtractTables(t *testing.T) {
	tables, err := ExtractTables(sampleMarkdown)
	if err != nil {
		t.Fatalf("ExtractTables returned error: %v", err)
	}
	if len(tables) != 1 {
		t.Fatalf("Expected 1 table, got %d", len(tables))
	}

	table := tables[0]
	if len(table.Header) != 3 {
		t.Errorf("Expected header length 3, got %d", len(table.Header))
	}
	if len(table.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(table.Rows))
	}
	if table.Header[0] != "Name" || table.Header[1] != "Age" || table.Header[2] != "City" {
		t.Errorf("Unexpected header values: %v", table.Header)
	}
	if table.Rows[0][0] != "Alice" || table.Rows[1][0] != "Bob" {
		t.Errorf("Unexpected row values: %v", table.Rows)
	}
}

func TestSanitizeCellContent(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Math-mode input
		{"$44.4 \\%$", "44.4%"},
		{"$\\mathbf{8 0 . 7 \\%}$", "80.7%"},
		{"$\\textbf{6 . 8 4}$ +/- $\\textbf{0 . 0 7}$", "6.84 +/- 0.07"},
		{"$\\mathrm{1 2 3}$", "123"},
		{"$\\mathbf{2 9 . 0 \\%}$", "29.0%"},

		// Literal dollar input — should remain unchanged
		{"$5.00", "$5.00"},
		{"Price is $5.00", "Price is $5.00"},
		{"\\$5.00", "\\$5.00"},

		// No formatting
		{"Plain text", "Plain text"},
		{"123.45%", "123.45%"},
	}

	for _, tt := range tests {
		result := sanitizeCellContent(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeCellContent(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}