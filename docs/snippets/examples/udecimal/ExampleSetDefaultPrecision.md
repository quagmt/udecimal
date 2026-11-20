```go
udecimal.SetDefaultPrecision(10)
defer udecimal.SetDefaultPrecision(19)

a := udecimal.MustParse("1.23")
b := udecimal.MustParse("4.12475")
c, _ := a.Div(b)
fmt.Println(c)
// Output:
// 0.2981998909
```
