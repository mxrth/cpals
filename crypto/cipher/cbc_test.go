package cipher

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"testing"
)

func TestCBCMode(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	plain := make([]byte, 100*aes.BlockSize)
	rand.Read(plain)
	orig := make([]byte, len(plain))
	copy(orig, plain)
	iv := make([]byte, aes.BlockSize)
	rand.Read(iv)

	block, _ := aes.NewCipher(key)

	encrypter := cipher.NewCBCEncrypter(block, iv)
	decrypter := cipher.NewCBCDecrypter(block, iv)

	encrypter.CryptBlocks(plain, plain)

	if bytes.Equal(orig, plain) {
		t.Error("encryption does nothing")
	}

	decrypter.CryptBlocks(plain, plain)

	if !bytes.Equal(plain, orig) {
		t.Error("dec enc != id")
	}

}
