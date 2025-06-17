package mdtable

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/yuin/goldmark"
	gmast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

func ExtractTables(markdown string) ([]Table, error) {
	var tables []Table

	mdParser := goldmark.New(
		goldmark.WithExtensions(extension.Table),
	)
	source := []byte(markdown) // <-- define source here

	reader := text.NewReader(source)
	doc := mdParser.Parser().Parse(reader)

	err := gmast.Walk(doc, func(n gmast.Node, entering bool) (gmast.WalkStatus, error) {
		if tbl, ok := n.(*extast.Table); ok && entering {
			var header []string
			var rows [][]string
			isHeader := true

			for row := tbl.FirstChild(); row != nil; row = row.NextSibling() {
				var rowData []string
				for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
					buf := bytes.Buffer{}
					err := gmast.Walk(cell, func(n gmast.Node, entering bool) (gmast.WalkStatus, error) {
						if entering {
							if tn, ok := n.(*gmast.Text); ok {
								buf.Write(tn.Segment.Value(source))
							}
						}
						return gmast.WalkContinue, nil
					})
					if err != nil {
						return gmast.WalkStop, fmt.Errorf("walking cell: %w", err)
					}
                    clean := sanitizeCellContent(buf.String())
                    rowData = append(rowData, clean)
				}
				if isHeader {
					header = rowData
					isHeader = false
				} else {
					rows = append(rows, rowData)
				}
			}

			tables = append(tables, Table{Header: header, Rows: rows})
		}
		return gmast.WalkContinue, nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking document: %w", err)
	}

	return tables, nil
}

var (
	mathExprRE     = regexp.MustCompile(`\$(.+?)\$`)
	latexCmdRE     = regexp.MustCompile(`\\[a-zA-Z]+\{([^{}]+)\}`)
	whitespaceRE   = regexp.MustCompile(`\s+`)
	backslashRE    = regexp.MustCompile(`\\`)
)

func sanitizeCellContent(s string) string {
	// Replace all $...$ spans
	return mathExprRE.ReplaceAllStringFunc(s, func(match string) string {
		content := match[1 : len(match)-1] // remove surrounding $
		
		// Unwrap LaTeX commands like \textbf{...}, \mathrm{...}, \mathbf{...}
		for latexCmdRE.MatchString(content) {
			content = latexCmdRE.ReplaceAllString(content, `$1`)
		}

		// Remove any remaining backslashes
		content = backslashRE.ReplaceAllString(content, "")

		// Remove all whitespace
		content = whitespaceRE.ReplaceAllString(content, "")

		return content
	})
}