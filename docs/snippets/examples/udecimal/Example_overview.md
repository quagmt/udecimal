```go
price, err := udecimal.Parse("19.99")
if err != nil {
	panic(err)
}

// three items and 4.99 for shipping
total := price.Mul(udecimal.MustParse("3")).Add(udecimal.MustParse("4.99"))

// split three ways
share, err := total.Div(udecimal.MustParse("3"))
if err != nil {
	panic(err) // only when dividing by zero
}

fmt.Println(total)
fmt.Println(share)
fmt.Println(share.RoundHAZ(2))
// Output:
// 64.96
// 21.6533333333333333333
// 21.65
```
