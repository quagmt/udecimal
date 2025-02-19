package udecimal

const (
	maxHighPrec = 41
)

type u1024 [16]uint64

type ubig struct {
	neg  bool
	coef u1024
}

func (u ubig) IsZero() bool {
	return u.coef == u1024{}
}

var (
	ubigOne = ubig{coef: u1024{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}}
)

func ubigFromBint(coef bint) (ubig, error) {
	if !coef.overflow() {
		return ubig{coef: u1024{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, coef.u128.hi, coef.u128.lo}}, nil
	}

	// dBigBytes is in big-endian order
	dBigBytes := coef.GetBig().Bytes()

	if len(dBigBytes) > 128 {
		return ubig{}, errOverflow
	}

	var u u1024
	for i := 0; i < len(dBigBytes); i++ {
		u[i/8] |= uint64(dBigBytes[i]) << uint((i%8)*8)
	}

	return ubig{coef: u}, nil
}

func (u ubig) ToBint() bint {
	return bint{}
}

func (u ubig) Mul64(v uint64) (ubig, error) {
	return ubig{}, nil
}

func (u ubig) MulU128(v u128) (ubig, error) {
	return ubig{}, nil
}

func (u ubig) Mul(v ubig) (ubig, error) {
	return ubig{}, nil
}

func (u ubig) Add(v ubig) (ubig, error) {
	return ubig{}, nil
}

func (u ubig) Sub(v ubig) (ubig, error) {
	return ubig{}, nil
}

func (u ubig) Div(v ubig) (ubig, error) {
	return ubig{}, nil
}

func (u ubig) DivU128(v u128) (ubig, error) {
	return ubig{}, nil
}
