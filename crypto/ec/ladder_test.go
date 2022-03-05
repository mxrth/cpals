package ec

import (
	"math/big"
	"testing"
)

func TestLadder(t *testing.T) {
	a := big.NewInt(534)
	b := big.NewInt(1)
	p, _ := new(big.Int).SetString("233970423115425145524320034830162017933", 10)

	c := MCurve{A: a, B: b, P: p}

	base := big.NewInt(4)
	//order, _ := new(big.Int).SetString("29246302889428143187362802287225875743", 10)

	four := c.Ladder(base, big.NewInt(1))

	if base.Cmp(four) != 0 {
		t.Errorf("Ladder(4, 1) != 4, %v", four)
	}

	//o := c.Ladder(base, order)
	//if o.Cmp(big.NewInt(0)) != 0 {
	//	t.Errorf("Ladder(base, order) != 0, %v", o)
	//}
}

func BenchmarkLadder(b *testing.B) {
	a := big.NewInt(534)
	b2 := big.NewInt(1)
	p, _ := new(big.Int).SetString("233970423115425145524320034830162017933", 10)

	c := MCurve{A: a, B: b2, P: p}

	base := big.NewInt(4)
	order, _ := new(big.Int).SetString("29246302889428143187362802287225875743", 10)
	// run the Ladder function b.N times
	for n := 0; n < b.N; n++ {

		c.Ladder(base, order)
	}
}
