```go
for _, s := range []string{"1e5", ".5", "1."} {
	a, aErr := decimal.NewFromString(s)
	b, bErr := udecimal.Parse(s)
	fmt.Printf("%-4s shopspring: %v %v, udecimal: %v %v\n", s, a, aErr, b, bErr)
}
// Output:
// 1e5  shopspring: 100000 <nil>, udecimal: 0 invalid format: can't parse '1e5'
// .5   shopspring: 0.5 <nil>, udecimal: 0 invalid format: can't parse '.5'
// 1.   shopspring: 1 <nil>, udecimal: 0 invalid format: can't parse '1.'
```
