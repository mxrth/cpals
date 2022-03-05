package numt

import (
	"fmt"
	"math/big"
)

//CRT solves a system of congruences via the chinese remainder theorem:
//given numbers r_1,...,r_n pairwise coprime and m = r_1*...*r_n as well as congruences
//	x = a_1 mod r_1
//	...
//	x = a_n mod r_n
// this function computes the (unique up to a multiple of m) solution x
func CRT(as, rs []*big.Int) *big.Int {
	m := big.NewInt(1)
	for _, r := range rs {
		m.Mul(m, r)
	}
	//set m_i = m/r_i
	//since r_i and m/m_i are coprime, there are s_i, t_i s.t.:
	//	s_i*r_i + t_i*(m/m_i) = 1
	// Hence \Sum_i^n a_i*t_i*(m/m_i) is a solution
	sum := new(big.Int)
	for i, a := range as {
		r := rs[i]

		//m_i = m/r_i
		mI := new(big.Int)
		mI.Div(m, r)

		// s*r + t*m_i = 1
		var gcd, t big.Int
		(&gcd).GCD(nil, &t, r, mI)

		if (&gcd).Cmp(big.NewInt(1)) != 0 {
			panic(fmt.Sprintf("CRT: moduli not pairwise coprime? (gcd %v)", &gcd))
		}

		tmp := new(big.Int)
		tmp.Mul(mI, a)   //tmp = a_i*m_i
		tmp.Mul(tmp, &t) //tmp = a_i*t_i*m_i
		sum.Add(sum, tmp)

		sum.Mod(sum, m) //reduce sum modulo m to find minimal solution
	}
	return sum
}
