package udecimal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetDefaultParseMode(t *testing.T) {
	require.Equal(t, ParseModeError, defaultParseMode)

	SetDefaultParseMode(ParseModeTrunc)
	require.Equal(t, ParseModeTrunc, defaultParseMode)

	SetDefaultParseMode(ParseModeError)
	require.Equal(t, ParseModeError, defaultParseMode)

	// expect panic if prec is 0
	require.PanicsWithValue(t, "can't set default parse mode: invalid mode value", func() {
		SetDefaultParseMode(2)
	})
}

func TestParseModeTrunc(t *testing.T) {
	defer SetDefaultParseMode(ParseModeError)

	SetDefaultParseMode(ParseModeTrunc)

	testcases := []struct {
		input string
		want  string
	}{
		{
			input: "1.123456789012345678999",
			want:  "1.1234567890123456789",
		},
		{
			input: "-1.123456789012345678999",
			want:  "-1.1234567890123456789",
		},
		{
			input: "1.123",
			want:  "1.123",
		},
		{
			input: "-1.123",
			want:  "-1.123",
		},
		{
			input: "12324564654613213216546546132131265.123456789012345678999",
			want:  "12324564654613213216546546132131265.1234567890123456789",
		},
		{
			input: "-12324564654613213216546546132131265.123456789012345678999",
			want:  "-12324564654613213216546546132131265.1234567890123456789",
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("parse %s with mode trunc", tc.input), func(t *testing.T) {
			a := MustParse(tc.input)

			require.Equal(t, tc.want, a.String())
		})
	}
}

func TestInvalidParseMode(t *testing.T) {
	defer SetDefaultParseMode(ParseModeError)

	defaultParseMode = 2

	testcases := []string{
		"1.123456789012345678999",
		"-1.123456789012345678999",
		"12324564654613213216546546132131265.123456789012345678999",
		"-12324564654613213216546546132131265.123456789012345678999",
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("parse %s with mode trunc", tc), func(t *testing.T) {
			_, err := Parse(tc)
			require.EqualError(t, err, "invalid parse mode: 2. Make sure to use SetParseMode with a valid value")
		})
	}
}
