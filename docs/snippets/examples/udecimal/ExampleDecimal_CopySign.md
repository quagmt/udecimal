```go
fmt.Println(udecimal.MustParse("-10").CopySign(udecimal.MustParse("2")))
fmt.Println(udecimal.MustParse("10").CopySign(udecimal.MustParse("-3")))
// Output:
// 10
// -10
```
