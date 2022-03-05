package crypto

import (
	"bytes"
	"testing"

	"github.com/mxrth/cpals/encoding"
)

func TestFixedXOR(t *testing.T) {
	a := encoding.HexToBytes("1c0111001f010100061a024b53535009181c")
	b := encoding.HexToBytes("686974207468652062756c6c277320657965")
	c := encoding.HexToBytes("746865206b696420646f6e277420706c6179")

	if !bytes.Equal(FixedXOR(a, b), c) {
		t.Error("a XOR b != c")
	}
}

func TestRepeatingXOR(t *testing.T) {
	s := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	k := "ICE"
	e := encoding.HexToBytes("0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f")

	if !bytes.Equal(RepeatingXOR([]byte(s), []byte(k)), e) {
		t.Errorf("Did not match")
	}
}
