```go
fmt.Println(udecimal.MustParse("1.23").GreaterThan(udecimal.MustParse("4.12475")))
fmt.Println(udecimal.MustParse("1.23").GreaterThan(udecimal.MustParse("1.23")))
fmt.Println(udecimal.MustParse("1.23").GreaterThan(udecimal.MustParse("0.12475")))
// Output:
// false
// false
// true
```
