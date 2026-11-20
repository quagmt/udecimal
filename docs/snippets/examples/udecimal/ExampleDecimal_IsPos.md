```go
fmt.Println(udecimal.MustParse("1.23").IsPos())
fmt.Println(udecimal.MustParse("-1.23").IsPos())
fmt.Println(udecimal.MustParse("0").IsPos())
// Output:
// true
// false
// false
```
