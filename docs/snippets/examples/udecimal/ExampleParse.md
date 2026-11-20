```go
fmt.Println(udecimal.Parse("1234567890123456789.1234567890123456789"))
fmt.Println(udecimal.Parse("-1234567890123456789.1234567890123456789"))
fmt.Println(udecimal.Parse("-0.00007890123456789"))

// error cases
fmt.Println(udecimal.Parse("0.12345678901234567890123"))
fmt.Println(udecimal.Parse(""))
fmt.Println(udecimal.Parse("1.123.123"))

// Output:
// 1234567890123456789.1234567890123456789 <nil>
// -1234567890123456789.1234567890123456789 <nil>
// -0.00007890123456789 <nil>
// 0 precision out of range. Only support maximum 19 digits after the decimal point
// 0 can't parse empty string
// 0 invalid format: can't parse '1.123.123'
```
