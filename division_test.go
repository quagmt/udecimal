package udecimal

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func u256BigInt(u u256) *big.Int {
	return new(big.Int).Or(new(big.Int).Lsh(u.carry.ToBigInt(), 128), u128{hi: u.hi, lo: u.lo}.ToBigInt())
}

func TestFastQuoTrialOverflow(t *testing.T) {
	const highBit = uint64(1) << 63
	testcases := []struct {
		name    string
		u       u256
		v       u128
		wantErr bool
	}{
		{"three words", u256{carry: u128{lo: 1}}, u128{hi: 2}, false},
		{"initial equality", u256{carry: u128{lo: 1}}, u128{hi: 1}, false},
		{"four words", u256{carry: u128{hi: 1}, hi: 1}, u128{hi: highBit}, false},
		{"first trial overflow", u256{carry: u128{hi: highBit}}, u128{hi: highBit, lo: 1}, true},
		{"later trial overflow", u256{carry: u128{hi: 1}, hi: 1}, u128{hi: highBit, lo: 1}, true},
		{"trial B plus one", u256{carry: u128{hi: 1, lo: 1}, hi: 1}, u128{hi: highBit, lo: highBit + 1}, true},
		{"quotient overflow", u256{carry: u128{hi: 1}}, u128{hi: 1}, true},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var q, r u128
			var err error
			require.NotPanics(t, func() { q, r, err = tc.u.fastQuo(tc.v) })
			if tc.wantErr {
				require.ErrorIs(t, err, errOverflow)
				require.Equal(t, u128{}, q)
				require.Equal(t, u128{}, r)
				return
			}
			require.NoError(t, err)
			wantQ, wantR := new(big.Int), new(big.Int)
			wantQ.QuoRem(u256BigInt(tc.u), tc.v.ToBigInt(), wantR)
			require.Zero(t, wantQ.Cmp(q.ToBigInt()))
			require.Zero(t, wantR.Cmp(r.ToBigInt()))
		})
	}
}

func TestFastQuoNormalisation(t *testing.T) {
	for shift := range 64 {
		t.Run(fmt.Sprint(shift), func(t *testing.T) {
			v := u128{hi: uint64(1) << (63 - shift), lo: 1}
			wantQ := u128{hi: 1, lo: 3}
			n := new(big.Int).Mul(v.ToBigInt(), wantQ.ToBigInt())
			n.Add(n, big.NewInt(1))
			words := new(big.Int).Set(n)
			u := u256{lo: words.Uint64()}
			words.Rsh(words, 64)
			u.hi = words.Uint64()
			words.Rsh(words, 64)
			u.carry.lo = words.Uint64()
			words.Rsh(words, 64)
			u.carry.hi = words.Uint64()
			q, r, err := u.fastQuo(v)
			require.NoError(t, err)
			require.Equal(t, wantQ, q)
			require.Equal(t, u128{lo: 1}, r)
		})
	}
}

func TestDivTrialOverflow(t *testing.T) {
	testcases := []struct {
		name, a, b, want string
		prec             uint8
		fallback         bool
	}{
		{"reported", "340282366920938463463374607431768211455", "184467440737095516141553255926290448384", "1.8446744073709551617", 19, false},
		{"small", "3.4028236692093846347", "1.8446744073709551616", "1.8446744073709551616", 19, false},
		{"adjacent", "3.4028236692093846346", "1.8446744073709551616", "1.8446744073709551615", 19, false},
		{"later trial", "100000000000000000000", "27105054312137610850.1863200217485427857", "3.6893488147419103231", 19, true},
		{"reported precision 14", "340282366920938463463374607431768211455", "184467440737095516141553255926290448384", "1.84467440737095", 14, false},
		{"small precision 14", "34028236692.09384634633747", "184467.44073709551616", "184467.44073709551616", 14, false},
	}
	for _, tc := range testcases {
		for _, aneg := range []bool{false, true} {
			for _, bneg := range []bool{false, true} {
				for _, heap := range []bool{false, true} {
					t.Run(fmt.Sprintf("%s/%t/%t/heap=%t", tc.name, aneg, bneg, heap), func(t *testing.T) {
						previousPrec := defaultPrec
						SetDefaultPrecision(tc.prec)
						t.Cleanup(func() { SetDefaultPrecision(previousPrec) })
						a, b := MustParse(tc.a), MustParse(tc.b)
						a.neg, b.neg = aneg, bneg
						if heap {
							a.coef = bintFromBigInt(a.coef.GetBig())
							b.coef = bintFromBigInt(b.coef.GetBig())
						}
						aAlias, bAlias := a, b
						aCoef, bCoef := a.coef.GetBig(), b.coef.GetBig()
						var got Decimal
						var err error
						require.NotPanics(t, func() { got, err = a.Div(b) })
						require.NoError(t, err)
						want := tc.want
						if aneg != bneg {
							want = "-" + want
						}
						require.Equal(t, want, got.String())
						require.Equal(t, tc.prec, got.prec)
						assertOverflow(t, got, heap || tc.fallback)
						require.Equal(t, aAlias, a)
						require.Equal(t, bAlias, b)
						require.Equal(t, aCoef, aAlias.coef.GetBig())
						require.Equal(t, bCoef, bAlias.coef.GetBig())
					})
				}
			}
		}
	}
}

func TestQuoRemTrialOverflow(t *testing.T) {
	for _, aneg := range []bool{false, true} {
		for _, bneg := range []bool{false, true} {
			t.Run(fmt.Sprintf("%t/%t", aneg, bneg), func(t *testing.T) {
				a := MustParse("34028236692093846347")
				b := MustParse("1.8446744073709551616")
				a.neg, b.neg = aneg, bneg
				var q, r Decimal
				var err error
				require.NotPanics(t, func() { q, r, err = a.QuoRem(b) })
				require.NoError(t, err)
				wantQ, wantR := "18446744073709551616", "0.6625392568231788544"
				if aneg != bneg {
					wantQ = "-" + wantQ
				}
				if aneg {
					wantR = "-" + wantR
				}
				require.Equal(t, wantQ, q.String())
				require.Equal(t, wantR, r.String())
				require.Equal(t, uint8(0), q.prec)
				require.Equal(t, uint8(19), r.prec)
				assertOverflow(t, q, false)
				assertOverflow(t, r, false)
			})
		}
	}
}

var divisionBenchmarkResult Decimal

func TestDivTrialOverflowAllocations(t *testing.T) {
	for _, left := range []string{"3.4028236692093846346", "3.4028236692093846347", "4"} {
		t.Run(left, func(t *testing.T) {
			a, b := MustParse(left), MustParse("1.8446744073709551616")
			var err error
			allocs := testing.AllocsPerRun(100, func() { divisionBenchmarkResult, err = a.Div(b) })
			require.NoError(t, err)
			require.Zero(t, allocs)
		})
	}
	a := MustParse("340282366920938463463374607431768211455")
	b := MustParse("184467440737095516141553255926290448384")
	var err error
	allocs := testing.AllocsPerRun(100, func() { divisionBenchmarkResult, err = a.Div(b) })
	require.NoError(t, err)
	require.Zero(t, allocs)
}

func BenchmarkDivTrialOverflow(b *testing.B) {
	for _, tc := range []struct{ name, a, divisor string }{
		{"adjacent", "3.4028236692093846346", "1.8446744073709551616"},
		{"small", "3.4028236692093846347", "1.8446744073709551616"},
		{"reported", "340282366920938463463374607431768211455", "184467440737095516141553255926290448384"},
		{"fallback", "100000000000000000000", "27105054312137610850.1863200217485427857"},
	} {
		b.Run(tc.name, func(b *testing.B) {
			a, divisor := MustParse(tc.a), MustParse(tc.divisor)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var err error
				divisionBenchmarkResult, err = a.Div(divisor)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
