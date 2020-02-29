package crypto

import "bytes"

func PadPKCS7(blocks []byte, blockLen int) []byte {
	padding := blockLen - (len(blocks) % blockLen)
	return append(blocks, bytes.Repeat([]byte{byte(padding)}, padding)...)
}
