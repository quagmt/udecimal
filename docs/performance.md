# Performance

udecimal compared with shopspring, ericlagergren, apd and govalues on an Apple M4. Unless noted, inputs have at most
19 digits.

- **`Add`, `Sub`, `Mul` and `Div`** don't allocate and are 4.5–7.7× faster than shopspring, 1.9–5.1× faster than
  ericlagergren and 3.1–8.1× faster than apd.
- **govalues** is about as fast on `Add`, faster on `Sub`, `String` and `MarshalText`, and 20.6× slower on `Div`.
- **shopspring** is faster on `Pow` with an exponent of 100.
- **Over 19 digits**, govalues can't represent the values. udecimal is the fastest up to 38 digits, and in 1 of the 6
  cases beyond that.

Each library's time as a multiple of udecimal's. Every value is the geometric mean of the median times over a
benchmark's cases.

--8<-- "performance/summary.md"

## Arithmetic

By geometric mean, udecimal is the fastest on `Add`, `Mul` and `Div`, and govalues is the fastest on `Sub` (5.5 ns,
against 6.3 ns for udecimal). govalues is also slightly faster in three individual `Add` cases and three `Mul` cases.

--8<-- "performance/arithmetic.html"

govalues' `Mul` takes about 4 ns when the product fits in 19 digits. In the last two cases the product is longer,
so govalues rounds it to 19 digits, which takes 91–97 ns; the other libraries keep the full product. govalues' `Div`
takes 193–212 ns and allocates twice in every case.

??? note "All cases"
    --8<-- "performance/cases-arithmetic.md"

## Parsing and formatting

udecimal parses 2.0–4.8× faster than the other libraries. Formatting is closer: govalues' `String` is faster (0.62×
udecimal's time), and apd's is similar overall (1.1×) and faster in 3 of the 4 cases.

--8<-- "performance/parse-format.html"

??? note "All cases, including MarshalText and UnmarshalText"
    --8<-- "performance/cases-parse-format.md"

## JSON

Most of the time is spent in `encoding/json`: udecimal's `String` takes 17.6 ns, and `json.Marshal` of a udecimal
value takes 76.6 ns. The differences between libraries shrink to 1.1–1.8×, and udecimal is the fastest in both
directions.

--8<-- "performance/json.html"

??? note "All cases"
    --8<-- "performance/cases-json.md"

## Pow

The libraries keep different numbers of digits in a `Pow` result. 1.01^100 has 200 decimal places: shopspring's
`PowInt32` keeps all of them, udecimal keeps 19 decimal places, govalues 19 digits in total, and ericlagergren and
apd 38 significant digits. udecimal is the fastest for 1.01^10. shopspring is faster for 1.01^100, where udecimal
takes 433 ns and allocates 13 times.

--8<-- "performance/pow.md"

## Beyond 19 digits

These inputs have more than 19 significant digits, which is more than govalues can hold. govalues either rounds them
during parsing without returning an error, or rejects them, so it has no results in this section.

### 20 to 38 digits

These values still fit in udecimal's 128-bit representation. Parsing and arithmetic don't allocate, and udecimal is
the fastest in every case.

--8<-- "performance/wide-128.html"

??? note "All cases"
    --8<-- "performance/cases-wide-128.md"

??? note "govalues errors"
    The reasons recorded by the benchmark. `parsed X as Y` means govalues returned `Y` for the input `X`.

    --8<-- "performance/govalues-reasons.md"

### More than 38 digits

These cases have an operand of 40 to 48 digits, which doesn't fit in 128 bits, so udecimal uses `math/big`. It
allocates 4–7 times per operation, is the fastest in 1 of the 6 cases, and is up to 15% slower than the fastest
library in the others.

??? note "All cases"
    --8<-- "performance/wide-big.md"

## Allocations

Allocations per operation. udecimal's `Add`, `Sub`, `Mul`, `Div` and parsing don't allocate. `String`,
`MarshalText` and `MarshalJSON` allocate the string or byte slice they return. A range means the number differs
between cases.

--8<-- "performance/allocs.md"

## Methodology

The benchmarks ran on an Apple M4 (darwin/arm64) with Go 1.27.4 on 2026-11-20, from 09:26 to 11:16 UTC, using
`go test -run NONE -bench=./<library> -benchmem -count=10 -timeout=2h -v` once per library
([Makefile](https://github.com/quagmt/udecimal/blob/master/benchmarks/Makefile)). Each value is the median of 10
runs. The charts and the summary table use the geometric mean of these medians over each benchmark's cases.

Results depend on the hardware. To run the benchmarks yourself, see
[benchmarks/README.md](https://github.com/quagmt/udecimal/blob/master/benchmarks/README.md).

??? note "Library versions and configuration"
    | Library | Version | Configuration |
    |---|---|---|
    | udecimal | v1.12.0 | defaults: up to 19 decimal places, extra digits truncated |
    | [shopspring/decimal](https://github.com/shopspring/decimal) | v1.6.0 | `DivisionPrecision = 19`, so `Div` keeps 19 decimal places like udecimal; shopspring rounds the last digit half away from zero. `Pow` is `PowInt32`. |
    | [ericlagergren/decimal](https://github.com/ericlagergren/decimal) | v0.0.0-20240411145413-00de7ca16731 | results: 38 significant digits, `ToZero` (truncate). Inputs are parsed at precision 64, so they are read exactly. Results are written to a separate `z`, never to an operand. |
    | [cockroachdb/apd](https://github.com/cockroachdb/apd) | v3.2.3 | `BaseContext.WithPrecision(38)`, `RoundDown` (truncate) |
    | [govalues/decimal](https://github.com/govalues/decimal) | v0.1.36 | not configurable: results are rounded to 19 significant digits, half to even |

All benchmarks follow these rules ([source](https://github.com/quagmt/udecimal/tree/master/benchmarks)):

- Inputs are parsed exactly. Before timing, each library's parsed value is compared with the input. A mismatch fails
  the benchmark, and the report generator rejects the run. The exception is [Beyond 19 digits](#beyond-19-digits),
  where govalues is reported as "can't represent".
- Only operations that all five libraries provide are compared. JSON is encoded with `json.Marshal` and
  `json.Unmarshal`, because apd has no JSON methods and ericlagergren has no `MarshalJSON`. Binary encoding isn't
  compared, because apd and ericlagergren don't implement it.
- Results aren't compared between libraries, only inputs. Results can differ in the last digits because the
  libraries keep different numbers of digits.

Raw results:

- Per library: [udecimal](https://github.com/quagmt/udecimal/blob/master/benchmarks/bench-udec.txt),
  [shopspring](https://github.com/quagmt/udecimal/blob/master/benchmarks/bench-ss.txt),
  [ericlagergren](https://github.com/quagmt/udecimal/blob/master/benchmarks/bench-ed.txt),
  [apd](https://github.com/quagmt/udecimal/blob/master/benchmarks/bench-apd.txt),
  [govalues](https://github.com/quagmt/udecimal/blob/master/benchmarks/bench-gov.txt)
- `benchstat` against udecimal:
  [shopspring](https://github.com/quagmt/udecimal/blob/master/benchmarks/bench-vs-shopspring.txt),
  [ericlagergren](https://github.com/quagmt/udecimal/blob/master/benchmarks/bench-vs-ericlagergren.txt),
  [apd](https://github.com/quagmt/udecimal/blob/master/benchmarks/bench-vs-apd.txt),
  [govalues](https://github.com/quagmt/udecimal/blob/master/benchmarks/bench-vs-govalues.txt)
