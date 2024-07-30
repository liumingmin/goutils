package algorithm

import (
	"bytes"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var testTempDirPath = filepath.Join(os.TempDir(), "goutils_xor")
var testXorOrigFilePath = filepath.Join(testTempDirPath, "iofile.orig")
var testXorCrpytoFilePath = filepath.Join(testTempDirPath, "iofile.xor")
var testXorRecoverFilePath = filepath.Join(testTempDirPath, "iofile.recover")
var testXorKey = []byte("goutils_is_great")

func TestXorIO(t *testing.T) {
	data := []byte("1234567890abcdefhijklmn")

	w := &bytes.Buffer{}

	xw := NewXORWriter(w, testXorKey)
	_, err := io.Copy(xw, bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	cipherBs := w.Bytes()
	xr := NewXORReader(bytes.NewReader(cipherBs), testXorKey)
	rdata, err := io.ReadAll(xr)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, rdata) {
		t.FailNow()
	}
}

func cipherXor(t *testing.T) {
	writerIndex := uint64(0)

	f, err := os.Open(testXorOrigFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	cf, err := os.OpenFile(testXorCrpytoFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer cf.Close()

	w := NewXORWriterWithOffset(cf, testXorKey, &writerIndex)

	_, err = io.Copy(w, f)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeCipherXor(t *testing.T) {
	readerIndex := uint64(0)

	cipherXor(t)

	func() {
		f, err := os.Open(testXorCrpytoFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		r := NewXORReaderWithOffset(f, testXorKey, &readerIndex)

		rf, err := os.OpenFile(testXorRecoverFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			t.Fatal(err)
		}
		defer rf.Close()

		_, err = io.Copy(rf, r)
		if err != nil {
			t.Fatal(err)
		}
	}()

	bs1, _ := os.ReadFile(testXorOrigFilePath)
	bs2, _ := os.ReadFile(testXorRecoverFilePath)

	if !bytes.Equal(bs1, bs2) {
		t.FailNow()
	}
}

func TestXORReaderAt(t *testing.T) {
	cipherXor(t)
	func() {

		f, err := os.Open(testXorCrpytoFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		rf, err := os.OpenFile(testXorRecoverFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			t.Fatal(err)
		}
		defer rf.Close()

		fileInfo, err := os.Stat(testXorCrpytoFilePath)
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
			sr := io.NewSectionReader(NewXORReaderAt(f, testXorKey), offset, rdsize)

			_, err = io.Copy(rf, sr)
			if err != nil {
				t.Fatal(err)
			}

			offset += rdsize
		}
	}()

	bs1, _ := os.ReadFile(testXorOrigFilePath)
	bs2, _ := os.ReadFile(testXorRecoverFilePath)

	if !bytes.Equal(bs1, bs2) {
		t.FailNow()
	}
}

func TestMain(m *testing.M) {
	os.MkdirAll(testTempDirPath, 0666)

	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	file, err := os.OpenFile(testXorOrigFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	block := make([]byte, 32)
	blockCnt := 1024*32 + rd.Intn(16)
	for i := 0; i < blockCnt; i++ {
		for j := 0; j < len(block); j++ {
			block[j] = byte(rd.Intn(256))
		}
		file.Write(block)
	}
	file.Close()

	m.Run()

	os.RemoveAll(testTempDirPath)
}
