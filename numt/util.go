package numt

import "math/big"

//Product computes the product of given big.Ints
func Product(ns []*big.Int) *big.Int {
	prod := big.NewInt(1)
	for _, n := range ns {
		prod.Mul(prod, n)
	}
	return prod
}
