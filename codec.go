package udecimal

import (
	"unsafe"
)

// func (d Decimal) Scan(value interface{}) error {
// 	return nil
// }

// func (d Decimal) MarshalText() ([]byte, error) {
// 	buf := []byte("0000000000000000000000000000000000000000")
// 	n := d.writeToBytes(buf, true)
// 	return buf[n:], nil
// }

// func (d Decimal) UnmarshalText(text []byte) error {

// 	return nil
// }

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

	var intPart, fractionalPart string

	// NOTE(vadim): this cast to int will cause bugs if d.exp == INT_MIN
	// and you are on a 32-bit machine. Won't fix this super-edge case.
	dExpInt := int(d.scale)
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
	if len(fractionalPart) > 0 {
		number += "." + fractionalPart
	}

	if d.neg {
		return "-" + number
	}

	return number

}
func (d Decimal) stringU128(trimTrailingZeros bool) string {
	// max 40 bytes: 1 sign + 19 whole + 1 dot + 19 fraction
	buf := []byte("0000000000000000000000000000000000000000")

	quo, rem, _ := d.coef.u128.QuoRem(pow10[d.scale])
	l := len(buf)
	n := 0
	scale := d.scale

	if !rem.IsZero() {
		if trimTrailingZeros {
			// remove trailing zeros, e.g. 1.2300 -> 1.23
			// both scale and rem will be adjusted
			zeros := getTrailingZeros(rem)
			rem, _, _ = rem.QuoRem(pow10[zeros])
			scale -= zeros
		}

		for {
			q, r := rem.QuoRem64(10)
			n++
			buf[l-n] += byte(r)
			if q.IsZero() {
				break
			}

			rem = q
		}

		buf[l-1-int(scale)] = '.'
		n = int(scale + 1)
	}

	if quo.IsZero() {
		// quo is zero, so we need to print at least one zero
		n++
	} else {
		for {
			q, r := quo.QuoRem64(10)
			n++
			buf[l-n] += byte(r)
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

	return unsafeBytesToString(buf[l-n:])
}

func unsafeBytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
