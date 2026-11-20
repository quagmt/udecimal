# Benchmarks

Benchmarks of udecimal against shopspring/decimal, ericlagergren/decimal, cockroachdb/apd and govalues/decimal.
Results, configuration and charts: <https://quagmt.github.io/udecimal/performance/>.

This is a separate Go module, so the other libraries never become dependencies of udecimal. It measures the
udecimal source in the parent directory (`replace github.com/quagmt/udecimal => ../`).

## Run the benchmarks

From this directory:

```sh
go install golang.org/x/perf/cmd/benchstat@latest
make bench-all stats
go run ./cmd/report
```

- `make bench-all` runs every benchmark 10 times per library and writes `bench-udec.txt`, `bench-ss.txt`,
  `bench-ed.txt`, `bench-apd.txt` and `bench-gov.txt`. It takes about two hours on an Apple M4.
- `make stats` compares each library with udecimal using benchstat and writes `bench-vs-<library>.txt`.
- `go run ./cmd/report` regenerates `../docs/snippets/`, which the documentation site includes:
  - `performance/`: the summary, allocation and per-case tables, and the charts, from the five `bench-*.txt` files
  - `examples/udecimal/` and `examples/migration/`: the Example functions of `../doc_test.go` and
    `migration_test.go`, rendered as code blocks

To run one library only, use its target, for example `make bench-udec`.

## Tests

`go test ./...` (without `-bench`) runs in CI for every change to Go code or the docs. It fails when:

- a library doesn't parse a benchmark input exactly
- `../docs/snippets/` differs from what `go run ./cmd/report` writes: run it and commit the result
- a row of the "Migrating from shopspring" page no longer matches what both libraries do
- an exported identifier of udecimal is mentioned on no page, or a page's hand-written Go code prints output or
  calls a package function without `udecimal.`
