package analysis

import "testing"

func TestHammingDist(t *testing.T) {
	if HammingDist([]byte("this is a test"), []byte("wokka wokka!!!")) != 37 {
		t.Errorf("Wrong Haming-Distance")
	}
	s := "this is a testwokka wokka!!!"
	if HammingDist([]byte(s[:14]), []byte(s[14:28])) != 37 {
		t.Errorf("Wat?")
	}
}
