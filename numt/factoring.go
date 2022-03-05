package numt

import (
	"math/big"
)

//SmallFactors finds (non-repeated) prime factors of n smaller than max
func SmallFactors(n, max *big.Int) []*big.Int {
	two := big.NewInt(2)
	factors := []*big.Int{}
	//run while d < max
	for d := big.NewInt(3); d.Cmp(max) == -1; d.Add(d, two) {
		rem := new(big.Int).Rem(n, d)
		if rem.Cmp(new(big.Int)) == 0 && d.ProbablyPrime(20) {
			//fmt.Println(d)
			f := new(big.Int)
			f.Set(d)
			factors = append(factors, f)
		}
	}
	return factors
}
