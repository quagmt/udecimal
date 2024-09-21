package udecimal

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringFixed(t *testing.T) {
	testcases := []struct {
		in    string
		scale uint8
		want  string
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
		{"123456789012345678901234567890123.123456789", 15, "123456789012345678901234567890123.123456789000000"},
		{"-123456789012345678901234567890123.123456789", 15, "-123456789012345678901234567890123.123456789000000"},
		{"-123456789012345678901234567890123.123456789", 20, "-123456789012345678901234567890123.1234567890000000000"},
	}

	for _, tc := range testcases {
		t.Run(tc.in, func(t *testing.T) {
			d := MustParse(tc.in)
			require.Equal(t, tc.want, d.StringFixed(tc.scale))
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

func BenchmarkMarshalBinary(b *testing.B) {
	b.StopTimer()
	a := MustParse("123456789.123456789")

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = a.MarshalBinary()
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
