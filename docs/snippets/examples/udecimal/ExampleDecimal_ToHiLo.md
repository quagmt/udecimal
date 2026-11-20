```go
fmt.Println(udecimal.MustParse("1.23").ToHiLo())
fmt.Println(udecimal.MustParse("1234567890123456789.1234567890123456789").ToHiLo())
fmt.Println(udecimal.MustParse("-1234567890123456789.1234567890123456789").ToHiLo())
fmt.Println(udecimal.MustParse("0.123").ToHiLo())
fmt.Println(udecimal.MustParse("-0.123").ToHiLo())
fmt.Println(udecimal.MustParse("12345678901234567890123456.123456").ToHiLo())
fmt.Println(udecimal.MustParse("0").ToHiLo())
fmt.Println(udecimal.MustParse("18446744073709551615").ToHiLo())
fmt.Println(udecimal.MustParse("18446744073709551617").ToHiLo())
fmt.Println(udecimal.MustParse("184467440737095516.15").ToHiLo())
fmt.Println(udecimal.MustParse("184467440737095516.16").ToHiLo())
fmt.Println(udecimal.MustParse("184467440737095516160.1844674407370955161").ToHiLo())
// Output:
// false 0 123 2 true
// false 669260594276348691 15255105882844922133 19 true
// true 669260594276348691 15255105882844922133 19 true
// false 0 123 3 true
// true 0 123 3 true
// false 669260594276 6432227781799973440 6 true
// false 0 0 0 true
// false 0 18446744073709551615 0 true
// false 1 1 0 true
// false 0 18446744073709551615 2 true
// false 1 0 2 true
// false 0 0 0 false
```
