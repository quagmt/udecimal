```go
fmt.Println("shopspring:", decimal.NewFromFloat32(0.1))
fmt.Println("udecimal:  ", udecimal.MustFromFloat64(float64(float32(0.1))))
// Output:
// shopspring: 0.1
// udecimal:   0.10000000149011612
```
