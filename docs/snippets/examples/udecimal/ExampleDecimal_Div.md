```go
fmt.Println(udecimal.MustParse("1.23").Div(udecimal.MustParse("4.12475")))
fmt.Println(udecimal.MustParse("1.23").Div(udecimal.MustParse("0")))
// Output:
// 0.2981998909024789381 <nil>
// 0 can't divide by zero
```
