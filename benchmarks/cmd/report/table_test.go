package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// tableSuite: Add and Sub with one case each for udecimal, shopspring and govalues; govalues is faster on Sub.
func tableSuite() *suite {
	return &suite{
		Benches: []string{"Add", "Sub"},
		Cases:   map[string][]string{"Add": {"3.Add(7)"}, "Sub": {"3.Sub(7)"}},
		Results: map[string]map[string]map[string]result{
			"Add": {"3.Add(7)": {"udecimal": {NsPerOp: 4.4}, "shopspring": {NsPerOp: 23.76, AllocsPerOp: 2}, "govalues": {NsPerOp: 4.7}}},
			"Sub": {"3.Sub(7)": {"udecimal": {NsPerOp: 6.3}, "shopspring": {NsPerOp: 28.6, AllocsPerOp: 2}, "govalues": {NsPerOp: 5.481}}},
		},
		Skips: map[string]map[string]map[string]string{},
	}
}

var tableLibs = []string{"udecimal", "shopspring", "govalues"}

func TestSummaryTable(t *testing.T) {
	out, err := summaryTable(tableSuite(), tableLibs, []string{"Add", "Sub"})
	require.NoError(t, err)

	require.Contains(t, out, "| Benchmark | udecimal | shopspring | govalues |\n")
	require.Contains(t, out, "| Add | 4.4 ns | 5.4× | 1.1× |\n")
	require.Contains(t, out, "| Sub | 6.3 ns | 4.5× | **0.87×** |\n")
	require.Contains(t, out, "Values in bold are faster than udecimal.")
}

func TestSummaryTableMissingResult(t *testing.T) {
	s := tableSuite()
	delete(s.Results["Sub"]["3.Sub(7)"], "govalues")

	out, err := summaryTable(s, tableLibs, []string{"Add", "Sub"})
	require.ErrorContains(t, err, "Sub")
	require.ErrorContains(t, err, "govalues")
	require.NotContains(t, out, "NaN")
	require.NotContains(t, out, "Inf")
}

func TestAllocTable(t *testing.T) {
	s := tableSuite()
	s.Cases["Add"] = append(s.Cases["Add"], "1.Add(2)")
	s.Results["Add"]["1.Add(2)"] = map[string]result{"udecimal": {NsPerOp: 5}, "shopspring": {NsPerOp: 20, AllocsPerOp: 6}, "govalues": {NsPerOp: 5}}

	out := allocTable(s, tableLibs, []string{"Add", "Sub"})

	require.Contains(t, out, "| Add | 0 | 2–6 | 0 |\n")
	require.Contains(t, out, "| Sub | 0 | 2 | 0 |\n")
}

func TestCaseTable(t *testing.T) {
	s := &suite{
		Cases: map[string][]string{"AddWide": {"9.9.Add(1)", "1234.Div(12345679012345679890123456789.1234567890123456789)"}},
		Results: map[string]map[string]map[string]result{"AddWide": {
			"9.9.Add(1)": {"udecimal": {NsPerOp: 21, AllocsPerOp: 0}},
			"1234.Div(12345679012345679890123456789.1234567890123456789)": {"udecimal": {NsPerOp: 74.3, AllocsPerOp: 7}},
		}},
		Skips: map[string]map[string]map[string]string{"AddWide": {
			"9.9.Add(1)": {"govalues": "cannot represent: parsed 9.9 as 9"},
			"1234.Div(12345679012345679890123456789.1234567890123456789)": {"govalues": "cannot represent: overflow"},
		}},
	}

	out, err := caseTable(s, []string{"udecimal", "govalues"}, []string{"AddWide"}, classWide128)
	require.NoError(t, err)
	require.Contains(t, out, "| `AddWide` `9.9.Add(1)` | 21.0 ns / 0 | can't represent |\n")
	require.NotContains(t, out, "12345679012345679890123456789") // a classWideBig case
}

func TestReasonsList(t *testing.T) {
	s := &suite{
		Benches: []string{"AddWide", "SubWide"},
		Cases:   map[string][]string{"AddWide": {"9.9.Add(1)"}, "SubWide": {"9.9.Sub(1)"}},
		Skips: map[string]map[string]map[string]string{
			"AddWide": {"9.9.Add(1)": {"govalues": "cannot represent: parsed 9.9 as 9"}},
			"SubWide": {"9.9.Sub(1)": {"govalues": "cannot represent: parsed 9.9 as 9"}},
		},
	}

	require.Equal(t, "- `parsed 9.9 as 9`\n", reasonsList(s))
}
