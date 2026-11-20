```go
a := udecimal.MustParse("1.23")
b := udecimal.MustParse("4.12475")
c := a.Sub(b)
fmt.Println(c)
// Output:
// -2.89475
```
