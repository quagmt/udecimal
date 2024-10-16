package udecimal

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"math/bits"
	"unsafe"
)

var (
	_ fmt.Stringer             = (*Decimal)(nil)
	_ sql.Scanner              = (*Decimal)(nil)
	_ driver.Valuer            = (*Decimal)(nil)
	_ encoding.TextMarshaler   = (*Decimal)(nil)
	_ encoding.TextUnmarshaler = (*Decimal)(nil)
	_ json.Marshaler           = (*Decimal)(nil)
	_ json.Unmarshaler         = (*Decimal)(nil)
)

// String returns the string representation of the decimal.
// Trailing zeros will be removed.
func (d Decimal) String() string {
	if d.IsZero() {
		return "0"
	}

	if !d.coef.overflow() {
		return d.stringU128(true, false)
	}

	return d.stringBigInt(true)
}

// StringFixed returns the string representation of the decimal with fixed prec.
// Trailing zeros will not be removed.
//
// Special case: if the decimal is zero, it will return "0" regardless of the prec.
func (d Decimal) StringFixed(prec uint8) string {
	d1 := d.rescale(prec)

	if !d1.coef.overflow() {
		return d1.stringU128(false, false)
	}

	return d1.stringBigInt(false)
}

func (d Decimal) stringBigInt(trimTrailingZeros bool) string {
	str := d.coef.bigInt.String()
	dExpInt := int(d.prec)
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
	maxByteMap = [129]byte{
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
		37, 37, 37, 38, 38, 38, 38, 39, 39, // 120-128 bits
	}

	digitBytes = [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
)

func (d Decimal) stringU128(trimTrailingZeros bool, withQuote bool) string {
	return unsafeBytesToString(d.bytesU128(trimTrailingZeros, withQuote))
}

// bytesU128 returns the byte representation of the decimal if the coefficient is u128.
func (d Decimal) bytesU128(trimTrailingZeros bool, withQuote bool) []byte {
	var totalLen uint8
	byteLen := maxByteMap[d.coef.u128.bitLen()]

	if trimTrailingZeros {
		// if d.prec > byteLen, that means we need to allocate upto d.prec to cover all the zeros of the fraction part
		// e.g. 0.00000123, prec = 8, byteLen = 3 --> we need to allocate 8 bytes for the fraction part
		if byteLen <= d.prec {
			byteLen = d.prec + 1 // 1 for zero in the whole part
		}

		totalLen = byteLen + 2
	} else {
		// if not trimming trailing zeros, we can safely allocate 41 bytes
		// 1 sign + 1 dot + len(u128) (which is max to 39 bytes)
		// buf = []byte("00000000000000000000000000000000000000000")
		totalLen = 41
	}

	if withQuote {
		// if withQuote is true, we need to add quotes at the beginning and the end
		totalLen += 2
		buf := make([]byte, totalLen)
		n := d.fillBuffer(buf[1:len(buf)-1], trimTrailingZeros)

		n += 2 // 1 for quote offset at buf[l-1], 1 for moving the index to next position
		l := len(buf)
		buf[l-1], buf[l-n] = '"', '"'
		return buf[l-n:]
	}

	buf := make([]byte, totalLen)
	n := d.fillBuffer(buf, trimTrailingZeros)

	return buf[len(buf)-n:]
}

func (d Decimal) fillBuffer(buf []byte, trimTrailingZeros bool) int {
	quo, rem := d.coef.u128.QuoRem64(pow10[d.prec].lo) // max prec is 19, so we can safely use QuoRem64

	prec := d.prec
	l := len(buf)
	n := 0

	if rem != 0 {
		if trimTrailingZeros {
			// remove trailing zeros, e.g. 1.2300 -> 1.23
			// both prec and rem will be adjusted
			zeros := getTrailingZeros64(rem)
			rem /= pow10[zeros].lo
			prec -= zeros
		}

		for ; rem != 0; rem /= 10 {
			n++
			buf[l-n] = digitBytes[rem%10]
		}

		// fill remaining zeros
		for i := n + 1; i <= int(prec); i++ {
			buf[l-i] = '0'
		}

		buf[l-1-int(prec)] = '.'
		n = int(prec + 1)
	}

	if quo.IsZero() {
		// quo is zero, we need to print at least one zero
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

	return n
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

func unssafeStringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// MarshalJSON implements the [json.Marshaler] interface.
func (d Decimal) MarshalJSON() ([]byte, error) {
	if !d.coef.overflow() {
		return d.bytesU128(true, true), nil
	}

	return []byte(`"` + d.stringBigInt(true) + `"`), nil
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (d *Decimal) UnmarshalJSON(data []byte) error {
	// Remove quotes if they exist.
	if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}

	return d.UnmarshalText(data)
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (d Decimal) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (d *Decimal) UnmarshalText(data []byte) error {
	var err error
	*d, err = parseBytes(data)
	return err
}

// MarshalBinary implements [encoding.BinaryMarshaler] interface with custom binary format.
//
//	Binary format: [overflow + neg] [prec] [total bytes] [coef]
//
//	 example 1: -1.2345
//	 1st byte: 0b0001_0000 (overflow = true, neg = false)
//	 2nd byte: 0b0000_0100 (prec = 4)
//	 3rd byte: 0b0000_1101 (total bytes = 11)
//	 4th-11th bytes: 0x0000_0000_0000_3039 (coef = 12345, only stores the coef.lo part)
//
//	 example 2: 1234567890123456789.1234567890123456789
//	 1st byte: 0b0000_0000 (overflow = false, neg = false)
//	 2nd byte: 0b0001_0011 (prec = 19)
//	 3rd byte: 0b0001_0011 (total bytes = 19)
//	 4th-11th bytes: 0x0949_b0f6_f002_3313 (coef.hi)
//	 12th-19th bytes: 0xd3b5_05f9_b5f1_8115 (coef.lo)
func (d Decimal) MarshalBinary() ([]byte, error) {
	if !d.coef.overflow() {
		return d.marshalBinaryU128()
	}

	return d.marshalBinaryBigInt()
}

func (d Decimal) marshalBinaryU128() ([]byte, error) {
	coef := d.coef.u128
	totalBytes := 19

	if coef.hi == 0 {
		totalBytes = 11
	}

	buf := make([]byte, totalBytes)
	var neg int
	if d.neg {
		neg = 1
	}

	// overflow + neg with overflow = false (always 0)
	buf[0] = byte(neg)
	buf[1] = byte(d.prec)
	buf[2] = byte(totalBytes)

	if coef.hi != 0 {
		copyUint64ToBytes(buf[3:], coef.hi)
		copyUint64ToBytes(buf[11:], coef.lo)
	} else {
		copyUint64ToBytes(buf[3:], coef.lo)
	}

	return buf, nil
}

func copyUint64ToBytes(b []byte, n uint64) {
	// use big endian to make it consistent with big.Int.FillBytes, which also uses big endian
	binary.BigEndian.PutUint64(b, n)
}

func (d *Decimal) UnmarshalBinary(data []byte) error {
	if len(data) < 3 {
		return ErrInvalidBinaryData
	}

	overflow := data[0] >> 4 & 1
	if overflow == 0 {
		return d.unmarshalBinaryU128(data)
	}

	return d.unmarshalBinaryBigInt(data)
}

func (d *Decimal) unmarshalBinaryU128(data []byte) error {
	d.neg = data[0]&1 == 1
	d.prec = data[1]

	totalBytes := data[2]

	// for u128, totalBytes must be 11 or 19
	if totalBytes != 11 && totalBytes != 19 {
		return ErrInvalidBinaryData
	}

	coef := u128{}
	if totalBytes == 11 {
		coef.lo = binary.BigEndian.Uint64(data[3:])
	} else {
		coef.hi = binary.BigEndian.Uint64(data[3:])
		coef.lo = binary.BigEndian.Uint64(data[11:])
	}

	d.coef.u128 = coef
	return nil
}

func (d *Decimal) unmarshalBinaryBigInt(data []byte) error {
	d.neg = data[0]&1 == 1
	d.prec = data[1]

	totalBytes := data[2]

	if totalBytes < 3 {
		return ErrInvalidBinaryData
	}

	d.coef.bigInt = new(big.Int).SetBytes(data[3:totalBytes])
	return nil
}

func (d Decimal) marshalBinaryBigInt() ([]byte, error) {
	var neg int
	if d.neg {
		neg = 1
	}

	if d.coef.bigInt == nil {
		return nil, ErrInvalidBinaryData
	}

	words := d.coef.bigInt.Bits()
	totalBytes := 3 + len(words)*(bits.UintSize/8)
	buf := make([]byte, totalBytes)

	// overflow + neg with overflow = true (always 1)
	buf[0] = byte(1<<4 | neg)
	buf[1] = byte(d.prec)
	buf[2] = byte(totalBytes)
	d.coef.bigInt.FillBytes(buf[3:])

	return buf, nil
}

// Scan implements sql.Scanner interface.
func (d *Decimal) Scan(src any) error {
	var err error
	switch v := src.(type) {
	case []byte:
		*d, err = Parse(unsafeBytesToString(v))
	case string:
		*d, err = Parse(v)
	case uint64:
		*d, err = NewFromUint64(v, 0)
	case int64:
		*d, err = NewFromInt64(v, 0)
	case int:
		*d, err = NewFromInt64(int64(v), 0)
	case int32:
		*d, err = NewFromInt64(int64(v), 0)
	case float64:
		*d, err = NewFromFloat64(v)
	case nil:
		err = fmt.Errorf("can't scan nil to Decimal")
	default:
		err = fmt.Errorf("can't scan %T to Decimal: %T is not supported", src, src)
	}

	return err
}

// Value implements [driver.Valuer] interface.
func (d Decimal) Value() (driver.Value, error) {
	return d.String(), nil
}

// NullDecimal is a nullable Decimal.
type NullDecimal struct {
	Decimal Decimal
	Valid   bool
}

// Scan implements [sql.Scanner] interface.
//
// [sql.Scanner]: https://pkg.go.dev/database/sql#Scanner
func (d *NullDecimal) Scan(src any) error {
	if src == nil {
		d.Decimal, d.Valid = Decimal{}, false
		return nil
	}

	var err error
	switch v := src.(type) {
	case []byte:
		d.Decimal, err = Parse(string(v))
	case string:
		d.Decimal, err = Parse(v)
	case uint64:
		d.Decimal, err = NewFromUint64(v, 0)
	case int64:
		d.Decimal, err = NewFromInt64(v, 0)
	case int:
		d.Decimal, err = NewFromInt64(int64(v), 0)
	case int32:
		d.Decimal, err = NewFromInt64(int64(v), 0)
	case float64:
		d.Decimal, err = NewFromFloat64(v)
	default:
		err = fmt.Errorf("can't scan %T to Decimal: %T is not supported", src, src)
	}

	d.Valid = err == nil
	return err
}

// Value implements the [driver.Valuer] interface.
//
// [driver.Valuer]: https://pkg.go.dev/database/sql/driver#Valuer
func (d NullDecimal) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}

	return d.Decimal.String(), nil
}
