package main

import (
	"fmt"
	"html"
	"strings"
)

// bar is one library's value in a chart: the geometric mean of its median ns/op over the chart's cases.
type bar struct {
	Lib     string
	Ns      float64
	Missing bool // the library skipped every case: it can't represent these inputs
}

// Chart geometry, in viewBox units. The page lays charts out about 400px wide, so one unit is about one CSS px:
// 15 units stay at least 12px down to a 320px-wide column.
const (
	chartWidth  = 400.0
	chartHeight = 200.0
	fontSize    = 15.0
	lineHeight  = 18.0
	topMargin   = fontSize + 10    // room for the value label above the tallest bar
	nameArea    = 2*lineHeight + 8 // two lines of library name under the baseline
	plotHeight  = chartHeight - topMargin - nameArea
	baseline    = topMargin + plotHeight
	sidePadding = 4.0
	charWidth   = 0.56 * fontSize // a generous average glyph width for the theme's sans-serif font
)

// nameBreaks splits a library name onto two lines when it is wider than its slot.
var nameBreaks = map[string][2]string{
	"shopspring":    {"shop-", "spring"},
	"ericlagergren": {"eric-", "lagergren"},
}

const (
	// bars need 3:1 against the page (WCAG non-text contrast): --lighter is about 2.2:1 on white, --light passes; the
	// accent, not the primary colour, for udecimal: indigo primary is 2.3:1 against the dark background, the accent 3.8:1
	barStyle   = "fill:var(--md-default-fg-color--light)"
	ownStyle   = "fill:var(--md-accent-fg-color)"
	textStyle  = "fill:var(--md-default-fg-color)"
	mutedStyle = "fill:var(--md-default-fg-color--light)"
	axisStyle  = "stroke:var(--md-default-fg-color--lighter)"
)

// renderChart draws one vertical bar per library, udecimal in the theme's accent colour and the others neutral, with
// the value above each bar and the library name under it. A missing library gets "can't represent" in its slot.
// Colours are CSS variables, so the chart follows the site's light and dark palettes when it is inlined into a page.
// The output has no blank lines: Markdown would end the raw HTML block there.
func renderChart(title string, bars []bar) string {
	maxNs := 0.0
	for _, b := range bars {
		maxNs = max(maxNs, b.Ns)
	}

	slot := (chartWidth - 2*sidePadding) / float64(len(bars))
	barWidth := slot * 0.56

	var labels []string
	for _, b := range bars {
		labels = append(labels, b.Lib+" "+valueText(b))
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "<figure class=\"perf-chart\">\n<figcaption>%s</figcaption>\n", html.EscapeString(title))
	fmt.Fprintf(&sb, "<svg viewBox=\"0 0 %g %g\" width=\"100%%\" role=\"img\" aria-label=\"%s, geometric mean of median ns/op: %s\" font-size=\"%g\" text-anchor=\"middle\">\n",
		chartWidth, chartHeight, html.EscapeString(title), html.EscapeString(strings.Join(labels, ", ")), fontSize)
	fmt.Fprintf(&sb, "<line x1=\"%g\" y1=\"%g\" x2=\"%g\" y2=\"%g\" style=\"%s\"/>\n", sidePadding, baseline, chartWidth-sidePadding, baseline, axisStyle)

	for i, b := range bars {
		center := sidePadding + slot*(float64(i)+0.5)

		if b.Missing {
			// two short lines fit a slot where "can't represent" on one line would not
			fmt.Fprintf(&sb, "<text x=\"%.1f\" y=\"%.1f\" style=\"%s\">%s</text>\n", center, baseline-lineHeight-6, mutedStyle, html.EscapeString("can't"))
			fmt.Fprintf(&sb, "<text x=\"%.1f\" y=\"%.1f\" style=\"%s\">represent</text>\n", center, baseline-6, mutedStyle)
		} else {
			h := 0.0
			if maxNs > 0 {
				h = b.Ns / maxNs * plotHeight
			}

			style := barStyle
			if b.Lib == "udecimal" {
				style = ownStyle
			}

			fmt.Fprintf(&sb, "<rect data-lib=\"%s\" x=\"%.1f\" y=\"%.2f\" width=\"%.1f\" height=\"%.2f\" rx=\"2\" style=\"%s\"/>\n",
				html.EscapeString(b.Lib), center-barWidth/2, baseline-h, barWidth, h, style)
			fmt.Fprintf(&sb, "<text x=\"%.1f\" y=\"%.1f\" style=\"%s\">%s</text>\n", center, baseline-h-6, textStyle, valueText(b))
		}

		for j, line := range nameLines(b.Lib, slot) {
			fmt.Fprintf(&sb, "<text x=\"%.1f\" y=\"%.1f\" style=\"%s\">%s</text>\n", center, baseline+lineHeight*float64(j+1), textStyle, html.EscapeString(line))
		}
	}

	sb.WriteString("</svg>\n</figure>")

	return sb.String()
}

// valueText is the label above a bar, e.g. "4.4 ns".
func valueText(b bar) string {
	if b.Missing {
		return "can't represent"
	}

	return formatNs(b.Ns) + " ns"
}

// nameLines returns the library name as one line, or two when it is wider than the slot, e.g. "eric-", "lagergren".
func nameLines(lib string, slot float64) []string {
	parts, ok := nameBreaks[lib]
	if !ok || float64(len(lib))*charWidth <= slot-2*sidePadding {
		return []string{lib}
	}

	return parts[:]
}

// renderGroup lays charts out with Material's grid: two per row on a desktop, one per row on a phone.
func renderGroup(figures []string) string {
	return "<div class=\"grid\">\n" + strings.Join(figures, "\n") + "\n</div>\n"
}

// chartBars returns one bar per library for bench, over its cases in class c. A library that skipped every case is
// Missing. A case with neither a result nor a skip reason, or a library that skipped only some cases, is an error:
// the run is broken and a mean over it would mislead.
func chartBars(s *suite, names []string, bench string, c class) ([]bar, error) {
	var cases []string
	for _, cas := range s.Cases[bench] {
		if classify(bench, cas) == c {
			cases = append(cases, cas)
		}
	}

	if len(cases) == 0 {
		return nil, fmt.Errorf("%s: no cases in this class", bench)
	}

	bars := make([]bar, 0, len(names))
	for _, lib := range names {
		var ns []float64
		skipped := 0

		for _, cas := range cases {
			if r, ok := s.Results[bench][cas][lib]; ok {
				ns = append(ns, r.NsPerOp)
				continue
			}

			if _, ok := s.Skips[bench][cas][lib]; ok {
				skipped++
				continue
			}

			return nil, fmt.Errorf("%s/%s: no result or skip reason for %s", bench, cas, lib)
		}

		switch {
		case skipped == len(cases):
			bars = append(bars, bar{Lib: lib, Missing: true})
		case skipped > 0:
			return nil, fmt.Errorf("%s: %s skipped %d of %d cases; a mean over the rest would mislead", bench, lib, skipped, len(cases))
		default:
			bars = append(bars, bar{Lib: lib, Ns: geomean(ns)})
		}
	}

	return bars, nil
}
