package udecimal

import (
	"unsafe"
)

// String returns the string representation of the decimal.
// Trailing zeros will be removed.
func (d Decimal) String() string {
	if d.IsZero() {
		return "0"
	}

	if !d.coef.overflow {
		return d.stringU128(true)
	}

	return d.stringBigInt(true)
}

func (d Decimal) stringBigInt(trimTrailingZeros bool) string {
	str := d.coef.bigInt.String()
	dExpInt := int(d.scale)
	if dExpInt > len(str) {
		// pad with zeros
		l := len(str)
		for i := 0; i < dExpInt-l; i++ {
			str = "0" + str
		}
	}

	var intPart, fractionalPart string
	intPart = str[:len(str)-dExpInt]
	fractionalPart = str[len(str)-dExpInt:]

	if trimTrailingZeros {
		i := len(fractionalPart) - 1
		for ; i >= 0; i-- {
			if fractionalPart[i] != '0' {
				break
			}
		}
		fractionalPart = fractionalPart[:i+1]
	}

	number := intPart
	if number == "" {
		number = "0"
	}

	if len(fractionalPart) > 0 {
		number += "." + fractionalPart
	}

	if d.neg {
		return "-" + number
	}

	return number
}

var (
	// maxByteMap is a map of maximum byte needed to represent an u128 number, indexed by the number of bits.
	maxByteMap = [128]byte{
		1, 1, 1, 1, 2, 2, 2, 3, 3, 3, // 0-9 bits
		4, 4, 4, 4, 5, 5, 5, 6, 6, 6, // 10-19 bits
		7, 7, 7, 7, 8, 8, 8, 9, 9, 9, // 20-29 bits
		10, 10, 10, 10, 11, 11, 11, 12, 12, 12, // 30-39 bits
		13, 13, 13, 13, 14, 14, 14, 15, 15, 15, // 40-49 bits
		16, 16, 16, 16, 17, 17, 17, 18, 18, 18, // 50-59 bits
		19, 19, 19, 19, 20, 20, 20, 21, 21, 21, // 60-69 bits
		22, 22, 22, 22, 23, 23, 23, 24, 24, 24, // 70-79 bits
		25, 25, 25, 25, 26, 26, 26, 27, 27, 27, // 80-89 bits
		28, 28, 28, 28, 29, 29, 29, 30, 30, 30, // 90-99 bits
		31, 31, 31, 32, 32, 32, 32, 33, 33, 33, // 100-109 bits
		34, 34, 34, 35, 35, 35, 35, 36, 36, 36, // 110-119 bits
		37, 37, 37, 38, 38, 38, 38, 39, // 120-127 bits
	}

	digitBytes = [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
)

func (d Decimal) stringU128(trimTrailingZeros bool) string {
	return unsafeBytesToString(d.bytesU128(trimTrailingZeros))
}

// bytesU128 returns the byte representation of the decimal if the coefficient is u128.
func (d Decimal) bytesU128(trimTrailingZeros bool) []byte {
	byteLen := maxByteMap[d.coef.u128.bitLen()]
	var buf []byte
	if trimTrailingZeros {
		// if d.scale > byteLen, that means we need to allocate upto d.scale to cover all the zeros of the fraction part
		// e.g. 0.00000123, scale = 8, byteLen = 3 --> we need to allocate 8 bytes
		if byteLen <= d.scale {
			byteLen = d.scale + 1 // 1 for zero in the whole part
		}

		buf = make([]byte, byteLen+2) // 1 sign + 1 dot
	} else {
		// if not trimming trailing zeros, we can safely allocate 40 bytes
		// 1 sign + 1 dot + len(u128) (which is max to 38 bytes)
		buf = []byte("0000000000000000000000000000000000000000")
	}

	quo, rem := d.coef.u128.QuoRem64(pow10[d.scale].lo) // max scale is 19, so we can safely use QuoRem64
	l := len(buf)
	n := 0
	scale := d.scale

	if rem != 0 {
		if trimTrailingZeros {
			// remove trailing zeros, e.g. 1.2300 -> 1.23
			// both scale and rem will be adjusted
			zeros := getTrailingZeros64(rem)
			rem /= pow10[zeros].lo
			scale -= zeros
		}

		for ; rem != 0; rem /= 10 {
			n++
			buf[l-n] = digitBytes[rem%10]
		}

		// fill remaining zeros
		for i := n + 1; i <= int(scale); i++ {
			buf[l-i] = '0'
		}

		buf[l-1-int(scale)] = '.'
		n = int(scale + 1)
	}

	if quo.IsZero() {
		// quo is zero, so we need to print at least one zero
		n++
		buf[l-n] = '0'
	} else {
		for {
			q, r := quoRem64(quo, 10)

			n++
			buf[l-n] = digitBytes[r]

			if q.IsZero() {
				break
			}

			quo = q
		}
	}

	if d.neg {
		n++
		buf[l-n] = '-'
	}

	return buf[l-n:]
}

func quoRem64(u u128, v uint64) (q u128, r uint64) {
	if u.hi == 0 {
		return u128{lo: u.lo / v}, u.lo % v
	}

	return u.QuoRem64(v)
}

func unsafeBytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// MarshalText implements encoding.TextMarshaler interface.
func (d Decimal) MarshalText() ([]byte, error) {
	if !d.coef.overflow {
		return d.bytesU128(true), nil
	}

	return []byte(d.stringBigInt(true)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler interface.
func (d *Decimal) UnmarshalText(text []byte) error {
	var err error
	*d, err = Parse(string(text))
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (d Decimal) MarshalBinary() ([]byte, error) {
	return nil, nil
}
