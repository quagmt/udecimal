package udecimal

import (
	"fmt"
	"math/big"
	"strings"
)

var (
	bigZero = big.NewInt(0)
	bigOne  = big.NewInt(1)
	bigTen  = big.NewInt(10)
)

// bint represents a whole decimal number without a decimal point.
// The value is always positive and is stored in a 128-bit unsigned integer.
// If the value exceeds the 128-bit limit, it falls back to using big.Int.
type bint struct {
	// Indicates if the value has overflowed and is stored in big.Int.
	overflow bool

	// Stores small numbers with high performance and zero allocation operations.
	// The value range is 0 <= u128 <= 2^128 - 1
	u128 u128

	// fall back, in case the value is our of u128 range
	bigInt *big.Int
}

func bintFromBigInt(b *big.Int) bint {
	return bint{overflow: true, bigInt: b}
}

func bintFromU128(u u128) bint {
	return bint{u128: u}
}

func bintFromU64(u uint64) bint {
	return bint{u128: u128{lo: u}}
}

func (u bint) GetBig() *big.Int {
	if u.overflow {
		return new(big.Int).Set(u.bigInt)
	}

	return u.u128.ToBigInt()
}

func (u bint) IsZero() bool {
	if !u.overflow {
		return u.u128.IsZero()
	}

	return u.bigInt.Cmp(bigZero) == 0
}

func (u bint) Cmp(v bint) int {
	if !u.overflow && !v.overflow {
		return u.u128.Cmp(v.u128)
	}

	return u.GetBig().Cmp(v.GetBig())
}

func errInvalidFormat(s []byte) error {
	return fmt.Errorf("%w: can't parse '%s' to Decimal", ErrInvalidFormat, s)
}

func parseBint(s []byte) (bool, bint, uint8, error) {
	if len(s) > maxStrLen {
		return false, bint{}, 0, ErrMaxStrLen
	}

	// if s has less than 40 characters, it can fit into u128
	if len(s) <= 40 {
		neg, bint, scale, err := parseBintFromU128(s)
		if err == nil || err != errOverflow {
			return neg, bint, scale, err
		}

		// overflow, try to parse into big.Int
	}

	// parse into big.Int
	var (
		width      = len(s)
		intString  string
		scale, pos int
		neg        bool
		value      = s
	)

	switch s[0] {
	case '.':
		return false, bint{}, 0, errInvalidFormat(s)
	case '-':
		neg = true
		value = s[1:]
		pos++
	case '+':
		pos++
	default:
		// do nothing
	}

	// prevent "+" or "-"
	if pos == width {
		return false, bint{}, 0, errInvalidFormat(s)
	}

	// prevent "-.123" or "+.123"
	if s[pos] == '.' {
		return false, bint{}, 0, errInvalidFormat(s)
	}

	pIndex := -1
	vLen := len(value)
	for i := 0; i < vLen; i++ {
		if value[i] == '.' {
			if pIndex > -1 {
				// input has more than 1 decimal point
				return false, bint{}, 0, errInvalidFormat(s)
			}
			pIndex = i
		}
	}

	switch {
	case pIndex == -1:
		// There is no decimal point, we can just parse the original string as an int
		intString = string(value)
	case pIndex >= vLen-1:
		// prevent "123." or "-123."
		return false, bint{}, 0, errInvalidFormat(s)
	default:
		b := strings.Builder{}
		_, err := b.Write(value[:pIndex])
		if err != nil {
			return false, bint{}, 0, err
		}

		_, err = b.Write(value[pIndex+1:])
		if err != nil {
			return false, bint{}, 0, err
		}

		// intString = value[:pIndex] + value[pIndex+1:]
		intString = b.String()
		scale = len(value[pIndex+1:])
	}

	if scale > int(defaultScale) {
		return false, bint{}, 0, ErrScaleOutOfRange
	}

	dValue := new(big.Int)
	_, ok := dValue.SetString(intString, 10)
	if !ok {
		return false, bint{}, 0, errInvalidFormat(s)
	}

	// the value should always be positive, as we already extracted the sign
	if dValue.Sign() == -1 {
		return false, bint{}, 0, errInvalidFormat(s)
	}

	// nolint: gosec
	return neg, bintFromBigInt(dValue), uint8(scale), nil
}

func parseBintFromU128(s []byte) (bool, bint, uint8, error) {
	width := len(s)

	var (
		pos int
		neg bool
	)

	switch s[0] {
	case '.':
		return false, bint{}, 0, errInvalidFormat(s)
	case '-':
		neg = true
		pos++
	case '+':
		pos++
	default:
		// do nothing
	}

	// prevent "+" or "-"
	if pos == width {
		return false, bint{}, 0, errInvalidFormat(s)
	}

	// prevent "-.123" or "+.123"
	if s[pos] == '.' {
		return false, bint{}, 0, errInvalidFormat(s)
	}

	var (
		err   error
		coef  u128
		scale uint8
	)
	for ; pos < width; pos++ {
		if s[pos] == '.' {
			// return err if we encounter the '.' more than once
			if scale != 0 {
				return false, bint{}, 0, errInvalidFormat(s)
			}

			// nolint: gosec
			scale = uint8(width - pos - 1)

			// prevent "123." or "-123."
			if scale == 0 {
				return false, bint{}, 0, errInvalidFormat(s)
			}

			if scale > defaultScale {
				return false, bint{}, 0, ErrScaleOutOfRange
			}

			continue
		}

		if s[pos] < '0' || s[pos] > '9' {
			return false, bint{}, 0, errInvalidFormat(s)
		}

		coef, err = coef.Mul64(10)
		if err != nil {
			return false, bint{}, 0, err
		}

		coef, err = coef.Add64(uint64(s[pos] - '0'))
		if err != nil {
			return false, bint{}, 0, err
		}
	}

	if coef.IsZero() {
		return false, bint{}, 0, nil
	}

	return neg, bint{u128: coef}, scale, nil
}

// GT returns true if u > v
func (u bint) GT(v bint) bool {
	return u.Cmp(v) == 1
}

func (u bint) Add(v bint) bint {
	if !u.overflow && !v.overflow {
		c, err := u.u128.Add(v.u128)
		if err == nil {
			return bint{u128: c}
		}

		// overflow, fallback to big.Int
	}

	return bintFromBigInt(new(big.Int).Add(u.GetBig(), v.GetBig()))
}

func (u bint) Sub(v bint) (bint, error) {
	if !u.overflow && !v.overflow {
		c, err := u.u128.Sub(v.u128)
		if err == nil {
			return bint{u128: c}, nil
		}
	}

	uBig := u.GetBig()
	vBig := v.GetBig()

	// make sure the result is always positive
	if uBig.Cmp(vBig) < 0 {
		return bint{}, errOverflow
	}

	return bintFromBigInt(new(big.Int).Sub(uBig, vBig)), nil
}

func (u bint) Mul(v bint) bint {
	if !u.overflow && !v.overflow {
		c, err := u.u128.Mul(v.u128)
		if err == nil {
			return bint{u128: c}
		}
	}

	return bintFromBigInt(new(big.Int).Mul(u.GetBig(), v.GetBig()))
}
