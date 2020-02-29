package crypto

import (
	"bytes"
	"testing"
)

func TestPKCS7(t *testing.T) {
	//already a multiple of block size: want a full block of padding
	b := []byte("YELLOW SUBMARINE")
	want := []byte("YELLOW SUBMARINE\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10")
	b = PadPKCS7(b, len(b))
	if !bytes.Equal(b, want) {
		t.Error("wrong padding")
	}
}
