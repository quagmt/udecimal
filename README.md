# udecimal

[![GoDoc](https://pkg.go.dev/badge/github.com/quagmt/udecimal)](https://pkg.go.dev/github.com/quagmt/udecimal)
[![Go Report Card](https://goreportcard.com/badge/github.com/quagmt/udecimal)](https://goreportcard.com/report/github.com/quagmt/udecimal)

Blazing fast, high precision fixed-point decimal number library. Specifically designed for high-traffic financial applications.

## Features

- High precision (up to 19 decimal places) and no precision loss during arithmetic operations
- Panic-free operations
- Optimized for speed (see [benchmarks](#benchmarks)) and zero memory allocation (in most cases, check out [why](#faq))
- Various rounding methods: HALF AWAY FROM ZERO, HALF TOWARD ZERO, and Banker's rounding
- Intuitive API and easy to use

## Installation

```sh
go get github.com/quagmt/udecimal
```
