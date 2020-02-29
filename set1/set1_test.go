package set1

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"testing"

	"github.com/drak3/cpals/encoding"
	"github.com/drak3/cpals/file"
)

func TestChallenge001(t *testing.T) {
	var haveHex = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	var wantBase = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	gotBase, err := encoding.HexToBase64(haveHex)
	if err != nil {
		t.Error(err)
	}
	if gotBase != wantBase {
		t.Errorf("Wrong hex (%s) from base64", gotBase)
	}
}

func TestChallenge002(t *testing.T) {
	a, err := hex.DecodeString("1c0111001f010100061a024b53535009181c")
	b, err := hex.DecodeString("686974207468652062756c6c277320657965")
	want, err := hex.DecodeString("746865206b696420646f6e277420706c6179")
	if err != nil {
		t.Errorf("wut?")
	}
	c := bytewiseXOR(a, b)
	if !bytes.Equal(want, c) {
		t.Errorf("FixedXOR failed")
	}
}

func TestChallenge003(t *testing.T) {
	c, _ := hex.DecodeString("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")
	k, s := breakSingleByteXOR(c)
	t.Logf("Found key %d (%c) with score %f\n", k, k, s)
	t.Logf("decrypts to: \"%s\"\n", singleByteXOR(c, 'X'))
	if k != byte('X') {
		t.Errorf("Got wrong key %d (%c)\n", k, k)
	}
}

func TestChallenge004(t *testing.T) {
	lines := file.ReadLines("testdata/4.txt")
	bestLine, key, maxScore := 0, byte(0), 0.0
	for i, line := range lines {
		l, err := hex.DecodeString(line)
		if err != nil {
			t.Fail()
		}
		k, score := breakSingleByteXOR(l)
		if score > maxScore {
			//t.Log(string(singleByteXOR(line, k)))
			//t.Log(i)
			bestLine, key, maxScore = i, k, score
		}
	}
	t.Log(maxScore)
	t.Logf("Best line is line no. %d, key=%d (%c)\n", bestLine, key, key)
	l, _ := hex.DecodeString(lines[bestLine])
	t.Logf("It decrypts to %s\n", string(singleByteXOR(l, key)))
	if bestLine != 170 {
		t.Errorf("Got wrong line %d, expected 170", bestLine)
	}
	if key != byte('5') {
		t.Errorf("Got wrong key %d, expected '5'", key)
	}
}

func TestChallenge005(t *testing.T) {
	plain := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	c := multibyteXOR([]byte(plain), []byte("ICE"))
	want, _ := hex.DecodeString("0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f")
	if !bytes.Equal(c, want) {
		t.Errorf("Failed multibyteXOR encryption")
	}
}

func TestChallenge006(t *testing.T) {
	dist := hammingDist([]byte("this is a test"), []byte("wokka wokka!!!"))
	if dist != 37 {
		t.Errorf("Expected Hamming distance 37, got %d", dist)
	}
	base, _ := ioutil.ReadFile("testdata/6.txt")
	cipher, _ := base64.StdEncoding.DecodeString(string(base))
	t.Log(len(cipher))
	key := breakMultibyteXOR(cipher)

	t.Logf("Found key '%s' with length %d", key, len(key))
	//plain := multibyteXOR(cipher, key)
	//t.Logf("Plaintext:\n %s\n", string(plain))

}

func TestChallenge007(t *testing.T) {
	bs, _ := ioutil.ReadFile("testdata/7.txt")
	ciph, _ := base64.StdEncoding.DecodeString(string(bs))
	key := []byte("YELLOW SUBMARINE")

	c, _ := aes.NewCipher(key)
	crypter := newECBDecrypter(c)
	plain := make([]byte, len(ciph))
	crypter.CryptBlocks(plain, ciph)
	//t.Log(string(plain))
}

func TestECB(t *testing.T) {
	r := make([]byte, 100*aes.BlockSize)
	orig := make([]byte, len(r))
	rand.Read(r)
	copy(orig, r)
	key := []byte("YELLOW SUBMARINE")
	c, _ := aes.NewCipher(key)
	dec := newECBDecrypter(c)
	enc := newECBEncrypter(c)
	enc.CryptBlocks(r, r)

	if bytes.Equal(r, orig) {
		t.Error("enc = id (probability of fluke is really low)")
	}

	dec.CryptBlocks(r, r)
	if !bytes.Equal(r, orig) {
		t.Error("dec enc != id")
	}

}

func TestChallenge008(t *testing.T) {
	lines := readLines("testdata/8.txt")

loop:
	for i, l := range lines {
		//if we model the ciphertexts as random bytes, the probability that there are two blocks that are exactly the same is
		// roughly (numBlocks^2)/2^(8*aes.BlockSize) the numerator comes out as 2^128 so whenever we find two blocks that are the same we can be virtually certain
		// that ECB-Mode was used
		// detection is done in the naive O(n^2) way: iterate over the blocks twice and compare
		c, _ := hex.DecodeString(l)

		for j := 0; j < len(c)/aes.BlockSize; j++ {
			for k := j + 1; k < len(c)/aes.BlockSize; k++ {
				if bytes.Equal(c[j*aes.BlockSize:(j+1)*aes.BlockSize], c[k*aes.BlockSize:(k+1)*aes.BlockSize]) {
					t.Logf("Line %d is probably encrypted in ECB mode.", i)
					break loop
				}
			}
		}
	}
}
