package set8

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/mxrth/cpals/crypto"
	"github.com/mxrth/cpals/crypto/ec"
	"github.com/mxrth/cpals/numt"
)

//Bob simulates one end of an ECDH key-exchange
//Internally holds a secret (and public params)
//generates a message plus the messages tag under the shared secret when given a message by Alice
type Bob func(h ec.Point) ([]byte, crypto.Tag)

//ECDHBob is an (yet to be curried) implementation of bob
func ECDHBob(curve ec.Curve, secret *big.Int, h ec.Point) ([]byte, crypto.Tag) {
	msg := []byte("crazy flamboyant for the rap enjoyment")
	//fmt.Println(secret)

	sharedSecret := curve.Scale(h, secret)
	//fmt.Println(sharedSecret)

	t := crypto.MAC(extractKey(sharedSecret), msg)

	return msg, t
}

func extractKey(s ec.Point) *big.Int {
	return s.X
	//return new(big.Int).Xor(s.X, s.Y)
}

func TestChallenge003(t *testing.T) {

	var publicCurve ec.Curve
	var publicBase ec.Point
	var publicOrder *big.Int

	paramsChallenge03(&publicCurve, &publicBase, &publicOrder)

	rnd := rand.New(rand.NewSource(43))

	//get 0 <= secret < publicOrder
	secret := new(big.Int).Rand(rnd, publicOrder)
	t.Logf("Secret: %v\n", secret)

	//create Bob
	bob := Bob(func(h ec.Point) ([]byte, crypto.Tag) {
		return ECDHBob(publicCurve, secret, h)
	})

	guess := invalidPointsAttack(publicCurve, publicBase, publicOrder, bob)
	guess.Mod(guess, publicOrder)

	t.Logf("Guess: %v\n", guess)

	if guess.Cmp(secret) == 0 {
		t.Log("SUCCESS!")
	} else {
		t.Error("FAILURE!")
	}
}

func invalidPointsAttack(c ec.Curve, b ec.Point, o *big.Int, bob Bob) *big.Int {
	curves, orders := loadAttackCurves(c)
	//fmt.Println(curves, orders)

	//for every attack curve, find
	primes := []*big.Int{}
	residues := []*big.Int{}

	//use primes < 2^20
	max := new(big.Int).Exp(big.NewInt(2), big.NewInt(20), nil)

	for i, c := range curves {
		order := orders[i]
		factors := numt.SmallFactors(order, max)

		for _, r := range factors {
			if contains(primes, r) {
				continue //allready have that residue
			}

			//get random point on c  of order r
			p := randomPoint(c, order, r)
			candidate := guessSecret(c, p, r, bob)
			//candidate is either secret or -secret, hence
			//candidate^2 == secret^2 mod r
			candidate.Exp(candidate, big.NewInt(2), r)

			primes = append(primes, r)
			residues = append(residues, candidate)
		}
	}

	//fmt.Println(residues, primes)
	//fmt.Println(len(primes))
	if square(o).Cmp(numt.Product(primes)) != 1 {
		//fmt.Println("have enough primes")
	} else {
		panic("not enough primes")
	}

	res := numt.CRT(residues, primes)
	//fmt.Println(res)
	//res.Mod(res, square(publicOrder))
	//res = s^2 mod product(primes)

	sqrt := new(big.Int).Sqrt(res)

	if res.Cmp(new(big.Int).Exp(sqrt, big.NewInt(2), nil)) != 0 {
		panic("res wasn't a square")
	}

	//fmt.Println(sqrt)

	return sqrt
	//return big.NewInt(0)
}

func square(n *big.Int) *big.Int {
	return new(big.Int).Exp(n, big.NewInt(2), nil)
}

func guessSecret(c ec.Curve, p ec.Point, r *big.Int, bob Bob) *big.Int {
	m, t := bob(p) //send p to bob

	//bruteforce to find b s.t. secret = b mod r
	b := big.NewInt(0)
	candidate := big.NewInt(0)
	for ; b.Cmp(r) < 0; b.Add(b, big.NewInt(1)) {
		secretGuess := c.Scale(p, b)
		if crypto.VerifyMAC(extractKey(secretGuess), m, t) {
			//when just extracting x-coordinate as key: secretGuess \in {secret, -secret}
			//fmt.Printf("Found candidate %v for prime %v\n", b, r)
			candidate.Set(b)
			break
		}
	}
	if candidate.Cmp(r) == 0 {
		panic("No suitable residue found")
	}
	return candidate
}

func loadAttackCurves(c ec.Curve) ([]ec.Curve, []*big.Int) {
	//alternative curves:
	//    y^2 = x^3 - 95051*x + 210
	//    y^2 = x^3 - 95051*x + 504
	//    y^2 = x^3 - 95051*x + 727
	bs := []*big.Int{big.NewInt(210), big.NewInt(504), big.NewInt(727), big.NewInt(21), big.NewInt(22)}

	curves := []ec.Curve{}
	for _, b := range bs {
		c := ec.Curve{A: c.A, B: b, P: c.P}
		curves = append(curves, c)
	}

	orderStrings := []string{"233970423115425145550826547352470124412", "233970423115425145544350131142039591210", "233970423115425145545378039958152057148",
		"233970423115425145528846051484558945812", "233970423115425145541000636256027604090"}
	orders := []*big.Int{}
	for _, s := range orderStrings {
		o, ok := new(big.Int).SetString(s, 10)
		if !ok {
			panic("Failed parsing order")
		}
		orders = append(orders, o)
	}

	return curves, orders
}

func paramsChallenge03(publicCurve *ec.Curve, publicBase *ec.Point, publicOrder **big.Int) {
	a := big.NewInt(-95051)
	b := big.NewInt(11279326)

	p := new(big.Int)
	_, ok := p.SetString("233970423115425145524320034830162017933", 10)

	if !ok {
		panic("Parsing p failed")
	}

	*publicCurve = ec.Curve{A: a, B: b, P: p}

	x := big.NewInt(182)

	y := new(big.Int)
	_, ok = y.SetString("85518893674295321206118380980485522083", 10)

	*publicBase = ec.Point{X: x, Y: y}

	if !ok || !publicCurve.IsOnCurve(*publicBase) {
		panic("Failed parsing public base point")
	}

	*publicOrder = new(big.Int)
	_, ok = (*publicOrder).SetString("29246302889428143187362802287225875743", 10)

	if !ok || !ec.Equal(publicCurve.Scale(*publicBase, *publicOrder), ec.O) {
		panic("Failed parsing order of public base point")
	}
}
