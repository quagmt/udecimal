```go
// splitting 10.05 between two people
share, _ := udecimal.MustParse("10.05").Div64(2)

fmt.Println(share)
fmt.Println(share.RoundHAZ(2))
fmt.Println(share.RoundBank(2))
// Output:
// 5.025
// 5.03
// 5.02
```
