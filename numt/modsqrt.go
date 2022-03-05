package numt

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

//ModSqrt calculates the square root of n in Z/pZ where p is some odd prime
func ModSqrt(n, p *big.Int) *big.Int {
	//Preprocessing
	//make sure 0 <= n <p
	n.Mod(n, p)
	if n.Cmp(zero) == 0 {
		return big.NewInt(0)
	}
	//if not a quadratic residue, abort
	if LegendreSymbol(n, p) == -1 {
		return nil
	}

	pSubOne := new(big.Int).Sub(p, one)
	s, q := factorOutTwos(new(big.Int).Sub(p, one))
	x := new(big.Int).Exp(n, q, p)
	l := big.NewInt(0)

	//choose z s.t. (z/p) == -1
	z := big.NewInt(0)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for LegendreSymbol(z, p) != -1 {
		z.Rand(rnd, p)
	}
	fmt.Println(z)

	//g = z^q mod p
	g := z.Exp(z, q, p)

	sSubOne := new(big.Int).Sub(s, big.NewInt(1))
	for j := big.NewInt(1); j.Cmp(sSubOne) < 1; j.Add(j, big.NewInt(1)) {
		//t = (x*g^-l)^(2^(s-1-j)
		//t = g^((p-1)-l) mod (p)
		t := new(big.Int).Exp(g, new(big.Int).Sub(pSubOne, l), p)
		//t = t*x
		t.Mul(t, x)
		//t = t^(2^(s-1-j))
		t.Exp(t, new(big.Int).Exp(two, new(big.Int).Sub(sSubOne, j), p), p)
		if t.Cmp(new(big.Int).Sub(p, one)) == 0 {
			l.Add(l, new(big.Int).Exp(two, j, nil))
		}
	}

	q.Add(q, one)
	q.Div(q, two)

	//a = n^((q+1)/2)*g^(-l/2)
	a := new(big.Int).Exp(n, q, p)
	a.Mul(a, g.Exp(g, new(big.Int).Sub(pSubOne, l.Div(l, two)), p))
	a.Mod(a, p)

	return a
}

func factorOutTwos(m *big.Int) (s *big.Int, q *big.Int) {
	q = new(big.Int).Set(m)
	s = big.NewInt(0)

	for new(big.Int).Mod(q, two).Cmp(zero) == 0 {
		s.Add(s, one)
		q.Div(q, two)
	}
	return s, q
}

//LegendreSymbol calculates the legendresymbol (a/p)
// (a/p) == 0 iff p divides a
// (a/p) == 1 iff a is a quadratic residue mod p
// (a/p) == -1 iff a is a quadratic nonresidue mod p
func LegendreSymbol(a, p *big.Int) int {

	rem := new(big.Int).Mod(a, p)
	if rem.Cmp(zero) == 0 {
		return 0
	}

	l := new(big.Int).Sub(p, big.NewInt(1))
	l.Div(l, two)
	l.Exp(a, l, p)
	if l.Cmp(one) == 0 {
		return 1
	} else if l.Cmp(new(big.Int).Sub(p, one)) == 0 {
		return -1
	}
	fmt.Println(l)
	panic("Bug calculating legendre symbol")
}
