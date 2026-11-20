```go
fmt.Println(udecimal.Parse("0.12345678901234567890123"))

udecimal.SetDefaultParseMode(udecimal.ParseModeTrunc)
defer udecimal.SetDefaultParseMode(udecimal.ParseModeError)

fmt.Println(udecimal.Parse("0.12345678901234567890123"))
// Output:
// 0 precision out of range. Only support maximum 19 digits after the decimal point
// 0.1234567890123456789 <nil>
```
