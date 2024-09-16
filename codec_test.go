package udecimal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type A struct {
	P Decimal `json:"a"`
}

func TestMarshalJSON(t *testing.T) {

	a := A{
		P: MustParse("1.23"),
	}

	b, err := json.Marshal(a)
	require.NoError(t, err)

	// unmarshal back
	var c A
	require.NoError(t, json.Unmarshal(b, &c))

	require.Equal(t, a, c)
}
