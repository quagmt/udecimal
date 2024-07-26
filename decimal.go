package udecimal

import (
	"fmt"
	"unsafe"
)

const (
	MaxScale = 19
)

var pow10 = [39]bint{
	{lo: 1},                                  // 10^0
	{lo: 10},                                 // 10^1
	{lo: 1e2},                                // 10^2
	{lo: 1e3},                                // 10^3
	{lo: 1e4},                                // 10^4
	{lo: 1e5},                                // 10^5
	{lo: 1e6},                                // 10^6
	{lo: 1e7},                                // 10^7
	{lo: 1e8},                                // 10^8
	{lo: 1e9},                                // 10^9
	{lo: 1e10},                               // 10^10
	{lo: 1e11},                               // 10^11
	{lo: 1e12},                               // 10^12
	{lo: 1e13},                               // 10^13
	{lo: 1e14},                               // 10^14
	{lo: 1e15},                               // 10^15
	{lo: 1e16},                               // 10^16
	{lo: 1e17},                               // 10^17
	{lo: 1e18},                               // 10^18
	{lo: 1e19},                               // 10^19
	{lo: 7_766_279_631_452_241_920, hi: 5},   // 10^20
	{lo: 3_875_820_019_684_212_736, hi: 54},  // 10^21
	{lo: 1_864_712_049_423_024_128, hi: 542}, // 10^22
	{lo: 200_376_420_520_689_664, hi: 5_421}, // 10^23
	{lo: 2_003_764_205_206_896_640, hi: 54_210},                  // 10^24
	{lo: 1_590_897_978_359_414_784, hi: 542_101},                 // 10^25
	{lo: 15_908_979_783_594_147_840, hi: 5_421_010},              // 10^26
	{lo: 11_515_845_246_265_065_472, hi: 54_210_108},             // 10^27
	{lo: 4_477_988_020_393_345_024, hi: 542_101_086},             // 10^28
	{lo: 7_886_392_056_514_347_008, hi: 5_421_010_862},           // 10^29
	{lo: 5_076_944_270_305_263_616, hi: 54_210_108_624},          // 10^30
	{lo: 1_387_595_455_563_353_2928, hi: 542_101_086_242},        // 10^31
	{lo: 9_632_337_040_368_467_968, hi: 5_421_010_862_427},       // 10^32
	{lo: 4_089_650_035_136_921_600, hi: 54_210_108_624_275},      // 10^33
	{lo: 4_003_012_203_950_112_768, hi: 542_101_086_242_752},     // 10^34
	{lo: 3_136_633_892_082_024_448, hi: 5_421_010_862_427_522},   // 10^35
	{lo: 12_919_594_847_110_692_864, hi: 54_210_108_624_275_221}, // 10^36
	{lo: 68_739_955_140_067_328, hi: 542_101_086_242_752_217},    // 10^37
	{lo: 687_399_551_400_673_280, hi: 5_421_010_862_427_522_170}, // 10^38
}

// Decimal represents a fixed-point decimal number.
// The number is represented as a coefficient and a scale.
// Number = coef / 10^(scale)
// For efficiency, both whole and fractional parts can only have 19 digits at most.
type Decimal struct {
	coef  bint
	neg   bool // true if number is negative
	scale int8
}

func isOverflow(coef bint, scale int8) bool {
	// scale = frac digits
	// whole part has at most 19 digits
	// consider it's overflow when total digits > scale + 19 --> coef >= 10^(scale+19)
	return !coef.LessThan(pow10[scale+19])
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
		err   error
		coef  bint
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

		coef, err = coef.Mul64(10)
		if err != nil {
			return Decimal{}, err
		}

		coef, err = coef.Add64(uint64(s[pos] - '0'))
		if err != nil {
			return Decimal{}, err
		}
	}

	if isOverflow(coef, scale) {
		return Decimal{}, ErrOverflow
	}

	return Decimal{neg: neg, coef: coef, scale: scale}, nil
}

func (u Decimal) Div(v Decimal) (Decimal, error) {
	var neg bool
	if u.neg != v.neg {
		neg = true
	}

	if v.coef.IsZero() {
		return Decimal{}, fmt.Errorf("division by zero")
	}

	// Need to multiply divident with factor
	// to make sure the total decimal number after the decimal point is MaxScale
	factor := MaxScale - (u.scale - v.scale)
	divident := u.coef.MulPow10(int(factor))
	quo, _, err := divident.QuoRem(v.coef)
	if err != nil {
		return Decimal{}, err
	}

	// rescale the fraction part by removing trailing zeros
	trailingZeros := getTrailingZeros(quo)
	if trailingZeros > 0 {
		quo, _ = quo.QuoRem64(pow10[trailingZeros-1].lo)
	}

	return Decimal{neg: neg, coef: quo, scale: int8(MaxScale - trailingZeros)}, nil
}

func (u Decimal) String() string {
	if u.coef.IsZero() {
		return "0"
	}

	quo, rem := u.coef.QuoRem64(pow10[u.scale].lo)
	buf := []byte("0000000000000000000000000000000000000000") // log10(2^128) < 40

	i := len(buf)
	n := 0

	if rem != 0 {
		for ; rem != 0; rem /= 10 {
			n++
			buf[i-n] += byte(rem % 10)
		}

		buf[i-1-int(u.scale)] = '.'
		n = int(u.scale + 1)
	}

	qlo := quo.lo
	for ; qlo != 0; qlo /= 10 {
		n++
		buf[i-n] += byte(qlo % 10)
	}

	if u.neg {
		n++
		buf[i-n] = '-'
	}

	p := buf[i-n:]
	return unsafe.String(unsafe.SliceData(p), len(p))
}

func getTrailingZeros(coef bint) int8 {
	var z int8 = 0

	q, rem := coef.QuoRem64(pow10[16].lo)
	if rem == 0 {
		z += 16

		// short path because max scale is 19
		q, rem = q.QuoRem64(100)
		if rem == 0 {
			z += 2
		}

		q, rem = q.QuoRem64(10)
		if rem == 0 {
			z += 1
		}

		return z
	}

	q, rem = q.QuoRem64(pow10[8].lo)
	if rem == 0 {
		z += 8
	}

	q, rem = q.QuoRem64(10_000)
	if rem == 0 {
		z += 4
	}

	q, rem = q.QuoRem64(100)
	if rem == 0 {
		z += 2
	}

	q, rem = coef.QuoRem64(10)
	if rem == 0 {
		z += 1
	}

	return z
}
