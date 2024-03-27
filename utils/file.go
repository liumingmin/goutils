package utils

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetCurrPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}

	path, err := filepath.Abs(file)
	if err != nil {
		return ""
	}

	return filepath.Dir(path)
}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}

func FileExt(filePath string) string {
	idx := strings.LastIndex(filePath, ".")
	ext := ""
	if idx >= 0 {
		ext = strings.ToLower(filePath[idx:])
	}

	return ext
}

func FileCopy(src, dst string) error {
	if IsSameFilePath(src, dst) {
		return nil
	}

	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	// Set destination file attributes
	if err := os.Chmod(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if err := os.Chtimes(dst, srcinfo.ModTime(), srcinfo.ModTime()); err != nil {
		return err
	}

	return nil
}

func IsPathTravOut(path, base string) bool {
	return !strings.HasPrefix(UniformPathStyleCase(filepath.Clean(path)), UniformPathStyleCase(filepath.Clean(base)))
}

func UniformPathStyle(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func UniformPathStyleCase(path string) string {
	return strings.ToLower(UniformPathStyle(path))
}

func UniformPathListStyleCase(pathList []string) []string {
	resultPathList := make([]string, 0, len(pathList))
	for _, path := range pathList {
		resultPathList = append(resultPathList, UniformPathStyleCase(path))
	}
	return resultPathList
}

func IsSameFilePath(src, dst string) bool {
	srcPath, _ := filepath.Abs(src)
	dstPath, _ := filepath.Abs(dst)

	return UniformPathStyleCase(srcPath) == UniformPathStyleCase(dstPath)
}
