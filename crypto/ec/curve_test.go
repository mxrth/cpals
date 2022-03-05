package ec

import (
	"math/big"
	"testing"
)

var c Curve
var base Point
var order *big.Int

func init() {
	A := new(big.Int)
	A.SetString("-95051", 10)

	B := new(big.Int)
	B.SetString("11279326", 10)

	P := new(big.Int)
	P.SetString("233970423115425145524320034830162017933", 10)

	x := new(big.Int)
	x.SetString("182", 10)

	y := new(big.Int)
	y.SetString("85518893674295321206118380980485522083", 10)

	base = Point{x, y}

	order = new(big.Int)
	order.SetString("29246302889428143187362802287225875743", 10)

	c = Curve{A, B, P}
}

func TestEC(t *testing.T) {
	if !c.IsOnCurve(base) {
		t.Errorf("Basepoint not on curve")
	}

	if !Equal(O, c.Scale(base, big.NewInt(0))) {
		t.Errorf("Scale(base, 0) != O")
	}

	if !Equal(base, c.Scale(base, big.NewInt(1))) {
		t.Errorf("Scale(base, 1) != base")
	}

	if !Equal(O, c.Scale(base, order)) {
		t.Errorf("Scale(base, order)  != O")
	}
}

func BenchmarkScale(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c.Scale(base, order)
	}
}

func TestAdd(t *testing.T) {
	A := new(big.Int)
	A.SetString("2", 10)

	B := new(big.Int)
	B.SetString("2", 10)

	P := new(big.Int)
	P.SetString("11", 10)

	x := new(big.Int)
	x.SetString("1", 10)

	y := new(big.Int)
	y.SetString("4", 10)

	p1 := Point{x, y}

	x2 := big.NewInt(2)
	y2 := big.NewInt(5)

	p2 := Point{x2, y2}

	c := &Curve{A, B, P}
	e := Point{big.NewInt(9), big.NewInt(10)}

	e2 := Point{big.NewInt(2), big.NewInt(5)}

	res := c.Add(p1, p2)
	res2 := c.Add(p1, p1)
	res3 := c.Scale(p1, big.NewInt(9))

	if !Equal(e, res) {
		t.Errorf("Does not add up")
	}

	if !Equal(e2, res2) {
		t.Errorf("Double trouble")
	}

	//fmt.Println(res3)

	if !Equal(res3, O) {
		t.Errorf("Scaling doesn't check out")
	}
}
