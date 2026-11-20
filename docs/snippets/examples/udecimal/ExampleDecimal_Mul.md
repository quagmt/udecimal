```go
a := udecimal.MustParse("1.23")
b := udecimal.MustParse("4.12475")
c := a.Mul(b)
fmt.Println(c)
// Output:
// 5.0734425
```
