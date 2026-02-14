package udecimal

import (
	"fmt"
	"testing"
)

// Benchmark the fast-path operations that use reciprocal multiplication.
// These cover rounding, truncation, floor, ceil, string formatting,
// and trailing-zeros detection.

var (
	benchDecimals = []struct {
		name string
		d    Decimal
	}{
		{"small", MustParse("123.456")},
		{"medium", MustParse("123456789.123456789")},
		{"maxprec", MustParse("0.1234567890123456789")},
		{"integer", MustParse("123456789")},
	}
	sinkDecimal Decimal
	sinkString  string
)

func BenchmarkRoundBank(b *testing.B) {
	for _, tc := range benchDecimals {
		if tc.d.IsZero() || tc.d.prec == 0 {
			continue
		}
		target := tc.d.prec / 2
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				sinkDecimal = tc.d.RoundBank(uint8(target))
			}
		})
	}
}

func BenchmarkRoundHAZ(b *testing.B) {
	for _, tc := range benchDecimals {
		if tc.d.IsZero() || tc.d.prec == 0 {
			continue
		}
		target := tc.d.prec / 2
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				sinkDecimal = tc.d.RoundHAZ(uint8(target))
			}
		})
	}
}

func BenchmarkTrunc(b *testing.B) {
	for _, tc := range benchDecimals {
		if tc.d.IsZero() || tc.d.prec == 0 {
			continue
		}
		target := tc.d.prec / 2
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				sinkDecimal = tc.d.Trunc(uint8(target))
			}
		})
	}
}

func BenchmarkFloor(b *testing.B) {
	for _, tc := range benchDecimals {
		if tc.d.IsZero() || tc.d.prec == 0 {
			continue
		}
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				sinkDecimal = tc.d.Floor()
			}
		})
	}
}

func BenchmarkCeil(b *testing.B) {
	for _, tc := range benchDecimals {
		if tc.d.IsZero() || tc.d.prec == 0 {
			continue
		}
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				sinkDecimal = tc.d.Ceil()
			}
		})
	}
}

func BenchmarkDecimalString(b *testing.B) {
	for _, tc := range benchDecimals {
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				sinkString = tc.d.String()
			}
		})
	}
}

func BenchmarkTrimTrailingZeros(b *testing.B) {
	// Decimals with trailing zeros of varying lengths
	cases := []struct {
		name string
		d    Decimal
	}{
		{"no_trailing", MustParse("123.456")},
		{"some_trailing", MustParse("123.45600000")},
		{"all_trailing", MustParse("123.000000000")},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				sinkDecimal = tc.d.trimTrailingZeros()
			}
		})
	}
}

// TestReciprocal verifies the reciprocal table against hardware division
// for all pow10 indices and a selection of edge-case dividends.
func TestReciprocal(t *testing.T) {
	edgeCases := []uint64{
		0, 1, 2, 9, 10, 99, 100, 999, 1000,
		1<<32 - 1, 1 << 32, 1<<32 + 1,
		1<<63 - 1, 1 << 63, 1<<64 - 1,
		1e18, 1e18 + 1, 1e18 - 1,
	}

	for k := uint8(0); k <= 19; k++ {
		d := pow10[k].lo
		t.Run(fmt.Sprintf("pow10_%d", k), func(t *testing.T) {
			for _, n := range edgeCases {
				u := u128{lo: n}
				gotQ, gotR := u.quoRem64Pow10(k)
				wantQ, wantR := u.QuoRem64(d)
				if gotQ != wantQ || gotR != wantR {
					t.Errorf("quoRem64Pow10(%d) for n=%d: got (%v,%d), want (%v,%d)",
						k, n, gotQ, gotR, wantQ, wantR)
				}
			}
		})
	}
}
