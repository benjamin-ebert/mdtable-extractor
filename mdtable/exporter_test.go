package mdtable

import (
    "os"
    "path/filepath"
    "testing"
)

func TestExporters(t *testing.T) {
    tables := []Table{
        {
            Header: []string{"Col1", "Col2"},
            Rows: [][]string{
                {"val11", "val12"},
                {"val21", "val22"},
            },
        },
    }

    tempDir := t.TempDir()

    if err := ExportToCSV(tables, tempDir); err != nil {
        t.Fatalf("ExportToCSV failed: %v", err)
    }
    if err := ExportToJSON(tables, tempDir); err != nil {
        t.Fatalf("ExportToJSON failed: %v", err)
    }
    if err := ExportToExcel(tables, tempDir); err != nil {
        t.Fatalf("ExportToExcel failed: %v", err)
    }

    // Check files exist
    csvPath := filepath.Join(tempDir, "table_1.csv")
    if _, err := os.Stat(csvPath); err != nil {
        t.Errorf("CSV file not found: %v", err)
    }
    jsonPath := filepath.Join(tempDir, "tables.json")
    if _, err := os.Stat(jsonPath); err != nil {
        t.Errorf("JSON file not found: %v", err)
    }
    excelPath := filepath.Join(tempDir, "tables.xlsx")
    if _, err := os.Stat(excelPath); err != nil {
        t.Errorf("Excel file not found: %v", err)
    }
}