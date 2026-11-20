```go
var a udecimal.Decimal
_ = a.UnmarshalText([]byte("1.23"))
fmt.Println(a)
// Output:
// 1.23
```
