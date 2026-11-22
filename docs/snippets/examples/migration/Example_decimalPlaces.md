```go
a, err := decimal.NewFromString("0.12345678901234567890123")
fmt.Println("shopspring:", a, err)

b, err := udecimal.Parse("0.12345678901234567890123")
fmt.Println("udecimal:  ", b, err)
// Output:
// shopspring: 0.12345678901234567890123 <nil>
// udecimal:   0 precision out of range: more decimal places than the default precision
```
