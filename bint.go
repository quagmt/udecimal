package udecimal

import (
	"fmt"
	"math/big"
)

var (
	zero = big.NewInt(0)
)

type bint struct {
	// flag to indicate if the value is overflow and stored in big.Int
	overflow bool

	// use for storing small number, with high performance
	u128 u128

	// fall back
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
		return u.bigInt
	}

	return u.u128.ToBigInt()
}

func (u bint) IsZero() bool {
	if !u.overflow {
		return u.u128.IsZero()
	}

	return u.bigInt.Cmp(zero) == 0
}

func parseBint(s string) (bool, bint, uint8, error) {
	// if s has less than 40 characters, it can fit into u128
	if len(s) <= 40 {
		neg, bint, scale, err := parseBintFromU128(s)
		if err == nil || err != ErrOverflow {
			return neg, bint, scale, err
		}

		// overflow, try to parse into big.Int
	}

	// parse into big.Int
	errInvalidFormat := fmt.Errorf("%w: can't parse '%s' to decimal", ErrInvalidFormat, s)

	var (
		width      = len(s)
		intString  string
		scale, pos int
		neg        bool
		value      = s
	)

	switch s[0] {
	case '.':
		return false, bint{}, 0, errInvalidFormat
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
		return false, bint{}, 0, errInvalidFormat
	}

	// prevent "-.123" or "+.123"
	if s[pos] == '.' {
		return false, bint{}, 0, errInvalidFormat
	}

	pIndex := -1
	vLen := len(value)
	for i := 0; i < vLen; i++ {
		if value[i] == '.' {
			if pIndex > -1 {
				return false, bint{}, 0, errInvalidFormat
			}
			pIndex = i
		}
	}

	if pIndex == -1 {
		// There is no decimal point, we can just parse the original string as
		// an int
		intString = value
	} else {
		if pIndex+1 < vLen {
			intString = value[:pIndex] + value[pIndex+1:]
		} else {
			intString = value[:pIndex]
		}

		scale = len(value[pIndex+1:])
	}

	if scale > maxScale {
		return false, bint{}, 0, ErrMaxScale
	}

	// strconv.ParseInt is faster than new(big.Int).SetString so this is just a shortcut for strings we know won't overflow
	dValue := new(big.Int)
	_, ok := dValue.SetString(intString, 10)
	if !ok {
		return false, bint{}, 0, errInvalidFormat
	}

	return neg, bintFromBigInt(dValue), uint8(scale), nil
}

func parseBintFromU128(s string) (bool, bint, uint8, error) {
	errInvalidFormat := fmt.Errorf("%w: can't parse '%s' to decimal", ErrInvalidFormat, s)
	width := len(s)

	var (
		pos int
		neg bool
	)

	switch s[0] {
	case '.':
		return false, bint{}, 0, errInvalidFormat
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
		return false, bint{}, 0, errInvalidFormat
	}

	// prevent "-.123" or "+.123"
	if s[pos] == '.' {
		return false, bint{}, 0, errInvalidFormat
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
				return false, bint{}, 0, errInvalidFormat
			}

			scale = uint8(width - pos - 1)

			// prevent "123." or "-123."
			if scale == 0 {
				return false, bint{}, 0, errInvalidFormat
			}

			if scale > maxScale {
				return false, bint{}, 0, ErrMaxScale
			}

			continue
		}

		if s[pos] < '0' || s[pos] > '9' {
			return false, bint{}, 0, errInvalidFormat
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

	if coef.isOverflow() {
		return false, bint{}, 0, ErrOverflow
	}

	return neg, bint{u128: coef}, scale, nil
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

func (u bint) Mul(v bint) bint {
	if !u.overflow && v.overflow {
		c, err := u.u128.Mul(v.u128)
		if err == nil {
			return bint{u128: c}
		}
	}

	return bintFromBigInt(new(big.Int).Mul(u.GetBig(), v.GetBig()))
}
