```go
a := decimal.RequireFromString("-1.225")
fmt.Println("shopspring:", a.Round(2), a.RoundBank(2), a.Truncate(2))

b := udecimal.MustParse("-1.225")
fmt.Println("udecimal:  ", b.RoundHAZ(2), b.RoundBank(2), b.Trunc(2))
// Output:
// shopspring: -1.23 -1.22 -1.22
// udecimal:   -1.23 -1.22 -1.22
```
