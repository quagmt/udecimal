```go
fmt.Println(udecimal.NewFromUint64(123, 3))
fmt.Println(udecimal.NewFromUint64(12345, 2))
fmt.Println(udecimal.NewFromUint64(12345, 20))
// Output:
// 0.123 <nil>
// 123.45 <nil>
// 0 precision out of range: more decimal places than the default precision
```
