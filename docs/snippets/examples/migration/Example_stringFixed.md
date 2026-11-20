```go
fmt.Println("shopspring:", decimal.RequireFromString("1.235").StringFixed(2))
fmt.Println("udecimal:  ", udecimal.MustParse("1.235").StringFixed(2))
fmt.Println("udecimal:  ", udecimal.MustParse("1.235").RoundHAZ(2).StringFixed(2))
// Output:
// shopspring: 1.24
// udecimal:   1.235
// udecimal:   1.24
```
