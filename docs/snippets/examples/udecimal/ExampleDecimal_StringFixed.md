```go
fmt.Println(udecimal.MustParse("1").StringFixed(2))
fmt.Println(udecimal.MustParse("1.23").StringFixed(4))
fmt.Println(udecimal.MustParse("-1.230000").StringFixed(5))
// Output:
// 1.00
// 1.2300
// -1.23000
```
