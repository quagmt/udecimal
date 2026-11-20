# How it works

A `Decimal` is a fixed-point number:

```go
// value = (neg ? -1 : 1) * coef * 10^(-prec)
type Decimal struct {
	coef bint
	neg  bool
	prec uint8 // 0 <= prec <= 19
}

// 123.456  -> neg = false, coef = 123456, prec = 3
// -123.456 -> neg = true,  coef = 123456, prec = 3
```

The coefficient is a `bint`, which holds one of two representations:

```go
type bint struct {
	// set when the value doesn't fit in u128, or was computed with math/big
	bigInt *big.Int

	// 0 <= u128 <= 2^128 - 1
	u128 u128
}
```

## Range

`u128` is a pair of `uint64` values and holds coefficients up to 2^128 − 1 = 340282366920938463463374607431768211455.
That covers every coefficient of up to 38 digits and some of 39. Because the coefficient includes the decimal places,
the largest value depends on `prec`:

| `prec` | Largest magnitude in 128 bits |
|---|---|
| 0 | 340282366920938463463374607431768211455 |
| 2 | 3402823669209384634633746074317682114.55 |
| 19 | 34028236692093846346.3374607431768211455 |

## Fast path

On the fast path, operations use 64-bit and 128-bit integer arithmetic, with a 256-bit intermediate result for `Mul`
and `Div`, and don't allocate. Each operation uses the fast path under these conditions:

| Operation | Condition |
|---|---|
| `Parse` | the coefficient fits in 128 bits, which is always true for up to 38 digits |
| `Add`, `Sub` | both operands, scaled to the larger `prec`, and the result fit |
| `Mul` | both operands and the result, truncated to at most 19 decimal places, fit |
| `Div` | both operands fit, and the quotient fits at `prec = 19`, the default precision, so it is below about 3.4 × 10^19 |

`Div` always returns 19 decimal places, so its limit doesn't depend on the precision of the operands:

```go
// 34028236692093846346 / 1 -> fast path, no allocation
// 34028236692093846347 / 1 -> the quotient's coefficient exceeds 2^128 - 1: math/big, allocates
```

Lowering the default precision with `SetDefaultPrecision` raises this limit. At a default precision of 2, the
quotient can reach about 3.4 × 10^36.

A value computed with `math/big` keeps using `math/big`, even if it would fit in 128 bits. Operations on it allocate,
`ToHiLo` returns `ok = false`, and `Int64` returns `ErrIntPartOverflow`. For example, `10^39 - (10^39 - 1)` returns
`1` stored as `math/big`.

## math/big fallback

If an operand or the result doesn't fit, the whole operation is computed again with `math/big`. This is correct but
slower, and it allocates. In the [benchmarks](performance.md#more-than-38-digits), the fallback allocates 4–7 times
per operation, and udecimal is the fastest library in only 1 of those 6 cases.

## Truncation

When a result has more decimal places than the precision allows, udecimal truncates the extra digits, which rounds
toward zero. It doesn't choose a rounding mode for you, because the modes give different results for the same input
and financial code often has to follow a specific one. To round, call one of the
[rounding methods](precision.md#rounding).
