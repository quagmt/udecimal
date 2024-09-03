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
			wantErr: ErrOverflow,
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

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s*%s", tc.u, tc.v), func(t *testing.T) {
			got, err := tc.u.Mul(tc.v)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr, err)
				return
			}

			require.Equal(t, tc.want, got)
		})
	}
}
