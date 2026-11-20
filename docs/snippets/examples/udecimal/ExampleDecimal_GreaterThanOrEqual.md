```go
fmt.Println(udecimal.MustParse("1.23").GreaterThanOrEqual(udecimal.MustParse("4.12475")))
fmt.Println(udecimal.MustParse("1.23").GreaterThanOrEqual(udecimal.MustParse("1.23")))
fmt.Println(udecimal.MustParse("1.23").GreaterThanOrEqual(udecimal.MustParse("0.12475")))
// Output:
// false
// true
// true
```
