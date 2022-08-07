package algorithm

import "testing"

func TestCrc16(t *testing.T) {
	t.Log(Crc16([]byte("abcdefg汉字")))
}

func TestCrc16s(t *testing.T) {
	t.Log(Crc16s("abcdefg汉字") == Crc16([]byte("abcdefg汉字")))
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
