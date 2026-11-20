```go
var a udecimal.Decimal
_ = a.UnmarshalBinary([]byte{0, 2, 11, 0, 0, 0, 0, 0, 0, 0, 123})
fmt.Println(a)
// Output:
// 1.23
```
