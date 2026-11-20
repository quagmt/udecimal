package benchmarks_test

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/quagmt/udecimal"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// The Examples below are the code on the "Coming from shopspring" page; cmd/report renders them into
// docs/snippets/examples/migration/. Each prints shopspring's result, then udecimal's.

func Example_divideByZero() {
	func() {
		defer func() { fmt.Println("shopspring:", recover()) }()
		decimal.NewFromInt(1).Div(decimal.Zero)
	}()

	_, err := udecimal.MustParse("1").Div(udecimal.Zero)
	fmt.Println("udecimal:  ", err)
	// Output:
	// shopspring: decimal division by 0
	// udecimal:   can't divide by zero
}

func Example_decimalPlaces() {
	a, err := decimal.NewFromString("0.12345678901234567890123")
	fmt.Println("shopspring:", a, err)

	b, err := udecimal.Parse("0.12345678901234567890123")
	fmt.Println("udecimal:  ", b, err)
	// Output:
	// shopspring: 0.12345678901234567890123 <nil>
	// udecimal:   0 precision out of range. Only support maximum 19 digits after the decimal point
}

func Example_division() {
	// shopspring's DivisionPrecision is 16 by default
	fmt.Println("shopspring:", decimal.NewFromInt(2).Div(decimal.NewFromInt(3)))

	q, _ := udecimal.MustParse("2").Div(udecimal.MustParse("3"))
	fmt.Println("udecimal:  ", q)
	// Output:
	// shopspring: 0.6666666666666667
	// udecimal:   0.6666666666666666666
}

func Example_multiplication() {
	fmt.Println("shopspring:", decimal.RequireFromString("0.1234567890123456789").Mul(decimal.RequireFromString("0.1")))
	fmt.Println("udecimal:  ", udecimal.MustParse("0.1234567890123456789").Mul(udecimal.MustParse("0.1")))
	// Output:
	// shopspring: 0.01234567890123456789
	// udecimal:   0.0123456789012345678
}

func Example_stringFixed() {
	fmt.Println("shopspring:", decimal.RequireFromString("1.235").StringFixed(2))
	fmt.Println("udecimal:  ", udecimal.MustParse("1.235").StringFixed(2))
	fmt.Println("udecimal:  ", udecimal.MustParse("1.235").RoundHAZ(2).StringFixed(2))
	// Output:
	// shopspring: 1.24
	// udecimal:   1.235
	// udecimal:   1.24
}

func Example_rounding() {
	a := decimal.RequireFromString("-1.225")
	fmt.Println("shopspring:", a.Round(2), a.RoundBank(2), a.Truncate(2))

	b := udecimal.MustParse("-1.225")
	fmt.Println("udecimal:  ", b.RoundHAZ(2), b.RoundBank(2), b.Trunc(2))
	// Output:
	// shopspring: -1.23 -1.22 -1.22
	// udecimal:   -1.23 -1.22 -1.22
}

func Example_inputFormat() {
	for _, s := range []string{"1e5", ".5", "1."} {
		a, aErr := decimal.NewFromString(s)
		b, bErr := udecimal.Parse(s)
		fmt.Printf("%-4s shopspring: %v %v, udecimal: %v %v\n", s, a, aErr, b, bErr)
	}
	// Output:
	// 1e5  shopspring: 100000 <nil>, udecimal: 0 invalid format: can't parse '1e5'
	// .5   shopspring: 0.5 <nil>, udecimal: 0 invalid format: can't parse '.5'
	// 1.   shopspring: 1 <nil>, udecimal: 0 invalid format: can't parse '1.'
}

func Example_powInt32() {
	_, err := decimal.Zero.PowInt32(0)
	fmt.Println("shopspring:", err)

	one, err := udecimal.Zero.PowInt32(0)
	fmt.Println("udecimal:  ", one, err)

	// a negative exponent: shopspring keeps PowPrecisionNegativeExponent (16) digits after the decimal point
	a, _ := decimal.NewFromInt(3).PowInt32(-1)
	b, _ := udecimal.MustParse("3").PowInt32(-1)
	fmt.Println("shopspring:", a)
	fmt.Println("udecimal:  ", b)
	// Output:
	// shopspring: cannot represent undefined value of 0**0
	// udecimal:   1 <nil>
	// shopspring: 0.3333333333333333
	// udecimal:   0.3333333333333333333
}

func Example_nullDecimalJSON() {
	type row struct {
		A decimal.NullDecimal
		B udecimal.NullDecimal
	}

	b, _ := json.Marshal(row{})
	fmt.Println(string(b))
	// Output:
	// {"A":null,"B":{"Decimal":"0","Valid":false}}
}

func Example_float32() {
	fmt.Println("shopspring:", decimal.NewFromFloat32(0.1))
	fmt.Println("udecimal:  ", udecimal.MustFromFloat64(float64(float32(0.1))))
	// Output:
	// shopspring: 0.1
	// udecimal:   0.10000000149011612
}

// mapping is the name table on the page. ud is the text of the page's udecimal column for that row. in and arg are
// a sample input and argument (a second operand, a number of places or an exponent); same says whether both
// libraries give the same string for them.
var mapping = []struct {
	ss, ud  string
	in, arg string
	same    bool
}{
	{"NewFromString", "Parse", "123.45", "", true},
	{"NewFromString", "Parse", "1e5", "", false}, // shopspring: 100000; udecimal: ErrInvalidFormat
	{"NewFromString", "Parse", ".5", "", false},  // shopspring: 0.5; udecimal: ErrInvalidFormat
	{"RequireFromString", "MustParse", "123.45", "", true},
	{"NewFromInt", "NewFromInt64(v, 0)", "-42", "", true},
	{"NewFromInt32", "NewFromInt64(int64(v), 0)", "-42", "", true},
	{"NewFromUint64", "NewFromUint64(v, 0)", "42", "", true},
	{"New", "NewFromInt64(value, uint8(-exp))", "12345", "-2", true},
	{"NewFromFloat", "NewFromFloat64", "0.1", "", true},
	{"Zero", "Zero", "", "", true},
	{"Add", "Add", "1.5", "2.25", true},
	{"Sub", "Sub", "1.5", "2.25", true},
	{"Mul", "Mul", "1.5", "2.25", true},
	{"Div", "Div", "1", "4", true},
	{"DivRound", "Div", "2", "3", true},
	{"Mod", "Mod", "-7", "3", true},
	{"QuoRem", "QuoRem", "7.5", "2", true},
	{"PowInt32", "PowInt32", "1.5", "3", true},
	{"PowInt32", "PowInt32", "0", "0", false},     // shopspring: an error; udecimal: 1
	{"PowInt32", "PowInt32", "3", "-1", false},    // shopspring keeps 16 digits after the decimal point, udecimal 19
	{"PowInt32", "PowInt32", "1.01", "20", false}, // shopspring exact (40 digits after the decimal point), udecimal 19
	{"Sqrt", "Sqrt", "2.25", "", true},
	{"Neg", "Neg", "1.5", "", true},
	{"Abs", "Abs", "-1.5", "", true},
	{"Cmp", "Cmp", "1", "2", true},
	{"Compare", "Cmp", "1", "2", true},
	{"Equal", "Equal", "1.50", "1.5", true},
	{"Equals", "Equal", "1.50", "1.5", true},
	{"GreaterThan", "GreaterThan", "1", "2", true},
	{"GreaterThanOrEqual", "GreaterThanOrEqual", "1", "2", true},
	{"LessThan", "LessThan", "1", "2", true},
	{"LessThanOrEqual", "LessThanOrEqual", "1", "2", true},
	{"IsZero", "IsZero", "0", "", true},
	{"IsNegative", "IsNeg", "-1", "", true},
	{"IsPositive", "IsPos", "1", "", true},
	{"Sign", "Sign", "-1.5", "", true},
	{"Round", "RoundHAZ", "-1.225", "2", true},
	{"RoundBank", "RoundBank", "1.225", "2", true},
	{"RoundUp", "RoundAwayFromZero", "1.231", "2", true},
	{"RoundDown", "Trunc", "-1.239", "2", true},
	{"Truncate", "Trunc", "1.239", "2", true},
	{"Floor", "Floor", "-1.2", "", true},
	{"Ceil", "Ceil", "-1.2", "", true},
	{"String", "String", "1.20", "", true},
	{"StringFixed", "StringFixed", "1.235", "2", false},
	{"StringFixedBank", "RoundBank(p).StringFixed(p)", "1.225", "2", true},
	{"IntPart", "Int64", "-1.9", "", true},
	{"InexactFloat64", "InexactFloat64", "1.5", "", true},
	{"Exponent", "-Prec()", "1.230", "", true},
	{"Avg", "Avg", "1,2,3", "", true},
	{"Max", "Max", "1,3,2", "", true},
	{"Min", "Min", "1,3,2", "", true},
	{"MarshalJSON", "MarshalJSON", "1.5", "", true},
	{"UnmarshalJSON", "UnmarshalJSON", `"1.5"`, "", true},
	{"MarshalText", "MarshalText", "1.5", "", true},
	{"UnmarshalText", "UnmarshalText", "1.5", "", true},
	{"MarshalBinary", "MarshalBinary", "1.5", "", false},
	{"UnmarshalBinary", "UnmarshalBinary", "1.5", "", false},
	{"Scan", "Scan", "1.5", "", true},
	{"Scan", "Scan", `"1.5"`, "", false}, // a quoted string: shopspring strips the quotes; udecimal: an error
	{"Value", "Value", "1.50", "", true},
	{"NewNullDecimal", "NullDecimal{Decimal: d, Valid: true}", "1.5", "", true},
	{"NullDecimal.Scan", "Scan", "1.5", "", true},
	{"NullDecimal.Value", "Value", "1.50", "", true},
}

// noEquivalentMethods are shopspring methods udecimal has no method for; the page says what to do instead.
var noEquivalentMethods = []string{
	"Pow", "PowWithPrecision", "PowBigInt", "SqrtRound", "Root", "RootRound", "Ln", "ExpTaylor", "ExpHullAbrham",
	"Sin", "Cos", "Tan", "Atan", "Shift", "RoundCeil", "RoundFloor", "RoundCash", "StringFixedCash", "StringScaled",
	"ScientificNotationString", "Float64", "BigInt", "BigFloat", "Rat", "Coefficient", "CoefficientInt64",
	"NumDigits", "IsInteger", "Copy", "Normalize", "Format", "GobEncode", "GobDecode", "MarshalBSONValue",
	"UnmarshalBSONValue", "EncodeSpanner", "DecodeSpanner",
}

// noEquivalentNull are shopspring NullDecimal methods udecimal's NullDecimal lacks: it implements only Scan and Value.
var noEquivalentNull = []string{
	"NullDecimal.MarshalJSON", "NullDecimal.UnmarshalJSON", "NullDecimal.MarshalText", "NullDecimal.UnmarshalText",
	"NullDecimal.MarshalBSONValue", "NullDecimal.UnmarshalBSONValue", "NullDecimal.EncodeSpanner", "NullDecimal.DecodeSpanner",
}

// noEquivalentFuncs are shopspring package functions udecimal has no function for. Referencing each one proves at
// compile time that shopspring has it.
var noEquivalentFuncs = map[string]any{
	"NewFromFloat32":           decimal.NewFromFloat32,
	"NewFromFloatWithExponent": decimal.NewFromFloatWithExponent,
	"NewFromBigInt":            decimal.NewFromBigInt,
	"NewFromBigRat":            decimal.NewFromBigRat,
	"NewFromBigFloat":          decimal.NewFromBigFloat,
	"NewFromFormattedString":   decimal.NewFromFormattedString,
	"Sum":                      decimal.Sum,
	"RescalePair":              decimal.RescalePair,
}

func split(in string) (first string, rest []string) {
	parts := strings.Split(in, ",")
	return parts[0], parts[1:]
}

// ssCall runs the shopspring side of a mapping row and returns its result as a string.
func ssCall(name, in, arg string) (string, error) {
	p := decimal.RequireFromString
	places := func() int32 { n, _ := strconv.Atoi(arg); return int32(n) }

	switch name {
	case "NewFromString":
		d, err := decimal.NewFromString(in)
		return d.String(), err
	case "RequireFromString":
		return p(in).String(), nil
	case "NewFromInt":
		v, _ := strconv.ParseInt(in, 10, 64)
		return decimal.NewFromInt(v).String(), nil
	case "NewFromInt32":
		v, _ := strconv.ParseInt(in, 10, 32)
		return decimal.NewFromInt32(int32(v)).String(), nil
	case "NewFromUint64":
		v, _ := strconv.ParseUint(in, 10, 64)
		return decimal.NewFromUint64(v).String(), nil
	case "New":
		v, _ := strconv.ParseInt(in, 10, 64)
		return decimal.New(v, places()).String(), nil
	case "NewFromFloat":
		f, _ := strconv.ParseFloat(in, 64)
		return decimal.NewFromFloat(f).String(), nil
	case "Zero":
		return decimal.Zero.String(), nil
	case "Add":
		return p(in).Add(p(arg)).String(), nil
	case "Sub":
		return p(in).Sub(p(arg)).String(), nil
	case "Mul":
		return p(in).Mul(p(arg)).String(), nil
	case "Div":
		return p(in).Div(p(arg)).String(), nil
	case "DivRound":
		return p(in).DivRound(p(arg), 2).String(), nil
	case "Mod":
		return p(in).Mod(p(arg)).String(), nil
	case "QuoRem":
		q, r := p(in).QuoRem(p(arg), 0)
		return q.String() + " " + r.String(), nil
	case "PowInt32":
		d, err := p(in).PowInt32(places())
		return d.String(), err
	case "Sqrt":
		d, err := p(in).Sqrt()
		return d.String(), err
	case "Neg":
		return p(in).Neg().String(), nil
	case "Abs":
		return p(in).Abs().String(), nil
	case "Cmp":
		return fmt.Sprint(p(in).Cmp(p(arg))), nil
	case "Compare":
		return fmt.Sprint(p(in).Compare(p(arg))), nil
	case "Equal":
		return fmt.Sprint(p(in).Equal(p(arg))), nil
	case "Equals":
		return fmt.Sprint(p(in).Equals(p(arg))), nil
	case "GreaterThan":
		return fmt.Sprint(p(in).GreaterThan(p(arg))), nil
	case "GreaterThanOrEqual":
		return fmt.Sprint(p(in).GreaterThanOrEqual(p(arg))), nil
	case "LessThan":
		return fmt.Sprint(p(in).LessThan(p(arg))), nil
	case "LessThanOrEqual":
		return fmt.Sprint(p(in).LessThanOrEqual(p(arg))), nil
	case "IsZero":
		return fmt.Sprint(p(in).IsZero()), nil
	case "IsNegative":
		return fmt.Sprint(p(in).IsNegative()), nil
	case "IsPositive":
		return fmt.Sprint(p(in).IsPositive()), nil
	case "Sign":
		return fmt.Sprint(p(in).Sign()), nil
	case "Round":
		return p(in).Round(places()).String(), nil
	case "RoundBank":
		return p(in).RoundBank(places()).String(), nil
	case "RoundUp":
		return p(in).RoundUp(places()).String(), nil
	case "RoundDown":
		return p(in).RoundDown(places()).String(), nil
	case "Truncate":
		return p(in).Truncate(places()).String(), nil
	case "Floor":
		return p(in).Floor().String(), nil
	case "Ceil":
		return p(in).Ceil().String(), nil
	case "String":
		return p(in).String(), nil
	case "StringFixed":
		return p(in).StringFixed(places()), nil
	case "StringFixedBank":
		return p(in).StringFixedBank(places()), nil
	case "IntPart":
		return fmt.Sprint(p(in).IntPart()), nil
	case "InexactFloat64":
		return fmt.Sprint(p(in).InexactFloat64()), nil
	case "Exponent":
		return fmt.Sprint(p(in).Exponent()), nil
	case "Avg", "Max", "Min":
		first, rest := split(in)
		var ds []decimal.Decimal
		for _, r := range rest {
			ds = append(ds, p(r))
		}
		agg := map[string]func(decimal.Decimal, ...decimal.Decimal) decimal.Decimal{"Avg": decimal.Avg, "Max": decimal.Max, "Min": decimal.Min}
		return agg[name](p(first), ds...).String(), nil
	case "MarshalJSON":
		b, err := p(in).MarshalJSON()
		return string(b), err
	case "UnmarshalJSON":
		var d decimal.Decimal
		err := d.UnmarshalJSON([]byte(in))
		return d.String(), err
	case "MarshalText":
		b, err := p(in).MarshalText()
		return string(b), err
	case "UnmarshalText":
		var d decimal.Decimal
		err := d.UnmarshalText([]byte(in))
		return d.String(), err
	case "MarshalBinary":
		b, err := p(in).MarshalBinary()
		return fmt.Sprint(b), err
	case "UnmarshalBinary":
		b, _ := p(in).MarshalBinary()
		var d decimal.Decimal
		err := d.UnmarshalBinary(b)
		return d.String(), err
	case "Scan":
		var d decimal.Decimal
		err := d.Scan(in)
		return d.String(), err
	case "Value":
		v, err := p(in).Value()
		return fmt.Sprint(v), err
	case "NewNullDecimal":
		v, err := decimal.NewNullDecimal(p(in)).Value()
		return fmt.Sprint(v), err
	case "NullDecimal.Scan":
		var d decimal.NullDecimal
		err := d.Scan(in)
		return fmt.Sprint(d.Decimal, d.Valid), err
	case "NullDecimal.Value":
		v, err := decimal.NullDecimal{Decimal: p(in), Valid: true}.Value()
		return fmt.Sprint(v), err
	}

	return "", fmt.Errorf("no shopspring case for %s", name)
}

// udCall runs the udecimal equivalent of the shopspring name, as the page's udecimal column describes it.
func udCall(name, in, arg string) (string, error) {
	p := udecimal.MustParse
	places := func() int { n, _ := strconv.Atoi(arg); return n }

	switch name {
	case "NewFromString":
		d, err := udecimal.Parse(in)
		return d.String(), err
	case "RequireFromString":
		return p(in).String(), nil
	case "NewFromInt", "NewFromInt32":
		v, _ := strconv.ParseInt(in, 10, 64)
		d, err := udecimal.NewFromInt64(v, 0)
		return d.String(), err
	case "NewFromUint64":
		v, _ := strconv.ParseUint(in, 10, 64)
		d, err := udecimal.NewFromUint64(v, 0)
		return d.String(), err
	case "New":
		v, _ := strconv.ParseInt(in, 10, 64)
		d, err := udecimal.NewFromInt64(v, uint8(-places()))
		return d.String(), err
	case "NewFromFloat":
		f, _ := strconv.ParseFloat(in, 64)
		d, err := udecimal.NewFromFloat64(f)
		return d.String(), err
	case "Zero":
		return udecimal.Zero.String(), nil
	case "Add":
		return p(in).Add(p(arg)).String(), nil
	case "Sub":
		return p(in).Sub(p(arg)).String(), nil
	case "Mul":
		return p(in).Mul(p(arg)).String(), nil
	case "Div":
		d, err := p(in).Div(p(arg))
		return d.String(), err
	case "DivRound":
		d, err := p(in).Div(p(arg))
		return d.RoundHAZ(2).String(), err
	case "Mod":
		d, err := p(in).Mod(p(arg))
		return d.String(), err
	case "QuoRem":
		q, r, err := p(in).QuoRem(p(arg))
		return q.String() + " " + r.String(), err
	case "PowInt32":
		d, err := p(in).PowInt32(int32(places()))
		return d.String(), err
	case "Sqrt":
		d, err := p(in).Sqrt()
		return d.String(), err
	case "Neg":
		return p(in).Neg().String(), nil
	case "Abs":
		return p(in).Abs().String(), nil
	case "Cmp", "Compare":
		return fmt.Sprint(p(in).Cmp(p(arg))), nil
	case "Equal", "Equals":
		return fmt.Sprint(p(in).Equal(p(arg))), nil
	case "GreaterThan":
		return fmt.Sprint(p(in).GreaterThan(p(arg))), nil
	case "GreaterThanOrEqual":
		return fmt.Sprint(p(in).GreaterThanOrEqual(p(arg))), nil
	case "LessThan":
		return fmt.Sprint(p(in).LessThan(p(arg))), nil
	case "LessThanOrEqual":
		return fmt.Sprint(p(in).LessThanOrEqual(p(arg))), nil
	case "IsZero":
		return fmt.Sprint(p(in).IsZero()), nil
	case "IsNegative":
		return fmt.Sprint(p(in).IsNeg()), nil
	case "IsPositive":
		return fmt.Sprint(p(in).IsPos()), nil
	case "Sign":
		return fmt.Sprint(p(in).Sign()), nil
	case "Round":
		return p(in).RoundHAZ(uint8(places())).String(), nil
	case "RoundBank":
		return p(in).RoundBank(uint8(places())).String(), nil
	case "RoundUp":
		return p(in).RoundAwayFromZero(uint8(places())).String(), nil
	case "RoundDown", "Truncate":
		return p(in).Trunc(uint8(places())).String(), nil
	case "Floor":
		return p(in).Floor().String(), nil
	case "Ceil":
		return p(in).Ceil().String(), nil
	case "String":
		return p(in).String(), nil
	case "StringFixed":
		return p(in).StringFixed(uint8(places())), nil
	case "StringFixedBank":
		return p(in).RoundBank(uint8(places())).StringFixed(uint8(places())), nil
	case "IntPart":
		v, err := p(in).Int64()
		return fmt.Sprint(v), err
	case "InexactFloat64":
		return fmt.Sprint(p(in).InexactFloat64()), nil
	case "Exponent":
		return fmt.Sprint(-p(in).Prec()), nil
	case "Avg", "Max", "Min":
		first, rest := split(in)
		var ds []udecimal.Decimal
		for _, r := range rest {
			ds = append(ds, p(r))
		}
		agg := map[string]func(udecimal.Decimal, ...udecimal.Decimal) udecimal.Decimal{"Avg": udecimal.Avg, "Max": udecimal.Max, "Min": udecimal.Min}
		return agg[name](p(first), ds...).String(), nil
	case "MarshalJSON":
		b, err := p(in).MarshalJSON()
		return string(b), err
	case "UnmarshalJSON":
		var d udecimal.Decimal
		err := d.UnmarshalJSON([]byte(in))
		return d.String(), err
	case "MarshalText":
		b, err := p(in).MarshalText()
		return string(b), err
	case "UnmarshalText":
		var d udecimal.Decimal
		err := d.UnmarshalText([]byte(in))
		return d.String(), err
	case "MarshalBinary":
		b, err := p(in).MarshalBinary()
		return fmt.Sprint(b), err
	case "UnmarshalBinary":
		// shopspring's bytes, read by udecimal: the layouts differ
		b, _ := decimal.RequireFromString(in).MarshalBinary()
		var d udecimal.Decimal
		err := d.UnmarshalBinary(b)
		return d.String(), err
	case "Scan":
		var d udecimal.Decimal
		err := d.Scan(in)
		return d.String(), err
	case "Value":
		v, err := p(in).Value()
		return fmt.Sprint(v), err
	case "NewNullDecimal":
		v, err := udecimal.NullDecimal{Decimal: p(in), Valid: true}.Value()
		return fmt.Sprint(v), err
	case "NullDecimal.Scan":
		var d udecimal.NullDecimal
		err := d.Scan(in)
		return fmt.Sprint(d.Decimal, d.Valid), err
	case "NullDecimal.Value":
		v, err := udecimal.NullDecimal{Decimal: p(in), Valid: true}.Value()
		return fmt.Sprint(v), err
	}

	return "", fmt.Errorf("no udecimal case for %s", name)
}

func TestShopspringMapping(t *testing.T) {
	for _, row := range mapping {
		ss, ssErr := ssCall(row.ss, row.in, row.arg)
		ud, udErr := udCall(row.ss, row.in, row.arg)

		if row.same {
			require.NoError(t, ssErr, row.ss)
			require.NoError(t, udErr, row.ss)
			require.Equal(t, ss, ud, row.ss)

			continue
		}

		require.False(t, ss == ud && ssErr == nil && udErr == nil, "%s: documented as different, but both give %q", row.ss, ss)
	}
}

// udecimalFuncs returns the exported package-level function and variable names of udecimal, from its source.
func udecimalFuncs(t *testing.T) map[string]bool {
	t.Helper()

	names := map[string]bool{}
	files, err := filepath.Glob("../*.go")
	require.NoError(t, err)

	fset := token.NewFileSet()
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}

		f, err := parser.ParseFile(fset, path, nil, 0)
		require.NoError(t, err)

		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Recv == nil && d.Name.IsExported() {
					names[d.Name.Name] = true
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					if v, ok := spec.(*ast.ValueSpec); ok {
						for _, n := range v.Names {
							if n.IsExported() {
								names[n.Name] = true
							}
						}
					}
				}
			}
		}
	}

	return names
}

func TestShopspringNoEquivalent(t *testing.T) {
	ssType, udType := reflect.TypeOf(decimal.Decimal{}), reflect.TypeOf(udecimal.Decimal{})
	ssPtr := reflect.PointerTo(ssType)

	for _, name := range noEquivalentMethods {
		_, onValue := ssType.MethodByName(name)
		_, onPointer := ssPtr.MethodByName(name)
		require.True(t, onValue || onPointer, "shopspring has no method %s", name)

		_, onUdValue := udType.MethodByName(name)
		_, onUdPointer := reflect.PointerTo(udType).MethodByName(name)
		require.False(t, onUdValue || onUdPointer, "udecimal has %s: move it to the mapping", name)
	}

	ssNull, udNull := reflect.PointerTo(reflect.TypeOf(decimal.NullDecimal{})), reflect.PointerTo(reflect.TypeOf(udecimal.NullDecimal{}))
	for _, name := range noEquivalentNull {
		method := strings.TrimPrefix(name, "NullDecimal.")

		_, ok := ssNull.MethodByName(method)
		require.True(t, ok, "shopspring's NullDecimal has no method %s", method)

		_, ok = udNull.MethodByName(method)
		require.False(t, ok, "udecimal's NullDecimal has %s: move it to the mapping", method)
	}

	funcs := udecimalFuncs(t)
	for name := range noEquivalentFuncs {
		require.False(t, funcs[name], "udecimal has %s: move it to the mapping", name)
	}
}

// pageRow matches a row of the mapping tables on the page, e.g. "| `Round` | `RoundHAZ` | ... |".
var pageRow = regexp.MustCompile("^\\| `([A-Za-z0-9.]+)` \\| ([^|]*) \\|")

func TestMigrationPageCovers(t *testing.T) {
	page, err := os.ReadFile("../docs/from-shopspring.md")
	require.NoError(t, err)

	onPage := map[string]string{} // shopspring name -> udecimal cell
	for _, line := range strings.Split(string(page), "\n") {
		if m := pageRow.FindStringSubmatch(line); m != nil {
			onPage[m[1]] = m[2]
		}
	}

	tested := map[string]bool{}
	for _, row := range mapping {
		tested[row.ss] = true

		cell, ok := onPage[row.ss]
		require.True(t, ok, "%s is tested but not on the page", row.ss)
		require.Contains(t, cell, "`"+row.ud+"`", "%s: the page's udecimal column differs from the test", row.ss)
	}

	for _, name := range noEquivalentMethods {
		tested[name] = true

		cell, ok := onPage[name]
		require.True(t, ok, "%s is tested but not on the page", name)
		require.Contains(t, cell, "no equivalent", name)
	}

	for _, name := range noEquivalentNull {
		tested[name] = true

		cell, ok := onPage[name]
		require.True(t, ok, "%s is tested but not on the page", name)
		require.Contains(t, cell, "no equivalent", name)
	}

	for name := range noEquivalentFuncs {
		tested[name] = true

		cell, ok := onPage[name]
		require.True(t, ok, "%s is tested but not on the page", name)
		require.Contains(t, cell, "no equivalent", name)
	}

	for name := range onPage {
		require.True(t, tested[name], "%s is on the page but has no test", name)
	}
}

// shopspringAPI returns shopspring's exported package functions, its exported variables, and its exported methods,
// named "Method" for Decimal and "NullDecimal.Method" for NullDecimal, read from the module's source.
func shopspringAPI(t *testing.T) (funcs, vars, methods []string) {
	t.Helper()

	out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/shopspring/decimal").Output()
	require.NoError(t, err)

	paths, err := filepath.Glob(filepath.Join(strings.TrimSpace(string(out)), "*.go"))
	require.NoError(t, err)

	fset := token.NewFileSet()
	for _, path := range paths {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}

		f, err := parser.ParseFile(fset, path, nil, 0)
		require.NoError(t, err)

		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if !d.Name.IsExported() {
					continue
				}

				if d.Recv == nil {
					funcs = append(funcs, d.Name.Name)
					continue
				}

				recv := d.Recv.List[0].Type
				if star, ok := recv.(*ast.StarExpr); ok {
					recv = star.X
				}

				switch recv.(*ast.Ident).Name {
				case "Decimal":
					methods = append(methods, d.Name.Name)
				case "NullDecimal":
					methods = append(methods, "NullDecimal."+d.Name.Name)
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					if v, ok := spec.(*ast.ValueSpec); ok {
						for _, n := range v.Names {
							if n.IsExported() {
								vars = append(vars, n.Name)
							}
						}
					}
				}
			}
		}
	}

	return funcs, vars, methods
}

// TestShopspringAPICovered fails when shopspring has an exported function or method the page doesn't account for,
// so a shopspring upgrade can't add a name the mapping silently skips.
func TestShopspringAPICovered(t *testing.T) {
	covered := map[string]bool{}
	for _, row := range mapping {
		covered[row.ss] = true
	}

	for _, name := range noEquivalentMethods {
		covered[name] = true
	}

	for name := range noEquivalentFuncs {
		covered[name] = true
	}

	for _, name := range noEquivalentNull {
		covered[name] = true
	}

	funcs, vars, methods := shopspringAPI(t)

	var missing []string
	for _, name := range append(funcs, methods...) {
		if !covered[name] {
			missing = append(missing, name)
		}
	}

	require.Empty(t, missing, "shopspring names neither in the mapping nor in a no-equivalent list")

	// package variables are covered by "Package settings" on the page
	page, err := os.ReadFile("../docs/from-shopspring.md")
	require.NoError(t, err)

	for _, name := range vars {
		if name == "Zero" {
			continue // a mapping row
		}

		require.Contains(t, string(page), "`"+name+"`", "the page doesn't mention shopspring's %s", name)
	}
}
