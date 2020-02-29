package cipher

import ciph "crypto/cipher"

type ecb struct {
	b ciph.Block
}

type ecbEncrypter ecb
type ecbDecrypter ecb

func NewECBDecrypter(b ciph.Block) ciph.BlockMode {
	return &ecbDecrypter{b: b}
}

func (x *ecbDecrypter) BlockSize() int {
	return x.b.BlockSize()
}

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(dst) < len(src) {
		panic("dst < src")
	}
	if len(src)%x.BlockSize() != 0 {
		panic("src not a multiple of blocksize")
	}
	for i := 0; i < len(src); i += x.BlockSize() {
		x.b.Decrypt(dst[i:], src[i:])
	}
}

func NewECBEncrypter(b ciph.Block) ciph.BlockMode {
	return &ecbEncrypter{b: b}
}

func (x *ecbEncrypter) BlockSize() int {
	return x.b.BlockSize()
}

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(dst) < len(src) {
		panic("dst < src")
	}
	if len(src)%x.BlockSize() != 0 {
		panic("src not a multiple of blocksize")
	}
	for i := 0; i < len(src); i += x.BlockSize() {
		x.b.Encrypt(dst[i:], src[i:])
	}
}
