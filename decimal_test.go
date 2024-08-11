package udecimal

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	testcases := []struct {
		s       string
		wantErr error
	}{
		{"0.123", nil},
		{"-0.123", nil},
		{"0", nil},
		{"0.9999999999999999999", nil},
		{"-0.9999999999999999999", nil},
		{"1", nil},
		{"123", nil},
		{"123.456", nil},
		{"123.456789012345678901", nil},
		{"123456789.123456789", nil},
		{"-1", nil},
		{"-123", nil},
		{"-123.456", nil},
		{"-123.456789012345678901", nil},
		{"-123456789.123456789", nil},
		{"-123456789123456789.123456789123456789", nil},
		{"-123456.123456", nil},
		{"1234567891234567890.0123456879123456789", nil},
		{"9999999999999999999.9999999999999999999", nil},
		{"-9999999999999999999.9999999999999999999", nil},
		{"123456.0000000000000000001", nil},
		{"-123456.0000000000000000001", nil},
		{"+123456.123456", nil},
		{"+123.123", nil},
		{"", ErrEmptyString},
		{".", fmt.Errorf("%w: can't parse '.' to decimal", ErrInvalidFormat)},
		{"123.", fmt.Errorf("%w: can't parse '123.' to decimal", ErrInvalidFormat)},
		{"-123.", fmt.Errorf("%w: can't parse '-123.' to decimal", ErrInvalidFormat)},
		{"-.123456", fmt.Errorf("%w: can't parse '-.123456' to decimal", ErrInvalidFormat)},
		{"12c45.123456", fmt.Errorf("%w: can't parse '12c45.123456' to decimal", ErrInvalidFormat)},
		{"1245.-123456", fmt.Errorf("%w: can't parse '1245.-123456' to decimal", ErrInvalidFormat)},
		{"1245.123.456", fmt.Errorf("%w: can't parse '1245.123.456' to decimal", ErrInvalidFormat)},
		{"12345..123456", fmt.Errorf("%w: can't parse '12345..123456' to decimal", ErrInvalidFormat)},
		{"123456.123c456", fmt.Errorf("%w: can't parse '123456.123c456' to decimal", ErrInvalidFormat)},
		{"+.", fmt.Errorf("%w: can't parse '+.' to decimal", ErrInvalidFormat)},
		{"-12345678912345678901.1234567890123456789", fmt.Errorf("%w: overflow. string length is greater than 40", ErrInvalidFormat)},
		{"12345678901234567890.123456789", ErrOverflow},
		{"1234567890123456789123456789012345678901", ErrOverflow},
		{"340282366920938463463374607431768211459", ErrOverflow},
		{"1.234567890123456789012348901", ErrMaxScale},
		{"+", fmt.Errorf("%w: can't parse '+' to decimal", ErrInvalidFormat)},
		{"-", fmt.Errorf("%w: can't parse '-' to decimal", ErrInvalidFormat)},
	}

	for _, tc := range testcases {
		t.Run(tc.s, func(t *testing.T) {
			d, err := Parse(tc.s)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)
			if tc.s[0] == '+' {
				require.Equal(t, tc.s[1:], d.String())
			} else {
				require.Equal(t, tc.s, d.String())
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	testcases := []struct {
		s       string
		wantErr error
	}{
		{"0.123", nil},
		{"-0.123", nil},
		{"0", nil},
		{"0.9999999999999999999", nil},
		{"-0.9999999999999999999", nil},
		{"1", nil},
		{"123", nil},
		{"123.456", nil},
		{"123.456789012345678901", nil},
		{"123456789.123456789", nil},
		{"-123456789123456789.123456789123456789", nil},
		{"-123456.123456", nil},
		{"1234567891234567890.0123456879123456789", nil},
		{"9999999999999999999.9999999999999999999", nil},
		{"-9999999999999999999.9999999999999999999", nil},
		{"123456.0000000000000000001", nil},
		{"-123456.0000000000000000001", nil},
		{"+123456.123456", nil},
		{"+123.123", nil},
		{"", ErrEmptyString},
		{".", fmt.Errorf("%w: can't parse '.' to decimal", ErrInvalidFormat)},
		{"123.", fmt.Errorf("%w: can't parse '123.' to decimal", ErrInvalidFormat)},
		{"-123.", fmt.Errorf("%w: can't parse '-123.' to decimal", ErrInvalidFormat)},
		{"-.123456", fmt.Errorf("%w: can't parse '-.123456' to decimal", ErrInvalidFormat)},
		{"12c45.123456", fmt.Errorf("%w: can't parse '12c45.123456' to decimal", ErrInvalidFormat)},
		{"12345..123456", fmt.Errorf("%w: can't parse '12345..123456' to decimal", ErrInvalidFormat)},
		{"+.", fmt.Errorf("%w: can't parse '+.' to decimal", ErrInvalidFormat)},
		{"-12345678912345678901.1234567890123456789", fmt.Errorf("%w: overflow. string length is greater than 40", ErrInvalidFormat)},
		{"12345678901234567890.123456789", ErrOverflow},
		{"1234567890123456789123456789012345678901", ErrOverflow},
		{"340282366920938463463374607431768211459", ErrOverflow},
		{"1.234567890123456789012348901", ErrMaxScale},
		{"+", fmt.Errorf("%w: can't parse '+' to decimal", ErrInvalidFormat)},
		{"-", fmt.Errorf("%w: can't parse '-' to decimal", ErrInvalidFormat)},
	}

	for _, tc := range testcases {
		t.Run(tc.s, func(t *testing.T) {
			if tc.wantErr != nil {
				require.PanicsWithError(t, tc.wantErr.Error(), func() {
					_ = MustParse(tc.s)
				})
				return
			}

			var d Decimal
			require.NotPanics(t, func() {
				d = MustParse(tc.s)
			})

			if tc.s[0] == '+' {
				require.Equal(t, tc.s[1:], d.String())
			} else {
				require.Equal(t, tc.s, d.String())
			}
		})
	}
}

// func TestNewFromInt64(t *testing.T) {
// 	testcases := []struct {
// 		i       int64
// 		s       string
// 		wantErr error
// 	}{}
// }

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

func TestAdd64(t *testing.T) {
	testcases := []struct {
		a string
		b uint64
	}{
		{"1234567890123456789", 1},
		{"1234567890123456789", 2},
		{"123456789012345678.9", 1},
		{"111111111111", 1111},
		{"1.1234567890123456789", 123456789},
		{"-123.456", 123456789},
		{"-1234567890123456789.123456789", 123456789},
	}

	for _, tc := range testcases {
		t.Run(tc.a, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			c, err := a.Add64(tc.b)
			require.NoError(t, err)

			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.NewFromInt(int64(tc.b))

			prec := int32(c.scale)
			cc := aa.Add(bb).Truncate(prec)

			fmt.Println("c = ", c.String())

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

func TestMul(t *testing.T) {
	testcases := []struct {
		a, b string
	}{
		{"123456.1234567890123456789", "0"},
		{"123456.1234567890123456789", "123456.1234567890123456789"},
		{"123456.1234567890123456789", "-123456.1234567890123456789"},
		{"-123456.1234567890123456789", "123456.1234567890123456789"},
		{"-123456.1234567890123456789", "-123456.1234567890123456789"},
		{"9999999999999999999", "0.999"},
		{"1234567890123456789", "1"},
		{"1234567890123456789", "2"},
		{"123456789012345678.9", "0.1"},
		{"1111111111111", "1111.123456789123456789"},
		{"123456789", "1.1234567890123456789"},
		{"1", "1111.123456789123456789"},
		{"1", "1.123456789123456789"},
		{"1", "2"},
		{"1", "3"},
		{"1", "4"},
		{"1", "5"},
		{"123456789123456789.123456789", "3.123456789"},
		{"123456789123456789.123456789", "3"},
		{"1.123456789123456789", "1.123456789123456789"},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c, err := a.Mul(b)
			require.NoError(t, err)

			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.RequireFromString(tc.b)

			prec := int32(c.scale)
			cc := aa.Mul(bb).Truncate(prec)

			require.Equal(t, cc.String(), c.String())
		})
	}
}

func TestDiv(t *testing.T) {
	testcases := []struct {
		a, b string
	}{
		{"123456.1234567890123456789", "234567.1234567890123456789"},
		{"-123456.1234567890123456789", "234567.1234567890123456789"},
		{"123456.1234567890123456789", "-234567.1234567890123456789"},
		{"-123456.1234567890123456789", "-234567.1234567890123456789"},
		{"9999999999999999999", "1.0001"},
		{"-9999999999999999999.9999999999999999999", "9999999999999999999"},
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

func BenchmarkMul(b *testing.B) {
	a, err := Parse("1234567890")
	require.NoError(b, err)

	bb, err := Parse("1111.123456789123456789")
	require.NoError(b, err)

	b.ResetTimer()
	for range b.N {
		_, _ = a.Mul(bb)
	}
}

func BenchmarkShopspringMul(b *testing.B) {
	a, err := decimal.NewFromString("1234567890")
	require.NoError(b, err)

	bb, err := decimal.NewFromString("1111.123456789123456789")
	require.NoError(b, err)

	b.ResetTimer()
	for range b.N {
		a.Mul(bb)
	}
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
