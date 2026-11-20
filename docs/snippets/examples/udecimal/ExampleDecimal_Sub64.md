```go
a := udecimal.MustParse("1.23")
c := a.Sub64(4)
fmt.Println(c)
// Output:
// -2.77
```
