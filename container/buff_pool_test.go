package container

import (
	"testing"
	"unsafe"
)

func TestBuffPool(t *testing.T) {
	buf1 := PoolBuffer128K.Get()
	ptr1 := uintptr(unsafe.Pointer(&buf1[0]))
	len1 := len(buf1)

	PoolBuffer128K.Put(buf1)

	buf2 := PoolBuffer128K.Get()
	ptr2 := uintptr(unsafe.Pointer(&buf2[0]))
	len2 := len(buf2)
	PoolBuffer128K.Put(buf2)

	if len1 != 128*1024 {
		t.Error("pool get BUFF_128K len failed")
	}

	if len1 != len2 {
		t.Error("pool get BUFF_128K len failed")
	}

	if ptr1 != ptr2 {
		t.Error("pool get BUFF_128K failed")
	}

	//4M
	buf3 := PoolBuffer4M.Get()
	ptr3 := uintptr(unsafe.Pointer(&buf3[0]))
	len3 := len(buf3)
	PoolBuffer4M.Put(buf3)

	buf4 := PoolBuffer4M.Get()
	ptr4 := uintptr(unsafe.Pointer(&buf4[0]))
	len4 := len(buf4)
	PoolBuffer4M.Put(buf4)

	if len3 != 4*1024*1024 {
		t.Error("pool get BUFF_4M len failed")
	}

	if len3 != len4 {
		t.Error("pool get BUFF_4M len failed")
	}

	if ptr3 != ptr4 {
		t.Error("pool get BUFF_4M failed")
	}
}
