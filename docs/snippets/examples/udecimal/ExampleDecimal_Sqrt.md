```go
fmt.Println(udecimal.MustParse("1.21").Sqrt())
fmt.Println(udecimal.MustParse("0").Sqrt())
fmt.Println(udecimal.MustParse("-1.21").Sqrt())
// Output:
// 1.1 <nil>
// 0 <nil>
// 0 can't calculate square root of negative number
```
