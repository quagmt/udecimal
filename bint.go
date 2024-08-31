package udecimal

import "math/big"

type BInt struct {
	// flag to indicate if the value is overflow and stored in big.Int
	overflow bool

	// use for storing small number, with high performance
	u128 u128

	// fall back
	bigInt *big.Int
}

func (u BInt) Add(v BInt) BInt {
	if !u.overflow && !v.overflow {
		c, err := u.u128.Add(v.u128)
		if err == nil {
			return BInt{
				u128: c,
			}
		}
	}

	uBig, vBig := u.bigInt, v.bigInt

	if !u.overflow {
		uBig = u.u128.ToBigInt()
	}

	if !v.overflow {
		vBig = v.u128.ToBigInt()
	}

	return BInt{
		overflow: true,
		bigInt:   new(big.Int).Add(uBig, vBig),
	}
}
