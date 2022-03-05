package set8

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/mxrth/cpals/set8"
)

//Set in init()
var p, q, j, g, y1, y2, qSubOne *big.Int

var one = big.NewInt(1)

func TestChallenge002(t *testing.T) {
	//Test the cangaroo
	testKangaroo(y1)
	testKangaroo(y2)

	//fmt.Println(time.Since(start))
	//DH-Breaking begins here\\

	//Generate a secret between 1 and q
	rnd := rand.New(rand.NewSource(44))
	secret := new(big.Int).Rand(rnd, qSubOne)
	secret.Add(secret, one)

	public := new(big.Int).Exp(g, secret, p)

	//Bob is the curried DHBob which only takes a public param h and then outputs a tag
	bob := set8.Bob(func(h *big.Int) ([]byte, set8.Tag) {
		return set8.DHBob(p, secret, h)
	})

	guess := breakDH(p, g, q, j, public, bob) //Let's break it

	fmt.Println(guess)
	fmt.Println(secret)

	if guess.Cmp(secret) == 0 {
		fmt.Println("SUCCESS!")
	} else {
		fmt.Println("FAILURE!")
	}
}

func breakDH(p, g, q, j, public *big.Int, bob set8.Bob) *big.Int {
	n, r := set8.DHSmallSubgroupAttack(p, g, q, j, bob)
	fmt.Println("done smallsubgroupattack")

	gPrime := new(big.Int).Exp(g, r, p)
	yPrime := new(big.Int).Mul(public, new(big.Int).Exp(g, new(big.Int).Sub(q, n), p))

	upper := new(big.Int).Sub(q, big.NewInt(1))
	upper.Div(upper, r)

	m := set8.Kangaroo(p, q, j, gPrime, yPrime, big.NewInt(0), upper)
	m.Mul(m, r)
	m.Add(m, n)
	return m
}

func testKangaroo(y *big.Int) {
	start := time.Now()
	a := new(big.Int)
	b := new(big.Int).Exp(big.NewInt(2), big.NewInt(40), nil)
	result := set8.Kangaroo(p, q, j, g, y, a, b)
	fmt.Println(result)
	tmp := new(big.Int)

	fmt.Println(tmp.Exp(g, result, p))
	fmt.Println(y)

	if tmp.Cmp(y) == 0 {
		fmt.Println("SUCCESS!")
	} else {
		panic("FAILURE!")
	}
	fmt.Println(time.Since(start))
}

func init() {
	p = new(big.Int)
	p.SetString("11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623", 10)

	g = new(big.Int)
	g.SetString("622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357", 10)

	q = new(big.Int)
	q.SetString("335062023296420808191071248367701059461", 10)

	qSubOne = new(big.Int)
	qSubOne.Sub(q, big.NewInt(1))

	j = new(big.Int)
	j.SetString("34233586850807404623475048381328686211071196701374230492615844865929237417097514638999377942356150481334217896204702", 10)

	y1 = new(big.Int)
	y1.SetString("7760073848032689505395005705677365876654629189298052775754597607446617558600394076764814236081991643094239886772481052254010323780165093955236429914607119", 10)

	y2 = new(big.Int)
	y2.SetString("9388897478013399550694114614498790691034187453089355259602614074132918843899833277397448144245883225611726912025846772975325932794909655215329941809013733", 10)
}
