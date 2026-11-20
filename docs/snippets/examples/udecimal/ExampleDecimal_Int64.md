```go
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
```
