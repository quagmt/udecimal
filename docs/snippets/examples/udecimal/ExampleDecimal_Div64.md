```go
fmt.Println(udecimal.MustParse("1.23").Div64(4))
fmt.Println(udecimal.MustParse("1.23").Div64(0))
// Output:
// 0.3075 <nil>
// 0 can't divide by zero
```
