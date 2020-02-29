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
