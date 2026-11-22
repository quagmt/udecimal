package udecimal_test

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/quagmt/udecimal"
)

func Example() {
	a := udecimal.MustParse("123.456")
	b := udecimal.MustParse("-12.3456")

	fmt.Println(a.Add(b))
	fmt.Println(a.Sub(b))
	fmt.Println(a.Mul(b))
	fmt.Println(a.Div(b))
	fmt.Println(a.Div(udecimal.MustParse("7")))
	// Output:
	// 111.1104
	// 135.8016
	// -1524.1383936
	// -10 <nil>
	// 17.6365714285714285714 <nil>
}

func ExampleSetDefaultParseMode() {
	fmt.Println(udecimal.Parse("0.12345678901234567890123"))

	udecimal.SetDefaultParseMode(udecimal.ParseModeTrunc)
	defer udecimal.SetDefaultParseMode(udecimal.ParseModeError)

	fmt.Println(udecimal.Parse("0.12345678901234567890123"))
	// Output:
	// 0 precision out of range: more decimal places than the default precision
	// 0.1234567890123456789 <nil>
}

func ExampleSetDefaultPrecision() {
	udecimal.SetDefaultPrecision(10)
	defer udecimal.SetDefaultPrecision(19)

	a := udecimal.MustParse("1.23")
	b := udecimal.MustParse("4.12475")
	c, _ := a.Div(b)
	fmt.Println(c)
	// Output:
	// 0.2981998909
}

func ExampleMustFromFloat64() {
	fmt.Println(udecimal.MustFromFloat64(1.234))

	// cautious: result will lose some precision when converting to decimal
	fmt.Println(udecimal.MustFromFloat64(123456789.1234567890123456789))
	// Output:
	// 1.234
	// 123456789.12345679
}

func ExampleNewFromFloat64() {
	fmt.Println(udecimal.NewFromFloat64(1.234))

	// cautious: result will lose some precision when converting to decimal
	fmt.Println(udecimal.NewFromFloat64(123456789.1234567890123456789))
	// Output:
	// 1.234 <nil>
	// 123456789.12345679 <nil>
}

func ExampleMustFromInt64() {
	fmt.Println(udecimal.MustFromInt64(123, 3))
	fmt.Println(udecimal.MustFromInt64(-12345, 2))
	// Output:
	// 0.123
	// -123.45
}

func ExampleNewFromInt64() {
	fmt.Println(udecimal.NewFromInt64(123, 3))
	fmt.Println(udecimal.NewFromInt64(12345, 20))
	// Output:
	// 0.123 <nil>
	// 0 precision out of range: more decimal places than the default precision
}

func ExampleMustFromUint64() {
	fmt.Println(udecimal.MustFromUint64(123, 3))
	fmt.Println(udecimal.MustFromUint64(12345, 2))
	// Output:
	// 0.123
	// 123.45
}

func ExampleNewFromUint64() {
	fmt.Println(udecimal.NewFromUint64(123, 3))
	fmt.Println(udecimal.NewFromUint64(12345, 2))
	fmt.Println(udecimal.NewFromUint64(12345, 20))
	// Output:
	// 0.123 <nil>
	// 123.45 <nil>
	// 0 precision out of range: more decimal places than the default precision
}

func ExampleMustParse() {
	fmt.Println(udecimal.MustParse("1234567890123456789.1234567890123456789"))
	fmt.Println(udecimal.MustParse("-1234567890123456789.1234567890123456789"))
	fmt.Println(udecimal.MustParse("-0.00007890123456789"))
	// Output:
	// 1234567890123456789.1234567890123456789
	// -1234567890123456789.1234567890123456789
	// -0.00007890123456789
}

func ExampleParse() {
	fmt.Println(udecimal.Parse("1234567890123456789.1234567890123456789"))
	fmt.Println(udecimal.Parse("-1234567890123456789.1234567890123456789"))
	fmt.Println(udecimal.Parse("-0.00007890123456789"))

	// error cases
	fmt.Println(udecimal.Parse("0.12345678901234567890123"))
	fmt.Println(udecimal.Parse(""))
	fmt.Println(udecimal.Parse("1.123.123"))

	// Output:
	// 1234567890123456789.1234567890123456789 <nil>
	// -1234567890123456789.1234567890123456789 <nil>
	// -0.00007890123456789 <nil>
	// 0 precision out of range: more decimal places than the default precision
	// 0 can't parse empty string
	// 0 invalid format: can't parse '1.123.123'
}

func ExampleNewFromHiLo() {
	fmt.Println(udecimal.NewFromHiLo(false, 1, 1, 10))
	fmt.Println(udecimal.NewFromHiLo(true, 0, 123456, 4))
	// Output:
	// 1844674407.3709551617 <nil>
	// -12.3456 <nil>
}

func ExampleDecimal_Abs() {
	fmt.Println(udecimal.MustParse("-123.45").Abs())
	// Output:
	// 123.45
}

func ExampleDecimal_Add() {
	a := udecimal.MustParse("1.23")
	b := udecimal.MustParse("4.12475")
	c := a.Add(b)
	fmt.Println(c)
	// Output:
	// 5.35475
}

func ExampleDecimal_Add64() {
	a := udecimal.MustParse("1.23")
	c := a.Add64(4)
	fmt.Println(c)
	// Output:
	// 5.23
}

func ExampleDecimal_Ceil() {
	fmt.Println(udecimal.MustParse("1.23").Ceil())
	// Output:
	// 2
}

func ExampleDecimal_Cmp() {
	fmt.Println(udecimal.MustParse("1.23").Cmp(udecimal.MustParse("4.12475")))
	fmt.Println(udecimal.MustParse("1.23").Cmp(udecimal.MustParse("1.23")))
	fmt.Println(udecimal.MustParse("1.23").Cmp(udecimal.MustParse("0.12475")))
	// Output:
	// -1
	// 0
	// 1
}

func ExampleDecimal_LessThan() {
	fmt.Println(udecimal.MustParse("1.23").LessThan(udecimal.MustParse("4.12475")))
	fmt.Println(udecimal.MustParse("1.23").LessThan(udecimal.MustParse("1.23")))
	fmt.Println(udecimal.MustParse("1.23").LessThan(udecimal.MustParse("0.12475")))
	// Output:
	// true
	// false
	// false
}

func ExampleDecimal_LessThanOrEqual() {
	fmt.Println(udecimal.MustParse("1.23").LessThanOrEqual(udecimal.MustParse("4.12475")))
	fmt.Println(udecimal.MustParse("1.23").LessThanOrEqual(udecimal.MustParse("1.23")))
	fmt.Println(udecimal.MustParse("1.23").LessThanOrEqual(udecimal.MustParse("0.12475")))
	// Output:
	// true
	// true
	// false
}

func ExampleDecimal_GreaterThan() {
	fmt.Println(udecimal.MustParse("1.23").GreaterThan(udecimal.MustParse("4.12475")))
	fmt.Println(udecimal.MustParse("1.23").GreaterThan(udecimal.MustParse("1.23")))
	fmt.Println(udecimal.MustParse("1.23").GreaterThan(udecimal.MustParse("0.12475")))
	// Output:
	// false
	// false
	// true
}

func ExampleDecimal_GreaterThanOrEqual() {
	fmt.Println(udecimal.MustParse("1.23").GreaterThanOrEqual(udecimal.MustParse("4.12475")))
	fmt.Println(udecimal.MustParse("1.23").GreaterThanOrEqual(udecimal.MustParse("1.23")))
	fmt.Println(udecimal.MustParse("1.23").GreaterThanOrEqual(udecimal.MustParse("0.12475")))
	// Output:
	// false
	// true
	// true
}

func ExampleMax() {
	fmt.Println(udecimal.Max(udecimal.MustParse("1.23"), udecimal.MustParse("4.12475")))
	// Output:
	// 4.12475
}

func ExampleMin() {
	fmt.Println(udecimal.Min(udecimal.MustParse("1.23"), udecimal.MustParse("4.12475")))
	// Output:
	// 1.23
}

func ExampleDecimal_Equal() {
	fmt.Println(udecimal.MustParse("1.123").Equal(udecimal.MustParse("-1.123")))
	fmt.Println(udecimal.MustParse("1.123").Equal(udecimal.MustParse("1.1230000")))
	// Output:
	// false
	// true
}

func ExampleDecimal_Div() {
	fmt.Println(udecimal.MustParse("1.23").Div(udecimal.MustParse("4.12475")))
	fmt.Println(udecimal.MustParse("1.23").Div(udecimal.MustParse("0")))
	// Output:
	// 0.2981998909024789381 <nil>
	// 0 can't divide by zero
}

func ExampleDecimal_MustDiv() {
	fmt.Println(udecimal.MustParse("1.23").MustDiv(udecimal.MustParse("4.12475")))
	// Output:
	// 0.2981998909024789381
}

func ExampleDecimal_Div64() {
	fmt.Println(udecimal.MustParse("1.23").Div64(4))
	fmt.Println(udecimal.MustParse("1.23").Div64(0))
	// Output:
	// 0.3075 <nil>
	// 0 can't divide by zero
}

func ExampleDecimal_QuoRem() {
	fmt.Println(udecimal.MustParse("1.23").QuoRem(udecimal.MustParse("0.5")))
	fmt.Println(udecimal.MustParse("1.23").QuoRem(udecimal.MustParse("0")))
	// Output:
	// 2 0.23 <nil>
	// 0 0 can't divide by zero
}

func ExampleDecimal_MustQuoRem() {
	fmt.Println(udecimal.MustParse("1.23").MustQuoRem(udecimal.MustParse("0.5")))
	// Output:
	// 2 0.23
}

func ExampleDecimal_Mod() {
	fmt.Println(udecimal.MustParse("1.23").Mod(udecimal.MustParse("0.5")))
	fmt.Println(udecimal.MustParse("1.23").Mod(udecimal.MustParse("0")))
	// Output:
	// 0.23 <nil>
	// 0 can't divide by zero
}

func ExampleDecimal_MustMod() {
	fmt.Println(udecimal.MustParse("1.23").MustMod(udecimal.MustParse("0.5")))
	// Output:
	// 0.23
}

func ExampleDecimal_Sub() {
	a := udecimal.MustParse("1.23")
	b := udecimal.MustParse("4.12475")
	c := a.Sub(b)
	fmt.Println(c)
	// Output:
	// -2.89475
}

func ExampleDecimal_Sub64() {
	a := udecimal.MustParse("1.23")
	c := a.Sub64(4)
	fmt.Println(c)
	// Output:
	// -2.77
}

func ExampleDecimal_Mul() {
	a := udecimal.MustParse("1.23")
	b := udecimal.MustParse("4.12475")
	c := a.Mul(b)
	fmt.Println(c)
	// Output:
	// 5.0734425
}

func ExampleDecimal_Mul64() {
	a := udecimal.MustParse("1.23")
	c := a.Mul64(4)
	fmt.Println(c)
	// Output:
	// 4.92
}

func ExampleDecimal_Floor() {
	fmt.Println(udecimal.MustParse("1.23").Floor())
	fmt.Println(udecimal.MustParse("-1.23").Floor())
	// Output:
	// 1
	// -2
}

func ExampleDecimal_Int64() {
	fmt.Println(udecimal.MustParse("1.23").Int64())
	fmt.Println(udecimal.MustParse("1234567890123456789.1234567890123456789").Int64())
	fmt.Println(udecimal.MustParse("-1234567890123456789.1234567890123456789").Int64())
	fmt.Println(udecimal.MustParse("0.123").Int64())
	fmt.Println(udecimal.MustParse("-0.123").Int64())
	fmt.Println(udecimal.MustParse("12345678901234567890123456.123456").Int64())
	// Output:
	// 1 <nil>
	// 1234567890123456789 <nil>
	// -1234567890123456789 <nil>
	// 0 <nil>
	// 0 <nil>
	// 0 integer part is too large to fit in int64
}

func ExampleDecimal_InexactFloat64() {
	fmt.Println(udecimal.MustParse("1.23").InexactFloat64())
	fmt.Println(udecimal.MustParse("123456789.123456789").InexactFloat64())
	// Output:
	// 1.23
	// 1.2345678912345679e+08
}

func ExampleDecimal_ToHiLo() {
	fmt.Println(udecimal.MustParse("1.23").ToHiLo())
	fmt.Println(udecimal.MustParse("1234567890123456789.1234567890123456789").ToHiLo())
	fmt.Println(udecimal.MustParse("-1234567890123456789.1234567890123456789").ToHiLo())
	fmt.Println(udecimal.MustParse("0.123").ToHiLo())
	fmt.Println(udecimal.MustParse("-0.123").ToHiLo())
	fmt.Println(udecimal.MustParse("12345678901234567890123456.123456").ToHiLo())
	fmt.Println(udecimal.MustParse("0").ToHiLo())
	fmt.Println(udecimal.MustParse("18446744073709551615").ToHiLo())
	fmt.Println(udecimal.MustParse("18446744073709551617").ToHiLo())
	fmt.Println(udecimal.MustParse("184467440737095516.15").ToHiLo())
	fmt.Println(udecimal.MustParse("184467440737095516.16").ToHiLo())
	fmt.Println(udecimal.MustParse("184467440737095516160.1844674407370955161").ToHiLo())
	// Output:
	// false 0 123 2 true
	// false 669260594276348691 15255105882844922133 19 true
	// true 669260594276348691 15255105882844922133 19 true
	// false 0 123 3 true
	// true 0 123 3 true
	// false 669260594276 6432227781799973440 6 true
	// false 0 0 0 true
	// false 0 18446744073709551615 0 true
	// false 1 1 0 true
	// false 0 18446744073709551615 2 true
	// false 1 0 2 true
	// false 0 0 0 false
}

func ExampleDecimal_IsNeg() {
	fmt.Println(udecimal.MustParse("1.23").IsNeg())
	fmt.Println(udecimal.MustParse("-1.23").IsNeg())
	fmt.Println(udecimal.MustParse("0").IsNeg())
	// Output:
	// false
	// true
	// false
}

func ExampleDecimal_IsPos() {
	fmt.Println(udecimal.MustParse("1.23").IsPos())
	fmt.Println(udecimal.MustParse("-1.23").IsPos())
	fmt.Println(udecimal.MustParse("0").IsPos())
	// Output:
	// true
	// false
	// false
}

func ExampleDecimal_IsZero() {
	fmt.Println(udecimal.MustParse("1.23").IsZero())
	fmt.Println(udecimal.MustParse("0").IsZero())
	// Output:
	// false
	// true
}

func ExampleDecimal_MarshalBinary() {
	fmt.Println(udecimal.MustParse("1.23").MarshalBinary())
	fmt.Println(udecimal.MustParse("-1.2345").MarshalBinary())
	fmt.Println(udecimal.MustParse("1234567890123456789.1234567890123456789").MarshalBinary())
	// Output:
	// [0 2 11 0 0 0 0 0 0 0 123] <nil>
	// [1 4 11 0 0 0 0 0 0 48 57] <nil>
	// [0 19 19 9 73 176 246 240 2 51 19 211 181 5 249 181 241 129 21] <nil>
}

func ExampleDecimal_AppendBinary() {
	a, _ := udecimal.MustParse("1.23").AppendBinary(nil)
	b, _ := udecimal.MustParse("-1.2345").AppendBinary(nil)
	c, _ := udecimal.MustParse("1234567890123456789.1234567890123456789").AppendBinary(nil)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	// Output:
	// [0 2 11 0 0 0 0 0 0 0 123]
	// [1 4 11 0 0 0 0 0 0 48 57]
	// [0 19 19 9 73 176 246 240 2 51 19 211 181 5 249 181 241 129 21]
}

func ExampleDecimal_MarshalJSON() {
	a, _ := udecimal.MustParse("1.23").MarshalJSON()
	b, _ := udecimal.MustParse("-1.2345").MarshalJSON()
	c, _ := udecimal.MustParse("1234567890123456789.1234567890123456789").MarshalJSON()
	fmt.Println(string(a))
	fmt.Println(string(b))
	fmt.Println(string(c))
	// Output:
	// "1.23"
	// "-1.2345"
	// "1234567890123456789.1234567890123456789"
}

func ExampleDecimal_MarshalText() {
	a, _ := udecimal.MustParse("1.23").MarshalText()
	b, _ := udecimal.MustParse("-1.2345").MarshalText()
	c, _ := udecimal.MustParse("1234567890123456789.1234567890123456789").MarshalText()
	fmt.Println(string(a))
	fmt.Println(string(b))
	fmt.Println(string(c))
	// Output:
	// 1.23
	// -1.2345
	// 1234567890123456789.1234567890123456789
}

func ExampleDecimal_AppendText() {
	a, _ := udecimal.MustParse("1.23").AppendText(nil)
	b, _ := udecimal.MustParse("-1.2345").AppendText(nil)
	c, _ := udecimal.MustParse("1234567890123456789.1234567890123456789").AppendText(nil)
	fmt.Println(string(a))
	fmt.Println(string(b))
	fmt.Println(string(c))
	// Output:
	// 1.23
	// -1.2345
	// 1234567890123456789.1234567890123456789
}

func ExampleDecimal_Neg() {
	fmt.Println(udecimal.MustParse("1.23").Neg())
	fmt.Println(udecimal.MustParse("-1.23").Neg())
	// Output:
	// -1.23
	// 1.23
}

func ExampleDecimal_PowInt() {
	fmt.Println(udecimal.MustParse("1.23").PowInt(2))
	fmt.Println(udecimal.MustParse("1.23").PowInt(0))
	fmt.Println(udecimal.MustParse("1.23").PowInt(-2))
	// Output:
	// 1.5129
	// 1
	// 0.6609822195782933439
}

func ExampleDecimal_PowToIntPart() {
	fmt.Println(udecimal.MustParse("1.23").PowToIntPart(udecimal.MustParse("2.123")))
	fmt.Println(udecimal.MustParse("1.23").PowToIntPart(udecimal.MustParse("0.123")))
	fmt.Println(udecimal.MustParse("1.23").PowToIntPart(udecimal.MustParse("-2.123")))
	fmt.Println(udecimal.MustParse("0").PowToIntPart(udecimal.MustParse("0.123")))
	fmt.Println(udecimal.MustParse("0").PowToIntPart(udecimal.MustParse("-2.123")))
	// Output:
	// 1.5129 <nil>
	// 1 <nil>
	// 0.6609822195782933439 <nil>
	// 1 <nil>
	// 0 can't raise zero to a negative power

}

func ExampleDecimal_PowInt32() {
	fmt.Println(udecimal.MustParse("1.23").PowInt32(2))
	fmt.Println(udecimal.MustParse("1.23").PowInt32(0))
	fmt.Println(udecimal.MustParse("1.23").PowInt32(-2))
	fmt.Println(udecimal.MustParse("0").PowInt32(0))
	fmt.Println(udecimal.MustParse("0").PowInt32(-2))
	// Output:
	// 1.5129 <nil>
	// 1 <nil>
	// 0.6609822195782933439 <nil>
	// 1 <nil>
	// 0 can't raise zero to a negative power
}

func ExampleDecimal_Prec() {
	fmt.Println(udecimal.MustParse("1.23").Prec())
	// Output:
	// 2
}

func ExampleDecimal_PrecUint() {
	fmt.Println(udecimal.MustParse("1.23456").PrecUint())
	// Output:
	// 5
}

func ExampleDecimal_RoundBank() {
	fmt.Println(udecimal.MustParse("1.12345").RoundBank(4))
	fmt.Println(udecimal.MustParse("1.12335").RoundBank(4))
	fmt.Println(udecimal.MustParse("1.5").RoundBank(0))
	fmt.Println(udecimal.MustParse("-1.5").RoundBank(0))
	// Output:
	// 1.1234
	// 1.1234
	// 2
	// -2
}

func ExampleDecimal_RoundAwayFromZero() {
	fmt.Println(udecimal.MustParse("1.12").RoundAwayFromZero(1))
	fmt.Println(udecimal.MustParse("1.15").RoundAwayFromZero(1))
	fmt.Println(udecimal.MustParse("-1.12").RoundAwayFromZero(1))
	fmt.Println(udecimal.MustParse("-1.15").RoundAwayFromZero(1))
	// Output:
	// 1.2
	// 1.2
	// -1.2
	// -1.2
}

func ExampleDecimal_RoundHAZ() {
	fmt.Println(udecimal.MustParse("1.12345").RoundHAZ(4))
	fmt.Println(udecimal.MustParse("1.12335").RoundHAZ(4))
	fmt.Println(udecimal.MustParse("1.5").RoundHAZ(0))
	fmt.Println(udecimal.MustParse("-1.5").RoundHAZ(0))
	// Output:
	// 1.1235
	// 1.1234
	// 2
	// -2
}

func ExampleDecimal_RoundHTZ() {
	fmt.Println(udecimal.MustParse("1.12345").RoundHTZ(4))
	fmt.Println(udecimal.MustParse("1.12335").RoundHTZ(4))
	fmt.Println(udecimal.MustParse("1.5").RoundHTZ(0))
	fmt.Println(udecimal.MustParse("-1.5").RoundHTZ(0))
	// Output:
	// 1.1234
	// 1.1233
	// 1
	// -1
}

func ExampleDecimal_Scan() {
	var a udecimal.Decimal
	_ = a.Scan("1.23")
	fmt.Println(a)
	// Output:
	// 1.23
}

func ExampleDecimal_Sign() {
	fmt.Println(udecimal.MustParse("1.23").Sign())
	fmt.Println(udecimal.MustParse("-1.23").Sign())
	fmt.Println(udecimal.MustParse("0").Sign())
	// Output:
	// 1
	// -1
	// 0
}

func ExampleDecimal_Sqrt() {
	fmt.Println(udecimal.MustParse("1.21").Sqrt())
	fmt.Println(udecimal.MustParse("0").Sqrt())
	fmt.Println(udecimal.MustParse("-1.21").Sqrt())
	// Output:
	// 1.1 <nil>
	// 0 <nil>
	// 0 can't calculate square root of negative number
}

func ExampleDecimal_String() {
	fmt.Println(udecimal.MustParse("1.23").String())
	fmt.Println(udecimal.MustParse("-1.230000").String())
	// Output:
	// 1.23
	// -1.23
}

func ExampleDecimal_StringFixed() {
	fmt.Println(udecimal.MustParse("1").StringFixed(2))
	fmt.Println(udecimal.MustParse("1.23").StringFixed(4))
	fmt.Println(udecimal.MustParse("-1.230000").StringFixed(5))
	// Output:
	// 1.00
	// 1.2300
	// -1.23000
}

func ExampleDecimal_Trunc() {
	fmt.Println(udecimal.MustParse("1.23").Trunc(1))
	fmt.Println(udecimal.MustParse("-1.23").Trunc(5))
	// Output:
	// 1.2
	// -1.23
}

func ExampleDecimal_UnmarshalBinary() {
	var a udecimal.Decimal
	_ = a.UnmarshalBinary([]byte{0, 2, 11, 0, 0, 0, 0, 0, 0, 0, 123})
	fmt.Println(a)
	// Output:
	// 1.23
}

func ExampleDecimal_UnmarshalJSON() {
	var a udecimal.Decimal
	_ = a.UnmarshalJSON([]byte("1.23"))
	fmt.Println(a)
	// Output:
	// 1.23
}

func ExampleDecimal_UnmarshalText() {
	var a udecimal.Decimal
	_ = a.UnmarshalText([]byte("1.23"))
	fmt.Println(a)
	// Output:
	// 1.23
}

func ExampleDecimal_Value() {
	fmt.Println(udecimal.MustParse("1.2345").Value())
	// Output:
	// 1.2345 <nil>
}

func ExampleNullDecimal_Scan() {
	var a, b udecimal.NullDecimal
	_ = a.Scan("1.23")
	_ = b.Scan(nil)

	fmt.Println(a)
	fmt.Println(b)
	// Output:
	// {1.23 true}
	// {0 false}
}

func ExampleNullDecimal_Value() {
	fmt.Println(udecimal.NullDecimal{Decimal: udecimal.MustParse("1.2345"), Valid: true}.Value())
	fmt.Println(udecimal.NullDecimal{}.Value())
	// Output:
	// 1.2345 <nil>
	// <nil> <nil>
}

func Example_money() {
	// splitting 10.05 between two people
	share, _ := udecimal.MustParse("10.05").Div64(2)

	fmt.Println(share)
	fmt.Println(share.RoundHAZ(2))
	fmt.Println(share.RoundBank(2))
	// Output:
	// 5.025
	// 5.03
	// 5.02
}

func Example_jsonStruct() {
	type Order struct {
		Price udecimal.Decimal `json:"price"`
		Qty   udecimal.Decimal `json:"qty"`
	}

	b, _ := json.Marshal(Order{Price: udecimal.MustParse("12.50"), Qty: udecimal.MustParse("3")})
	fmt.Println(string(b))

	// quoted and unquoted numbers both decode
	var o Order
	err := json.Unmarshal([]byte(`{"price": 12.5, "qty": "3"}`), &o)
	fmt.Println(o.Price, o.Qty, err)
	// Output:
	// {"price":"12.5","qty":"3"}
	// 12.5 3 <nil>
}

func Example_errors() {
	_, err := udecimal.MustParse("1").Div(udecimal.Zero)
	fmt.Println(errors.Is(err, udecimal.ErrDivideByZero))

	_, err = udecimal.Parse("1e5")
	fmt.Println(errors.Is(err, udecimal.ErrInvalidFormat))
	// Output:
	// true
	// true
}

func ExampleAvg() {
	fmt.Println(udecimal.Avg(udecimal.MustParse("1"), udecimal.MustParse("2"), udecimal.MustParse("4")))
	// Output:
	// 2.3333333333333333333
}

func ExampleParts() {
	fmt.Println(udecimal.Parts(udecimal.MustParse("123.456")))
	fmt.Println(udecimal.Parts(udecimal.MustParse("-123.456")))
	// Output:
	// 123 0.456
	// -123 -0.456
}

func ExampleDecimal_CopySign() {
	fmt.Println(udecimal.MustParse("-10").CopySign(udecimal.MustParse("2")))
	fmt.Println(udecimal.MustParse("10").CopySign(udecimal.MustParse("-3")))
	// Output:
	// 10
	// -10
}

func Example_overview() {
	price, err := udecimal.Parse("19.99")
	if err != nil {
		panic(err)
	}

	// three items and 4.99 for shipping
	total := price.Mul(udecimal.MustParse("3")).Add(udecimal.MustParse("4.99"))

	// split three ways
	share, err := total.Div(udecimal.MustParse("3"))
	if err != nil {
		panic(err) // only when dividing by zero
	}

	fmt.Println(total)
	fmt.Println(share)
	fmt.Println(share.RoundHAZ(2))
	// Output:
	// 64.96
	// 21.6533333333333333333
	// 21.65
}
