```go
fmt.Println(udecimal.MustParse("1.12345").RoundBank(4))
fmt.Println(udecimal.MustParse("1.12335").RoundBank(4))
fmt.Println(udecimal.MustParse("1.5").RoundBank(0))
fmt.Println(udecimal.MustParse("-1.5").RoundBank(0))
// Output:
// 1.1234
// 1.1234
// 2
// -2
```
