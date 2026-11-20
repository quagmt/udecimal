```go
fmt.Println(udecimal.MustParse("1.23").PowInt(2))
fmt.Println(udecimal.MustParse("1.23").PowInt(0))
fmt.Println(udecimal.MustParse("1.23").PowInt(-2))
// Output:
// 1.5129
// 1
// 0.6609822195782933439
```
