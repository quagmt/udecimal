```go
fmt.Println(udecimal.MustParse("1.23").InexactFloat64())
fmt.Println(udecimal.MustParse("123456789.123456789").InexactFloat64())
// Output:
// 1.23
// 1.2345678912345679e+08
```
