package numt

import (
	"math/big"
)

var minusOne, zero, one, two *big.Int

func init() {
	minusOne = big.NewInt(-1)
	zero = big.NewInt(0)
	one = big.NewInt(1)
	two = big.NewInt(2)
}
