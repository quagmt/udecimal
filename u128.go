package udecimal

import (
	"encoding/binary"
	"math/big"
	"math/bits"
)

var (
	one128 = u128{lo: 1}
)

// u128 (big unsigned-integer) is a 128-bits unsigned integer
// represents by two 64-bits unsigned integer
// value = hi*2^64 + lo
type u128 struct {
	hi uint64
	lo uint64
}

// bitLen returns the number of bits required to represent u
func (u u128) bitLen() int {
	if u.hi != 0 {
		return bits.Len64(u.hi) + 64
	}

	return bits.Len64(u.lo)
}

// IsZero returns true if u is zero
func (u u128) IsZero() bool {
	return u == u128{}
}

// isOverflow returns true if coef is greater than or equal to 10^38
// coef should be less than 10^38 to take advantage of 128-bits unsigned integer
func (u u128) isOverflow() bool {
	// scale = frac digits
	// whole part has at most 19 digits
	// consider it's overflow when total digits > scale + 19, which means coef >= 10^(scale+19)
	return !u.LessThan(pow10[38])
}

// Cmp compares u, v and returns:
//
//	-1 if u < v
//	0 if u == v
//	1 if u > v
func (u u128) Cmp(v u128) int {
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

func (u u128) Cmp64(v uint64) int {
	if u.hi != 0 {
		return 1
	}

	switch {
	case u.lo < v:
		return -1
	case u.lo > v:
		return 1
	default:
		return 0
	}
}

func (u u128) LessThan(v u128) bool {
	if u.hi < v.hi || (u.hi == v.hi && u.lo < v.lo) {
		return true
	}

	return false
}

func (u u128) Add(v u128) (u128, error) {
	q, err := u.addRaw(v)
	if err != nil {
		return u128{}, err
	}

	if q.isOverflow() {
		return u128{}, ErrOverflow
	}

	return q, nil
}

func (u u128) addRaw(v u128) (u128, error) {
	lo, carry := bits.Add64(u.lo, v.lo, 0)
	hi, carry := bits.Add64(u.hi, v.hi, carry)
	if carry != 0 {
		return u128{}, ErrOverflow
	}

	return u128{hi: hi, lo: lo}, nil
}

// Add64 returns u+v.
func (u u128) Add64(v uint64) (u128, error) {
	q, err := u.add64Raw(v)
	if err != nil {
		return u128{}, err
	}

	if q.isOverflow() {
		return u128{}, ErrOverflow
	}

	return q, nil
}

func (u u128) add64Raw(v uint64) (u128, error) {
	lo, carry := bits.Add64(u.lo, v, 0)
	hi, carry := bits.Add64(u.hi, 0, carry)
	if carry != 0 {
		return u128{}, ErrOverflow
	}

	return u128{hi: hi, lo: lo}, nil
}

func (u u128) Sub(v u128) (u128, error) {
	q, err := u.subRaw(v)
	if err != nil {
		return u128{}, err
	}

	if q.isOverflow() {
		return u128{}, ErrOverflow
	}

	return q, nil
}

func (u u128) subRaw(v u128) (u128, error) {
	lo, borrow := bits.Sub64(u.lo, v.lo, 0)
	hi, borrow := bits.Sub64(u.hi, v.hi, borrow)
	if borrow != 0 {
		// borrow != 0 means u < v and this must not happen
		return u128{}, ErrOverflow
	}

	return u128{hi: hi, lo: lo}, nil
}

// Sub64 returns u-v.
func (u u128) Sub64(v uint64) (u128, error) {
	q, err := u.sub64Raw(v)
	if err != nil {
		return u128{}, err
	}

	if q.isOverflow() {
		return u128{}, ErrOverflow
	}

	return q, nil
}

func (u u128) sub64Raw(v uint64) (u128, error) {
	lo, borrow := bits.Sub64(u.lo, v, 0)
	hi, borrow := bits.Sub64(u.hi, 0, borrow)
	if borrow != 0 {
		return u128{}, ErrOverflow
	}

	return u128{hi: hi, lo: lo}, nil
}

// Mul64 returns u*v.
// Return overflow if the result is greater than 10^38-1
func (u u128) Mul64(v uint64) (u128, error) {
	q, err := u.mul64Raw(v)
	if err != nil {
		return u128{}, err
	}

	if q.isOverflow() {
		return u128{}, ErrOverflow
	}

	return q, nil
}

func (u u128) mul64Raw(v uint64) (u128, error) {
	hi, lo := bits.Mul64(u.lo, v)
	p0, p1 := bits.Mul64(u.hi, v)
	hi, c0 := bits.Add64(hi, p1, 0)
	if p0 != 0 || c0 != 0 {
		return u128{}, ErrOverflow
	}

	return u128{hi: hi, lo: lo}, nil
}

func (u u128) Mul(v u128) (u128, error) {
	if u.hi != 0 && v.hi != 0 {
		return u128{}, ErrOverflow
	}

	if v.hi == 0 {
		return u.Mul64(v.lo)
	}

	// u.hi == 0
	return v.Mul64(u.lo)
}

// raw version of Mul, doesn't check custom overflow condition
func (u u128) mulRaw(v u128) (u128, error) {
	if u.hi != 0 && v.hi != 0 {
		return u128{}, ErrOverflow
	}

	if v.hi == 0 {
		return u.mul64Raw(v.lo)
	}

	// u.hi == 0
	return v.mul64Raw(u.lo)
}

func (u u128) MulToU256(v u128) U256 {
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

	carry := u128{
		hi: e0 + d0 + d1 + e3, // result can't be overflow because max(pow) = pow[39] < 2^128
		lo: e2,
	}

	return U256{
		lo:    lo,
		hi:    hi,
		carry: carry,
	}
}

func u128FromU64(v uint64) u128 {
	return u128{lo: v}
}

func u128FromHiLo(hi, lo uint64) u128 {
	return u128{hi: hi, lo: lo}
}

// QuoRem returns q = u/v and r = u%v.
func (u u128) QuoRem(v u128) (q, r u128, err error) {
	if v.hi == 0 {
		var r64 uint64
		q, r64 = u.QuoRem64(v.lo)
		r = u128FromU64(r64)
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

		q = u128FromU64(tq)
		vq, err := v.mul64Raw(tq)
		if err != nil {
			return q, r, err
		}

		r, err = u.subRaw(vq)
		if err != nil {
			return q, r, err
		}

		if r.Cmp(v) >= 0 {
			q, err = q.add64Raw(1)
			if err != nil {
				return q, r, err
			}

			r, err = r.subRaw(v)
			if err != nil {
				return q, r, err
			}
		}
	}
	return
}

// QuoRem64 returns q = u/v and r = u%v.
func (u u128) QuoRem64(v uint64) (q u128, r uint64) {
	if u.hi < v {
		q.lo, r = bits.Div64(u.hi, u.lo, v)
	} else {
		q.hi, r = bits.Div64(0, u.hi, v)
		q.lo, r = bits.Div64(r, u.lo, v)
	}
	return
}

// Lsh returns u<<n.
func (u u128) Lsh(n uint) (s u128) {
	if n >= 64 {
		s.lo = 0
		s.hi = u.lo << (n - 64)
	} else {
		s.lo = u.lo << n
		s.hi = u.hi<<n | u.lo>>(64-n)
	}
	return
}

func (u u128) String() string {
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

// Rsh returns u>>n.
func (u u128) Rsh(n uint) (s u128) {
	if n >= 64 {
		s.lo = u.hi >> (n - 64)
		s.hi = 0
	} else {
		s.lo = u.lo>>n | u.hi<<(64-n)
		s.hi = u.hi >> n
	}
	return
}

func (u u128) ToBigInt() *big.Int {
	bytes := make([]byte, 16)
	binary.BigEndian.PutUint64(bytes, u.hi)
	binary.BigEndian.PutUint64(bytes[8:], u.lo)

	return new(big.Int).SetBytes(bytes)
}

// getTrailingZeros64 returns the number of trailing zeros in u
// NOTE: this only works when maxScale is 19
func getTrailingZeros64(u uint64) uint8 {
	var z uint8
	if u%1e16 == 0 {
		z = 16

		if u%pow10[z+2].lo == 0 {
			z += 2
		}

		if u%pow10[z+1].lo == 0 {
			z++
		}

		return z
	}

	if u%pow10[8].lo == 0 {
		z += 8
	}

	if u%pow10[z+4].lo == 0 {
		z += 4
	}

	if u%pow10[z+2].lo == 0 {
		z += 2
	}

	if u%pow10[z+1].lo == 0 {
		z++
	}

	return z
}
