```go
fmt.Println(udecimal.Parts(udecimal.MustParse("123.456")))
fmt.Println(udecimal.Parts(udecimal.MustParse("-123.456")))
// Output:
// 123 0.456
// -123 -0.456
```
