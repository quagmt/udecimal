```go
fmt.Println(udecimal.MustParse("1.23").Mod(udecimal.MustParse("0.5")))
fmt.Println(udecimal.MustParse("1.23").Mod(udecimal.MustParse("0")))
// Output:
// 0.23 <nil>
// 0 can't divide by zero
```
