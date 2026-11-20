# Parsing and encoding

## JSON

A `Decimal` can be used directly as a struct field. `MarshalJSON` always writes a quoted string, so JSON parsers in
other languages don't read the value as a floating-point number. `UnmarshalJSON` accepts quoted and unquoted
numbers, and the JSON literal `null` leaves the `Decimal` unchanged.

--8<-- "examples/udecimal/Example_jsonStruct.md"

Decoding uses the same rules as [`Parse`](#parsing), including the parse mode.

## Databases

`Decimal` implements `sql.Scanner` and `driver.Valuer`. `Value` stores the decimal as a string. `Scan` accepts
`string`, `[]byte`, `int`, `int32`, `int64`, `uint64` and `float64`, and returns an error for `nil`. For nullable
columns, use `NullDecimal`; its `Value` returns `nil` when `Valid` is false.

--8<-- "examples/udecimal/ExampleNullDecimal_Scan.md"

--8<-- "examples/udecimal/ExampleNullDecimal_Value.md"

## Parsing

--8<-- "examples/udecimal/ExampleParse.md"

`Parse` accepts:

- an optional sign, `+` or `-`
- one or more digits
- optionally, a `.` followed by one or more digits, up to the [default precision](precision.md#default-precision)
  (19 unless changed)
- at most 200 characters in total

| Input | Result |
|---|---|
| `+1.5` | `1.5` |
| `007.50` | `7.5` |
| `.5`, `-.5`, `1.`, `+` | `ErrInvalidFormat` |
| `1e5` | `ErrInvalidFormat`; scientific notation isn't supported |
| `""` | `ErrEmptyString` |
| 201 characters or more | `ErrMaxStrLen` |

`MustParse` works like `Parse` but panics on error.

### Parse modes

By default, input with more decimal places than the default precision returns an error (`ParseModeError`). If your
data has more decimal places and losing them is acceptable, use `ParseModeTrunc`, which truncates the extra digits:

--8<-- "examples/udecimal/ExampleSetDefaultParseMode.md"

The parse mode (`ParseMode`) applies to the whole package: `Parse`, `MustParse` and the text decoders
(`UnmarshalJSON`, `UnmarshalText` and `Scan`). Set it once at startup. `SetDefaultParseMode` panics for any value
other than `ParseModeError` or `ParseModeTrunc`.

## Text, query strings and binary

`MarshalText` and `UnmarshalText` use the same format as `String` and `Parse`, without quotes. `EncodeValues` lets
[go-querystring](https://github.com/google/go-querystring) encode a `Decimal` field as a string.

`MarshalBinary` writes a compact binary form:

--8<-- "examples/udecimal/ExampleDecimal_MarshalBinary.md"

The layout is `[flags] [prec] [total length] [coefficient bytes]`: one flag byte for the sign and for whether the
coefficient is stored as `math/big`, the precision, the total length in bytes, and the coefficient in big-endian
order. A 128-bit coefficient takes 8 bytes, or 16 if it doesn't fit in 64 bits. `UnmarshalBinary` returns
`ErrInvalidBinaryData` for data that doesn't follow this layout.

## Interfaces

| Interface | Methods | Format |
|---|---|---|
| `json.Marshaler`, `json.Unmarshaler` | `MarshalJSON`, `UnmarshalJSON` | quoted string, for example `"1.23"` |
| `encoding.TextMarshaler`, `TextUnmarshaler`, `TextAppender` | `MarshalText`, `UnmarshalText`, `AppendText` | plain string, for example `1.23` |
| `encoding.BinaryMarshaler`, `BinaryUnmarshaler`, `BinaryAppender` | `MarshalBinary`, `UnmarshalBinary`, `AppendBinary` | compact binary (see above) |
| `sql.Scanner`, `driver.Valuer` | `Scan`, `Value` | string |
| go-querystring `query.Encoder` | `EncodeValues` | string |
