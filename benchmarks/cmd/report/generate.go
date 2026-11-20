package main

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// chartGroup is one file of charts: a bench per chart, all over cases of one class.
type chartGroup struct {
	file    string
	benches []string
	class   class
}

var wideBenches = []string{"ParseWide", "StringWide", "AddWide", "SubWide", "MulWide", "DivWide"}

var chartGroups = []chartGroup{
	{"performance/arithmetic.html", []string{"Add", "Sub", "Mul", "Div"}, classHead},
	{"performance/parse-format.html", []string{"Parse", "String"}, classHead},
	{"performance/json.html", []string{"MarshalJSON", "UnmarshalJSON"}, classHead},
	{"performance/wide-128.html", wideBenches, classWide128},
}

// caseGroup is one collapsed "every case" table.
type caseGroup struct {
	file    string
	benches []string
	class   class
}

var caseGroups = []caseGroup{
	{"performance/cases-arithmetic.md", []string{"Add", "Sub", "Mul", "Div"}, classHead},
	{"performance/cases-parse-format.md", []string{"Parse", "String", "MarshalText", "UnmarshalText"}, classHead},
	{"performance/cases-json.md", []string{"MarshalJSON", "UnmarshalJSON"}, classHead},
	{"performance/pow.md", []string{"Pow"}, classHead},
	{"performance/cases-wide-128.md", wideBenches, classWide128},
	{"performance/wide-big.md", wideBenches, classWideBig},
}

// generate writes every snippet the docs include into outDir: tables and charts from the bench-*.txt files in
// benchDir, and the Example functions of each file in examples (source name -> path) under examples/<name>/.
func generate(benchDir, outDir string, examples map[string]string) error {
	paths := map[string]string{}
	for lib, file := range files {
		paths[lib] = filepath.Join(benchDir, file)
	}

	s, err := load(paths)
	if err != nil {
		return err
	}

	if err := checkClasses(s); err != nil {
		return err
	}

	out := map[string]string{}

	if out["performance/summary.md"], err = summaryTable(s, libs, headBenches); err != nil {
		return err
	}

	out["performance/allocs.md"] = allocTable(s, libs, headBenches)
	out["performance/govalues-reasons.md"] = reasonsList(s)

	for _, g := range chartGroups {
		var figures []string

		for _, bench := range g.benches {
			bars, err := chartBars(s, libs, bench, g.class)
			if err != nil {
				return err
			}

			figures = append(figures, renderChart(strings.TrimSuffix(bench, "Wide"), bars))
		}

		out[g.file] = renderGroup(figures)
	}

	for _, g := range caseGroups {
		if out[g.file], err = caseTable(s, libs, g.benches, g.class); err != nil {
			return err
		}
	}

	for name, path := range examples {
		rendered, err := renderExamples(path)
		if err != nil {
			return err
		}

		for fn, block := range rendered {
			out[filepath.Join("examples", name, fn+".md")] = block
		}
	}

	for rel, content := range out {
		path := filepath.Join(outDir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}

		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}

	return nil
}

// diffDirs returns the relative paths of files that are only in want, only in got, or differ between them, sorted.
func diffDirs(want, got string) ([]string, error) {
	a, err := readTree(want)
	if err != nil {
		return nil, err
	}

	b, err := readTree(got)
	if err != nil {
		return nil, err
	}

	var diffs []string
	for p, content := range a {
		if other, ok := b[p]; !ok || !bytes.Equal(content, other) {
			diffs = append(diffs, p)
		}
	}

	for p := range b {
		if _, ok := a[p]; !ok {
			diffs = append(diffs, p)
		}
	}

	slices.Sort(diffs)

	return diffs, nil
}

// readTree maps every regular file under root, by slash-separated relative path, to its content.
func readTree(root string) (map[string][]byte, error) {
	tree := map[string][]byte{}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		tree[filepath.ToSlash(rel)] = content

		return nil
	})

	return tree, err
}
