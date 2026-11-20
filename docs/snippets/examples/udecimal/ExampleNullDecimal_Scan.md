```go
var a, b udecimal.NullDecimal
_ = a.Scan("1.23")
_ = b.Scan(nil)

fmt.Println(a)
fmt.Println(b)
// Output:
// {1.23 true}
// {0 false}
```
