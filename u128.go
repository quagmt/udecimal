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

func (u bint) Mul(v bint) (bint, error) {
	if u.hi&v.hi != 0 {
		return bint{}, ErrOverflow
	}

	p0, p1 := bits.Mul64(u.hi, v.lo)
	p2, p3 := bits.Mul64(u.lo, v.hi)

	if p0 != 0 || p2 != 0 {
		return bint{}, ErrOverflow
	}

	hi, lo := bits.Mul64(u.lo, v.lo)
	hi, c0 := bits.Add64(hi, p1, 0)
	hi, c1 := bits.Add64(hi, p3, c0)
	if c1 != 0 {
		return bint{}, ErrOverflow
	}

	return overflow(hi, lo)
}

func (u bint) Quo(v bint) (bint, error) {
	return bint{}, nil
}

func (u bint) QuoRem(v bint) (bint, error) {
	return bint{}, nil
}
