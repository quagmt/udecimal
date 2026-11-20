```go
fmt.Println(udecimal.NewFromFloat64(1.234))

// cautious: result will lose some precision when converting to decimal
fmt.Println(udecimal.NewFromFloat64(123456789.1234567890123456789))
// Output:
// 1.234 <nil>
// 123456789.12345679 <nil>
```
