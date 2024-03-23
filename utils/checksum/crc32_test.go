package checksum

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var testTempDirPath = filepath.Join(os.TempDir(), "goutils_checksum")
var testChecksumName = "goutils"
var testChecmsumFileName = "goutils.checksum"

func TestCompareChecksumFiles(t *testing.T) {
	checkSumPath, err := GenerateChecksumFile(context.Background(), testTempDirPath, testChecksumName)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(checkSumPath)

	checksumMd5Path, err := GenerateChecksumMd5File(context.Background(), checkSumPath)
	if err != nil {
		t.Error(err)
		return
	}

	valid := IsChecksumFileValid(context.Background(), checkSumPath, checksumMd5Path)
	if !valid {
		t.Error(valid)
		return
	}

	err = CompareChecksumFiles(context.Background(), testTempDirPath, checkSumPath)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMain(m *testing.M) {
	os.MkdirAll(testTempDirPath, 0666)

	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for l1 := 0; l1 < 5; l1++ {
		for l2 := 0; l2 < 5; l2++ {
			filePath := testTempDirPath + fmt.Sprintf("/%v/%v", l1, l2)
			os.MkdirAll(filepath.Dir(filePath), 0666)
			file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return
			}
			block := make([]byte, 32*1024)
			blockCnt := 32 + rd.Intn(16)
			for i := 0; i < blockCnt; i++ {
				for j := 0; j < len(block); j++ {
					block[j] = byte(rd.Intn(256))
				}
				file.Write(block)
			}
			file.Close()
		}
	}

	m.Run()

	os.RemoveAll(testTempDirPath)
}
