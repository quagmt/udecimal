package benchmarks

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/quagmt/udecimal"
	"github.com/stretchr/testify/require"
)

// udecimalDoc parses the root package's non-test files.
func udecimalDoc(t *testing.T) *doc.Package {
	t.Helper()

	paths, err := filepath.Glob("../*.go")
	require.NoError(t, err)

	fset := token.NewFileSet()
	var files []*ast.File

	for _, path := range paths {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}

		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		require.NoError(t, err)

		files = append(files, f)
	}

	pkg, err := doc.NewFromFiles(fset, files, "github.com/quagmt/udecimal")
	require.NoError(t, err)

	return pkg
}

// exported returns every exported identifier of the package whose doc comment doesn't mark it deprecated.
func exported(pkg *doc.Package) []string {
	var names []string

	add := func(name, comment string) {
		if !strings.Contains(comment, "Deprecated:") {
			names = append(names, name)
		}
	}

	values := func(vs []*doc.Value) {
		for _, v := range vs {
			for _, n := range v.Names {
				if ast.IsExported(n) {
					add(n, v.Doc)
				}
			}
		}
	}

	values(pkg.Consts)
	values(pkg.Vars)

	for _, f := range pkg.Funcs {
		add(f.Name, f.Doc)
	}

	for _, typ := range pkg.Types {
		add(typ.Name, typ.Doc)
		values(typ.Consts)
		values(typ.Vars)

		for _, f := range typ.Funcs {
			add(f.Name, f.Doc)
		}

		for _, m := range typ.Methods {
			add(m.Name, m.Doc)
		}
	}

	slices.Sort(names)

	return slices.Compact(names)
}

// include matches a snippet include on a page, e.g. `--8<-- "examples/udecimal/Example.md"`.
var include = regexp.MustCompile(`--8<-- "([^"]+)"`)

// pages returns every site page, by path, with the snippet files it includes appended to its text.
func pages(t *testing.T) map[string]string {
	t.Helper()

	paths, err := filepath.Glob("../docs/*.md")
	require.NoError(t, err)

	out := map[string]string{}
	for _, path := range paths {
		b, err := os.ReadFile(path)
		require.NoError(t, err)

		text := string(b)
		for _, m := range include.FindAllStringSubmatch(text, -1) {
			s, err := os.ReadFile(filepath.Join("../docs/snippets", m[1]))
			require.NoError(t, err, m[1])

			text += "\n" + string(s)
		}

		out[path] = text
	}

	return out
}

func TestExportedIdentifiersDocumented(t *testing.T) {
	var all strings.Builder
	for _, text := range pages(t) {
		all.WriteString(text)
	}

	var missing []string
	for _, name := range exported(udecimalDoc(t)) {
		if !regexp.MustCompile(`\b` + name + `\b`).MatchString(all.String()) {
			missing = append(missing, name)
		}
	}

	require.Empty(t, missing, "exported identifiers no page mentions")
}

// goBlock matches a hand-written Go code block on a page.
var goBlock = regexp.MustCompile("(?s)```go\n(.*?)```")

func TestCodeBlocks(t *testing.T) {
	pkg := udecimalDoc(t)

	// package-level functions, e.g. Parse, MustParse, NewFromInt64, Max
	var funcs []string
	for _, f := range pkg.Funcs {
		funcs = append(funcs, f.Name)
	}

	for _, typ := range pkg.Types {
		for _, f := range typ.Funcs {
			funcs = append(funcs, f.Name)
		}
	}

	unqualified := regexp.MustCompile(`(^|[^.\w])(` + strings.Join(funcs, "|") + `)\(`)

	paths, err := filepath.Glob("../docs/*.md")
	require.NoError(t, err)

	for _, path := range paths {
		page, err := os.ReadFile(path)
		require.NoError(t, err)

		for _, m := range goBlock.FindAllStringSubmatch(string(page), -1) {
			block := m[1]

			require.NotContains(t, block, "fmt.Print", "%s: a block that prints output must be an included Example:\n%s", path, block)
			require.False(t, unqualified.MatchString(block), "%s: write package functions as udecimal.X:\n%s", path, block)
		}
	}
}

// TestBigResultsFitAgain pins what How it works says: a result computed with math/big is stored in 128 bits again
// if it fits, so ToHiLo and Int64 work on it and later operations don't allocate.
func TestBigResultsFitAgain(t *testing.T) {
	a := udecimal.MustParse("1" + strings.Repeat("0", 39)) // 10^39: more than 2^128 - 1
	b := udecimal.MustParse(strings.Repeat("9", 39))       // 10^39 - 1: more than 2^128 - 1
	one := a.Sub(b)
	require.Equal(t, "1", one.String())

	_, _, lo, _, ok := one.ToHiLo()
	require.True(t, ok)
	require.Equal(t, uint64(1), lo)

	n, err := one.Int64()
	require.NoError(t, err)
	require.Equal(t, int64(1), n)

	two := udecimal.MustParse("2")
	require.Zero(t, testing.AllocsPerRun(10, func() { _ = one.Add(two) }))
}
