```go
fmt.Println(udecimal.MustFromFloat64(1.234))

// cautious: result will lose some precision when converting to decimal
fmt.Println(udecimal.MustFromFloat64(123456789.1234567890123456789))
// Output:
// 1.234
// 123456789.12345679
```
