package benchmarks

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/apd/v3"
	ed "github.com/ericlagergren/decimal"
	gv "github.com/govalues/decimal"
	ss "github.com/shopspring/decimal"

	"github.com/quagmt/udecimal"
)

// opPrecision is the number of significant digits ericlagergren and apd keep in a result.
// It covers udecimal's u128 range (up to 38 digits).
const opPrecision = 38

// parsePrecision is large enough for every benchmark input (at most 48 digits), so ericlagergren parses exactly.
const parsePrecision = 64

// apdCtx rounds results toward zero, like udecimal's truncation.
var apdCtx = func() *apd.Context {
	c := apd.BaseContext.WithPrecision(opPrecision)
	c.Rounding = apd.RoundDown
	return c
}()

// ssDivPrecision makes shopspring's Div keep 19 digits after the decimal point, like udecimal (its default is 16),
// for one benchmark only: the migration examples need the default.
func ssDivPrecision(b *testing.B) {
	old := ss.DivisionPrecision
	ss.DivisionPrecision = 19

	b.Cleanup(func() { ss.DivisionPrecision = old })
}

// exact returns an error unless plain (a library's non-exponent string form of its parsed value)
// is the same number as input.
func exact(input, plain string) error {
	want, err := udecimal.Parse(input)
	if err != nil {
		return err
	}

	got, err := udecimal.Parse(plain)
	if err != nil {
		return err
	}

	if !got.Equal(want) {
		return fmt.Errorf("parsed %s as %s", input, plain)
	}

	return nil
}

// fail reports a setup problem. In a Wide benchmark govalues may be unable to represent an input,
// which is recorded as a skip with the reason. Anything else is a bug in the benchmark and stops the run.
func fail(b *testing.B, lib string, wide bool, err error) {
	b.Helper()

	if wide && lib == "gov" {
		b.Skipf("cannot represent: %v", err)
	}

	b.Fatalf("%s: %v", lib, err)
}

func ssParse(s string) (ss.Decimal, string, error) {
	d, err := ss.NewFromString(s)
	return d, d.String(), err
}

func edNew() *ed.Big {
	z := new(ed.Big)
	z.Context.Precision = opPrecision
	z.Context.RoundingMode = ed.ToZero

	return z
}

func edParse(s string) (*ed.Big, string, error) {
	z := new(ed.Big)
	z.Context.Precision = parsePrecision

	if _, ok := z.SetString(s); !ok {
		return nil, "", fmt.Errorf("can't parse %s", s)
	}

	plain := edPlain(z)
	z.Context.Precision = opPrecision
	z.Context.RoundingMode = ed.ToZero

	return z, plain, nil
}

// edPlain formats z without an exponent.
func edPlain(z *ed.Big) string { return fmt.Sprintf("%f", z) }

func apdParse(s string) (*apd.Decimal, string, error) {
	d, _, err := apd.NewFromString(s)
	if err != nil {
		return nil, "", err
	}

	return d, d.Text('f'), nil
}

func govParse(s string) (gv.Decimal, string, error) {
	d, err := gv.Parse(s)
	return d, d.String(), err
}

func udecParse(s string) (udecimal.Decimal, string, error) {
	d, err := udecimal.Parse(s)
	return d, d.String(), err
}
