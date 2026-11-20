```go
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
```
