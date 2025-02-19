package udecimal

import "fmt"

var (
	// internal error
	errOverflow = fmt.Errorf("overflow")
)

var (
	// ErrPrecOutOfRange is returned when the decimal precision is greater than the default precision
	// default precision can be configured using SetDefaultPrecision, and its value is up to 19
	ErrPrecOutOfRange = fmt.Errorf("precision out of range. Only support maximum %d digits after the decimal point", defaultPrec)

	// ErrEmptyString is returned when the input string is empty
	ErrEmptyString = fmt.Errorf("can't parse empty string")

	// ErrMaxStrLen is returned when the input string exceeds the maximum length
	// Maximum length is arbitrarily set to 200 so string length value can fit in 1 byte (for MarshalBinary).
	// Also such that big number (more than 200 digits) is unrealistic in financial system
	// which this library is mainly designed for.
	ErrMaxStrLen = fmt.Errorf("string input exceeds maximum length %d", maxStrLen)

	// ErrInvalidFormat is returned when the input string is not in the correct format
	// It doesn't support scientific notation, such as 1e-2, 1.23e4, etc.
	ErrInvalidFormat = fmt.Errorf("invalid format")

	// ErrDivideByZero is returned when dividing by zero
	ErrDivideByZero = fmt.Errorf("can't divide by zero")

	// ErrSqrtNegative is returned when calculating square root of negative number
	ErrSqrtNegative = fmt.Errorf("can't calculate square root of negative number")

	// ErrInvalidBinaryData is returned when unmarshalling invalid binary data
	// The binary data should follow the format as described in MarshalBinary
	ErrInvalidBinaryData = fmt.Errorf("invalid binary data")

	// ErrZeroPowNegative is returned when raising zero to a negative power
	ErrZeroPowNegative = fmt.Errorf("can't raise zero to a negative power")

	// ErrExponentTooLarge is returned when the exponent is too large and becomes impractical.
	ErrExponentTooLarge = fmt.Errorf("exponent is too large. Must be less than or equal math.MaxInt32")

	// ErrIntPartOverflow is returned when the integer part of the decimal is too large to fit in int64
	ErrIntPartOverflow = fmt.Errorf("integer part is too large to fit in int64")

	// ErrLnNonPositive is returned when calculating natural logarithm of non-positive number
	ErrLnNonPositive = fmt.Errorf("can't calculate natural logarithm of non-positive number")
)
