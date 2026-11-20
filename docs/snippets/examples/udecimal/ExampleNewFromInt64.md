```go
fmt.Println(udecimal.NewFromInt64(123, 3))
fmt.Println(udecimal.NewFromInt64(12345, 20))
// Output:
// 0.123 <nil>
// 0 precision out of range. Only support maximum 19 digits after the decimal point
```
