```go
a, _ := udecimal.MustParse("1.23").AppendBinary(nil)
b, _ := udecimal.MustParse("-1.2345").AppendBinary(nil)
c, _ := udecimal.MustParse("1234567890123456789.1234567890123456789").AppendBinary(nil)
fmt.Println(a)
fmt.Println(b)
fmt.Println(c)
// Output:
// [0 2 11 0 0 0 0 0 0 0 123]
// [1 4 11 0 0 0 0 0 0 48 57]
// [0 19 19 9 73 176 246 240 2 51 19 211 181 5 249 181 241 129 21]
```
