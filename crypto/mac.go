package crypto

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
