```go
fmt.Println(udecimal.MustParse("1.23").Sign())
fmt.Println(udecimal.MustParse("-1.23").Sign())
fmt.Println(udecimal.MustParse("0").Sign())
// Output:
// 1
// -1
// 0
```
