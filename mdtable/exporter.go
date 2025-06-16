package mdtable

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"

    "github.com/xuri/excelize/v2"
)

func ExportToCSV(tables []Table, outputDir string) error {
    if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
        return fmt.Errorf("creating output dir: %w", err)
    }

    for i, table := range tables {
        filename := filepath.Join(outputDir, fmt.Sprintf("table_%d.csv", i+1))
        file, err := os.Create(filename)
        if err != nil {
            return fmt.Errorf("creating CSV file %s: %w", filename, err)
        }

        writer := csv.NewWriter(file)

        if err := writer.Write(table.Header); err != nil {
            file.Close()
            return fmt.Errorf("writing header to %s: %w", filename, err)
        }
        for _, row := range table.Rows {
            if err := writer.Write(row); err != nil {
                file.Close()
                return fmt.Errorf("writing row to %s: %w", filename, err)
            }
        }
        writer.Flush()
        if err := writer.Error(); err != nil {
            file.Close()
            return fmt.Errorf("flushing CSV writer for %s: %w", filename, err)
        }

        if err := file.Close(); err != nil {
            return fmt.Errorf("closing CSV file %s: %w", filename, err)
        }
    }
    return nil
}

func ExportToJSON(tables []Table, outputDir string) error {
    if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
        return fmt.Errorf("creating output dir: %w", err)
    }

    filename := filepath.Join(outputDir, "tables.json")
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("creating JSON file %s: %w", filename, err)
    }
    defer func() {
        if cerr := file.Close(); cerr != nil {
            err = fmt.Errorf("closing JSON file %s: %w", filename, cerr)
        }
    }()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(tables); err != nil {
        return fmt.Errorf("encoding JSON to %s: %w", filename, err)
    }

    return nil
}

func ExportToExcel(tables []Table, outputDir string) error {
    if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
        return fmt.Errorf("creating output dir: %w", err)
    }

    f := excelize.NewFile()
    for i, table := range tables {
        sheetName := fmt.Sprintf("Table%d", i+1)
        if i == 0 {
            if err := f.SetSheetName("Sheet1", sheetName); err != nil {
                return fmt.Errorf("setting sheet name: %w", err)
            }
        } else {
            if _, err := f.NewSheet(sheetName); err != nil {
                return fmt.Errorf("creating new sheet: %w", err)
            }
        }

        for j, col := range table.Header {
            cell, _ := excelize.CoordinatesToCellName(j+1, 1)
            if err := f.SetCellValue(sheetName, cell, col); err != nil {
                return fmt.Errorf("setting cell value: %w", err)
            }
        }

        for rowIdx, row := range table.Rows {
            for colIdx, cellVal := range row {
                cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
                if err := f.SetCellValue(sheetName, cell, cellVal); err != nil {
                    return fmt.Errorf("setting cell value: %w", err)
                }
            }
        }
    }

    filename := filepath.Join(outputDir, "tables.xlsx")
    if err := f.SaveAs(filename); err != nil {
        return fmt.Errorf("saving Excel file %s: %w", filename, err)
    }

    return nil
}