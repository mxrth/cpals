package set8

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/mxrth/cpals/crypto"
	"github.com/mxrth/cpals/crypto/analysis"
)

//takes around 3.5 minutes on my machine
func TestChallenge002(t *testing.T) {
	var p, q, j, g, y1, y2, qSubOne = new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)

	paramsChallenge002(p, q, j, g, y1, y2, qSubOne)

	var one = big.NewInt(1)
	//Test the cangaroo
	t.Run("Kangaroo: y1", func(t *testing.T) {
		testKangaroo(y1, p, q, j, g, t)
	})
	t.Run("Kangaroo: y2", func(t *testing.T) {
		testKangaroo(y2, p, q, j, g, t)
	})

	//fmt.Println(time.Since(start))
	//DH-Breaking begins here

	//Generate a secret between 1 and q
	rnd := rand.New(rand.NewSource(44))
	secret := new(big.Int).Rand(rnd, qSubOne)
	secret.Add(secret, one)

	public := new(big.Int).Exp(g, secret, p)

	//Bob is the curried DHBob which only takes a public param h and then outputs a tag
	bob := analysis.Bob(func(h *big.Int) ([]byte, crypto.Tag) {
		return DHBob(p, secret, h)
	})

	guess := breakDH(p, g, q, j, public, bob, t) //Let's break it

	t.Logf("guess : %v\n", guess)
	t.Logf("secret: %v\n", secret)

	if guess.Cmp(secret) == 0 {
		t.Logf("SUCCESS!\n")
	} else {
		t.Errorf("FAILURE!")
	}
}

func breakDH(p, g, q, j, public *big.Int, bob analysis.Bob, t *testing.T) *big.Int {
	n, r := analysis.DHSmallSubgroupAttack(p, g, q, j, bob)
	t.Log("done smallsubgroupattack")

	gPrime := new(big.Int).Exp(g, r, p)
	yPrime := new(big.Int).Mul(public, new(big.Int).Exp(g, new(big.Int).Sub(q, n), p))

	upper := new(big.Int).Sub(q, big.NewInt(1))
	upper.Div(upper, r)

	m := analysis.Kangaroo(p, q, j, gPrime, yPrime, big.NewInt(0), upper)
	m.Mul(m, r)
	m.Add(m, n)
	return m
}

func testKangaroo(y, p, q, j, g *big.Int, t *testing.T) {
	start := time.Now()
	a := new(big.Int)
	b := new(big.Int).Exp(big.NewInt(2), big.NewInt(40), nil)
	result := analysis.Kangaroo(p, q, j, g, y, a, b)
	t.Log(result)
	tmp := new(big.Int)

	t.Log(tmp.Exp(g, result, p))
	t.Log(y)

	if tmp.Cmp(y) == 0 {
		t.Logf("SUCCESS!")
	} else {
		t.Errorf("FAILURE!")
	}
	t.Log(time.Since(start))
}

func paramsChallenge002(p, q, j, g, y1, y2, qSubOne *big.Int) {
	p.SetString("11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623", 10)

	g.SetString("622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357", 10)

	q.SetString("335062023296420808191071248367701059461", 10)

	qSubOne.Sub(q, big.NewInt(1))

	j.SetString("34233586850807404623475048381328686211071196701374230492615844865929237417097514638999377942356150481334217896204702", 10)

	y1.SetString("7760073848032689505395005705677365876654629189298052775754597607446617558600394076764814236081991643094239886772481052254010323780165093955236429914607119", 10)

	y2.SetString("9388897478013399550694114614498790691034187453089355259602614074132918843899833277397448144245883225611726912025846772975325932794909655215329941809013733", 10)
}
