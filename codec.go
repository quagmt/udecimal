package udecimal

func (d Decimal) Scan(value interface{}) error {
	return nil
}

func (d Decimal) MarshalText() ([]byte, error) {
	buf := []byte("0000000000000000000000000000000000000000")
	n := d.writeToBytes(buf, true)
	return buf[n:], nil
}

func (d Decimal) UnmarshalText(text []byte) error {

	return nil
}
