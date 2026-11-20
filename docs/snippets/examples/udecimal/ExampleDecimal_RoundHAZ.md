```go
fmt.Println(udecimal.MustParse("1.12345").RoundHAZ(4))
fmt.Println(udecimal.MustParse("1.12335").RoundHAZ(4))
fmt.Println(udecimal.MustParse("1.5").RoundHAZ(0))
fmt.Println(udecimal.MustParse("-1.5").RoundHAZ(0))
// Output:
// 1.1235
// 1.1234
// 2
// -2
```
