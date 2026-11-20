```go
fmt.Println("shopspring:", decimal.RequireFromString("0.1234567890123456789").Mul(decimal.RequireFromString("0.1")))
fmt.Println("udecimal:  ", udecimal.MustParse("0.1234567890123456789").Mul(udecimal.MustParse("0.1")))
// Output:
// shopspring: 0.01234567890123456789
// udecimal:   0.0123456789012345678
```
