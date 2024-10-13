package benchmarks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/quagmt/udecimal"

	gv "github.com/govalues/decimal"
	ss "github.com/shopspring/decimal"
)

func BenchmarkParse(b *testing.B) {
	testcases := []string{
		"1234567890123456789.1234567890123456879",
		"123",
		"123456.123456",
		"1234567890",
		"0.1234567890123456879",
	}

	for _, tc := range testcases {
		// shopspring benchmark
		b.Run(fmt.Sprintf("ss/%s", tc), func(b *testing.B) {
			b.ResetTimer()
			for range b.N {
				_, _ = ss.NewFromString(tc)
			}
		})

		b.Run(fmt.Sprintf("udec/%s", tc), func(b *testing.B) {
			b.ResetTimer()
			for range b.N {
				_, _ = udecimal.Parse(tc)
			}
		})
	}
}

func BenchmarkParseFallBack(b *testing.B) {
	testcases := []string{
		"123456789123456789123456.1234567890123456",
		"111222333444555666777888999.1234567890123456789",
	}

	for _, tc := range testcases {
		// shopspring benchmark
		b.Run(fmt.Sprintf("ss/%s", tc), func(b *testing.B) {
			b.ResetTimer()
			for range b.N {
				_, _ = ss.NewFromString(tc)
			}
		})

		b.Run(fmt.Sprintf("udec/%s", tc), func(b *testing.B) {
			b.ResetTimer()
			for range b.N {
				_, _ = udecimal.Parse(tc)
			}
		})
	}
}

func BenchmarkString(b *testing.B) {
	testcases := []string{
		"1234567890123456789.1234567890123456879",
		"123",
		"123456.123456",
		"1234567890",
		"0.1234567890123456879",
	}

	for _, tc := range testcases {
		// shopspring benchmark
		b.Run(fmt.Sprintf("ss/%s", tc), func(b *testing.B) {
			bb := ss.RequireFromString(tc)

			b.ResetTimer()
			for range b.N {
				_ = bb.String()
			}
		})

		b.Run(fmt.Sprintf("udec/%s", tc), func(b *testing.B) {
			bb := udecimal.MustParse(tc)

			b.ResetTimer()
			for range b.N {
				_ = bb.String()
			}
		})
	}
}

func BenchmarkStringFallBack(b *testing.B) {
	testcases := []string{
		"123456789123456789123456.1234567890123456",
		"111222333444555666777888999.1234567890123456789",
	}

	for _, tc := range testcases {
		// shopspring benchmark
		b.Run(fmt.Sprintf("ss/%s", tc), func(b *testing.B) {
			bb := ss.RequireFromString(tc)

			b.ResetTimer()
			for range b.N {
				_ = bb.String()
			}
		})

		b.Run(fmt.Sprintf("udec/%s", tc), func(b *testing.B) {
			bb := udecimal.MustParse(tc)

			b.ResetTimer()
			for range b.N {
				_ = bb.String()
			}
		})
	}
}

func BenchmarkAdd(b *testing.B) {
	testcases := []struct {
		a, b string
	}{
		{"1234567890123456789.1234567890123456879", "1111.1789"},
		{"1234567890123456789.1234567890123456879", "1234567890123456789.1234567890123456789"},
		{"123.456", "0.123"},
		{"3", "7"},
		{"123456.123456", "999999"},
		{"123456.123456", "456781244.1324897546"},
		{"548751.15465466546", "1542.456487"},
	}

	for _, tc := range testcases {
		// shopspring benchmark
		b.Run(fmt.Sprintf("ss/%s.Add(%s)", tc.a, tc.b), func(b *testing.B) {
			a := ss.RequireFromString(tc.a)
			bb := ss.RequireFromString(tc.b)

			b.ResetTimer()
			for range b.N {
				_ = a.Add(bb)
			}
		})

		b.Run(fmt.Sprintf("udec/%s.Add(%s)", tc.a, tc.b), func(b *testing.B) {
			a, err := udecimal.Parse(tc.a)
			require.NoError(b, err)

			bb, err := udecimal.Parse(tc.b)
			require.NoError(b, err)

			b.ResetTimer()
			for range b.N {
				_ = a.Add(bb)
			}
		})
	}
}

func BenchmarkSub(b *testing.B) {
	testcases := []struct {
		a, b string
	}{
		{"3", "7"},
		{"1234567890123456789.1234567890123456879", "1111.1789"},
		{"1234567890123456789.1234567890123456879", "1234567890123456789.1234567890123456789"},
		{"123.456", "0.123"},
		{"123456.123456", "456781244.1324897546"},
		{"548751.15465466546", "1542.456487"},
	}

	for _, tc := range testcases {
		// shopspring benchmark
		b.Run(fmt.Sprintf("ss/%s.Sub(%s)", tc.a, tc.b), func(b *testing.B) {
			a := ss.RequireFromString(tc.a)
			bb := ss.RequireFromString(tc.b)

			b.ResetTimer()
			for range b.N {
				_ = a.Sub(bb)
			}
		})

		b.Run(fmt.Sprintf("udec/%s.Sub(%s)", tc.a, tc.b), func(b *testing.B) {
			a, err := udecimal.Parse(tc.a)
			require.NoError(b, err)

			bb, err := udecimal.Parse(tc.b)
			require.NoError(b, err)

			b.ResetTimer()
			for range b.N {
				_ = a.Sub(bb)
			}
		})
	}
}

func BenchmarkMul(b *testing.B) {
	testcases := []struct {
		a, b string
	}{
		{"1234.1234567890123456879", "1111.1789"},
		{"1234.1234567890123456879", "1111.1234567890123456789"},
		{"123.456", "0.123"},
		{"3", "7"},
		{"123456.123456", "999999"},
		{"123456.123456", "456781244.1324897546"},
		{"548751.15465466546", "1542.456487"},
	}

	for _, tc := range testcases {
		// shopspring benchmark
		b.Run(fmt.Sprintf("ss/%s.Mul(%s)", tc.a, tc.b), func(b *testing.B) {
			a := ss.RequireFromString(tc.a)
			bb := ss.RequireFromString(tc.b)

			b.ResetTimer()
			for range b.N {
				_ = a.Mul(bb)
			}
		})

		b.Run(fmt.Sprintf("udec/%s.Mul(%s)", tc.a, tc.b), func(b *testing.B) {
			a, err := udecimal.Parse(tc.a)
			require.NoError(b, err)

			bb, err := udecimal.Parse(tc.b)
			require.NoError(b, err)

			b.ResetTimer()
			for range b.N {
				_ = a.Mul(bb)
			}
		})
	}
}

func BenchmarkDiv(b *testing.B) {
	testcases := []struct {
		a, b string
	}{
		{"1234567890123456789.1234567890123456879", "1111.1789"},
		{"12345.1234567890123456879", "1111.1234567890123456789"},
		{"1234567890123456789.1234567890123456879", "9876543210987654321.1234567890123456789"},
		{"123.456", "0.123"},
		{"3", "7"},
		{"123456.123456", "999999"},
		{"123456.123456", "456781244.1324897546"},
		{"548751.15465466546", "1542.456487"},
	}

	for _, tc := range testcases {
		// shopspring benchmark
		b.Run(fmt.Sprintf("ss/%s.Div(%s)", tc.a, tc.b), func(b *testing.B) {
			a := ss.RequireFromString(tc.a)
			bb := ss.RequireFromString(tc.b)

			b.ResetTimer()
			for range b.N {
				_ = a.Div(bb)
			}
		})

		// govalues benchmark
		b.Run(fmt.Sprintf("gv/%s.Div(%s)", tc.a, tc.b), func(b *testing.B) {
			a, err := gv.Parse(tc.a)
			if err != nil {
				return
			}

			bb, err := gv.Parse(tc.b)
			if err != nil {
				return
			}

			b.ResetTimer()
			for range b.N {
				_, _ = a.Quo(bb)
			}
		})

		b.Run(fmt.Sprintf("udec/%s.Div(%s)", tc.a, tc.b), func(b *testing.B) {
			a, err := udecimal.Parse(tc.a)
			require.NoError(b, err)

			bb, err := udecimal.Parse(tc.b)
			require.NoError(b, err)

			b.ResetTimer()
			for range b.N {
				_, _ = a.Div(bb)
			}
		})
	}
}

func BenchmarkDivFallback(b *testing.B) {
	testcases := []struct {
		a, b string
	}{
		{"12345679012345679890123456789.1234567890123456789", "999999"},
		{"1234", "12345679012345679890123456789.1234567890123456789"},
	}

	for _, tc := range testcases {
		// shopspring benchmark
		b.Run(fmt.Sprintf("ss/%s.Div(%s)", tc.a, tc.b), func(b *testing.B) {
			a := ss.RequireFromString(tc.a)
			bb := ss.RequireFromString(tc.b)

			b.ResetTimer()
			for range b.N {
				_ = a.Div(bb)
			}
		})

		b.Run(fmt.Sprintf("udec/%s.Div(%s)", tc.a, tc.b), func(b *testing.B) {
			a, err := udecimal.Parse(tc.a)
			require.NoError(b, err)

			bb, err := udecimal.Parse(tc.b)
			require.NoError(b, err)

			b.ResetTimer()
			for range b.N {
				_, _ = a.Div(bb)
			}
		})
	}
}

func BenchmarkPow(b *testing.B) {
	testcases := []struct {
		a string
		b int
	}{
		{"1.01", 10},
		{"1.01", 100},
	}

	for _, tc := range testcases {
		// shopspring benchmark
		b.Run(fmt.Sprintf("ss/%s.Pow(%d)", tc.a, tc.b), func(b *testing.B) {
			a := ss.RequireFromString(tc.a)
			bb := ss.NewFromInt(int64(tc.b))

			b.ResetTimer()
			for range b.N {
				_ = a.Pow(bb)
			}
		})

		b.Run(fmt.Sprintf("udec/%s.Pow(%d)", tc.a, tc.b), func(b *testing.B) {
			a := udecimal.MustParse(tc.a)

			b.ResetTimer()
			for range b.N {
				_ = a.PowInt(tc.b)
			}
		})
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	testcases := []string{
		"1234567890123456789.1234567890123456879",
		"123",
		"123456.123456",
		"1234567890",
		"0.1234567890123456879",
		"12345678901234567891234567890123456789.1234567890123456879",
	}

	for _, tc := range testcases {
		b.Run(fmt.Sprintf("ss/%s", tc), func(b *testing.B) {
			bb := ss.RequireFromString(tc)

			b.ResetTimer()
			for range b.N {
				_, _ = bb.MarshalJSON()
			}
		})

		b.Run(fmt.Sprintf("udec/%s", tc), func(b *testing.B) {
			bb := udecimal.MustParse(tc)

			b.ResetTimer()
			for range b.N {
				_, _ = bb.MarshalJSON()
			}
		})
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	testcases := []string{
		"1234567890123456789.1234567890123456879",
		"123",
		"123456.123456",
		"1234567890",
		"0.1234567890123456879",
		"12345678901234567891234567890123456789.1234567890123456879",
	}

	for _, tc := range testcases {
		b.Run(fmt.Sprintf("ss/%s", tc), func(b *testing.B) {
			data, _ := ss.RequireFromString(tc).MarshalJSON()

			b.ResetTimer()
			for range b.N {
				var d ss.Decimal
				_ = d.UnmarshalJSON(data)
			}
		})

		b.Run(fmt.Sprintf("udec/%s", tc), func(b *testing.B) {
			data, _ := udecimal.MustParse(tc).MarshalJSON()

			b.ResetTimer()
			for range b.N {
				var d udecimal.Decimal
				_ = d.UnmarshalJSON(data)
			}
		})
	}
}

func BenchmarkMarshalBinary(b *testing.B) {
	testcases := []string{
		"1234567890123456789.1234567890123456879",
		"123",
		"123456.123456",
		"1234567890",
		"0.1234567890123456879",
		"12345678901234567891234567890123456789.1234567890123456879",
	}

	for _, tc := range testcases {
		b.Run(fmt.Sprintf("ss/%s", tc), func(b *testing.B) {
			bb := ss.RequireFromString(tc)

			b.ResetTimer()
			for range b.N {
				_, _ = bb.MarshalBinary()
			}
		})

		b.Run(fmt.Sprintf("udec/%s", tc), func(b *testing.B) {
			bb := udecimal.MustParse(tc)

			b.ResetTimer()
			for range b.N {
				_, _ = bb.MarshalBinary()
			}
		})
	}
}

func BenchmarkUnmarshalBinary(b *testing.B) {
	testcases := []string{
		"1234567890123456789.1234567890123456879",
		"123",
		"123456.123456",
		"1234567890",
		"0.1234567890123456879",
		"12345678901234567891234567890123456789.1234567890123456879",
	}

	for _, tc := range testcases {
		b.Run(fmt.Sprintf("ss/%s", tc), func(b *testing.B) {
			data, _ := ss.RequireFromString(tc).MarshalBinary()

			b.ResetTimer()
			for range b.N {
				var d ss.Decimal
				_ = d.UnmarshalBinary(data)
			}
		})

		b.Run(fmt.Sprintf("udec/%s", tc), func(b *testing.B) {
			data, _ := udecimal.MustParse(tc).MarshalBinary()

			b.ResetTimer()
			for range b.N {
				var d udecimal.Decimal
				_ = d.UnmarshalBinary(data)
			}
		})
	}
}
