package set8

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/mxrth/cpals/encoding"
)

//public params p,g,q and j are set in params.go
//p: prime
//g: element of (Z/pZ)* of order q
//q: order of g
//j: j = (p-1) / q
var p, pSubOne, g, q, qSubOne, j *big.Int

var one = big.NewInt(1)

func TestChallenge001(t *testing.T) {

	//Generate a secret between 1 and q
	rnd := rand.New(rand.NewSource(44))
	secret := new(big.Int).Rand(rnd, qSubOne)
	secret.Add(secret, one)

	//Bob is the curried DHBob which only takes a public param h and then outputs a tag
	bob := Bob(func(h *big.Int) ([]byte, Tag) {
		return DHBob(p, secret, h)
	})

	guess, _ := DHSmallSubgroupAttack(p, g, q, j, bob) //Let's break it

	fmt.Println(guess)
	fmt.Println(secret)

	if guess.Cmp(secret) == 0 {
		fmt.Println("SUCCESS!")
	} else {
		fmt.Println("FAILURE!")
	}
	var haveHex = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	var wantBase = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	gotBase, err := encoding.HexToBase64(haveHex)
	if err != nil {
		t.Error(err)
	}
	if gotBase != wantBase {
		t.Errorf("Wrong hex (%s) from base64", gotBase)
	}
}

func init() {
	p = new(big.Int)
	p.SetString("7199773997391911030609999317773941274322764333428698921736339643928346453700085358802973900485592910475480089726140708102474957429903531369589969318716771", 10)

	pSubOne = new(big.Int)
	pSubOne.Sub(p, big.NewInt(1))

	g = new(big.Int)
	g.SetString("4565356397095740655436854503483826832136106141639563487732438195343690437606117828318042418238184896212352329118608100083187535033402010599512641674644143", 10)

	q = new(big.Int)
	q.SetString("236234353446506858198510045061214171961", 10)

	qSubOne = new(big.Int)
	qSubOne.Sub(q, big.NewInt(1))

	j = new(big.Int)
	j.SetString("30477252323177606811760882179058908038824640750610513771646768011063128035873508507547741559514324673960576895059570", 10)
}
