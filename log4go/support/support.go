package support

import (
	"bufio"
	"os"
	"time"
)

type support interface {
	StatTimes(filepath string) (atime, ctime, mtime time.Time, err error)
}

var _support support

// GetStatTime returns the times properties corresponding to the given filepath
// NOTE: the atime under windows system may not correct, it maybe the same with
// ctime. (2016-02-26 golang version 1.5.3)
func GetStatTime(filepath string) (atime, ctime, mtime time.Time, err error) {
	return _support.StatTimes(filepath)
}

func GetLines(filepath string) int {
	fd, err := os.Open(filepath)
	defer fd.Close()
	if err != nil {
		return -1
	}
	count := 0

	reader := bufio.NewReader(fd)
	for {
		if _, err := reader.ReadString('\n'); err == nil {
			count++
		} else {
			break
		}
	}

	return count + 1
}

func GetSize(filepath string) int64 {
	if fi, err := os.Stat(filepath); err == nil {
		return fi.Size()
	}
	return -1
}
