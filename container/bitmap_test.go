package container

import (
	"fmt"
	"testing"
)

func initTestData() Bitmap {
	bitmap := Bitmap{}
	bitmap.Init(65)
	fmt.Println(bitmap)
	bitmap.Sets([]uint32{4, 16, 66, 64, 32, 128, 122})
	//fmt.Println(bitmap)
	return bitmap
}

func TestBitmapExists(t *testing.T) {
	bitmap := initTestData()
	t.Log(bitmap)

	if !bitmap.Exists(122) {
		t.FailNow()
	}

	if bitmap.Exists(123) {
		t.FailNow()
	}
}

func TestBitmapSet(t *testing.T) {
	bitmap := initTestData()

	if bitmap.Exists(1256) {
		t.FailNow()
	}

	bitmap.Set(1256)

	if !bitmap.Exists(1256) {
		t.FailNow()
	}
}

func TestBitmapUnionOr(t *testing.T) {
	bitmap := initTestData()
	bitmap2 := initTestData()
	bitmap2.Set(256)

	bitmap3 := bitmap.Union(&bitmap2)
	if !bitmap3.Exists(256) {
		t.FailNow()
	}

	bitmap3.Set(562)

	if !bitmap3.Exists(562) {
		t.FailNow()
	}

	if bitmap.Exists(562) {
		t.FailNow()
	}
}

func TestBitmapBitInverse(t *testing.T) {
	bitmap := initTestData()

	if !bitmap.Exists(66) {
		t.FailNow()
	}

	bitmap.Inverse()

	if bitmap.Exists(66) {
		t.FailNow()
	}
}

func BenchmarkBitmap_Exists(b *testing.B) {
	bitmap := Bitmap{}
	bitmap.Init(1000000000)
	for i := 0; i < 1000000000; i++ {
		if i%3 == 0 {
			bitmap.Set(uint32(i)) //3333w
		}
	}

	//b.Log("start compare")
	//b.StartTimer()
	b.N = 1000000000
	for i := 0; i < b.N; i++ {
		bitmap.Exists(uint32(i))
	}
	//b.StopTimer()
}
