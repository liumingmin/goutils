package checksum

import (
	"context"
	"testing"
)

func TestGenerateCheckSumFile(t *testing.T) {
	src := "D:\\gitea_ws\\repair_dir\\dev_test_01\\1.0.0.1\\product"
	checksumName := "nwjs"
	checkSumPath, err := GenerateChecksumFile(context.Background(), src, checksumName)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(checkSumPath)
}

func TestGenerateChecksumMd5File(t *testing.T) {
	src := "D:\\gitea_ws\\repair_dir\\dev_test_01\\1.0.0.1\\product\\nwjs.checksum"
	checksumMd5Path, err := GenerateChecksumMd5File(context.Background(), src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(checksumMd5Path)
}

func TestIsChecksumFileValid(t *testing.T) {
	src := "D:\\gitea_ws\\repair_dir\\dev_test_01\\1.0.0.1\\product\\nwjs.checksum"
	md5Path := "D:\\gitea_ws\\repair_dir\\dev_test_01\\1.0.0.1\\product\\nwjs.checksum.md5"
	valid := IsChecksumFileValid(context.Background(), src, md5Path)
	if !valid {
		t.Error(valid)
		return
	}
	t.Log(valid)
}

func TestCompareChecksumFiles(t *testing.T) {
	src := "D:\\gitea_ws\\repair_dir\\dev_test_01\\1.0.0.1\\product\\nwjs.checksum"
	root := "D:\\gitea_ws\\repair_dir\\dev_test_01\\1.0.0.1\\product"
	err := CompareChecksumFiles(context.Background(), root, src)
	if err != nil {
		t.Error(err)
		return
	}
}
