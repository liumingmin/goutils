package algorithm

import "testing"

func TestCrc16(t *testing.T) {
	t.Log(Crc16([]byte("abcdefg")))
}
