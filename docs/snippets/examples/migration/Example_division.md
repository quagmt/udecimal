```go
// shopspring's DivisionPrecision is 16 by default
fmt.Println("shopspring:", decimal.NewFromInt(2).Div(decimal.NewFromInt(3)))

q, _ := udecimal.MustParse("2").Div(udecimal.MustParse("3"))
fmt.Println("udecimal:  ", q)
// Output:
// shopspring: 0.6666666666666667
// udecimal:   0.6666666666666666666
```
