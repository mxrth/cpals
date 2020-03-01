package set2

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"

	"github.com/mxrth/cpals/crypto"
	ccipher "github.com/mxrth/cpals/crypto/cipher"
	"github.com/mxrth/cpals/crypto/rand"
)

func playChallenge011() bool {
	payload := make([]byte, 10*aes.BlockSize) //send all zeroes
	enc, ecb := encryptionOracle(payload)

	//use code from set1
	guessEcb := false
	for j := 0; j < len(enc)/aes.BlockSize; j++ {
		for k := j + 1; k < len(enc)/aes.BlockSize; k++ {
			if bytes.Equal(enc[j*aes.BlockSize:(j+1)*aes.BlockSize], enc[k*aes.BlockSize:(k+1)*aes.BlockSize]) {
				guessEcb = true
			}
		}
	}

	return guessEcb == ecb
}

func encryptionOracle(in []byte) ([]byte, bool) {
	var key = make([]byte, 16)
	_, err := rand.Read(key)

	if err != nil {
		panic(err)
	}

	pre := rand.Intn(5) + 5
	post := rand.Intn(5) + 5

	in = prependRandom(pre, in)
	in = appendRandom(post, in)

	c, _ := aes.NewCipher(key)

	var crypter cipher.BlockMode

	ecb := rand.Bool()

	if ecb {
		crypter = ccipher.NewECBEncrypter(c)
	} else {
		iv := make([]byte, aes.BlockSize)
		rand.Read(iv)
		crypter = ccipher.NewCBCEncrypter(c, iv)
	}

	in = crypto.PadPKCS7(in, crypter.BlockSize())
	crypter.CryptBlocks(in, in)

	return in, ecb
}

func prependRandom(n int, b []byte) []byte {
	out := make([]byte, n)
	_, err := rand.Read(out)
	if err != nil {
		panic(err)
	}
	out = append(out, b...)
	return out
}

func appendRandom(n int, b []byte) []byte {
	a := make([]byte, n)
	_, err := rand.Read(a)
	if err != nil {
		panic(err)
	}
	return append(b, a...)
}
