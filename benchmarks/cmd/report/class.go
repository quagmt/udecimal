package main

import (
	"fmt"
	"strings"
)

// class says which part of the Performance page a case belongs to.
type class int

const (
	classHead    class = iota // head-to-head: every library takes part, at most 19 digits
	classWide128              // more than 19 digits, at most 38: udecimal stays on its 128-bit path
	classWideBig              // an operand of more than 38 digits: udecimal falls back to math/big
)

// classify puts a case of a "*Wide" benchmark into classWide128 or classWideBig by its longest operand,
// counted in significant digits, e.g. "0.000123" has 3. Every other benchmark is classHead.
func classify(bench, cas string) class {
	if !strings.HasSuffix(bench, "Wide") {
		return classHead
	}

	longest := 0
	for _, n := range number.FindAllString(cas, -1) {
		digits := strings.TrimLeft(strings.ReplaceAll(n, ".", ""), "0")
		longest = max(longest, len(digits))
	}

	if longest > 38 {
		return classWideBig
	}

	return classWide128
}

// checkedWide are the wide benchmarks whose classWide128 cases must not allocate in udecimal.
// StringWide is left out: String always allocates the string it returns.
var checkedWide = []string{"ParseWide", "AddWide", "SubWide", "MulWide", "DivWide"}

// checkClasses fails when a case classified as 128-bit allocates in udecimal, so the page never puts a math/big
// case under the "still 128-bit" heading.
func checkClasses(s *suite) error {
	for _, bench := range checkedWide {
		for _, cas := range s.Cases[bench] {
			if classify(bench, cas) != classWide128 {
				continue
			}

			if r, ok := s.Results[bench][cas]["udecimal"]; ok && r.AllocsPerOp > 0 {
				return fmt.Errorf("%s/%s: classified as 128-bit but udecimal allocates %v times", bench, cas, r.AllocsPerOp)
			}
		}
	}

	return nil
}
