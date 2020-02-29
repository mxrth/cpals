package cpals

import (
	"bytes"
	"crypto/rand"
)

func padPKCS7(blocks []byte, blockLen int) []byte {
	padding := blockLen - (len(blocks) % blockLen)
	return append(blocks, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

// func encryptionOracle(in []byte) (out []byte, ecb bool) {
// 	//ecb = rand.
// }

func randBool() bool {
	b := make([]byte, 1)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	if b[0]%2 == 1 {
		return true
	}
	return false
}
