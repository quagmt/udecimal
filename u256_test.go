package udecimal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitlen(t *testing.T) {
	testcases := []struct {
		u    u256
		want int
	}{
		{
			u:    u256{hi: 0, lo: 0, carry: u128{hi: 0, lo: 0}},
			want: 0,
		},
		{
			u:    u256{hi: 0, lo: 1, carry: u128{hi: 0, lo: 0}},
			want: 1,
		},
		{
			u:    u256{hi: 0, lo: 0, carry: u128{hi: 0, lo: 1}},
			want: 129,
		},
		{
			u:    u256{hi: 0, lo: 1, carry: u128{hi: 0, lo: 1}},
			want: 129,
		},
		{
			u:    u256{hi: 0, lo: 1, carry: u128{hi: 123456789, lo: 1}},
			want: 219,
		},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := tc.u.bitLen()
			require.Equal(t, tc.want, got)
		})
	}
}

func TestRsh(t *testing.T) {
	testcases := []struct {
		u     u256
		shift uint
	}{
		{
			u:     u256{hi: 1234567890123456, lo: 1234567890123456, carry: u128{hi: 1234567890123456, lo: 1234567890123456}},
			shift: 49,
		},
		{
			u:     u256{hi: 1234567890123456, lo: 1234567890123456, carry: u128{hi: 1234567890123456, lo: 1234567890123456}},
			shift: 64,
		},
		{
			u:     u256{hi: 1234567890123456, lo: 1234567890123456, carry: u128{hi: 1234567890123456, lo: 1234567890123456}},
			shift: 113,
		},
		{
			u:     u256{hi: 1234567890123456, lo: 1234567890123456, carry: u128{hi: 1234567890123456, lo: 1234567890123456}},
			shift: 157,
		},

		{
			u:     u256{hi: 1234567890123456, lo: 1234567890123456, carry: u128{hi: 1234567890123456, lo: 1234567890123456}},
			shift: 212,
		},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			u := tc.u
			v := u.rsh(tc.shift)

			if tc.shift <= 64 {
				a := u128{hi: u.hi, lo: u.lo}.Rsh(tc.shift)
				require.Equal(t, a.lo, v.lo)

				b := u128{hi: u.carry.lo, lo: u.hi}.Rsh(tc.shift)
				require.Equal(t, b.lo, v.hi)

				c := u128{hi: u.carry.hi, lo: u.carry.lo}.Rsh(tc.shift)
				require.Equal(t, c, v.carry)
				return
			}

			if tc.shift <= 128 {
				a := u128{hi: u.carry.lo, lo: u.hi}.Rsh(tc.shift - 64)
				require.Equal(t, a.lo, v.lo)

				b := u128{hi: u.carry.hi, lo: u.carry.lo}.Rsh(tc.shift - 64)
				require.Equal(t, b.lo, v.hi)

				c := u128{hi: 0, lo: u.carry.hi}.Rsh(tc.shift - 64)
				require.Equal(t, c, v.carry)
				return
			}

			a := u128{hi: u.carry.hi, lo: u.carry.lo}.Rsh(tc.shift - 128)
			require.Equal(t, a.lo, v.lo)
			require.Equal(t, a.hi, v.hi)
			require.Equal(t, u128{}, v.carry)
		})
	}
}

// 0000000000000100011000101101010100111100100010101011101011000000.
// 0000000000000100011000101101010100111100100010101011101011000000.
// 0000000000000100011000101101010100111100100010101011101011000000.
// 0000000000000100011000101101010100111100100010101011101011000000
