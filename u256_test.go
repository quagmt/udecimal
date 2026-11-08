package udecimal

import (
	"fmt"
	"math/big"
	"math/rand"
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

func u256ToBig(u u256) *big.Int {
	x := new(big.Int).SetUint64(u.carry.hi)
	for _, w := range []uint64{u.carry.lo, u.hi, u.lo} {
		x.Lsh(x, 64).Or(x, new(big.Int).SetUint64(w))
	}

	return x
}

func u256FromBig(x *big.Int) u256 {
	m := new(big.Int).SetUint64(^uint64(0))
	w := func(k uint) uint64 { return new(big.Int).And(new(big.Int).Rsh(x, k), m).Uint64() }

	return u256{carry: u128{hi: w(192), lo: w(128)}, hi: w(64), lo: w(0)}
}

// compareWithBigInt compares div256by128 with big.Int. Inputs outside its contract
// (v.hi == 0 or u.carry >= v) are skipped.
func compareWithBigInt(t *testing.T, u u256, v u128) bool {
	t.Helper()

	U, V := u256ToBig(u), v.ToBigInt()
	if v.hi == 0 || U.Cmp(new(big.Int).Lsh(V, 128)) >= 0 {
		return false
	}

	q, r := u.div256by128(v)

	wantQ, wantR := new(big.Int).QuoRem(U, V, new(big.Int))
	if q.ToBigInt().Cmp(wantQ) != 0 || r.ToBigInt().Cmp(wantR) != 0 {
		t.Fatalf("div256by128(%v, %v) = (%v, %v), want (%v, %v)", U, V, q.ToBigInt(), r.ToBigInt(), wantQ, wantR)
	}

	return true
}

// fromNormalized builds u from its normalized words [a3, a2, a1, a0], so that u << n == a.
// The low n bits of a0 are dropped.
func fromNormalized(a3, a2, a1, a0 uint64, n uint) u256 {
	x := new(big.Int)
	for _, w := range []uint64{a3, a2, a1, a0 &^ (1<<n - 1)} {
		x.Lsh(x, 64).Or(x, new(big.Int).SetUint64(w))
	}

	return u256FromBig(x.Rsh(x, n))
}

func TestDiv256by128(t *testing.T) {
	testcases := []struct {
		u            u256
		v            u128
		wantQ, wantR u128
	}{
		// issue #46: (2^128-1)*10^19 / (9999999999999999999*2^64), 2 quotient digits with
		// a2 == v.hi and [a2, a1] >= v (e.g. 896/87)
		{
			u:     u256{hi: 18446744073709551615, lo: 8446744073709551616, carry: u128FromHiLo(0, 9999999999999999999)},
			v:     u128FromHiLo(9999999999999999999, 0),
			wantQ: u128FromHiLo(1, 1),
			wantR: u128FromHiLo(8446744073709551616, 8446744073709551616),
		},
		// 2 quotient digits, a2 == v.hi and [a2, a1] >= v, n = 47
		{
			u:     u256{hi: 1618127530, lo: 0, carry: u128FromHiLo(0, 65536)},
			v:     u128FromHiLo(65536, 1618127530),
			wantQ: u128FromHiLo(1, 0),
			wantR: u128FromHiLo(0, 0),
		},
		// 1 quotient digit, u2 == v.hi, no rOverflow (e.g. 866/87)
		{
			u:     u256{hi: 6148820866244280319, lo: 11259375, carry: u128FromHiLo(0, 9223372036854788153)},
			v:     u128FromHiLo(9223372036854788153, 6148820866244280320),
			wantQ: u128FromHiLo(0, 18446744073709551615),
			wantR: u128FromHiLo(9223372036854788152, 6148820866255539695),
		},
		// 1 quotient digit, u2 == v.hi, rOverflow, n = 13
		{
			u:     u256{hi: 9353976426048520190, lo: 18444492273895867742, carry: u128FromHiLo(0, 1125899906842625)},
			v:     u128FromHiLo(1125899906842625, 9353976426048520191),
			wantQ: u128FromHiLo(0, 18446744073709551615),
			wantR: u128FromHiLo(1125899906842625, 9351724626234836317),
		},
		// u2 == v.hi, 2^64-1 is right (e.g. 500/51)
		{
			u:     u256{hi: 0, lo: 0, carry: u128FromHiLo(0, 9223372036854788153)},
			v:     u128FromHiLo(9223372036854788153, 9223372036854775808),
			wantQ: u128FromHiLo(0, 18446744073709551615),
			wantR: u128FromHiLo(12345, 9223372036854775808),
		},
		// u2 == v.hi, 2^64-1 is 1 too big (e.g. 500/56), n = 13
		{
			u:     u256{hi: 9351724626234834944, lo: 0, carry: u128FromHiLo(0, 1125899906842625)},
			v:     u128FromHiLo(1125899906842625, 9353976426048520191),
			wantQ: u128FromHiLo(0, 18446744073709551614),
			wantR: u128FromHiLo(4, 261208778387488766),
		},
		// 1st window u2 == v.hi (a3 == v.hi), no rOverflow, n = 47
		{
			u:     u256{hi: 18446744073709420544, lo: 610795520, carry: u128FromHiLo(65536, 1618127529)},
			v:     u128FromHiLo(65536, 1618127530),
			wantQ: u128FromHiLo(18446744073709551615, 18446744073709551614),
			wantR: u128FromHiLo(0, 3847050580),
		},
		// 1st window u2 == v.hi (a3 == v.hi), rOverflow; then 2nd window u2 == v.hi, 2^64-1 is 1 too big
		{
			u:     u256{hi: 4660, lo: 11259375, carry: u128FromHiLo(9223372036854788153, 18446744073709551614)},
			v:     u128FromHiLo(9223372036854788153, 18446744073709551615),
			wantQ: u128FromHiLo(18446744073709551615, 18446744073709551614),
			wantR: u128FromHiLo(29352, 11259373),
		},
		// 2nd window u2 == v.hi (e.g. 6046/87: 604 -> rem 82, next window 826), n = 13
		{
			u:     u256{hi: 888360040388362239, lo: 18444492273895867742, carry: u128FromHiLo(1, 15987778677165306890)},
			v:     u128FromHiLo(1125899906842625, 9352475214719483904),
			wantQ: u128FromHiLo(30583, 18446744073709551615),
			wantR: u128FromHiLo(1125899906842625, 9350223414905800030),
		},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			q, r := tc.u.div256by128(tc.v)
			require.Equal(t, tc.wantQ, q)
			require.Equal(t, tc.wantR, r)
		})
	}
}

func TestDiv256by128Random(t *testing.T) {
	r := rand.New(rand.NewSource(1))

	for i := 0; i < 500_000; i++ {
		v := u128{hi: r.Uint64() >> (r.Uint64() % 64), lo: r.Uint64()}
		c := u128{hi: r.Uint64() >> (r.Uint64() % 64), lo: r.Uint64()}

		if c.Cmp(v) >= 0 && v.hi != 0 {
			c.hi %= v.hi
		}

		compareWithBigInt(t, u256{carry: c, hi: r.Uint64(), lo: r.Uint64()}, v)
	}
}
