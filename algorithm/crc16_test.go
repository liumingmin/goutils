package algorithm

import "testing"

func TestCrc16(t *testing.T) {
	a := Crc16([]byte("abcdefg汉字"))
	b := Crc16([]byte("abcdefg汉字"))
	if a != b {
		t.Error(Crc16([]byte("abcdefg汉字")))
	}
}

func TestCrc16s(t *testing.T) {
	a := Crc16s("abcdefg汉字")
	b := Crc16([]byte("abcdefg汉字"))

	if a != b {
		t.Error(Crc16([]byte("abcdefg汉字")))
	}
}

func BenchmarkCrc16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Crc16([]byte("abcdefg汉字"))
	}
}

func BenchmarkCrc16s(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Crc16s("abcdefg汉字")
	}
}
