```go
fmt.Println(udecimal.Parse("0.12345678901234567890123"))

udecimal.SetDefaultParseMode(udecimal.ParseModeTrunc)
defer udecimal.SetDefaultParseMode(udecimal.ParseModeError)

fmt.Println(udecimal.Parse("0.12345678901234567890123"))
// Output:
// 0 precision out of range: more decimal places than the default precision
// 0.1234567890123456789 <nil>
```
