```go
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

```
