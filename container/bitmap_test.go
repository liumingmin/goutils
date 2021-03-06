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
	fmt.Println(bitmap)
	return bitmap
}

func TestBitmapExists(t *testing.T) {
	bitmap := initTestData()
	t.Log(bitmap)

	t.Log(bitmap.Exists(122))
	t.Log(bitmap.Exists(123))

	//data1 := []byte{1, 2, 4, 7}
	//data2 := []byte{0, 1, 5}

}

func TestBitmapSet(t *testing.T) {
	bitmap := initTestData()

	t.Log(bitmap.Exists(1256))

	bitmap.Set(1256)

	t.Log(bitmap.Exists(1256))
}

func TestBitmapUnionOr(t *testing.T) {
	bitmap := initTestData()
	bitmap2 := initTestData()
	bitmap2.Set(256)

	bitmap3 := bitmap.Union(&bitmap2)
	t.Log(bitmap3.Exists(256))

	bitmap3.Set(562)
	t.Log(bitmap3.Exists(562))

	t.Log(bitmap.Exists(562))
}

func TestBitmapBitInverse(t *testing.T) {
	bitmap := initTestData()

	t.Log(bitmap.Exists(66))

	bitmap.Inverse()

	t.Log(bitmap.Exists(66))

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
	b.StopTimer()
}
