```go
fmt.Println(udecimal.MustParse("1.23").Cmp(udecimal.MustParse("4.12475")))
fmt.Println(udecimal.MustParse("1.23").Cmp(udecimal.MustParse("1.23")))
fmt.Println(udecimal.MustParse("1.23").Cmp(udecimal.MustParse("0.12475")))
// Output:
// -1
// 0
// 1
```
