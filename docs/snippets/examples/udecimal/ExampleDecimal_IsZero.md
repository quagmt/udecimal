```go
fmt.Println(udecimal.MustParse("1.23").IsZero())
fmt.Println(udecimal.MustParse("0").IsZero())
// Output:
// false
// true
```
