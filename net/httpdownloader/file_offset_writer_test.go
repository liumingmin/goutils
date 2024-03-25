package httpdownloader

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

var testTempDirPath = filepath.Join(os.TempDir(), "goutils_httpdownloader")

func TestFileOffsetWriter(t *testing.T) {
	os.MkdirAll(testTempDirPath, 0666)
	savePath := filepath.Join(testTempDirPath, "testFileOffsetWriter.txt")
	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.FailNow()
	}

	var mutex sync.Mutex
	writer := NewOffsetWriter(file, &mutex, 0, 1024)
	writer.Write([]byte("1234567890"))
	writer.ResetOffset()
	writer.Write([]byte("1234567890"))
	file.Close()

	bs, err := os.ReadFile(savePath)
	if err != nil {
		t.FailNow()
	}

	if string(bs) != "1234567890" {
		t.FailNow()
	}
}
