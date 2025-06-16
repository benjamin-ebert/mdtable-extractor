package mdtable

import (
    "encoding/csv"
    "encoding/json"
    "os"
    "path/filepath"
    "strings"
    "testing"

    "github.com/xuri/excelize/v2"
)

const integrationMarkdown = `
| Product | Price | Stock |
|---------|-------|-------|
| Pen     | 1.20  | 100   |
| Paper   | 0.50  | 500   |
`

func TestIntegration_ExtractAndExport(t *testing.T) {
    tables, err := ExtractTables(integrationMarkdown)
    if err != nil {
        t.Fatalf("ExtractTables error: %v", err)
    }

    if len(tables) != 1 {
        t.Fatalf("Expected 1 table, got %d", len(tables))
    }

    outputDir := t.TempDir()

    // Export all formats
    if err := ExportToCSV(tables, outputDir); err != nil {
        t.Fatalf("ExportToCSV error: %v", err)
    }
    if err := ExportToJSON(tables, outputDir); err != nil {
        t.Fatalf("ExportToJSON error: %v", err)
    }
    if err := ExportToExcel(tables, outputDir); err != nil {
        t.Fatalf("ExportToExcel error: %v", err)
    }

    // --- Validate CSV ---
    csvPath := filepath.Join(outputDir, "table_1.csv")
    csvFile, err := os.Open(csvPath)
    if err != nil {
        t.Fatalf("Opening CSV failed: %v", err)
    }
    defer csvFile.Close()

    reader := csv.NewReader(csvFile)
    records, err := reader.ReadAll()
    if err != nil {
        t.Fatalf("Reading CSV failed: %v", err)
    }

    expectedHeader := []string{"Product", "Price", "Stock"}
    if len(records) < 3 {
        t.Fatalf("CSV records too few: got %d", len(records))
    }
    for i, val := range expectedHeader {
        if records[0][i] != val {
            t.Errorf("CSV header mismatch at col %d: got %q, want %q", i, records[0][i], val)
        }
    }

    expectedRow1 := []string{"Pen", "1.20", "100"}
    for i, val := range expectedRow1 {
        if records[1][i] != val {
            t.Errorf("CSV row 1 col %d mismatch: got %q, want %q", i, records[1][i], val)
        }
    }

    expectedRow2 := []string{"Paper", "0.50", "500"}
    for i, val := range expectedRow2 {
        if records[2][i] != val {
            t.Errorf("CSV row 2 col %d mismatch: got %q, want %q", i, records[2][i], val)
        }
    }

    // --- Validate JSON ---
    jsonPath := filepath.Join(outputDir, "tables.json")
    jsonData, err := os.ReadFile(jsonPath)
    if err != nil {
        t.Fatalf("Reading JSON failed: %v", err)
    }

    var loadedTables []Table
    if err := json.Unmarshal(jsonData, &loadedTables); err != nil {
        t.Fatalf("Unmarshal JSON failed: %v", err)
    }

    if len(loadedTables) != 1 {
        t.Fatalf("JSON table count mismatch: got %d", len(loadedTables))
    }

    if strings.TrimSpace(loadedTables[0].Header[0]) != "Product" {
        t.Errorf("JSON header mismatch: got %v", loadedTables[0].Header)
    }

    // --- Validate Excel content ---
    excelPath := filepath.Join(outputDir, "tables.xlsx")
    if _, err := os.Stat(excelPath); err != nil {
        t.Fatalf("Excel file not found: %v", err)
    }

    f, err := excelize.OpenFile(excelPath)
    if err != nil {
        t.Fatalf("Open Excel failed: %v", err)
    }
    defer f.Close()

    sheetName := "Table1"
    rows, err := f.GetRows(sheetName)
    if err != nil {
        t.Fatalf("GetRows failed: %v", err)
    }

    if len(rows) != 3 {
        t.Fatalf("Expected 3 rows in Excel, got %d", len(rows))
    }

    for i, val := range expectedHeader {
        if rows[0][i] != val {
            t.Errorf("Excel header col %d mismatch: got %q, want %q", i, rows[0][i], val)
        }
    }

    for i, val := range expectedRow1 {
        if rows[1][i] != val {
            t.Errorf("Excel row 1 col %d mismatch: got %q, want %q", i, rows[1][i], val)
        }
    }

    for i, val := range expectedRow2 {
        if rows[2][i] != val {
            t.Errorf("Excel row 2 col %d mismatch: got %q, want %q", i, rows[2][i], val)
        }
    }
}