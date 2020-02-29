package set2

import (
	"bytes"
	"crypto/aes"

	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/mxrth/cpals/crypto"
	"github.com/mxrth/cpals/crypto/cipher"
)

func TestChallenge009(t *testing.T) {
	b := []byte("YELLOW SUBMARINE")
	want := []byte("YELLOW SUBMARINE\x04\x04\x04\x04")
	b = crypto.PadPKCS7(b, 20)
	if !bytes.Equal(b, want) {
		t.Error("wrong padding")
	}
}

func TestChallenge010(t *testing.T) {
	s, _ := ioutil.ReadFile("testdata/10.txt")
	b, _ := base64.StdEncoding.DecodeString(string(s))
	//b = padPKCS7(b, aes.BlockSize)
	key := []byte("YELLOW SUBMARINE")
	block, _ := aes.NewCipher(key)

	iv := bytes.Repeat([]byte{0}, block.BlockSize())

	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(b, b)
	t.Log(string(b))
}
