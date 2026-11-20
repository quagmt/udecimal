// Command report turns the raw `go test -bench` output in benchmarks/bench-*.txt and the Example functions into
// the tables, charts and code snippets the documentation site includes. Run it from the benchmarks directory:
//
//	go run ./cmd/report
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// libs is the display order. Each name maps to the file the benchmarks Makefile writes.
var libs = []string{"udecimal", "shopspring", "ericlagergren", "apd", "govalues"}

var files = map[string]string{
	"udecimal":      "bench-udec.txt",
	"shopspring":    "bench-ss.txt",
	"ericlagergren": "bench-ed.txt",
	"apd":           "bench-apd.txt",
	"govalues":      "bench-gov.txt",
}

type result struct {
	NsPerOp     float64
	AllocsPerOp float64
}

// suite holds the medians of every run, keyed by benchmark, case and library.
type suite struct {
	Benches  []string                                  // benchmark names in first-seen order, e.g. "Add"
	Cases    map[string][]string                       // bench -> case names in first-seen order
	Results  map[string]map[string]map[string]result   // bench -> case -> lib -> median
	Skips    map[string]map[string]map[string]string   // bench -> case -> lib -> "cannot represent: ..."
	RunCount int                                       // most runs seen for any one result
	runs     map[string]map[string]map[string][]result // bench -> case -> lib -> every run
}

// e.g. "BenchmarkAdd/3.Add(7)-10   100   10.0 ns/op   0 B/op   0 allocs/op"
var resultLine = regexp.MustCompile(`^(?:Benchmark)?([^/\s]+)/(\S+)-\d+\s+\d+\s+([\d.]+) ns/op\s+[\d.]+ B/op\s+([\d.]+) allocs/op`)

// e.g. "--- SKIP: BenchmarkAddWide/9.9.Add(1)"
var skipLine = regexp.MustCompile(`^--- SKIP: (?:Benchmark)?([^/\s]+)/(\S+)$`)

const skipPrefix = "cannot represent: "

func parseLine(line string) (bench, cas string, r result, ok bool) {
	m := resultLine.FindStringSubmatch(line)
	if m == nil {
		return "", "", result{}, false
	}

	ns, err1 := strconv.ParseFloat(m[3], 64)
	allocs, err2 := strconv.ParseFloat(m[4], 64)
	if err1 != nil || err2 != nil {
		return "", "", result{}, false
	}

	return m[1], m[2], result{NsPerOp: ns, AllocsPerOp: allocs}, true
}

// load reads one bench file per library (display name -> path).
func load(paths map[string]string) (*suite, error) {
	s := &suite{
		Cases:   map[string][]string{},
		Results: map[string]map[string]map[string]result{},
		Skips:   map[string]map[string]map[string]string{},
		runs:    map[string]map[string]map[string][]result{},
	}

	// read in display order so Benches and Cases come out the same on every run
	for _, lib := range libs {
		path, ok := paths[lib]
		if !ok {
			continue
		}

		if err := s.read(lib, path); err != nil {
			return nil, err
		}
	}

	for bench, cases := range s.runs {
		for cas, byLib := range cases {
			for lib, runs := range byLib {
				set(s.Results, bench, cas, lib, median(runs))
				s.RunCount = max(s.RunCount, len(runs))
			}
		}
	}

	return s, nil
}

func (s *suite) read(lib, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// with -v, the skip reason is printed on the line before "--- SKIP:"
	var reason string

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()

		if i := strings.Index(line, skipPrefix); i >= 0 {
			reason = line[i:]
			continue
		}

		if m := skipLine.FindStringSubmatch(line); m != nil {
			if reason == "" {
				return fmt.Errorf("%s: skip without a reason: %s", path, line)
			}

			s.see(m[1], m[2])
			set(s.Skips, m[1], m[2], lib, reason)
			reason = ""

			continue
		}

		if bench, cas, r, ok := parseLine(line); ok {
			s.see(bench, cas)

			if s.runs[bench] == nil {
				s.runs[bench] = map[string]map[string][]result{}
			}

			if s.runs[bench][cas] == nil {
				s.runs[bench][cas] = map[string][]result{}
			}

			s.runs[bench][cas][lib] = append(s.runs[bench][cas][lib], r)
		}
	}

	return sc.Err()
}

func (s *suite) see(bench, cas string) {
	if !slices.Contains(s.Benches, bench) {
		s.Benches = append(s.Benches, bench)
	}

	if !slices.Contains(s.Cases[bench], cas) {
		s.Cases[bench] = append(s.Cases[bench], cas)
	}
}

func set[V any](m map[string]map[string]map[string]V, a, b, c string, v V) {
	if m[a] == nil {
		m[a] = map[string]map[string]V{}
	}

	if m[a][b] == nil {
		m[a][b] = map[string]V{}
	}

	m[a][b][c] = v
}

func median(runs []result) result {
	ns := make([]float64, len(runs))
	allocs := make([]float64, len(runs))

	for i, r := range runs {
		ns[i], allocs[i] = r.NsPerOp, r.AllocsPerOp
	}

	return result{NsPerOp: medianOf(ns), AllocsPerOp: medianOf(allocs)}
}

func medianOf(v []float64) float64 {
	slices.Sort(v)

	n := len(v)
	if n%2 == 1 {
		return v[n/2]
	}

	return (v[n/2-1] + v[n/2]) / 2
}

func geomean(v []float64) float64 {
	var sum float64
	for _, x := range v {
		sum += math.Log(x)
	}

	return math.Exp(sum / float64(len(v)))
}

// number matches one decimal operand in a case name, e.g. "1234.5678" in "1234.5678.Add(1.5)".
var number = regexp.MustCompile(`\d+(?:\.\d+)?`)

// formatNs prints one decimal below 100 and none from 100 up, e.g. 4.44 -> "4.4", 99.97 -> "100", 123.4 -> "123".
func formatNs(ns float64) string {
	if ns < 99.95 {
		return strconv.FormatFloat(ns, 'f', 1, 64)
	}

	return strconv.FormatFloat(ns, 'f', 0, 64)
}

// sources is a repeatable -examples flag: name=path, e.g. udecimal=../doc_test.go. Default: udecimal and migration.
type sources map[string]string

func (s sources) String() string { return fmt.Sprint(map[string]string(s)) }

func (s sources) Set(v string) error {
	name, path, ok := strings.Cut(v, "=")
	if !ok {
		return fmt.Errorf("want name=path, got %q", v)
	}

	s[name] = path

	return nil
}

func main() {
	benchDir := flag.String("bench", ".", "directory holding the bench-*.txt files")
	out := flag.String("out", "../docs/snippets", "directory to write the snippets to; it is replaced")
	examples := sources{}
	flag.Var(examples, "examples", "name=path of a Go file whose Example functions to render; repeatable")
	flag.Parse()

	if len(examples) == 0 {
		examples["udecimal"] = "../doc_test.go"
		examples["migration"] = "migration_test.go"
	}

	if err := regenerate(*benchDir, *out, examples); err != nil {
		log.Fatal(err)
	}
}

// regenerate replaces out with freshly generated snippets, so renamed or removed snippets don't linger. It writes
// into a sibling directory first and only then swaps it in, so a failed run leaves out as it was. out must be a
// directory named snippets: it is deleted, and a typo must not delete something else.
func regenerate(benchDir, out string, examples map[string]string) error {
	out = filepath.Clean(out)
	if filepath.Base(out) != "snippets" {
		return fmt.Errorf("the output must be a directory named snippets, got %q", out)
	}

	tmp := out + ".new"
	if err := os.RemoveAll(tmp); err != nil {
		return err
	}

	if err := generate(benchDir, tmp, examples); err != nil {
		_ = os.RemoveAll(tmp)
		return err
	}

	if err := os.RemoveAll(out); err != nil {
		return err
	}

	return os.Rename(tmp, out)
}
