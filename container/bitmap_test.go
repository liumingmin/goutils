package container

import (
	"fmt"
	"testing"
)

func initTestData() Bitmap {
	bitmap := Bitmap{}
	bitmap.Init(65)
	fmt.Println(bitmap)
	bitmap.SetAll([]uint32{4, 16, 66, 64, 32, 128, 122})

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

	bitmap3 := bitmap.UnionOr(&bitmap2)
	t.Log(bitmap3.Exists(256))

	bitmap3.Set(562)
	t.Log(bitmap3.Exists(562))

	t.Log(bitmap.Exists(562))
}

func TestBitmapBitInverse(t *testing.T) {
	bitmap := initTestData()

	t.Log(bitmap.Exists(66))

	bitmap.BitInverse()

	t.Log(bitmap.Exists(66))

}
