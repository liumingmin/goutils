package container

type Bitmap []uint32

func (b *Bitmap) Init(maxNum uint32) {
	needLen := b.calcNeedSize(maxNum)
	*b = make([]uint32, needLen)
}

func (b *Bitmap) Exists(item uint32) bool {
	index := b.calcIndex(item)
	if len(*b) < int(index+1) {
		return false
	}

	offsetItem := b.calcPosition(item)
	return ((*b)[index] & (1 << offsetItem)) > 0
}

func (b *Bitmap) Set(item uint32) bool {
	index := b.calcIndex(item)
	if len(*b) < int(index+1) {
		b.Grow(item)
	}

	offsetItem := b.calcPosition(item)

	(*b)[index] = (*b)[index] | (1 << offsetItem)
	return true
}

func (b *Bitmap) Sets(items []uint32) {
	if len(items) == 0 {
		return
	}

	max := uint32(0)
	for _, item := range items {
		max = b.MaxU(max, item)
	}

	index := b.calcIndex(max)
	if len(*b) < int(index+1) {
		b.Grow(max)
	}

	for _, item := range items {
		b.Set(item)
	}
}

func (b *Bitmap) MaxU(x, y uint32) uint32 {
	if x > y {
		return x
	}
	return y
}

func (b *Bitmap) Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (b *Bitmap) Union(b2 *Bitmap) *Bitmap {
	var maxData *Bitmap
	if len(*b) > len(*b2) {
		maxData = b.Clone()
	} else {
		maxData = b2.Clone()
	}

	minLen := b.Min(len(*b), len(*b2))

	for i := 0; i < minLen; i++ {
		(*maxData)[i] = (*b)[i] | (*b2)[i]
	}

	return maxData
}

func (b *Bitmap) Clone() *Bitmap {
	bLen := len(*b)
	copyBitmap := make([]uint32, bLen)
	for i := 0; i < bLen; i++ {
		copyBitmap[i] = (*b)[i]
	}

	bitmap := Bitmap(copyBitmap)
	return &bitmap
}

func (b *Bitmap) Inverse() {
	bLen := len(*b)
	for i := 0; i < bLen; i++ {
		(*b)[i] = ^(*b)[i]
	}
}

func (b *Bitmap) Grow(item uint32) {
	needLen := b.calcNeedSize(item)
	dataLen := uint32(len(*b))
	if dataLen < needLen {
		newData := make([]uint32, needLen-dataLen)
		*b = append(*b, newData...)
	}
}

func (b *Bitmap) calcIndex(i uint32) uint32 {
	return i >> 5
}

func (b *Bitmap) calcPosition(i uint32) uint32 {
	return i & 0x1F
}

func (b *Bitmap) calcNeedSize(i uint32) uint32 {
	return b.calcIndex(i) + 1
}
