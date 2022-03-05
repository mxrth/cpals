package encoding

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

//TODO: document, write tests

func HexToBase64(s string) (string, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("invalid hex")
	}
	return base64.StdEncoding.EncodeToString(b), nil

}

func Base64ToHex(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("invalid hex")
	}
	return hex.EncodeToString(b), nil
}

func HexToBytes(hex string) (bs []byte) {
	if len(hex)%2 != 0 {
		panic("length of hex string must be multiple of two")
	}
	for i := 0; i < len(hex); i += 2 {
		b := 16*charToByte(hex[i]) + charToByte(hex[i+1])
		bs = append(bs, b)
	}
	return
}

func charToByte(c byte) byte {
	if c >= '0' && c <= '9' {
		return c - '0'
	}
	if c >= 'a' && c <= 'f' {
		return c - 'a' + 10
	}
	if c >= 'A' && c <= 'F' {
		return c - 'A' + 10
	}
	panic("Invalid hex character")
}

func BytesToBase64(bs []byte) string {
	return base64.StdEncoding.EncodeToString(bs)
}
