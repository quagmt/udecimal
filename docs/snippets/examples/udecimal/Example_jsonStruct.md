```go
type Order struct {
	Price udecimal.Decimal `json:"price"`
	Qty   udecimal.Decimal `json:"qty"`
}

b, _ := json.Marshal(Order{Price: udecimal.MustParse("12.50"), Qty: udecimal.MustParse("3")})
fmt.Println(string(b))

// quoted and unquoted numbers both decode
var o Order
err := json.Unmarshal([]byte(`{"price": 12.5, "qty": "3"}`), &o)
fmt.Println(o.Price, o.Qty, err)
// Output:
// {"price":"12.5","qty":"3"}
// 12.5 3 <nil>
```
