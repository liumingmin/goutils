package checksum

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
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
	dirMap["jsex"] = []string{"E:\\game\\jsex\\base"}
	for code, dirs := range dirMap {
		t.Log("game: ", code)
		for _, dir := range dirs {
			t.Log("dir: ", dir)
			start := time.Now() // 获取当前时间
			checksumName := fmt.Sprintf("%v", code)
			checksumMd5Path, err := GenerateChecksumFileWithIgnore(context.Background(), dir, checksumName, []string{fmt.Sprintf("%v.checksum", code), "pak", "locales\\pak"})
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

func TestRelPath(t *testing.T) {
	repos := []string{"", "a", "b", "a\\b", "a/c", "a\\b/c", "a/d/c", "d/a", "d/c"}

	for _, repo1 := range repos {
		t.Log(">>>", repo1)
		for _, repo2 := range repos {
			rel, _ := filepath.Rel(repo1, repo2)
			if !strings.Contains(rel, ".") {
				t.Log(repo2, ":", rel)
			}
		}
	}

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
