package udecimal

// Decimal represents a fixed-point decimal number.
// The number is represented as a coefficient and a scale.
// Number = coef / 10^(scale)
// For efficiency, max decimal places is 19 and
type Decimal struct {
	neg  bool // true if number is negative
	coef uint
}
