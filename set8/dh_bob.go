package set8

import (
	"crypto/hmac"
	"crypto/sha256"
	"math/big"
)

//Tag is the type of a tag for a message under MAC
type Tag []byte

//MAC is the MAC-function used by Bob
func MAC(k *big.Int, m []byte) Tag {
	mac := hmac.New(sha256.New, k.Bytes())
	mac.Write(m)
	return mac.Sum(nil)
}

//VerifyMAC checks if t is a valid tag for m under key key
func VerifyMAC(k *big.Int, m []byte, t Tag) bool {
	return hmac.Equal([]byte(MAC(k, m)), []byte(t))
}

//DHBob simulates one party of a Diffie-Hellman Key exchange:
//given the public parameter p (prime) and his secret key x (an integer mod p)
//he computes on input h (= g^y where y is Alice's secret) K := h^y mod p and outputs a pair (m,t)
//where m is some message, and t := MAC(K, m) is some valid tag for m under K
func DHBob(p *big.Int, secret *big.Int, h *big.Int) ([]byte, Tag) {
	message := []byte("crazy flamboyant for the rap enjoyment") //message is fixed

	K := new(big.Int).Exp(h, secret, p) //K = h^x mod p
	//fmt.Println("Bob says:")
	//fmt.Println(K)

	return message, MAC(K, message)
}
