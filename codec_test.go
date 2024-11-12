package udecimal

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringFixed(t *testing.T) {
	testcases := []struct {
		in   string
		prec uint8
		want string
	}{
		{"0", 3, "0.000"},
		{"0", 1, "0.0"},
		{"0", 0, "0"},
		{"0", 2, "0.00"},
		{"1", 3, "1.000"},
		{"-10", 10, "-10.0000000000"},
		{"123.123", 0, "123.123"},
		{"123.123", 1, "123.123"},
		{"123.123", 2, "123.123"},
		{"123.123", 3, "123.123"},
		{"123.123", 4, "123.1230"},
		{"-123.123", 4, "-123.1230"},
		{"-123.123", 5, "-123.12300"},
		{"123.123000000", 1, "123.123"},
		{"123.123000000", 2, "123.123"},
		{"123.123000000", 3, "123.123"},
		{"123.123000000", 4, "123.1230"},
		{"123.123000000", 5, "123.12300"},
		{"123.123000000", 6, "123.123000"},
		{"-123.123000000", 1, "-123.123"},
		{"-123.123000000", 2, "-123.123"},
		{"-123.123000000", 3, "-123.123"},
		{"-123.123000000", 4, "-123.1230"},
		{"-123.123000000", 5, "-123.12300"},
		{"-123.123000000", 6, "-123.123000"},
		{"123456789012345678901234567890123.123456789", 15, "123456789012345678901234567890123.123456789000000"},
		{"-123456789012345678901234567890123.123456789", 15, "-123456789012345678901234567890123.123456789000000"},
		{"-123456789012345678901234567890123.123456789", 20, "-123456789012345678901234567890123.1234567890000000000"},
		{"-34028236692093846346.3374607431768211455", 19, "-34028236692093846346.3374607431768211455"},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s.StringFixed(%d)", tc.in, tc.prec), func(t *testing.T) {
			d := MustParse(tc.in)
			require.Equal(t, tc.want, d.StringFixed(tc.prec))
		})
	}
}

func TestMarshalText(t *testing.T) {
	testcases := []struct {
		in string
	}{
		{"123456789.123456789"},
		{"0"},
		{"1"},
		{"-1"},
		{"-123456789.123456789"},
		{"0.000000001"},
		{"-0.000000001"},
		{"123.123"},
		{"-123.123"},
		{"12345678901234567890123456789.1234567890123456789"},
		{"-12345678901234567890123456789.1234567890123456789"},
	}

	for _, tc := range testcases {
		t.Run(tc.in, func(t *testing.T) {
			a := MustParse(tc.in)

			b, err := a.MarshalText()
			require.NoError(t, err)

			var c Decimal
			require.NoError(t, c.UnmarshalText(b))

			require.Equal(t, a, c)
		})
	}
}

func TestUnmarshalText(t *testing.T) {
	testcases := []struct {
		in      string
		wantErr error
	}{
		{"", ErrEmptyString},
		{" ", ErrInvalidFormat},
		{"abc", ErrInvalidFormat},
		{"1234567890123.1234567890123", nil},
		{"1234567890123.12345678901234567899", ErrPrecOutOfRange},
	}

	for _, tc := range testcases {
		t.Run(tc.in, func(t *testing.T) {
			var d Decimal
			err := d.UnmarshalText([]byte(tc.in))
			require.ErrorIs(t, err, tc.wantErr)

			if tc.wantErr == nil {
				require.Equal(t, MustParse(tc.in), d)
			}
		})
	}
}

type A struct {
	P Decimal `json:"a"`
}

func TestMarshalJSON(t *testing.T) {
	testcases := []struct {
		in string
	}{
		{"123456789.123456789"},
		{"0"},
		{"1"},
		{"-1"},
		{"-123456789.123456789"},
		{"0.000000001"},
		{"-0.000000001"},
		{"123.123"},
		{"-123.123"},
		{"12345678901234567890123456789.1234567890123456789"},
		{"-12345678901234567890123456789.1234567890123456789"},
	}

	for _, tc := range testcases {
		t.Run(tc.in, func(t *testing.T) {
			a := A{P: MustParse(tc.in)}

			b, err := json.Marshal(a)
			require.NoError(t, err)

			require.Equal(t, fmt.Sprintf(`{"a":"%s"}`, tc.in), string(b))

			// unmarshal back
			var c A
			require.NoError(t, json.Unmarshal(b, &c))

			require.Equal(t, a, c)
		})
	}
}

type Test struct {
	Test Decimal `json:"price"`
}

func TestUnmarshalJSON(t *testing.T) {
	testcases := []struct {
		in      string
		wantErr error
	}{
		{`""`, ErrEmptyString},
		{`" "`, ErrInvalidFormat},
		{`"abc"`, ErrInvalidFormat},
		{"1234567890123.1234567890123", nil},
		{"1234567890123.12345678901234567899", ErrPrecOutOfRange},
		{`"1234567890123.1234567890123"`, nil},
	}

	for _, tc := range testcases {
		t.Run(tc.in, func(t *testing.T) {
			s := fmt.Sprintf(`{"price":%s}`, tc.in)

			var test Test
			err := json.Unmarshal([]byte(s), &test)
			require.ErrorIs(t, err, tc.wantErr)

			if tc.wantErr == nil {
				require.Equal(t, strings.Trim(tc.in, `"`), test.Test.String())
			}
		})
	}
}

func TestUnmarshalJSONNull(t *testing.T) {
	var test Test
	err := json.Unmarshal([]byte(`{"price": null}`), &test)
	require.NoError(t, err)
	require.True(t, test.Test.IsZero())
}

func TestMarshalBinary(t *testing.T) {
	testcases := []struct {
		in string
	}{
		{"0"},
		{"1"},
		{"-1"},
		{"123456789.123456789"},
		{"-123456789.123456789"},
		{"0.000000001"},
		{"-0.000000001"},
		{"123.123"},
		{"-123.123"},
		{"1234567890123456789.1234567890123456789"},
		{"-1234567890123456789.1234567890123456789"},
		{"12345678901234567890123456789.1234567890123456789"},
		{"-12345678901234567890123456789.1234567890123456789"},
		{"0.0000000000000000001"},
		{"-0.0000000000000000001"},
	}

	for _, tc := range testcases {
		t.Run(tc.in, func(t *testing.T) {
			a := A{P: MustParse(tc.in)}

			var buffer bytes.Buffer
			encoder := gob.NewEncoder(&buffer)
			require.NoError(t, encoder.Encode(a))

			var c A
			decoder := gob.NewDecoder(&buffer)
			require.NoError(t, decoder.Decode(&c))

			require.Equal(t, a, c)
		})
	}
}

func TestInvalidUnmarshalBinary(t *testing.T) {
	testcases := []struct {
		name    string
		data    []byte
		wantErr error
	}{
		{"empty", []byte{}, fmt.Errorf("invalid binary data")},
		{"invalid", []byte{0x01, 0x02, 0x03}, fmt.Errorf("invalid binary data")},
		{"total len mismatched", []byte{0x01, 0x02, 0x01, 0x04, 0x05}, fmt.Errorf("invalid binary data")},
		{"len is less than 3", []byte{0x01, 0x02}, fmt.Errorf("invalid binary data")},
		{"len is less than 3, bigInt", []byte{0x11, 0x02, 0x01}, fmt.Errorf("invalid binary data")},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var d Decimal
			err := d.UnmarshalBinary(tc.data)
			require.Equal(t, tc.wantErr, err)
		})
	}
}

func TestScan(t *testing.T) {
	testcases := []struct {
		in      any
		want    Decimal
		wantErr error
	}{
		{int(0), MustParse("0"), nil},
		{int(-1234567), MustParse("-1234567"), nil},
		{int32(1), MustParse("1"), nil},
		{int64(0), MustParse("0"), nil},
		{int64(1), MustParse("1"), nil},
		{uint64(1234567890123456789), MustParse("1234567890123456789"), nil},
		{int64(-1), MustParse("-1"), nil},
		{float64(1.123), MustParse("1.123"), nil},
		{float64(-1.123), MustParse("-1.123"), nil},
		{"123.123", MustParse("123.123"), nil},
		{[]byte("123456789.123456789"), MustParse("123456789.123456789"), nil},
		{[]byte("-123456789.123456789"), MustParse("-123456789.123456789"), nil},
		{"-12345678901234567890123456789.1234567890123456789", MustParse("-12345678901234567890123456789.1234567890123456789"), nil},
		{nil, Decimal{}, fmt.Errorf("can't scan nil to Decimal")},
		{byte('a'), Decimal{}, fmt.Errorf("can't scan uint8 to Decimal: uint8 is not supported")},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			var d Decimal
			err := d.Scan(tc.in)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, d)

			// test that the value is the same after scanning
			val, err := d.Value()
			require.NoError(t, err)

			require.Equal(t, tc.want.String(), val)
		})
	}
}

func TestNullScan(t *testing.T) {
	testcases := []struct {
		in      any
		want    NullDecimal
		wantErr error
	}{
		{int(0), NullDecimal{Valid: true, Decimal: MustParse("0")}, nil},
		{int(-1234567), NullDecimal{Valid: true, Decimal: MustParse("-1234567")}, nil},
		{int32(1), NullDecimal{Valid: true, Decimal: MustParse("1")}, nil},
		{int64(0), NullDecimal{Valid: true, Decimal: MustParse("0")}, nil},
		{int64(1), NullDecimal{Valid: true, Decimal: MustParse("1")}, nil},
		{uint64(1234567890123456789), NullDecimal{Valid: true, Decimal: MustParse("1234567890123456789")}, nil},
		{int64(-1), NullDecimal{Valid: true, Decimal: MustParse("-1")}, nil},
		{float64(1.123), NullDecimal{Valid: true, Decimal: MustParse("1.123")}, nil},
		{float64(-1.123), NullDecimal{Valid: true, Decimal: MustParse("-1.123")}, nil},
		{"123.123", NullDecimal{Valid: true, Decimal: MustParse("123.123")}, nil},
		{[]byte("123456789.123456789"), NullDecimal{Valid: true, Decimal: MustParse("123456789.123456789")}, nil},
		{[]byte("-123456789.123456789"), NullDecimal{Valid: true, Decimal: MustParse("-123456789.123456789")}, nil},
		{"-12345678901234567890123456789.1234567890123456789", NullDecimal{Valid: true, Decimal: MustParse("-12345678901234567890123456789.1234567890123456789")}, nil},
		{nil, NullDecimal{Valid: false}, nil},
		{byte('a'), NullDecimal{Valid: false}, fmt.Errorf("can't scan uint8 to Decimal: uint8 is not supported")},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			var d NullDecimal
			err := d.Scan(tc.in)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, d)

			// test that the value is the same after scanning
			val, err := d.Value()
			require.NoError(t, err)

			if !d.Valid {
				require.Nil(t, val)
				return
			}

			require.Equal(t, tc.want.Decimal.String(), val)
		})
	}
}

func BenchmarkMarshalBinary(b *testing.B) {
	b.StopTimer()
	a := MustParse("123456789.123456789")

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = a.MarshalBinary()
	}
}

func BenchmarkUnmarshalBinary(b *testing.B) {
	b.StopTimer()
	data, _ := MustParse("123456789.123456789").MarshalBinary()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var d Decimal
		_ = d.UnmarshalBinary(data)
	}
}

func BenchmarkMarshalBinaryBigInt(b *testing.B) {
	b.StopTimer()
	a := MustParse("12345678901234567890123456789.1234567890123456789")

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = a.MarshalBinary()
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	b.StopTimer()
	data := []byte("123456789.123456789")

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var d Decimal
		_ = d.UnmarshalJSON(data)
	}
}

func BenchmarkString(b *testing.B) {
	b.StopTimer()
	a := MustParse("123456.123456")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = a.String()
	}
}
