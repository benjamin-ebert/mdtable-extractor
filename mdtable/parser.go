package mdtable

import (
	"bytes"
	"fmt"
	"strings"

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
					rowData = append(rowData, strings.TrimSpace(buf.String()))
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
