package udecimal

import (
	"fmt"
	"math/bits"
)

var (
	// maxBint is the maximum value that can be represented by a bint.
	// consists of 38 digits (1 digit less than 2^128-1) to avoid overflow.
	// maxBInt = 99_999_999_999_999_999_999_999_999_999_999_999_999
	maxBint = bint{
		hi: 5_421_010_862_427_522_170,
		lo: 687_399_551_400_673_279,
	}
)

var (
	ErrOverflow = fmt.Errorf("overflow")
)

// bint (big unsigned-integer) is a 128-bits unsigned integer
// represents by two 64-bits unsigned integer
// value = hi*2^64 + lo
type bint struct {
	hi uint64
	lo uint64
}

func (b bint) IsZero() bool {
	return b == bint{}
}

// Cmp compares u, v and retuns:
//
//	-1 if u < v
//	0 if u == v
//	1 if u > v
func (u bint) Cmp(v bint) int {
	if u.hi < v.hi {
		return -1
	}

	if u.hi > v.hi {
		return 1
	}

	// u.hi == v.hi
	switch {
	case u.lo < v.lo:
		return -1
	case u.lo > v.lo:
		return 1
	default:
		return 0
	}
}

func (u bint) GreaterThan(v bint) bool {
	if u.hi > v.hi || (u.hi == v.hi && u.lo > v.lo) {
		return true
	}

	return false
}

func (u bint) LessThan(v bint) bool {
	if u.hi < v.hi || (u.hi == v.hi && u.lo < v.lo) {
		return true
	}

	return false
}

func overflow(hi, lo uint64) (bint, error) {
	u := bint{
		hi: hi,
		lo: lo,
	}

	if u.GreaterThan(maxBint) {
		return bint{}, ErrOverflow
	}

	return u, nil
}

func (u bint) Add(v bint) (bint, error) {
	lo, carry := bits.Add64(u.lo, v.lo, 0)
	hi, carry := bits.Add64(u.hi, v.hi, carry)
	if carry != 0 {
		return bint{}, ErrOverflow
	}

	return overflow(hi, lo)
}

// Add64 returns u+v.
func (u bint) Add64(v uint64) (bint, error) {
	lo, carry := bits.Add64(u.lo, v, 0)
	hi, carry := bits.Add64(u.hi, 0, carry)
	if carry != 0 {
		return bint{}, ErrOverflow
	}

	return overflow(hi, lo)
}

func (u bint) Sub(v bint) (bint, error) {
	lo, borrow := bits.Sub64(u.lo, v.lo, 0)
	hi, borrow := bits.Sub64(u.hi, v.hi, borrow)
	if borrow != 0 {
		return bint{}, ErrOverflow
	}

	return overflow(hi, lo)
}

// Sub64 returns u-v.
func (u bint) Sub64(v uint64) (bint, error) {
	lo, borrow := bits.Sub64(u.lo, v, 0)
	hi, borrow := bits.Sub64(u.hi, 0, borrow)
	if borrow != 0 {
		return bint{}, ErrOverflow
	}

	return overflow(hi, lo)
}

func (u bint) Mul64(v uint64) (bint, error) {
	hi, lo := bits.Mul64(u.lo, v)
	p0, p1 := bits.Mul64(u.hi, v)
	hi, c0 := bits.Add64(hi, p1, 0)
	if p0 != 0 || c0 != 0 {
		return bint{}, ErrOverflow
	}

	return overflow(hi, lo)
}

func (u bint) Mul(v bint) (bint, error) {
	if u.hi&v.hi != 0 {
		return bint{}, ErrOverflow
	}

	if v.hi == 0 {
		return u.Mul64(v.lo)
	}

	p0, p1 := bits.Mul64(u.hi, v.lo)
	p2, p3 := bits.Mul64(u.lo, v.hi)

	if p0 != 0 || p2 != 0 {
		return bint{}, ErrOverflow
	}

	hi, lo := bits.Mul64(u.lo, v.lo)
	hi, c0 := bits.Add64(hi, p1, 0)
	hi, c1 := bits.Add64(hi, p3, 0)
	if c1&c0 != 0 {
		return bint{}, ErrOverflow
	}

	return overflow(hi, lo)
}

// MulPow10 is similar to Mul. However it doesn't return error
// but 256-bits unsigned integer instead. The carry will be stored to U256.carry
func (u bint) MulPow10(pow int) U256 {
	v := pow10[pow]

	// TODO: can speed up with SIMD
	hi, lo := bits.Mul64(u.lo, v.lo)
	p0, p1 := bits.Mul64(u.hi, v.lo)
	p2, p3 := bits.Mul64(u.lo, v.hi)

	// calculate hi,lo
	// NOTE: carryIn doesn't mean the carry value from previous calculation
	// under the hood: sum = hi + p1 + carryIn (add 0/1 to the total sum???)
	// hence, hi, c1 := bits.Add64(hi, p3, c0) is incorrect
	// The total overflow carry must be sum of all carries
	hi, c0 := bits.Add64(hi, p1, 0)
	hi, c1 := bits.Add64(hi, p3, 0)
	c1 += c0

	// calculate carry
	// TODO: can speed up with SIMD
	e0, e1 := bits.Mul64(u.hi, v.hi)
	d, d0 := bits.Add64(p0, p2, 0)
	d, d1 := bits.Add64(d, c1, 0)
	e2, e3 := bits.Add64(d, e1, 0)

	carry := bint{
		hi: e0 + d0 + d1 + e3, // result can't be overflow because max(pow) = pow[39] < 2^128
		lo: e2,
	}

	return U256{
		lo:    lo,
		hi:    hi,
		carry: carry,
	}
}

func (u bint) Mul64Pow10(v uint64) U256 {
	return U256{}
}

func FromU64(v uint64) bint {
	return bint{lo: v}
}

// QuoRem returns q = u/v and r = u%v.
func (u bint) QuoRem(v bint) (q, r bint, err error) {
	if v.hi == 0 {
		var r64 uint64
		q, r64 = u.QuoRem64(v.lo)
		r = FromU64(r64)
	} else {
		// generate a "trial quotient," guaranteed to be within 1 of the actual
		// quotient, then adjust.
		n := uint(bits.LeadingZeros64(v.hi))
		v1 := v.Lsh(n)
		u1 := u.Rsh(1)
		tq, _ := bits.Div64(u1.hi, u1.lo, v1.hi)
		tq >>= 63 - n
		if tq != 0 {
			tq--
		}

		q = FromU64(tq)
		// calculate remainder using trial quotient, then adjust if remainder is
		// greater than divisor
		// r = u.Sub(v.Mul64(tq))
		// if r.Cmp(v) >= 0 {
		// 	q = q.Add64(1)
		// 	r = r.Sub(v)
		// }

		vq, err := v.Mul64(tq)
		if err != nil {
			return q, r, err
		}

		r, err = u.Sub(vq)
		if err != nil {
			return q, r, err
		}

		if r.Cmp(v) >= 0 {
			q, err = q.Add64(1)
			if err != nil {
				return q, r, err
			}

			r, err = r.Sub(v)
			if err != nil {
				return q, r, err
			}
		}
	}
	return
}

// QuoRem64 returns q = u/v and r = u%v.
func (u bint) QuoRem64(v uint64) (q bint, r uint64) {
	if u.hi < v {
		q.lo, r = bits.Div64(u.hi, u.lo, v)
	} else {
		q.hi, r = bits.Div64(0, u.hi, v)
		q.lo, r = bits.Div64(r, u.lo, v)
	}
	return
}

// Lsh returns u<<n.
func (u bint) Lsh(n uint) (s bint) {
	if n > 64 {
		s.lo = 0
		s.hi = u.lo << (n - 64)
	} else {
		s.lo = u.lo << n
		s.hi = u.hi<<n | u.lo>>(64-n)
	}
	return
}

// Rsh returns u>>n.
func (u bint) Rsh(n uint) (s bint) {
	if n > 64 {
		s.lo = u.hi >> (n - 64)
		s.hi = 0
	} else {
		s.lo = u.lo>>n | u.hi<<(64-n)
		s.hi = u.hi >> n
	}
	return
}

// String returns the base-10 representation of u as a string.
func (u bint) String() string {
	if u.IsZero() {
		return "0"
	}
	buf := []byte("0000000000000000000000000000000000000000") // log10(2^128) < 40
	for i := len(buf); ; i -= 19 {
		q, r := u.QuoRem64(1e19) // largest power of 10 that fits in a uint64
		var n int
		for ; r != 0; r /= 10 {
			n++
			buf[i-n] += byte(r % 10)
		}
		if q.IsZero() {
			return string(buf[i-n:])
		}
		u = q
	}
}
