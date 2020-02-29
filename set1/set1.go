//Package set1 implements solutions for the cryptopals challenges, Set 1. See https://cryptopals.com for infos.
//TODO: document running tests
package set1

import (
	"bytes"
	"math"
	"math/bits"
)

func bytewiseXOR(a, b []byte) []byte {
	if len(a) != len(b) {
		panic("bytewiseXOR needs buffers of same length")
	}
	var r = make([]byte, len(a))
	for i := range a {
		r[i] = a[i] ^ b[i]
	}
	return r
}

func singleByteXOR(a []byte, k byte) []byte {
	r := make([]byte, len(a))

	for i := range a {
		r[i] = a[i] ^ k
	}

	return r
}

func breakSingleByteXOR(c []byte) (solution byte, maxScore float64) {
	for k := 0; k < 256; k++ {
		candidate := singleByteXOR(c, byte(k))
		if scoreText(candidate) > maxScore {
			maxScore = scoreText(candidate)
			solution = byte(k)
		}
	}
	return
}

//taken from http://pi.math.cornell.edu/~mec/2003-2004/cryptography/subs/frequencies.html
var frequencies = map[byte]float64{
	'a': 0.08167,
	'b': 0.01492,
	'c': 0.02782,
	'd': 0.04253,
	'e': 0.1270,
	'f': 0.02228,
	'g': 0.02015,
	'h': 0.06094,
	'i': 0.06966,
	'j': 0.00153,
	'k': 0.00772,
	'l': 0.04025,
	'm': 0.02406,
	'n': 0.06749,
	'o': 0.07507,
	'p': 0.01929,
	'q': 0.00095,
	'r': 0.05987,
	's': 0.06327,
	't': 0.09056,
	'u': 0.02758,
	'v': 0.00978,
	'w': 0.02360,
	'x': 0.00150,
	'y': 0.01974,
	'z': 0.00074,
	' ': 0.02, //ad-hoc probability of spaces roughly equals P(w) or P(y)
}

//implement Bhattacharyya coefficient (https://en.wikipedia.org/wiki/Bhattacharyya_distance#Bhattacharyya_coefficient)
//higher is better
func scoreText(text []byte) float64 {
	//c = \sum \sqrt(p_i * p_q)
	text = bytes.ToLower(text)
	c := 0.0
	for b, p := range frequencies {
		q := float64(bytes.Count(text, []byte{b})) / float64(len(text))
		c += math.Sqrt(p * q)
	}
	return c
}

func multibyteXOR(plain, key []byte) []byte {
	cipher := make([]byte, len(plain))
	for i, b := range plain {
		cipher[i] = b ^ key[i%len(key)]
	}
	return cipher
}

func breakMultibyteXOR(cipher []byte) (key []byte) {
	minDist := 100.0
	keysizeGuess := 2
	for keysize := 2; keysize <= 40; keysize++ {
		blocks := make([][]byte, 4)
		for i := 0; i < 4; i++ {
			blocks[i] = cipher[i*keysize : (i+1)*keysize]
		}
		dist := 0.0
		for i := 0; i < 4; i++ {
			for j := i + 1; j < 4; j++ {
				dist += float64(hammingDist(blocks[i], blocks[j])) / float64(keysize)
			}
		}
		if dist < minDist {
			keysizeGuess = keysize
			minDist = dist
		}
	}
	//fmt.Println(keysizeGuess)

	//byte i goes into block (i % keysizeGuess)
	blocks := make([][]byte, keysizeGuess)

	for i, b := range cipher {
		blocks[i%keysizeGuess] = append(blocks[i%keysizeGuess], b)
	}

	for _, block := range blocks {
		k, _ := breakSingleByteXOR(block)
		key = append(key, k)
	}

	return
}

func hammingDist(as, bs []byte) int {
	if len(as) != len(bs) {
		panic("hammingDist expects slices of equal length")
	}
	dist := 0
	for i, a := range as {
		dist += bits.OnesCount8(a ^ bs[i])
	}
	return dist
}
