package set8

import (
	"fmt"
	"math/big"
	"math/rand"

	"github.com/mxrth/cpals/numt"
)

//Bob is given alice's PK, calculates shared secret and computes tag of a fixed message
type Bob func(h *big.Int) ([]byte, Tag)

//DHSmallSubgroupAttack small subgroup confinement attack
//Given the public params p,g,q,j and a function Bob that simulates a (honest) DH-party,
//this function tries to recover bobs secret key
//It returns x and r s.t. if bob has secret key g^y then y = x mod r
//r might be smaller than y
func DHSmallSubgroupAttack(p, g, q, j *big.Int, bob Bob) (x, r *big.Int) {
	one := big.NewInt(1)
	rnd := rand.New(rand.NewSource(422)) //fixed seed
	primes := numt.SmallFactors(j, new(big.Int).Exp(big.NewInt(2), big.NewInt(20), nil))
	//fmt.Println("found small prime factors")
	residues := []*big.Int{}
	for _, r := range primes {
		//h = rand(1,p)^((p-1)/r) mod p
		h := big.NewInt(1)
		//search as long as h == 1
		for h.Cmp(one) == 0 {
			e := new(big.Int).Sub(p, one)
			h.Rand(rnd, e)
			h.Add(h, one)
			h.Exp(h, e.Div(e, r), p) // h = rnd^((p-1)/r) mod p
		}
		//fmt.Println(h)
		m, t := bob(new(big.Int).Set(h))
		//find b such that t = MAC(h^b,m), then we know secret = b mod r
		b := big.NewInt(0)
		for ; b.Cmp(r) != 0; b.Add(b, one) {
			k := new(big.Int).Exp(h, b, p)
			//fmt.Println("K Candidate")
			//fmt.Println(k)
			if VerifyMAC(k, m, t) {
				//fmt.Println("found valid k")
				break
			}
		}
		if b.Cmp(r) == 0 {
			panic(fmt.Sprintf("No suitable residue found for %v!, b: %v", r, b))
		}
		cpy := new(big.Int)
		cpy.Set(b)
		residues = append(residues, cpy)
	}
	//fmt.Println(residues, primes)
	return numt.CRT(residues, primes), numt.Product(primes)
}
