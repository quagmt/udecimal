package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var fixtures = map[string]string{
	"udecimal": "testdata/bench-udec.txt",
	"govalues": "testdata/bench-gov.txt",
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		line   string
		bench  string
		cas    string
		ns     float64
		allocs float64
		ok     bool
	}{
		{"BenchmarkAdd/123.456.Add(0.123)-10   100000000   10.75 ns/op   0 B/op   0 allocs/op", "Add", "123.456.Add(0.123)", 10.75, 0, true},
		{"BenchmarkDivWide/1234.Div(5)-10 \t 3901846 \t 306.7 ns/op \t 320 B/op \t 7 allocs/op", "DivWide", "1234.Div(5)", 306.7, 7, true},
		{"Add/3.Add(7)-10   100   10.0 ns/op   0 B/op   0 allocs/op", "Add", "3.Add(7)", 10, 0, true},
		{"BenchmarkAdd/3.Add(7)", "", "", 0, 0, false},
		{"--- SKIP: BenchmarkAddWide/9.9.Add(1)", "", "", 0, 0, false},
		{"goos: darwin", "", "", 0, 0, false},
		{"PASS", "", "", 0, 0, false},
		{"", "", "", 0, 0, false},
	}

	for _, tc := range tests {
		bench, cas, r, ok := parseLine(tc.line)
		require.Equal(t, tc.ok, ok, tc.line)
		require.Equal(t, tc.bench, bench, tc.line)
		require.Equal(t, tc.cas, cas, tc.line)
		require.Equal(t, tc.ns, r.NsPerOp, tc.line)
		require.Equal(t, tc.allocs, r.AllocsPerOp, tc.line)
	}
}

func TestLoadMedian(t *testing.T) {
	s, err := load(fixtures)
	require.NoError(t, err)

	// runs of 10 and 12 ns/op
	require.Equal(t, 11.0, s.Results["Add"]["1.Add(2)"]["udecimal"].NsPerOp)
	require.Equal(t, 55.0, s.Results["Add"]["1.Add(2)"]["govalues"].NsPerOp)
	require.Equal(t, 1.0, s.Results["Add"]["1.Add(2)"]["govalues"].AllocsPerOp)
	require.Equal(t, []string{"Add", "AddWide"}, s.Benches)
	require.Equal(t, []string{"1.Add(2)"}, s.Cases["Add"])
}

func TestLoadSkipReason(t *testing.T) {
	s, err := load(fixtures)
	require.NoError(t, err)

	require.Equal(t, "cannot represent: parsed 9.9 as 9", s.Skips["AddWide"]["9.9.Add(1)"]["govalues"])
	_, ok := s.Results["AddWide"]["9.9.Add(1)"]["govalues"]
	require.False(t, ok)
}

func TestLoadMissingFile(t *testing.T) {
	_, err := load(map[string]string{"udecimal": "testdata/nope.txt"})
	require.Error(t, err)
}

func TestRunCount(t *testing.T) {
	s, err := load(fixtures)
	require.NoError(t, err)
	require.Equal(t, 2, s.RunCount)
}

func TestGeomean(t *testing.T) {
	tests := []struct {
		in   []float64
		want float64
	}{
		{[]float64{4}, 4},
		{[]float64{2, 8}, 4},
		{[]float64{1, 10, 100}, 10},
	}

	for _, tc := range tests {
		require.InDelta(t, tc.want, geomean(tc.in), 1e-9)
	}
}

func TestClassify(t *testing.T) {
	tests := []struct {
		bench, cas string
		want       class
	}{
		{"Add", "3.Add(7)", classHead},
		{"Parse", "0.1234567890123456879", classHead},
		{"ParseWide", "1234567890123456789.1234567890123456879", classWide128},                          // 38 digits
		{"DivWide", "22773757910726981402256170801141121114.Div(811656739243220271.159)", classWide128}, // 38 digits
		{"ParseWide", "123456789123456789123456.1234567890123456", classWideBig},                        // 40 digits
		{"DivWide", "1234.Div(12345679012345679890123456789.1234567890123456789)", classWideBig},        // 48 digits
		{"ParseWide", "0.00000000000000000000000000000000000000001", classWide128},                      // leading zeros don't count
	}

	for _, tc := range tests {
		require.Equal(t, tc.want, classify(tc.bench, tc.cas), tc.cas)
	}
}

func TestCheckClasses(t *testing.T) {
	// a 20-to-38-digit AddWide case where udecimal allocates is an error; StringWide is exempt
	s := &suite{
		Benches: []string{"AddWide", "StringWide"},
		Cases: map[string][]string{
			"AddWide":    {"1234567890123456789.1234567890123456879.Add(1111.1789)"},
			"StringWide": {"12345.1234567890123456789"},
		},
		Results: map[string]map[string]map[string]result{
			"AddWide":    {"1234567890123456789.1234567890123456879.Add(1111.1789)": {"udecimal": {NsPerOp: 5, AllocsPerOp: 1}}},
			"StringWide": {"12345.1234567890123456789": {"udecimal": {NsPerOp: 29, AllocsPerOp: 1}}},
		},
	}
	require.ErrorContains(t, checkClasses(s), "AddWide")

	s.Results["AddWide"]["1234567890123456789.1234567890123456879.Add(1111.1789)"]["udecimal"] = result{NsPerOp: 5}
	require.NoError(t, checkClasses(s))
}

func TestFormatNs(t *testing.T) {
	tests := []struct {
		in   float64
		want string
	}{
		{4.44, "4.4"},
		{99.94, "99.9"},
		{99.97, "100"}, // rounds to 100.0 at one decimal, so it is printed like any other value of 100 or more
		{123.4, "123"},
	}

	for _, tc := range tests {
		require.Equal(t, tc.want, formatNs(tc.in), tc.in)
	}
}

func TestRegenerateKeepsOldOnError(t *testing.T) {
	out := filepath.Join(t.TempDir(), "snippets")
	writeFile(t, out, "old.md", "kept")

	// no bench files in an empty directory: generation fails
	require.Error(t, regenerate(t.TempDir(), out, nil))

	content, err := os.ReadFile(filepath.Join(out, "old.md"))
	require.NoError(t, err)
	require.Equal(t, "kept", string(content))

	_, err = os.Stat(out + ".new")
	require.True(t, os.IsNotExist(err), "the temporary directory is left behind")
}

func TestRegenerateRefusesOtherDirs(t *testing.T) {
	require.ErrorContains(t, regenerate(".", t.TempDir(), nil), "snippets")
}
