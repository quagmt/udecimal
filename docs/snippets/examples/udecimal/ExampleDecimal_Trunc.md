```go
fmt.Println(udecimal.MustParse("1.23").Trunc(1))
fmt.Println(udecimal.MustParse("-1.23").Trunc(5))
// Output:
// 1.2
// -1.23
```
