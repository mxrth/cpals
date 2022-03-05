package ec

import (
	"math/big"
)

//Curve is a Elliptic curve in weierstrass-form over F_P defined by
// y^2 = x^3 + Ax +B
type Curve struct {
	A *big.Int
	B *big.Int
	P *big.Int
}

//Point is a point on a elliptic curve
type Point struct {
	X *big.Int
	Y *big.Int
}

//O is neutral element ("Point at infinity")
var O Point

func init() {
	O = Point{big.NewInt(-1), big.NewInt(0)}
}

//Equal tests for equality of two points
func Equal(p1, p2 Point) bool {
	if p1.X.Cmp(p2.X) != 0 {
		return false
	}
	if p1.Y.Cmp(p2.Y) != 0 {
		return false
	}
	return true
}

//IsOnCurve checks if point p is on curve c
func (c *Curve) IsOnCurve(p Point) bool {
	if Equal(O, p) {
		return true
	}

	//Check if y^2 = x^3 + Ax + B
	ySquared := new(big.Int).Exp(p.Y, big.NewInt(2), c.P)

	//rhs = X^3 (mod P)
	rhs := new(big.Int).Exp(p.X, big.NewInt(3), c.P)
	//rhs += Ax
	tmp := new(big.Int).Mul(c.A, p.X)
	rhs.Add(rhs, tmp)
	//rhs += B
	rhs.Add(rhs, c.B)
	rhs.Mod(rhs, c.P)
	return rhs.Cmp(ySquared) == 0
}

//Add adds the two points p1 and p2 on c
func (c *Curve) Add(p1, p2 Point) Point {
	if Equal(O, p1) {
		return p2
	}
	if Equal(p2, O) {
		return p1
	}
	if Equal(p1, c.Invert(p2)) {
		return O
	}

	var m *big.Int

	if Equal(p1, p2) {
		//m := (3*x1^2 + a) / 2*y1
		//m = 3x_1^2
		m = new(big.Int)
		m.Mul(p1.X, p1.X).Mul(m, big.NewInt(3))
		//m += a
		m.Add(m, c.A)
		tmp := new(big.Int).Lsh(p1.Y, 1)
		tmp.ModInverse(tmp, c.P)
		m.Mul(m, tmp)
	} else {
		//m := (y2 - y1) / (x2 - x1)
		m = new(big.Int).Sub(p2.Y, p1.Y)
		tmp := new(big.Int).Sub(p2.X, p1.X)
		tmp.ModInverse(tmp, c.P)
		m.Mul(m, tmp)
	}

	//x3 := m^2 - x1 - x2
	x3 := new(big.Int).Mul(m, m)
	x3.Sub(x3, p1.X)
	x3.Sub(x3, p2.X)
	x3.Mod(x3, c.P)

	//y3 := m*(x1 - x3) - y1
	y3 := new(big.Int).Mul(m, new(big.Int).Sub(p1.X, x3))
	y3.Sub(y3, p1.Y)
	y3.Mod(y3, c.P)
	return Point{x3, y3}
}

var cThree, cTwo *big.Int

func fastAdd(p1, p2, res Point, A, P, m, t1 *big.Int) {
	if Equal(p1, O) {
		res.X.Set(p2.X)
		res.Y.Set(p2.Y)
		return
	}
	if Equal(p2, O) {
		res.X.Set(p1.X)
		res.Y.Set(p1.Y)
		return
	}
	//p1 = -p2?
	if p1.X.Cmp(p2.X) == 0 && t1.Sub(P, p1.Y).Cmp(p2.Y) == 0 {
		res.X.Set(O.X)
		res.Y.Set(O.Y)
		return
	}

	if Equal(p1, p2) {
		//m := (3*x1^2 + a) / 2*y1
		//m = 3x_1^2
		m.Mul(p1.X, p1.X).Mul(m, cThree)
		//m += a
		m.Add(m, A)
		t1.Lsh(p1.Y, 1)
		t1.Exp(t1, new(big.Int).Sub(P, cTwo), P)
		m.Mul(m, t1)
	} else {
		//m := (y2 - y1) / (x2 - x1)
		m.Sub(p2.Y, p1.Y)
		t1.Sub(p2.X, p1.X)
		t1.Exp(t1, new(big.Int).Sub(P, cTwo), P)
		m.Mul(m, t1)
	}
	//m.Mod(m, P)

	//x3 := m^2 - x1 - x2
	res.X.Mul(m, m)
	res.X.Sub(res.X, p1.X)
	res.X.Sub(res.X, p2.X)
	res.X.Mod(res.X, P)

	//y3 := m*(x1 - x3) - y1
	res.Y.Mul(m, t1.Sub(p1.X, res.X))
	res.Y.Sub(res.Y, p1.Y)
	res.Y.Mod(res.Y, P)
}

func init() {
	cThree = big.NewInt(3)
	cTwo = big.NewInt(2)
}

//Scale computes n*P in c
//Double and add
func (c *Curve) Scale(p Point, n *big.Int) Point {
	result := Point{new(big.Int), new(big.Int)}
	result.X.Set(O.X)
	result.Y.Set(O.Y)

	m, t1 := new(big.Int), new(big.Int)
	r := Point{new(big.Int), new(big.Int)}
	x := Point{new(big.Int).Set(p.X), new(big.Int).Set(p.Y)}
	for i := 0; i < n.BitLen(); i++ {
		if n.Bit(i) == 1 {
			fastAdd(result, x, r, c.A, c.P, m, t1)
			result.X.Set(r.X)
			result.Y.Set(r.Y)
			//result = c.Add(result, p)
		}
		fastAdd(x, x, r, c.A, c.P, m, t1)
		x.X.Set(r.X)
		x.Y.Set(r.Y)
	}
	return result
}

//Invert inverts a point on c
func (c *Curve) Invert(p Point) Point {
	newY := new(big.Int).Sub(c.P, p.Y)
	return Point{new(big.Int).Set(p.X), newY}
}
