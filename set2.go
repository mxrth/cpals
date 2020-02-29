package cpals

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
)

func padPKCS7(blocks []byte, blockLen int) []byte {
	padding := blockLen - (len(blocks) % blockLen)
	return append(blocks, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

type cbc struct {
	iv []byte
	b  cipher.Block
}

type encryptCBC cbc
type decryptCBC cbc

func newCBCEncrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	if len(iv) != b.BlockSize() {
		panic("invalid IV")
	}
	e := &encryptCBC{b: b, iv: make([]byte, b.BlockSize())}
	copy(e.iv, iv)
	return e
}

func newCBCDecrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	if len(iv) != b.BlockSize() {
		panic("invalid IV")
	}
	e := &decryptCBC{b: b, iv: make([]byte, b.BlockSize())}
	copy(e.iv, iv)
	return e
}

// BlockSize returns the mode's block size.
func (x *encryptCBC) BlockSize() int {

	return x.b.BlockSize()
}

// CryptBlocks encrypts or decrypts a number of blocks. The length of
// src must be a multiple of the block size. Dst and src must overlap
// entirely or not at all.
//
// If len(dst) < len(src), CryptBlocks should panic. It is acceptable
// to pass a dst bigger than src, and in that case, CryptBlocks will
// only update dst[:len(src)] and will not touch the rest of dst.
//
// Multiple calls to CryptBlocks behave as if the concatenation of
// the src buffers was passed in a single run. That is, BlockMode
// maintains state and does not reset at each CryptBlocks call.
func (x *encryptCBC) CryptBlocks(dst []byte, src []byte) {
	if len(dst) < len(src) {
		panic("dst < src")
	}
	if len(src)%x.BlockSize() != 0 {
		panic("src not a multiple of blocksize")
	}

	for i := 0; i < len(src); i += x.BlockSize() {
		xor(x.iv, src[i:], x.BlockSize()) //iv' = iv \xor srcBlock
		x.b.Encrypt(dst[i:], x.iv)        // dstBlock = enc(iv') = enc(iv \xor srcBlock)
		copy(x.iv, dst[i:])               //newIV = dstBlock
	}
}

func xor(dst, src []byte, l int) {
	for i := 0; i < l; i++ {
		dst[i] ^= src[i]
	}
}

// BlockSize returns the mode's block size.
func (x *decryptCBC) BlockSize() int {
	return x.b.BlockSize()
}

// CryptBlocks encrypts or decrypts a number of blocks. The length of
// src must be a multiple of the block size. Dst and src must overlap
// entirely or not at all.
//
// If len(dst) < len(src), CryptBlocks should panic. It is acceptable
// to pass a dst bigger than src, and in that case, CryptBlocks will
// only update dst[:len(src)] and will not touch the rest of dst.
//
// Multiple calls to CryptBlocks behave as if the concatenation of
// the src buffers was passed in a single run. That is, BlockMode
// maintains state and does not reset at each CryptBlocks call.
func (x *decryptCBC) CryptBlocks(dst []byte, src []byte) {
	if len(dst) < len(src) {
		panic("dst < src")
	}
	if len(src)%x.BlockSize() != 0 {
		panic("src not a multiple of blocksize")
	}
	tmp := make([]byte, x.BlockSize()) //dst and src might be the same
	for i := 0; i < len(src); i += x.BlockSize() {
		x.b.Decrypt(tmp, src[i:])     // dstBlock = dec(srcBlock)
		xor(tmp, x.iv, x.BlockSize()) //dstBlock = dec(srcBlock) \xor x.iv
		copy(x.iv, src[i:])           //newIV = srcBlock
		copy(dst[i:], tmp)
	}
}

func encryptionOracle(in []byte) (out []byte, ecb bool) {
	//ecb = rand.
}

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
