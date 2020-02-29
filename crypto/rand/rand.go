//Package rand provides convenience wrappers around crypto/rand and math/rand so you only have to import one package
//it exports top level functions similar to those packages but are all fed from the crypto/rand source
package rand

import (
	cr "crypto/rand"
	"encoding/binary"
	"math/big"
	mr "math/rand"
)

//ok to use from multiple go routines
var source mr.Source64 = src{}

type src struct{}

var rnd = mr.New(source)

func (s src) Seed(seed int64) {}

func (s src) Uint64() (v uint64) {
	err := binary.Read(cr.Reader, binary.BigEndian, &v)
	if err != nil {
		panic(err)
	}
	return
}

func (s src) Int63() int64 {
	//mask highest bit
	return int64(s.Uint64() & ^uint64(1<<63))
}

// Read is a helper function that calls Reader.Read using io.ReadFull.
// On return, n == len(b) if and only if err == nil.
func Read(buf []byte) (int, error) {
	return cr.Read(buf)
}

// BigInt returns a uniform random value in [0, max). It panics if max <= 0.
func BigInt(max *big.Int) (n *big.Int, err error) {
	return cr.Int(cr.Reader, max)
}

// Prime returns a number, p, of the given size, such that p is prime
// with high probability.
// Prime will panic if bits < 2
func Prime(bits int) (p *big.Int, err error) {
	return cr.Prime(cr.Reader, bits)
}

// Float32 returns, as a float32, a random number in [0.0,1.0)
// from the default Source.
func Float32() float32 {
	return rnd.Float32()
}

// Float64 returns, as a float64, a random number in [0.0,1.0)
// from the default Source.
func Float64() float64 {
	return rnd.Float64()
}

// Int returns a non-negative random int from the default Source.
func Int() int {
	return rnd.Int()
}

// Intn returns, as an int, a non-negative random number in [0,n)
// from the default Source.
// It panics if n <= 0.
func Intn(n int) int {
	return rnd.Intn(n)
}

// Uint32 returns a pseudo-random 32-bit value as a uint32
// from the default Source.
func Uint32() uint32 {
	return rnd.Uint32()
}

// Uint64 returns a random 64-bit value as a uint64
// from the default Source.
func Uint64() uint64 {
	return rnd.Uint64()
}

// Bool returns a random boolean
// from the default Source.
func Bool() bool {
	b := make([]byte, 1)
	_, err := cr.Read(b)
	if err != nil {
		panic(err)
	}
	if b[0]%2 == 1 {
		return true
	}
	return false
}
