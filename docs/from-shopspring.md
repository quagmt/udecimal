# Migrating from shopspring

This page describes how to move from [shopspring/decimal](https://github.com/shopspring/decimal) v1.6.0 to udecimal.
Most method names are the same. The main differences are that udecimal returns errors where shopspring panics,
keeps at most 19 decimal places, truncates instead of rounding, and accepts fewer input formats.

The examples and the [API mapping](#api-mapping) are tested against both libraries.

## Behavior differences

### Division by zero

`Div`, `Mod` and `QuoRem` return an error when the divisor is zero, where shopspring panics. To panic instead, use
`MustDiv`, `MustMod` or `MustQuoRem`.

--8<-- "examples/migration/Example_divideByZero.md"

### Decimal places

shopspring keeps every digit of its input. udecimal's `Parse` returns an error for more than 19 decimal places, or
truncates them in [`ParseModeTrunc`](parsing-and-encoding.md#parse-modes).

--8<-- "examples/migration/Example_decimalPlaces.md"

### Input formats

shopspring's `NewFromString` also accepts exponents (`1e5`, `1.5E-3`) and a decimal point without a digit on one
side (`.5`, `1.`). udecimal's `Parse` returns `ErrInvalidFormat` for these. shopspring's `Scan` also removes the
quotes from a quoted string such as `"1.5"`; udecimal's `Scan` returns an error.

--8<-- "examples/migration/Example_inputFormat.md"

### Division

shopspring's `Div` rounds the quotient half away from zero to `DivisionPrecision` decimal places, 16 by default.
udecimal's `Div` keeps 19 decimal places and truncates the rest.

--8<-- "examples/migration/Example_division.md"

### Multiplication

shopspring's `Mul` is exact. udecimal's `Mul` keeps at most 19 decimal places and truncates the rest.

--8<-- "examples/migration/Example_multiplication.md"

### `StringFixed`

shopspring's `StringFixed(2)` rounds to 2 decimal places. udecimal's `StringFixed` only pads with zeros, so round the
value first if it has more decimal places.

--8<-- "examples/migration/Example_stringFixed.md"

### `PowInt32`

For 0^0, shopspring returns an error and udecimal returns 1. For a negative exponent, shopspring keeps
`PowPrecisionNegativeExponent` decimal places, 16 by default, and udecimal keeps 19. For a positive exponent,
shopspring is exact and udecimal keeps 19 decimal places.

--8<-- "examples/migration/Example_powInt32.md"

### Rounding method names

`Round` is `RoundHAZ` (half away from zero), `Truncate` and `RoundDown` are `Trunc`, and `RoundUp` is
`RoundAwayFromZero`. udecimal also has `RoundHTZ` (half toward zero).

--8<-- "examples/migration/Example_rounding.md"

### `NullDecimal` encoding

udecimal's `NullDecimal` implements only `Scan` and `Value`. shopspring's also implements `MarshalJSON`,
`UnmarshalJSON`, `MarshalText` and `UnmarshalText`, so a `NullDecimal` field is encoded differently in JSON:

--8<-- "examples/migration/Example_nullDecimalJSON.md"

### Binary encoding

Both libraries implement `MarshalBinary` and `UnmarshalBinary`, but the byte layouts differ, and udecimal returns
`invalid binary data` for shopspring's output. Text, JSON and `database/sql` values are plain decimal strings in
both libraries and can be exchanged if they have at most 19 decimal places, at most 200 characters and no exponent.

### Package settings

shopspring's package variables and their udecimal equivalents:

- `DivisionPrecision`: [`SetDefaultPrecision`](precision.md#default-precision) sets the number of decimal places for
  every operation and for parsing, not only for `Div`.
- `MarshalJSONWithoutQuotes`: none. `MarshalJSON` always writes a quoted string, and `UnmarshalJSON` accepts both
  forms.
- `TrimTrailingZeros`: none. `String` always removes trailing zeros; use `StringFixed` for a fixed number of decimal
  places.
- `PowPrecisionNegativeExponent`: none. `PowInt32` uses the default precision for negative exponents too.
- `UseScientificNotation`: none. udecimal doesn't read or write scientific notation.
- `MaxDecodeExponent`: none. udecimal doesn't accept exponents, and input is limited to 200 characters.
- `ExpMaxIterations`: none. udecimal has no exponential functions.

## API mapping

### Creating values

| shopspring | udecimal | Note |
|---|---|---|
| `NewFromString` | `Parse` | returns an error for more than 19 decimal places and for [some formats](#input-formats) shopspring accepts |
| `RequireFromString` | `MustParse` | |
| `NewFromInt` | `NewFromInt64(v, 0)` | returns an error |
| `NewFromInt32` | `NewFromInt64(int64(v), 0)` | |
| `NewFromUint64` | `NewFromUint64(v, 0)` | returns an error |
| `New` | `NewFromInt64(value, uint8(-exp))` | for `exp` from 0 down to −19 |
| `NewFromFloat` | `NewFromFloat64` | returns an error |
| `Zero` | `Zero` | |
| `NewNullDecimal` | `NullDecimal{Decimal: d, Valid: true}` | |
| `NewFromFloat32` | no equivalent | `NewFromFloat64(float64(f))` gives the float64 digits: `0.10000000149011612` for `float32(0.1)`, where shopspring gives `0.1` |
| `NewFromFloatWithExponent` | no equivalent | `NewFromFloat64`, then a rounding method |
| `NewFromBigInt` | no equivalent | |
| `NewFromBigRat` | no equivalent | |
| `NewFromBigFloat` | no equivalent | |
| `NewFromFormattedString` | no equivalent | remove the formatting, then `Parse` |

### Arithmetic

| shopspring | udecimal | Note |
|---|---|---|
| `Add` | `Add` | |
| `Sub` | `Sub` | |
| `Mul` | `Mul` | keeps at most 19 decimal places |
| `Div` | `Div` | returns an error; keeps 19 decimal places, truncated |
| `DivRound` | `Div`, then `RoundHAZ(p)` | the same result for `p` below 19 |
| `Mod` | `Mod` | returns an error |
| `QuoRem` | `QuoRem` | no precision argument: the quotient is an integer, as with shopspring's precision 0; returns an error |
| `PowInt32` | `PowInt32` | [differs at 0^0 and in precision](#powint32) |
| `Sqrt` | `Sqrt` | keeps 19 decimal places; shopspring keeps `DivisionPrecision` (16 by default) |
| `Neg` | `Neg` | |
| `Abs` | `Abs` | |
| `Pow` | no equivalent | `PowInt32` for an integer exponent, `PowToIntPart` for the integer part of a `Decimal` exponent |
| `PowWithPrecision` | no equivalent | |
| `PowBigInt` | no equivalent | `PowInt32` |
| `SqrtRound` | no equivalent | `Sqrt`, then a rounding method |
| `Root` | no equivalent | |
| `RootRound` | no equivalent | |
| `Ln` | no equivalent | |
| `ExpTaylor` | no equivalent | |
| `ExpHullAbrham` | no equivalent | |
| `Sin` | no equivalent | |
| `Cos` | no equivalent | |
| `Tan` | no equivalent | |
| `Atan` | no equivalent | |
| `Shift` | no equivalent | `Mul` by a power of ten |

### Comparison

| shopspring | udecimal | Note |
|---|---|---|
| `Cmp` | `Cmp` | |
| `Compare` | `Cmp` | |
| `Equal` | `Equal` | |
| `Equals` | `Equal` | deprecated in shopspring |
| `GreaterThan` | `GreaterThan` | |
| `GreaterThanOrEqual` | `GreaterThanOrEqual` | |
| `LessThan` | `LessThan` | |
| `LessThanOrEqual` | `LessThanOrEqual` | |
| `IsZero` | `IsZero` | |
| `IsNegative` | `IsNeg` | |
| `IsPositive` | `IsPos` | |
| `Sign` | `Sign` | |
| `IsInteger` | no equivalent | `d.Equal(d.Trunc(0))` |

### Rounding

| shopspring | udecimal | Note |
|---|---|---|
| `Round` | `RoundHAZ` | half away from zero |
| `RoundBank` | `RoundBank` | half to even |
| `RoundUp` | `RoundAwayFromZero` | |
| `RoundDown` | `Trunc` | |
| `Truncate` | `Trunc` | |
| `Floor` | `Floor` | |
| `Ceil` | `Ceil` | |
| `RoundCeil` | no equivalent | `Ceil` rounds to an integer only |
| `RoundFloor` | no equivalent | `Floor` rounds to an integer only |
| `RoundCash` | no equivalent | |

### Conversion

| shopspring | udecimal | Note |
|---|---|---|
| `String` | `String` | both drop trailing zeros |
| `StringFixed` | `StringFixed` | pads with zeros and doesn't round; [see above](#stringfixed) |
| `StringFixedBank` | `RoundBank(p).StringFixed(p)` | |
| `StringFixedCash` | no equivalent | |
| `StringScaled` | no equivalent | |
| `ScientificNotationString` | no equivalent | |
| `IntPart` | `Int64` | returns an error when the integer part doesn't fit in an `int64` |
| `InexactFloat64` | `InexactFloat64` | |
| `Float64` | no equivalent | `InexactFloat64`, without the exactness flag |
| `Exponent` | `-Prec()` | `Prec` is the number of decimal places |
| `BigInt` | no equivalent | |
| `BigFloat` | no equivalent | |
| `Rat` | no equivalent | |
| `Coefficient` | no equivalent | `ToHiLo` returns the 128-bit coefficient |
| `CoefficientInt64` | no equivalent | |
| `NumDigits` | no equivalent | |
| `Copy` | no equivalent | values are immutable, so assign instead |
| `Normalize` | no equivalent | |
| `Format` | no equivalent | `fmt` prints a `Decimal` with `String` |

### Aggregates

| shopspring | udecimal | Note |
|---|---|---|
| `Avg` | `Avg` | keeps 19 decimal places, like `Div` |
| `Max` | `Max` | |
| `Min` | `Min` | |
| `Sum` | no equivalent | `Add` in a loop |
| `RescalePair` | no equivalent | |

### Encoding

| shopspring | udecimal | Note |
|---|---|---|
| `MarshalJSON` | `MarshalJSON` | always a quoted string |
| `UnmarshalJSON` | `UnmarshalJSON` | |
| `MarshalText` | `MarshalText` | |
| `UnmarshalText` | `UnmarshalText` | |
| `MarshalBinary` | `MarshalBinary` | a different byte layout |
| `UnmarshalBinary` | `UnmarshalBinary` | can't read shopspring's bytes |
| `Scan` | `Scan` | doesn't accept a quoted string |
| `Value` | `Value` | both store the string form |
| `NullDecimal.Scan` | `Scan` | |
| `NullDecimal.Value` | `Value` | |
| `NullDecimal.MarshalJSON` | no equivalent | [encoded as a struct in JSON](#nulldecimal-encoding) |
| `NullDecimal.UnmarshalJSON` | no equivalent | |
| `NullDecimal.MarshalText` | no equivalent | |
| `NullDecimal.UnmarshalText` | no equivalent | |
| `NullDecimal.MarshalBSONValue` | no equivalent | |
| `NullDecimal.UnmarshalBSONValue` | no equivalent | |
| `NullDecimal.EncodeSpanner` | no equivalent | |
| `NullDecimal.DecodeSpanner` | no equivalent | |
| `GobEncode` | no equivalent | |
| `GobDecode` | no equivalent | |
| `MarshalBSONValue` | no equivalent | |
| `UnmarshalBSONValue` | no equivalent | |
| `EncodeSpanner` | no equivalent | |
| `DecodeSpanner` | no equivalent | |
