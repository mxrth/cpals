package analysis

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/mxrth/cpals/crypto"
	"github.com/mxrth/cpals/crypto/ec"
	"github.com/mxrth/cpals/numt"
)

//ECBob simulates one endpoint of a DH key exchange
type ECBob func(h ec.Point) ([]byte, crypto.Tag)

//ECDHInvalidPointsAttack carries out attack...
//c: Curve we operate on
//base: generator of cyclic subgroup
//order: order of base
//NOT GENERIC but explicitly bound to the curve y^2 = x^3 - 95051x + 11279326
//tries to guess bobs secret
func ECDHInvalidPointsAttack(curve ec.Curve, base ec.Point, bob ECBob) *big.Int {
	//Idea: for each curve
	//factorize the order, if we already have that factor 'r', continue, else,
	//generate point of order r, send it to bob, bruteforce secret key => know secret mod r
	//In the end: piece it all together with CRT

	A := big.NewInt(-95051)
	Bs := []*big.Int{big.NewInt(210), big.NewInt(504), big.NewInt(727)}

	orders := []*big.Int{s2big("233970423115425145550826547352470124412"),
		s2big("233970423115425145544350131142039591210"),
		s2big("233970423115425145545378039958152057148")}

	factors := []*big.Int{}
	residues := []*big.Int{}
	upper := new(big.Int).Exp(big.NewInt(2), big.NewInt(22), nil)

	for i, b := range Bs {
		primes := numt.SmallFactors(orders[i], upper)
		fmt.Println(primes)
		c := ec.Curve{A: A, B: b, P: curve.P}
		fmt.Println(c)
		for _, r := range primes {
			if contains(factors, r) {
				continue
			}
			factors = append(factors, r)

			p := randomPoint(c, orders[i], r)
			m, t := bob(p)
			//find l s.t. t = MAC((l*p).X, m)
			l := big.NewInt(0)
			one := big.NewInt(1)
			//fmt.Println(one)
			for ; l.Cmp(r) != 0; l.Add(l, one) {
				k := curve.Scale(p, l).X
				//normalize(k, curve.P)
				if crypto.VerifyMAC(k, m, t) {
					fmt.Printf("found valid k l:%v\n", l)
					break
				}
			}
			if l.Cmp(r) == 0 {
				panic(fmt.Sprintf("No suitable residue found for %v! l: %v", r, l))
			}
			cpy := new(big.Int)
			cpy.Set(l)
			residues = append(residues, cpy)
		}
	}
	fmt.Println(residues, factors)
	return numt.CRT(residues, factors)
}

func normalize(n, p *big.Int) {
	n.Mod(n, p)
	//fmt.Println(n)
	if n.Cmp(big.NewInt(0)) == -1 {
		n.Sub(p, n)
	}
}

func randomPoint(c ec.Curve, curveOrder, pointOrder *big.Int) ec.Point {
	p := ec.O
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	j := new(big.Int).Div(curveOrder, pointOrder)
	for ec.Equal(p, ec.O) {
		x := new(big.Int).Rand(rnd, c.P)
		y := computeY(c, x)
		if y == nil { //x not a quadratic residue
			continue
		}
		p = c.Scale(ec.Point{X: x, Y: y}, j)
	}
	fmt.Println(p)
	if !ec.Equal(c.Scale(p, pointOrder), ec.O) {
		panic("LOGIC ERROR")
	}
	if !c.IsOnCurve(p) {
		panic("POINT NOT ON CURVE")
	}
	return p
}

func computeY(c ec.Curve, x *big.Int) *big.Int {
	//rhs = X^3 (mod P)
	rhs := new(big.Int).Exp(x, big.NewInt(3), c.P)
	//rhs += Ax
	tmp := new(big.Int).Mul(c.A, x)
	rhs.Add(rhs, tmp)
	//rhs += B
	rhs.Add(rhs, c.B)
	rhs.Mod(rhs, c.P)

	rhs = rhs.ModSqrt(rhs, c.P)
	if rhs != nil && !c.IsOnCurve(ec.Point{X: x, Y: rhs}) {
		panic("NOT ON CURVE")
	}
	return rhs
}

func contains(s []*big.Int, n *big.Int) bool {
	for _, m := range s {
		if m.Cmp(n) == 0 {
			return true
		}
	}
	return false
}
func s2big(s string) *big.Int {
	r := new(big.Int)
	r.SetString(s, 10)
	return r
}
