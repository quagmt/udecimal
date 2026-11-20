```go
a, _ := udecimal.MustParse("1.23").AppendText(nil)
b, _ := udecimal.MustParse("-1.2345").AppendText(nil)
c, _ := udecimal.MustParse("1234567890123456789.1234567890123456789").AppendText(nil)
fmt.Println(string(a))
fmt.Println(string(b))
fmt.Println(string(c))
// Output:
// 1.23
// -1.2345
// 1234567890123456789.1234567890123456789
```
