package udecimal

import "math/bits"

// pow10Reciprocal contains precomputed magic numbers for fast division by pow10[k].lo.
// For a uint64 n, n / pow10[k].lo can be computed as:
//
//	hi, _ := bits.Mul64(n, magic)
//	t := ((n - hi) >> 1) + hi
//	q := t >> shift
//
// Based on "Division by Invariant Integers using Multiplication"
// (Granlund & Montgomery, 1994) and Hacker's Delight Chapter 10.
//
// All pow10 values > 1 use the "add-shift" variant because their magic numbers
// require 65 bits (overflow uint64). The algorithm handles this by computing:
//
//	q = (mulhi(n, m_lo) + ((n - mulhi(n, m_lo)) >> 1)) >> (l - 1)
//
// where m_lo is the lower 64 bits of the full 65-bit magic number.
var pow10Reciprocal = [20]struct {
	magic uint64
	shift uint8
}{
	{magic: 0, shift: 0},                     // 10^0 (identity, use k==0 fast path)
	{magic: 11068046444225730970, shift: 3},  // 10^1
	{magic: 5165088340638674453, shift: 6},   // 10^2
	{magic: 442721857769029239, shift: 9},    // 10^3
	{magic: 11776401416656177752, shift: 13}, // 10^4
	{magic: 5731772318583031879, shift: 16},  // 10^5
	{magic: 896069040124515180, shift: 19},   // 10^6
	{magic: 12501756908424955257, shift: 23}, // 10^7
	{magic: 6312056711998053882, shift: 26},  // 10^8
	{magic: 1360296554856532783, shift: 29},  // 10^9
	{magic: 13244520931996183422, shift: 33}, // 10^10
	{magic: 6906267930855036414, shift: 36},  // 10^11
	{magic: 1835665529942118808, shift: 39},  // 10^12
	{magic: 14005111292133121063, shift: 43}, // 10^13
	{magic: 7514740218964586527, shift: 46},  // 10^14
	{magic: 2322443360429758899, shift: 49},  // 10^15
	{magic: 14783955820913345207, shift: 53}, // 10^16
	{magic: 8137815841988765843, shift: 56},  // 10^17
	{magic: 2820903858849102351, shift: 59},  // 10^18
	{magic: 15581492618384294731, shift: 63}, // 10^19
}

// quoRem64Pow10 computes u / pow10[k] returning quotient and remainder.
// For the common case (u.hi == 0), it uses reciprocal multiplication
// which avoids the expensive hardware DIV instruction (~30-90 cycles → ~4 cycles).
// For the rare u.hi != 0 case, it falls back to the standard QuoRem64.
//
// k must be in [0, 19]. This is always true in the library since k is derived
// from prec values (uint8 bounded to 0-19 at construction).
func (u u128) quoRem64Pow10(k uint8) (q u128, r uint64) {
	if k == 0 {
		return u, 0
	}

	// Fast path: when hi == 0, the value fits in a uint64 and we can use
	// reciprocal multiplication (MUL + SHR) instead of DIV.
	if u.hi == 0 {
		rec := pow10Reciprocal[k]
		hi, _ := bits.Mul64(u.lo, rec.magic)
		t := ((u.lo - hi) >> 1) + hi
		quo := t >> rec.shift
		rem := u.lo - quo*pow10[k].lo
		return u128{lo: quo}, rem
	}

	// Slow path: u.hi != 0, use standard 2-limb division.
	return u.QuoRem64(pow10[k].lo)
}

// quo64Pow10 computes u / pow10[k] returning only the quotient.
// Same optimization as quoRem64Pow10 but skips remainder computation overhead
// for call sites that don't need it (e.g., Trunc, trailing zeros checks).
func (u u128) quo64Pow10(k uint8) u128 {
	if k == 0 {
		return u
	}

	if u.hi == 0 {
		rec := pow10Reciprocal[k]
		hi, _ := bits.Mul64(u.lo, rec.magic)
		t := ((u.lo - hi) >> 1) + hi
		quo := t >> rec.shift
		return u128{lo: quo}
	}

	q, _ := u.QuoRem64(pow10[k].lo)
	return q
}

// rem64Pow10 computes u % pow10[k] returning only the remainder.
// Used for trailing zeros detection where we only need to check if remainder is 0.
func (u u128) rem64Pow10(k uint8) uint64 {
	if k == 0 {
		return 0
	}

	if u.hi == 0 {
		rec := pow10Reciprocal[k]
		hi, _ := bits.Mul64(u.lo, rec.magic)
		t := ((u.lo - hi) >> 1) + hi
		quo := t >> rec.shift
		return u.lo - quo*pow10[k].lo
	}

	_, r := u.QuoRem64(pow10[k].lo)
	return r
}
