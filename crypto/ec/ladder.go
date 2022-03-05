package ec

import (
	"fmt"
	"math/big"
)

//MCurve represents a elliptic curve in montgomery form, that is curves of the form
// B*v^2 = u^3 + Au^2 + u
// over GF(P)
type MCurve struct {
	A *big.Int
	B *big.Int
	P *big.Int
}

//Ladder implements the Montgomery ladder
func (c MCurve) Ladder(u, k *big.Int) *big.Int {
	//u2, w2 := (1, 0)
	u2, w2 := big.NewInt(1), big.NewInt(0)
	//u3, w3 := (u, 1)
	u3, w3 := new(big.Int).Set(u), big.NewInt(1)

	var s, t, z = new(big.Int), new(big.Int), new(big.Int)
	u2sq, w2sq, u2w2 := new(big.Int), new(big.Int), new(big.Int)

	var q big.Int

	//for i in reverse(range(bitlen(p))):
	for i := c.P.BitLen() - 1; i >= 0; i-- {
		//b := 1 & (k >> i)
		b := k.Bit(i)
		fmt.Printf("i = %v, b = %v\n", i, b)
		//u2, u3 := cswap(u2, u3, b)
		u2, u3 = cswap(u2, u3, b)
		//w2, w3 := cswap(w2, w3, b)
		w2, w3 = cswap(w2, w3, b)

		//u3, w3 := ((u2*u3 - w2*w3)^2,
		//           u * (u2*w3 - w2*u3)^2)

		//s = (u2*u3 - w2*w3)^2
		s.Mul(u2, u3)
		s.Sub(s, z.Mul(w2, w3))
		s.Mul(s, s)

		//big.Int.Mod calls QuoRem internally anyways, so we take the direkt path
		q.QuoRem(s, c.P, s)

		//t = u*(u2*w3 - w2*u3)^2
		t.Mul(u2, w3)
		t.Sub(t, z.Mul(w2, u3))

		t.Mul(t, t)
		t.Mul(t, u)
		//t.Mod(t, c.P)
		q.QuoRem(t, c.P, t)

		u3.Set(s)
		w3.Set(t)

		//u2, w2 := ((u2^2 - w2^2)^2,
		//           4*u2*w2 * (u2^2 + A*u2*w2 + w2^2))
		u2sq.Mul(u2, u2)
		w2sq.Mul(w2, w2)

		//s = (u2^2 - w2^2)^2
		s.Sub(u2sq, w2sq)
		s.Mul(s, s)
		//s.Mod(s, c.P)
		q.QuoRem(s, c.P, s)

		//t =  4*u2*w2 * (u2^2 + A*u2*w2 + w2^2))

		u2w2.Mul(u2, w2)

		// z = A*u2*w2 + u2^2 + w2^2
		t.Mul(c.A, u2w2)
		t.Add(t, u2sq)
		t.Add(t, w2sq)

		t.Mul(u2w2, t)
		t.Lsh(t, 2)

		//t.Mod(t, c.P)
		q.QuoRem(t, c.P, t)

		u2.Set(s)
		w2.Set(t)
		//u2, u3 := cswap(u2, u3, b)
		u2, u3 = cswap(u2, u3, b)
		//w2, w3 := cswap(w2, w3, b)
		w2, w3 = cswap(w2, w3, b)
		fmt.Printf("u2 = %v, u3 = %v\nw2 = %v, w3 = %v\n\n", u2, u3, w2, w3)
	}

	//return u2 * w2^(p-2)
	s.Mul(u2, t.Exp(w2, z.Sub(c.P, lTwo), c.P))
	//s.Mul(u2, t.ModInverse(w2, c.P))
	return s.Mod(s, c.P)
}

var lTwo *big.Int

func init() {
	lTwo = big.NewInt(2)
}

//swaps u,v if b==1, leaves if b==0
//return  (1-b)*u + b*v, b*u + (1-b)*v
func cswap(u, v *big.Int, b uint) (*big.Int, *big.Int) {
	if b == 0 {
		return u, v
	}
	return v, u
}
