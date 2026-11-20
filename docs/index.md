---
title: Fast fixed-point decimals for Go
---

# udecimal

Fast fixed-point decimal numbers for Go, designed for financial applications.

[![build](https://github.com/quagmt/udecimal/actions/workflows/test.yaml/badge.svg)](https://github.com/quagmt/udecimal/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/quagmt/udecimal/graph/badge.svg?token=662ET843EZ)](https://codecov.io/gh/quagmt/udecimal)
[![GoDoc](https://pkg.go.dev/badge/github.com/quagmt/udecimal)](https://pkg.go.dev/github.com/quagmt/udecimal)
[![Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#financial)

## Features

- **No allocations.** Values are stored in 128 bits, and arithmetic on them doesn't allocate memory as long as the
  results fit. Calculations in hot paths, such as pricing or order matching, put no load on the garbage collector.
  Larger values fall back to `math/big`; [How it works](how-it-works.md) lists the exact limits.
- **Fast.** Addition, subtraction, multiplication and division are 4.5–7.7× faster than shopspring/decimal for values
  of up to 19 digits on an Apple M4. See [Performance](performance.md).
- **No implicit rounding.** A result keeps up to 19 decimal places and extra digits are dropped. Values are rounded
  only when you ask, with the [rounding mode](precision.md#rounding) you choose.
- **Explicit errors.** Operations that can fail, such as dividing by zero or parsing invalid input, return an error
  instead of panicking.
- **Enough precision for money.** 19 decimal places cover every currency and common crypto assets such as Bitcoin
  (8 decimal places) and Ether (18).
- **Safe to share.** Values are immutable, so they can be shared between goroutines.
- **Ready for storage.** Built-in JSON, text, binary and `database/sql` support.
- **Tested.** Arithmetic is fuzz tested against shopspring/decimal.

## Installation

```sh
go get github.com/quagmt/udecimal
```

Requires Go 1.23 or later.

## Usage

--8<-- "examples/udecimal/Example_overview.md"

## Performance

Time per operation for inputs of up to 19 digits, on an Apple M4. Each bar is the geometric mean of the median
times over the benchmark's cases. Lower is better.

--8<-- "performance/arithmetic.html"

See [Performance](performance.md) for all benchmarks, including the cases where other libraries are faster.

## Next steps

- [Comparison](comparison.md) with shopspring, apd, govalues and ericlagergren
- [Migrating from shopspring](from-shopspring.md)
- [Basics](basics.md)
