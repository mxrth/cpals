package analysis

func HammingDist(a []byte, b []byte) int {
	if len(a) != len(b) {
		panic("Hamming Distance only defined for slices of equal length")
	}
	dist := 0
	for i, _ := range a {
		dist += byteDist(a[i], b[i])
	}
	return dist
}

func byteDist(a byte, b byte) int {
	dist := 0
	for c := a ^ b; c != 0; c &= c - 1 {
		dist++
	}
	return dist
}

func avgHamming(c []byte, len int) float64 {
	n := 4
	chunks := make([][]byte, n)
	for i := 0; i < n; i++ {
		chunks[i] = c[i*len : (i+1)*len]
	}
	dist := 0
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			dist += HammingDist(chunks[i], chunks[j])
		}
	}
	numDists := float64((n * (n - 1)) / 2)
	return float64(dist) / numDists
}
