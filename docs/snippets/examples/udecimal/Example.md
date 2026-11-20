```go
a := udecimal.MustParse("123.456")
b := udecimal.MustParse("-12.3456")

fmt.Println(a.Add(b))
fmt.Println(a.Sub(b))
fmt.Println(a.Mul(b))
fmt.Println(a.Div(b))
fmt.Println(a.Div(udecimal.MustParse("7")))
// Output:
// 111.1104
// 135.8016
// -1524.1383936
// -10 <nil>
// 17.6365714285714285714 <nil>
```
