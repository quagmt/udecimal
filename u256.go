package udecimal

import (
	"math/bits"
)

// U256 represents a 256-bits unsigned integer
// U256 = carry * 2^128 + hi*2^64 + lo
// carry = u*2^64 + v
type U256 struct {
	hi, lo uint64

	// store overflow
	carry bint
}

// for debugging
// func (u U256) PrintBit() string {
// 	b1 := strconv.FormatUint(u.carry.hi, 2)
// 	b2 := strconv.FormatUint(u.carry.lo, 2)
// 	b3 := strconv.FormatUint(u.hi, 2)
// 	b4 := strconv.FormatUint(u.lo, 2)

// 	return fmt.Sprintf("%s.%s.%s.%s", apz(b1), apz(b2), apz(b3), apz(b4))
// }

// func apz(s string) string {
// 	for range 64 - len(s) {
// 		s = "0" + s
// 	}

// 	return s
// }

// Compare 2 U256, returns:
//
//	+1 when u > v
//	 0 when u = v
//	-1 when u < v
func (u U256) cmp(v U256) int {
	if k := u.carry.Cmp(v.carry); k != 0 {
		return k
	}

	return bintFromHiLo(u.hi, u.lo).Cmp(bintFromHiLo(v.hi, v.lo))
}

// Compare U256 and U128, returns:
//
//	+1 when u > v
//	 0 when u = v
//	-1 when u < v
func (u U256) Cmp128(v bint) int {
	if !u.carry.IsZero() {
		return 1
	}

	return bintFromHiLo(u.hi, u.lo).Cmp(v)
}

func (u U256) Sub(v U256) (U256, error) {
	lo, borrow := bits.Sub64(u.lo, v.lo, 0)
	hi, borrow := bits.Sub64(u.hi, v.hi, borrow)

	c, err := v.carry.Add64(borrow)
	if err != nil {
		return U256{}, err
	}

	c1, err := u.carry.Sub(c)
	if err != nil {
		return U256{}, err
	}

	return U256{lo: lo, hi: hi, carry: c1}, nil
}

func (u U256) Rsh(n uint) (v U256) {
	switch {
	case n < 64:
		v.carry = u.carry.Rsh(n)
		v.hi = u.carry.lo<<(64-n) | u.hi>>n
		v.lo = u.hi<<(64-n) | u.lo>>n

	case 64 <= n && n < 128:
		v.carry.hi = 0
		v.carry.lo = u.carry.hi >> (n - 64)
		v.hi = u.carry.hi<<(128-n) | u.carry.lo>>(n-64)
		v.lo = u.carry.lo<<(128-n) | u.hi>>(n-64)

	case n >= 128:
		v.carry = bint{}
		v.hi = u.carry.hi >> (n - 128)
		v.lo = u.carry.hi<<(196-n) | u.carry.lo>>(n-128)

	default:
		// n < 0, can't happen
	}

	return
}

// Quo only returns quotient of u/v
// The implementation follows Hacker's Delight multiword division algorithm
// with some constraints regarding max coef and scale value, including:
//
//	max(coef) = 10^38-1
//	max(scale) = 19
//	max(whole_part) = 10^19-1
func (u U256) quo(v bint) (bint, error) {
	if u.carry.IsZero() {
		q, _, err := bintFromHiLo(u.hi, u.lo).QuoRem(v)
		return q, err
	}

	if v.hi == 0 {
		q, _, err := u.quoRem64ToBint(v.lo)
		return q, err
	}

	// if u >= 2^192, the quotient won't fit in 128-bits number (overflow).
	// Put in both here and inside QuoRem64, in case we call QuoRem64 directly
	if u.carry.hi != 0 {
		return bint{}, ErrOverflow
	}

	// 1 <= n <= 63 (as bint < 10^38)
	n := uint(bits.LeadingZeros64(v.hi))
	v1 := v.Lsh(n)
	u1 := u.Rsh(64 - n)

	// let q are final quotient and remainder and tq = q + k (k >= 0)
	// calculate 'trial quotient' tq (q <= tq < q + 2^64)
	tq, _, err := u1.quoRem64ToBint(v1.hi)
	if err != nil {
		return bint{}, err
	}

	vq := v.MulToU256(tq)

	// Some pre-conditions:
	// We only allow bint to have 38-digits max, then:
	// max(u) = (10^38-1) * 10^19 = 10^58 - 10^19 < 2^190 --> u < 2^190
	//
	// max(v*k) = u * [2^(64-n) - 1]/2^(127-n) (with n is v's leading zeros, 1 <= n <= 63)
	// --> max(v*k) = u / 2^63 < 2^190 / 2^63
	// --> v*k < 2^127
	// vqu = vq - u = (q+k)*v - (q*v + r) = k*v - r
	// with v*k < 2^127 --> vqu < 2^128 and can be represented by a 128-bit uint (no overflow)
	if vq.cmp(u) <= 0 {
		// vq <= u means tq = q
		return tq, nil
	}

	vqu, err := vq.Sub(u)
	if err != nil {
		return bint{}, err
	}

	// techically this can't happen, just put it here to do fuzz test and cross-check with other libs
	if vqu.carry.hi&vqu.carry.lo != 0 {
		return bint{}, ErrOverflow
	}

	// k1 = k - 1
	k1, _, err := bintFromHiLo(vqu.hi, vqu.lo).QuoRem(v)
	if err != nil {
		return bint{}, err
	}

	// adjust the result, with tq = q_final + k = q_final + (k1 + 1) --> q_final = tq - (k1 + 1)
	tq, err = tq.Sub(k1)
	if err != nil {
		return bint{}, err
	}

	tq, err = tq.Sub64(1)
	if err != nil {
		return bint{}, err
	}

	// we don't really need the remainder, might un-comment later if needed
	// r, err := v.Sub(r1)
	// if err != nil {
	// 	return bint{}, bint{}, err
	// }

	return tq, nil
}

// quoRem64ToBint return q,r which:
//
//	q must be a bint
//	u = q*v + r
//	Return overflow if the result q doesn't fit in a bint
func (u U256) quoRem64ToBint(v uint64) (bint, uint64, error) {
	// obvious case that the result won't fit in 128-bits number
	if u.carry.hi != 0 {
		return bint{}, 0, ErrOverflow
	}

	if u.carry.lo == 0 {
		q, r := bintFromHiLo(u.hi, u.lo).QuoRem64(v)
		return q, r, nil
	}

	quo, rem := bintFromHiLo(u.carry.lo, u.hi).QuoRem64(v)
	if quo.hi != 0 {
		return bint{}, 0, ErrOverflow
	}

	hi := quo.lo
	lo, r := bits.Div64(rem, u.lo, v)

	return bintFromHiLo(hi, lo), r, nil
}
