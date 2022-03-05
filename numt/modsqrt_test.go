package numt

import (
	"math/big"
	"testing"
)

func TestCrt(t *testing.T) {
	p := big.NewInt(41)
	n := big.NewInt(2)

	sqrt := ModSqrt(n, p)

	if n.Cmp(new(big.Int).Exp(sqrt, two, p)) != 0 {
		t.Errorf("Wrong sqaure root %v", sqrt)
	}
}
