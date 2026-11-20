package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderExamples(t *testing.T) {
	got, err := renderExamples("testdata/example_test.go.txt")
	require.NoError(t, err)
	require.Len(t, got, 2)

	require.Equal(t, "```go\n"+
		"a := 1\n"+
		"\n"+
		"// a comment inside the body stays\n"+
		"fmt.Println(a)\n"+
		"fmt.Println(a + 1)\n"+
		"// Output:\n"+
		"// 1\n"+
		"// 2\n"+
		"```\n", got["ExampleTwo"])

	require.Equal(t, "```go\nfmt.Println(\"x\")\n```\n", got["Example_noOutput"])
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

func TestDiffDirs(t *testing.T) {
	want, got := t.TempDir(), t.TempDir()
	writeFile(t, want, "same.md", "a")
	writeFile(t, got, "same.md", "a")
	writeFile(t, want, "sub/changed.md", "old")
	writeFile(t, got, "sub/changed.md", "new")
	writeFile(t, want, "missing.md", "x")
	writeFile(t, got, "extra.md", "y")

	diffs, err := diffDirs(want, got)
	require.NoError(t, err)
	require.Equal(t, []string{"extra.md", "missing.md", "sub/changed.md"}, diffs)

	diffs, err = diffDirs(want, want)
	require.NoError(t, err)
	require.Empty(t, diffs)
}

func TestSnippetsUpToDate(t *testing.T) {
	tmp := t.TempDir()
	require.NoError(t, generate("../..", tmp, map[string]string{"udecimal": "../../../doc_test.go", "migration": "../../migration_test.go"}))

	diffs, err := diffDirs("../../../docs/snippets", tmp)
	require.NoError(t, err)
	require.Empty(t, diffs, "docs/snippets is stale: run `go run ./cmd/report` in benchmarks/ and commit the result")
}
