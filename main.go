package main

import (
	"fmt"
	"os"

	"github.com/benjamin-ebert/mdtable-extractor/mdtable"
)

func main() {
	markdown, err := os.ReadFile("sample.md")
	if err != nil {
		panic(err)
	}

	tables, err := mdtable.ExtractTables(string(markdown))
	if err != nil {
		panic(err)
	}

	outputDir := "./output"

	if err := mdtable.ExportToCSV(tables, outputDir); err != nil {
		panic(err)
	}
	if err := mdtable.ExportToJSON(tables, outputDir); err != nil {
		panic(err)
	}
	if err := mdtable.ExportToExcel(tables, outputDir); err != nil {
		panic(err)
	}

	fmt.Println("✅ Tables exported to CSV, JSON, and Excel in ./output")
}
