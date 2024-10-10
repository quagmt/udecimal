package udecimal

import (
	"fmt"
	"math/bits"
)

// u256 represents a 256-bits unsigned integer
// u256 = carry * 2^128 + hi*2^64 + lo
// carry = u*2^64 + v
type u256 struct {
	hi, lo uint64

	// store overflow
	carry u128
}

func (u u256) bitLen() int {
	if u.carry.hi != 0 {
		return 192 + bits.Len64(u.carry.hi)
	}

	if u.carry.lo != 0 {
		return 128 + bits.Len64(u.carry.lo)
	}

	if u.hi != 0 {
		return 64 + bits.Len64(u.hi)
	}

	return bits.Len64(u.lo)
}

// for debugging
// func (u u256) PrintBit() {
// 	b1 := strconv.FormatUint(u.carry.hi, 2)
// 	b2 := strconv.FormatUint(u.carry.lo, 2)
// 	b3 := strconv.FormatUint(u.hi, 2)
// 	b4 := strconv.FormatUint(u.lo, 2)

// 	fmt.Printf("%s.%s.%s.%s\n", apz(b1), apz(b2), apz(b3), apz(b4))
// }

// func apz(s string) string {
// 	if len(s) == 64 {
// 		return s
// 	}

// 	l := len(s)

// 	for range 64 - l {
// 		s = "0" + s
// 	}

// 	return s
// }

// Compare 2 u256, returns:
//
//	+1 when u > v
//	 0 when u = v
//	-1 when u < v
func (u u256) cmp(v u256) int {
	if k := u.carry.Cmp(v.carry); k != 0 {
		return k
	}

	return u128FromHiLo(u.hi, u.lo).Cmp(u128FromHiLo(v.hi, v.lo))
}

// Compare u256 and U128, returns:
//
//	+1 when u > v
//	 0 when u = v
//	-1 when u < v
func (u u256) cmp128(v u128) int {
	if !u.carry.IsZero() {
		return 1
	}

	return u128FromHiLo(u.hi, u.lo).Cmp(v)
}

func (u u256) pow(e int) (u256, error) {
	if e <= 0 {
		return u256{}, fmt.Errorf("invalid exponent %d. Must be greater than 0", e)
	}

	result := u256{lo: 1}
	d256 := u
	var err error

	for ; e > 0; e >>= 1 {
		if e&1 == 1 {
			if !result.carry.IsZero() {
				return u256{}, errOverflow
			}

			// result = result * u (with u = (d256)^(2^i))
			result, err = d256.mul128(u128{lo: result.lo, hi: result.hi})
			if err != nil {
				return u256{}, err
			}
		}

		// d256 = (d256)^2 each time
		d256, err = d256.mul128(u128{lo: d256.lo, hi: d256.hi})
		if err != nil {
			return u256{}, err
		}

		// if there's a carry, next iteration will overflow
		if !d256.carry.IsZero() && e > 1 {
			return u256{}, errOverflow
		}
	}

	return result, nil
}

func (u u256) mul128(v u128) (u256, error) {
	a := u128FromHiLo(u.hi, u.lo).MulToU256(v)
	b, err := u.carry.Mul(v)
	if err != nil {
		return u256{}, err
	}

	c, err := a.carry.Add(b)
	if err != nil {
		return u256{}, err
	}

	return u256{hi: a.hi, lo: a.lo, carry: c}, nil
}

func (u u256) sub(v u256) (u256, error) {
	lo, borrow := bits.Sub64(u.lo, v.lo, 0)
	hi, borrow := bits.Sub64(u.hi, v.hi, borrow)

	c, err := v.carry.Add64(borrow)
	if err != nil {
		return u256{}, err
	}

	c1, err := u.carry.Sub(c)
	if err != nil {
		return u256{}, err
	}

	return u256{lo: lo, hi: hi, carry: c1}, nil
}

func (u u256) rsh(n uint) (v u256) {
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
		v.carry = u128{}
		c := u128{hi: u.carry.hi, lo: u.carry.lo}.Rsh(n - 128)
		v.hi, v.lo = c.hi, c.lo
	default:
		// n < 0, can't happen
	}

	return
}

// Quo only returns quotient of u/v
// Fast divsion for U192 divided by U128 using Hacker's Delight multiword division algorithm
// with some constraints regarding max coef and scale value, including:
//
//	max(coef) = 2^128-1
//	max(scale) = 19
//	max(u) = 2^192-1
func (u u256) fastQuo(v u128) (u128, error) {
	// if u >= 2^192, the quotient might won't fit in 128-bits number (overflow).
	if u.carry.hi != 0 {
		return u128{}, errOverflow
	}

	if u.carry.IsZero() {
		q, _, err := u128FromHiLo(u.hi, u.lo).QuoRem(v)
		return q, err
	}

	if v.hi == 0 {
		q, _, err := u.quoRem64Tou128(v.lo)
		return q, err
	}

	// Let q be the final quotient and tq be the 'trial quotient'
	// The trial quotient tq is defined as tq = q + k, where k is a correction factor.
	// Here's how we determine k using the following steps:

	// 1. Compute the trial quotient tq as tq = [u1 / v1], where:
	//    - u1 = [u / 2^(64 - shift)]
	//    - v1 = [v / 2^(64 - shift)]
	//    (This follows the Hacker's Delight multiword division algorithm)
	// 2. Calculate vq = v * tq
	// 3. If vq <= u, then q = tq
	// 4. If vq > u, compute the difference vqu = vq - u, which can be expressed as:
	//    vqu = v * (q + k) - (vq + r) = v * k - r = v * (k1 + 1) - r = v * k1 + v - r
	// 5. Determine k1 as k1 = [vqu / v] and adjust k1 if necessary
	//
	// However, given that u < 2^192 and v < 2^128, and 0 <= k <= 2^64, it's possible for v * k to exceed 2^128,
	// causing an overflow in vqu.
	// To mitigate this, we need to find the minimum k.
	// If even the minimum k leads to v * k > 2^128, we fall back to big.Int division
	// due to the lack of a fast algorithm for dividing U192 by U128.
	//
	// tq = [u1 / v1] = [u / (v - rem(v, 2^(64 - n))]
	// The minimum k is achieved when tq is minimized, which happens when rem(v, 2^(64 - n)) is minimized,
	// --> 2^(64 - n) being minimized --> n should be maximized.
	// n is the number of leading zeros in v, so 0 <= n <= 63.
	// Technically, using n = 63 provides the optimal k.

	// However, when n = 63:
	// - u1 = [u / 2^(64 - 63)] = u >> 1
	// - v1 = [v / 2^(64 - 63)] = v >> 1
	// And if u1 is U192 and v1 is U128, we cannot find tq = [u1 / v1],
	// since there's no U192/U128 division algorithm currently available.
	//
	// What we do have are fast algorithms for U192/U64 or U128/U128 division.
	// Therefore, we can only compute tq by adjusting u and v to fit either U128/U128 or U192/U64.
	// This might not be the best optimization, but it's the best we can achieve for now.
	// If we later find a fast U192/U128 division algorithm, we can improve this process.
	//
	// As previously mentioned, if after finding the minimum k, v * k still exceeds 2^128, we will fall back to big.Int division.

	// nolint: gosec
	n := uint(bits.LeadingZeros64(v.hi))

	// nolint: gosec
	m := uint(bits.LeadingZeros64(u.carry.lo))

	var (
		v1, tq u128
		u1     u256
		err    error
	)

	if n >= m {
		v1 = v.Lsh(n)
		u1 = u.rsh(64 - n)
		tq, _, err = u1.quoRem64Tou128(v1.hi)
		if err != nil {
			return u128{}, err
		}
	} else {
		// n < m
		v1 = v.Rsh(64 - m)
		u1 = u.rsh(64 - m)

		tq, _, err = u128FromHiLo(u1.hi, u1.lo).QuoRem(v1)
		if err != nil {
			return u128{}, err
		}
	}

	vq := v.MulToU256(tq)

	// let k = 1 + [(u*rem(v, 2^(64-n))) / (v*(v-rem(v, 2^(64-n)))]
	// vq = v*tq = v(q + k)
	if vq.cmp(u) <= 0 {
		// vq <= u means tq = q
		return tq, nil
	}

	// vqu = vq - u = v*(q+k) - (vq + r) = v*k - r
	vqu, err := vq.sub(u)
	if err != nil {
		return u128{}, err
	}

	if !vqu.carry.IsZero() {
		// v * k > 2^128, we can't find k
		// fall back to big.Int division
		return u128{}, errOverflow
	}

	vqu128 := u128FromHiLo(vqu.hi, vqu.lo)

	// k1 = k - 1
	// vqu = v*k - r = v*(k1 + 1) - r = v*k1 + v - r
	// k1 <= [vqu / v] <= k1 + 1
	k1, _, err := vqu128.QuoRem(v)
	if err != nil {
		return u128{}, err
	}

	// adjust k1
	vqu1, err := v.Mul(k1)
	if err != nil {
		return u128{}, err
	}

	// if [vqu / v] = k1 + 1, then we don't have to adjust because final k = k1 + 1
	// if [vqu / v] = k1, then final k = k1 + 1
	if vqu1.Cmp(vqu128) < 0 {
		k1, err = k1.Add64(1)
		if err != nil {
			return u128{}, err
		}
	}

	// final q = tq - k
	tq, err = tq.Sub(k1)
	if err != nil {
		return u128{}, err
	}

	// we don't really need the remainder, might un-comment later if needed
	// r, err := v.Sub(r1)
	// if err != nil {
	// 	return u128{}, u128{}, err
	// }

	return tq, nil
}

// quoRem64Tou128 return q,r which:
//
//	q must be a u128
//	u = q*v + r
//	Return overflow if the result q doesn't fit in a u128
func (u u256) quoRem64Tou128(v uint64) (u128, uint64, error) {
	if u.carry.lo == 0 {
		q, r := u128FromHiLo(u.hi, u.lo).QuoRem64(v)
		return q, r, nil
	}

	quo, rem := u128FromHiLo(u.carry.lo, u.hi).QuoRem64(v)
	if quo.hi != 0 {
		return u128{}, 0, errOverflow
	}

	hi := quo.lo

	// can't panic because rem < v
	lo, r := bits.Div64(rem, u.lo, v)

	return u128FromHiLo(hi, lo), r, nil
}
