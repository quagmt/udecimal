```go
a := udecimal.MustParse("1.23")
c := a.Mul64(4)
fmt.Println(c)
// Output:
// 4.92
```
