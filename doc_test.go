package udecimal

import "fmt"

func ExampleSetDefaultPrecision() {
	a := MustParse("1.23")
	b := MustParse("4.12475")

	c, _ := a.Div(b)
	SetDefaultPrecision(10)

	fmt.Println(c)

	// Output:
	// 0.2981998909
}
