package benchmarks

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"
	ed "github.com/ericlagergren/decimal"
	gv "github.com/govalues/decimal"
	ss "github.com/shopspring/decimal"

	"github.com/quagmt/udecimal"
)

// Head-to-head inputs have at most 19 significant digits in total, so all five libraries hold them exactly.
// Wide inputs have 20 or more digits; up to 38 digits udecimal stays on its u128 path, beyond that it falls back to big.Int.

var parseCases = []string{
	"123",
	"123456.123456",
	"1234567890",
	"0.1234567890123456879",
}

var parseWideCases = []string{
	"1234567890123456789.1234567890123456879",
	"12345678901234567890123456789.123",
	"123456789123456789123456.1234567890123456",
	"111222333444555666777888999.1234567890123456789",
}

var stringWideCases = []string{
	"1234567890123456789.1234567890123456879",
	"12345.1234567890123456789",
	"123456789123456789123456.1234567890123456",
	"111222333444555666777888999.1234567890123456789",
}

type pair struct{ a, b string }

func (p pair) name(op string) string { return fmt.Sprintf("%s.%s(%s)", p.a, op, p.b) }

var arithCases = []pair{
	{"123.456", "0.123"},
	{"3", "7"},
	{"123456.123456", "999999"},
	{"123456.123456", "456781244.1324897546"},
	{"548751.15465466546", "1542.456487"},
}

var addWideCases = []pair{
	{"1234567890123456789.1234567890123456879", "1111.1789"},
	{"1234567890123456789.1234567890123456879", "1234567890123456789.1234567890123456789"},
}

var mulWideCases = []pair{
	{"1234.1234567890123456879", "1111.1789"},
	{"1234.1234567890123456879", "1111.1234567890123456789"},
}

var divWideCases = []pair{
	{"1234567890123456789.1234567890123456879", "1111.1789"},
	{"12345.1234567890123456879", "1111.1234567890123456789"},
	{"1234567890123456789.1234567890123456879", "9876543210987654321.1234567890123456789"},
	{"22773757910726981402256170801141121114", "811656739243220271.159"},
	{"12345679012345679890123456789.1234567890123456789", "999999"},
	{"1234", "12345679012345679890123456789.1234567890123456789"},
}

type powCase struct {
	a string
	n int
}

var powCases = []powCase{
	{"1.01", 10},
	{"1.01", 100},
}

// setup parses in with a library's parse function and checks the parsed value is exactly in.
func setup[T any](b *testing.B, lib string, wide bool, parse func(string) (T, string, error), in string) T {
	b.Helper()

	v, plain, err := parse(in)
	if err == nil {
		err = exact(in, plain)
	}

	if err != nil {
		fail(b, lib, wide, err)
	}

	return v
}

// check stops the benchmark if a library's text output is not exactly in.
func check(b *testing.B, lib string, in string, out []byte, err error) {
	b.Helper()

	if err == nil {
		err = exact(in, strings.Trim(string(out), `"`))
	}

	if err != nil {
		b.Fatalf("%s: %v", lib, err)
	}
}

func BenchmarkParse(b *testing.B)     { benchParse(b, parseCases, false) }
func BenchmarkParseWide(b *testing.B) { benchParse(b, parseWideCases, true) }

func benchParse(b *testing.B, cases []string, wide bool) {
	for _, tc := range cases {
		b.Run("ss/"+tc, func(b *testing.B) {
			setup(b, "ss", wide, ssParse, tc)
			b.ResetTimer()
			for range b.N {
				_, _ = ss.NewFromString(tc)
			}
		})

		b.Run("eric/"+tc, func(b *testing.B) {
			setup(b, "eric", wide, edParse, tc)
			z := new(ed.Big)
			z.Context.Precision = parsePrecision
			b.ResetTimer()
			for range b.N {
				_, _ = z.SetString(tc)
			}
		})

		b.Run("apd/"+tc, func(b *testing.B) {
			setup(b, "apd", wide, apdParse, tc)
			z := new(apd.Decimal)
			b.ResetTimer()
			for range b.N {
				_, _, _ = z.SetString(tc)
			}
		})

		b.Run("gov/"+tc, func(b *testing.B) {
			setup(b, "gov", wide, govParse, tc)
			b.ResetTimer()
			for range b.N {
				_, _ = gv.Parse(tc)
			}
		})

		b.Run("udec/"+tc, func(b *testing.B) {
			setup(b, "udec", wide, udecParse, tc)
			b.ResetTimer()
			for range b.N {
				_, _ = udecimal.Parse(tc)
			}
		})
	}
}

func BenchmarkString(b *testing.B)     { benchString(b, parseCases, false) }
func BenchmarkStringWide(b *testing.B) { benchString(b, stringWideCases, true) }

func benchString(b *testing.B, cases []string, wide bool) {
	for _, tc := range cases {
		b.Run("ss/"+tc, func(b *testing.B) {
			d := setup(b, "ss", wide, ssParse, tc)
			b.ResetTimer()
			for range b.N {
				_ = d.String()
			}
		})

		b.Run("eric/"+tc, func(b *testing.B) {
			d := setup(b, "eric", wide, edParse, tc)
			b.ResetTimer()
			for range b.N {
				_ = d.String()
			}
		})

		b.Run("apd/"+tc, func(b *testing.B) {
			d := setup(b, "apd", wide, apdParse, tc)
			b.ResetTimer()
			for range b.N {
				_ = d.String()
			}
		})

		b.Run("gov/"+tc, func(b *testing.B) {
			d := setup(b, "gov", wide, govParse, tc)
			b.ResetTimer()
			for range b.N {
				_ = d.String()
			}
		})

		b.Run("udec/"+tc, func(b *testing.B) {
			d := setup(b, "udec", wide, udecParse, tc)
			b.ResetTimer()
			for range b.N {
				_ = d.String()
			}
		})
	}
}

func BenchmarkAdd(b *testing.B)     { benchAdd(b, arithCases, false) }
func BenchmarkAddWide(b *testing.B) { benchAdd(b, addWideCases, true) }

func benchAdd(b *testing.B, cases []pair, wide bool) {
	for _, tc := range cases {
		name := tc.name("Add")

		b.Run("ss/"+name, func(b *testing.B) {
			x, y := setup(b, "ss", wide, ssParse, tc.a), setup(b, "ss", wide, ssParse, tc.b)
			b.ResetTimer()
			for range b.N {
				_ = x.Add(y)
			}
		})

		b.Run("eric/"+name, func(b *testing.B) {
			x, y := setup(b, "eric", wide, edParse, tc.a), setup(b, "eric", wide, edParse, tc.b)
			z := edNew()
			b.ResetTimer()
			for range b.N {
				_ = z.Add(x, y)
			}
		})

		b.Run("apd/"+name, func(b *testing.B) {
			x, y := setup(b, "apd", wide, apdParse, tc.a), setup(b, "apd", wide, apdParse, tc.b)
			z := new(apd.Decimal)
			if _, err := apdCtx.Add(z, x, y); err != nil {
				fail(b, "apd", wide, err)
			}
			b.ResetTimer()
			for range b.N {
				_, _ = apdCtx.Add(z, x, y)
			}
		})

		b.Run("gov/"+name, func(b *testing.B) {
			x, y := setup(b, "gov", wide, govParse, tc.a), setup(b, "gov", wide, govParse, tc.b)
			if _, err := x.Add(y); err != nil {
				fail(b, "gov", wide, err)
			}
			b.ResetTimer()
			for range b.N {
				_, _ = x.Add(y)
			}
		})

		b.Run("udec/"+name, func(b *testing.B) {
			x, y := setup(b, "udec", wide, udecParse, tc.a), setup(b, "udec", wide, udecParse, tc.b)
			b.ResetTimer()
			for range b.N {
				_ = x.Add(y)
			}
		})
	}
}

func BenchmarkSub(b *testing.B)     { benchSub(b, arithCases, false) }
func BenchmarkSubWide(b *testing.B) { benchSub(b, addWideCases, true) }

func benchSub(b *testing.B, cases []pair, wide bool) {
	for _, tc := range cases {
		name := tc.name("Sub")

		b.Run("ss/"+name, func(b *testing.B) {
			x, y := setup(b, "ss", wide, ssParse, tc.a), setup(b, "ss", wide, ssParse, tc.b)
			b.ResetTimer()
			for range b.N {
				_ = x.Sub(y)
			}
		})

		b.Run("eric/"+name, func(b *testing.B) {
			x, y := setup(b, "eric", wide, edParse, tc.a), setup(b, "eric", wide, edParse, tc.b)
			z := edNew()
			b.ResetTimer()
			for range b.N {
				_ = z.Sub(x, y)
			}
		})

		b.Run("apd/"+name, func(b *testing.B) {
			x, y := setup(b, "apd", wide, apdParse, tc.a), setup(b, "apd", wide, apdParse, tc.b)
			z := new(apd.Decimal)
			if _, err := apdCtx.Sub(z, x, y); err != nil {
				fail(b, "apd", wide, err)
			}
			b.ResetTimer()
			for range b.N {
				_, _ = apdCtx.Sub(z, x, y)
			}
		})

		b.Run("gov/"+name, func(b *testing.B) {
			x, y := setup(b, "gov", wide, govParse, tc.a), setup(b, "gov", wide, govParse, tc.b)
			if _, err := x.Sub(y); err != nil {
				fail(b, "gov", wide, err)
			}
			b.ResetTimer()
			for range b.N {
				_, _ = x.Sub(y)
			}
		})

		b.Run("udec/"+name, func(b *testing.B) {
			x, y := setup(b, "udec", wide, udecParse, tc.a), setup(b, "udec", wide, udecParse, tc.b)
			b.ResetTimer()
			for range b.N {
				_ = x.Sub(y)
			}
		})
	}
}

func BenchmarkMul(b *testing.B)     { benchMul(b, arithCases, false) }
func BenchmarkMulWide(b *testing.B) { benchMul(b, mulWideCases, true) }

func benchMul(b *testing.B, cases []pair, wide bool) {
	for _, tc := range cases {
		name := tc.name("Mul")

		b.Run("ss/"+name, func(b *testing.B) {
			x, y := setup(b, "ss", wide, ssParse, tc.a), setup(b, "ss", wide, ssParse, tc.b)
			b.ResetTimer()
			for range b.N {
				_ = x.Mul(y)
			}
		})

		b.Run("eric/"+name, func(b *testing.B) {
			x, y := setup(b, "eric", wide, edParse, tc.a), setup(b, "eric", wide, edParse, tc.b)
			z := edNew()
			b.ResetTimer()
			for range b.N {
				_ = z.Mul(x, y)
			}
		})

		b.Run("apd/"+name, func(b *testing.B) {
			x, y := setup(b, "apd", wide, apdParse, tc.a), setup(b, "apd", wide, apdParse, tc.b)
			z := new(apd.Decimal)
			if _, err := apdCtx.Mul(z, x, y); err != nil {
				fail(b, "apd", wide, err)
			}
			b.ResetTimer()
			for range b.N {
				_, _ = apdCtx.Mul(z, x, y)
			}
		})

		b.Run("gov/"+name, func(b *testing.B) {
			x, y := setup(b, "gov", wide, govParse, tc.a), setup(b, "gov", wide, govParse, tc.b)
			if _, err := x.Mul(y); err != nil {
				fail(b, "gov", wide, err)
			}
			b.ResetTimer()
			for range b.N {
				_, _ = x.Mul(y)
			}
		})

		b.Run("udec/"+name, func(b *testing.B) {
			x, y := setup(b, "udec", wide, udecParse, tc.a), setup(b, "udec", wide, udecParse, tc.b)
			b.ResetTimer()
			for range b.N {
				_ = x.Mul(y)
			}
		})
	}
}

func BenchmarkDiv(b *testing.B)     { benchDiv(b, arithCases, false) }
func BenchmarkDivWide(b *testing.B) { benchDiv(b, divWideCases, true) }

func benchDiv(b *testing.B, cases []pair, wide bool) {
	for _, tc := range cases {
		name := tc.name("Div")

		b.Run("ss/"+name, func(b *testing.B) {
			ssDivPrecision(b)
			x, y := setup(b, "ss", wide, ssParse, tc.a), setup(b, "ss", wide, ssParse, tc.b)
			b.ResetTimer()
			for range b.N {
				_ = x.Div(y)
			}
		})

		b.Run("eric/"+name, func(b *testing.B) {
			x, y := setup(b, "eric", wide, edParse, tc.a), setup(b, "eric", wide, edParse, tc.b)
			z := edNew()
			b.ResetTimer()
			for range b.N {
				_ = z.Quo(x, y)
			}
		})

		b.Run("apd/"+name, func(b *testing.B) {
			x, y := setup(b, "apd", wide, apdParse, tc.a), setup(b, "apd", wide, apdParse, tc.b)
			z := new(apd.Decimal)
			if _, err := apdCtx.Quo(z, x, y); err != nil {
				fail(b, "apd", wide, err)
			}
			b.ResetTimer()
			for range b.N {
				_, _ = apdCtx.Quo(z, x, y)
			}
		})

		b.Run("gov/"+name, func(b *testing.B) {
			x, y := setup(b, "gov", wide, govParse, tc.a), setup(b, "gov", wide, govParse, tc.b)
			if _, err := x.Quo(y); err != nil {
				fail(b, "gov", wide, err)
			}
			b.ResetTimer()
			for range b.N {
				_, _ = x.Quo(y)
			}
		})

		b.Run("udec/"+name, func(b *testing.B) {
			x, y := setup(b, "udec", wide, udecParse, tc.a), setup(b, "udec", wide, udecParse, tc.b)
			if _, err := x.Div(y); err != nil {
				fail(b, "udec", wide, err)
			}
			b.ResetTimer()
			for range b.N {
				_, _ = x.Div(y)
			}
		})
	}
}

func BenchmarkPow(b *testing.B) {
	for _, tc := range powCases {
		name := fmt.Sprintf("%s.Pow(%d)", tc.a, tc.n)

		b.Run("ss/"+name, func(b *testing.B) {
			x := setup(b, "ss", false, ssParse, tc.a)
			if _, err := x.PowInt32(int32(tc.n)); err != nil {
				fail(b, "ss", false, err)
			}
			b.ResetTimer()
			for range b.N {
				_, _ = x.PowInt32(int32(tc.n))
			}
		})

		b.Run("eric/"+name, func(b *testing.B) {
			x := setup(b, "eric", false, edParse, tc.a)
			n := ed.New(int64(tc.n), 0)
			z := edNew()
			b.ResetTimer()
			for range b.N {
				_ = z.Context.Pow(z, x, n)
			}
		})

		b.Run("apd/"+name, func(b *testing.B) {
			x := setup(b, "apd", false, apdParse, tc.a)
			n := apd.New(int64(tc.n), 0)
			z := new(apd.Decimal)
			if _, err := apdCtx.Pow(z, x, n); err != nil {
				fail(b, "apd", false, err)
			}
			b.ResetTimer()
			for range b.N {
				_, _ = apdCtx.Pow(z, x, n)
			}
		})

		b.Run("gov/"+name, func(b *testing.B) {
			x := setup(b, "gov", false, govParse, tc.a)
			if _, err := x.PowInt(tc.n); err != nil {
				fail(b, "gov", false, err)
			}
			b.ResetTimer()
			for range b.N {
				_, _ = x.PowInt(tc.n)
			}
		})

		b.Run("udec/"+name, func(b *testing.B) {
			x := setup(b, "udec", false, udecParse, tc.a)
			if _, err := x.PowInt32(int32(tc.n)); err != nil {
				fail(b, "udec", false, err)
			}
			b.ResetTimer()
			for range b.N {
				_, _ = x.PowInt32(int32(tc.n))
			}
		})
	}
}

func BenchmarkMarshalText(b *testing.B) {
	for _, tc := range parseCases {
		b.Run("ss/"+tc, func(b *testing.B) {
			d := setup(b, "ss", false, ssParse, tc)
			out, err := d.MarshalText()
			check(b, "ss", tc, out, err)
			b.ResetTimer()
			for range b.N {
				_, _ = d.MarshalText()
			}
		})

		b.Run("eric/"+tc, func(b *testing.B) {
			d := setup(b, "eric", false, edParse, tc)
			out, err := d.MarshalText()
			check(b, "eric", tc, out, err)
			b.ResetTimer()
			for range b.N {
				_, _ = d.MarshalText()
			}
		})

		b.Run("apd/"+tc, func(b *testing.B) {
			d := setup(b, "apd", false, apdParse, tc)
			out, err := d.MarshalText()
			check(b, "apd", tc, out, err)
			b.ResetTimer()
			for range b.N {
				_, _ = d.MarshalText()
			}
		})

		b.Run("gov/"+tc, func(b *testing.B) {
			d := setup(b, "gov", false, govParse, tc)
			out, err := d.MarshalText()
			check(b, "gov", tc, out, err)
			b.ResetTimer()
			for range b.N {
				_, _ = d.MarshalText()
			}
		})

		b.Run("udec/"+tc, func(b *testing.B) {
			d := setup(b, "udec", false, udecParse, tc)
			out, err := d.MarshalText()
			check(b, "udec", tc, out, err)
			b.ResetTimer()
			for range b.N {
				_, _ = d.MarshalText()
			}
		})
	}
}

func BenchmarkUnmarshalText(b *testing.B) {
	for _, tc := range parseCases {
		data := []byte(tc)

		b.Run("ss/"+tc, func(b *testing.B) {
			var d ss.Decimal
			err := d.UnmarshalText(data)
			check(b, "ss", tc, []byte(d.String()), err)
			b.ResetTimer()
			for range b.N {
				var d ss.Decimal
				_ = d.UnmarshalText(data)
			}
		})

		b.Run("eric/"+tc, func(b *testing.B) {
			z := new(ed.Big)
			z.Context.Precision = parsePrecision
			err := z.UnmarshalText(data)
			check(b, "eric", tc, []byte(edPlain(z)), err)
			b.ResetTimer()
			for range b.N {
				_ = z.UnmarshalText(data)
			}
		})

		b.Run("apd/"+tc, func(b *testing.B) {
			z := new(apd.Decimal)
			err := z.UnmarshalText(data)
			check(b, "apd", tc, []byte(z.Text('f')), err)
			b.ResetTimer()
			for range b.N {
				_ = z.UnmarshalText(data)
			}
		})

		b.Run("gov/"+tc, func(b *testing.B) {
			var d gv.Decimal
			err := d.UnmarshalText(data)
			check(b, "gov", tc, []byte(d.String()), err)
			b.ResetTimer()
			for range b.N {
				var d gv.Decimal
				_ = d.UnmarshalText(data)
			}
		})

		b.Run("udec/"+tc, func(b *testing.B) {
			var d udecimal.Decimal
			err := d.UnmarshalText(data)
			check(b, "udec", tc, []byte(d.String()), err)
			b.ResetTimer()
			for range b.N {
				var d udecimal.Decimal
				_ = d.UnmarshalText(data)
			}
		})
	}
}

// JSON goes through encoding/json for every library: apd has no JSON methods and ericlagergren has no MarshalJSON.
// Each library gets a pointer, so boxing the argument into an interface doesn't allocate.

func BenchmarkMarshalJSON(b *testing.B) {
	for _, tc := range parseCases {
		b.Run("ss/"+tc, func(b *testing.B) {
			d := setup(b, "ss", false, ssParse, tc)
			out, err := json.Marshal(&d)
			check(b, "ss", tc, out, err)
			b.ResetTimer()
			for range b.N {
				_, _ = json.Marshal(&d)
			}
		})

		b.Run("eric/"+tc, func(b *testing.B) {
			d := setup(b, "eric", false, edParse, tc)
			out, err := json.Marshal(d)
			check(b, "eric", tc, out, err)
			b.ResetTimer()
			for range b.N {
				_, _ = json.Marshal(d)
			}
		})

		b.Run("apd/"+tc, func(b *testing.B) {
			d := setup(b, "apd", false, apdParse, tc)
			out, err := json.Marshal(d)
			check(b, "apd", tc, out, err)
			b.ResetTimer()
			for range b.N {
				_, _ = json.Marshal(d)
			}
		})

		b.Run("gov/"+tc, func(b *testing.B) {
			d := setup(b, "gov", false, govParse, tc)
			out, err := json.Marshal(&d)
			check(b, "gov", tc, out, err)
			b.ResetTimer()
			for range b.N {
				_, _ = json.Marshal(&d)
			}
		})

		b.Run("udec/"+tc, func(b *testing.B) {
			d := setup(b, "udec", false, udecParse, tc)
			out, err := json.Marshal(&d)
			check(b, "udec", tc, out, err)
			b.ResetTimer()
			for range b.N {
				_, _ = json.Marshal(&d)
			}
		})
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	for _, tc := range parseCases {
		data := []byte(`"` + tc + `"`)

		b.Run("ss/"+tc, func(b *testing.B) {
			var d ss.Decimal
			err := json.Unmarshal(data, &d)
			check(b, "ss", tc, []byte(d.String()), err)
			b.ResetTimer()
			for range b.N {
				_ = json.Unmarshal(data, &d)
			}
		})

		b.Run("eric/"+tc, func(b *testing.B) {
			z := new(ed.Big)
			z.Context.Precision = parsePrecision
			err := json.Unmarshal(data, z)
			check(b, "eric", tc, []byte(edPlain(z)), err)
			b.ResetTimer()
			for range b.N {
				_ = json.Unmarshal(data, z)
			}
		})

		b.Run("apd/"+tc, func(b *testing.B) {
			z := new(apd.Decimal)
			err := json.Unmarshal(data, z)
			check(b, "apd", tc, []byte(z.Text('f')), err)
			b.ResetTimer()
			for range b.N {
				_ = json.Unmarshal(data, z)
			}
		})

		b.Run("gov/"+tc, func(b *testing.B) {
			var d gv.Decimal
			err := json.Unmarshal(data, &d)
			check(b, "gov", tc, []byte(d.String()), err)
			b.ResetTimer()
			for range b.N {
				_ = json.Unmarshal(data, &d)
			}
		})

		b.Run("udec/"+tc, func(b *testing.B) {
			var d udecimal.Decimal
			err := json.Unmarshal(data, &d)
			check(b, "udec", tc, []byte(d.String()), err)
			b.ResetTimer()
			for range b.N {
				_ = json.Unmarshal(data, &d)
			}
		})
	}
}
