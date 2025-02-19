//go:build fuzz

package udecimal

import (
	"encoding/binary"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"

	ss "github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

var corpus = []struct {
	neg  bool
	hi   uint64
	lo   uint64
	prec uint8
}{
	{false, 0, 0, 0},
	{false, 1, 0, 0},
	{false, 1234567890123456789, 0, 0},
	{true, 1, 0, 0},
	{false, 1123, 0, 3},
	{true, 1123, 0, 3},
	{false, 123123, 0, 6},
	{true, 123123, 0, 6},
	{false, 123456789123456789, 1234567890123456789, 9},
	{true, 123456789123456789, 1234567890123456789, 9},
	{false, 0, 1234567890123456789, 19},
	{true, 0, 1234567890123456789, 19},
	{false, 0, 1, 19},
	{true, 0, 1, 19},
	{false, math.MaxUint64, math.MaxUint64, 0},
	{false, math.MaxUint64, math.MaxUint64, 10},
	{true, math.MaxUint64, math.MaxUint64, 0},
	{true, math.MaxUint64, math.MaxUint64, 10},
}

func ssDecimal(neg bool, hi, lo uint64, prec uint8) ss.Decimal {
	bytes := make([]byte, 16)
	binary.BigEndian.PutUint64(bytes, hi)
	binary.BigEndian.PutUint64(bytes[8:], lo)

	bint := new(big.Int).SetBytes(bytes)

	if neg {
		bint = bint.Neg(bint)
	}

	d := ss.NewFromBigInt(bint, -int32(prec))
	return d
}

func FuzzParse(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		require.NoError(t, err)

		c := a.Mul(b)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		cstr := c.String()
		d := MustParse(cstr)

		require.Equal(t, cstr, d.String())
	})
}

func FuzzAddDec(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		require.NoError(t, err)

		c := a.Add(b)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)
		bb := ssDecimal(bneg, bhi, blo, bprec)
		cc := aa.Add(bb)

		require.Equal(t, cc.String(), c.String(), "add %s %s", a, b)
	})
}

func FuzzAdd64(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.lo)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, blo uint64) {
		aprec = aprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		c := a.Add64(blo)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)
		bb := ssDecimal(false, 0, blo, 0)
		cc := aa.Add(bb)

		require.Equal(t, cc.String(), c.String(), "add64 %s %d", a, blo)
	})
}

func FuzzSubDec(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		require.NoError(t, err)

		c := a.Sub(b)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)
		bb := ssDecimal(bneg, bhi, blo, bprec)
		cc := aa.Sub(bb)

		require.Equal(t, cc.String(), c.String(), "sub %s %s", a, b)
	})
}

func FuzzSub64(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.lo)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, blo uint64) {
		aprec = aprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		c := a.Sub64(blo)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)
		bb := ssDecimal(false, 0, blo, 0)
		cc := aa.Sub(bb)

		require.Equal(t, cc.String(), c.String(), "sub64 %s %d", a, blo)
	})
}

func FuzzMulDec(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		if err == ErrPrecOutOfRange {
			t.Skip()
		} else {
			require.NoError(t, err)
		}

		c := a.Mul(b)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)
		bb := ssDecimal(bneg, bhi, blo, bprec)

		prec := int32(c.Prec())
		cc := aa.Mul(bb).Truncate(prec)

		require.Equal(t, cc.String(), c.String(), "mul %s %s", a, b)
	})
}

func FuzzMul64(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.lo)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, blo uint64) {
		aprec = aprec % maxPrec
		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		c := a.Mul64(blo)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)
		bb := ssDecimal(false, 0, blo, 0)

		prec := int32(c.Prec())
		cc := aa.Mul(bb).Truncate(prec)

		require.Equal(t, cc.String(), c.String(), "mul64 %s %d", a, blo)
	})
}

func FuzzDivDec(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		require.NoError(t, err)

		c, err := a.Div(b)
		if b.IsZero() {
			require.Equal(t, ErrDivideByZero, err)
			return
		}

		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		require.NoError(t, err)

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)
		bb := ssDecimal(bneg, bhi, blo, bprec)

		prec := int32(c.Prec())
		cc := aa.DivRound(bb, 28).Truncate(prec)

		// sometimes shopspring/decimal does rounding differently
		// e.g. 0.099999999999999 -> 0.1
		// so to check the result, we can check the difference
		// between our result and shopspring/decimal result
		// valid result should be less than or equal to 1e-19, which is our smallest unit
		d := MustParse(cc.String())
		e := c.Sub(d)

		require.LessOrEqual(t, e.Abs().Cmp(ulp), 0, "a: %s, b: %s, expected %s, got %s", a, b, cc.String(), c.String())
	})
}

func FuzzDiv64(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.lo)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, blo uint64) {
		aprec = aprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		c, err := a.Div64(blo)
		if blo == 0 {
			require.Equal(t, ErrDivideByZero, err)
			return
		}

		require.NoError(t, err)

		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)
		bb := ssDecimal(false, 0, blo, 0)

		prec := int32(c.Prec())
		cc := aa.DivRound(bb, 24).Truncate(prec)

		// sometimes shopspring/decimal does rounding differently
		// e.g. 0.099999999999999 -> 0.1
		// so to check the result, we can check the difference
		// between our result and shopspring/decimal result
		// valid result should have the difference less than or equal to 1e-19, which is our smallest unit
		d := MustParse(cc.String())
		e := c.Sub(d)
		require.LessOrEqual(t, e.Abs().Cmp(ulp), 0)
	})
}

func FuzzQuoRem(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		require.NoError(t, err)

		q, r, err := a.QuoRem(b)

		if b.IsZero() {
			require.Equal(t, ErrDivideByZero, err)
			return
		}

		require.NoError(t, err)

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)
		bb := ssDecimal(bneg, bhi, blo, bprec)

		qq, rr := aa.QuoRem(bb, 0)

		require.Equal(t, qq.String(), q.String(), "quo %s %s", a, b)
		require.Equal(t, rr.String(), r.String(), "rem %s %s", a, b)
	})
}

func FuzzRoundBank(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			roundPrecision := uint8(rand.N(20))
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec, roundPrecision)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8, roundPrecision uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		require.NoError(t, err)

		c := a.Mul(b)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		cstr := c.String()

		u8Round := uint8(roundPrecision % maxPrec)
		cround := c.RoundBank(u8Round)

		// compare with shopspring/decimal
		i32Round := int32(roundPrecision % maxPrec)
		cc := ss.RequireFromString(cstr).RoundBank(i32Round)

		require.Equal(t, cc.String(), cround.String(), "roundBank %s %d", c, roundPrecision)
	})
}

func FuzzRoundAwayFromZero(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			roundPrecision := uint8(rand.N(20))
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec, roundPrecision)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8, roundPrecision uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		require.NoError(t, err)

		c := a.Mul(b)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		cstr := c.String()

		u8Round := uint8(roundPrecision % maxPrec)
		cround := c.RoundAwayFromZero(u8Round)

		// compare with shopspring/decimal
		i32Round := int32(roundPrecision % maxPrec)
		cc := ss.RequireFromString(cstr).RoundUp(i32Round)

		require.Equal(t, cc.String(), cround.String(), "roundAwayFromZero %s %d", c, roundPrecision)
	})
}

func FuzzFloor(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		require.NoError(t, err)

		c := a.Mul(b)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		cstr := c.String()

		cround := c.Floor()

		// compare with shopspring/decimal
		cc := ss.RequireFromString(cstr).Floor()

		require.Equal(t, cc.String(), cround.String(), "floor %s", c)
	})
}

func FuzzCeil(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		require.NoError(t, err)

		c := a.Mul(b)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		cstr := c.String()

		cround := c.Ceil()

		// compare with shopspring/decimal
		cc := ss.RequireFromString(cstr).Ceil()

		require.Equal(t, cc.String(), cround.String(), "ceil %s", c)
	})
}

func FuzzTrunc(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			roundPrecision := uint8(rand.N(20))
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec, roundPrecision)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8, roundPrecision uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		require.NoError(t, err)

		c := a.Mul(b)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		cstr := c.String()

		u8Round := uint8(roundPrecision) % maxPrec
		cround := c.Trunc(u8Round)

		// compare with shopspring/decimal
		i32Round := int32(roundPrecision % maxPrec)
		cc := ss.RequireFromString(cstr).Truncate(i32Round)

		require.Equal(t, cc.String(), cround.String(), "trunc %s %d", c, roundPrecision)
	})
}

func FuzzDepcrecatedPowInt(f *testing.F) {
	for _, c := range corpus {
		f.Add(c.neg, c.hi, c.lo, c.prec, rand.Int())
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, pow int) {
		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		if err == ErrPrecOutOfRange {
			t.Skip()
		} else {
			require.NoError(t, err)
		}

		// use pow less than 10000
		p := pow % 10000

		c := a.PowInt(p)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)

		prec := int32(c.Prec())
		aa = aa.Pow(ss.NewFromInt(int64(p))).Truncate(prec)

		require.Equal(t, aa.String(), c.String(), "powInt %s %d", a, p)
	})
}

func FuzzPowToIntPart(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		require.NoError(t, err)

		// use pow less than 10000
		b, err = b.Mod(MustFromInt64(10000, 0))
		require.NoError(t, err)

		c, err := a.PowToIntPart(b)
		if a.IsZero() && b.IsNeg() {
			require.Equal(t, ErrZeroPowNegative, err, "zero pow negative, %s %s", a, b)
			return
		}

		d := b.Trunc(0)
		if d.coef.overflow() || d.coef.u128.Cmp64(math.MaxInt32) > 0 {
			require.Equal(t, ErrExponentTooLarge, err, "exponent too large, %s %s", a, b)
			return
		}

		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)
		bb := ssDecimal(bneg, bhi, blo, bprec).Mod(ss.NewFromInt(10000))
		aa, err = aa.PowWithPrecision(bb.Truncate(0), int32(c.prec)+4)

		// special case for 0^0
		// udecimal: 0^0 = 1
		// shopspring/decimal: 0^0 is undefined and will return an error
		if a.IsZero() && d.IsZero() {
			require.EqualError(t, err, "cannot represent undefined value of 0**0", "a = %s, b = %s", a, b)
			require.Equal(t, "1", c.String())
			return
		}

		require.NoError(t, err)
		aa = aa.Truncate(int32(c.prec))

		require.Equal(t, aa.String(), c.String(), "PowToIntPart %s %s", a, b)
	})
}

func FuzzPowInt32(f *testing.F) {
	for _, c := range corpus {
		f.Add(c.neg, c.hi, c.lo, c.prec, rand.Int())
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, pow int) {
		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		if err == ErrPrecOutOfRange {
			t.Skip()
		} else {
			require.NoError(t, err)
		}

		// use pow less than 10000
		p := pow % 10000

		c, err := a.PowInt32(int32(p))
		if a.IsZero() && p < 0 {
			require.Equal(t, ErrZeroPowNegative, err)
			return
		}

		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)
		aa, err = aa.PowWithPrecision(ss.New(int64(p), 0), int32(c.prec)+4)

		// special case for 0^0
		// udecimal: 0^0 = 1
		// shopspring/decimal: 0^0 is undefined and will return an error
		if a.IsZero() && p == 0 {
			require.EqualError(t, err, "cannot represent undefined value of 0**0")
			require.Equal(t, "1", c.String())
			return
		}

		require.NoError(t, err)
		aa = aa.Truncate(int32(c.prec))

		require.Equal(t, aa.String(), c.String(), "powInt %s %d", a, p)
	})
}

func FuzzPowNegative(f *testing.F) {
	for _, c := range corpus {
		f.Add(c.neg, c.hi, c.lo, c.prec, rand.Int64())
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, pow int64) {
		aprec = aprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		p := -(pow % int64(maxPrec))

		c := a.PowInt(int(p))
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		if a.IsZero() {
			require.Equal(t, "0", c.String())
			return
		}

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)

		ssPow := ss.NewFromInt(p)
		aa, err = aa.PowWithPrecision(ssPow, int32(c.prec)+8)
		require.NoError(t, err)

		prec := int32(c.Prec())
		aa = aa.Truncate(prec)

		require.Equal(t, aa.String(), c.String(), "powIntNegative %s %d, %s", a, p, ssPow)
	})
}

func FuzzMarshalJSON(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		require.NoError(t, err)

		c := a.Mul(b)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		data, err := c.MarshalJSON()
		require.NoError(t, err)

		var e Decimal
		require.NoError(t, e.UnmarshalJSON(data))

		require.Equal(t, c.String(), e.String())
	})
}

func FuzzMarshalBinary(f *testing.F) {
	for _, c := range corpus {
		for _, d := range corpus {
			f.Add(c.neg, c.hi, c.lo, c.prec, d.neg, d.hi, d.lo, d.prec)
		}
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8, bneg bool, bhi uint64, blo uint64, bprec uint8) {
		aprec = aprec % maxPrec
		bprec = bprec % maxPrec

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		b, err := NewFromHiLo(bneg, bhi, blo, bprec)
		require.NoError(t, err)

		c := a.Mul(b)
		if c.coef.overflow() {
			require.NotNil(t, c.coef.bigInt)
			require.Equal(t, u128{}, c.coef.u128)
		} else {
			require.Nil(t, c.coef.bigInt)
		}

		data, err := c.MarshalBinary()
		require.NoError(t, err)

		var e Decimal
		require.NoError(t, e.UnmarshalBinary(data))

		require.Equal(t, c.String(), e.String())
	})
}

func FuzzLn(f *testing.F) {
	for _, c := range corpus {
		f.Add(c.neg, c.hi, c.lo, c.prec)
	}

	f.Fuzz(func(t *testing.T, aneg bool, ahi uint64, alo uint64, aprec uint8) {
		aprec = aprec % maxPrec
		aneg = false

		a, err := NewFromHiLo(aneg, ahi, alo, aprec)
		require.NoError(t, err)

		if a.IsZero() {
			return
		}

		c, err := a.Ln()
		require.NoError(t, err)
		c = c.trimTrailingZeros()

		// compare with shopspring/decimal
		aa := ssDecimal(aneg, ahi, alo, aprec)
		cc, err := aa.Ln(int32(c.prec))
		require.NoError(t, err)

		d := MustParse(cc.String())
		e := c.Sub(d)

		require.LessOrEqual(t, e.Abs().Cmp(ulp), 0, "ln %s, expected %s, got %s", a, cc.String(), c.String())

	})
}
