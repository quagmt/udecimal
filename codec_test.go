package udecimal

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestStringFixed(t *testing.T) {
	testcases := []struct {
		in   string
		prec uint8
		want string
	}{
		{"0", 1, "0"},
		{"0", 0, "0"},
		{"0", 2, "0"},
		{"0", 3, "0"},
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
	}

	for _, tc := range testcases {
		t.Run(tc.in, func(t *testing.T) {
			d := MustParse(tc.in)
			require.Equal(t, tc.want, d.StringFixed(tc.prec))
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
		{"0"},
		{"1"},
		{"-1"},
		{"123456789.123456789"},
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

func TestUnmarshalNumber(t *testing.T) {
	testcases := []struct {
		in      string
		wantErr error
	}{
		{"1234567890123.1234567890123", nil},
		{"1234567890123.12345678901234567899", fmt.Errorf("precision out of range. Only support maximum 19 digits after the decimal point")},
		{`"1234567890123.1234567890123"`, nil},
	}

	for _, tc := range testcases {
		t.Run(tc.in, func(t *testing.T) {
			s := fmt.Sprintf(`{"price":%s}`, tc.in)

			var test Test
			err := json.Unmarshal([]byte(s), &test)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, strings.Trim(tc.in, `"`), test.Test.String())
		})
	}
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

func TestDynamodbMarshal(t *testing.T) {
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
		{"12345678901234567890123456789.1234567890123456789"},
		{"-12345678901234567890123456789.1234567890123456789"},
	}

	for _, tc := range testcases {
		t.Run(tc.in, func(t *testing.T) {
			d := MustParse(tc.in)

			av, err := d.MarshalDynamoDBAttributeValue()
			require.NoError(t, err)

			avN, ok := av.(*types.AttributeValueMemberN)
			require.True(t, ok)

			require.Equal(t, tc.in, avN.Value)
		})
	}
}

func TestDynamodbUnmarshal(t *testing.T) {
	testcases := []struct {
		in      types.AttributeValue
		want    Decimal
		wantErr error
	}{
		{&types.AttributeValueMemberN{Value: "0"}, MustParse("0"), nil},
		{&types.AttributeValueMemberN{Value: "1"}, MustParse("1"), nil},
		{&types.AttributeValueMemberN{Value: "-1"}, MustParse("-1"), nil},
		{&types.AttributeValueMemberN{Value: "123456789.123456789"}, MustParse("123456789.123456789"), nil},
		{&types.AttributeValueMemberN{Value: "-123456789.123456789"}, MustParse("-123456789.123456789"), nil},
		{&types.AttributeValueMemberN{Value: "0.000000001"}, MustParse("0.000000001"), nil},
		{&types.AttributeValueMemberN{Value: "-0.000000001"}, MustParse("-0.000000001"), nil},
		{&types.AttributeValueMemberN{Value: "123.123"}, MustParse("123.123"), nil},
		{&types.AttributeValueMemberN{Value: "-123.123"}, MustParse("-123.123"), nil},
		{&types.AttributeValueMemberN{Value: "12345678901234567890123456789.1234567890123456789"}, MustParse("12345678901234567890123456789.1234567890123456789"), nil},
		{&types.AttributeValueMemberN{Value: "-12345678901234567890123456789.1234567890123456789"}, MustParse("-12345678901234567890123456789.1234567890123456789"), nil},
		{&types.AttributeValueMemberS{Value: "0"}, MustParse("0"), nil},
		{&types.AttributeValueMemberS{Value: "1"}, MustParse("1"), nil},
		{&types.AttributeValueMemberS{Value: "-1"}, MustParse("-1"), nil},
		{&types.AttributeValueMemberS{Value: "123456789.123456789"}, MustParse("123456789.123456789"), nil},
		{&types.AttributeValueMemberS{Value: "-123456789.123456789"}, MustParse("-123456789.123456789"), nil},
		{&types.AttributeValueMemberS{Value: "0.000000001"}, MustParse("0.000000001"), nil},
		{&types.AttributeValueMemberS{Value: "-0.000000001"}, MustParse("-0.000000001"), nil},
		{&types.AttributeValueMemberS{Value: "123.123"}, MustParse("123.123"), nil},
		{&types.AttributeValueMemberS{Value: "-123.123"}, MustParse("-123.123"), nil},
		{&types.AttributeValueMemberS{Value: "12345678901234567890123456789.1234567890123456789"}, MustParse("12345678901234567890123456789.1234567890123456789"), nil},
		{&types.AttributeValueMemberS{Value: "-12345678901234567890123456789.1234567890123456789"}, MustParse("-12345678901234567890123456789.1234567890123456789"), nil},
		{&types.AttributeValueMemberBOOL{Value: true}, Decimal{}, fmt.Errorf("can't unmarshal %T to Decimal: %T is not supported", &types.AttributeValueMemberBOOL{}, &types.AttributeValueMemberBOOL{})},
		{&types.AttributeValueMemberN{Value: "a"}, Decimal{}, fmt.Errorf("invalid format: can't parse 'a' to Decimal")},
		{&types.AttributeValueMemberS{Value: "a"}, Decimal{}, fmt.Errorf("invalid format: can't parse 'a' to Decimal")},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			var d Decimal
			err := d.UnmarshalDynamoDBAttributeValue(tc.in)
			if tc.wantErr != nil {
				require.EqualError(t, tc.wantErr, err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, d)
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

func BenchmarkUnmarshalBinaryBigInt(b *testing.B) {
	b.StopTimer()
	data, _ := MustParse("12345678901234567890123456789.1234567890123456789").MarshalBinary()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var d Decimal
		_ = d.UnmarshalBinary(data)
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
