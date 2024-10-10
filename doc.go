// Package udecimal provides a high-performance, zero-allocation, high-precision
// (up to 19 digits after the decimal point) decimal arithmetic library.
// It includes functions for parsing and performing arithmetic
// operations such as addition, subtraction, multiplication, and division
// on decimal numbers. The package is designed to handle decimal numbers
// with a high degree of precision and efficiency, making it suitable for
// high-traffic financial applications where both precision and performance are critical.
//
// Maximum and default precision is 19 digits after the decimal point. The default precision
// can be changed globally to any value between 1 and 19 to suit your use case and make sure
// that the precision is consistent across the entire application.
//
// # How it works
//
// The Decimal type represents a fixed-point decimal number. It is composed of sign, coefficient, and scale,
// where the number is represented as:
//
//	number = (neg ? -1 : 1) * coef / 10^(scale)
//	e.g. 123.456 = 123456 / 10^3 -> neg = false, coef = 123456, scale = 3
//	    -123.456 = -123456 / 10^3 -> neg = true, coef = 123456, scale = 3
//
// Fields:
// - neg: A boolean indicating if the number is negative.
// - scale: The scale of the decimal number, representing the number of digits after the decimal point (up to 19). The scale is always non-negative.
// - coef: The coefficient of the decimal number. The coefficient is always non-negative and is stored in a special format that allows for efficient arithmetic operations.
//
// # Codec
//
// The udecimal package supports various encoding and decoding mechanisms to facilitate easy integration with
// different data storage and transmission systems.
//
// - Marshal/UnmarshalText: json, string
// - Marshal/UnmarshalBinary: gob, protobuf
// - SQL: The Decimal type implements the sql.Scanner interface, enabling seamless integration with SQL databases.
// - DynamoDB: The package supports parsing DynamoDB number (regarless number or string) to Decimal and marshal Decimal back to DynamoDB number.
// For more details, see the documentation for each method.
package udecimal
