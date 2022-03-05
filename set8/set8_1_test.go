package set8

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/mxrth/cpals/crypto"
	"github.com/mxrth/cpals/crypto/analysis"
)

func TestChallenge001(t *testing.T) {

	//public params p,g,q and j are set in params.go
	//p: prime
	//g: element of (Z/pZ)* of order q
	//q: order of g
	//j: j = (p-1) / q
	var p, pSubOne, g, q, qSubOne, j = new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)
	paramsChallenge001(p, pSubOne, g, q, qSubOne, j)

	var one = big.NewInt(1)

	//Generate a secret between 1 and q
	rnd := rand.New(rand.NewSource(44))
	secret := new(big.Int).Rand(rnd, qSubOne)
	secret.Add(secret, one)

	//Bob is the curried DHBob which only takes a public param h and then outputs a tag
	bob := analysis.Bob(func(h *big.Int) ([]byte, crypto.Tag) {
		return DHBob(p, secret, h)
	})

	guess, _ := analysis.DHSmallSubgroupAttack(p, g, q, j, bob) //Let's break it

	t.Logf("guess : %v\n", guess)
	t.Logf("secret: %v\n", secret)

	if guess.Cmp(secret) == 0 {
		t.Log("SUCCESS!")
	} else {
		t.Errorf("FAILURE!")
	}
}

func paramsChallenge001(p, pSubOne, g, q, qSubOne, j *big.Int) {

	p.SetString("7199773997391911030609999317773941274322764333428698921736339643928346453700085358802973900485592910475480089726140708102474957429903531369589969318716771", 10)

	pSubOne.Sub(p, big.NewInt(1))

	g.SetString("4565356397095740655436854503483826832136106141639563487732438195343690437606117828318042418238184896212352329118608100083187535033402010599512641674644143", 10)

	q.SetString("236234353446506858198510045061214171961", 10)

	qSubOne.Sub(q, big.NewInt(1))

	j.SetString("30477252323177606811760882179058908038824640750610513771646768011063128035873508507547741559514324673960576895059570", 10)
}
