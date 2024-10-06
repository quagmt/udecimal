package udecimal

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
)

var (
	// defaultScale is the default number of digits after the decimal point
	// if not specified
	defaultScale uint8 = 19

	// maxScale is the maximum number of digits after the decimal point
	maxScale uint8 = 19

	// maxStrLen is the maximum length of string input when using Parse/MustParse
	// set it to 200 so string length value can fit in 1 byte (for MarshalBinary).
	// Also such that big number (more than 200 digits) is unrealistic in financial system
	// which this library is mainly designed for
	maxStrLen = 200
)

// pre-computed values
var pow10 = [39]u128{
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

var pow10Big = [20]*big.Int{
	big.NewInt(1),        // 10^0
	big.NewInt(10),       // 10^1
	big.NewInt(1e2),      // 10^2
	big.NewInt(1e3),      // 10^3
	big.NewInt(1e4),      // 10^4
	big.NewInt(1e5),      // 10^5
	big.NewInt(1e6),      // 10^6
	big.NewInt(1e7),      // 10^7
	big.NewInt(1e8),      // 10^8
	big.NewInt(1e9),      // 10^9
	big.NewInt(1e10),     // 10^10
	big.NewInt(1e11),     // 10^11
	big.NewInt(1e12),     // 10^12
	big.NewInt(1e13),     // 10^13
	big.NewInt(1e14),     // 10^14
	big.NewInt(1e15),     // 10^15
	big.NewInt(1e16),     // 10^16
	big.NewInt(1e17),     // 10^17
	big.NewInt(1e18),     // 10^18
	pow10[19].ToBigInt(), // 10^19
}

var (
	ErrOverflow        = fmt.Errorf("overflow")
	ErrScaleOutOfRange = fmt.Errorf("scale out of range. Only support maximum %d digits after the decimal point", defaultScale)
	ErrEmptyString     = fmt.Errorf("parse empty string")
	ErrMaxStrLen       = fmt.Errorf("string input exceeds maximum length %d", maxStrLen)
	ErrInvalidFormat   = fmt.Errorf("invalid format")
	ErrDivideByZero    = fmt.Errorf("can't divide by zero")
	ErrSqrtNegative    = fmt.Errorf("can't calculate square root of negative number")
)

var (
	Zero    = Decimal{}
	One     = MustFromInt64(1, 0)
	OneUint = MustFromUint64(1, 19)
)

// Decimal represents a fixed-point decimal number.
// The number is represented as a coefficient and a scale.
//
// Number = coef / 10^(scale)
type Decimal struct {
	coef  bint
	neg   bool // true if number is negative
	scale uint8
}

func SetDefaultScale(scale uint8) {
	if scale > maxScale {
		panic(fmt.Sprintf("scale out of range. Only allow maximum %d digits after the decimal points", maxScale))
	}

	if scale == 0 {
		panic("scale must be greater than 0")
	}

	defaultScale = scale
}

func NewFromHiLo(neg bool, hi uint64, lo uint64, scale uint8) (Decimal, error) {
	if scale > defaultScale {
		return Decimal{}, ErrScaleOutOfRange
	}

	coef := u128{hi: hi, lo: lo}
	return newDecimal(neg, bintFromU128(coef), scale), nil
}

// newDecimal return the decimal
func newDecimal(neg bool, coef bint, scale uint8) Decimal {
	return Decimal{neg: neg, coef: coef, scale: scale}
}

// NewFromUint64 returns a decimal which equals to coef / 10^scale and coef is an uint64
// Trailing zeros wll be removed and the scale will also be adjusted
func NewFromUint64(coef uint64, scale uint8) (Decimal, error) {
	if scale > defaultScale {
		return Decimal{}, ErrScaleOutOfRange
	}

	return newDecimal(false, bintFromU64(coef), scale), nil
}

// MustFromUint64 similars to NewFromUint64, but panics instead of returning error
func MustFromUint64(coef uint64, scale uint8) Decimal {
	d, err := NewFromUint64(coef, scale)
	if err != nil {
		panic(err)
	}

	return d
}

// NewFromInt64 returns a decimal which equals to coef / 10^scale and coef is an int64
// Trailing zeros wll be removed and the scale will also be adjusted
func NewFromInt64(coef int64, scale uint8) (Decimal, error) {
	var neg bool
	if coef < 0 {
		neg = true
		coef = -coef
	}

	if scale > defaultScale {
		return Decimal{}, ErrScaleOutOfRange
	}

	// nolint: gosec
	return newDecimal(neg, bintFromU64(uint64(coef)), scale), nil
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
		return Decimal{}, fmt.Errorf("%w: can't parse float '%v' to Decimal", ErrInvalidFormat, f)
	}

	s := strconv.FormatFloat(f, 'f', -1, 64)
	d, err := Parse(s)
	if err != nil {
		return Decimal{}, fmt.Errorf("can't parse float: %w", err)
	}

	return d, nil
}

// MustFromFloat64 similars to NewFromFloat64, but panics instead of returning error
func MustFromFloat64(f float64) Decimal {
	d, err := NewFromFloat64(f)
	if err != nil {
		panic(err)
	}

	return d
}

// InexactFloat64 returns the float64 representation of the decimal.
// The result may not be 100% accurate due to the limitation of float64 (less decimal precision).
// !Caution: this method will not return the exact number if the decimal is too large.
// !e.g. 123456789012345678901234567890123456789.9999999999999999999 -> 123456789012345680000000000000000000000
func (d Decimal) InexactFloat64() (float64, error) {
	return strconv.ParseFloat(d.String(), 64)
}

// Parse parses a number in string to Decimal.
// The string must be in the format of: [+-]d{1,19}[.d{1,19}]
// e.g. "123", "-123", "123.456", "-123.456", "+123.456", "0.123".
//
// Returns error if:
//  1. empty/invalid string
//  2. the number has more than 19 digits after the decimal point
func Parse(s string) (Decimal, error) {
	if len(s) == 0 {
		return Decimal{}, ErrEmptyString
	}

	neg, bint, scale, err := parseBint(s)
	if err != nil {
		return Decimal{}, err
	}

	return newDecimal(neg, bint, scale), nil
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
func (d Decimal) Add(e Decimal) Decimal {
	dcoef, ecoef := d.coef, e.coef

	var (
		scale uint8
	)

	switch {
	case d.scale == e.scale:
		scale = d.scale
	case d.scale > e.scale:
		scale = d.scale
		ecoef = ecoef.Mul(bintFromU128(pow10[d.scale-e.scale]))
	case d.scale < e.scale:
		scale = e.scale
		dcoef = dcoef.Mul(bintFromU128(pow10[e.scale-d.scale]))
	}

	if d.neg == e.neg {
		return newDecimal(d.neg, dcoef.Add(ecoef), scale)
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
func (d Decimal) Add64(e uint64) Decimal {
	ecoef := bintFromU64(e).Mul(bintFromU128(pow10[d.scale]))

	if d.neg {
		var (
			dcoef bint
			neg   bool
		)

		if d.coef.GT(ecoef) {
			// can ignore the error as we already check if dcoef > ecoef
			dcoef, _ = d.coef.Sub(ecoef)
			neg = true
		} else {
			dcoef, _ = ecoef.Sub(d.coef)
			neg = false
		}

		return newDecimal(neg, dcoef, d.scale)
	}

	dcoef := d.coef.Add(ecoef)
	return newDecimal(false, dcoef, d.scale)
}

// Sub returns d - e
func (d Decimal) Sub(e Decimal) Decimal {
	dcoef, ecoef := d.coef, e.coef

	var (
		scale uint8
	)

	switch {
	case d.scale == e.scale:
		scale = d.scale
	case d.scale > e.scale:
		scale = d.scale
		ecoef = ecoef.Mul(bintFromU128(pow10[d.scale-e.scale]))
	case d.scale < e.scale:
		scale = e.scale
		dcoef = dcoef.Mul(bintFromU128(pow10[e.scale-d.scale]))
	}

	if d.neg != e.neg {
		// different sign
		coef := dcoef.Add(ecoef)
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
func (d Decimal) Sub64(e uint64) Decimal {
	ecoef := bintFromU64(e).Mul(bintFromU128(pow10[d.scale]))

	if !d.neg {
		var (
			dcoef bint
			neg   bool
		)

		if d.coef.GT(ecoef) {
			// dcoef > ecoef, subtract can't overflow
			dcoef, _ = d.coef.Sub(ecoef)
			neg = false
		} else {
			dcoef, _ = ecoef.Sub(d.coef)
			neg = true
		}

		return newDecimal(neg, dcoef, d.scale)
	}

	return newDecimal(true, d.coef.Add(ecoef), d.scale)
}

// Mul returns d * e.
// The result will have at most defaultScale digits after the decimal point.
func (d Decimal) Mul(e Decimal) Decimal {
	if e.coef.IsZero() {
		return Decimal{}
	}

	scale := d.scale + e.scale
	neg := d.neg != e.neg

	v, err := tryMulU128(d, e, neg, scale)
	if err == nil {
		return v
	}

	// overflow, try with *big.Int
	dBig := d.coef.GetBig()
	eBig := e.coef.GetBig()

	dBig = dBig.Mul(dBig, eBig)
	if scale <= defaultScale {
		return newDecimal(neg, bintFromBigInt(dBig), scale)
	}

	q, _ := new(big.Int).QuoRem(dBig, pow10[scale-defaultScale].ToBigInt(), new(big.Int))
	return newDecimal(neg, bintFromBigInt(q), defaultScale)
}

func tryMulU128(d, e Decimal, neg bool, scale uint8) (Decimal, error) {
	if d.coef.overflow || e.coef.overflow {
		return Decimal{}, ErrOverflow
	}

	rcoef := d.coef.u128.MulToU256(e.coef.u128)
	if scale <= defaultScale {
		if !rcoef.carry.IsZero() {
			return Decimal{}, ErrOverflow
		}

		coef := u128{hi: rcoef.hi, lo: rcoef.lo}

		return newDecimal(neg, bintFromU128(coef), scale), nil
	}

	q, err := rcoef.fastQuo(pow10[scale-defaultScale])
	if err != nil {
		return Decimal{}, err
	}

	return newDecimal(neg, bintFromU128(q), defaultScale), nil
}

// Mul64 returns d * e where e is a uint64.
// The result will have at most defaultScale digits after the decimal point.
func (d Decimal) Mul64(v uint64) Decimal {
	if v == 0 {
		return Decimal{}
	}

	if v == 1 {
		return d
	}

	if !d.coef.overflow {
		coef, err := d.coef.u128.Mul64(v)
		if err == nil {
			return newDecimal(d.neg, bintFromU128(coef), d.scale)
		}
	}

	// overflow, try with *big.Int
	dBig := d.coef.GetBig()
	dBig = dBig.Mul(dBig, new(big.Int).SetUint64(v))

	return newDecimal(d.neg, bintFromBigInt(dBig), d.scale)
}

// Div returns d / e
// If the result has more than 19 fraction digits, it will be truncated to 19 digits.
// Returns divide by zero error when e is zero
func (d Decimal) Div(e Decimal) (Decimal, error) {
	if e.coef.IsZero() {
		return Decimal{}, ErrDivideByZero
	}

	neg := d.neg != e.neg

	q, err := tryDivU128(d, e, neg)
	if err == nil {
		return q, nil
	}

	// Need to multiply divident with factor
	// to make sure the total decimal number after the decimal point is defaultScale
	factor := defaultScale - (d.scale - e.scale)

	// overflow, try with *big.Int
	dBig := d.coef.GetBig()
	eBig := e.coef.GetBig()

	dBig = dBig.Mul(dBig, pow10[factor].ToBigInt())
	dBig = dBig.Div(dBig, eBig)
	return newDecimal(neg, bintFromBigInt(dBig), defaultScale), nil
}

// DivExact returns d / e with exact scale. The result will be truncated to the specified scale. No rounding is performed.
func (d Decimal) DivExact(e Decimal, scale uint8) (Decimal, error) {
	q, err := d.Div(e)
	if err != nil {
		return Decimal{}, err
	}

	return q.Trunc(scale), nil
}

func tryDivU128(d, e Decimal, neg bool) (Decimal, error) {
	if d.coef.overflow || e.coef.overflow {
		return Decimal{}, ErrOverflow
	}

	// Need to multiply divident with factor
	// to make sure the total decimal number after the decimal point is defaultScale
	factor := defaultScale - (d.scale - e.scale)

	d256 := d.coef.u128.MulToU256(pow10[factor])
	quo, err := d256.fastQuo(e.coef.u128)
	if err != nil {
		return Decimal{}, err
	}

	return newDecimal(neg, bintFromU128(quo), defaultScale), nil
}

// Div64 returns d / e where e is a uint64
// If the result has more than 19 fraction digits, it will be truncated to 19 digits.
// Returns divide by zero error when e is zero
func (d Decimal) Div64(v uint64) (Decimal, error) {
	if v == 0 {
		return Decimal{}, ErrDivideByZero
	}

	if v == 1 {
		return d, nil
	}

	if !d.coef.overflow {
		d256 := d.coef.u128.MulToU256(pow10[defaultScale-d.scale])
		quo, _, err := d256.quoRem64Tou128(v)
		if err == nil {
			return newDecimal(d.neg, bintFromU128(quo), defaultScale), nil
		}
	}

	// overflow, try with *big.Int
	dBig := d.coef.GetBig()
	dBig.Mul(dBig, pow10[defaultScale-d.scale].ToBigInt())
	dBig.Div(dBig, new(big.Int).SetUint64(v))

	return newDecimal(d.neg, bintFromBigInt(dBig), defaultScale), nil
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
		return -d.cmpDecSameSign(e)
	}

	return d.cmpDecSameSign(e)
}

func (d Decimal) cmpDecSameSign(e Decimal) int {
	result, err := tryCmpU128(d, e)
	if err == nil {
		return result
	}

	// overflow, fallback to big.Int
	dBig := d.coef.GetBig()
	eBig := e.coef.GetBig()

	if d.scale == e.scale {
		return dBig.Cmp(eBig)
	}

	if d.scale < e.scale {
		dBig = dBig.Mul(dBig, pow10[e.scale-d.scale].ToBigInt())
	} else {
		eBig = eBig.Mul(eBig, pow10[d.scale-e.scale].ToBigInt())
	}

	return dBig.Cmp(eBig)
}

func tryCmpU128(d, e Decimal) (int, error) {
	if d.coef.overflow || e.coef.overflow {
		return 0, ErrOverflow
	}

	if d.scale == e.scale {
		return d.coef.u128.Cmp(e.coef.u128), nil
	}

	// scale is different
	// e has more fraction digits
	if d.scale < e.scale {
		// d has more fraction digits
		d256 := d.coef.u128.MulToU256(pow10[e.scale-d.scale])
		return d256.cmp128(e.coef.u128), nil
	}

	// d has more fraction digits
	// we need to compare d with e * 10^(d.scale - e.scale)
	e256 := e.coef.u128.MulToU256(pow10[d.scale-e.scale])

	// remember to reverse the result because e256.cmp128(d.coef) returns the opposite
	return -e256.cmp128(d.coef.u128), nil
}

// Rescale returns the decimal with the new scale only if the new scale is greater than the current scale.
// Useful when you want to increase the scale of the decimal.
//
// Example:
//
//	d := MustParse("123.456")
//	d.Rescale(5) // 123.45600
func (d Decimal) Rescale(scale uint8) Decimal {
	if scale > maxScale {
		scale = maxScale
	}

	if scale <= d.scale {
		return d
	}

	diff := scale - d.scale
	coef := d.coef.Mul(bintFromU128(pow10[diff]))

	return newDecimal(d.neg, coef, scale)
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

// RoundBank uses half up to even (banker's rounding) to round the decimal to the specified scale.
//
// Examples:
//
//	Round(1.12345, 4) = 1.1234
//	Round(1.12335, 4) = 1.1234
//	Round(1.5, 0) = 2
//	Roung(-1.5, 0) = -2
func (d Decimal) RoundBank(scale uint8) Decimal {
	if scale >= d.scale {
		return d
	}

	factor := pow10[d.scale-scale]
	lo := factor.lo / 2

	if !d.coef.overflow {
		var err error
		q, r := d.coef.u128.QuoRem64(factor.lo)
		if lo < r || (lo == r && q.lo%2 == 1) {
			q, err = q.Add64(1)
		}

		// no overflow, return the result
		if err == nil {
			return newDecimal(d.neg, bintFromU128(q), scale)
		}
	}

	// overflow, fallback to big.Int
	dBig := d.coef.GetBig()
	q, r := new(big.Int).QuoRem(dBig, factor.ToBigInt(), new(big.Int))

	loBig := new(big.Int).SetUint64(lo)
	if r.Cmp(loBig) > 0 || (r.Cmp(loBig) == 0 && q.Bit(0) == 1) {
		q.Add(q, bigOne)
	}

	return newDecimal(d.neg, bintFromBigInt(q), scale)
}

// RoundHAZ rounds the decimal to the specified scale using HALF AWAY FROM ZERO method (https://en.wikipedia.org/wiki/Rounding#Rounding_half_away_from_zero).
//
// Examples:
//
//	Round(1.12345, 4) = 1.1235
//	Round(1.12335, 4) = 1.1234
//	Round(1.5, 0) = 2
//	Round(-1.5, 0) = -2
func (d Decimal) RoundHAZ(scale uint8) Decimal {
	if scale >= d.scale {
		return d
	}

	factor := pow10[d.scale-scale]
	lo, _ := factor.QuoRem64(2)

	if !d.coef.overflow {
		var err error
		q, r := d.coef.u128.QuoRem64(factor.lo)
		if lo.Cmp64(r) <= 0 {
			q, err = q.Add64(1)
		}

		if err == nil {
			return newDecimal(d.neg, bintFromU128(q), scale)
		}
	}

	// overflow, fallback to big.Int
	dBig := d.coef.GetBig()
	q, r := new(big.Int).QuoRem(dBig, factor.ToBigInt(), new(big.Int))

	loBig := lo.ToBigInt()
	if r.Cmp(loBig) >= 0 {
		q.Add(q, bigOne)
	}

	return newDecimal(d.neg, bintFromBigInt(q), scale)
}

// RoundHTZ rounds the decimal to the specified scale using HALF TOWARD ZERO method (https://en.wikipedia.org/wiki/Rounding#Rounding_half_toward_zero).
//
// Examples:
//
//	Round(1.12345, 4) = 1.1234
//	Round(1.12335, 4) = 1.1233
//	Round(1.5, 0) = 1
//	Round(-1.5, 0) = -1
func (d Decimal) RoundHTZ(scale uint8) Decimal {
	if scale >= d.scale {
		return d
	}

	factor := pow10[d.scale-scale]
	lo, _ := factor.QuoRem64(2)

	if !d.coef.overflow {
		var err error
		q, r := d.coef.u128.QuoRem64(factor.lo)
		if lo.Cmp64(r) < 0 {
			q, err = q.Add64(1)
		}

		if err == nil {
			return newDecimal(d.neg, bintFromU128(q), scale)
		}
	}

	// overflow, fallback to big.Int
	dBig := d.coef.GetBig()
	q, r := new(big.Int).QuoRem(dBig, factor.ToBigInt(), new(big.Int))

	loBig := lo.ToBigInt()
	if r.Cmp(loBig) > 0 {
		q.Add(q, bigOne)
	}

	return newDecimal(d.neg, bintFromBigInt(q), scale)
}

// Floor returns the largest integer value less than or equal to d.
//
// Examples:
//
//	Floor(1.12345) = 1
//	Floor(1.12335) = 1
//	Floor(1.5, 0) = 1
//	Floor(-1.5, 0) = -2
func (d Decimal) Floor() Decimal {
	if d.scale == 0 {
		return d
	}

	if !d.coef.overflow {
		var err error
		q, r := d.coef.u128.QuoRem64(pow10[d.scale].lo)

		// add 1 if it's negative and there's a remainder, e.g. -1.5 -> -2
		if d.neg && r != 0 {
			q, err = q.Add64(1)
		}

		if err == nil {
			return newDecimal(d.neg, bintFromU128(q), 0)
		}
	}

	// overflow, fallback to big.Int
	dBig := d.coef.GetBig()
	q, r := new(big.Int).QuoRem(dBig, pow10[d.scale].ToBigInt(), new(big.Int))

	// add 1 if it's negative and there's a remainder, e.g. -1.5 -> -2
	if d.neg && r.Cmp(bigZero) != 0 {
		q.Add(q, bigOne)
	}

	return newDecimal(d.neg, bintFromBigInt(q), 0)
}

// Ceil returns the smallest integer value greater than or equal to d.
//
// Examples:
//
//	Ceil(1.12345, 4) = 1.1235
//	Ceil(1.12335, 4) = 1.1234
//	Ceil(1.5, 0) = 2
//	Ceil(-1.5, 0) = -1
func (d Decimal) Ceil() Decimal {
	if d.scale == 0 {
		return d
	}

	if !d.coef.overflow {
		var err error
		q, r := d.coef.u128.QuoRem64(pow10[d.scale].lo)

		// add 1 if it's positive and there's a remainder, e.g. 1.5 -> 2
		if !d.neg && r != 0 {
			q, err = q.Add64(1)
		}

		if err == nil {
			return newDecimal(d.neg, bintFromU128(q), 0)
		}
	}

	// overflow, fallback to big.Int
	dBig := d.coef.GetBig()
	q, r := new(big.Int).QuoRem(dBig, pow10[d.scale].ToBigInt(), new(big.Int))

	// add 1 if it's positive and there's a remainder, e.g. 1.5 -> 2
	if !d.neg && r.Cmp(bigZero) != 0 {
		q.Add(q, bigOne)
	}

	return newDecimal(d.neg, bintFromBigInt(q), 0)
}

// Trunc returns d after truncating the decimal to the specified scale.
//
// Examples:
//
//	Trunc(1.12345, 4) = 1.1234
//	Trunc(1.12335, 4) = 1.1233
func (d Decimal) Trunc(scale uint8) Decimal {
	if scale >= d.scale {
		return d
	}

	factor := pow10[d.scale-scale]

	if !d.coef.overflow {
		q, _ := d.coef.u128.QuoRem64(factor.lo)
		return newDecimal(d.neg, bintFromU128(q), scale)
	}

	// overflow, fallback to big.Int
	dBig := d.coef.GetBig()
	q := new(big.Int).Quo(dBig, factor.ToBigInt())
	return newDecimal(d.neg, bintFromBigInt(q), scale)
}

func (d Decimal) trimTrailingZeros() Decimal {
	if d.coef.overflow {
		zeros := trailingZerosBigInt(d.coef.bigInt)

		var (
			dBig  = d.coef.GetBig()
			scale uint8
		)

		if zeros == 0 {
			return newDecimal(d.neg, bintFromBigInt(dBig), d.scale)
		}

		if zeros >= d.scale {
			dBig.Div(dBig, pow10[d.scale].ToBigInt())
			scale = 0
		} else {
			scale = d.scale - uint8(zeros)
			dBig.Div(dBig, pow10[zeros].ToBigInt())
		}

		return newDecimal(d.neg, bintFromBigInt(dBig), scale)
	}

	zeros := trailingZerosU128(d.coef.u128)
	if zeros == 0 {
		return newDecimal(d.neg, bintFromU128(d.coef.u128), d.scale)
	}

	var (
		coef  u128
		scale uint8
	)

	if zeros >= d.scale {
		coef, _, _ = d.coef.u128.QuoRem(pow10[d.scale])
		scale = 0
	} else {
		scale = d.scale - zeros
		coef, _, _ = d.coef.u128.QuoRem(pow10[zeros])
	}

	d.coef = bintFromU128(coef)
	d.scale = scale
	return d
}

func trailingZerosBigInt(n *big.Int) uint8 {
	var (
		zeros uint8
		z, m  = new(big.Int), new(big.Int)
	)

	_, m = z.QuoRem(n, pow10Big[16], m)
	if m.Cmp(bigZero) == 0 {
		zeros += 16

		// shortcut because maxScale = 19
		_, m = z.QuoRem(n, pow10Big[zeros+2], m)
		if m.Cmp(bigZero) == 0 {
			zeros += 2
		}

		_, m = z.QuoRem(n, pow10Big[zeros+1], m)
		if m.Cmp(bigZero) == 0 {
			zeros++
		}

		return zeros
	}

	_, m = z.QuoRem(n, pow10Big[8], m)
	if m.Cmp(bigZero) == 0 {
		zeros += 8
	}

	_, m = z.QuoRem(n, pow10Big[zeros+4], m)
	if m.Cmp(bigZero) == 0 {
		zeros += 4
	}

	_, m = z.QuoRem(n, pow10Big[zeros+2], m)
	if m.Cmp(bigZero) == 0 {
		zeros += 2
	}

	_, m = z.QuoRem(n, pow10Big[zeros+1], m)
	if m.Cmp(bigZero) == 0 {
		zeros++
	}

	return zeros
}

func trailingZerosU128(n u128) uint8 {
	var zeros uint8

	_, rem := n.QuoRem64(1e16)
	if rem == 0 {
		zeros += 16

		_, rem = n.QuoRem64(pow10[zeros+2].lo)
		if rem == 0 {
			zeros += 2
		}

		_, rem = n.QuoRem64(pow10[zeros+1].lo)
		if rem == 0 {
			zeros++
		}

		return zeros
	}

	_, rem = n.QuoRem64(1e8)
	if rem == 0 {
		zeros += 8
	}

	_, rem = n.QuoRem64(pow10[zeros+4].lo)
	if rem == 0 {
		zeros += 4
	}

	_, rem = n.QuoRem64(pow10[zeros+2].lo)
	if rem == 0 {
		zeros += 2
	}

	_, rem = n.QuoRem64(pow10[zeros+1].lo)
	if rem == 0 {
		zeros++
	}

	return zeros
}

// PowInt returns d^e where e is an integer.
//
// Examples:
//
//	PowInt(2.5, 2) = 6.25
//	PowInt(2.5, -2) = 0.16
func (d Decimal) PowInt(e int) Decimal {
	// check 0 first to avoid 0^0 = 1
	if d.coef.IsZero() {
		return Zero
	}

	if e == 0 {
		return One
	}

	if e == 1 {
		return d
	}

	// rescale fist to remove trailing zeros
	dTrim := d.trimTrailingZeros()

	if e < 0 {
		return dTrim.powIntInverse(-e)
	}

	// e > 1 && d != 0
	q, err := dTrim.tryPowIntU128(e)
	if err == nil {
		return q
	}

	// overflow, fallback to big.Int
	dBig := dTrim.coef.GetBig()
	factor := 0
	powScale := int(dTrim.scale) * e
	if powScale >= int(defaultScale) {
		factor = powScale - int(defaultScale)
		powScale = int(defaultScale)
	}

	m := new(big.Int).Exp(bigTen, big.NewInt(int64(factor)), nil)
	dBig = new(big.Int).Exp(dBig, big.NewInt(int64(e)), nil)
	qBig := dBig.Quo(dBig, m)

	neg := d.neg
	if e%2 == 0 {
		neg = false
	}

	// nolint: gosec
	return newDecimal(neg, bintFromBigInt(qBig), uint8(powScale))
}

// powIntInverse returns d^(-e), with e > 0
func (d Decimal) powIntInverse(e int) Decimal {
	q, err := d.tryInversePowIntU128(e)
	if err == nil {
		return q
	}

	// overflow, fallback to big.Int
	dBig := d.coef.GetBig()
	powScale := int(d.scale) * e

	// d^(-e) = 10^(defaultScale + e) / d^e (with defaultScale digits after the decimal point)
	m := new(big.Int).Exp(bigTen, big.NewInt(int64(powScale+int(defaultScale))), nil)
	dBig = new(big.Int).Exp(dBig, big.NewInt(int64(e)), nil)
	qBig := dBig.Quo(m, dBig)

	neg := d.neg
	if e%2 == 0 {
		neg = false
	}

	return newDecimal(neg, bintFromBigInt(qBig), defaultScale)
}

func (d Decimal) tryPowIntU128(e int) (Decimal, error) {
	if d.coef.overflow {
		return Decimal{}, ErrOverflow
	}

	if d.coef.u128.hi != 0 && e >= 3 {
		// e > 3 and u128.hi != 0 means the result will >= 2^192,
		// which we can't use fast division. So we need to use big.Int instead
		return Decimal{}, ErrOverflow
	}

	neg := d.neg
	if e%2 == 0 {
		neg = false
	}

	powScale := int(d.scale) * e
	if powScale > int(defaultScale)+38 {
		return Decimal{}, ErrOverflow
	}

	factor := 0
	if powScale > int(defaultScale) {
		factor = powScale - int(defaultScale)
		powScale = int(defaultScale)
	}

	d256 := u256{lo: d.coef.u128.lo, hi: d.coef.u128.hi}
	result, err := d256.pow(e)
	if err != nil {
		return Decimal{}, err
	}

	if factor == 0 {
		if !result.carry.IsZero() {
			return Decimal{}, ErrOverflow
		}

		// nolint: gosec
		return newDecimal(neg, bintFromU128(u128{hi: result.hi, lo: result.lo}), uint8(powScale)), nil
	}

	if result.carry.hi != 0 {
		return Decimal{}, ErrOverflow
	}

	q, err := result.fastQuo(pow10[factor]) // it's safe to use pow10[factor] as factor <= 38
	if err != nil {
		return Decimal{}, err
	}

	return newDecimal(neg, bintFromU128(q), defaultScale), nil
}

func (d Decimal) tryInversePowIntU128(e int) (Decimal, error) {
	if d.coef.overflow {
		return Decimal{}, ErrOverflow
	}

	if d.coef.u128.hi != 0 && e >= 3 {
		// e > 3 and u128.hi != 0 means the result will >= 2^192,
		// which we can't use fast division. So we need to use big.Int instead
		return Decimal{}, ErrOverflow
	}

	neg := d.neg
	if e%2 == 0 {
		neg = false
	}

	powScale := int(d.scale) * e
	if powScale > int(defaultScale)+38 {
		return Decimal{}, ErrOverflow
	}

	factor := 0
	if powScale > int(defaultScale) {
		factor = powScale - int(defaultScale)
		powScale = int(defaultScale)
	}

	d256 := u256{lo: d.coef.u128.lo, hi: d.coef.u128.hi}
	result, err := d256.pow(e)
	if err != nil {
		return Decimal{}, err
	}

	if factor == 0 {
		if !result.carry.IsZero() {
			return Decimal{}, ErrOverflow
		}

		// nolint: gosec
		a256 := one128.MulToU256(pow10[defaultScale+uint8(powScale)])

		q, err := a256.fastQuo(u128{hi: result.hi, lo: result.lo})
		if err != nil {
			return Decimal{}, err
		}

		return newDecimal(neg, bintFromU128(q), defaultScale), nil
	}

	// if result is not u128, one solution is adjusting it to u128 by dividing it with 10^factor
	// in some cases, this adjustment creates a big difference in the final result
	// so to be safe, use big.Int instead
	if !result.carry.IsZero() {
		return Decimal{}, ErrOverflow
	}

	// a256 = 10^(powScale + factor + defaultScale)
	// nolint: gosec
	a256 := pow10[factor].MulToU256(pow10[defaultScale+uint8(powScale)])
	q, err := a256.fastQuo(u128{hi: result.hi, lo: result.lo})
	if err != nil {
		return Decimal{}, err
	}

	return newDecimal(neg, bintFromU128(q), defaultScale), nil
}

// Sqrt returns the square root of d using Newton-Raphson method.
// The result will have at most defaultScale digits after the decimal point.
// Returns error if d < 0
//
// Examples:
//
//	Sqrt(4) = 2
//	Sqrt(2) = 1.4142135623730950488
func (d Decimal) Sqrt() (Decimal, error) {
	if d.neg {
		return Decimal{}, ErrSqrtNegative
	}

	if d.coef.IsZero() {
		return Zero, nil
	}

	if d.Cmp(One) == 0 {
		return One, nil
	}

	if !d.coef.overflow {
		q, err := d.sqrtU128()
		if err == nil {
			return q, nil
		}
	}

	// overflow, fallback to big.Int
	dBig := d.coef.GetBig()
	factor := 2*defaultScale - d.scale
	coef := dBig.Mul(dBig, pow10[factor].ToBigInt())
	return newDecimal(false, bintFromBigInt(coef.Sqrt(coef)), defaultScale), nil
}

func (d Decimal) sqrtU128() (Decimal, error) {
	factor := 2*defaultScale - d.scale

	coef := d.coef.u128.MulToU256(pow10[factor])
	if coef.carry.hi != 0 {
		return Decimal{}, ErrOverflow
	}

	// nolint: gosec
	bitLen := uint(coef.bitLen()) // bitLen < 192

	// initial guess = 2^((bitLen + 1) / 2) ≥ √coef
	x := one128.Lsh((bitLen + 1) / 2)

	// Newton-Raphson method
	for {
		// calculate x1 = (x + coef/x) / 2
		y, err := coef.fastQuo(x)
		if err != nil {
			return Decimal{}, err
		}

		x1, err := x.Add(y)
		if err != nil {
			return Decimal{}, err
		}

		x1 = x1.Rsh(1)
		if x1.Cmp(x) == 0 {
			break
		}

		x = x1
	}

	return newDecimal(false, bintFromU128(x), defaultScale), nil
}
