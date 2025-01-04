package udecimal

import (
	"fmt"
)

func ExampleSetDefaultPrecision() {
	SetDefaultPrecision(10)
	defer SetDefaultPrecision(19)

	a := MustParse("1.23")
	b := MustParse("4.12475")
	c, _ := a.Div(b)
	fmt.Println(c)
	// Output:
	// 0.2981998909
}

func ExampleMustFromFloat64() {
	fmt.Println(MustFromFloat64(1.234))

	// cautious: result will lose some precision when converting to decimal
	fmt.Println(MustFromFloat64(123456789.1234567890123456789))
	// Output:
	// 1.234
	// 123456789.12345679
}

func ExampleNewFromFloat64() {
	fmt.Println(NewFromFloat64(1.234))

	// cautious: result will lose some precision when converting to decimal
	fmt.Println(NewFromFloat64(123456789.1234567890123456789))
	// Output:
	// 1.234 <nil>
	// 123456789.12345679 <nil>
}

func ExampleMustFromInt64() {
	fmt.Println(MustFromInt64(123, 3))
	fmt.Println(MustFromInt64(-12345, 2))
	// Output:
	// 0.123
	// -123.45
}

func ExampleNewFromInt64() {
	fmt.Println(NewFromInt64(123, 3))
	fmt.Println(NewFromInt64(12345, 20))
	// Output:
	// 0.123 <nil>
	// 0 precision out of range. Only support maximum 19 digits after the decimal point
}

func ExampleMustFromUint64() {
	fmt.Println(MustFromUint64(123, 3))
	fmt.Println(MustFromUint64(12345, 2))
	// Output:
	// 0.123
	// 123.45
}

func ExampleNewFromUint64() {
	fmt.Println(NewFromUint64(123, 3))
	fmt.Println(NewFromUint64(12345, 2))
	fmt.Println(NewFromUint64(12345, 20))
	// Output:
	// 0.123 <nil>
	// 123.45 <nil>
	// 0 precision out of range. Only support maximum 19 digits after the decimal point
}

func ExampleMustParse() {
	fmt.Println(MustParse("1234567890123456789.1234567890123456789"))
	fmt.Println(MustParse("-1234567890123456789.1234567890123456789"))
	fmt.Println(MustParse("-0.00007890123456789"))
	// Output:
	// 1234567890123456789.1234567890123456789
	// -1234567890123456789.1234567890123456789
	// -0.00007890123456789
}

func ExampleParse() {
	fmt.Println(Parse("1234567890123456789.1234567890123456789"))
	fmt.Println(Parse("-1234567890123456789.1234567890123456789"))
	fmt.Println(Parse("-0.00007890123456789"))

	// error cases
	fmt.Println(Parse("0.12345678901234567890123"))
	fmt.Println(Parse(""))
	fmt.Println(Parse("1.123.123"))

	// Output:
	// 1234567890123456789.1234567890123456789 <nil>
	// -1234567890123456789.1234567890123456789 <nil>
	// -0.00007890123456789 <nil>
	// 0 precision out of range. Only support maximum 19 digits after the decimal point
	// 0 can't parse empty string
	// 0 invalid format: can't parse '1.123.123'
}

func ExampleNewFromHiLo() {
	fmt.Println(NewFromHiLo(false, 1, 1, 10))
	fmt.Println(NewFromHiLo(true, 0, 123456, 4))
	// Output:
	// 1844674407.3709551617 <nil>
	// -12.3456 <nil>
}

func ExampleDecimal_Abs() {
	fmt.Println(MustParse("-123.45").Abs())
	// Output:
	// 123.45
}

func ExampleDecimal_Add() {
	a := MustParse("1.23")
	b := MustParse("4.12475")
	c := a.Add(b)
	fmt.Println(c)
	// Output:
	// 5.35475
}

func ExampleDecimal_Add64() {
	a := MustParse("1.23")
	c := a.Add64(4)
	fmt.Println(c)
	// Output:
	// 5.23
}

func ExampleDecimal_Ceil() {
	fmt.Println(MustParse("1.23").Ceil())
	// Output:
	// 2
}

func ExampleDecimal_Cmp() {
	fmt.Println(MustParse("1.23").Cmp(MustParse("4.12475")))
	fmt.Println(MustParse("1.23").Cmp(MustParse("1.23")))
	fmt.Println(MustParse("1.23").Cmp(MustParse("0.12475")))
	// Output:
	// -1
	// 0
	// 1
}

func ExampleDecimal_LessThan() {
	fmt.Println(MustParse("1.23").LessThan(MustParse("4.12475")))
	fmt.Println(MustParse("1.23").LessThan(MustParse("1.23")))
	fmt.Println(MustParse("1.23").LessThan(MustParse("0.12475")))
	// Output:
	// true
	// false
	// false
}

func ExampleDecimal_LessThanOrEqual() {
	fmt.Println(MustParse("1.23").LessThanOrEqual(MustParse("4.12475")))
	fmt.Println(MustParse("1.23").LessThanOrEqual(MustParse("1.23")))
	fmt.Println(MustParse("1.23").LessThanOrEqual(MustParse("0.12475")))
	// Output:
	// true
	// true
	// false
}

func ExampleDecimal_GreaterThan() {
	fmt.Println(MustParse("1.23").GreaterThan(MustParse("4.12475")))
	fmt.Println(MustParse("1.23").GreaterThan(MustParse("1.23")))
	fmt.Println(MustParse("1.23").GreaterThan(MustParse("0.12475")))
	// Output:
	// false
	// false
	// true
}

func ExampleDecimal_GreaterThanOrEqual() {
	fmt.Println(MustParse("1.23").GreaterThanOrEqual(MustParse("4.12475")))
	fmt.Println(MustParse("1.23").GreaterThanOrEqual(MustParse("1.23")))
	fmt.Println(MustParse("1.23").GreaterThanOrEqual(MustParse("0.12475")))
	// Output:
	// false
	// true
	// true
}

func ExampleMax() {
	fmt.Println(Max(MustParse("1.23"), MustParse("4.12475")))
	// Output:
	// 4.12475
}

func ExampleMin() {
	fmt.Println(Min(MustParse("1.23"), MustParse("4.12475")))
	// Output:
	// 1.23
}

func ExampleDecimal_Equal() {
	fmt.Println(MustParse("1.123").Equal(MustParse("-1.123")))
	fmt.Println(MustParse("1.123").Equal(MustParse("1.1230000")))
	// Output:
	// false
	// true
}

func ExampleDecimal_Div() {
	fmt.Println(MustParse("1.23").Div(MustParse("4.12475")))
	fmt.Println(MustParse("1.23").Div(MustParse("0")))
	// Output:
	// 0.2981998909024789381 <nil>
	// 0 can't divide by zero
}

func ExampleDecimal_Div64() {
	fmt.Println(MustParse("1.23").Div64(4))
	fmt.Println(MustParse("1.23").Div64(0))
	// Output:
	// 0.3075 <nil>
	// 0 can't divide by zero
}

func ExampleDecimal_QuoRem() {
	fmt.Println(MustParse("1.23").QuoRem(MustParse("0.5")))
	fmt.Println(MustParse("1.23").QuoRem(MustParse("0")))
	// Output:
	// 2 0.23 <nil>
	// 0 0 can't divide by zero
}

func ExampleDecimal_Mod() {
	fmt.Println(MustParse("1.23").Mod(MustParse("0.5")))
	fmt.Println(MustParse("1.23").Mod(MustParse("0")))
	// Output:
	// 0.23 <nil>
	// 0 can't divide by zero
}

func ExampleDecimal_Sub() {
	a := MustParse("1.23")
	b := MustParse("4.12475")
	c := a.Sub(b)
	fmt.Println(c)
	// Output:
	// -2.89475
}

func ExampleDecimal_Sub64() {
	a := MustParse("1.23")
	c := a.Sub64(4)
	fmt.Println(c)
	// Output:
	// -2.77
}

func ExampleDecimal_Mul() {
	a := MustParse("1.23")
	b := MustParse("4.12475")
	c := a.Mul(b)
	fmt.Println(c)
	// Output:
	// 5.0734425
}

func ExampleDecimal_Mul64() {
	a := MustParse("1.23")
	c := a.Mul64(4)
	fmt.Println(c)
	// Output:
	// 4.92
}

func ExampleDecimal_Floor() {
	fmt.Println(MustParse("1.23").Floor())
	fmt.Println(MustParse("-1.23").Floor())
	// Output:
	// 1
	// -2
}

func ExampleDecimal_Int64() {
	fmt.Println(MustParse("1.23").Int64())
	fmt.Println(MustParse("1234567890123456789.1234567890123456789").Int64())
	fmt.Println(MustParse("-1234567890123456789.1234567890123456789").Int64())
	fmt.Println(MustParse("0.123").Int64())
	fmt.Println(MustParse("-0.123").Int64())
	fmt.Println(MustParse("12345678901234567890123456.123456").Int64())
	// Output:
	// 1 <nil>
	// 1234567890123456789 <nil>
	// -1234567890123456789 <nil>
	// 0 <nil>
	// 0 <nil>
	// 0 integer part is too large to fit in int64
}

func ExampleDecimal_InexactFloat64() {
	fmt.Println(MustParse("1.23").InexactFloat64())
	fmt.Println(MustParse("123456789.123456789").InexactFloat64())
	// Output:
	// 1.23
	// 1.2345678912345679e+08
}

func ExampleDecimal_ToHiLo() {
	fmt.Println(MustParse("1.23").ToHiLo())
	fmt.Println(MustParse("1234567890123456789.1234567890123456789").ToHiLo())
	fmt.Println(MustParse("-1234567890123456789.1234567890123456789").ToHiLo())
	fmt.Println(MustParse("0.123").ToHiLo())
	fmt.Println(MustParse("-0.123").ToHiLo())
	fmt.Println(MustParse("12345678901234567890123456.123456").ToHiLo())
	fmt.Println(MustParse("0").ToHiLo())
	fmt.Println(MustParse("18446744073709551615").ToHiLo())
	fmt.Println(MustParse("18446744073709551617").ToHiLo())
	fmt.Println(MustParse("184467440737095516.15").ToHiLo())
	fmt.Println(MustParse("184467440737095516.16").ToHiLo())
	fmt.Println(MustParse("184467440737095516160.1844674407370955161").ToHiLo())
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
	fmt.Println(MustParse("1.23").IsNeg())
	fmt.Println(MustParse("-1.23").IsNeg())
	fmt.Println(MustParse("0").IsNeg())
	// Output:
	// false
	// true
	// false
}

func ExampleDecimal_IsPos() {
	fmt.Println(MustParse("1.23").IsPos())
	fmt.Println(MustParse("-1.23").IsPos())
	fmt.Println(MustParse("0").IsPos())
	// Output:
	// true
	// false
	// false
}

func ExampleDecimal_IsZero() {
	fmt.Println(MustParse("1.23").IsZero())
	fmt.Println(MustParse("0").IsZero())
	// Output:
	// false
	// true
}

func ExampleDecimal_MarshalBinary() {
	fmt.Println(MustParse("1.23").MarshalBinary())
	fmt.Println(MustParse("-1.2345").MarshalBinary())
	fmt.Println(MustParse("1234567890123456789.1234567890123456789").MarshalBinary())
	// Output:
	// [0 2 11 0 0 0 0 0 0 0 123] <nil>
	// [1 4 11 0 0 0 0 0 0 48 57] <nil>
	// [0 19 19 9 73 176 246 240 2 51 19 211 181 5 249 181 241 129 21] <nil>
}

func ExampleDecimal_MarshalJSON() {
	a, _ := MustParse("1.23").MarshalJSON()
	b, _ := MustParse("-1.2345").MarshalJSON()
	c, _ := MustParse("1234567890123456789.1234567890123456789").MarshalJSON()
	fmt.Println(string(a))
	fmt.Println(string(b))
	fmt.Println(string(c))
	// Output:
	// "1.23"
	// "-1.2345"
	// "1234567890123456789.1234567890123456789"
}

func ExampleDecimal_MarshalText() {
	a, _ := MustParse("1.23").MarshalText()
	b, _ := MustParse("-1.2345").MarshalText()
	c, _ := MustParse("1234567890123456789.1234567890123456789").MarshalText()
	fmt.Println(string(a))
	fmt.Println(string(b))
	fmt.Println(string(c))
	// Output:
	// 1.23
	// -1.2345
	// 1234567890123456789.1234567890123456789
}

func ExampleDecimal_Neg() {
	fmt.Println(MustParse("1.23").Neg())
	fmt.Println(MustParse("-1.23").Neg())
	// Output:
	// -1.23
	// 1.23
}

func ExampleDecimal_PowInt() {
	fmt.Println(MustParse("1.23").PowInt(2))
	fmt.Println(MustParse("1.23").PowInt(0))
	fmt.Println(MustParse("1.23").PowInt(-2))
	// Output:
	// 1.5129
	// 1
	// 0.6609822195782933439
}

func ExampleDecimal_PowToIntPart() {
	fmt.Println(MustParse("1.23").PowToIntPart(MustParse("2.123")))
	fmt.Println(MustParse("1.23").PowToIntPart(MustParse("0.123")))
	fmt.Println(MustParse("1.23").PowToIntPart(MustParse("-2.123")))
	fmt.Println(MustParse("0").PowToIntPart(MustParse("0.123")))
	fmt.Println(MustParse("0").PowToIntPart(MustParse("-2.123")))
	// Output:
	// 1.5129 <nil>
	// 1 <nil>
	// 0.6609822195782933439 <nil>
	// 1 <nil>
	// 0 can't raise zero to a negative power

}

func ExampleDecimal_PowInt32() {
	fmt.Println(MustParse("1.23").PowInt32(2))
	fmt.Println(MustParse("1.23").PowInt32(0))
	fmt.Println(MustParse("1.23").PowInt32(-2))
	fmt.Println(MustParse("0").PowInt32(0))
	fmt.Println(MustParse("0").PowInt32(-2))
	// Output:
	// 1.5129 <nil>
	// 1 <nil>
	// 0.6609822195782933439 <nil>
	// 1 <nil>
	// 0 can't raise zero to a negative power
}

func ExampleDecimal_Prec() {
	fmt.Println(MustParse("1.23").Prec())
	// Output:
	// 2
}

func ExampleDecimal_PrecUint() {
	fmt.Println(MustParse("1.23456").PrecUint())
	// Output:
	// 5
}

func ExampleDecimal_RoundBank() {
	fmt.Println(MustParse("1.12345").RoundBank(4))
	fmt.Println(MustParse("1.12335").RoundBank(4))
	fmt.Println(MustParse("1.5").RoundBank(0))
	fmt.Println(MustParse("-1.5").RoundBank(0))
	// Output:
	// 1.1234
	// 1.1234
	// 2
	// -2
}

func ExampleDecimal_RoundAwayFromZero() {
	fmt.Println(MustParse("1.12").RoundAwayFromZero(1))
	fmt.Println(MustParse("1.15").RoundAwayFromZero(1))
	fmt.Println(MustParse("-1.12").RoundAwayFromZero(1))
	fmt.Println(MustParse("-1.15").RoundAwayFromZero(1))
	// Output:
	// 1.2
	// 1.2
	// -1.2
	// -1.2
}

func ExampleDecimal_RoundHAZ() {
	fmt.Println(MustParse("1.12345").RoundHAZ(4))
	fmt.Println(MustParse("1.12335").RoundHAZ(4))
	fmt.Println(MustParse("1.5").RoundHAZ(0))
	fmt.Println(MustParse("-1.5").RoundHAZ(0))
	// Output:
	// 1.1235
	// 1.1234
	// 2
	// -2
}

func ExampleDecimal_RoundHTZ() {
	fmt.Println(MustParse("1.12345").RoundHTZ(4))
	fmt.Println(MustParse("1.12335").RoundHTZ(4))
	fmt.Println(MustParse("1.5").RoundHTZ(0))
	fmt.Println(MustParse("-1.5").RoundHTZ(0))
	// Output:
	// 1.1234
	// 1.1233
	// 1
	// -1
}

func ExampleDecimal_Scan() {
	var a Decimal
	_ = a.Scan("1.23")
	fmt.Println(a)
	// Output:
	// 1.23
}

func ExampleDecimal_Sign() {
	fmt.Println(MustParse("1.23").Sign())
	fmt.Println(MustParse("-1.23").Sign())
	fmt.Println(MustParse("0").Sign())
	// Output:
	// 1
	// -1
	// 0
}

func ExampleDecimal_Sqrt() {
	fmt.Println(MustParse("1.21").Sqrt())
	fmt.Println(MustParse("0").Sqrt())
	fmt.Println(MustParse("-1.21").Sqrt())
	// Output:
	// 1.1 <nil>
	// 0 <nil>
	// 0 can't calculate square root of negative number
}

func ExampleDecimal_String() {
	fmt.Println(MustParse("1.23").String())
	fmt.Println(MustParse("-1.230000").String())
	// Output:
	// 1.23
	// -1.23
}

func ExampleDecimal_StringFixed() {
	fmt.Println(MustParse("1").StringFixed(2))
	fmt.Println(MustParse("1.23").StringFixed(4))
	fmt.Println(MustParse("-1.230000").StringFixed(5))
	// Output:
	// 1.00
	// 1.2300
	// -1.23000
}

func ExampleDecimal_Trunc() {
	fmt.Println(MustParse("1.23").Trunc(1))
	fmt.Println(MustParse("-1.23").Trunc(5))
	// Output:
	// 1.2
	// -1.23
}

func ExampleDecimal_UnmarshalBinary() {
	var a Decimal
	_ = a.UnmarshalBinary([]byte{0, 2, 11, 0, 0, 0, 0, 0, 0, 0, 123})
	fmt.Println(a)
	// Output:
	// 1.23
}

func ExampleDecimal_UnmarshalJSON() {
	var a Decimal
	_ = a.UnmarshalJSON([]byte("1.23"))
	fmt.Println(a)
	// Output:
	// 1.23
}

func ExampleDecimal_UnmarshalText() {
	var a Decimal
	_ = a.UnmarshalText([]byte("1.23"))
	fmt.Println(a)
	// Output:
	// 1.23
}

func ExampleDecimal_Value() {
	fmt.Println(MustParse("1.2345").Value())
	// Output:
	// 1.2345 <nil>
}

func ExampleNullDecimal_Scan() {
	var a, b NullDecimal
	_ = a.Scan("1.23")
	_ = b.Scan(nil)

	fmt.Println(a)
	fmt.Println(b)
	// Output:
	// {1.23 true}
	// {0 false}
}

func ExampleNullDecimal_Value() {
	fmt.Println(NullDecimal{Decimal: MustParse("1.2345"), Valid: true}.Value())
	fmt.Println(NullDecimal{}.Value())
	// Output:
	// 1.2345 <nil>
	// <nil> <nil>
}
