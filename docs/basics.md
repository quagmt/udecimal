# Basics

## Installation

```sh
go get github.com/quagmt/udecimal
```

udecimal requires Go 1.23 or later.

## Usage

--8<-- "examples/udecimal/Example.md"

`Add`, `Sub` and `Mul` return a `Decimal`. `Div` returns `(Decimal, error)`, because the divisor can be zero.
`123.456 / 7` has more than 19 decimal places, so the result is truncated to 19. See
[Precision and rounding](precision.md).

## Creating values

| Function | Example | Result |
|---|---|---|
| `Parse(s)` | `udecimal.Parse("-0.00007890123456789")` | `-0.00007890123456789`, `nil` |
| `NewFromInt64(coef, prec)` | `udecimal.NewFromInt64(-12345, 2)` | `-123.45`, `nil` |
| `NewFromUint64(coef, prec)` | `udecimal.NewFromUint64(12345, 2)` | `123.45`, `nil` |
| `NewFromFloat64(f)` | `udecimal.NewFromFloat64(1.234)` | `1.234`, `nil` |
| `NewFromHiLo(neg, hi, lo, prec)` | `udecimal.NewFromHiLo(true, 0, 123456, 4)` | `-12.3456`, `nil` |

`prec` is the number of decimal places. `Parse` and the functions that take `prec` return `ErrPrecOutOfRange` for
more decimal places than the [default precision](precision.md#default-precision), which is 19 unless changed.
`MustParse`, `MustFromInt64`, `MustFromUint64` and `MustFromFloat64` panic instead of returning an error. `Zero`
and `One` are predefined values, and the zero value of `Decimal` is 0.

!!! warning "Floats"
    `NewFromFloat64` converts a binary floating-point number, which may not hold every digit you wrote:
    `udecimal.NewFromFloat64(123456789.1234567890123456789)` is `123456789.12345679`. Use `Parse` when the exact
    value matters.

See [Parsing](parsing-and-encoding.md#parsing) for the formats `Parse` accepts.

## Arithmetic

| Method | Returns | Returns an error when |
|---|---|---|
| `Add`, `Sub`, `Mul` | `Decimal` | never |
| `Add64`, `Sub64`, `Mul64` | `Decimal`; the operand is a `uint64` | never |
| `Div`, `Div64` | `(Decimal, error)` | the divisor is zero |
| `QuoRem` | `(quotient, remainder Decimal, error)`; the quotient is an integer | the divisor is zero |
| `Mod` | `(Decimal, error)`, the remainder of `QuoRem` | the divisor is zero |
| `PowInt32(e)` | `(Decimal, error)` | 0 is raised to a negative power |
| `PowToIntPart(e)` | `(Decimal, error)`; raises to the integer part of the `Decimal` `e` | 0 is raised to a negative power, or the integer part of `e` is beyond ±`math.MaxInt32` |
| `Sqrt` | `(Decimal, error)` | the value is negative |

`MustDiv`, `MustMod` and `MustQuoRem` panic instead of returning an error. `PowInt` is deprecated; use `PowInt32`,
which handles 0^0 and negative exponents correctly.

--8<-- "examples/udecimal/ExampleDecimal_QuoRem.md"

--8<-- "examples/udecimal/ExampleDecimal_Sqrt.md"

## Comparison

| Method | Returns |
|---|---|
| `a.Cmp(b)` | -1, 0 or 1 |
| `Equal`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan`, `LessThanOrEqual` | `bool` |
| `Sign` | -1, 0 or 1 |
| `IsZero`, `IsNeg`, `IsPos` | `bool` |
| `udecimal.Max(a, b...)`, `udecimal.Min(a, b...)` | the largest or smallest value |
| `udecimal.Avg(a, b...)` | the mean, with at most 19 decimal places, like `Div` |

`Equal` compares values, so `1.50` equals `1.5`.

--8<-- "examples/udecimal/ExampleDecimal_Cmp.md"

--8<-- "examples/udecimal/ExampleAvg.md"

## Sign and precision

`Abs` returns the absolute value and `Neg` the negated value. `a.CopySign(b)` returns `a` with the sign of `b`.
`udecimal.Parts(d)` returns the integer part and the fractional part of `d`, which add up to `d`. The fractional part
has the sign of `d`, so `Parts(-0.5)` returns `0` and `-0.5`. `Prec` (an `int`) and `PrecUint` (a `uint8`) return
the number of decimal places.

--8<-- "examples/udecimal/ExampleParts.md"

--8<-- "examples/udecimal/ExampleDecimal_CopySign.md"

## Conversion

| Method | Returns |
|---|---|
| `String` | the shortest form, without trailing zeros |
| `StringFixed(prec)` | at least `prec` decimal places, padded with zeros; it doesn't round |
| `Int64` | the integer part as an `int64`; an error if it doesn't fit |
| `InexactFloat64` | the nearest `float64` |
| `ToHiLo` | the sign, the 128-bit coefficient as two `uint64` values, and the precision; `ok` is false if the coefficient doesn't fit in 128 bits |

--8<-- "examples/udecimal/ExampleDecimal_StringFixed.md"

--8<-- "examples/udecimal/ExampleDecimal_Int64.md"

## Errors

All errors are package variables, so you can check them with `errors.Is`:

| Error | Returned by |
|---|---|
| `ErrPrecOutOfRange` | `Parse`, `NewFromInt64`, `NewFromUint64`, `NewFromHiLo` and `NewFromFloat64`, for more decimal places than the default precision |
| `ErrEmptyString` | `Parse("")` |
| `ErrMaxStrLen` | `Parse` of a string longer than 200 characters |
| `ErrInvalidFormat` | `Parse` of a string that isn't a decimal number, and `NewFromFloat64` of NaN or an infinity |
| `ErrDivideByZero` | `Div`, `Div64`, `Mod` and `QuoRem` with a zero divisor |
| `ErrZeroPowNegative` | `PowInt32` and `PowToIntPart` when 0 is raised to a negative power |
| `ErrExponentTooLarge` | `PowToIntPart` when the integer part of the exponent is beyond ±`math.MaxInt32` |
| `ErrSqrtNegative` | `Sqrt` of a negative value |
| `ErrIntPartOverflow` | `Int64` when the integer part doesn't fit |
| `ErrInvalidBinaryData` | `UnmarshalBinary` of data in the wrong layout, or with more decimal places than the default precision |

The text decoders (`UnmarshalJSON`, `UnmarshalText` and `Scan`) wrap the errors from `Parse`, and `NewFromFloat64`
wraps its own errors, so use `errors.Is` instead of `==`.

--8<-- "examples/udecimal/Example_errors.md"

Apart from the `Must*` functions, only `SetDefaultPrecision` and `SetDefaultParseMode` panic, and only for an
invalid argument.
