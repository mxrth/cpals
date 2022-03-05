package encoding

import (
	"bytes"
	"testing"
)

func TestHexToBytes(t *testing.T) {
	data := map[string][]byte{
		"00":       {0},
		"":         {},
		"000102":   {0, 1, 2},
		"deadbeef": {222, 173, 190, 239},
		"111213ff": {17, 18, 19, 255},
	}
	for hex, b := range data {
		gotBytes := HexToBytes(hex)
		if !bytes.Equal(b, gotBytes) {
			t.Errorf("%v got decoded wrong. Got %v, expected %v", hex, gotBytes, b)
		}
	}
}

func TestHexToBase64(t *testing.T) {
	hex := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	base := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	if BytesToBase64(HexToBytes(hex)) != base {
		t.Error("hex was encoded wrong")
	}
}
