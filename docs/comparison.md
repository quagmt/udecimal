# Comparison

This page compares udecimal with four other Go decimal libraries, using the versions benchmarked on the
[Performance](performance.md) page.

## Choosing a library

- **udecimal** suits code where 19 decimal places are enough and arithmetic should be fast and allocation-free.
- **shopspring/decimal** suits code that needs unbounded precision. `Add`, `Sub` and `Mul` are exact. See
  [Migrating from shopspring](from-shopspring.md) to switch.
- **cockroachdb/apd** implements the General Decimal Arithmetic model, with contexts, conditions and traps.
- **govalues/decimal** suits code where every value has at most 19 digits in total, including the integer part.
- **ericlagergren/decimal** provides arbitrary precision with contexts, using pointers.

## Summary

| Library | Representation | Digits kept | Extra digits | Allocations¹ |
|---|---|---|---|---|
| **udecimal** | 128-bit, `math/big` fallback | 19 decimal places | truncated | 0 / 0 / 0 / 0 |
| [shopspring](https://github.com/shopspring/decimal) v1.6.0 | `big.Int` and exponent | all | `Div` [rounds to 16](https://pkg.go.dev/github.com/shopspring/decimal#pkg-variables) | 2 / 2 / 2 / 7 |
| [apd](https://github.com/cockroachdb/apd) v3.2.3 | `big.Int` and context | [set by context](https://pkg.go.dev/github.com/cockroachdb/apd/v3#Context) | rounded by context | 0 / 0 / 0 / 2 |
| [govalues](https://github.com/govalues/decimal) v0.1.36 | `uint64` | [19 in total](https://pkg.go.dev/github.com/govalues/decimal#hdr-Constraints_Overview) | [rounded, half to even](https://pkg.go.dev/github.com/govalues/decimal#hdr-Rounding_Methods) | 0 / 0 / 0 / 2 |
| [ericlagergren](https://github.com/ericlagergren/decimal) | `big.Int` and context | set by context, [16 by default](https://pkg.go.dev/github.com/ericlagergren/decimal#pkg-constants) | rounded by context | 0 / 0 / 0 / 2 |

¹ Allocations per `Add` / `Sub` / `Mul` / `Div` for inputs of up to 19 digits. See
[Arithmetic](performance.md#arithmetic).

## Libraries

### udecimal

udecimal is a fixed-point type with a 128-bit coefficient, which holds values of up to 38 digits in total.
Arithmetic on these values doesn't allocate within the limits described in [How it works](how-it-works.md); larger
values use `math/big`. Results have at most 19 decimal places. Extra digits are truncated, and other rounding modes
are explicit method calls. Values are immutable, and operations that can fail, such as `Div`, return an error.

### shopspring/decimal

shopspring/decimal is a general-purpose library with no upper limit on size. `Div` rounds the quotient half away
from zero to `DivisionPrecision` decimal places, 16 by default. In the benchmarks it allocates twice per `Add`, `Sub`
and `Mul`, and 7 times per `Div`. Values are immutable.

### cockroachdb/apd

apd implements much of the General Decimal Arithmetic specification. A `Context` sets the precision, the rounding
mode and the traps, and results are rounded to that precision. Operations take pointers, for example
`ctx.Add(z, x, y)`, and return `(Condition, error)`.

### govalues/decimal

govalues/decimal stores the coefficient in a `uint64`, so a value has at most 19 digits in total, counting both sides
of the decimal point. Longer values are rounded half to even, and a result that would lose an integer digit returns
an overflow error. `Add`, `Sub` and `Mul` don't allocate. `Quo` falls back to `big.Int` when the `uint64` computation
overflows; in the benchmarks, `Div` allocates twice and is the slowest of the five libraries. Operations return
`(Decimal, error)`.

### ericlagergren/decimal

ericlagergren/decimal provides arbitrary precision with a context, similar to apd. Operations work on pointers and
overwrite the receiver, for example `z.Add(x, y)`. If no precision is set, a context keeps 16 significant digits.

## govalues and inputs over 19 digits

The govalues documentation states that [`Parse`](https://pkg.go.dev/github.com/govalues/decimal#Parse) returns a
"possibly rounded" decimal. In v0.1.36, an input with more than 19 digits in total is rounded to 19 digits without
an error. If the integer part has 19 digits, the fractional part is lost entirely:

```go
import govalues "github.com/govalues/decimal"

d, err := govalues.Parse("1234567890123456789.1234567890123456879")
// d = 1234567890123456789, err = nil

d, err = udecimal.Parse("1234567890123456789.1234567890123456879")
// d = 1234567890123456789.1234567890123456879, err = nil
```

If the integer part is shorter, the fraction is rounded instead: `12345.1234567890123456789` becomes
`12345.12345678901235`. An integer part of more than 19 digits returns an error. For this reason, the
[Beyond 19 digits](performance.md#beyond-19-digits) benchmarks list govalues as "can't represent" instead of timing
it on a different number. This behavior is covered by
[`TestGovaluesLosesFraction`](https://github.com/quagmt/udecimal/blob/master/benchmarks/exact_test.go).
