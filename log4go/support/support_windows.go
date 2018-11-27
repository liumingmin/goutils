package support

import (
	"os"
	"syscall"
	"time"
)

func init() {
	_support = &supportWin{}
}

type supportWin struct{}

func (t *supportWin) StatTimes(filepath string) (atime, ctime, mtime time.Time, err error) {
	fi, err := os.Lstat(filepath)
	if err != nil {
		return
	}
	data := fi.Sys().(*syscall.Win32FileAttributeData)
	atime = time.Unix(0, data.LastAccessTime.Nanoseconds())
	ctime = time.Unix(0, data.CreationTime.Nanoseconds())
	mtime = fi.ModTime()
	err = nil
	return
}
