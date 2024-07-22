package udecimal

import "fmt"

// Decimal represents a fixed-point decimal number.
// The number is represented as a coefficient and a scale.
// Number = coef / 10^(scale)
// For efficiency, max decimal places is 19 and
type Decimal struct {
	neg   bool // true if number is negative
	coef  Uint128
	scale int8
}

func Parse(s string) (Decimal, error) {
	var pos int
	width := len(s)

	if width == 0 {
		return Decimal{}, fmt.Errorf("empty string")
	}

	var neg bool
	switch s[0] {
	case '-':
		neg = true
		pos++
	case '+':
		pos++
	default:
		// do nothing
	}

	var (
		coef  Uint128
		scale int8
	)

	for ; pos < width; pos++ {
		if s[pos] == '.' {
			scale = int8(width - pos - 1)
			continue
		}

		if s[pos] < '0' || s[pos] > '9' {
			return Decimal{}, fmt.Errorf("invalid character: %c", s[pos])
		}

		coef = coef.Mul64(10).Add64(uint64(s[pos] - '0'))
	}

	return Decimal{neg: neg, coef: coef, scale: scale}, nil
}
