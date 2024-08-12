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

			prec := int32(c.scale)
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

			prec := int32(c.scale)
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

			prec := int32(c.scale)
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

			prec := int32(c.scale)
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

			prec := int32(c.scale)
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

			prec := int32(c.scale)
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
		{"1", "1111.123456789123456789", nil},
		{"1", "1.123456789123456789", nil},
		{"1", "2", nil},
		{"1", "3", nil},
		{"1", "4", nil},
		{"1", "5", nil},
		{"123456789123456789.123456789", "3.123456789", nil},
		{"123456789123456789.123456789", "3", nil},
		{"123456789123456789.123456789", "0", ErrDivideByZero},
		{"1000000000000", "0.0000001", ErrOverflow},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			a, err := Parse(tc.a)
			require.NoError(t, err)

			b, err := Parse(tc.b)
			require.NoError(t, err)

			c, err := a.Div(b)
			if tc.wantErr != nil {
				require.EqualError(t, tc.wantErr, err.Error())
				return
			}

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
