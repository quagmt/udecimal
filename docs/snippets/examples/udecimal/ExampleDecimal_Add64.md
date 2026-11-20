```go
a := udecimal.MustParse("1.23")
c := a.Add64(4)
fmt.Println(c)
// Output:
// 5.23
```
