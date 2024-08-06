package udecimal

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	testcases := []string{
		"0.123",
		"-0.123",
		"0",
		"0.9999999999999999999",
		"-0.9999999999999999999",
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

func TestAdd(t *testing.T) {
	testcases := []struct {
		a, b string
	}{
		{"1234567890123456789", "1"},
		{"1234567890123456789", "2"},
		{"123456789012345678.9", "0.1"},
		{"1111111111111", "1111.123456789123456789"},
		{"123456789", "1.1234567890123456789"},
		{"2345678901234567899", "1234567890123456789.1234567890123456789"},
		{"-1111111111111", "1111.123456789123456789"},
		{"-123456789", "1.1234567890123456789"},
		{"-2345678901234567899", "1234567890123456789.1234567890123456789"},
		{"1111111111111", "-1111.123456789123456789"},
		{"123456789", "-1.1234567890123456789"},
		{"2345678901234567899", "-1234567890123456789.1234567890123456789"},
		{"-1111111111111", "-1111.123456789123456789"},
		{"-123456789", "-1.1234567890123456789"},
		{"-2345678901234567899", "-1234567890123456789.1234567890123456789"},
		{"1", "1111.123456789123456789"},
		{"1", "1.123456789123456789"},
		{"1", "2"},
		{"1", "3"},
		{"1", "4"},
		{"1", "5"},
		{"123456789123456789.123456789", "3.123456789"},
		{"123456789123456789.123456789", "3"},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c, err := a.Add(b)
			require.NoError(t, err)

			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.RequireFromString(tc.b)

			prec := int32(c.scale)
			cc := aa.Add(bb).Truncate(prec)

			require.Equal(t, cc.String(), c.String())
		})
	}
}

func TestSub(t *testing.T) {
	testcases := []struct {
		a, b string
	}{
		{"1234567890123456789", "1"},
		{"1234567890123456789", "2"},
		{"123456789012345678.9", "0.1"},
		{"1111111111111", "1111.123456789123456789"},
		{"123456789", "1.1234567890123456789"},
		{"2345678901234567899", "1234567890123456789.1234567890123456789"},
		{"-1111111111111", "1111.123456789123456789"},
		{"-123456789", "1.1234567890123456789"},
		{"-2345678901234567899", "1234567890123456789.1234567890123456789"},
		{"1111111111111", "-1111.123456789123456789"},
		{"123456789", "-1.1234567890123456789"},
		{"2345678901234567899", "-1234567890123456789.1234567890123456789"},
		{"-1111111111111", "-1111.123456789123456789"},
		{"-123456789", "-1.1234567890123456789"},
		{"-2345678901234567899", "-1234567890123456789.1234567890123456789"},
		{"1", "1111.123456789123456789"},
		{"1", "1.123456789123456789"},
		{"1", "2"},
		{"1", "3"},
		{"1", "4"},
		{"1", "5"},
		{"123456789123456789.123456789", "3.123456789"},
		{"123456789123456789.123456789", "3"},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c, err := a.Sub(b)
			require.NoError(t, err)

			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.RequireFromString(tc.b)

			prec := int32(c.scale)
			cc := aa.Sub(bb).Truncate(prec)

			require.Equal(t, cc.String(), c.String())
		})
	}
}

func TestDiv(t *testing.T) {
	testcases := []struct {
		a, b string
	}{
		{"1234567890123456789", "1"},
		{"1234567890123456789", "2"},
		{"123456789012345678.9", "0.1"},
		{"1111111111111", "1111.123456789123456789"},
		{"123456789", "1.1234567890123456789"},
		{"2345678901234567899", "1234567890123456789.1234567890123456789"},
		{"1", "1111.123456789123456789"},
		{"1", "1.123456789123456789"},
		{"1", "2"},
		{"1", "3"},
		{"1", "4"},
		{"1", "5"},
		{"123456789123456789.123456789", "3.123456789"},
		{"123456789123456789.123456789", "3"},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c, err := a.Div(b)
			require.NoError(t, err)

			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.RequireFromString(tc.b)

			prec := int32(c.scale)
			cc := aa.DivRound(bb, prec+1).Truncate(prec)

			require.Equal(t, cc.String(), c.String())
		})
	}
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
	// {"2345678901234567899", "1234567890123456789.1234567890123456789"},

	a, err := Parse("1234567890123456789")
	require.NoError(b, err)

	bb, err := Parse("1111.1789")
	require.NoError(b, err)

	b.ResetTimer()
	for range b.N {
		_, _ = a.Div(bb)
	}
}

func BenchmarkShopspringDiv(b *testing.B) {
	a, err := decimal.NewFromString("1234567890123456789")
	require.NoError(b, err)

	bb, err := decimal.NewFromString("1111.123456789123456789")
	require.NoError(b, err)

	b.ResetTimer()
	for range b.N {
		a.Div(bb)
	}
}

func TestGovaluesDiv(t *testing.T) {
	a := decimal.RequireFromString("1")
	b := decimal.RequireFromString("2")

	c := a.Div(b)

	require.Equal(t, "0.5", c.String())
}
