```go
_, err := decimal.Zero.PowInt32(0)
fmt.Println("shopspring:", err)

one, err := udecimal.Zero.PowInt32(0)
fmt.Println("udecimal:  ", one, err)

// a negative exponent: shopspring keeps PowPrecisionNegativeExponent (16) digits after the decimal point
a, _ := decimal.NewFromInt(3).PowInt32(-1)
b, _ := udecimal.MustParse("3").PowInt32(-1)
fmt.Println("shopspring:", a)
fmt.Println("udecimal:  ", b)
// Output:
// shopspring: cannot represent undefined value of 0**0
// udecimal:   1 <nil>
// shopspring: 0.3333333333333333
// udecimal:   0.3333333333333333333
```
