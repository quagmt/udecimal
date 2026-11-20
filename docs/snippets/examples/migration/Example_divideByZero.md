```go
func() {
	defer func() { fmt.Println("shopspring:", recover()) }()
	decimal.NewFromInt(1).Div(decimal.Zero)
}()

_, err := udecimal.MustParse("1").Div(udecimal.Zero)
fmt.Println("udecimal:  ", err)
// Output:
// shopspring: decimal division by 0
// udecimal:   can't divide by zero
```
