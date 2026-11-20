```go
fmt.Println(udecimal.NullDecimal{Decimal: udecimal.MustParse("1.2345"), Valid: true}.Value())
fmt.Println(udecimal.NullDecimal{}.Value())
// Output:
// 1.2345 <nil>
// <nil> <nil>
```
