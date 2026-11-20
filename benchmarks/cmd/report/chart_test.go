package main

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// rectOf returns the <rect> drawn for lib, found by its data-lib attribute.
func rectOf(t *testing.T, svg, lib string) string {
	t.Helper()

	m := regexp.MustCompile(`<rect data-lib="` + lib + `"[^>]*>`).FindString(svg)
	require.NotEmpty(t, m, "no rect for %s", lib)

	return m
}

func rectHeight(t *testing.T, rect string) float64 {
	t.Helper()

	m := regexp.MustCompile(`height="([\d.]+)"`).FindStringSubmatch(rect)
	require.NotNil(t, m, rect)

	h, err := strconv.ParseFloat(m[1], 64)
	require.NoError(t, err)

	return h
}

func TestRenderChartHighlightsUdecimal(t *testing.T) {
	out := renderChart("Add", []bar{{Lib: "udecimal", Ns: 4.4}, {Lib: "shopspring", Ns: 23.7}})

	require.Contains(t, rectOf(t, out, "udecimal"), "var(--md-accent-fg-color)")
	require.Contains(t, rectOf(t, out, "shopspring"), "var(--md-default-fg-color--light)")
	require.Contains(t, out, "4.4 ns")
	require.Contains(t, out, "23.7 ns")
}

func TestRenderChartMissing(t *testing.T) {
	out := renderChart("AddWide", []bar{{Lib: "udecimal", Ns: 5.5}, {Lib: "govalues", Missing: true}})

	require.NotContains(t, out, `<rect data-lib="govalues"`)
	require.Contains(t, out, "can&#39;t")
	require.Contains(t, out, "represent")
}

func TestRenderChartScale(t *testing.T) {
	out := renderChart("Div", []bar{{Lib: "udecimal", Ns: 10}, {Lib: "govalues", Ns: 20}})

	half, full := rectHeight(t, rectOf(t, out, "udecimal")), rectHeight(t, rectOf(t, out, "govalues"))
	require.InDelta(t, plotHeight, full, 0.01)
	require.InDelta(t, full/2, half, 1)
}

func TestRenderChartNoBlankLines(t *testing.T) {
	out := renderChart("Add", []bar{{Lib: "udecimal", Ns: 4.4}, {Lib: "ericlagergren", Ns: 12.5}, {Lib: "govalues", Missing: true}})

	for i, line := range strings.Split(out, "\n") {
		require.NotEmpty(t, strings.TrimSpace(line), "line %d is blank", i)
	}
}

func TestRenderChartEscapes(t *testing.T) {
	out := renderChart("A<B&C", []bar{{Lib: "udecimal", Ns: 1}})

	require.Contains(t, out, "A&lt;B&amp;C")
	require.NotContains(t, out, "A<B&C")
}

func TestRenderChartThemeVars(t *testing.T) {
	out := renderChart("Add", []bar{{Lib: "udecimal", Ns: 4.4}, {Lib: "shopspring", Ns: 23.7}, {Lib: "govalues", Missing: true}})

	allowed := []string{"--md-accent-fg-color", "--md-default-fg-color", "--md-default-fg-color--light", "--md-default-fg-color--lighter"}
	for _, m := range regexp.MustCompile(`var\((--[\w-]+)\)`).FindAllStringSubmatch(out, -1) {
		require.Contains(t, allowed, m[1])
	}
}

func TestRenderChartAriaLabel(t *testing.T) {
	out := renderChart("Add", []bar{{Lib: "udecimal", Ns: 4.4}, {Lib: "shopspring", Ns: 23.7}, {Lib: "govalues", Missing: true}})

	require.Contains(t, out, `role="img"`)
	require.Contains(t, out, "Add, geometric mean of median ns/op: ")
	require.Contains(t, out, "udecimal 4.4 ns, shopspring 23.7 ns, govalues can&#39;t represent")
}

func TestRenderGroup(t *testing.T) {
	out := renderGroup([]string{"<figure>a</figure>", "<figure>b</figure>"})

	require.Equal(t, "<div class=\"grid\">\n<figure>a</figure>\n<figure>b</figure>\n</div>\n", out)
}

func TestChartBars(t *testing.T) {
	s, err := load(fixtures)
	require.NoError(t, err)

	bars, err := chartBars(s, []string{"udecimal", "govalues"}, "Add", classHead)
	require.NoError(t, err)
	require.Len(t, bars, 2)
	require.Equal(t, "udecimal", bars[0].Lib)
	require.InDelta(t, 11, bars[0].Ns, 1e-9)
	require.Equal(t, "govalues", bars[1].Lib)
	require.InDelta(t, 55, bars[1].Ns, 1e-9)

	bars, err = chartBars(s, []string{"udecimal", "govalues"}, "AddWide", classWide128)
	require.NoError(t, err)
	require.Len(t, bars, 2)
	require.InDelta(t, 21, bars[0].Ns, 1e-9)
	require.Equal(t, bar{Lib: "govalues", Missing: true}, bars[1])
}

func TestChartBarsMissingResult(t *testing.T) {
	// shopspring has neither a result nor a skip reason for Add/1.Add(2)
	s, err := load(fixtures)
	require.NoError(t, err)

	_, err = chartBars(s, []string{"udecimal", "shopspring"}, "Add", classHead)
	require.ErrorContains(t, err, "Add/1.Add(2)")
	require.ErrorContains(t, err, "shopspring")
}
