package common

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func TestXorIO(t *testing.T) {
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
		t.Fail()
	}
}

func TestCipherXor(t *testing.T) {
	key := []byte("goutils_is_great")
	writerIndex := uint64(0)

	f, err := os.Open("testxor")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	cf, err := os.OpenFile("testxor.cipher", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer cf.Close()

	w := NewXORWriterWithOffset(cf, key, &writerIndex)

	_, err = io.Copy(w, f)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeCipherXor(t *testing.T) {
	key := []byte("goutils_is_great")
	readerIndex := uint64(0)

	func() {
		f, err := os.Open("testxor.cipher")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		r := NewXORReaderWithOffset(f, key, &readerIndex)

		rf, err := os.OpenFile("testxor.recover", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			t.Fatal(err)
		}
		defer rf.Close()

		_, err = io.Copy(rf, r)
		if err != nil {
			t.Fatal(err)
		}
	}()

	bs1, _ := ioutil.ReadFile("testxor")
	bs2, _ := ioutil.ReadFile("testxor.recover")

	if bytes.Compare(bs1, bs2) != 0 {
		t.Fatal("not eq")
	}
}

func TestXORReaderAt(t *testing.T) {
	key := []byte("goutils_is_great")

	func() {

		f, err := os.Open("testxor.cipher")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		rf, err := os.OpenFile("testxor.recover", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			t.Fatal(err)
		}
		defer rf.Close()

		size := GetFileSize("testxor.cipher")

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

	bs1, _ := ioutil.ReadFile("testxor")
	bs2, _ := ioutil.ReadFile("testxor.recover")

	if bytes.Compare(bs1, bs2) != 0 {
		t.Fatal("not eq")
	}
}
