package udecimal

import (
	"math/bits"
)

// bint (big unsigned-integer) is a 128-bits unsigned integer
// represents by two 64-bits unsigned integer
// value = hi*2^64 + lo
type bint struct {
	hi uint64
	lo uint64
}

func (u bint) IsZero() bool {
	return u == bint{}
}

// Cmp compares u, v and returns:
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

func (u bint) Add(v bint) (bint, error) {
	lo, carry := bits.Add64(u.lo, v.lo, 0)
	hi, carry := bits.Add64(u.hi, v.hi, carry)
	if carry != 0 {
		return bint{}, ErrOverflow
	}

	return bint{hi: hi, lo: lo}, nil
}

// Add64 returns u+v.
func (u bint) Add64(v uint64) (bint, error) {
	lo, carry := bits.Add64(u.lo, v, 0)
	hi, carry := bits.Add64(u.hi, 0, carry)
	if carry != 0 {
		return bint{}, ErrOverflow
	}

	return bint{hi: hi, lo: lo}, nil
}

func (u bint) Sub(v bint) (bint, error) {
	lo, borrow := bits.Sub64(u.lo, v.lo, 0)
	hi, borrow := bits.Sub64(u.hi, v.hi, borrow)
	if borrow != 0 {
		// borrow != 0 means u < v and this must not happen
		return bint{}, ErrOverflow
	}

	return bint{hi: hi, lo: lo}, nil
}

// Sub64 returns u-v.
func (u bint) Sub64(v uint64) (bint, error) {
	lo, borrow := bits.Sub64(u.lo, v, 0)
	hi, borrow := bits.Sub64(u.hi, 0, borrow)
	if borrow != 0 {
		return bint{}, ErrOverflow
	}

	return bint{hi: hi, lo: lo}, nil
}

func (u bint) Mul64(v uint64) (bint, error) {
	hi, lo := bits.Mul64(u.lo, v)
	p0, p1 := bits.Mul64(u.hi, v)
	hi, c0 := bits.Add64(hi, p1, 0)
	if p0 != 0 || c0 != 0 {
		return bint{}, ErrOverflow
	}

	return bint{hi: hi, lo: lo}, nil
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

	return bint{hi: hi, lo: lo}, nil
}

func (u bint) MulToU256(v bint) U256 {
	hi, lo := bits.Mul64(u.lo, v.lo)
	p0, p1 := bits.Mul64(u.hi, v.lo)
	p2, p3 := bits.Mul64(u.lo, v.hi)

	// calculate hi + p1 + p3
	// total carry = carry(hi+p1) + carry(hi+p1+p3)
	hi, c0 := bits.Add64(hi, p1, 0)
	hi, c1 := bits.Add64(hi, p3, 0)
	c1 += c0

	// calculate upper part of U256
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

func bintFromU64(v uint64) bint {
	return bint{lo: v}
}

func bintFromHiLo(hi, lo uint64) bint {
	return bint{hi: hi, lo: lo}
}

// QuoRem returns q = u/v and r = u%v.
func (u bint) QuoRem(v bint) (q, r bint, err error) {
	if v.hi == 0 {
		var r64 uint64
		q, r64 = u.QuoRem64(v.lo)
		r = bintFromU64(r64)
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

		q = bintFromU64(tq)
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
	if n >= 64 {
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
	if n >= 64 {
		s.lo = u.hi >> (n - 64)
		s.hi = 0
	} else {
		s.lo = u.lo>>n | u.hi<<(64-n)
		s.hi = u.hi >> n
	}
	return
}
