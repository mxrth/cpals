package numt

import "testing"
import "math/big"

func TestCRT(t *testing.T) {
	//x = 1 mod 7
	//x = 4 mod 9
	//x = 3 mod 5
	//SOLUTION: x = 148 mod 7*9*5
	a := []*big.Int{big.NewInt(1), big.NewInt(4), big.NewInt(3)}
	r := []*big.Int{big.NewInt(7), big.NewInt(9), big.NewInt(5)}

	x := CRT(a, r)

	if x.Cmp(big.NewInt(148)) != 0 {
		t.Errorf("Wrong solution %v", x)
	}
}

func TestNonPrimeModuli(t *testing.T) {
	a := []*big.Int{big.NewInt(2), big.NewInt(4)}
	r := []*big.Int{big.NewInt(8), big.NewInt(7)}

	x := CRT(a, r)

	if x.Cmp(big.NewInt(18)) != 0 {
		t.Errorf("Wrong solution %v", x)
	}
}
