package udecimal

import "testing"

func TestParse(t *testing.T) {
	testcases := []string{
		"0",
		"1",
		"123",
		"123.456",
		"123.456789012345678901",
		"123456789.123456789",
		"-123456789123456789.123456789123456789",
		"-123456.123456",
	}

	for _, tc := range testcases {
		d, err := Parse(tc)
		if err != nil {
			t.Errorf("Parse(%s) failed: %v", tc, err)
		}
		t.Logf("%s -> %v", tc, d)
	}
}
