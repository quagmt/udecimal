package main

import (
	"bytes"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

// renderExamples returns every Example function in the Go file at path as a fenced Markdown code block, keyed by
// the function's name, e.g. "ExampleDecimal_Add64". The block is the function body, gofmt-formatted and dedented,
// which ends with the expected output comment `go test` checks:
//
//	a := udecimal.MustParse("1.23")
//	fmt.Println(a.Add64(4))
//	// Output:
//	// 5.23
//
// The site shows exactly what `go test` checks, so a snippet can't drift from the code.
func renderExamples(path string) (map[string]string, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	out := map[string]string{}
	for _, ex := range doc.Examples(f) {
		body, err := printBody(fset, ex.Code, ex.Comments)
		if err != nil {
			return nil, err
		}

		var sb strings.Builder
		sb.WriteString("```go\n")
		sb.WriteString(body)

		sb.WriteString("```\n")
		out["Example"+ex.Name] = sb.String()
	}

	return out, nil
}

// printBody formats a function body without its braces and one level of indentation, e.g. "{\n\ta := 1\n}" -> "a := 1\n".
func printBody(fset *token.FileSet, code ast.Node, comments []*ast.CommentGroup) (string, error) {
	var buf bytes.Buffer

	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}
	if err := cfg.Fprint(&buf, fset, &printer.CommentedNode{Node: code, Comments: comments}); err != nil {
		return "", err
	}

	lines := strings.Split(buf.String(), "\n")
	lines = lines[1 : len(lines)-1] // drop "{" and "}"

	var sb strings.Builder
	for _, line := range lines {
		sb.WriteString(strings.TrimPrefix(line, "\t") + "\n")
	}

	return sb.String(), nil
}
