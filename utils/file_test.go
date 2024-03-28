package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var testTempDirPath = filepath.Join(os.TempDir(), "goutils_file")

func TestGetCurrPath(t *testing.T) {
	path := GetCurrPath()
	t.Log(path)
}

func TestFileExist(t *testing.T) {
	runFile := os.Args[0]

	if !FileExist(runFile) {
		t.Error(runFile)
	}
}

func TestFileExt(t *testing.T) {
	if FileExt("aaa.txt") != ".txt" {
		t.Error(FileExt("aaa.txt"))
	}

	if FileExt("aaa.txt.zip") != ".zip" {
		t.Error(FileExt("aaa.txt.zip"))
	}

	if FileExt("aaa.txt.") != "." {
		t.Error(FileExt("aaa.txt."))
	}

	if FileExt("aaa") != "" {
		t.Error(FileExt("aaa"))
	}
}

func TestFileCopy(t *testing.T) {
	runFile := os.Args[0]
	err := FileCopy(runFile, filepath.Join(testTempDirPath, "test_file"))
	if err != nil {
		t.Error()
	}

	err = FileCopy(filepath.Join(testTempDirPath, "test_file"), filepath.Join(testTempDirPath, "test_file"))
	if err != nil {
		t.Error()
	}

	err = FileCopy(filepath.Join(testTempDirPath, "test_file"), ".")
	if err == nil {
		t.Error()
	}
}

func TestIsPathTravOut(t *testing.T) {
	if IsPathTravOut(`C:\a\b`, `C:\a`) {
		t.FailNow()
	}

	if IsPathTravOut(`C:\A\B`, `C:\a`) {
		t.FailNow()
	}

	if IsPathTravOut(`C:\a\b\..`, `C:\a`) {
		t.FailNow()
	}

	if !IsPathTravOut(`C:\a\b\..\..`, `C:\a`) {
		t.FailNow()
	}

	if !IsPathTravOut(`C:\A\B\..\..`, `C:\a`) {
		t.FailNow()
	}

	if !IsPathTravOut(`C:\a\b`, `C:\c`) {
		t.FailNow()
	}
}

func TestUniformPathStyle(t *testing.T) {
	if UniformPathStyle(`C:\a\b`) != `C:/a/b` {
		t.FailNow()
	}

	if UniformPathStyleCase(`C:\A\B`) != `c:/a/b` {
		t.FailNow()
	}

	if !reflect.DeepEqual(UniformPathListStyleCase([]string{`C:\A\B`}), []string{`c:/a/b`}) {
		t.FailNow()
	}
}

func TestIsSameFilePath(t *testing.T) {
	if !IsSameFilePath(`C:\a\b`, `C:/a/b`) {
		t.FailNow()
	}

	if !IsSameFilePath(`C:\a\..\a\b`, `C:/a/b`) {
		t.FailNow()
	}

	if IsSameFilePath(`C:\a\..\a\b\c`, `C:/a/b`) {
		t.FailNow()
	}
}

func TestMain(m *testing.M) {
	os.MkdirAll(testTempDirPath, 0666)

	m.Run()

	os.RemoveAll(testTempDirPath)
}
