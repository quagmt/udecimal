```go
type row struct {
	A decimal.NullDecimal
	B udecimal.NullDecimal
}

b, _ := json.Marshal(row{})
fmt.Println(string(b))
// Output:
// {"A":null,"B":{"Decimal":"0","Valid":false}}
```
