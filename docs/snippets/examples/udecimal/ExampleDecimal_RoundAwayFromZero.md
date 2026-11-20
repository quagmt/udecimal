```go
fmt.Println(udecimal.MustParse("1.12").RoundAwayFromZero(1))
fmt.Println(udecimal.MustParse("1.15").RoundAwayFromZero(1))
fmt.Println(udecimal.MustParse("-1.12").RoundAwayFromZero(1))
fmt.Println(udecimal.MustParse("-1.15").RoundAwayFromZero(1))
// Output:
// 1.2
// 1.2
// -1.2
// -1.2
```
