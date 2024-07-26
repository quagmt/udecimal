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

func (u U256) QuoRem(v bint) (bint, bint, error) {
	if u.carry.IsZero() {
		return bint{hi: u.hi, lo: u.lo}.QuoRem(v)
	}

	if v.hi == 0 {
		return u.QuoRem64(v.lo)
	}

	// TODO
	return bint{}, bint{}, nil
}

func (u U256) QuoRem64(v uint64) (q bint, r bint, err error) {
	if u.carry.hi != 0 {
		err = ErrOverflow
		return
	}

	// u.carry.hi = 0
	// if v <= u.lo {
	var k uint64
	q.hi, k = bits.Div64(u.carry.lo, u.hi, v)
	q.lo, r.lo = bits.Div64(k, u.lo, v)
	return
	// }

	// // v > u.lo --> need to borrow 1 from u.hi
	// if u.hi <= 1 {
	// 	err = fmt.Errorf("unexpected u.hi value: %d. Should be greater than 1", u.hi)
	// 	return
	// }

	// q.lo, r.lo = bits.Div64(1, u.lo, v)
	// q.hi, r.hi = bits.Div64(u.carry.lo, u.hi-1, v)
	// return
}
