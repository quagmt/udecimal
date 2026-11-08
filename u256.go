package udecimal

import (
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

// pow returns u^e (with e > 0).
// Use int instead of uint to avoid unnecessary checks for type conversion.
//
// NOTE: Caller must ensure that e > 0 before calling this function.
func (u u256) pow(e int) (u256, error) {
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

// fastQuo returns quotient and remainder of u/v
func (u u256) fastQuo(v u128) (u128, u128, error) {
	if u.carry.IsZero() {
		q, r, err := u128FromHiLo(u.hi, u.lo).QuoRem(v)
		return q, r, err
	}

	if v.hi == 0 && u.carry.hi == 0 {
		q, r, err := u.div192by64(v.lo)
		return q, u128{lo: r}, err
	}

	// now we have u192 / u128 or u256 / u128
	if u.carry.Cmp(v) >= 0 {
		// obviously the result won't fit into u128
		return u128{}, u128{}, errOverflow
	}

	q, r := u.div256by128(v)
	return q, r, nil
}

// div192by64 return q,r which:
// q must be a u128
// u = q*v + r
// Returns error if u.carry >= v, because the result can't fit into u128
func (u u256) div192by64(v uint64) (u128, uint64, error) {
	if u.carry.Cmp64(v) >= 0 {
		return u128{}, 0, errOverflow
	}

	// can't panic because we already check u.carry < v (u.carry.hi == 0 && u.carry.lo < v)
	hi, rem := bits.Div64(u.carry.lo, u.hi, v)

	// can't panic because rem < v
	lo, r := bits.Div64(rem, u.lo, v)
	return u128FromHiLo(hi, lo), r, nil
}

// div256by128 returns q, r such that u = q*v + r and 0 <= r < v.
// Requires v.hi != 0 and u.carry < v (so q fits in 128 bits). fastQuo checks both.
//
// It is long division (Knuth's Algorithm D) where each digit is a uint64, explained in
// https://ridiculousfish.com/blog/posts/labor-of-division-episode-iv.html
//
// Comments below use base-10 examples: one decimal digit stands for one uint64,
// so "10" stands for 2^64 and "9" for 2^64-1.
func (u u256) div256by128(v u128) (u128, u128) {
	// 1. Normalize: shift u and v left by n bits so the top bit of v is set.
	// Step 3a guesses each quotient digit from v's top digit alone, so that digit must be big:
	//   e.g. 192 / 29 (q = 6): guess 19/2 = 9, off by 3
	//   e.g. x3 -> 576 / 87 (q = 6): guess 57/8 = 7, off by 1
	// With the top bit set, the guess is never off by more than 2.
	// q is unchanged, but r is scaled too (e.g. 54 = 18*3); step 4 shifts it back.
	n := bits.LeadingZeros64(v.hi)

	//nolint:gosec // 0 <= n <= 63, so it's safe to convert to uint
	v = v.Lsh(uint(n))

	// Shift u left by the same n, split into 4 words a = [a3, a2, a1, a0].
	// Each word = its own bits << n, plus the top n bits carried up from the word below.
	// The result still fits in 4 words because u.carry < v.
	a := [4]uint64{}
	a[0] = u.lo << n
	a[1] = u.lo>>(64-n) | u.hi<<n
	a[2] = u.hi>>(64-n) | u.carry.lo<<n
	a[3] = u.carry.lo>>(64-n) | u.carry.hi<<n

	// 2. Count quotient digits (1 or 2).
	// Long division by hand always does 2 steps here, and a leading 0 is fine:
	//   e.g. 0576 / 87:  057 / 87 = 0,  then 576 / 87 = 6  -> 06
	// When the first digit is 0 we skip that step to save time. The skip must be exact:
	// q has 2 digits exactly when the top 3 digits [a3, a2, a1] >= v, i.e.
	// a3 != 0 or [a2, a1] >= v. e.g. with v = 87:
	//   0576 / 87 = 6   -> 1 digit  (a3 = 0, 57 < 87)
	//   5766 / 87 = 66  -> 2 digits (a3 != 0)
	//   0896 / 87 = 10  -> 2 digits (a3 = 0, but 89 >= 87)
	//   0866 / 87 = 9   -> 1 digit  (a3 = 0, 86 < 87)
	// Don't compare only a2 with v.hi: when a2 == v.hi, q can have 1 or 2 digits depending
	// on a1 (e.g. 0896 vs 0866 above). The old check (a2 > v.hi) counted 0896 as 1 digit
	// (issue #46).
	qLen := 1
	a21 := u128{hi: a[2], lo: a[1]}
	if a[3] != 0 || a21.Cmp(v) >= 0 {
		qLen = 2
	}

	q := [2]uint64{}

	// 3. Compute q[i] from high to low. Each step divides the 3-digit window
	// [u2, u1, u0] = [a[i+2], a[i+1], a[i]] by the 2-digit v and gives 1 quotient digit:
	//   e.g. 5766 / 87:  i = 1: 576 / 87 = 6, rem 54
	//                    i = 0: 546 / 87 = 6, rem 24   -> q = 66, r = 24
	// Every window has [u2, u1] < v, so its quotient digit fits in one uint64:
	//   - first window: from u.carry < v (or from the qLen check above)
	//   - later windows: [u2, u1] is the previous remainder, which is < v
	for i := qLen - 1; i >= 0; i-- {
		u2, u1, u0 := a[i+2], a[i+1], a[i]

		// 3a. Guess the digit from the top digits only:
		//   tq = [u2, u1] / v.hi
		//   r  = [u2, u1] - tq*v.hi   (remainder of the guess, NOT the true remainder)
		// e.g. 576 / 87: tq = 57/8 = 7, r = 57 - 7*8 = 1   (true digit is 6)
		// The guess is never too small, and at most 2 too big (thanks to step 1).
		//
		// rOverflow is set below when r doesn't fit in a uint64 (only possible when u2 == v.hi).
		var tq, r uint64
		rOverflow := false

		if u2 == v.hi {
			// Problem: here the guess [u2, u1] / v.hi is 2^64 or 2^64+1, which doesn't fit in
			// a uint64, e.g. 866 / 87: guess = 86 / 8 = 10.
			// bits.Div64 panics when its result doesn't fit in one digit (issue #46).
			//
			// Fix: start the guess at the largest digit, 2^64-1 (Knuth: qhat = min(qhat, b-1)).
			// This works because when u2 == v.hi, the true digit is always 2^64-2 or 2^64-1:
			//   - At most 2^64-1, because [u2, u1] < v means the window < 2^64*v.
			//     e.g. 866 / 87: 866 < 870 = 10*87, so the digit is < 10.
			//   - At least 2^64-2, because the window and v both start with v.hi, and
			//     v.hi >= 2^63 (step 1).
			//     e.g. v.hi = 8: the window is >= 800 and v <= 89, so the digit is >= 800/89 = 8.98.
			// So 2^64-1 is either right or 1 too big. Step 3b checks it like any other guess:
			//   e.g. 866 / 87: 9*87 = 783 <= 866  -> 9 is right   (rem 83)
			//   e.g. 500 / 51: 9*51 = 459 <= 500  -> 9 is right   (rem 41)
			//   e.g. 500 / 56: 9*56 = 504 >  500  -> 1 too big, 8 (rem 52)
			// (u2 > v.hi can't happen, because [u2, u1] < v.)
			tq = ^uint64(0)

			// Overflow: with tq = 2^64-1, the remainder of the guess might not fit in a uint64:
			//   r = [u2, u1] - (2^64-1)*v.hi
			//     = u2*2^64 + u1 - (2^64-1)*v.hi
			//     = u1 + v.hi   (because u2 == v.hi)
			//   e.g. 500 / 51: r = 50 - 9*5 = 5  = 0 + 5  -> 1 digit, fine
			//   e.g. 866 / 87: r = 86 - 9*8 = 14 = 6 + 8  -> 2 digits, overflow
			// (In the normal case below, r < v.hi, so it always fits.)
			// bits.Add64 keeps the low 64 bits and returns the lost top bit as carry
			// (e.g. 14 -> 4, carry 1). We remember it in rOverflow because it matters later:
			//   - 3b: c2 = [r, u0] loses that bit (e.g. 146 is stored as 46), so comparing
			//     with it is wrong
			//   - 3c: the remainder still comes out right, because subtraction wraps around
			var carry uint64
			r, carry = bits.Add64(u1, v.hi, 0)
			rOverflow = carry != 0
		} else {
			// Normal case: u2 < v.hi, so the guess fits in one digit.
			tq, r = bits.Div64(u2, u1, v.hi)
		}

		// 3b. Check the guess against the full window. Since [u2, u1] = tq*v.hi + r:
		//   [u2, u1, u0] = tq*v.hi*2^64 + [r, u0]
		//   tq*v         = tq*v.hi*2^64 + tq*v.lo
		// so the true remainder of the guess is [u2, u1, u0] - tq*v = c2 - c1, where
		//   c1 = tq*v.lo   (the part of tq*v the guess ignored)
		//   c2 = [r, u0]   (what is left of the window to cover it)
		// e.g. 576 / 87: c1 = 7*7 = 49, c2 = [1,6] = 16 -> 16 - 49 < 0, so 7 is too big
		c1h, c1l := bits.Mul64(tq, v.lo)
		c1 := u128{hi: c1h, lo: c1l}
		// If r overflowed, the true c2 = r*2^64 + u0 is >= 2^128, and this u128 is missing
		// the top bit, e.g. 866 / 87: true c2 = [14, 6] = 146, stored c2 = [4, 6] = 46
		c2 := u128{hi: r, lo: u0}

		// k = how much tq overshoots the true digit (0, 1 or 2):
		//   c1 <= c2           -> k = 0   e.g. 500/51: tq = 9, c1 = 9, c2 = 50 -> 9
		//   d = c1 - c2 <= v   -> k = 1   e.g. 576/87: d = 49 - 16 = 33 <= 87 -> 7-1 = 6
		//   d > v              -> k = 2   e.g. 450/59: tq = 45/5 = 9, r = 0, c1 = 81, c2 = 0,
		//                                      d = 81 > 59 -> 9-2 = 7
		// (k = 2: after subtracting one v the remainder is v - d, still negative if d > v.)
		var k uint64
		switch {
		case rOverflow:
			// r overflowed (see 3a): k is always 0, so tq = 2^64-1 is right.
			// The true c2 is >= 2^128 (because r >= 2^64), while c1 = tq*v.lo < 2^128,
			// so c1 < c2.
			//   e.g. 866 / 87: c1 = 9*7 = 63 < 146 -> 9 is right
			// Don't compare c1 with the stored c2: it lost its top bit, so it can look smaller
			// than c1 (e.g. 63 > 46 would wrongly say 9 is too big).
		case c1.Cmp(c2) > 0:
			k = 1

			if subUnsafe(c1, c2).Cmp(v) > 0 {
				k = 2
			}
		}

		q[i] = tq - k

		// 3c. True remainder rem = [u2, u1, u0] - q[i]*v = c2 - c1 + k*v.
		// It is always in [0, v), so it fits in a u128. Each case is ordered so that no
		// intermediate value overflows.
		var rem u128
		switch k {
		case 0:
			// rem = c2 - c1
			// e.g. 500/51: 50 - 9 = 41  (500 - 9*51 = 41)
			//
			// When r overflowed, the stored c2 is missing its top bit, but this still works:
			// subUnsafe wraps around (mod 2^128), and the true answer (< v) fits in a u128,
			// e.g. 866 / 87: 46 - 63 wraps to 83 (mod 100) = 146 - 63 ✓
			rem = subUnsafe(c2, c1)
		case 1:
			// rem = v - d
			// e.g. 576/87: 87 - 33 = 54  (576 - 6*87 = 54)
			rem = subUnsafe(c1, c2)
			rem = subUnsafe(v, rem)
		case 2:
			// rem = 2v - d, but 2v may overflow. Use not(v) = max - v (max = 2^128-1):
			//   2v - d = (v - not(v)) + (max - d)
			// both parts >= 0 (v normalized -> v >= not(v)), sum = rem < v, so no overflow.
			// e.g. 450/59 (max = 99, not(v) = 40): (59 - 40) + (99 - 81) = 19 + 18 = 37
			//      (450 - 7*59 = 37)
			c12 := subUnsafe(c1, c2)
			c12 = subUnsafe(max128, c12)
			rem = subUnsafe(v, u128{hi: ^v.hi, lo: ^v.lo})

			rem, _ = rem.Add(c12)
		}

		// rem replaces [u2, u1]; next window is [rem.hi, rem.lo, a[i-1]] (bring down next digit).
		// e.g. 5766/87: 576 -> q1 = 6, rem 54 -> next window 546 -> q0 = 6, rem 24
		a[i+1], a[i] = rem.hi, rem.lo
	}

	// 4. Undo normalization on the remainder: [a1, a0] >> n.
	// e.g. 54 / 3 = 18  (192 = 6*29 + 18)
	//nolint:gosec // 0 <= n <= 63, so it's safe to convert to uint
	r := u128{hi: a[1], lo: a[0]}.Rsh(uint(n))

	return u128{hi: q[1], lo: q[0]}, r
}

// subUnsafe returns u - v with u >= v
// must be called only when u >= v or the result will be incorrect
func subUnsafe(u, v u128) u128 {
	lo, borrow := bits.Sub64(u.lo, v.lo, 0)
	hi, _ := bits.Sub64(u.hi, v.hi, borrow)
	return u128{hi: hi, lo: lo}
}
