```go
fmt.Println(udecimal.MustParse("1.123").Equal(udecimal.MustParse("-1.123")))
fmt.Println(udecimal.MustParse("1.123").Equal(udecimal.MustParse("1.1230000")))
// Output:
// false
// true
```
