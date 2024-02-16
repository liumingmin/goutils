

<!-- toc -->

- [algorithm 算法模块](#algorithm-%E7%AE%97%E6%B3%95%E6%A8%A1%E5%9D%97)
  * [circ2buffer_test.go](#circ2buffer_testgo)
    + [TestCreateC2Buffer](#testcreatec2buffer)
    + [TestWriteBlock](#testwriteblock)
    + [TestWritingUnderCapacityGivesEmptyEvicted](#testwritingundercapacitygivesemptyevicted)
    + [TestWritingMultipleBytesWhenBufferIsNotFull](#testwritingmultiplebyteswhenbufferisnotfull)
    + [TestEvictedRegession1](#testevictedregession1)
    + [TestGetBlock](#testgetblock)
    + [TestWriteTwoBlocksGet](#testwritetwoblocksget)
    + [TestWriteSingleByteGetSingleByte](#testwritesinglebytegetsinglebyte)
    + [TestWriteTwoBlocksGetEvicted](#testwritetwoblocksgetevicted)
    + [TestWriteSingleByteReturnsSingleEvictedByte](#testwritesinglebytereturnssingleevictedbyte)
    + [TestTruncatingAfterWriting](#testtruncatingafterwriting)
    + [TestWritingAfterTruncating](#testwritingaftertruncating)
  * [crc16_test.go crc16算法](#crc16_testgo-crc16%E7%AE%97%E6%B3%95)
    + [TestCrc16](#testcrc16)
    + [TestCrc16s](#testcrc16s)
  * [descartes_test.go 笛卡尔组合](#descartes_testgo-%E7%AC%9B%E5%8D%A1%E5%B0%94%E7%BB%84%E5%90%88)
    + [TestDescartes](#testdescartes)
  * [xor_io_test.go](#xor_io_testgo)
    + [TestXorIO](#testxorio)
    + [TestCipherXor](#testcipherxor)
    + [TestDeCipherXor](#testdecipherxor)
    + [TestXORReaderAt](#testxorreaderat)

<!-- tocstop -->

# algorithm 算法模块
## circ2buffer_test.go
### TestCreateC2Buffer
```go

MakeC2Buffer(BLOCK_SIZE)
```
### TestWriteBlock
```go

b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock)
```
### TestWritingUnderCapacityGivesEmptyEvicted
```go

b := MakeC2Buffer(2)
b.Write([]byte{1, 2})

if len(b.Evicted()) != 0 {
	t.Fatal("Evicted should have been empty:", b.Evicted())
}
```
### TestWritingMultipleBytesWhenBufferIsNotFull
```go

b := MakeC2Buffer(3)
b.Write([]byte{1, 2})
b.Write([]byte{3, 4})

ev := b.Evicted()

if len(ev) != 1 || ev[0] != 1 {
	t.Fatal("Evicted should have been [1,]:", ev)
}
```
### TestEvictedRegession1
```go

b := MakeC2Buffer(4)

b.Write([]byte{7, 6})
b.Write([]byte{5, 1, 2})
b.Write([]byte{3, 4})

ev := b.Evicted()
if len(ev) != 2 || ev[0] != 6 || ev[1] != 5 {
	t.Fatalf("Unexpected evicted [6,5]: %v", ev)
}
```
### TestGetBlock
```go

b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock)

block := b.GetBlock()

if len(block) != BLOCK_SIZE {
	t.Fatal("Wrong block size returned")
}

for i, by := range block {
	if byte(i) != by {
		t.Errorf("byte %v does not match", i)
	}
}
```
### TestWriteTwoBlocksGet
```go

b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock)
b.Write(incrementBlock2)

if bytes.Compare(b.GetBlock(), incrementBlock2) != 0 {
	t.Errorf("Get block did not return the right value: %s", b.GetBlock())
}
```
### TestWriteSingleByteGetSingleByte
```go

b := MakeC2Buffer(BLOCK_SIZE)
singleByte := []byte{0}
b.Write(singleByte)

if bytes.Compare(b.GetBlock(), singleByte) != 0 {
	t.Errorf("Get block did not return the right value: %s", b.GetBlock())
}
```
### TestWriteTwoBlocksGetEvicted
```go

b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock)
b.Write(incrementBlock2)

if bytes.Compare(b.Evicted(), incrementBlock) != 0 {
	t.Errorf("Evicted did not return the right value: %s", b.Evicted())
}
```
### TestWriteSingleByteReturnsSingleEvictedByte
```go

b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock2)
singleByte := []byte{0}

b.Write(singleByte)
e := b.Evicted()

if len(e) != 1 {
	t.Fatalf("Evicted length is not correct: %s", e)
}

if e[0] != byte(10) {
	t.Errorf("Evicted content is not correct: %s", e)
}
```
### TestTruncatingAfterWriting
```go

b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock)

evicted := b.Truncate(2)

if len(evicted) != 2 {
	t.Fatalf("Truncate did not return expected evicted length: %v", evicted)
}

if evicted[0] != 0 || evicted[1] != 1 {
	t.Errorf("Unexpected content in evicted: %v", evicted)
}
```
### TestWritingAfterTruncating
```go

// test that after we truncate some content, the next operations
// on the buffer give us the expected results
b := MakeC2Buffer(BLOCK_SIZE)
b.Write(incrementBlock)
b.Truncate(4)

b.Write([]byte{34, 46})

block := b.GetBlock()

if len(block) != BLOCK_SIZE-2 {
	t.Fatalf(
		"Unexpected block length after truncation: %v (%v)",
		block,
		len(block),
	)
}

if bytes.Compare(block, []byte{4, 5, 6, 7, 8, 9, 34, 46}) != 0 {
	t.Errorf(
		"Unexpected block content after truncation: %v (%v)",
		block,
		len(block))
}
```
## crc16_test.go crc16算法
### TestCrc16
```go

t.Log(Crc16([]byte("abcdefg汉字")))
```
### TestCrc16s
```go

t.Log(Crc16s("abcdefg汉字") == Crc16([]byte("abcdefg汉字")))
```
## descartes_test.go 笛卡尔组合
### TestDescartes
```go

result := DescartesCombine([][]string{{"A", "B"}, {"1", "2", "3"}, {"a", "b", "c", "d"}})
for _, item := range result {
	t.Log(item)
}
```
## xor_io_test.go
### TestXorIO
```go

key := []byte("goutils_is_great")
data := []byte("1234567890abcdefhijklmn")

w := &bytes.Buffer{}

xw := NewXORWriter(w, key)
_, err := io.Copy(xw, bytes.NewReader(data))
if err != nil {
	t.Fatal(err)
}

cipherBs := w.Bytes()
xr := NewXORReader(bytes.NewReader(cipherBs), key)
rdata, err := ioutil.ReadAll(xr)
if err != nil {
	t.Fatal(err)
}

if bytes.Compare(data, rdata) != 0 {
	t.FailNow()
}
```
### TestCipherXor
```go

key := []byte("goutils_is_great")
writerIndex := uint64(0)

f, err := os.Open("../.tools/protoc-3.19.4-win64.zip")
if err != nil {
	t.Fatal(err)
}
defer f.Close()

cf, err := os.OpenFile("../.tools/protoc-3.19.4-win64.zip.xor", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
if err != nil {
	t.Fatal(err)
}
defer cf.Close()

w := NewXORWriterWithOffset(cf, key, &writerIndex)

_, err = io.Copy(w, f)
if err != nil {
	t.Fatal(err)
}
```
### TestDeCipherXor
```go

key := []byte("goutils_is_great")
readerIndex := uint64(0)

func() {
	f, err := os.Open("../.tools/protoc-3.19.4-win64.zip.xor")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	r := NewXORReaderWithOffset(f, key, &readerIndex)

	rf, err := os.OpenFile("../.tools/protoc-3.19.4-win64.zip.recover", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer rf.Close()

	_, err = io.Copy(rf, r)
	if err != nil {
		t.Fatal(err)
	}
}()

bs1, _ := ioutil.ReadFile("../.tools/protoc-3.19.4-win64.zip")
bs2, _ := ioutil.ReadFile("../.tools/protoc-3.19.4-win64.zip.recover")

if bytes.Compare(bs1, bs2) != 0 {
	t.FailNow()
}
```
### TestXORReaderAt
```go

key := []byte("goutils_is_great")

func() {

	f, err := os.Open("../.tools/protoc-3.19.4-win64.zip.xor")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	rf, err := os.OpenFile("../.tools/protoc-3.19.4-win64.zip.recover", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer rf.Close()

	fileInfo, err := os.Stat("../.tools/protoc-3.19.4-win64.zip.xor")
	if err != nil {
		return
	}
	size := fileInfo.Size()

	for offset := int64(0); offset < size; {
		rdsize := int64(rand.Intn(int(size) / 2))
		if offset+rdsize > size {
			rdsize = size - offset
		}

		t.Logf("read section offset: %v, size: %v\n", offset, rdsize)
		sr := io.NewSectionReader(NewXORReaderAt(f, key), offset, rdsize)

		_, err = io.Copy(rf, sr)
		if err != nil {
			t.Fatal(err)
		}

		offset += rdsize
	}
}()

bs1, _ := ioutil.ReadFile("../.tools/protoc-3.19.4-win64.zip")
bs2, _ := ioutil.ReadFile("../.tools/protoc-3.19.4-win64.zip.recover")

if bytes.Compare(bs1, bs2) != 0 {
	t.FailNow()
}
```
