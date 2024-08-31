package udecimal

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	testcases := []struct {
		input, want string
		wantErr     error
	}{
		{"123.456000", "123.456", nil},
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
		{"", "", ErrEmptyString},
		{".", "", fmt.Errorf("%w: can't parse '.' to decimal", ErrInvalidFormat)},
		{"123.", "", fmt.Errorf("%w: can't parse '123.' to decimal", ErrInvalidFormat)},
		{"-123.", "", fmt.Errorf("%w: can't parse '-123.' to decimal", ErrInvalidFormat)},
		{"-.123456", "", fmt.Errorf("%w: can't parse '-.123456' to decimal", ErrInvalidFormat)},
		{"12c45.123456", "", fmt.Errorf("%w: can't parse '12c45.123456' to decimal", ErrInvalidFormat)},
		{"1245.-123456", "", fmt.Errorf("%w: can't parse '1245.-123456' to decimal", ErrInvalidFormat)},
		{"1245.123.456", "", fmt.Errorf("%w: can't parse '1245.123.456' to decimal", ErrInvalidFormat)},
		{"12345..123456", "", fmt.Errorf("%w: can't parse '12345..123456' to decimal", ErrInvalidFormat)},
		{"123456.123c456", "", fmt.Errorf("%w: can't parse '123456.123c456' to decimal", ErrInvalidFormat)},
		{"+.", "", fmt.Errorf("%w: can't parse '+.' to decimal", ErrInvalidFormat)},
		{"-12345678912345678901.1234567890123456789", "", fmt.Errorf("%w: string length is greater than 40", ErrInvalidFormat)},
		{"12345678901234567890.123456789", "", ErrOverflow},
		{"1234567890123456789123456789012345678901", "", ErrOverflow},
		{"340282366920938463463374607431768211459", "", ErrOverflow},
		{"1.234567890123456789012348901", "", ErrMaxScale},
		{"+", "", fmt.Errorf("%w: can't parse '+' to decimal", ErrInvalidFormat)},
		{"-", "", fmt.Errorf("%w: can't parse '-' to decimal", ErrInvalidFormat)},
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
		{"", ErrEmptyString},
		{".", fmt.Errorf("%w: can't parse '.' to decimal", ErrInvalidFormat)},
		{"123.", fmt.Errorf("%w: can't parse '123.' to decimal", ErrInvalidFormat)},
		{"-123.", fmt.Errorf("%w: can't parse '-123.' to decimal", ErrInvalidFormat)},
		{"-.123456", fmt.Errorf("%w: can't parse '-.123456' to decimal", ErrInvalidFormat)},
		{"12c45.123456", fmt.Errorf("%w: can't parse '12c45.123456' to decimal", ErrInvalidFormat)},
		{"12345..123456", fmt.Errorf("%w: can't parse '12345..123456' to decimal", ErrInvalidFormat)},
		{"+.", fmt.Errorf("%w: can't parse '+.' to decimal", ErrInvalidFormat)},
		{"-12345678912345678901.1234567890123456789", fmt.Errorf("%w: string length is greater than 40", ErrInvalidFormat)},
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
		{0, 20, "0", ErrMaxScale},
		{0, 40, "0", ErrMaxScale},
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
		{0, 20, "0", ErrMaxScale},
		{0, 40, "0", ErrMaxScale},
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
		{math.NaN(), "0", fmt.Errorf("%w: can't parse float 'NaN' to decimal", ErrInvalidFormat)},
		{math.Inf(1), "0", fmt.Errorf("%w: can't parse float '+Inf' to decimal", ErrInvalidFormat)},
		{math.Inf(-1), "0", fmt.Errorf("%w: can't parse float '-Inf' to decimal", ErrInvalidFormat)},
		{math.MaxFloat64, "0", fmt.Errorf("can't parse float: %w: string length is greater than 40", ErrInvalidFormat)},
		{-math.MaxFloat64, "0", fmt.Errorf("can't parse float: %w: string length is greater than 40", ErrInvalidFormat)},
	}

	for _, tc := range testcases {
		t.Run(strconv.FormatFloat(tc.input, 'f', -1, 64), func(t *testing.T) {
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
		{math.NaN(), "0", fmt.Errorf("%w: can't parse float 'NaN' to decimal", ErrInvalidFormat)},
		{math.Inf(1), "0", fmt.Errorf("%w: can't parse float '+Inf' to decimal", ErrInvalidFormat)},
		{math.Inf(-1), "0", fmt.Errorf("%w: can't parse float '-Inf' to decimal", ErrInvalidFormat)},
		{math.MaxFloat64, "0", fmt.Errorf("can't parse float: %w: string length is greater than 40", ErrInvalidFormat)},
		{-math.MaxFloat64, "0", fmt.Errorf("can't parse float: %w: string length is greater than 40", ErrInvalidFormat)},
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

func TestAdd(t *testing.T) {
	testcases := []struct {
		a, b    string
		wantErr error
	}{
		{"1", "2", nil},
		{"1234567890123456789", "1234567890123456879", nil},
		{"-1234567890123456789", "-1234567890123456879", nil},
		{"-1234567890123456789", "1234567890123456879", nil},
		{"1234567890123456789", "-1234567890123456879", nil},
		{"1111111111111", "1111.123456789123456789", nil},
		{"-1111111111111", "1111.123456789123456789", nil},
		{"1111111111111", "-1111.123456789123456789", nil},
		{"-1111111111111", "-1111.123456789123456789", nil},
		{"123456789012345678.9", "0.1", nil},
		{"123456789", "1.1234567890123456789", nil},
		{"1234567890123456789.1234567890123456789", "1234567890123456789.1234567890123456789", nil},
		{"1234567890123456789.1234567890123456789", "-1234567890123456789.1234567890123456789", nil},
		{"-1234567890123456789.1234567890123456789", "1234567890123456789.1234567890123456789", nil},
		{"-1234567890123456789.1234567890123456789", "-1234567890123456789.1234567890123456789", nil},
		{"2345678901234567899", "1234567890123456789.1234567890123456789", nil},
		{"-1111111111111", "1111.123456789123456789", nil},
		{"-123456789", "1.1234567890123456789", nil},
		{"-2345678901234567899", "1234567890123456789.1234567890123456789", nil},
		{"1111111111111", "-1111.123456789123456789", nil},
		{"123456789", "-1.1234567890123456789", nil},
		{"2345678901234567899", "-1234567890123456789.1234567890123456789", nil},
		{"-1111111111111", "-1111.123456789123456789", nil},
		{"-123456789", "-1.1234567890123456789", nil},
		{"-2345678901234567899", "-1234567890123456789.1234567890123456789", nil},
		{"1", "1111.123456789123456789", nil},
		{"1", "1.123456789123456789", nil},
		{"123456789123456789.123456789", "3.123456789", nil},
		{"123456789123456789.123456789", "3", nil},
		{"9999999999999999999.9999999999999999999", "-0.999", nil},
		{"-9999999999999999999.9999999999999999999", "0.999", nil},
		{"0.999", "-9999999999999999999.9999999999999999999", nil},
		{"-0.999", "9999999999999999999.9999999999999999999", nil},
		{"9999999999999999999.9999999999999999999", "-9999999999999999999.9999999999999999999", nil},
		{"-9999999999999999999.9999999999999999999", "9999999999999999999.9999999999999999999", nil},
		{"9999999999999999999.9999999999999999999", "0.999", ErrOverflow},
		{"0.999", "9999999999999999999.9999999999999999999", ErrOverflow},
		{"9999999999999999999.9999999999999999999", "999999999999999999.999", ErrOverflow},
		{"9999999999999999999.9999999999999999999", "9999999999999999999.9999999999999999999", ErrOverflow},
		{"-9999999999999999999.9999999999999999999", "-9999999999999999999.9999999999999999999", ErrOverflow},
		{"9999999999999999999", "1", ErrOverflow},
		{"-9999999999999999999", "-1", ErrOverflow},
		{"9999999999999999999.9999999999999999999", "0.0000000000000000001", ErrOverflow},
		{"-9999999999999999999.9999999999999999999", "-0.0000000000000000001", ErrOverflow},
		{"9999999999999999999.99999999999999", "0.00000000000001", ErrOverflow},
		{"0.0000000000000000001", "9999999999999999999.9999999999999999999", ErrOverflow},
		{"-0.0000000000000000001", "-9999999999999999999.9999999999999999999", ErrOverflow},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"+"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c, err := a.Add(b)
			if tc.wantErr != nil {
				require.EqualError(t, tc.wantErr, err.Error())
				return
			}

			require.NoError(t, err)

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
		a       string
		b       uint64
		wantErr error
	}{
		{"1234567890123456789", 1, nil},
		{"1234567890123456789", 2, nil},
		{"123456789012345678.9", 1, nil},
		{"111111111111", 1111, nil},
		{"1.1234567890123456789", 123456789, nil},
		{"-123.456", 123456789, nil},
		{"-1234567890123456789.123456789", 10_000_000_000_000_000_000, nil},
		{"-1234567890123456789.123456789", 123456789, nil},
		{"-1234567890123456789.123456789", math.MaxUint64, ErrOverflow},
		{"1234567890123456789.123456789", 10_000_000_000_000_000_000, ErrOverflow},
		{"9999999999999999999.9999999999999999999", 10_000_000_000_000_000_000, ErrOverflow},
		{"9999999999999999999", 1, ErrOverflow},
	}

	for _, tc := range testcases {
		t.Run(tc.a, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			c, err := a.Add64(tc.b)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)

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
		a, b    string
		wantErr error
	}{
		{"1", "1111.123456789123456789", nil},
		{"1", "1.123456789123456789", nil},
		{"1", "2", nil},
		{"1", "3", nil},
		{"1", "4", nil},
		{"1", "5", nil},
		{"1234567890123456789", "1", nil},
		{"1234567890123456789", "2", nil},
		{"123456789012345678.9", "0.1", nil},
		{"1111111111111", "1111.123456789123456789", nil},
		{"123456789", "1.1234567890123456789", nil},
		{"2345678901234567899", "1234567890123456789.1234567890123456789", nil},
		{"-1111111111111", "1111.123456789123456789", nil},
		{"-123456789", "1.1234567890123456789", nil},
		{"-2345678901234567899", "1234567890123456789.1234567890123456789", nil},
		{"1111111111111", "-1111.123456789123456789", nil},
		{"123456789", "-1.1234567890123456789", nil},
		{"2345678901234567899", "-1234567890123456789.1234567890123456789", nil},
		{"-1111111111111", "-1111.123456789123456789", nil},
		{"-123456789", "-1.1234567890123456789", nil},
		{"-2345678901234567899", "-1234567890123456789.1234567890123456789", nil},
		{"123456789123456789.123456789", "3.123456789", nil},
		{"123456789123456789.123456789", "3", nil},
		{"9999999999999999999.9999999999999999999", "0.999", nil},
		{"9999999999999999999.9999999999999999999", "-0.999", ErrOverflow},
		{"-9999999999999999999.9999999999999999999", "0.999", ErrOverflow},
		{"0.999", "-9999999999999999999.9999999999999999999", ErrOverflow},
		{"-0.999", "9999999999999999999.9999999999999999999", ErrOverflow},
		{"9999999999999999999", "-1", ErrOverflow},
		{"-9999999999999999999", "1", ErrOverflow},
		{"9999999999999999999.9999999999999999999", "-0.0000000000000000001", ErrOverflow},
		{"-9999999999999999999.9999999999999999999", "0.0000000000000000001", ErrOverflow},
		{"9999999999999999999.99999999999999", "-0.00000000000001", ErrOverflow},
		{"-0.0000000000000000001", "9999999999999999999.9999999999999999999", ErrOverflow},
		{"0.0000000000000000001", "-9999999999999999999.9999999999999999999", ErrOverflow},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c, err := a.Sub(b)
			if tc.wantErr != nil {
				require.EqualError(t, tc.wantErr, err.Error())
				return
			}

			require.NoError(t, err)

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
		a       string
		b       uint64
		wantErr error
	}{
		{"1234567890123456789", 1, nil},
		{"1234567890123456789", 2, nil},
		{"123456789012345678.9", 1, nil},
		{"111111111111", 1111, nil},
		{"1.1234567890123456789", 123456789, nil},
		{"-123.456", 123456789, nil},
		{"-1234567890123456789.123456789", 123456789, nil},
		{"1234567890123456789.123456789", 10_000_000_000_000_000_000, nil},
		{"1234567890123456789.123456789", math.MaxUint64, ErrOverflow},
		{"-1234567890123456789.123456789", math.MaxUint64, ErrOverflow},
		{"-1234567890123456789.123456789", 10_000_000_000_000_000_000, ErrOverflow},
		{"-9999999999999999999.9999999999999999999", 10_000_000_000_000_000_000, ErrOverflow},
		{"-9999999999999999999", 1, ErrOverflow},
	}

	for _, tc := range testcases {
		t.Run(tc.a, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			c, err := a.Sub64(tc.b)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)

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
		a, b    string
		wantErr error
	}{
		{"123456.1234567890123456789", "0", nil},
		{"123456.1234567890123456789", "123456.1234567890123456789", nil},
		{"123456.1234567890123456789", "-123456.1234567890123456789", nil},
		{"-123456.1234567890123456789", "123456.1234567890123456789", nil},
		{"-123456.1234567890123456789", "-123456.1234567890123456789", nil},
		{"9999999999999999999", "0.999", nil},
		{"1234567890123456789", "1", nil},
		{"1234567890123456789", "2", nil},
		{"123456789012345678.9", "0.1", nil},
		{"1111111111111", "1111.123456789123456789", nil},
		{"123456789", "1.1234567890123456789", nil},
		{"1", "1111.123456789123456789", nil},
		{"1", "1.123456789123456789", nil},
		{"1", "2", nil},
		{"1", "3", nil},
		{"1", "4", nil},
		{"1", "5", nil},
		{"123456789123456789.123456789", "3.123456789", nil},
		{"123456789123456789.123456789", "3", nil},
		{"1.123456789123456789", "1.123456789123456789", nil},
		{"1234567890123456789.1234567890123456789", "123456", ErrOverflow},
		{"1234567890123456789.1234567890123456789", "123456.1234567890123456789", ErrOverflow},
		{"1000000", "10000000000000", ErrOverflow},
		{"-1000000", "10000000000000", ErrOverflow},
		{"-1000000", "-10000000000000", ErrOverflow},
		{"1000000", "-10000000000000", ErrOverflow},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c, err := a.Mul(b)
			if tc.wantErr != nil {
				require.EqualError(t, tc.wantErr, err.Error())
				return
			}

			require.NoError(t, err)

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
		a       string
		b       uint64
		wantErr error
	}{
		{"1234567890123456789", 0, nil},
		{"0", 123456789, nil},
		{"1234567890123456789", 1, nil},
		{"1234567890123456789", 2, nil},
		{"123456789012345678.9", 1, nil},
		{"111111111111", 1111, nil},
		{"1.1234567890123456789", 123456789, nil},
		{"-123.456", 123456789, nil},
		{"0.1234567890123456789", 10_000_000_000_000_000_000, nil},
		{"1234567890123456789.123456789", math.MaxUint64, ErrOverflow},
		{"9999999999999999999.9999999999999999999", 10_000_000_000_000_000_000, ErrOverflow},
		{"123.9999999999999999999", 10_000_000_000_000_000_000, ErrOverflow},
		{"1000000", 10_000_000_000_000, ErrOverflow},
		{"-1000000", 10_000_000_000_000, ErrOverflow},
		{"10000000000000", 1_000_000, ErrOverflow},
		{"-10000000000000", 1_000_000, ErrOverflow},
	}

	for _, tc := range testcases {
		t.Run(tc.a, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			c, err := a.Mul64(tc.b)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)

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
		a, b    string
		wantErr error
	}{
		{"123456.1234567890123456789", "234567.1234567890123456789", nil},
		{"123456.1234567890123456789", "1", nil},
		{"-123456.1234567890123456789", "234567.1234567890123456789", nil},
		{"123456.1234567890123456789", "-234567.1234567890123456789", nil},
		{"-123456.1234567890123456789", "-234567.1234567890123456789", nil},
		{"9999999999999999999", "1.0001", nil},
		{"-9999999999999999999.9999999999999999999", "9999999999999999999", nil},
		{"1234567890123456789", "1", nil},
		{"1234567890123456789", "2", nil},
		{"123456789012345678.9", "0.1", nil},
		{"1111111111111", "1111.123456789123456789", nil},
		{"123456789", "1.1234567890123456789", nil},
		{"2345678901234567899", "1234567890123456789.1234567890123456789", nil},
		{"0.1234567890123456789", "0.04586201546101", nil},
		{"1", "1111.123456789123456789", nil},
		{"1", "1.123456789123456789", nil},
		{"1", "2", nil},
		{"1", "3", nil},
		{"1", "4", nil},
		{"1", "5", nil},
		{"123456789123456789.123456789", "3.123456789", nil},
		{"123456789123456789.123456789", "3", nil},
		{"9999999999999999999", "1234567890123456789.1234567890123456879", nil},
		{"9999999999999999999.999999999999999999", "1000000000000000000.1234567890123456789", nil},
		{"999999999999999999", "0.100000000000001", nil},
		{"123456789123456789.123456789", "0", ErrDivideByZero},
		{"1000000000000", "0.0000001", ErrOverflow},
		{"1234567890123456789.123456789123456789", "0.0000000000000000002", ErrOverflow},
		{"1234567890123456789.123456789123456789", "0.000000001", ErrOverflow},
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
		a       string
		b       uint64
		wantErr error
	}{
		{"1234567890123456789", 1, nil},
		{"1234567890123456789", 2, nil},
		{"123456789012345678.9", 1, nil},
		{"111111111111", 1111, nil},
		{"1.1234567890123456789", 123456789, nil},
		{"-123.456", 123456789, nil},
		{"1234567890123456789.123456789", 10_000_000_000_000_000_000, nil},
		{"1234567890123456789.123456789", 123456789, nil},
		{"1234567890123456789.123456789", math.MaxUint64, nil},
		{"9999999999999999999.9999999999999999999", 9999999999999999999, nil},
		{"9999999999999999999.9999999999999999999", 1, nil},
		{"0.1234567890123456789", 1, nil},
		{"0.1234567890123456789", 2, nil},
		{"9999999999999999999", 1, nil},
		{"9999999999999999999", 0, ErrDivideByZero},
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

			if (a.coef == bint{}) {
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
		a       string
		scale   uint8
		want    string
		wantErr error
	}{
		{"9999999999999999999.9999999999999999999", 3, "", ErrOverflow},
		{"-9999999999999999999.9999999999999999999", 3, "", ErrOverflow},
		{"123.456000", 0, "123", nil},
		{"123.456000", 1, "123.5", nil},
		{"123.456000", 2, "123.46", nil},
		{"123.456000", 3, "123.456", nil},
		{"123.456000", 4, "123.456", nil},
		{"123.456000", 5, "123.456", nil},
		{"123.456000", 6, "123.456", nil},
		{"123.456000", 7, "123.456", nil},
		{"-123.456000", 0, "-123", nil},
		{"-123.456000", 1, "-123.5", nil},
		{"-123.456000", 2, "-123.46", nil},
		{"-123.456000", 3, "-123.456", nil},
		{"-123.456000", 4, "-123.456", nil},
		{"-123.456000", 5, "-123.456", nil},
		{"-123.456000", 6, "-123.456", nil},
		{"-123.456000", 7, "-123.456", nil},
		{"123.1234567890987654321", 0, "123", nil},
		{"123.1234567890987654321", 1, "123.1", nil},
		{"123.1234567890987654321", 2, "123.12", nil},
		{"123.1234567890987654321", 3, "123.123", nil},
		{"123.1234567890987654321", 4, "123.1235", nil},
		{"123.1234567890987654321", 5, "123.12346", nil},
		{"123.1234567890987654321", 6, "123.123457", nil},
		{"123.1234567890987654321", 7, "123.1234568", nil},
		{"123.1234567890987654321", 8, "123.12345679", nil},
		{"123.1234567890987654321", 9, "123.123456789", nil},
		{"123.1234567890987654321", 10, "123.1234567891", nil},
		{"123.1234567890987654321", 11, "123.1234567891", nil},
		{"123.1234567890987654321", 12, "123.123456789099", nil},
		{"123.1234567890987654321", 13, "123.1234567890988", nil},
		{"123.1234567890987654321", 14, "123.12345678909877", nil},
		{"123.1234567890987654321", 15, "123.123456789098765", nil},
		{"123.1234567890987654321", 16, "123.1234567890987654", nil},
		{"123.1234567890987654321", 17, "123.12345678909876543", nil},
		{"123.1234567890987654321", 18, "123.123456789098765432", nil},
		{"123.1234567890987654321", 19, "123.1234567890987654321", nil},
		{"123.1234567890987654321", 20, "123.1234567890987654321", nil},
		{"-123.1234567890987654321", 0, "-123", nil},
		{"-123.1234567890987654321", 1, "-123.1", nil},
		{"-123.1234567890987654321", 2, "-123.12", nil},
		{"-123.1234567890987654321", 3, "-123.123", nil},
		{"-123.1234567890987654321", 4, "-123.1235", nil},
		{"-123.1234567890987654321", 5, "-123.12346", nil},
		{"-123.1234567890987654321", 6, "-123.123457", nil},
		{"-123.1234567890987654321", 7, "-123.1234568", nil},
		{"-123.1234567890987654321", 8, "-123.12345679", nil},
		{"-123.1234567890987654321", 9, "-123.123456789", nil},
		{"-123.1234567890987654321", 10, "-123.1234567891", nil},
		{"-123.1234567890987654321", 11, "-123.1234567891", nil},
		{"-123.1234567890987654321", 12, "-123.123456789099", nil},
		{"-123.1234567890987654321", 13, "-123.1234567890988", nil},
		{"-123.1234567890987654321", 14, "-123.12345678909877", nil},
		{"-123.1234567890987654321", 15, "-123.123456789098765", nil},
		{"-123.1234567890987654321", 16, "-123.1234567890987654", nil},
		{"-123.1234567890987654321", 17, "-123.12345678909876543", nil},
		{"-123.1234567890987654321", 18, "-123.123456789098765432", nil},
		{"-123.1234567890987654321", 19, "-123.1234567890987654321", nil},
		{"-123.1234567890987654321", 20, "-123.1234567890987654321", nil},
		{"123.12354", 3, "123.124", nil},
		{"-123.12354", 3, "-123.124", nil},
		{"123.12454", 3, "123.125", nil},
		{"-123.12454", 3, "-123.125", nil},
		{"123.1235", 3, "123.124", nil},
		{"-123.1235", 3, "-123.124", nil},
		{"123.1245", 3, "123.124", nil},
		{"-123.1245", 3, "-123.124", nil},
		{"1.12345", 4, "1.1234", nil},
		{"1.12335", 4, "1.1234", nil},
		{"1.5", 0, "2", nil},
		{"-1.5", 0, "-2", nil},
		{"2.5", 0, "2", nil},
		{"-2.5", 0, "-2", nil},
		{"1", 0, "1", nil},
		{"-1", 0, "-1", nil},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s.round(%d)", tc.a, tc.scale), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			a, err = a.RoundBank(tc.scale)
			if tc.wantErr != nil {
				require.EqualError(t, tc.wantErr, err.Error())
				return
			}

			require.NoError(t, err)
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
		a       string
		scale   uint8
		want    string
		wantErr error
	}{
		{"9999999999999999999.9999999999999999999", 3, "", ErrOverflow},
		{"-9999999999999999999.9999999999999999999", 3, "", ErrOverflow},
		{"123.456000", 0, "123", nil},
		{"123.456000", 1, "123.5", nil},
		{"123.456000", 2, "123.46", nil},
		{"123.456000", 3, "123.456", nil},
		{"123.456000", 4, "123.456", nil},
		{"123.456000", 5, "123.456", nil},
		{"123.456000", 6, "123.456", nil},
		{"123.456000", 7, "123.456", nil},
		{"-123.456000", 0, "-123", nil},
		{"-123.456000", 1, "-123.5", nil},
		{"-123.456000", 2, "-123.46", nil},
		{"-123.456000", 3, "-123.456", nil},
		{"-123.456000", 4, "-123.456", nil},
		{"-123.456000", 5, "-123.456", nil},
		{"-123.456000", 6, "-123.456", nil},
		{"-123.456000", 7, "-123.456", nil},
		{"123.1234567890987654321", 0, "123", nil},
		{"123.1234567890987654321", 1, "123.1", nil},
		{"123.1234567890987654321", 2, "123.12", nil},
		{"123.1234567890987654321", 3, "123.123", nil},
		{"123.1234567890987654321", 4, "123.1235", nil},
		{"123.1234567890987654321", 5, "123.12346", nil},
		{"123.1234567890987654321", 6, "123.123457", nil},
		{"123.1234567890987654321", 7, "123.1234568", nil},
		{"123.1234567890987654321", 8, "123.12345679", nil},
		{"123.1234567890987654321", 9, "123.123456789", nil},
		{"123.1234567890987654321", 10, "123.1234567891", nil},
		{"123.1234567890987654321", 11, "123.1234567891", nil},
		{"123.1234567890987654321", 12, "123.123456789099", nil},
		{"123.1234567890987654321", 13, "123.1234567890988", nil},
		{"123.1234567890987654321", 14, "123.12345678909877", nil},
		{"123.1234567890987654321", 15, "123.123456789098765", nil},
		{"123.1234567890987654321", 16, "123.1234567890987654", nil},
		{"123.1234567890987654321", 17, "123.12345678909876543", nil},
		{"123.1234567890987654321", 18, "123.123456789098765432", nil},
		{"123.1234567890987654321", 19, "123.1234567890987654321", nil},
		{"123.1234567890987654321", 20, "123.1234567890987654321", nil},
		{"-123.1234567890987654321", 0, "-123", nil},
		{"-123.1234567890987654321", 1, "-123.1", nil},
		{"-123.1234567890987654321", 2, "-123.12", nil},
		{"-123.1234567890987654321", 3, "-123.123", nil},
		{"-123.1234567890987654321", 4, "-123.1235", nil},
		{"-123.1234567890987654321", 5, "-123.12346", nil},
		{"-123.1234567890987654321", 6, "-123.123457", nil},
		{"-123.1234567890987654321", 7, "-123.1234568", nil},
		{"-123.1234567890987654321", 8, "-123.12345679", nil},
		{"-123.1234567890987654321", 9, "-123.123456789", nil},
		{"-123.1234567890987654321", 10, "-123.1234567891", nil},
		{"-123.1234567890987654321", 11, "-123.1234567891", nil},
		{"-123.1234567890987654321", 12, "-123.123456789099", nil},
		{"-123.1234567890987654321", 13, "-123.1234567890988", nil},
		{"-123.1234567890987654321", 14, "-123.12345678909877", nil},
		{"-123.1234567890987654321", 15, "-123.123456789098765", nil},
		{"-123.1234567890987654321", 16, "-123.1234567890987654", nil},
		{"-123.1234567890987654321", 17, "-123.12345678909876543", nil},
		{"-123.1234567890987654321", 18, "-123.123456789098765432", nil},
		{"-123.1234567890987654321", 19, "-123.1234567890987654321", nil},
		{"-123.1234567890987654321", 20, "-123.1234567890987654321", nil},
		{"123.12354", 3, "123.124", nil},
		{"-123.12354", 3, "-123.124", nil},
		{"123.12454", 3, "123.125", nil},
		{"-123.12454", 3, "-123.125", nil},
		{"123.1235", 3, "123.124", nil},
		{"-123.1235", 3, "-123.124", nil},
		{"123.1245", 3, "123.125", nil},
		{"-123.1245", 3, "-123.125", nil},
		{"1.12345", 4, "1.1235", nil},
		{"1.12335", 4, "1.1234", nil},
		{"1.5", 0, "2", nil},
		{"-1.5", 0, "-2", nil},
		{"2.5", 0, "3", nil},
		{"-2.5", 0, "-3", nil},
		{"1", 0, "1", nil},
		{"-1", 0, "-1", nil},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s.round(%d)", tc.a, tc.scale), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			a, err = a.RoundHAZ(tc.scale)
			if tc.wantErr != nil {
				require.EqualError(t, tc.wantErr, err.Error())
				return
			}

			require.NoError(t, err)
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
		a       string
		scale   uint8
		want    string
		wantErr error
	}{
		{"9999999999999999999.9999999999999999999", 3, "", ErrOverflow},
		{"-9999999999999999999.9999999999999999999", 3, "", ErrOverflow},
		{"123.456000", 0, "123", nil},
		{"123.456000", 1, "123.5", nil},
		{"123.456000", 2, "123.46", nil},
		{"123.456000", 3, "123.456", nil},
		{"123.456000", 4, "123.456", nil},
		{"123.456000", 5, "123.456", nil},
		{"123.456000", 6, "123.456", nil},
		{"123.456000", 7, "123.456", nil},
		{"-123.456000", 0, "-123", nil},
		{"-123.456000", 1, "-123.5", nil},
		{"-123.456000", 2, "-123.46", nil},
		{"-123.456000", 3, "-123.456", nil},
		{"-123.456000", 4, "-123.456", nil},
		{"-123.456000", 5, "-123.456", nil},
		{"-123.456000", 6, "-123.456", nil},
		{"-123.456000", 7, "-123.456", nil},
		{"123.1234567890987654321", 0, "123", nil},
		{"123.1234567890987654321", 1, "123.1", nil},
		{"123.1234567890987654321", 2, "123.12", nil},
		{"123.1234567890987654321", 3, "123.123", nil},
		{"123.1234567890987654321", 4, "123.1235", nil},
		{"123.1234567890987654321", 5, "123.12346", nil},
		{"123.1234567890987654321", 6, "123.123457", nil},
		{"123.1234567890987654321", 7, "123.1234568", nil},
		{"123.1234567890987654321", 8, "123.12345679", nil},
		{"123.1234567890987654321", 9, "123.123456789", nil},
		{"123.1234567890987654321", 10, "123.1234567891", nil},
		{"123.1234567890987654321", 11, "123.1234567891", nil},
		{"123.1234567890987654321", 12, "123.123456789099", nil},
		{"123.1234567890987654321", 13, "123.1234567890988", nil},
		{"123.1234567890987654321", 14, "123.12345678909877", nil},
		{"123.1234567890987654321", 15, "123.123456789098765", nil},
		{"123.1234567890987654321", 16, "123.1234567890987654", nil},
		{"123.1234567890987654321", 17, "123.12345678909876543", nil},
		{"123.1234567890987654321", 18, "123.123456789098765432", nil},
		{"123.1234567890987654321", 19, "123.1234567890987654321", nil},
		{"123.1234567890987654321", 20, "123.1234567890987654321", nil},
		{"-123.1234567890987654321", 0, "-123", nil},
		{"-123.1234567890987654321", 1, "-123.1", nil},
		{"-123.1234567890987654321", 2, "-123.12", nil},
		{"-123.1234567890987654321", 3, "-123.123", nil},
		{"-123.1234567890987654321", 4, "-123.1235", nil},
		{"-123.1234567890987654321", 5, "-123.12346", nil},
		{"-123.1234567890987654321", 6, "-123.123457", nil},
		{"-123.1234567890987654321", 7, "-123.1234568", nil},
		{"-123.1234567890987654321", 8, "-123.12345679", nil},
		{"-123.1234567890987654321", 9, "-123.123456789", nil},
		{"-123.1234567890987654321", 10, "-123.1234567891", nil},
		{"-123.1234567890987654321", 11, "-123.1234567891", nil},
		{"-123.1234567890987654321", 12, "-123.123456789099", nil},
		{"-123.1234567890987654321", 13, "-123.1234567890988", nil},
		{"-123.1234567890987654321", 14, "-123.12345678909877", nil},
		{"-123.1234567890987654321", 15, "-123.123456789098765", nil},
		{"-123.1234567890987654321", 16, "-123.1234567890987654", nil},
		{"-123.1234567890987654321", 17, "-123.12345678909876543", nil},
		{"-123.1234567890987654321", 18, "-123.123456789098765432", nil},
		{"-123.1234567890987654321", 19, "-123.1234567890987654321", nil},
		{"-123.1234567890987654321", 20, "-123.1234567890987654321", nil},
		{"123.12354", 3, "123.124", nil},
		{"-123.12354", 3, "-123.124", nil},
		{"123.12454", 3, "123.125", nil},
		{"-123.12454", 3, "-123.125", nil},
		{"123.1235", 3, "123.123", nil},
		{"-123.1235", 3, "-123.123", nil},
		{"123.1245", 3, "123.124", nil},
		{"-123.1245", 3, "-123.124", nil},
		{"1.12345", 4, "1.1234", nil},
		{"1.12335", 4, "1.1233", nil},
		{"1.5", 0, "1", nil},
		{"-1.5", 0, "-1", nil},
		{"2.5", 0, "2", nil},
		{"-2.5", 0, "-2", nil},
		{"1", 0, "1", nil},
		{"-1", 0, "-1", nil},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s.round(%d)", tc.a, tc.scale), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			a, err = a.RoundHTZ(tc.scale)
			if tc.wantErr != nil {
				require.EqualError(t, tc.wantErr, err.Error())
				return
			}

			require.NoError(t, err)
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
		a       string
		want    string
		wantErr error
	}{
		{"9999999999999999999.9999999999999999999", "9999999999999999999", nil},
		{"-9999999999999999999.9999999999999999999", "", ErrOverflow},
		{"123.456000", "123", nil},
		{"123.456000", "123", nil},
		{"123.456000", "123", nil},
		{"123.456000", "123", nil},
		{"123.456000", "123", nil},
		{"123.456000", "123", nil},
		{"123.456000", "123", nil},
		{"123.456000", "123", nil},
		{"-123.456000", "-124", nil},
		{"-123.456000", "-124", nil},
		{"-123.456000", "-124", nil},
		{"-123.456000", "-124", nil},
		{"-123.456000", "-124", nil},
		{"-123.456000", "-124", nil},
		{"-123.456000", "-124", nil},
		{"-123.456000", "-124", nil},
		{"123.1234567890987654321", "123", nil},
		{"-123.1234567890987654321", "-124", nil},
		{"123.12354", "123", nil},
		{"-123.12354", "-124", nil},
		{"123.12454", "123", nil},
		{"-123.12454", "-124", nil},
		{"123.1235", "123", nil},
		{"-123.1235", "-124", nil},
		{"123.1245", "123", nil},
		{"-123.1245", "-124", nil},
		{"1.12345", "1", nil},
		{"1.12335", "1", nil},
		{"1.5", "1", nil},
		{"-1.5", "-2", nil},
		{"2.5", "2", nil},
		{"-2.5", "-3", nil},
		{"1", "1", nil},
		{"-1", "-1", nil},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s.floor()", tc.a), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			a, err = a.Floor()
			if tc.wantErr != nil {
				require.EqualError(t, tc.wantErr, err.Error())
				return
			}

			require.NoError(t, err)
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
		a       string
		want    string
		wantErr error
	}{
		{"9999999999999999999.9999999999999999999", "", ErrOverflow},
		{"-9999999999999999999.9999999999999999999", "-9999999999999999999", nil},
		{"123.456000", "124", nil},
		{"123.456000", "124", nil},
		{"123.456000", "124", nil},
		{"123.456000", "124", nil},
		{"123.456000", "124", nil},
		{"123.456000", "124", nil},
		{"123.456000", "124", nil},
		{"123.456000", "124", nil},
		{"-123.456000", "-123", nil},
		{"-123.456000", "-123", nil},
		{"-123.456000", "-123", nil},
		{"-123.456000", "-123", nil},
		{"-123.456000", "-123", nil},
		{"-123.456000", "-123", nil},
		{"-123.456000", "-123", nil},
		{"-123.456000", "-123", nil},
		{"123.1234567890987654321", "124", nil},
		{"-123.1234567890987654321", "-123", nil},
		{"123.12354", "124", nil},
		{"-123.12354", "-123", nil},
		{"123.12454", "124", nil},
		{"-123.12454", "-123", nil},
		{"123.1235", "124", nil},
		{"-123.1235", "-123", nil},
		{"123.1245", "124", nil},
		{"-123.1245", "-123", nil},
		{"1.12345", "2", nil},
		{"1.12335", "2", nil},
		{"1.5", "2", nil},
		{"-1.5", "-1", nil},
		{"2.5", "3", nil},
		{"-2.5", "-2", nil},
		{"1", "1", nil},
		{"-1", "-1", nil},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s.floor()", tc.a), func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			a, err = a.Ceil()
			if tc.wantErr != nil {
				require.EqualError(t, tc.wantErr, err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, a.String())

			// cross check with shopspring/decimal
			aa := decimal.RequireFromString(tc.a)
			aa = aa.Ceil()

			require.Equal(t, aa.String(), a.String())
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
	a, err := decimal.NewFromString("1234567890123456789.1234567890123456879")
	require.NoError(b, err)

	bb, err := decimal.NewFromString("1111.1789")
	require.NoError(b, err)

	b.ResetTimer()
	for range b.N {
		a.Div(bb)
	}
}
