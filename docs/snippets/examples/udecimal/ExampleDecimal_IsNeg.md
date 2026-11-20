```go
fmt.Println(udecimal.MustParse("1.23").IsNeg())
fmt.Println(udecimal.MustParse("-1.23").IsNeg())
fmt.Println(udecimal.MustParse("0").IsNeg())
// Output:
// false
// true
// false
```
