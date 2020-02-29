package rand

import r "crypto/rand"

func Bool() bool {
	b := make([]byte, 1)
	_, err := r.Read(b)
	if err != nil {
		panic(err)
	}
	if b[0]%2 == 1 {
		return true
	}
	return false
}
