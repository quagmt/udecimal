package udecimal

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestSetDefaultScale(t *testing.T) {
	// NOTE: must be careful with tests that change the default scale
	// it can affect other tests, especially tests in different packages which can run in parallel
	defer SetDefaultScale(maxScale)

	require.Equal(t, uint8(19), defaultScale)

	SetDefaultScale(10)
	require.Equal(t, uint8(10), defaultScale)

	// expect panic if scale is 0
	require.PanicsWithValue(t, "scale must be greater than 0", func() {
		SetDefaultScale(0)
	})

	// expect panic if scale is > maxScale
	require.PanicsWithValue(t, fmt.Sprintf("scale out of range. Only allow maximum %d digits after the decimal points", maxScale), func() {
		SetDefaultScale(maxScale + 1)
	})

}

func TestParse(t *testing.T) {
	testcases := []struct {
		input, want string
		wantErr     error
	}{
		{"1234567890123456789.1234567890123456789", "1234567890123456789.1234567890123456789", nil},
		{"0.0000123456", "0.0000123456", nil},
		{"-0.0000123456", "-0.0000123456", nil},
		{"0.0101010101010101", "0.0101010101010101", nil},
		{"123.456000", "123.456", nil},
		{"1234567890123456789012345678901234567890", "1234567890123456789012345678901234567890", nil},
		{"-12345678912345678901.1234567890123456789", "-12345678912345678901.1234567890123456789", nil},
		{"123.0000", "123", nil},
		{"0.123", "0.123", nil},
		{"-0.123", "-0.123", nil},
		{"0", "0", nil},
		{"0.00000", "0", nil},
		{"-0", "0", nil},
		{"-0.00000", "0", nil},
		{"-123.0000", "-123", nil},
		{"0.9999999999999999999", "0.9999999999999999999", nil},
		{"-0.9999999999999999999", "-0.9999999999999999999", nil},
		{"1", "1", nil},
		{"123", "123", nil},
		{"123.456", "123.456", nil},
		{"123.456789012345678901", "123.456789012345678901", nil},
		{"123456789.123456789", "123456789.123456789", nil},
		{"-1", "-1", nil},
		{"-123", "-123", nil},
		{"-123.456", "-123.456", nil},
		{"-123.456789012345678901", "-123.456789012345678901", nil},
		{"-123456789.123456789", "-123456789.123456789", nil},
		{"-123456789123456789.123456789123456789", "-123456789123456789.123456789123456789", nil},
		{"-123456.123456", "-123456.123456", nil},
		{"1234567891234567890.0123456879123456789", "1234567891234567890.0123456879123456789", nil},
		{"9999999999999999999.9999999999999999999", "9999999999999999999.9999999999999999999", nil},
		{"-9999999999999999999.9999999999999999999", "-9999999999999999999.9999999999999999999", nil},
		{"123456.0000000000000000001", "123456.0000000000000000001", nil},
		{"-123456.0000000000000000001", "-123456.0000000000000000001", nil},
		{"+123456.123456", "123456.123456", nil},
		{"+123.123", "123.123", nil},
		{"12345678901234567890.123456789", "12345678901234567890.123456789", nil},
		{"1234567890123456789012345678901234567890", "1234567890123456789012345678901234567890", nil},
		{"1234567890123456789123456789012345678901", "1234567890123456789123456789012345678901", nil},
		{"340282366920938463463374607431768211459", "340282366920938463463374607431768211459", nil},
		{"340282366920938463463374607431768211459.123", "340282366920938463463374607431768211459.123", nil},
		{"+340282366920938463463374607431768211459", "340282366920938463463374607431768211459", nil},
		{"340282366920938463463374607431768211459.", "", fmt.Errorf("%w: can't parse '340282366920938463463374607431768211459.' to Decimal", ErrInvalidFormat)},
		{"--340282366920938463463374607431768211459", "", fmt.Errorf("%w: can't parse '--340282366920938463463374607431768211459' to Decimal", ErrInvalidFormat)},
		{".1234567890123456789012345678901234567890123456", "", fmt.Errorf("%w: can't parse '.1234567890123456789012345678901234567890123456' to Decimal", ErrInvalidFormat)},
		{"+.1234567890123456789012345678901234567890123456", "", fmt.Errorf("%w: can't parse '+.1234567890123456789012345678901234567890123456' to Decimal", ErrInvalidFormat)},
		{"-.1234567890123456789012345678901234567890123456", "", fmt.Errorf("%w: can't parse '-.1234567890123456789012345678901234567890123456' to Decimal", ErrInvalidFormat)},
		{"1.12345678901234567890123.45678901234567890123456", "", fmt.Errorf("%w: can't parse '1.12345678901234567890123.45678901234567890123456' to Decimal", ErrInvalidFormat)},
		{"340282366920938463463374607431768211459.123+--", "", fmt.Errorf("%w: can't parse '340282366920938463463374607431768211459.123+--' to Decimal", ErrInvalidFormat)},
		{"1.234567890123456789012348901", "", ErrScaleOutOfRange},
		{"1.123456789012345678912345678901234567890123456", "", ErrScaleOutOfRange},
		{"", "", ErrEmptyString},
		{".", "", fmt.Errorf("%w: can't parse '.' to Decimal", ErrInvalidFormat)},
		{"123.", "", fmt.Errorf("%w: can't parse '123.' to Decimal", ErrInvalidFormat)},
		{"-123.", "", fmt.Errorf("%w: can't parse '-123.' to Decimal", ErrInvalidFormat)},
		{"-.123456", "", fmt.Errorf("%w: can't parse '-.123456' to Decimal", ErrInvalidFormat)},
		{"12c45.123456", "", fmt.Errorf("%w: can't parse '12c45.123456' to Decimal", ErrInvalidFormat)},
		{"1245.-123456", "", fmt.Errorf("%w: can't parse '1245.-123456' to Decimal", ErrInvalidFormat)},
		{"1245.123.456", "", fmt.Errorf("%w: can't parse '1245.123.456' to Decimal", ErrInvalidFormat)},
		{"12345..123456", "", fmt.Errorf("%w: can't parse '12345..123456' to Decimal", ErrInvalidFormat)},
		{"123456.123c456", "", fmt.Errorf("%w: can't parse '123456.123c456' to Decimal", ErrInvalidFormat)},
		{"+.", "", fmt.Errorf("%w: can't parse '+.' to Decimal", ErrInvalidFormat)},
		{"+", "", fmt.Errorf("%w: can't parse '+' to Decimal", ErrInvalidFormat)},
		{"-", "", fmt.Errorf("%w: can't parse '-' to Decimal", ErrInvalidFormat)},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			d, err := Parse(tc.input)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, d.String())

			// compare with shopspring/decimal
			dd, err := decimal.NewFromString(tc.input)
			require.NoError(t, err)
			require.Equal(t, dd.String(), d.String())
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
		{"-12345678912345678901.1234567890123456789", nil},
		{"12345678901234567890.123456789", nil},
		{"1234567890123456789123456789012345678901", nil},
		{"340282366920938463463374607431768211459", nil},
		{"1.234567890123456789012348901", ErrScaleOutOfRange},
		{"", ErrEmptyString},
		{".", fmt.Errorf("%w: can't parse '.' to Decimal", ErrInvalidFormat)},
		{"123.", fmt.Errorf("%w: can't parse '123.' to Decimal", ErrInvalidFormat)},
		{"-123.", fmt.Errorf("%w: can't parse '-123.' to Decimal", ErrInvalidFormat)},
		{"-.123456", fmt.Errorf("%w: can't parse '-.123456' to Decimal", ErrInvalidFormat)},
		{"12c45.123456", fmt.Errorf("%w: can't parse '12c45.123456' to Decimal", ErrInvalidFormat)},
		{"12345..123456", fmt.Errorf("%w: can't parse '12345..123456' to Decimal", ErrInvalidFormat)},
		{"+.", fmt.Errorf("%w: can't parse '+.' to Decimal", ErrInvalidFormat)},
		{"+", fmt.Errorf("%w: can't parse '+' to Decimal", ErrInvalidFormat)},
		{"-", fmt.Errorf("%w: can't parse '-' to Decimal", ErrInvalidFormat)},
	}

	for _, tc := range testcases {
		t.Run(tc.s, func(t *testing.T) {
			if tc.wantErr != nil {
				require.PanicsWithError(t, tc.wantErr.Error(), func() {
					MustParse(tc.s)
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

func TestNewFromInt64(t *testing.T) {
	testcases := []struct {
		input   int64
		scale   uint8 // scale of decimal
		s       string
		wantErr error
	}{
		{0, 0, "0", nil},
		{0, 1, "0", nil},
		{0, 19, "0", nil},
		{1000000000000000000, 0, "1000000000000000000", nil},
		{10000, 4, "1", nil},
		{10000, 5, "0.1", nil},
		{123456000, 6, "123.456", nil},
		{0, 20, "0", ErrScaleOutOfRange},
		{0, 41, "0", ErrScaleOutOfRange},
		{1, 0, "1", nil},
		{-1, 0, "-1", nil},
		{1, 5, "0.00001", nil},
		{-1, 5, "-0.00001", nil},
		{1, 19, "0.0000000000000000001", nil},
		{-1, 19, "-0.0000000000000000001", nil},
		{math.MaxInt64, 0, "9223372036854775807", nil},
		{-math.MaxInt64, 0, "-9223372036854775807", nil},
		{math.MaxInt64, 19, "0.9223372036854775807", nil},
		{-math.MaxInt64, 19, "-0.9223372036854775807", nil},
	}

	for _, tc := range testcases {
		t.Run(strconv.FormatInt(tc.input, 10), func(t *testing.T) {
			d, err := NewFromInt64(tc.input, tc.scale)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.s, d.String())
		})
	}
}

func TestMustFromInt64(t *testing.T) {
	testcases := []struct {
		input   int64
		scale   uint8 // scale of decimal
		s       string
		wantErr error
	}{
		{0, 0, "0", nil},
		{0, 1, "0", nil},
		{0, 19, "0", nil},
		{1000000000000000000, 0, "1000000000000000000", nil},
		{10000, 4, "1", nil},
		{10000, 5, "0.1", nil},
		{123456000, 6, "123.456", nil},
		{0, 20, "0", ErrScaleOutOfRange},
		{0, 41, "0", ErrScaleOutOfRange},
		{1, 0, "1", nil},
		{-1, 0, "-1", nil},
		{1, 5, "0.00001", nil},
		{-1, 5, "-0.00001", nil},
		{1, 19, "0.0000000000000000001", nil},
		{-1, 19, "-0.0000000000000000001", nil},
		{math.MaxInt64, 0, "9223372036854775807", nil},
		{-math.MaxInt64, 0, "-9223372036854775807", nil},
		{math.MaxInt64, 19, "0.9223372036854775807", nil},
		{-math.MaxInt64, 19, "-0.9223372036854775807", nil},
	}

	for _, tc := range testcases {
		t.Run(strconv.FormatInt(tc.input, 10), func(t *testing.T) {
			if tc.wantErr != nil {
				require.PanicsWithError(t, tc.wantErr.Error(), func() {
					_ = MustFromInt64(tc.input, tc.scale)
				})
				return
			}

			d := MustFromInt64(tc.input, tc.scale)
			require.Equal(t, tc.s, d.String())
		})
	}
}

func TestNewFromUint64(t *testing.T) {
	testcases := []struct {
		input   uint64
		scale   uint8 // scale of decimal
		s       string
		wantErr error
	}{
		{0, 0, "0", nil},
		{0, 1, "0", nil},
		{0, 19, "0", nil},
		{1000000000000000000, 0, "1000000000000000000", nil},
		{10000, 4, "1", nil},
		{10000, 5, "0.1", nil},
		{123456000, 6, "123.456", nil},
		{0, 20, "0", ErrScaleOutOfRange},
		{0, 41, "0", ErrScaleOutOfRange},
		{1, 0, "1", nil},
		{1, 5, "0.00001", nil},
		{1, 19, "0.0000000000000000001", nil},
		{math.MaxUint64, 0, "18446744073709551615", nil},
		{math.MaxUint64, 19, "1.8446744073709551615", nil},
	}

	for _, tc := range testcases {
		t.Run(strconv.FormatUint(tc.input, 10), func(t *testing.T) {
			d, err := NewFromUint64(tc.input, tc.scale)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.s, d.String())
		})
	}
}

func TestMustFromUint64(t *testing.T) {
	testcases := []struct {
		input   uint64
		scale   uint8 // scale of decimal
		s       string
		wantErr error
	}{
		{0, 0, "0", nil},
		{0, 1, "0", nil},
		{0, 19, "0", nil},
		{1000000000000000000, 0, "1000000000000000000", nil},
		{10000, 4, "1", nil},
		{10000, 5, "0.1", nil},
		{123456000, 6, "123.456", nil},
		{0, 20, "0", ErrScaleOutOfRange},
		{0, 41, "0", ErrScaleOutOfRange},
		{1, 0, "1", nil},
		{1, 5, "0.00001", nil},
		{1, 19, "0.0000000000000000001", nil},
		{math.MaxUint64, 0, "18446744073709551615", nil},
		{math.MaxUint64, 19, "1.8446744073709551615", nil},
	}

	for _, tc := range testcases {
		t.Run(strconv.FormatUint(tc.input, 10), func(t *testing.T) {
			if tc.wantErr != nil {
				require.PanicsWithError(t, tc.wantErr.Error(), func() {
					_ = MustFromUint64(tc.input, tc.scale)
				})
				return
			}

			d := MustFromUint64(tc.input, tc.scale)
			require.Equal(t, tc.s, d.String())
		})
	}
}

func TestNewFromFloat64(t *testing.T) {
	testcases := []struct {
		input   float64
		s       string
		wantErr error
	}{
		{0, "0", nil},
		{0.123, "0.123", nil},
		{-0.123, "-0.123", nil},
		{1, "1", nil},
		{-1, "-1", nil},
		{1.00009, "1.00009", nil},
		{1000000.123456, "1000000.123456", nil},
		{-1000000.123456, "-1000000.123456", nil},
		{1.1234567890123456789123, "1.1234567890123457", nil},
		{123456789.1234567890123456789, "123456789.12345679", nil},
		{-1.1234567890123456789, "-1.1234567890123457", nil},
		{123.123000, "123.123", nil},
		{-123.123000, "-123.123", nil},
		{math.NaN(), "0", fmt.Errorf("%w: can't parse float 'NaN' to Decimal", ErrInvalidFormat)},
		{math.Inf(1), "0", fmt.Errorf("%w: can't parse float '+Inf' to Decimal", ErrInvalidFormat)},
		{math.Inf(-1), "0", fmt.Errorf("%w: can't parse float '-Inf' to Decimal", ErrInvalidFormat)},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("%d: %f", i, tc.input), func(t *testing.T) {
			d, err := NewFromFloat64(tc.input)
			if tc.wantErr != nil {
				require.EqualError(t, tc.wantErr, err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.s, d.String())
		})
	}
}

func TestMustFromFloat64(t *testing.T) {
	testcases := []struct {
		input   float64
		s       string
		wantErr error
	}{
		{0, "0", nil},
		{0.123, "0.123", nil},
		{-0.123, "-0.123", nil},
		{1, "1", nil},
		{-1, "-1", nil},
		{1000000.123456, "1000000.123456", nil},
		{-1000000.123456, "-1000000.123456", nil},
		{1.1234567890123456789123, "1.1234567890123457", nil},
		{123456789.1234567890123456789, "123456789.12345679", nil},
		{-1.1234567890123456789, "-1.1234567890123457", nil},
		{123.123000, "123.123", nil},
		{-123.123000, "-123.123", nil},
		{math.NaN(), "0", fmt.Errorf("%w: can't parse float 'NaN' to Decimal", ErrInvalidFormat)},
		{math.Inf(1), "0", fmt.Errorf("%w: can't parse float '+Inf' to Decimal", ErrInvalidFormat)},
		{math.Inf(-1), "0", fmt.Errorf("%w: can't parse float '-Inf' to Decimal", ErrInvalidFormat)},
	}

	for _, tc := range testcases {
		t.Run(strconv.FormatFloat(tc.input, 'f', -1, 64), func(t *testing.T) {
			if tc.wantErr != nil {
				require.PanicsWithError(t, tc.wantErr.Error(), func() {
					_ = MustFromFloat64(tc.input)
				})
				return
			}

			d := MustFromFloat64(tc.input)
			require.Equal(t, tc.s, d.String())
		})
	}
}

func assertOverflow(t *testing.T, d Decimal, isOverflow bool) {
	if isOverflow {
		require.True(t, d.coef.overflow)
		require.NotNil(t, d.coef.bigInt)
	} else {
		require.False(t, d.coef.overflow)
		require.Nil(t, d.coef.bigInt)
	}
}

func TestAdd(t *testing.T) {
	testcases := []struct {
		a, b     string
		overflow bool
	}{
		{"1", "2", false},
		{"1234567890123456789", "1234567890123456879", false},
		{"-1234567890123456789", "-1234567890123456879", false},
		{"-1234567890123456789", "1234567890123456879", false},
		{"1234567890123456789", "-1234567890123456879", false},
		{"1111111111111", "1111.123456789123456789", false},
		{"-1111111111111", "1111.123456789123456789", false},
		{"1111111111111", "-1111.123456789123456789", false},
		{"-1111111111111", "-1111.123456789123456789", false},
		{"123456789012345678.9", "0.1", false},
		{"123456789", "1.1234567890123456789", false},
		{"1234567890123456789.1234567890123456789", "1234567890123456789.1234567890123456789", false},
		{"1234567890123456789.1234567890123456789", "-1234567890123456789.1234567890123456789", false},
		{"-1234567890123456789.1234567890123456789", "1234567890123456789.1234567890123456789", false},
		{"-1234567890123456789.1234567890123456789", "-1234567890123456789.1234567890123456789", false},
		{"2345678901234567899", "1234567890123456789.1234567890123456789", false},
		{"-1111111111111", "1111.123456789123456789", false},
		{"-123456789", "1.1234567890123456789", false},
		{"-2345678901234567899", "1234567890123456789.1234567890123456789", false},
		{"1111111111111", "-1111.123456789123456789", false},
		{"123456789", "-1.1234567890123456789", false},
		{"2345678901234567899", "-1234567890123456789.1234567890123456789", false},
		{"-1111111111111", "-1111.123456789123456789", false},
		{"-123456789", "-1.1234567890123456789", false},
		{"-2345678901234567899", "-1234567890123456789.1234567890123456789", false},
		{"1", "1111.123456789123456789", false},
		{"1", "1.123456789123456789", false},
		{"123456789123456789.123456789", "3.123456789", false},
		{"123456789123456789.123456789", "3", false},
		{"9999999999999999999.9999999999999999999", "-0.999", false},
		{"-9999999999999999999.9999999999999999999", "0.999", false},
		{"0.999", "-9999999999999999999.9999999999999999999", false},
		{"-0.999", "9999999999999999999.9999999999999999999", false},
		{"9999999999999999999", "1", false},
		{"-9999999999999999999", "-1", false},
		{"9999999999999999999.99999999999999", "0.00000000000001", false},
		{"-9999999999999999999.9999999999999999999", "-0.0000000000000000001", true},
		{"9999999999999999999.9999999999999999999", "0.0000000000000000001", true},
		{"0.0000000000000000001", "9999999999999999999.9999999999999999999", true},
		{"-0.0000000000000000001", "-9999999999999999999.9999999999999999999", true},
		{"9999999999999999999.9999999999999999999", "-9999999999999999999.9999999999999999999", false},
		{"-9999999999999999999.9999999999999999999", "9999999999999999999.9999999999999999999", false},
		{"9999999999999999999.9999999999999999999", "0.999", true},
		{"0.999", "9999999999999999999.9999999999999999999", true},
		{"9999999999999999999.9999999999999999999", "999999999999999999.999", true},
		{"9999999999999999999.9999999999999999999", "9999999999999999999.9999999999999999999", true},
		{"-9999999999999999999.9999999999999999999", "-9999999999999999999.9999999999999999999", true},
		{"1234567890123456789012345678901234567890.1234567890123456789", "1234567890123456789012345678901234567890.1234567890123456789", true},
		{"1234567890123456789012345678901234567890.1234567890123456789", "-1234567890123456789012345678901234567890.1234567890123456789", true},
		{"-1234567890123456789012345678901234567890.1234567890123456789", "1234567890123456789012345678901234567890.1234567890123456789", true},
		{"-1234567890123456789012345678901234567890.1234567890123456789", "-1234567890123456789012345678901234567890.1234567890123456789", true},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"+"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c := a.Add(b)
			assertOverflow(t, c, tc.overflow)
			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.RequireFromString(tc.b)

			prec := int32(c.Scale())
			cc := aa.Add(bb).Truncate(prec)

			require.Equal(t, cc.String(), c.String())
		})
	}
}

func TestAdd64(t *testing.T) {
	testcases := []struct {
		a        string
		b        uint64
		overflow bool
	}{
		{"1234567890123456789", 1, false},
		{"1234567890123456789", 2, false},
		{"123456789012345678.9", 1, false},
		{"111111111111", 1111, false},
		{"1.1234567890123456789", 123456789, false},
		{"-123.456", 123456789, false},
		{"9999999999999999999", 1, false},
		{"-1234567890123456789.123456789", 123456789, false},
		{"1234567890123456789.123456789", math.MaxUint64, false},
		{"-1234567890123456789.123456789", math.MaxUint64, false},
		{"-1234567890123456789.123456789", 10_000_000_000_000_000_000, false},
		{"1234567890123456789.123456789", 10_000_000_000_000_000_000, false},
		{"9999999999999999999.9999999999999999999", 10_000_000_000_000_000_000, true},
		{"1234567890123456789012345678901234567890.1234567890123456789", 123456789, true},
		{"-1234567890123456789012345678901234567890.1234567890123456789", 123456789, true},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			c := a.Add64(tc.b)
			assertOverflow(t, c, tc.overflow)
			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.NewFromUint64(tc.b)

			prec := int32(c.Scale())
			cc := aa.Add(bb).Truncate(prec)

			require.Equal(t, cc.String(), c.String())
		})
	}
}

func TestSub(t *testing.T) {
	testcases := []struct {
		a, b     string
		overflow bool
	}{
		{"1", "1111.123456789123456789", false},
		{"1", "1.123456789123456789", false},
		{"1", "2", false},
		{"1", "3", false},
		{"1", "4", false},
		{"1", "5", false},
		{"1234567890123456789", "1", false},
		{"1234567890123456789", "2", false},
		{"123456789012345678.9", "0.1", false},
		{"1111111111111", "1111.123456789123456789", false},
		{"123456789", "1.1234567890123456789", false},
		{"2345678901234567899", "1234567890123456789.1234567890123456789", false},
		{"-1111111111111", "1111.123456789123456789", false},
		{"-123456789", "1.1234567890123456789", false},
		{"-2345678901234567899", "1234567890123456789.1234567890123456789", false},
		{"1111111111111", "-1111.123456789123456789", false},
		{"123456789", "-1.1234567890123456789", false},
		{"2345678901234567899", "-1234567890123456789.1234567890123456789", false},
		{"-1111111111111", "-1111.123456789123456789", false},
		{"-123456789", "-1.1234567890123456789", false},
		{"-2345678901234567899", "-1234567890123456789.1234567890123456789", false},
		{"123456789123456789.123456789", "3.123456789", false},
		{"123456789123456789.123456789", "3", false},
		{"9999999999999999999.9999999999999999999", "0.999", false},
		{"9999999999999999999", "-1", false},
		{"-9999999999999999999", "1", false},
		{"9999999999999999999.99999999999999", "-0.00000000000001", false},
		{"9999999999999999999.9999999999999999999", "-0.0000000000000000001", true},
		{"-9999999999999999999.9999999999999999999", "0.0000000000000000001", true},
		{"-0.0000000000000000001", "9999999999999999999.9999999999999999999", true},
		{"0.0000000000000000001", "-9999999999999999999.9999999999999999999", true},
		{"9999999999999999999.9999999999999999999", "-0.999", true},
		{"-9999999999999999999.9999999999999999999", "0.999", true},
		{"0.999", "-9999999999999999999.9999999999999999999", true},
		{"-0.999", "9999999999999999999.9999999999999999999", true},
		{"1234567890123456789012345678901234567890.1234567890123456789", "1234567890123456789012345678901234567890.1234567890123456789", true},
		{"1234567890123456789012345678901234567890.1234567890123456789", "-1234567890123456789012345678901234567890.1234567890123456789", true},
		{"-1234567890123456789012345678901234567890.1234567890123456789", "1234567890123456789012345678901234567890.1234567890123456789", true},
		{"-1234567890123456789012345678901234567890.1234567890123456789", "-1234567890123456789012345678901234567890.1234567890123456789", true},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c := a.Sub(b)
			assertOverflow(t, c, tc.overflow)
			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.RequireFromString(tc.b)

			prec := int32(c.Scale())
			cc := aa.Sub(bb).Truncate(prec)

			require.Equal(t, cc.String(), c.String())
		})
	}
}

func TestSub64(t *testing.T) {
	testcases := []struct {
		a        string
		b        uint64
		overflow bool
	}{
		{"1234567890123456789", 1, false},
		{"1234567890123456789", 2, false},
		{"123456789012345678.9", 1, false},
		{"111111111111", 1111, false},
		{"1.1234567890123456789", 123456789, false},
		{"-123.456", 123456789, false},
		{"-1234567890123456789.123456789", 123456789, false},
		{"1234567890123456789.123456789", 10_000_000_000_000_000_000, false},
		{"1234567890123456789.123456789", math.MaxUint64, false},
		{"-1234567890123456789.123456789", math.MaxUint64, false},
		{"-1234567890123456789.123456789", 10_000_000_000_000_000_000, false},
		{"-9999999999999999999", 1, false},
		{"-9999999999999999999.9999999999999999999", 10_000_000_000_000_000_000, true},
		{"1234567890123456789012345678901234567890.1234567890123456789", 123456789, true},
		{"-1234567890123456789012345678901234567890.1234567890123456789", 123456789, true},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s-%d", tc.a, tc.b), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			c := a.Sub64(tc.b)
			assertOverflow(t, c, tc.overflow)
			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.NewFromUint64(tc.b)

			prec := int32(c.Scale())
			cc := aa.Sub(bb).Truncate(prec)

			require.Equal(t, cc.String(), c.String())
		})
	}
}

func TestMul(t *testing.T) {
	testcases := []struct {
		a, b     string
		overflow bool
	}{
		{"123456.1234567890123456789", "0", false},
		{"123456.1234567890123456789", "123456.1234567890123456789", false},
		{"123456.1234567890123456789", "-123456.1234567890123456789", false},
		{"-123456.1234567890123456789", "123456.1234567890123456789", false},
		{"-123456.1234567890123456789", "-123456.1234567890123456789", false},
		{"9999999999999999999", "0.999", false},
		{"1234567890123456789", "1", false},
		{"1234567890123456789", "2", false},
		{"123456789012345678.9", "0.1", false},
		{"1111111111111", "1111.123456789123456789", false},
		{"123456789", "1.1234567890123456789", false},
		{"1", "1111.123456789123456789", false},
		{"1", "1.123456789123456789", false},
		{"1", "2", false},
		{"1", "3", false},
		{"1", "4", false},
		{"1", "5", false},
		{"1000000", "10000000000000", false},
		{"-1000000", "10000000000000", false},
		{"-1000000", "-10000000000000", false},
		{"1000000", "-10000000000000", false},
		{"123456789123456789.123456789", "3.123456789", false},
		{"123456789123456789.123456789", "3", false},
		{"1.123456789123456789", "1.123456789123456789", false},
		{"1234567890123456789.1234567890123456789", "123456", true},
		{"1234567890123456789.1234567890123456789", "123456.1234567890123456789", true},
		{"100000000000000000000", "100000000000000000000", true},
		{"-100000000000000000000", "100000000000000000000", true},
		{"-100000000000000000000", "-100000000000000000000", true},
		{"100000000000000000000", "-100000000000000000000", true},
		{"1000000000000000000000000.1234567890123456789", "-100000000000000000000", true},
		{"1234567890123456789012345678901234567890.1234567890123456789", "1234567890123456789012345678901234567890.1234567890123456789", true},
		{"1234567890123456789012345678901234567890.1234567890123456789", "-1234567890123456789012345678901234567890.1234567890123456789", true},
		{"-1234567890123456789012345678901234567890.1234567890123456789", "1234567890123456789012345678901234567890.1234567890123456789", true},
		{"-1234567890123456789012345678901234567890.1234567890123456789", "-1234567890123456789012345678901234567890.1234567890123456789", true},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c := a.Mul(b)
			assertOverflow(t, c, tc.overflow)
			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.RequireFromString(tc.b)

			prec := int32(c.Scale())
			cc := aa.Mul(bb).Truncate(prec)

			require.Equal(t, cc.String(), c.String())
		})
	}
}

func TestMul64(t *testing.T) {
	testcases := []struct {
		a        string
		b        uint64
		overflow bool
	}{
		{"1234567890123456789", 0, false},
		{"0", 123456789, false},
		{"1234567890123456789", 1, false},
		{"1234567890123456789", 2, false},
		{"123456789012345678.9", 1, false},
		{"111111111111", 1111, false},
		{"1.1234567890123456789", 123456789, false},
		{"-123.456", 123456789, false},
		{"0.1234567890123456789", 10_000_000_000_000_000_000, false},
		{"1000000", 10_000_000_000_000, false},
		{"-1000000", 10_000_000_000_000, false},
		{"10000000000000", 1_000_000, false},
		{"-10000000000000", 1_000_000, false},
		{"1234567890123456789.123456789", math.MaxUint64, true},
		{"9999999999999999999.9999999999999999999", 10_000_000_000_000_000_000, true},
		{"123.9999999999999999999", 10_000_000_000_000_000_000, true},
		{"1234567890123456789012345678901234567890.1234567890123456789", 123456789, true},
		{"-1234567890123456789012345678901234567890.1234567890123456789", 123456789, true},
	}

	for _, tc := range testcases {
		t.Run(tc.a, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			c := a.Mul64(tc.b)
			assertOverflow(t, c, tc.overflow)

			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.NewFromUint64(tc.b)

			prec := int32(c.Scale())
			cc := aa.Mul(bb).Truncate(prec)

			require.Equal(t, cc.String(), c.String())
		})
	}
}

func TestDiv(t *testing.T) {
	testcases := []struct {
		a, b     string
		overflow bool
		wantErr  error
	}{
		{"123456.1234567890123456789", "234567.1234567890123456789", false, nil},
		{"123456.1234567890123456789", "1", false, nil},
		{"-123456.1234567890123456789", "234567.1234567890123456789", false, nil},
		{"123456.1234567890123456789", "-234567.1234567890123456789", false, nil},
		{"-123456.1234567890123456789", "-234567.1234567890123456789", false, nil},
		{"9999999999999999999", "1.0001", false, nil},
		{"-9999999999999999999.9999999999999999999", "9999999999999999999", false, nil},
		{"1234567890123456789", "1", false, nil},
		{"1234567890123456789", "2", false, nil},
		{"123456789012345678.9", "0.1", false, nil},
		{"1111111111111", "1111.123456789123456789", false, nil},
		{"123456789", "1.1234567890123456789", false, nil},
		{"2345678901234567899", "1234567890123456789.1234567890123456789", false, nil},
		{"0.1234567890123456789", "0.04586201546101", false, nil},
		{"1", "1111.123456789123456789", false, nil},
		{"1", "1.123456789123456789", false, nil},
		{"1", "2", false, nil},
		{"1", "3", false, nil},
		{"1", "4", false, nil},
		{"1", "5", false, nil},
		{"1234567890123456789.1234567890123456879", "1111.1789", false, nil},
		{"123456789123456789.123456789", "3.123456789", false, nil},
		{"123456789123456789.123456789", "3", false, nil},
		{"9999999999999999999", "1234567890123456789.1234567890123456879", false, nil},
		{"9999999999999999999.999999999999999999", "1000000000000000000.1234567890123456789", false, nil},
		{"999999999999999999", "0.100000000000001", false, nil},
		{"123456789123456789.123456789", "0", false, ErrDivideByZero},
		{"1000000000000", "0.0000001", false, nil},
		{"1234567890123456789.1234567890123456789", "0.0000000000000000002", true, nil},
		{"1234567890123456789.1234567890123456789", "0.000000001", true, nil},
		{"1000000000000000000000000.1234567890123456789", "-100000000000000000000", true, nil},
		{"1234567890123456789012345678901234567890.1234567890123456789", "1234567890123456789012345678901234567890.1234567890123456789", true, nil},
		{"1234567890123456789012345678901234567890.1234567890123456789", "-1234567890123456789012345678901234567890.1234567890123456789", true, nil},
		{"-1234567890123456789012345678901234567890.1234567890123456789", "1234567890123456789012345678901234567890.1234567890123456789", true, nil},
		{"-1234567890123456789012345678901234567890.1234567890123456789", "-1234567890123456789012345678901234567890.1234567890123456789", true, nil},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c, err := a.Div(b)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)

			assertOverflow(t, c, tc.overflow)

			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.RequireFromString(tc.b)

			prec := int32(c.Scale())
			cc := aa.DivRound(bb, 24).Truncate(prec)

			require.Equal(t, cc.String(), c.String())
		})
	}
}

func TestDivExact(t *testing.T) {
	testcases := []struct {
		a, b     string
		scale    uint8
		overflow bool
		wantErr  error
	}{
		{"123456.1234567890123456789", "234567.1234567890123456789", 3, false, nil},
		{"123456.1234567890123456789", "1", 2, false, nil},
		{"-123456.1234567890123456789", "234567.1234567890123456789", 8, false, nil},
		{"123456.1234567890123456789", "-234567.1234567890123456789", 9, false, nil},
		{"-123456.1234567890123456789", "-234567.1234567890123456789", 1, false, nil},
		{"9999999999999999999", "1.0001", 12, false, nil},
		{"-9999999999999999999.9999999999999999999", "9999999999999999999", 16, false, nil},
		{"1234567890123456789", "1", 17, false, nil},
		{"1234567890123456789", "2", 9, false, nil},
		{"123456789012345678.9", "0.1", 8, false, nil},
		{"1111111111111", "1111.123456789123456789", 7, false, nil},
		{"123456789", "1.1234567890123456789", 10, false, nil},
		{"2345678901234567899", "1234567890123456789.1234567890123456789", 4, false, nil},
		{"0.1234567890123456789", "0.04586201546101", 5, false, nil},
		{"1", "1111.123456789123456789", 6, false, nil},
		{"1", "1.123456789123456789", 6, false, nil},
		{"1", "2", 6, false, nil},
		{"1", "3", 6, false, nil},
		{"1", "4", 6, false, nil},
		{"1", "5", 6, false, nil},
		{"1234567890123456789.1234567890123456879", "1111.1789", 11, false, nil},
		{"123456789123456789.123456789", "3.123456789", 11, false, nil},
		{"123456789123456789.123456789", "3", 11, false, nil},
		{"9999999999999999999", "1234567890123456789.1234567890123456879", 11, false, nil},
		{"9999999999999999999.999999999999999999", "1000000000000000000.1234567890123456789", 11, false, nil},
		{"999999999999999999", "0.100000000000001", 11, false, nil},
		{"123456789123456789.123456789", "0", 11, false, ErrDivideByZero},
		{"1000000000000", "0.0000001", 11, false, nil},
		{"1234567890123456789.1234567890123456789", "0.0000000000000000002", 12, true, nil},
		{"1234567890123456789.1234567890123456789", "0.000000001", 12, true, nil},
		{"1000000000000000000000000.1234567890123456789", "-100000000000000000000", 12, true, nil},
		{"1234567890123456789012345678901234567890.1234567890123456789", "1234567890123456789012345678901234567890.1234567890123456789", 12, true, nil},
		{"1234567890123456789012345678901234567890.1234567890123456789", "-1234567890123456789012345678901234567890.1234567890123456789", 12, true, nil},
		{"-1234567890123456789012345678901234567890.1234567890123456789", "1234567890123456789012345678901234567890.1234567890123456789", 12, true, nil},
		{"-1234567890123456789012345678901234567890.1234567890123456789", "-1234567890123456789012345678901234567890.1234567890123456789", 12, true, nil},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c, err := a.DivExact(b, tc.scale)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)

			assertOverflow(t, c, tc.overflow)

			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.RequireFromString(tc.b)

			cc := aa.DivRound(bb, 24).Truncate(int32(tc.scale))

			require.Equal(t, cc.String(), c.String())
		})
	}
}

func TestDivWithCustomScale(t *testing.T) {
	SetDefaultScale(14)
	defer SetDefaultScale(maxScale)

	testcases := []struct {
		a, b     string
		overflow bool
		wantErr  error
		parseErr error
	}{
		{"123456.1234567890123456789", "1", false, nil, ErrScaleOutOfRange},
		{"123456.1234567890123456789", "234567.1234567890123456789", false, nil, ErrScaleOutOfRange},
		{"-123456.1234567890123456789", "234567.1234567890123456789", false, nil, ErrScaleOutOfRange},
		{"123456.1234567890123456789", "-234567.1234567890123456789", false, nil, ErrScaleOutOfRange},
		{"-123456.1234567890123456789", "-234567.1234567890123456789", false, nil, ErrScaleOutOfRange},
		{"9999999999999999999", "1.0001", false, nil, nil},
		{"-9999999999999999999.99999999999999", "9999999999999999999", false, nil, nil},
		{"1234567890123456789", "1", false, nil, nil},
		{"1234567890123456789", "2", false, nil, nil},
		{"123456789012345678.9", "0.1", false, nil, nil},
		{"1111111111111", "1111.1234567891234", false, nil, nil},
		{"123456789", "1.12345678901234", false, nil, nil},
		{"2345678901234567899", "1234567890123456789.12345678901234", false, nil, nil},
		{"0.12345678901234", "0.04586201546101", false, nil, nil},
		{"1", "1111.1234567891234", false, nil, nil},
		{"1", "1.1234567891234", false, nil, nil},
		{"1", "2", false, nil, nil},
		{"1", "3", false, nil, nil},
		{"1", "4", false, nil, nil},
		{"1", "5", false, nil, nil},
		{"1234567890123456789.12345678901234", "1111.1789", false, nil, nil},
		{"123456789123456789.123456789", "3.123456789", false, nil, nil},
		{"123456789123456789.123456789", "3", false, nil, nil},
		{"9999999999999999999", "1234567890123456789.12345678901234", false, nil, nil},
		{"9999999999999999999.9999999999999", "1000000000000000000.12345678901234", false, nil, nil},
		{"999999999999999999", "0.1000000001", false, nil, nil},
		{"123456789123456789.123456789", "0", false, ErrDivideByZero, nil},
		{"1000000000000", "0.0000001", false, nil, nil},
		{"1234567890123456789.12345678901234", "0.00002", false, nil, nil},
		{"1234567890123456789.12345678901234", "0.000000001", true, nil, nil},
		{"1000000000000000000000000.12345678901234", "-100000000000000000000", true, nil, nil},
		{"1234567890123456789012345678901234567890.12345678901234", "1234567890123456789012345678901234567890.12345678901234", true, nil, nil},
		{"1234567890123456789012345678901234567890.12345678901234", "-1234567890123456789012345678901234567890.12345678901234", true, nil, nil},
		{"-1234567890123456789012345678901234567890.12345678901234", "1234567890123456789012345678901234567890.12345678901234", true, nil, nil},
		{"-1234567890123456789012345678901234567890.12345678901234", "-1234567890123456789012345678901234567890.12345678901234", true, nil, nil},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			if tc.parseErr != nil {
				require.Equal(t, tc.parseErr, err)
				return
			}

			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c, err := a.Div(b)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)

			assertOverflow(t, c, tc.overflow)

			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.RequireFromString(tc.b)

			prec := int32(c.Scale())
			cc := aa.DivRound(bb, 24).Truncate(prec)

			require.Equal(t, cc.String(), c.String())
		})
	}
}

func TestDiv64(t *testing.T) {
	testcases := []struct {
		a        string
		b        uint64
		overflow bool
		wantErr  error
	}{
		{"1234567890123456789", 1, false, nil},
		{"1234567890123456789", 2, false, nil},
		{"123456789012345678.9", 1, false, nil},
		{"111111111111", 1111, false, nil},
		{"1.1234567890123456789", 123456789, false, nil},
		{"-123.456", 123456789, false, nil},
		{"1234567890123456789.123456789", 10_000_000_000_000_000_000, false, nil},
		{"1234567890123456789.123456789", 123456789, false, nil},
		{"1234567890123456789.123456789", math.MaxUint64, false, nil},
		{"9999999999999999999.9999999999999999999", 9999999999999999999, false, nil},
		{"9999999999999999999.9999999999999999999", 1, false, nil},
		{"0.1234567890123456789", 1, false, nil},
		{"0.1234567890123456789", 2, false, nil},
		{"9999999999999999999", 1, false, nil},
		{"9999999999999999999", 0, false, ErrDivideByZero},
		{"1000000000000000000000000.1234567890123456789", 999_999_999_999_999, true, nil},
		{"1234567890123456789012345678901234567890.1234567890123456789", 123456789, true, nil},
		{"-1234567890123456789012345678901234567890.1234567890123456789", 123456789, true, nil},
	}

	for _, tc := range testcases {
		t.Run(tc.a, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			c, err := a.Div64(tc.b)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)

			assertOverflow(t, c, tc.overflow)

			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.NewFromUint64(tc.b)

			prec := int32(c.Scale())
			cc := aa.DivRound(bb, 24).Truncate(prec)

			require.Equal(t, cc.String(), c.String())
		})
	}
}

func TestCmp(t *testing.T) {
	testcases := []struct {
		a, b string
		want int
	}{
		{"1234567890123456789", "0", 1},
		{"123.123", "-123.123", 1},
		{"-123.123", "123.123", -1},
		{"-123.123", "-123.123", 0},
		{"-123.123", "-123.1234567890123456789", 1},
		{"123.123", "123.1234567890123456789", -1},
		{"123.123", "123.1230000000000000001", -1},
		{"-123.123", "-123.1230000000000000001", 1},
		{"123.1230000000000000002", "123.1230000000000000001", 1},
		{"-123.1230000000000000002", "-123.1230000000000000001", -1},
		{"123.1230000000000000002", "123.123000000001", -1},
		{"-123.1230000000000000002", "-123.123000000001", 1},
		{"123.123", "123.1230000", 0},
		{"123.101", "123.1001", 1},
		{"1000000000000000000000000.1234567890123456789", "1.1234567890123456789", 1},
		{"-1000000000000000000000000.1234567890123456789", "1.1234567890123456789", -1},
		{"-1000000000000000000000000.1234567890123456789", "-1.1234567890123456789", -1},
		{"1000000000000000000000000.1234567890123456789", "1000000000000000000000000.1234567890123456789", 0},
		{"-1000000000000000000000000.1234567890123456789", "-1000000000000000000000000.1234567890123456789", 0},
		{"1000000000000000000000000.1234567890123456789", "1000000000000000000000000.1234567890123456788", 1},
		{"-1000000000000000000000000.1234567890123456789", "-1000000000000000000000000.1234567890123456788", -1},
		{"1000000000000000000000000.12345678901234", "1000000000000000000000000.1234567890123456788", -1},
		{"-1000000000000000000000000.12345678901234", "-1000000000000000000000000.1234567890123456788", 1},
		{"1000000000000000000000000.1234567890123456788", "1000000000000000000000000.12345678901234", 1},
		{"-1000000000000000000000000.1234567890123456788", "-1000000000000000000000000.12345678901234", -1},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c := a.Cmp(b)
			require.Equal(t, tc.want, c)

			// compare with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			bb := decimal.RequireFromString(tc.b)

			cc := aa.Cmp(bb)
			require.Equal(t, cc, c)
		})
	}
}

func TestSign(t *testing.T) {
	testcases := []struct {
		a    string
		want int
	}{
		{"1234567890123456789", 1},
		{"123.123", 1},
		{"-123.123", -1},
		{"-123.1234567890123456789", -1},
		{"123.1234567890123456789", 1},
		{"123.1230000000000000001", 1},
		{"-123.1230000000000000001", -1},
		{"123.1230000000000000002", 1},
		{"-123.1230000000000000002", -1},
		{"123.123000000001", 1},
		{"-123.123000000001", -1},
		{"123.1230000", 1},
		{"123.1001", 1},
		{"0", 0},
		{"0.0", 0},
		{"-0", 0},
		{"-0.000", 0},
	}

	for _, tc := range testcases {
		t.Run(tc.a, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			c := a.Sign()
			require.Equal(t, tc.want, c)

			if a.coef.IsZero() {
				require.Equal(t, 0, a.Sign())
				require.True(t, a.IsZero())
				require.False(t, a.IsNeg())
				require.False(t, a.IsPos())
				return
			}

			// check neg and abs
			if a.neg {
				require.True(t, a.IsNeg())
				require.False(t, a.IsPos())
				require.Equal(t, a.Neg(), a.Abs())
			} else {
				require.True(t, a.IsPos())
				require.False(t, a.IsNeg())
				require.Equal(t, a, a.Abs())
			}
		})
	}
}

func TestRoundBank(t *testing.T) {
	testcases := []struct {
		a        string
		scale    uint8
		want     string
		overflow bool
	}{
		{"123456789012345678901234567890123456789.9999999999999999999", 3, "123456789012345678901234567890123456790", true},
		{"-123456789012345678901234567890123456789.9999999999999999999", 3, "-123456789012345678901234567890123456790", true},
		{"9999999999999999999.9999999999999999999", 3, "10000000000000000000", false},
		{"-9999999999999999999.9999999999999999999", 3, "-10000000000000000000", false},
		{"123.456000", 0, "123", false},
		{"123.456000", 1, "123.5", false},
		{"123.456000", 2, "123.46", false},
		{"123.456000", 3, "123.456", false},
		{"123.456000", 4, "123.456", false},
		{"123.456000", 5, "123.456", false},
		{"123.456000", 6, "123.456", false},
		{"123.456000", 7, "123.456", false},
		{"-123.456000", 0, "-123", false},
		{"-123.456000", 1, "-123.5", false},
		{"-123.456000", 2, "-123.46", false},
		{"-123.456000", 3, "-123.456", false},
		{"-123.456000", 4, "-123.456", false},
		{"-123.456000", 5, "-123.456", false},
		{"-123.456000", 6, "-123.456", false},
		{"-123.456000", 7, "-123.456", false},
		{"123.1234567890987654321", 0, "123", false},
		{"123.1234567890987654321", 1, "123.1", false},
		{"123.1234567890987654321", 2, "123.12", false},
		{"123.1234567890987654321", 3, "123.123", false},
		{"123.1234567890987654321", 4, "123.1235", false},
		{"123.1234567890987654321", 5, "123.12346", false},
		{"123.1234567890987654321", 6, "123.123457", false},
		{"123.1234567890987654321", 7, "123.1234568", false},
		{"123.1234567890987654321", 8, "123.12345679", false},
		{"123.1234567890987654321", 9, "123.123456789", false},
		{"123.1234567890987654321", 10, "123.1234567891", false},
		{"123.1234567890987654321", 11, "123.1234567891", false},
		{"123.1234567890987654321", 12, "123.123456789099", false},
		{"123.1234567890987654321", 13, "123.1234567890988", false},
		{"123.1234567890987654321", 14, "123.12345678909877", false},
		{"123.1234567890987654321", 15, "123.123456789098765", false},
		{"123.1234567890987654321", 16, "123.1234567890987654", false},
		{"123.1234567890987654321", 17, "123.12345678909876543", false},
		{"123.1234567890987654321", 18, "123.123456789098765432", false},
		{"123.1234567890987654321", 19, "123.1234567890987654321", false},
		{"123.1234567890987654321", 20, "123.1234567890987654321", false},
		{"-123.1234567890987654321", 0, "-123", false},
		{"-123.1234567890987654321", 1, "-123.1", false},
		{"-123.1234567890987654321", 2, "-123.12", false},
		{"-123.1234567890987654321", 3, "-123.123", false},
		{"-123.1234567890987654321", 4, "-123.1235", false},
		{"-123.1234567890987654321", 5, "-123.12346", false},
		{"-123.1234567890987654321", 6, "-123.123457", false},
		{"-123.1234567890987654321", 7, "-123.1234568", false},
		{"-123.1234567890987654321", 8, "-123.12345679", false},
		{"-123.1234567890987654321", 9, "-123.123456789", false},
		{"-123.1234567890987654321", 10, "-123.1234567891", false},
		{"-123.1234567890987654321", 11, "-123.1234567891", false},
		{"-123.1234567890987654321", 12, "-123.123456789099", false},
		{"-123.1234567890987654321", 13, "-123.1234567890988", false},
		{"-123.1234567890987654321", 14, "-123.12345678909877", false},
		{"-123.1234567890987654321", 15, "-123.123456789098765", false},
		{"-123.1234567890987654321", 16, "-123.1234567890987654", false},
		{"-123.1234567890987654321", 17, "-123.12345678909876543", false},
		{"-123.1234567890987654321", 18, "-123.123456789098765432", false},
		{"-123.1234567890987654321", 19, "-123.1234567890987654321", false},
		{"-123.1234567890987654321", 20, "-123.1234567890987654321", false},
		{"123.12354", 3, "123.124", false},
		{"-123.12354", 3, "-123.124", false},
		{"123.12454", 3, "123.125", false},
		{"-123.12454", 3, "-123.125", false},
		{"123.1235", 3, "123.124", false},
		{"-123.1235", 3, "-123.124", false},
		{"123.1245", 3, "123.124", false},
		{"-123.1245", 3, "-123.124", false},
		{"1.12345", 4, "1.1234", false},
		{"1.12335", 4, "1.1234", false},
		{"1.5", 0, "2", false},
		{"-1.5", 0, "-2", false},
		{"2.5", 0, "2", false},
		{"-2.5", 0, "-2", false},
		{"1", 0, "1", false},
		{"-1", 0, "-1", false},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s.round(%d)", tc.a, tc.scale), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			a = a.RoundBank(tc.scale)
			assertOverflow(t, a, tc.overflow)

			require.Equal(t, tc.want, a.String())

			// cross check with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			aa = aa.RoundBank(int32(tc.scale))

			require.Equal(t, aa.String(), a.String())
		})
	}
}

func TestRoundHalfAwayFromZero(t *testing.T) {
	testcases := []struct {
		a        string
		scale    uint8
		want     string
		overflow bool
	}{
		{"123456789012345678901234567890123456789.9999999999999999999", 3, "123456789012345678901234567890123456790", true},
		{"-123456789012345678901234567890123456789.9999999999999999999", 3, "-123456789012345678901234567890123456790", true},
		{"9999999999999999999.9999999999999999999", 3, "10000000000000000000", false},
		{"-9999999999999999999.9999999999999999999", 3, "-10000000000000000000", false},
		{"123.456000", 0, "123", false},
		{"123.456000", 1, "123.5", false},
		{"123.456000", 2, "123.46", false},
		{"123.456000", 3, "123.456", false},
		{"123.456000", 4, "123.456", false},
		{"123.456000", 5, "123.456", false},
		{"123.456000", 6, "123.456", false},
		{"123.456000", 7, "123.456", false},
		{"-123.456000", 0, "-123", false},
		{"-123.456000", 1, "-123.5", false},
		{"-123.456000", 2, "-123.46", false},
		{"-123.456000", 3, "-123.456", false},
		{"-123.456000", 4, "-123.456", false},
		{"-123.456000", 5, "-123.456", false},
		{"-123.456000", 6, "-123.456", false},
		{"-123.456000", 7, "-123.456", false},
		{"123.1234567890987654321", 0, "123", false},
		{"123.1234567890987654321", 1, "123.1", false},
		{"123.1234567890987654321", 2, "123.12", false},
		{"123.1234567890987654321", 3, "123.123", false},
		{"123.1234567890987654321", 4, "123.1235", false},
		{"123.1234567890987654321", 5, "123.12346", false},
		{"123.1234567890987654321", 6, "123.123457", false},
		{"123.1234567890987654321", 7, "123.1234568", false},
		{"123.1234567890987654321", 8, "123.12345679", false},
		{"123.1234567890987654321", 9, "123.123456789", false},
		{"123.1234567890987654321", 10, "123.1234567891", false},
		{"123.1234567890987654321", 11, "123.1234567891", false},
		{"123.1234567890987654321", 12, "123.123456789099", false},
		{"123.1234567890987654321", 13, "123.1234567890988", false},
		{"123.1234567890987654321", 14, "123.12345678909877", false},
		{"123.1234567890987654321", 15, "123.123456789098765", false},
		{"123.1234567890987654321", 16, "123.1234567890987654", false},
		{"123.1234567890987654321", 17, "123.12345678909876543", false},
		{"123.1234567890987654321", 18, "123.123456789098765432", false},
		{"123.1234567890987654321", 19, "123.1234567890987654321", false},
		{"123.1234567890987654321", 20, "123.1234567890987654321", false},
		{"-123.1234567890987654321", 0, "-123", false},
		{"-123.1234567890987654321", 1, "-123.1", false},
		{"-123.1234567890987654321", 2, "-123.12", false},
		{"-123.1234567890987654321", 3, "-123.123", false},
		{"-123.1234567890987654321", 4, "-123.1235", false},
		{"-123.1234567890987654321", 5, "-123.12346", false},
		{"-123.1234567890987654321", 6, "-123.123457", false},
		{"-123.1234567890987654321", 7, "-123.1234568", false},
		{"-123.1234567890987654321", 8, "-123.12345679", false},
		{"-123.1234567890987654321", 9, "-123.123456789", false},
		{"-123.1234567890987654321", 10, "-123.1234567891", false},
		{"-123.1234567890987654321", 11, "-123.1234567891", false},
		{"-123.1234567890987654321", 12, "-123.123456789099", false},
		{"-123.1234567890987654321", 13, "-123.1234567890988", false},
		{"-123.1234567890987654321", 14, "-123.12345678909877", false},
		{"-123.1234567890987654321", 15, "-123.123456789098765", false},
		{"-123.1234567890987654321", 16, "-123.1234567890987654", false},
		{"-123.1234567890987654321", 17, "-123.12345678909876543", false},
		{"-123.1234567890987654321", 18, "-123.123456789098765432", false},
		{"-123.1234567890987654321", 19, "-123.1234567890987654321", false},
		{"-123.1234567890987654321", 20, "-123.1234567890987654321", false},
		{"123.12354", 3, "123.124", false},
		{"-123.12354", 3, "-123.124", false},
		{"123.12454", 3, "123.125", false},
		{"-123.12454", 3, "-123.125", false},
		{"123.1235", 3, "123.124", false},
		{"-123.1235", 3, "-123.124", false},
		{"123.1245", 3, "123.125", false},
		{"-123.1245", 3, "-123.125", false},
		{"1.12345", 4, "1.1235", false},
		{"1.12335", 4, "1.1234", false},
		{"1.5", 0, "2", false},
		{"-1.5", 0, "-2", false},
		{"2.5", 0, "3", false},
		{"-2.5", 0, "-3", false},
		{"1", 0, "1", false},
		{"-1", 0, "-1", false},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s.roundHAZ(%d)", tc.a, tc.scale), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			a = a.RoundHAZ(tc.scale)
			assertOverflow(t, a, tc.overflow)
			require.Equal(t, tc.want, a.String())

			// cross check with shopspring/decimal
			// NOTE: shopspring/decimal roundup somehow similars to ceil, not round half up away from zero
			// Waiting this one to be merged: https://github.com/shopspring/decimal/pull/378
			// aa := decimal.RequireFromString(tc.a)
			// aa = aa.RoundUp(int32(tc.scale))

			// require.Equal(t, aa.String(), a.String())
		})
	}
}

func TestRoundHalfTowardZero(t *testing.T) {
	testcases := []struct {
		a        string
		scale    uint8
		want     string
		overflow bool
	}{
		{"123456789012345678901234567890123456789.9999999999999999999", 3, "123456789012345678901234567890123456790", true},
		{"-123456789012345678901234567890123456789.9999999999999999999", 3, "-123456789012345678901234567890123456790", true},
		{"1234567890123456789012345678912345678.5", 0, "1234567890123456789012345678912345678", false},
		{"-1234567890123456789012345678912345678.5", 0, "-1234567890123456789012345678912345678", false},
		{"9999999999999999999.9999999999999999999", 3, "10000000000000000000", false},
		{"-9999999999999999999.9999999999999999999", 3, "-10000000000000000000", false},
		{"123.456000", 0, "123", false},
		{"123.456000", 1, "123.5", false},
		{"123.456000", 2, "123.46", false},
		{"123.456000", 3, "123.456", false},
		{"123.456000", 4, "123.456", false},
		{"123.456000", 5, "123.456", false},
		{"123.456000", 6, "123.456", false},
		{"123.456000", 7, "123.456", false},
		{"-123.456000", 0, "-123", false},
		{"-123.456000", 1, "-123.5", false},
		{"-123.456000", 2, "-123.46", false},
		{"-123.456000", 3, "-123.456", false},
		{"-123.456000", 4, "-123.456", false},
		{"-123.456000", 5, "-123.456", false},
		{"-123.456000", 6, "-123.456", false},
		{"-123.456000", 7, "-123.456", false},
		{"123.1234567890987654321", 0, "123", false},
		{"123.1234567890987654321", 1, "123.1", false},
		{"123.1234567890987654321", 2, "123.12", false},
		{"123.1234567890987654321", 3, "123.123", false},
		{"123.1234567890987654321", 4, "123.1235", false},
		{"123.1234567890987654321", 5, "123.12346", false},
		{"123.1234567890987654321", 6, "123.123457", false},
		{"123.1234567890987654321", 7, "123.1234568", false},
		{"123.1234567890987654321", 8, "123.12345679", false},
		{"123.1234567890987654321", 9, "123.123456789", false},
		{"123.1234567890987654321", 10, "123.1234567891", false},
		{"123.1234567890987654321", 11, "123.1234567891", false},
		{"123.1234567890987654321", 12, "123.123456789099", false},
		{"123.1234567890987654321", 13, "123.1234567890988", false},
		{"123.1234567890987654321", 14, "123.12345678909877", false},
		{"123.1234567890987654321", 15, "123.123456789098765", false},
		{"123.1234567890987654321", 16, "123.1234567890987654", false},
		{"123.1234567890987654321", 17, "123.12345678909876543", false},
		{"123.1234567890987654321", 18, "123.123456789098765432", false},
		{"123.1234567890987654321", 19, "123.1234567890987654321", false},
		{"123.1234567890987654321", 20, "123.1234567890987654321", false},
		{"-123.1234567890987654321", 0, "-123", false},
		{"-123.1234567890987654321", 1, "-123.1", false},
		{"-123.1234567890987654321", 2, "-123.12", false},
		{"-123.1234567890987654321", 3, "-123.123", false},
		{"-123.1234567890987654321", 4, "-123.1235", false},
		{"-123.1234567890987654321", 5, "-123.12346", false},
		{"-123.1234567890987654321", 6, "-123.123457", false},
		{"-123.1234567890987654321", 7, "-123.1234568", false},
		{"-123.1234567890987654321", 8, "-123.12345679", false},
		{"-123.1234567890987654321", 9, "-123.123456789", false},
		{"-123.1234567890987654321", 10, "-123.1234567891", false},
		{"-123.1234567890987654321", 11, "-123.1234567891", false},
		{"-123.1234567890987654321", 12, "-123.123456789099", false},
		{"-123.1234567890987654321", 13, "-123.1234567890988", false},
		{"-123.1234567890987654321", 14, "-123.12345678909877", false},
		{"-123.1234567890987654321", 15, "-123.123456789098765", false},
		{"-123.1234567890987654321", 16, "-123.1234567890987654", false},
		{"-123.1234567890987654321", 17, "-123.12345678909876543", false},
		{"-123.1234567890987654321", 18, "-123.123456789098765432", false},
		{"-123.1234567890987654321", 19, "-123.1234567890987654321", false},
		{"-123.1234567890987654321", 20, "-123.1234567890987654321", false},
		{"123.12354", 3, "123.124", false},
		{"-123.12354", 3, "-123.124", false},
		{"123.12454", 3, "123.125", false},
		{"-123.12454", 3, "-123.125", false},
		{"123.1235", 3, "123.123", false},
		{"-123.1235", 3, "-123.123", false},
		{"123.1245", 3, "123.124", false},
		{"-123.1245", 3, "-123.124", false},
		{"1.12345", 4, "1.1234", false},
		{"1.12335", 4, "1.1233", false},
		{"1.5", 0, "1", false},
		{"-1.5", 0, "-1", false},
		{"2.5", 0, "2", false},
		{"-2.5", 0, "-2", false},
		{"1", 0, "1", false},
		{"-1", 0, "-1", false},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s.round(%d)", tc.a, tc.scale), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			a = a.RoundHTZ(tc.scale)
			assertOverflow(t, a, tc.overflow)

			require.Equal(t, tc.want, a.String())

			// cross check with shopspring/decimal
			// NOTE: shopspring/decimal roundup somehow similars to ceil, not round half up away from zero
			// Waiting this one to be merged: https://github.com/shopspring/decimal/pull/378
			// aa := decimal.RequireFromString(tc.a)
			// aa = aa.RoundUp(int32(tc.scale))

			// require.Equal(t, aa.String(), a.String())
		})
	}
}

func TestFloor(t *testing.T) {
	testcases := []struct {
		a        string
		want     string
		overflow bool
	}{
		{"123456789012345678901234567890123456789.9999999999999999999", "123456789012345678901234567890123456789", true},
		{"-123456789012345678901234567890123456789.9999999999999999999", "-123456789012345678901234567890123456790", true},
		{"1234567890123456789012345678912345678.5", "1234567890123456789012345678912345678", false},
		{"-1234567890123456789012345678912345678.5", "-1234567890123456789012345678912345679", false},
		{"9999999999999999999.9999999999999999999", "9999999999999999999", false},
		{"-9999999999999999999.9999999999999999999", "-10000000000000000000", false},
		{"123.456000", "123", false},
		{"123.456000", "123", false},
		{"123.456000", "123", false},
		{"123.456000", "123", false},
		{"123.456000", "123", false},
		{"123.456000", "123", false},
		{"123.456000", "123", false},
		{"123.456000", "123", false},
		{"-123.456000", "-124", false},
		{"-123.456000", "-124", false},
		{"-123.456000", "-124", false},
		{"-123.456000", "-124", false},
		{"-123.456000", "-124", false},
		{"-123.456000", "-124", false},
		{"-123.456000", "-124", false},
		{"-123.456000", "-124", false},
		{"123.1234567890987654321", "123", false},
		{"-123.1234567890987654321", "-124", false},
		{"123.12354", "123", false},
		{"-123.12354", "-124", false},
		{"123.12454", "123", false},
		{"-123.12454", "-124", false},
		{"123.1235", "123", false},
		{"-123.1235", "-124", false},
		{"123.1245", "123", false},
		{"-123.1245", "-124", false},
		{"1.12345", "1", false},
		{"1.12335", "1", false},
		{"1.5", "1", false},
		{"-1.5", "-2", false},
		{"2.5", "2", false},
		{"-2.5", "-3", false},
		{"1", "1", false},
		{"-1", "-1", false},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s.floor()", tc.a), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			a = a.Floor()
			assertOverflow(t, a, tc.overflow)

			require.Equal(t, tc.want, a.String())

			// cross check with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			aa = aa.Floor()

			require.Equal(t, aa.String(), a.String())
		})
	}
}

func TestCeil(t *testing.T) {
	testcases := []struct {
		a        string
		want     string
		overflow bool
	}{
		{"123456789012345678901234567890123456789.9999999999999999999", "123456789012345678901234567890123456790", true},
		{"-123456789012345678901234567890123456789.9999999999999999999", "-123456789012345678901234567890123456789", true},
		{"1234567890123456789012345678912345678.5", "1234567890123456789012345678912345679", false},
		{"-1234567890123456789012345678912345678.5", "-1234567890123456789012345678912345678", false},
		{"9999999999999999999.9999999999999999999", "10000000000000000000", false},
		{"-9999999999999999999.9999999999999999999", "-9999999999999999999", false},
		{"123.456000", "124", false},
		{"123.456000", "124", false},
		{"123.456000", "124", false},
		{"123.456000", "124", false},
		{"123.456000", "124", false},
		{"123.456000", "124", false},
		{"123.456000", "124", false},
		{"123.456000", "124", false},
		{"-123.456000", "-123", false},
		{"-123.456000", "-123", false},
		{"-123.456000", "-123", false},
		{"-123.456000", "-123", false},
		{"-123.456000", "-123", false},
		{"-123.456000", "-123", false},
		{"-123.456000", "-123", false},
		{"-123.456000", "-123", false},
		{"123.1234567890987654321", "124", false},
		{"-123.1234567890987654321", "-123", false},
		{"123.12354", "124", false},
		{"-123.12354", "-123", false},
		{"123.12454", "124", false},
		{"-123.12454", "-123", false},
		{"123.1235", "124", false},
		{"-123.1235", "-123", false},
		{"123.1245", "124", false},
		{"-123.1245", "-123", false},
		{"1.12345", "2", false},
		{"1.12335", "2", false},
		{"1.5", "2", false},
		{"-1.5", "-1", false},
		{"2.5", "3", false},
		{"-2.5", "-2", false},
		{"1", "1", false},
		{"-1", "-1", false},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s.floor()", tc.a), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			a = a.Ceil()
			assertOverflow(t, a, tc.overflow)

			require.Equal(t, tc.want, a.String())

			// cross check with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			aa = aa.Ceil()

			require.Equal(t, aa.String(), a.String())
		})
	}
}

func TestTrunc(t *testing.T) {
	testcases := []struct {
		a     string
		scale uint8
		want  string
	}{
		{"123456789012345678901234567890123456789.9999999999999999999", 0, "123456789012345678901234567890123456789"},
		{"-123456789012345678901234567890123456789.9999999999999999999", 0, "-123456789012345678901234567890123456789"},
		{"123456789012345678901234567890123456789.1234567890987654321", 0, "123456789012345678901234567890123456789"},
		{"123456789012345678901234567890123456789.1234567890987654321", 1, "123456789012345678901234567890123456789.1"},
		{"123456789012345678901234567890123456789.1234567890987654321", 2, "123456789012345678901234567890123456789.12"},
		{"123456789012345678901234567890123456789.1234567890987654321", 3, "123456789012345678901234567890123456789.123"},
		{"123456789012345678901234567890123456789.1234567890987654321", 4, "123456789012345678901234567890123456789.1234"},
		{"123456789012345678901234567890123456789.1234567890987654321", 5, "123456789012345678901234567890123456789.12345"},
		{"123456789012345678901234567890123456789.1234567890987654321", 6, "123456789012345678901234567890123456789.123456"},
		{"123456789012345678901234567890123456789.1234567890987654321", 7, "123456789012345678901234567890123456789.1234567"},
		{"123456789012345678901234567890123456789.1234567890987654321", 8, "123456789012345678901234567890123456789.12345678"},
		{"123456789012345678901234567890123456789.1234567890987654321", 9, "123456789012345678901234567890123456789.123456789"},
		{"123456789012345678901234567890123456789.1234567890987654321", 10, "123456789012345678901234567890123456789.123456789"},
		{"123456789012345678901234567890123456789.1234567890987654321", 11, "123456789012345678901234567890123456789.12345678909"},
		{"123456789012345678901234567890123456789.1234567890987654321", 12, "123456789012345678901234567890123456789.123456789098"},
		{"123456789012345678901234567890123456789.1234567890987654321", 13, "123456789012345678901234567890123456789.1234567890987"},
		{"123456789012345678901234567890123456789.1234567890987654321", 14, "123456789012345678901234567890123456789.12345678909876"},
		{"123456789012345678901234567890123456789.1234567890987654321", 15, "123456789012345678901234567890123456789.123456789098765"},
		{"123456789012345678901234567890123456789.1234567890987654321", 16, "123456789012345678901234567890123456789.1234567890987654"},
		{"123456789012345678901234567890123456789.1234567890987654321", 17, "123456789012345678901234567890123456789.12345678909876543"},
		{"123456789012345678901234567890123456789.1234567890987654321", 18, "123456789012345678901234567890123456789.123456789098765432"},
		{"123456789012345678901234567890123456789.1234567890987654321", 19, "123456789012345678901234567890123456789.1234567890987654321"},
		{"123456789012345678901234567890123456789.1234567890987654321", 20, "123456789012345678901234567890123456789.1234567890987654321"},
		{"1234567890123456789012345678912345678.5", 0, "1234567890123456789012345678912345678"},
		{"-1234567890123456789012345678912345678.5", 0, "-1234567890123456789012345678912345678"},
		{"9999999999999999999.9999999999999999999", 0, "9999999999999999999"},
		{"-9999999999999999999.9999999999999999999", 0, "-9999999999999999999"},
		{"123.456000", 0, "123"},
		{"123.456000", 1, "123.4"},
		{"123.456000", 2, "123.45"},
		{"123.456000", 3, "123.456"},
		{"123.456000", 4, "123.456"},
		{"123.456000", 5, "123.456"},
		{"123.456000", 6, "123.456"},
		{"123.456000", 7, "123.456"},
		{"-123.456000", 0, "-123"},
		{"-123.456000", 1, "-123.4"},
		{"-123.456000", 2, "-123.45"},
		{"-123.456000", 3, "-123.456"},
		{"-123.456000", 4, "-123.456"},
		{"-123.456000", 5, "-123.456"},
		{"-123.456000", 6, "-123.456"},
		{"-123.456000", 7, "-123.456"},
		{"123.1234567890987654321", 0, "123"},
		{"123.1234567890987654321", 1, "123.1"},
		{"123.1234567890987654321", 2, "123.12"},
		{"123.1234567890987654321", 3, "123.123"},
		{"123.1234567890987654321", 4, "123.1234"},
		{"123.1234567890987654321", 5, "123.12345"},
		{"123.1234567890987654321", 6, "123.123456"},
		{"123.1234567890987654321", 7, "123.1234567"},
		{"123.1234567890987654321", 8, "123.12345678"},
		{"123.1234567890987654321", 9, "123.123456789"},
		{"123.1234567890987654321", 10, "123.123456789"},
		{"123.1234567890987654321", 11, "123.12345678909"},
		{"123.1234567890987654321", 12, "123.123456789098"},
		{"123.1234567890987654321", 13, "123.1234567890987"},
		{"123.1234567890987654321", 14, "123.12345678909876"},
		{"123.1234567890987654321", 15, "123.123456789098765"},
		{"123.1234567890987654321", 16, "123.1234567890987654"},
		{"123.1234567890987654321", 17, "123.12345678909876543"},
		{"123.1234567890987654321", 18, "123.123456789098765432"},
		{"123.1234567890987654321", 19, "123.1234567890987654321"},
		{"123.1234567890987654321", 20, "123.1234567890987654321"},
		{"-123.1234567890987654321", 0, "-123"},
		{"-123.1234567890987654321", 1, "-123.1"},
		{"-123.1234567890987654321", 2, "-123.12"},
		{"-123.1234567890987654321", 3, "-123.123"},
		{"-123.1234567890987654321", 4, "-123.1234"},
		{"-123.1234567890987654321", 5, "-123.12345"},
		{"-123.1234567890987654321", 6, "-123.123456"},
		{"-123.1234567890987654321", 7, "-123.1234567"},
		{"-123.1234567890987654321", 8, "-123.12345678"},
		{"-123.1234567890987654321", 9, "-123.123456789"},
		{"-123.1234567890987654321", 10, "-123.123456789"},
		{"-123.1234567890987654321", 11, "-123.12345678909"},
		{"-123.1234567890987654321", 12, "-123.123456789098"},
		{"-123.1234567890987654321", 13, "-123.1234567890987"},
		{"-123.1234567890987654321", 14, "-123.12345678909876"},
		{"-123.1234567890987654321", 15, "-123.123456789098765"},
		{"-123.1234567890987654321", 16, "-123.1234567890987654"},
		{"-123.1234567890987654321", 17, "-123.12345678909876543"},
		{"-123.1234567890987654321", 18, "-123.123456789098765432"},
		{"-123.1234567890987654321", 19, "-123.1234567890987654321"},
		{"-123.1234567890987654321", 20, "-123.1234567890987654321"},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s.trunc(%d)", tc.a, tc.scale), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			a = a.Trunc(tc.scale)
			require.Equal(t, tc.want, a.String())

			// cross check with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			aa = aa.Truncate(int32(tc.scale))

			require.Equal(t, aa.String(), a.String())
		})
	}
}

func TestPowInt(t *testing.T) {
	testcases := []struct {
		a    string
		b    int
		want string
	}{
		{"0.12", 100, "0"},
		{"0", 1, "0"},
		{"0", 10, "0"},
		{"1.12345", 4, "1.5929971334827095062"},
		{"123456789012345678901234567890123456789.9999999999999999999", 2, "15241578753238836750495351562566681945252248135650053345652796829976527968319.753086421975308642"},
		{"123456789012345678901234567890123456789.9999999999999999999", 0, "1"},
		{"123456789012345678901234567890123456789.9999999999999999999", 1, "123456789012345678901234567890123456789.9999999999999999999"},
		{"1.5", 3, "3.375"},
		{"1.12345", 1, "1.12345"},
		{"1.12345", 2, "1.2621399025"},
		{"1.12345", 3, "1.417951073463625"},
		{"1.12345", 4, "1.5929971334827095062"},
		{"1.12345", 5, "1.7896526296111499947"},
		{"1.12345", 6, "2.0105852467366464616"},
		{"1.12345", 7, "2.2587919954462854673"},
		{"-1.12345", 4, "1.5929971334827095062"},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s.pow(%d)", tc.a, tc.b), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			a = a.PowInt(tc.b)
			require.Equal(t, tc.want, a.String())

			// cross check with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			aa, err = aa.PowWithPrecision(decimal.New(int64(tc.b), 0), int32(a.scale)+4)
			require.NoError(t, err)

			aa = aa.Truncate(int32(a.scale))

			require.Equal(t, aa.String(), a.String())
		})
	}
}

func TestPowIntRandom(t *testing.T) {
	inputs := []string{
		"0.1234",
		"-0.1234",
		"1.123456789012345679",
		"-1.123456789012345679",
		"1.12345",
		"-1.12345",
		"123456789012345678901234567890123456789.9999999999999999999",
		"123456789012345678901234567890123456789.9999999999999999999",
		"1.5",
		"123456.789",
		"123.4",
		"1234567890123456789.1234567890123456789",
		"-1234567890123456789.1234567890123456789",
	}

	for _, input := range inputs {
		t.Run(fmt.Sprintf("pow(%s)", input), func(t *testing.T) {
			a := MustParse(input)

			for i := 0; i <= 100; i++ {
				b := a.PowInt(i)

				aa := decimal.RequireFromString(input)
				aa, err := aa.PowWithPrecision(decimal.New(int64(i), 0), int32(b.scale)+4)
				require.NoError(t, err)

				aa = aa.Truncate(int32(b.scale))

				require.Equal(t, aa.String(), b.String(), "%s.pow(%d)", input, i)
			}
		})
	}

	for _, input := range inputs {
		t.Run(fmt.Sprintf("powInverse(%s)", input), func(t *testing.T) {
			a := MustParse(input)

			for i := 0; i >= -100; i-- {
				b := a.PowInt(i)

				aa := decimal.RequireFromString(input)
				aa, err := aa.PowWithPrecision(decimal.New(int64(i), 0), int32(b.scale)+4)
				require.NoError(t, err)

				aa = aa.Truncate(int32(b.scale))

				require.Equal(t, aa.String(), b.String(), "%s.pow(%d)", input, i)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	testcases := []struct {
		a       string
		want    string
		wantErr error
	}{
		{"10000000000", "100000", nil},
		{"3", "1.7320508075688772935", nil},
		{"-1", "", ErrSqrtNegative},
		{"0", "0", nil},
		{"1", "1", nil},
		{"2", "1.4142135623730950488", nil},
		{"1000", "31.6227766016837933199", nil},
		{"31.6227766016837933199", "5.6234132519034908039", nil},
		{"4", "2", nil},
		{"12345678901234567890.1234567890123456789", "3513641828.8201442531112223816", nil},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("sqrt(%s)", tc.a), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			a, err = a.Sqrt()
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, a.String())

			// cross check with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			aa, err = aa.PowWithPrecision(decimal.RequireFromString("0.5"), int32(a.scale)+4)
			require.NoError(t, err)

			a1 := decimal.RequireFromString(a.String()).Sub(aa).Truncate(int32(a.scale))
			require.True(t, a1.IsZero())
		})
	}
}

func TestRandomSqrt(t *testing.T) {
	// from 0.1 to 100
	for i := 1; i <= 1000; i++ {
		input := fmt.Sprintf("%f", float64(i)/10)

		a, err := Parse(input)
		require.NoError(t, err)

		a, err = a.Sqrt()
		require.NoError(t, err)

		// cross check with shopspring/decimal
		aa := decimal.RequireFromString(input)
		aa, err = aa.PowWithPrecision(decimal.RequireFromString("0.5"), int32(a.scale)+4)
		require.NoError(t, err)

		a1 := decimal.RequireFromString(a.String()).Sub(aa).Truncate(int32(a.scale))
		require.True(t, a1.IsZero())
	}
}

func TestInexactFloat64(t *testing.T) {
	testcases := []struct {
		a    string
		want float64
	}{
		{"0", 0},
		{"1", 1},
		{"1.12345", 1.12345},
		{"-1.12345", -1.12345},
		{"123456789.123456789", 123456789.123456789},
		{"-123456789.123456789", -123456789.123456789},
		{"1234567890123456789.1234567890123456789", 1234567890123456789.1234567890123456789},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("inexactFloat64(%s)", tc.a), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			got, err := a.InexactFloat64()
			require.NoError(t, err)

			require.Equal(t, tc.want, got)

			// cross check with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			got1, _ := aa.Float64()

			require.Equal(t, got1, got)
		})
	}
}

func BenchmarkString(b *testing.B) {
	a, err := Parse("12345.12345")
	require.NoError(b, err)

	b.ResetTimer()
	for range b.N {
		_ = a.String()
	}
}

func BenchmarkShopspringString(b *testing.B) {
	a, err := decimal.NewFromString("12345.12345")
	require.NoError(b, err)

	b.ResetTimer()
	for range b.N {
		_ = a.String()
	}
}

func BenchmarkMul(b *testing.B) {
	a, err := Parse("1234567890")
	require.NoError(b, err)

	bb, err := Parse("1111.123456789123456789")
	require.NoError(b, err)

	b.ResetTimer()
	for range b.N {
		a.Mul(bb)
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
	a, err := Parse("1234567890123456789.1234567890123456879")
	require.NoError(b, err)

	bb, err := Parse("1111.1789")
	require.NoError(b, err)

	b.ResetTimer()
	for range b.N {
		_, _ = a.Div(bb)
	}
}

func BenchmarkShopspringDiv(b *testing.B) {
	a, err := decimal.NewFromString("1234567890123456789.1234567890123456879")
	require.NoError(b, err)

	bb, err := decimal.NewFromString("1111.1789")
	require.NoError(b, err)

	b.ResetTimer()
	for range b.N {
		a.Div(bb)
	}
}

func BenchmarkPow(b *testing.B) {
	a, err := Parse("12.46")
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.PowInt(10)
	}
}
