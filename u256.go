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
// 	l := len(s)

// 	n := 64 - l

// 	for range n {
// 		s = "0" + s
// 	}

// 	return s
// }

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

func (u U256) Sub128(v bint) (U256, error) {
	lo, borrow := bits.Sub64(u.lo, v.lo, 0)
	hi, borrow := bits.Sub64(u.hi, v.hi, borrow)

	c1, err := u.carry.Sub64(borrow)
	if err != nil {
		return U256{}, err
	}

	return U256{lo: lo, hi: hi, carry: c1}, nil

}

// compare against bint
//
//	+1 when u > v
//	 0 when u = v
//	-1 when u < v
func (u U256) CmpBint(v bint) int {
	if !u.carry.IsZero() {
		return 1
	}

	b := bint{hi: u.hi, lo: u.lo}
	return b.Cmp(v)
}

func (u U256) Lsh(n uint) (v U256) {
	switch {
	case n < 64:
		v.carry = u.carry.Lsh(n)
		v.carry.lo = v.carry.lo | u.hi>>(64-n)
		c := bint{hi: u.hi, lo: u.lo}.Lsh(n)
		v.hi = c.hi
		v.lo = c.lo

	case 64 <= n && n < 128:
		v.lo = 0
		v.hi = u.lo << (n - 64)
		v.carry.lo = u.hi<<(n-64) | u.lo>>(128-n)
		v.carry.hi = u.carry.lo<<(n-64) | u.hi>>(128-n)

	case n >= 128:
		v.lo, v.hi = 0, 0
		v.carry = bint{hi: u.hi, lo: u.lo}.Lsh(n - 128)

	default:
		// n < 0, can't happen
	}

	return
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
func (u U256) Quo(v bint) (bint, error) {
	if u.carry.IsZero() {
		b := bint{hi: u.hi, lo: u.lo}
		q, _, err := b.QuoRem(v)
		return q, err
	}

	if v.hi == 0 {
		q, _, err := u.QuoRem64(v.lo)
		return q, err
	}

	// if u >= 2^192, the quotient won't fit in 128-bits number (overflow)
	// put in both here and inside QuoRem64, in case we call QuoRem64 directly
	if u.carry.hi != 0 {
		return bint{}, ErrOverflow
	}

	// 0 <= n <= 63
	n := uint(bits.LeadingZeros64(v.hi))
	v1 := v.Lsh(n)
	u1 := u.Lsh(n).Rsh(64)

	// let q, r are final quotient and remainder
	// calculate 'trial quotient' tq: tq < u/v + 2^64 --> tq < q + 2^64
	// let tq = q + k --> k < 2^64
	tq, _, err := u1.QuoRem64(v1.hi)
	if err != nil {
		return bint{}, err
	}

	vq := v.MulToU256(tq)

	// vqu = (q+k)*v - (q*v + r) = k*v - r
	// k*v < 2^128 --> vqu < 2^128 and can be represented in a bint (no overflow)
	vqu, err := vq.Sub(u)
	if err != nil {
		return bint{}, err
	}

	// techically this can't happen, just put it here to do fuzz test
	if vqu.carry.hi&vqu.carry.lo != 0 {
		return bint{}, ErrOverflow
	}

	// k1 = k - 1
	k1, _, err := bint{hi: vqu.hi, lo: vqu.lo}.QuoRem(v)
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

	// we don't really care abount the remainder, might uncomment later if needed
	// r, err := v.Sub(r1)
	// if err != nil {
	// 	return bint{}, bint{}, err
	// }

	return tq, nil
}

func (u U256) QuoRem64(v uint64) (q bint, r bint, err error) {
	// obvious case that the result won't fit in 128-bits number
	if u.carry.hi != 0 {
		err = ErrOverflow
		return
	}

	b := bint{hi: u.carry.lo, lo: u.hi}
	quo, rem := b.QuoRem64(v)
	if quo.hi != 0 {
		err = ErrOverflow
		return
	}

	q.hi = quo.lo
	q.lo, r.lo = bits.Div64(rem, u.lo, v)
	return
}
