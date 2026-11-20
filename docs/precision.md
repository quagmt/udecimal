# Precision and rounding

A `Decimal` consists of a sign, a coefficient and a precision:

```
value = (neg ? -1 : 1) × coef × 10^(-prec)
```

`prec` is the number of decimal places, at most 19. For example, `123.456` has `coef = 123456` and `prec = 3`.

## Result precision

| Operation | Decimal places in the result |
|---|---|
| `Add`, `Sub` | the larger `prec` of the two operands |
| `Mul` | the sum of both `prec` values, truncated to 19 |
| `Div` | 19, truncated |

Truncation drops the extra digits, which rounds toward zero. Results are not rounded unless you call a
[rounding method](#rounding).

--8<-- "examples/udecimal/ExampleDecimal_Div.md"

Operations whose operands and result fit in 128 bits don't allocate; other operations use `math/big`. See
[How it works](how-it-works.md) for the exact limits.

## Default precision

`SetDefaultPrecision(prec)` sets the maximum number of decimal places for the whole program, for example to match a
currency. Call it once at startup. It panics if `prec` is 0 or greater than 19.

--8<-- "examples/udecimal/ExampleSetDefaultPrecision.md"

The default precision also applies to input. `Parse`, `NewFromInt64`, `NewFromUint64` and `NewFromHiLo` return
`ErrPrecOutOfRange` for a value with more decimal places than the default precision:

--8<-- "examples/udecimal/ExampleNewFromInt64.md"

To truncate such input instead, see [Parse modes](parsing-and-encoding.md#parse-modes).

## Rounding

udecimal doesn't round implicitly. To round a value, call one of these methods. Each takes the number of decimal
places to keep.

| Method | Rule | `1.12345` → 4 places | `1.12335` → 4 places | `1.5` → 0 | `-1.5` → 0 |
|---|---|---|---|---|---|
| [`RoundBank`](https://en.wikipedia.org/wiki/Rounding#Rounding_half_to_even) | half to even | 1.1234 | 1.1234 | 2 | -2 |
| [`RoundHAZ`](https://en.wikipedia.org/wiki/Rounding#Rounding_half_away_from_zero) | half away from zero | 1.1235 | 1.1234 | 2 | -2 |
| [`RoundHTZ`](https://en.wikipedia.org/wiki/Rounding#Rounding_half_toward_zero) | half toward zero | 1.1234 | 1.1233 | 1 | -1 |

--8<-- "examples/udecimal/ExampleDecimal_RoundBank.md"

### Other rounding methods

| Method | Rule | Examples |
|---|---|---|
| [`RoundAwayFromZero(prec)`](https://en.wikipedia.org/wiki/Rounding#Rounding_away_from_zero) | any remainder rounds away from zero | `1.12` → `1.2`, `-1.12` → `-1.2` (1 place) |
| `Trunc(prec)` | drops the extra digits | `1.23` → `1.2` (1 place) |
| `Floor()` | toward negative infinity, to an integer | `1.23` → `1`, `-1.23` → `-2` |
| `Ceil()` | toward positive infinity, to an integer | `1.23` → `2` |

If the value has fewer decimal places than requested, it is returned unchanged.

--8<-- "examples/udecimal/ExampleDecimal_Trunc.md"

### Rounding currency amounts

Round once, at the end, to the number of decimal places of the currency. `RoundHAZ(2)` rounds ties away from zero.
`RoundBank(2)`, also called banker's rounding, rounds ties to the nearest even digit, which avoids a systematic bias
when many values are rounded.

--8<-- "examples/udecimal/Example_money.md"
