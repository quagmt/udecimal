package udecimal

import (
	"fmt"
	"math"
	"strconv"
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

var (
	ErrOverflow      = fmt.Errorf("overflow. Number is out of range [-9_999_999_999_999_999_999.9_999_999_999_999_999_999, 9_999_999_999_999_999_999.9_999_999_999_999_999_999]")
	ErrMaxScale      = fmt.Errorf("scale out of range. Max digits after decimal point is %d", MaxScale)
	ErrEmptyString   = fmt.Errorf("parse empty string")
	ErrInvalidFormat = fmt.Errorf("invalid format")
	ErrDivideByZero  = fmt.Errorf("can't divide by zero")
)

var (
	Zero = Decimal{}
	One  = MustFromInt64(1, 0)
)

// Decimal represents a fixed-point decimal number.
// The number is represented as a coefficient and a scale.
// Number = coef / 10^(scale)
// For efficiency, both whole and fractional parts can only have 19 digits at most.
// Hence, the decimal range is:
// -9_999_999_999_999_999_999.9_999_999_999_999_999_999 <= D <= 9_999_999_999_999_999_999.9_999_999_999_999_999_999
type Decimal struct {
	coef  bint
	neg   bool // true if number is negative
	scale uint8
}

// isOverflow return true if the whole or fraction part has more than 19 digits
func isOverflow(coef bint, scale uint8) bool {
	// scale = frac digits
	// whole part has at most 19 digits
	// consider it's overflow when total digits > scale + 19, which means coef >= 10^(scale+19)
	return !coef.LessThan(pow10[scale+MaxScale])
}

// NewFromInt64 returns a decimal which equals to coef / 10^scale
// Trailing zeros wll be removed and the scale will also be adjusted
func NewFromInt64(coef int64, scale uint8) (Decimal, error) {
	var neg bool
	if coef < 0 {
		neg = true
		coef = -coef
	}

	if scale > MaxScale {
		return Decimal{}, ErrMaxScale
	}

	return newWithRTZ(neg, bintFromHiLo(0, uint64(coef)), scale)
}

// MustFromInt64 similars to NewFromInt64, but panics instead of returning error
func MustFromInt64(coef int64, scale uint8) Decimal {
	d, err := NewFromInt64(coef, scale)
	if err != nil {
		panic(err)
	}

	return d
}

// NewFromFloat64 returns decimal from float64 f
// !!!NOTE: you'll expect to lose some precision for this method due to FormatFloat. See: https://stackoverflow.com/questions/21895756/why-are-floating-point-numbers-inaccurate
// This method is suitable for small numbers with small precision. e.g. 1.0001, 0.0001, -123.456, -1000000.123456
// If you don't want to lose any precision, use Parse with string input instead
//
// Returns error if:
//  1. f is NaN or Inf
//  2. error when parsing float to string and then to decimal
func NewFromFloat64(f float64) (Decimal, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return Decimal{}, fmt.Errorf("%w: can't parse float '%v' to decimal", ErrInvalidFormat, f)
	}

	s := strconv.FormatFloat(f, 'f', -1, 64)
	d, err := Parse(s)
	if err != nil {
		return Decimal{}, fmt.Errorf("can't parse float: %w", err)
	}

	return d, nil
}

func MustFromFloat64(f float64) Decimal {
	d, err := NewFromFloat64(f)
	if err != nil {
		panic(err)
	}

	return d
}

// TODO: improve with SIMD
// Parse parses a number in string to Decimal.
// The string must be in the format of: [+-]d{1,19}[.d{1,19}]
// e.g. "123", "-123", "123.456", "-123.456", "+123.456", "0.123"
// Returns error if:
//
//  1. empty/invalid string
//  2. the number has whole or fraction part greater than 10^19-1
func Parse(s string) (Decimal, error) {
	errInvalidFormat := fmt.Errorf("%w: can't parse '%s' to decimal", ErrInvalidFormat, s)
	width := len(s)

	if width == 0 {
		return Decimal{}, ErrEmptyString
	}

	// max width = 1 + 19 + 1 + 19 = 40 (sign + whole + dot + fraction)
	if width > 40 {
		return Decimal{}, fmt.Errorf("%w: string length is greater than 40", ErrInvalidFormat)
	}

	var (
		pos int
		neg bool
	)

	switch s[0] {
	case '.':
		return Decimal{}, errInvalidFormat
	case '-':
		neg = true
		pos++
	case '+':
		pos++
	default:
		// do nothing
	}

	// prevent "+" or "-"
	if pos == width {
		return Decimal{}, errInvalidFormat
	}

	// prevent "-.123" or "+.123"
	if s[pos] == '.' {
		return Decimal{}, errInvalidFormat
	}

	var (
		err   error
		coef  bint
		scale uint8
	)
	for ; pos < width; pos++ {
		if s[pos] == '.' {
			// return err if we encounter the '.' more than once
			if scale != 0 {
				return Decimal{}, errInvalidFormat
			}

			scale = uint8(width - pos - 1)

			// prevent "123." or "-123."
			if scale == 0 {
				return Decimal{}, errInvalidFormat
			}

			if scale > MaxScale {
				return Decimal{}, ErrMaxScale
			}

			continue
		}

		if s[pos] < '0' || s[pos] > '9' {
			return Decimal{}, errInvalidFormat
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

	return newWithRTZ(neg, coef, scale)
}

// MustParse parses a number in string to Decimal
// Panic on error
func MustParse(s string) Decimal {
	d, err := Parse(s)
	if err != nil {
		panic(err)
	}

	return d
}

// Add returns d + e
// Returns overflow error when:
//
//  1. either whole or fration part is greater than 10^19-1
//  2. coef >= 2^128
func (d Decimal) Add(e Decimal) (Decimal, error) {
	dcoef, ecoef := d.coef, e.coef

	var (
		scale uint8
		err   error
	)
	switch {
	case d.scale == e.scale:
		scale = d.scale
	case d.scale > e.scale:
		scale = d.scale
		// can't overflow because scale is limited to 19
		// and whole part also has at most 19 digits
		// keep this check for safety, in case we change the limit in the future
		ecoef, err = ecoef.Mul(pow10[d.scale-e.scale])
		if err != nil {
			return Decimal{}, err
		}
	case d.scale < e.scale:
		scale = e.scale

		dcoef, err = dcoef.Mul(pow10[e.scale-d.scale])
		if err != nil {
			return Decimal{}, err
		}
	}

	if d.neg == e.neg {
		// same sign
		coef, err := dcoef.Add(ecoef)
		if err != nil {
			return Decimal{}, err
		}

		return newWithRTZ(d.neg, coef, scale)
	}

	// different sign
	switch dcoef.Cmp(ecoef) {
	case 1:
		// dcoef > ecoef, subtract can't overflow
		coef, _ := dcoef.Sub(ecoef)
		return newWithRTZ(d.neg, coef, scale)
	default:
		// dcoef <= ecoef
		coef, _ := ecoef.Sub(dcoef)
		return newWithRTZ(e.neg, coef, scale)
	}
}

// Add64 returns d + e where e is a uint64
// Returns overflow error when:
//
//  1. either whole or fration part is greater than 10^19-1
//  2. coef >= 2^128
func (d Decimal) Add64(e uint64) (Decimal, error) {
	ecoef, err := bintFromHiLo(0, e).Mul(pow10[d.scale])
	if err != nil {
		return Decimal{}, err
	}

	if d.neg {
		var (
			dcoef bint
			neg   bool
		)

		if d.coef.GreaterThan(ecoef) {
			dcoef, err = d.coef.Sub(ecoef)
			neg = true
		} else {
			dcoef, err = ecoef.Sub(d.coef)
			neg = false
		}

		if err != nil {
			return Decimal{}, err
		}

		return newWithRTZ(neg, dcoef, d.scale)
	}

	dcoef, err := d.coef.Add(ecoef)
	if err != nil {
		return Decimal{}, err
	}

	return newWithRTZ(false, dcoef, d.scale)
}

// Sub returns d - e
// Returns overflow error when:
//
//  1. either whole or fration part is greater than 10^19-1
//  2. coef >= 2^128
func (d Decimal) Sub(e Decimal) (Decimal, error) {
	dcoef, ecoef := d.coef, e.coef

	var (
		scale uint8
		err   error
	)

	switch {
	case d.scale == e.scale:
		scale = d.scale
	case d.scale > e.scale:
		scale = d.scale

		// can't overflow because scale is limited to 19
		// and whole part also has at most 19 digits
		// keep this check for safety, in case we change the limit in the future
		ecoef, err = ecoef.Mul(pow10[d.scale-e.scale])
		if err != nil {
			return Decimal{}, err
		}
	case d.scale < e.scale:
		scale = e.scale

		dcoef, err = dcoef.Mul(pow10[e.scale-d.scale])
		if err != nil {
			return Decimal{}, err
		}
	}

	if d.neg != e.neg {
		// different sign
		coef, err := dcoef.Add(ecoef)
		if err != nil {
			return Decimal{}, err
		}

		return newWithRTZ(d.neg, coef, scale)
	}

	// same sign
	switch dcoef.Cmp(ecoef) {
	case 1:
		// dcoef > ecoef, subtract can't overflow
		coef, _ := dcoef.Sub(ecoef)
		return newWithRTZ(d.neg, coef, scale)
	default:
		// dcoef <= ecoef
		coef, _ := ecoef.Sub(dcoef)
		return newWithRTZ(!d.neg, coef, scale)
	}
}

// Sub64 returns d - e where e is a uint64
// Returns overflow error when:
//
//  1. either whole or fration part is greater than 10^19-1
//  2. coef >= 2^128
func (d Decimal) Sub64(e uint64) (Decimal, error) {
	ecoef, err := bintFromHiLo(0, e).Mul(pow10[d.scale])
	if err != nil {
		return Decimal{}, err
	}

	if !d.neg {
		var (
			dcoef bint
			neg   bool
		)

		if d.coef.GreaterThan(ecoef) {
			dcoef, err = d.coef.Sub(ecoef)
			neg = false
		} else {
			dcoef, err = ecoef.Sub(d.coef)
			neg = true
		}

		if err != nil {
			return Decimal{}, err
		}

		return newWithRTZ(neg, dcoef, d.scale)
	}

	dcoef, err := d.coef.Add(ecoef)
	if err != nil {
		return Decimal{}, err
	}

	return newWithRTZ(true, dcoef, d.scale)
}

// Mul returns d * e
// If the result has more than 19 fraction digits, it will be truncated to 19 digits.
// Returns overflow error when:
//
//  1. either whole or fration part is greater than 10^19-1
//  2. coef >= 2^128
func (d Decimal) Mul(e Decimal) (Decimal, error) {
	if e.coef.IsZero() {
		return Decimal{}, nil
	}

	scale := d.scale + e.scale
	neg := d.neg != e.neg
	coef := d.coef.MulToU256(e.coef)

	if scale <= MaxScale {
		if coef.Cmp128(pow10[scale+MaxScale]) >= 0 {
			return Decimal{}, ErrOverflow
		}

		return newWithRTZ(neg, bintFromHiLo(coef.hi, coef.lo), scale)
	}

	rcoef, err := coef.quo(pow10[scale-MaxScale])
	if err != nil {
		return Decimal{}, err
	}

	return newWithRTZ(neg, rcoef, MaxScale)
}

// Mul64 returns d * e where e is a uint64
// If the result has more than 19 fraction digits, it will be truncated to 19 digits.
// Returns overflow error when:
//
//  1. either whole or fration part is greater than 10^19-1
//  2. coef >= 2^128
func (d Decimal) Mul64(v uint64) (Decimal, error) {
	if v == 0 {
		return Decimal{}, nil
	}

	if v == 1 {
		return d, nil
	}

	coef, err := d.coef.Mul64(v)
	if err != nil {
		return Decimal{}, err
	}

	return newWithRTZ(d.neg, coef, d.scale)
}

// Div returns d / e
// If the result has more than 19 fraction digits, it will be truncated to 19 digits.
// Returns overflow error when:
//
// 1. either whole or fration part is greater than 10^19-1
// 2. coef >= 2^128
//
// Returns divide by zero error when e is zero
func (d Decimal) Div(e Decimal) (Decimal, error) {
	if e.coef.IsZero() {
		return Decimal{}, ErrDivideByZero
	}

	neg := d.neg != e.neg

	// Need to multiply divident with factor
	// to make sure the total decimal number after the decimal point is MaxScale
	factor := MaxScale - (d.scale - e.scale)

	d256 := d.coef.MulToU256(pow10[factor])
	quo, err := d256.quo(e.coef)
	if err != nil {
		return Decimal{}, err
	}

	return newWithRTZ(neg, quo, MaxScale)
}

func (d Decimal) Div64(v uint64) (Decimal, error) {
	if v == 0 {
		return Decimal{}, ErrDivideByZero
	}

	d256 := d.coef.MulToU256(pow10[MaxScale])
	quo, _, err := d256.quoRem64ToBint(v)
	if err != nil {
		return Decimal{}, err
	}

	return newWithRTZ(d.neg, quo, MaxScale)
}

// newWithRTZ return the decimal after removing all trailing zeros
func newWithRTZ(neg bool, coef bint, scale uint8) (Decimal, error) {
	if scale == 0 {
		if isOverflow(coef, 0) {
			return Decimal{}, ErrOverflow
		}

		return Decimal{neg: neg, coef: coef, scale: 0}, nil
	}

	trailingZeros := getTrailingZeros(coef)
	if trailingZeros > 0 {
		coef, _ = coef.QuoRem64(pow10[trailingZeros].lo)
	}

	rescale := scale - trailingZeros

	if isOverflow(coef, rescale) {
		return Decimal{}, ErrOverflow
	}

	return Decimal{neg: neg, coef: coef, scale: rescale}, nil
}

func (d Decimal) Scale() int {
	return int(d.scale)
}

func (d Decimal) Cmp(e Decimal) int {
	if d.neg && !e.neg {
		return -1
	}

	if !d.neg && e.neg {
		return 1
	}

	// d.neg = e.neg
	if d.neg {
		// both are negative, return the opposite
		return -d.cmpDec(e)
	}

	return d.cmpDec(e)
}

func (d Decimal) cmpDec(e Decimal) int {
	if d.scale == e.scale {
		return d.coef.Cmp(e.coef)
	}

	// scale is different
	if d.scale < e.scale {
		// d has more fraction digits
		d256 := d.coef.MulToU256(pow10[e.scale-d.scale])
		return d256.Cmp128(e.coef)
	}

	// v has more fraction digits
	v256 := e.coef.MulToU256(pow10[d.scale-e.scale])
	return v256.Cmp128(d.coef)
}

func (d Decimal) Neg() Decimal {
	return Decimal{neg: !d.neg, coef: d.coef, scale: d.scale}
}

func (d Decimal) Abs() Decimal {
	return Decimal{neg: false, coef: d.coef, scale: d.scale}
}

// Sign returns:
//
//	-1 if d < 0
//	 0 if d == 0
//	+1 if d > 0
func (d Decimal) Sign() int {
	if d.neg {
		return -1
	}

	if d.coef.IsZero() {
		return 0
	}

	return 1
}

// IsZero returns
//
//	true if d == 0
//	false if d != 0
func (d Decimal) IsZero() bool {
	return d.coef.IsZero()
}

// IsNeg returns
//
//	true if d < 0
//	false if d >= 0
func (d Decimal) IsNeg() bool {
	return d.neg
}

// IsPos returns
//
//	true if d > 0
//	false if d <= 0
func (d Decimal) IsPos() bool {
	return !d.neg && !d.coef.IsZero()
}

func (d Decimal) String() string {
	if d.coef.IsZero() {
		return "0"
	}

	buf := []byte("0000000000000000000000000000000000000000") // log10(2^128) < 40
	n := d.writeToBytes(buf)

	return unsafeBytesToString(buf[n:])
}

func (d Decimal) writeToBytes(b []byte) int {
	if d.coef.IsZero() {
		b[0] = '0'
		return 1
	}

	quo, rem := d.coef.QuoRem64(pow10[d.scale].lo)
	l := len(b)
	n := 0

	if rem != 0 {
		for ; rem != 0; rem /= 10 {
			n++
			b[l-n] += byte(rem % 10)
		}

		b[l-1-int(d.scale)] = '.'
		n = int(d.scale + 1)
	}

	qlo := quo.lo
	if qlo != 0 {
		for ; qlo != 0; qlo /= 10 {
			n++
			b[l-n] += byte(qlo % 10)
		}
	} else {
		// quo is zero, so we need to print at least one zero
		n++
	}

	if d.neg {
		n++
		b[l-n] = '-'
	}

	return l - n
}

func unsafeBytesToString(b []byte) string {
	// Ignore if your IDE shows an error here; it's a false positive.
	p := unsafe.SliceData(b)
	return unsafe.String(p, len(b))
}

// func unsafeStringToBytes(s string) []byte {
// 	// unsafe.StringData output is unspecified for empty string input so always
// 	// return nil.
// 	if len(s) == 0 {
// 		return nil
// 	}

// 	return unsafe.Slice(unsafe.StringData(s), len(s))
// }
