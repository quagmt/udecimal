package udecimal

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	testcases := []string{
		"0",
		"1",
		"123",
		"123.456",
		"123.456789012345678901",
		"123456789.123456789",
		"-123456789123456789.123456789123456789",
		"-123456.123456",
		"1234567891234567890.0123456879123456789",
		"9999999999999999999.9999999999999999999",
		"-9999999999999999999.9999999999999999999",
		"123456.0000000000000000001",
		"-123456.0000000000000000001",
	}

	for _, tc := range testcases {
		d, err := Parse(tc)
		require.NoError(t, err)

		require.Equal(t, tc, d.String())
	}
}

func TestDiv(t *testing.T) {
	a, err := Parse("1234567890123456000")
	require.NoError(t, err)

	b, err := Parse("0.9999999999999999999")
	require.NoError(t, err)

	c, err := a.Div(b)
	require.NoError(t, err)

	t.Logf("c = %s", c)
}

func BenchmarkString(b *testing.B) {
	a, err := Parse("1234567890123456789.123")
	require.NoError(b, err)

	bb, err := Parse("3.123456789")
	require.NoError(b, err)

	c, err := a.Div(bb)
	require.NoError(b, err)

	b.ResetTimer()
	for range b.N {
		_ = c.String()
	}
}

func BenchmarkShopspringString(b *testing.B) {
	a, err := decimal.NewFromString("1234567890123456789.123")
	require.NoError(b, err)

	bb, err := decimal.NewFromString("3.123456789")
	require.NoError(b, err)

	c := a.Div(bb)

	b.ResetTimer()
	for range b.N {
		_ = c.String()
	}
}

func TestShopspringDiv(t *testing.T) {
	a, err := decimal.NewFromString("1234567890123456789.123")
	require.NoError(t, err)

	b, err := decimal.NewFromString("3.123456789")
	require.NoError(t, err)

	c := a.DivRound(b, 20)

	t.Logf("c = %s", c.String())
}

func BenchmarkDiv(b *testing.B) {
	a, err := Parse("189")
	require.NoError(b, err)

	bb, err := Parse("3.1")
	require.NoError(b, err)

	b.ResetTimer()
	for range b.N {
		a.Div(bb)
	}
}

func BenchmarkShopspringDiv(b *testing.B) {
	a, err := decimal.NewFromString("1234567890123456789.123")
	require.NoError(b, err)

	bb, err := decimal.NewFromString("3")
	require.NoError(b, err)

	b.ResetTimer()
	for range b.N {
		a.Div(bb)
	}
}

func BenchmarkDiv2(b *testing.B) {
	b.StopTimer()

	a, err := Parse("2")
	require.NoError(b, err)

	bb, err := Parse("4")
	require.NoError(b, err)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.Div(bb)
	}
}

func TestGovaluesDiv(t *testing.T) {
	a := decimal.RequireFromString("1")
	b := decimal.RequireFromString("2")

	c := a.Div(b)

	require.Equal(t, "0.5", c.String())
}
