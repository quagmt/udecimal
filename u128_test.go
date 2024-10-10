package udecimal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestU128Mul(t *testing.T) {
	testcases := []struct {
		u, v    u128
		want    u128
		wantErr error
	}{
		{
			u:       u128FromHiLo(10, 10),
			v:       u128FromHiLo(5, 10),
			wantErr: errOverflow,
		},
		{
			u:    u128FromHiLo(0, 10),
			v:    u128FromHiLo(5, 10),
			want: u128FromHiLo(50, 100),
		},
		{
			u:    u128FromHiLo(5, 10),
			v:    u128FromHiLo(0, 10),
			want: u128FromHiLo(50, 100),
		},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got, err := tc.u.Mul(tc.v)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.Equal(t, tc.want, got)
		})
	}
}

func TestU128Cmp64(t *testing.T) {
	testcases := []struct {
		u    u128
		v    uint64
		want int
	}{
		{
			u:    u128FromHiLo(0, 10),
			v:    10,
			want: 0,
		},
		{
			u:    u128FromHiLo(0, 10),
			v:    100,
			want: -1,
		},
		{
			u:    u128FromHiLo(10, 10),
			v:    10,
			want: 1,
		},
		{
			u:    u128FromHiLo(10, 10),
			v:    20,
			want: 1,
		},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := tc.u.Cmp64(tc.v)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestSubOverflow(t *testing.T) {
	testcases := []struct {
		u, v u128
	}{
		{
			u: u128FromHiLo(0, 10),
			v: u128FromHiLo(1, 10),
		},
		{
			u: u128FromHiLo(1, 10),
			v: u128FromHiLo(2, 10),
		},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			_, err := tc.u.Sub(tc.v)
			require.Equal(t, errOverflow, err)
		})
	}
}

func TestRightShift(t *testing.T) {
	testcases := []struct {
		u     u128
		shift uint
		want  u128
	}{
		{
			u:     u128FromHiLo(0, 10),
			shift: 1,
			want:  u128FromHiLo(0, 5),
		},
		{
			u:     u128FromHiLo(0, 10),
			shift: 2,
			want:  u128FromHiLo(0, 2),
		},
		{
			u:     u128FromHiLo(0, 10),
			shift: 3,
			want:  u128FromHiLo(0, 1),
		},
		{
			u:     u128FromHiLo(0, 10),
			shift: 4,
			want:  u128FromHiLo(0, 0),
		},
		{
			u:     u128FromHiLo(10, 0),
			shift: 65,
			want:  u128FromHiLo(0, 5),
		},
		{
			u:     u128FromHiLo(10, 0),
			shift: 66,
			want:  u128FromHiLo(0, 2),
		},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := tc.u.Rsh(tc.shift)
			require.Equal(t, tc.want, got)
		})
	}
}
