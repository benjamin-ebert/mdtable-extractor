package mdtable

import (
    "bytes"
    "strings"

    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/extension"
    "github.com/yuin/goldmark/text"
    gmast "github.com/yuin/goldmark/ast"
    extast "github.com/yuin/goldmark/extension/ast"
)

func ExtractTables(markdown string) ([]Table, error) {
    var tables []Table

    mdParser := goldmark.New(
        goldmark.WithExtensions(extension.Table),
    )
    source := []byte(markdown)  // <-- define source here

    reader := text.NewReader(source)
    doc := mdParser.Parser().Parse(reader)

    gmast.Walk(doc, func(n gmast.Node, entering bool) (gmast.WalkStatus, error) {
        if tbl, ok := n.(*extast.Table); ok && entering {
            var header []string
            var rows [][]string
            isHeader := true

            for row := tbl.FirstChild(); row != nil; row = row.NextSibling() {
                var rowData []string
                for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
                    buf := bytes.Buffer{}
                    gmast.Walk(cell, func(n gmast.Node, entering bool) (gmast.WalkStatus, error) {
                        if entering {
                            if tn, ok := n.(*gmast.Text); ok {
                                buf.Write(tn.Segment.Value(source))
                            }
                        }
                        return gmast.WalkContinue, nil
                    })
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

    return tables, nil
}