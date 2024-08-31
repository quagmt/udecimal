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

// pre-computed values
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
//
// Number = coef / 10^(scale)
//
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

// NewFromInt64 returns a decimal which equals to coef / 10^scale.
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

	return newDecimal(neg, bintFromHiLo(0, uint64(coef)), scale)
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
// e.g. "123", "-123", "123.456", "-123.456", "+123.456", "0.123".
//
// Returns error if:
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

	if coef.IsZero() {
		return Zero, nil
	}

	return newDecimal(neg, coef, scale)
}

// MustParse parses a number in string to Decimal.
// Panic on error
func MustParse(s string) Decimal {
	d, err := Parse(s)
	if err != nil {
		panic(err)
	}

	return d
}

// Add returns d + e
//
// Returns overflow error when:
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

		return newDecimal(d.neg, coef, scale)
	}

	// different sign
	switch dcoef.Cmp(ecoef) {
	case 1:
		// dcoef > ecoef, subtract can't overflow
		coef, _ := dcoef.Sub(ecoef)
		return newDecimal(d.neg, coef, scale)
	default:
		// dcoef <= ecoef
		coef, _ := ecoef.Sub(dcoef)
		return newDecimal(e.neg, coef, scale)
	}
}

// Add64 returns d + e where e is a uint64
//
// Returns overflow error when:
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

		return newDecimal(neg, dcoef, d.scale)
	}

	dcoef, err := d.coef.Add(ecoef)
	if err != nil {
		return Decimal{}, err
	}

	return newDecimal(false, dcoef, d.scale)
}

// Sub returns d - e
//
// Returns overflow error when:
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

		return newDecimal(d.neg, coef, scale)
	}

	// same sign
	switch dcoef.Cmp(ecoef) {
	case 1:
		// dcoef > ecoef, subtract can't overflow
		coef, _ := dcoef.Sub(ecoef)
		return newDecimal(d.neg, coef, scale)
	default:
		// dcoef <= ecoef
		coef, _ := ecoef.Sub(dcoef)
		return newDecimal(!d.neg, coef, scale)
	}
}

// Sub64 returns d - e where e is a uint64
//
// Returns overflow error when:
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

		return newDecimal(neg, dcoef, d.scale)
	}

	dcoef, err := d.coef.Add(ecoef)
	if err != nil {
		return Decimal{}, err
	}

	return newDecimal(true, dcoef, d.scale)
}

// Mul returns d * e.
// If the result has more than 19 fraction digits, it will be truncated to 19 digits.
//
// Returns overflow error when:
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
		if !coef.carry.IsZero() {
			return Decimal{}, ErrOverflow
		}

		return newDecimal(neg, bintFromHiLo(coef.hi, coef.lo), scale)
	}

	rcoef, err := coef.quo(pow10[scale-MaxScale])
	if err != nil {
		return Decimal{}, err
	}

	return newDecimal(neg, rcoef, MaxScale)
}

// Mul64 returns d * e where e is a uint64.
// If the result has more than 19 fraction digits, it will be truncated to 19 digits.
//
// Returns overflow error when:
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

	return newDecimal(d.neg, coef, d.scale)
}

// Div returns d / e
// If the result has more than 19 fraction digits, it will be truncated to 19 digits.
//
// Returns overflow error when:
//  1. either whole or fration part is greater than 10^19-1
//  2. coef >= 2^128
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

	return newDecimal(neg, quo, MaxScale)
}

// Div64 returns d / e where e is a uint64
// If the result has more than 19 fraction digits, it will be truncated to 19 digits.
//
// Returns overflow error when:
//  1. either whole or fration part is greater than 10^19-1
//  2. coef >= 2^128
//
// Returns divide by zero error when e is zero
func (d Decimal) Div64(v uint64) (Decimal, error) {
	if v == 0 {
		return Decimal{}, ErrDivideByZero
	}

	d256 := d.coef.MulToU256(pow10[MaxScale-d.scale])
	quo, _, err := d256.quoRem64ToBint(v)
	if err != nil {
		return Decimal{}, err
	}

	return newDecimal(d.neg, quo, MaxScale)
}

// newDecimal return the decimal after removing all trailing zeros
func newDecimal(neg bool, coef bint, scale uint8) (Decimal, error) {
	if isOverflow(coef, scale) {
		return Decimal{}, ErrOverflow
	}

	return Decimal{neg: neg, coef: coef, scale: scale}, nil
}

// Scale returns decimal scale
func (d Decimal) Scale() int {
	return int(d.scale)
}

// cmp compares two decimals d,e and returns:
//
//	-1 if d < e
//	 0 if d == e
//	+1 if d > e
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
	// e has more fraction digits
	if d.scale < e.scale {
		// d has more fraction digits
		d256 := d.coef.MulToU256(pow10[e.scale-d.scale])
		return d256.Cmp128(e.coef)
	}

	// d has more fraction digits
	// we need to compare d with e * 10^(d.scale - e.scale)
	e256 := e.coef.MulToU256(pow10[d.scale-e.scale])

	// remember to reverse the result because e256.Cmp128(d.coef) returns the opposite
	return -e256.Cmp128(d.coef)
}

// Neg returns -d
func (d Decimal) Neg() Decimal {
	return Decimal{neg: !d.neg, coef: d.coef, scale: d.scale}
}

// Abs returns |d|
func (d Decimal) Abs() Decimal {
	return Decimal{neg: false, coef: d.coef, scale: d.scale}
}

// Sign returns:
//
//	-1 if d < 0
//	 0 if d == 0
//	+1 if d > 0
func (d Decimal) Sign() int {
	// check this first
	// because we allow parsing "-0" into decimal, which results in d.neg = true and d.coef = 0
	if d.coef.IsZero() {
		return 0
	}

	if d.neg {
		return -1
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
	return d.neg && !d.coef.IsZero()
}

// IsPos returns
//
//	true if d > 0
//	false if d <= 0
func (d Decimal) IsPos() bool {
	return !d.neg && !d.coef.IsZero()
}

// String returns the string representation of the decimal.
// Trailing zeros will be removed.
func (d Decimal) String() string {
	if d.IsZero() {
		return "0"
	}

	// max 40 bytes: 1 sign + 19 whole + 1 dot + 19 fraction
	buf := []byte("0000000000000000000000000000000000000000")
	n := d.writeToBytes(buf, true)

	return unsafeBytesToString(buf[n:])
}

func (d Decimal) writeToBytes(b []byte, trimTrailingZeros bool) int {
	if d.coef.IsZero() {
		return len(b) - 1
	}

	quo, rem := d.coef.QuoRem64(pow10[d.scale].lo)
	l := len(b)
	n := 0
	scale := d.scale

	if rem != 0 {
		if trimTrailingZeros {
			// remove trailing zeros, e.g. 1.2300 -> 1.23
			// both scale and rem will be adjusted
			zeros := getTrailingZeros64(rem)
			rem /= pow10[zeros].lo
			scale -= zeros
		}

		for ; rem != 0; rem /= 10 {
			n++
			b[l-n] += byte(rem % 10)
		}

		b[l-1-int(scale)] = '.'
		n = int(scale + 1)
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
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// RoundBank uses half up to even (banker's rounding) to round the decimal to the specified scale.
//
//	Examples:
//	Round(1.12345, 4) = 1.1234
//	Round(1.12335, 4) = 1.1234
//	Round(1.5, 0) = 2
//	Roung(-1.5, 0) = -2
func (d Decimal) RoundBank(scale uint8) (Decimal, error) {
	if scale >= d.scale {
		return d, nil
	}

	factor := pow10[d.scale-scale]
	lo := factor.lo / 2

	q, r := d.coef.QuoRem64(factor.lo)
	if lo < r || (lo == r && q.lo%2 == 1) {
		q, _ = q.Add64(1)
	}

	return newDecimal(d.neg, q, scale)
}

// RoundHAZ rounds the decimal to the specified scale using HALF AWAY FROM ZERO method (https://en.wikipedia.org/wiki/Rounding#Rounding_half_away_from_zero).
//
//	Examples:
//	Round(1.12345, 4) = 1.1235
//	Round(1.12335, 4) = 1.1234
//	Round(1.5, 0) = 2
//	Round(-1.5, 0) = -2
func (d Decimal) RoundHAZ(scale uint8) (Decimal, error) {
	if scale >= d.scale {
		return d, nil
	}

	factor := pow10[d.scale-scale]
	lo := factor.lo / 2

	q, r := d.coef.QuoRem64(factor.lo)
	if lo <= r {
		q, _ = q.Add64(1)
	}

	return newDecimal(d.neg, q, scale)
}

// RoundHTZ rounds the decimal to the specified scale using HALF TOWARD ZERO method (https://en.wikipedia.org/wiki/Rounding#Rounding_half_toward_zero).
//
//	Examples:
//	Round(1.12345, 4) = 1.1234
//	Round(1.12335, 4) = 1.1233
//	Round(1.5, 0) = 1
//	Round(-1.5, 0) = -1
func (d Decimal) RoundHTZ(scale uint8) (Decimal, error) {
	if scale >= d.scale {
		return d, nil
	}

	factor := pow10[d.scale-scale]
	lo := factor.lo / 2

	q, r := d.coef.QuoRem64(factor.lo)
	if lo < r {
		q, _ = q.Add64(1)
	}

	return newDecimal(d.neg, q, scale)
}

// Floor returns the largest integer value less than or equal to d.
//
//	Examples:
//	Floor(1.12345) = 1
//	Floor(1.12335) = 1
//	Floor(1.5, 0) = 1
//	Floor(-1.5, 0) = -2
func (d Decimal) Floor() (Decimal, error) {
	q, r := d.coef.QuoRem64(pow10[d.scale].lo)
	if d.neg && r != 0 {
		q, _ = q.Add64(1)
	}

	return newDecimal(d.neg, q, 0)
}

// Ceil returns the smallest integer value greater than or equal to d.
//
//	Examples:
//	Ceil(1.12345, 4) = 1.1235
//	Ceil(1.12335, 4) = 1.1234
//	Ceil(1.5, 0) = 2
//	Ceil(-1.5, 0) = -1
func (d Decimal) Ceil() (Decimal, error) {
	q, r := d.coef.QuoRem64(pow10[d.scale].lo)
	if !d.neg && r != 0 {
		q, _ = q.Add64(1)
	}

	return newDecimal(d.neg, q, 0)
}

// FMA (fused multiply-add) returns d*e + f in an efficient way
// and prevents intermediate rounding errors.
func (d Decimal) FMA(e Decimal, f Decimal) (Decimal, error) {
	// TODO: improve this
	return Decimal{}, nil
}

// func (d Decimal) Pow(e int) (Decimal, error) {
// 	if e == 0 {
// 		return One, nil
// 	}

// 	if e < 0 {
// 		return Decimal{}, fmt.Errorf("negative exponent is not supported")
// 	}

// 	res := One
// 	for i := 0; i < e; i++ {
// 		var err error
// 		res, err = res.Mul(d)
// 		if err != nil {
// 			return Decimal{}, err
// 		}
// 	}

// 	return res, nil
// }
