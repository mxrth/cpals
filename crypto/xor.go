package crypto

func FixedXOR(as []byte, bs []byte) []byte {
	var x []byte
	if len(as) != len(bs) {
		panic("Input not of equal length")
	}
	for i, a := range as {
		x = append(x, a^bs[i])
	}
	return x
}

func SingleXOR(as []byte, k byte) []byte {
	var x []byte
	for _, a := range as {
		x = append(x, a^k)
	}
	return x
}

func RepeatingXOR(ps []byte, k []byte) []byte {
	if len(k) == 0 {
		panic("Invalid key!")
	}
	c := make([]byte, len(ps))
	for i, p := range ps {
		c[i] = p ^ k[i%len(k)]
	}
	return c
}
