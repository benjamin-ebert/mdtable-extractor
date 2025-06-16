
# mdtable-extractor

**mdtable-extractor** is a Go library and CLI tool to extract tables from Markdown files (usually converted from PDFs) and export them as CSV, JSON, or Excel files.

## Features

- Parses Markdown tables using [goldmark](https://github.com/yuin/goldmark) parser  
- Supports exporting extracted tables to CSV, JSON, and Excel formats  
- Easy to use CLI and reusable Go library  
- Designed for automation and integration in pipelines  

## Installation

Make sure you have Go 1.24+ installed.

Clone this repository:

```bash
git clone https://github.com/benjamin-ebert/mdtable-extractor.git
cd mdtable-extractor
```

Build the CLI tool:

```bash
go build -o mdtable-extractor ./main.go
```

Or install directly with:

```bash
go install github.com/benjamin-ebert/mdtable-extractor@latest
```

## Usage

### CLI

```bash
./mdtable-extractor -input sample.md -output ./output
```

- `-input`: Path to the input Markdown file containing tables  
- `-output`: Directory where extracted CSV, JSON, and Excel files will be saved  

### Go Library

Import the package in your project:

```go
import "github.com/benjamin-ebert/mdtable-extractor/mdtable"
```

Example usage:

```go
data, err := os.ReadFile("sample.md")
if err != nil {
    log.Fatal(err)
}

tables, err := mdtable.ExtractTables(string(data))
if err != nil {
    log.Fatal(err)
}

err = mdtable.ExportToCSV(tables, "./output")
if err != nil {
    log.Fatal(err)
}
```

## Development

### Run tests with coverage info

```bash
go test -v -cover ./...
```

### Run tests with full coverage file and HTML report

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Lint code

```bash
golangci-lint run
```

## Contributing

Contributions are welcome! Feel free to open issues or pull requests.

## License

MIT © 2025 Benjamin Ebert

## Acknowledgments

- [goldmark](https://github.com/yuin/goldmark) for Markdown parsing  
- [excelize](https://github.com/xuri/excelize) for Excel file generation
