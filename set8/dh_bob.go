package set8

import (
	"math/big"

	"github.com/mxrth/cpals/crypto"
)

//DHBob simulates one party of a Diffie-Hellman Key exchange:
//given the public parameter p (prime) and his secret key x (an integer mod p)
//he computes on input h (= g^y where y is Alice's secret) K := h^y mod p and outputs a pair (m,t)
//where m is some message, and t := MAC(K, m) is some valid tag for m under K
func DHBob(p *big.Int, secret *big.Int, h *big.Int) ([]byte, crypto.Tag) {
	message := []byte("crazy flamboyant for the rap enjoyment") //message is fixed

	K := new(big.Int).Exp(h, secret, p) //K = h^x mod p
	//fmt.Println("Bob says:")
	//fmt.Println(K)

	return message, crypto.MAC(K, message)
}
