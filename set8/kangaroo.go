package set8

import (
	"fmt"
	"math/big"
)

//Kangaroo tries to solve the discrete logarithm problem
//Given a group G = <g> \subset Z_p^* of order q (and j = (p-1)/q)
//And a y \in G
//Kangaroo searches for an x s.t. y = g^x given that we know that a <= x <= b
func Kangaroo(p, q, j, g, y, a, b *big.Int) *big.Int {
	N := new(big.Int)
	N.SetString("2048576", 10) //parameter optimized for 'harder' challenges

	one := big.NewInt(1)

	//Generate tame sequence
	//xT := 0
	xT := new(big.Int)

	//yT := g^b
	yT := new(big.Int)
	yT.Exp(g, b, p)

	for i := big.NewInt(1); i.Cmp(N) <= 0; i.Add(i, one) { //for i = 1...N
		//xT := xT + f(yT)
		xT.Add(xT, f(yT))

		//yT := yT * g^(f(yT)) (mod p)
		tmp := new(big.Int).Exp(g, f(yT), p)
		yT.Mul(yT, tmp).Mod(yT, p)
	}

	fmt.Println("Tame sequence Done")

	//xW := 0
	xW := new(big.Int)
	//yW := y
	yW := new(big.Int).Set(y)

	intervalLength := new(big.Int).Sub(b, a)

	c := new(big.Int) //RHS of comparison

	//while xW < b - a + xT
	for xW.Cmp(c.Add(intervalLength, xW)) == -1 {
		xW.Add(xW, f(yW))

		tmp := new(big.Int).Exp(g, f(yW), p)
		yW.Mul(yW, tmp).Mod(yW, p)

		if yW.Cmp(yT) == 0 {
			xT.Add(xT, b)
			return xT.Sub(xT, xW)
		}
	}

	panic("No Rendezvous")

}

func f(y *big.Int) *big.Int {
	two := big.NewInt(2)

	k := big.NewInt(22)

	return new(big.Int).Exp(two, k.Mod(y, k), nil)
}
