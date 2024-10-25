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
