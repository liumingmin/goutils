package checksum

import (
	"context"
	"fmt"
	"testing"
	"time"
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

func TestGenerateChecksumFileWithIgnore(t *testing.T) {
	dirMap := make(map[string][]string)
	dirMap["fullClient"] = []string{"E:\\game\\dev_test_01", "E:\\game\\dev_test_01\\fullClient"}
	for code, dirs := range dirMap {
		t.Log("game: ", code)
		for _, dir := range dirs {
			t.Log("dir: ", dir)
			start := time.Now() // 获取当前时间
			checksumName := fmt.Sprintf("%v-62e204c376d4be7b1458d077", code)
			checksumMd5Path, err := GenerateChecksumFileWithIgnore(context.Background(), dir, checksumName, []string{fmt.Sprintf("%v-62e204c376d4be7b1458d077.checksum", code)})
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(checksumMd5Path)
			elapsed := time.Since(start)
			t.Log("time：", elapsed)
		}
	}
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
