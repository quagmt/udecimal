```go
fmt.Println(udecimal.MustParse("1.12345").RoundHTZ(4))
fmt.Println(udecimal.MustParse("1.12335").RoundHTZ(4))
fmt.Println(udecimal.MustParse("1.5").RoundHTZ(0))
fmt.Println(udecimal.MustParse("-1.5").RoundHTZ(0))
// Output:
// 1.1234
// 1.1233
// 1
// -1
```
