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
