package benchmarks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExact(t *testing.T) {
	tests := []struct {
		input   string
		plain   string
		wantErr bool
	}{
		{"123.45", "123.45", false},
		{"123.45", "123.4500", false},
		{"0.1234567890123456879", "0.1234567890123456879", false},
		{"1234567890123456789.1234567890123456879", "1234567890123456789", true},
		{"123.45", "123.4", true},
	}

	for _, tc := range tests {
		err := exact(tc.input, tc.plain)
		if tc.wantErr {
			require.Error(t, err, tc.input)
		} else {
			require.NoError(t, err, tc.input)
		}
	}
}

// allInputs returns every input of every benchmark: head-to-head inputs (at most 19 digits) and wide ones.
func allInputs() (head, wide []string) {
	head = append(head, parseCases...)
	for _, p := range arithCases {
		head = append(head, p.a, p.b)
	}
	for _, p := range powCases {
		head = append(head, p.a)
	}

	wide = append(wide, parseWideCases...)
	wide = append(wide, stringWideCases...)
	for _, ps := range [][]pair{addWideCases, mulWideCases, divWideCases} {
		for _, p := range ps {
			wide = append(wide, p.a, p.b)
		}
	}

	return head, wide
}

func TestParsersExact(t *testing.T) {
	head, wide := allInputs()

	for _, in := range head {
		_, plain, err := govParse(in)
		require.NoError(t, err, in)
		require.NoError(t, exact(in, plain), "gov "+in)
	}

	for _, in := range append(head, wide...) {
		_, plain, err := ssParse(in)
		require.NoError(t, err, in)
		require.NoError(t, exact(in, plain), "ss "+in)

		_, plain, err = edParse(in)
		require.NoError(t, err, in)
		require.NoError(t, exact(in, plain), "eric "+in)

		_, plain, err = apdParse(in)
		require.NoError(t, err, in)
		require.NoError(t, exact(in, plain), "apd "+in)

		_, plain, err = udecParse(in)
		require.NoError(t, err, in)
		require.NoError(t, exact(in, plain), "udec "+in)
	}
}

func TestGovaluesLosesFraction(t *testing.T) {
	// govalues v0.1.36 rounds input of more than 19 significant digits to 19, half to even, without an error.
	// 1234567890123456789.1234567890123456879 -> 1234567890123456789
	_, plain, err := govParse("1234567890123456789.1234567890123456879")
	require.NoError(t, err)
	require.Equal(t, "1234567890123456789", plain)
	require.Error(t, exact("1234567890123456789.1234567890123456879", plain))

	_, plain, err = govParse("0.1234567890123456879")
	require.NoError(t, err)
	require.NoError(t, exact("0.1234567890123456879", plain))

	_, plain, err = govParse("12345.1234567890123456789")
	require.NoError(t, err)
	require.Equal(t, "12345.12345678901235", plain)

	_, _, err = govParse("12345678901234567890")
	require.Error(t, err)
}
