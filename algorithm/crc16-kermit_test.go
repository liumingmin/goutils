package algorithm

import "testing"

func TestKermit(t *testing.T) {
	a := Kermit([]byte("abcdefg汉字"))
	b := Kermit([]byte("abcdefg汉字"))
	if a != b {
		t.Error(a, b)
	}
}
