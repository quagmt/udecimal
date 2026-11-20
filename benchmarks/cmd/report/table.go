package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// headBenches are the head-to-head benchmarks, in the order the summary table lists them.
var headBenches = []string{"Add", "Sub", "Mul", "Div", "Pow", "Parse", "String", "MarshalText", "UnmarshalText", "MarshalJSON", "UnmarshalJSON"}

const cantRepresent = "can't represent"

func header(first string, names []string) string {
	var sb strings.Builder
	sb.WriteString("| " + first + " | " + strings.Join(names, " | ") + " |\n|---")
	sb.WriteString(strings.Repeat("|---", len(names)))
	sb.WriteString("|\n")

	return sb.String()
}

// summaryTable has one row per bench: udecimal's time, then every other library's time as a multiple of it, e.g.
// "5.4×" for 5.4 times as long. A value below 1 is faster than udecimal and is set in bold, e.g. "**0.87×**".
// Each value is the geometric mean of the bench's per-case medians.
func summaryTable(s *suite, names, benches []string) (string, error) {
	var sb strings.Builder
	sb.WriteString(header("Benchmark", names))

	for _, bench := range benches {
		bars, err := chartBars(s, names, bench, classHead)
		if err != nil {
			return "", err
		}

		base := bars[0] // names starts with udecimal
		if base.Missing {
			return "", fmt.Errorf("%s: udecimal has no result", bench)
		}

		cells := []string{bench, formatNs(base.Ns) + " ns"}
		for _, b := range bars[1:] {
			cells = append(cells, ratio(b, base.Ns))
		}

		sb.WriteString("| " + strings.Join(cells, " | ") + " |\n")
	}

	sb.WriteString("\nValues in bold are faster than udecimal.\n")

	return sb.String(), nil
}

// ratio formats b's time relative to udecimal's, e.g. 23.7 / 4.4 -> "5.4×", 5.5 / 6.3 -> "**0.87×**".
func ratio(b bar, base float64) string {
	if b.Missing {
		return cantRepresent
	}

	r := b.Ns / base
	if r < 1 {
		return fmt.Sprintf("**%.2f×**", r)
	}

	return fmt.Sprintf("%.1f×", r)
}

// allocTable has one row per bench with each library's allocs/op: one number when every case agrees, else the
// range over the cases, e.g. "0–6".
func allocTable(s *suite, names, benches []string) string {
	var sb strings.Builder
	sb.WriteString(header("Benchmark", names))

	for _, bench := range benches {
		cells := []string{bench}

		for _, lib := range names {
			var allocs []float64
			for _, cas := range s.Cases[bench] {
				if r, ok := s.Results[bench][cas][lib]; ok {
					allocs = append(allocs, r.AllocsPerOp)
				}
			}

			cells = append(cells, allocRange(allocs))
		}

		sb.WriteString("| " + strings.Join(cells, " | ") + " |\n")
	}

	return sb.String()
}

func allocRange(allocs []float64) string {
	if len(allocs) == 0 {
		return cantRepresent
	}

	lo, hi := slices.Min(allocs), slices.Max(allocs)
	if lo == hi {
		return formatAllocs(lo)
	}

	return formatAllocs(lo) + "–" + formatAllocs(hi)
}

func formatAllocs(a float64) string {
	return strconv.FormatFloat(a, 'f', -1, 64)
}

// caseTable has one row per case of benches in class c, each cell "<ns> ns / <allocs>", e.g. "21.0 ns / 0".
func caseTable(s *suite, names, benches []string, c class) (string, error) {
	var sb strings.Builder
	sb.WriteString(header("Case", names))

	for _, bench := range benches {
		for _, cas := range s.Cases[bench] {
			if classify(bench, cas) != c {
				continue
			}

			cells := []string{"`" + bench + "` `" + cas + "`"}

			for _, lib := range names {
				if r, ok := s.Results[bench][cas][lib]; ok {
					cells = append(cells, fmt.Sprintf("%s ns / %s", formatNs(r.NsPerOp), formatAllocs(r.AllocsPerOp)))
					continue
				}

				if _, ok := s.Skips[bench][cas][lib]; ok {
					cells = append(cells, cantRepresent)
					continue
				}

				return "", fmt.Errorf("%s/%s: no result or skip reason for %s", bench, cas, lib)
			}

			sb.WriteString("| " + strings.Join(cells, " | ") + " |\n")
		}
	}

	sb.WriteString("\nEach cell shows the median ns/op and allocs/op.\n")

	return sb.String(), nil
}

// reasonsList lists govalues' distinct skip reasons in first-seen order, without the "cannot represent: " prefix.
func reasonsList(s *suite) string {
	var seen []string
	for _, bench := range s.Benches {
		for _, cas := range s.Cases[bench] {
			reason, ok := s.Skips[bench][cas]["govalues"]
			if !ok {
				continue
			}

			reason = strings.TrimPrefix(reason, skipPrefix)
			if !slices.Contains(seen, reason) {
				seen = append(seen, reason)
			}
		}
	}

	var sb strings.Builder
	for _, r := range seen {
		sb.WriteString("- `" + r + "`\n")
	}

	return sb.String()
}
