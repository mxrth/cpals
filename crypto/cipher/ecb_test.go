package cipher

import (
	"bytes"
	"crypto/aes"
	"testing"

	"github.com/mxrth/cpals/crypto/rand"
)

func TestECB(t *testing.T) {
	r := make([]byte, 100*aes.BlockSize)
	orig := make([]byte, len(r))
	rand.Read(r)
	copy(orig, r)
	key := []byte("YELLOW SUBMARINE")
	c, _ := aes.NewCipher(key)
	dec := NewECBDecrypter(c)
	enc := NewECBEncrypter(c)
	enc.CryptBlocks(r, r)

	if bytes.Equal(r, orig) {
		t.Error("enc = id (probability of fluke is really low)")
	}

	dec.CryptBlocks(r, r)
	if !bytes.Equal(r, orig) {
		t.Error("dec enc != id")
	}

}
