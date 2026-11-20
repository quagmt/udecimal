```go
_, err := udecimal.MustParse("1").Div(udecimal.Zero)
fmt.Println(errors.Is(err, udecimal.ErrDivideByZero))

_, err = udecimal.Parse("1e5")
fmt.Println(errors.Is(err, udecimal.ErrInvalidFormat))
// Output:
// true
// true
```
