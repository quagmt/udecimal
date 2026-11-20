```go
a := udecimal.MustParse("1.23")
b := udecimal.MustParse("4.12475")
c := a.Add(b)
fmt.Println(c)
// Output:
// 5.35475
```
