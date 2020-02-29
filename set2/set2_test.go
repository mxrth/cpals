package cpals

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"testing"
)

func TestChallenge009(t *testing.T) {
	b := []byte("YELLOW SUBMARINE")
	want := []byte("YELLOW SUBMARINE\x04\x04\x04\x04")
	b = padPKCS7(b, 20)
	if !bytes.Equal(b, want) {
		t.Error("wrong padding")
	}
}

func TestPKCS7(t *testing.T) {
	//already a multiple of block size: want a full block of padding
	b := []byte("YELLOW SUBMARINE")
	want := []byte("YELLOW SUBMARINE\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10")
	b = padPKCS7(b, len(b))
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

	decrypter := newCBCDecrypter(block, iv)
	decrypter.CryptBlocks(b, b)
	t.Log(string(b))
}

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

	encrypter := newCBCEncrypter(block, iv)
	decrypter := newCBCDecrypter(block, iv)

	encrypter.CryptBlocks(plain, plain)

	if bytes.Equal(orig, plain) {
		t.Error("encryption does nothing")
	}

	decrypter.CryptBlocks(plain, plain)

	if !bytes.Equal(plain, orig) {
		t.Error("dec enc != id")
	}

}
