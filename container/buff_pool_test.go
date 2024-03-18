package container

import (
	"testing"
	"unsafe"
)

func TestBuffPool(t *testing.T) {
	buf1 := GetPoolBuff(BUFF_128K)
	ptr1 := uintptr(unsafe.Pointer(&buf1[0]))
	len1 := len(buf1)
	PutPoolBuff(BUFF_128K, buf1)

	buf2 := GetPoolBuff(BUFF_128K)
	ptr2 := uintptr(unsafe.Pointer(&buf2[0]))
	len2 := len(buf2)
	PutPoolBuff(BUFF_128K, buf2)

	if len1 != BUFF_128K {
		t.Error("pool get BUFF_128K len failed")
	}

	if len1 != len2 {
		t.Error("pool get BUFF_128K len failed")
	}

	if ptr1 != ptr2 {
		t.Error("pool get BUFF_128K failed")
	}

	//4M
	buf1 = GetPoolBuff(BUFF_4M)
	ptr1 = uintptr(unsafe.Pointer(&buf1[0]))
	len1 = len(buf1)
	PutPoolBuff(BUFF_4M, buf1)

	buf2 = GetPoolBuff(BUFF_4M)
	ptr2 = uintptr(unsafe.Pointer(&buf2[0]))
	len2 = len(buf2)
	PutPoolBuff(BUFF_4M, buf2)

	if len1 != BUFF_4M {
		t.Error("pool get BUFF_4M len failed")
	}

	if len1 != len2 {
		t.Error("pool get BUFF_4M len failed")
	}

	if ptr1 != ptr2 {
		t.Error("pool get BUFF_4M failed")
	}
}
